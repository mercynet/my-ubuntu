// Package store persists profiles, preferences, run history, change journal
// and logs in XDG directories. Nothing here ever contains credentials.
package store

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"meuubuntu/internal/catalog"
	"meuubuntu/internal/safefile"
)

type Store struct {
	Config, State, Data string
	mu                  sync.Mutex
}

func xdg(env, def string) string {
	if v := os.Getenv(env); v != "" && filepath.IsAbs(v) {
		return v
	}
	h, _ := os.UserHomeDir()
	return filepath.Join(h, def)
}

// New opens the store; name is "meu-ubuntu" (or "meu-ubuntu-sim" in simulation).
func New(name string) (*Store, error) {
	s := &Store{
		Config: filepath.Join(xdg("XDG_CONFIG_HOME", ".config"), name),
		State:  filepath.Join(xdg("XDG_STATE_HOME", ".local/state"), name),
		Data:   filepath.Join(xdg("XDG_DATA_HOME", ".local/share"), name),
	}
	for _, d := range []string{filepath.Join(s.Config, "profiles"), filepath.Join(s.Config, "profiles", "imported"), filepath.Join(s.State, "runs"), filepath.Join(s.State, "logs"), s.Data} {
		if err := os.MkdirAll(d, 0o700); err != nil {
			return nil, err
		}
	}
	return s, nil
}

func writeJSON(path string, v any) error {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	return safefile.Write(path, b, 0o600)
}

func readJSON(path string, v any) error {
	b, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, v)
}

// ---- Preferences ----

type Prefs struct {
	Theme             string                       `json:"theme"` // system|light|dark
	Selections        map[string]catalog.Selection `json:"selections"`
	ShowAdv           bool                         `json:"showAdvanced"`
	ShowAllAreas      bool                         `json:"showAllAreas"`
	Reverts           []string                     `json:"reverts"`
	Origins           map[string][]string          `json:"origins"`
	Onboarded         bool                         `json:"onboarded"`
	ManualProgress    map[string][]string          `json:"manualProgress,omitempty"`
	CatalogProfile    string                       `json:"catalogProfile,omitempty"`
	ExcludedInstalled []string                     `json:"excludedInstalled,omitempty"`
}

func (s *Store) LoadPrefs() Prefs {
	p := Prefs{Theme: "system", Selections: map[string]catalog.Selection{}}
	readJSON(filepath.Join(s.Config, "prefs.json"), &p)
	if p.Selections == nil {
		p.Selections = map[string]catalog.Selection{}
	}
	return p
}

func (s *Store) SavePrefs(p Prefs) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return writeJSON(filepath.Join(s.Config, "prefs.json"), p)
}

// ---- Change journal (per-change reversal) ----

type Change struct {
	ID       string `json:"id"`
	Time     string `json:"time"`
	RunID    string `json:"runId"`
	Item     string `json:"item"`
	Title    string `json:"title"`
	Kind     string `json:"kind"` // gsettings | git | template | power | admin
	Schema   string `json:"schema,omitempty"`
	Key      string `json:"key,omitempty"`
	Before   string `json:"before"`
	After    string `json:"after"`
	HadValue bool   `json:"hadValue"`
	Reverted string `json:"reverted,omitempty"`
	Status   string `json:"status,omitempty"` // pending | applied; empty for old journals
}

func (s *Store) changesPath() string { return filepath.Join(s.State, "changes.json") }

func (s *Store) Changes() []Change {
	c, _ := s.ReadChanges()
	return c
}

func (s *Store) ReadChanges() ([]Change, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.readChanges()
}

func (s *Store) readChanges() ([]Change, error) {
	var c []Change
	err := readJSON(s.changesPath(), &c)
	if os.IsNotExist(err) {
		err = nil
	}
	return c, err
}

func (s *Store) AddChange(c Change) error {
	c.Status = "applied"
	_, err := s.addChange(c)
	return err
}

// BeginChange durably saves intent BEFORE touching the system. A pending
// record survives a failed command, crash or failed completion write.
func (s *Store) BeginChange(c Change) (string, error) {
	c.Status = "pending"
	return s.addChange(c)
}

func (s *Store) addChange(c Change) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	all, err := s.readChanges()
	if err != nil {
		return "", fmt.Errorf("histórico de alterações ilegível: %w", err)
	}
	c.ID = "c" + NewID()
	c.Time = time.Now().Format(time.RFC3339)
	all = append(all, c)
	return c.ID, writeJSON(s.changesPath(), all)
}

func (s *Store) FinishChange(id, after string) error {
	return s.updateChange(id, func(c *Change) { c.Status, c.After = "applied", after })
}

func (s *Store) MarkReverted(id string) error {
	return s.updateChange(id, func(c *Change) { c.Reverted = time.Now().Format(time.RFC3339) })
}

func (s *Store) FinishReset(id, after string) error {
	return s.updateChange(id, func(c *Change) { c.Status, c.After, c.Reverted = "applied", after, time.Now().Format(time.RFC3339) })
}

func (s *Store) updateChange(id string, update func(*Change)) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	all, err := s.readChanges()
	if err != nil {
		return err
	}
	for i := range all {
		if all[i].ID == id {
			update(&all[i])
			return writeJSON(s.changesPath(), all)
		}
	}
	return fmt.Errorf("alteração %s não encontrada", id)
}

func NewID() string {
	return strings.ReplaceAll(time.Now().UTC().Format("20060102-150405.000000000"), ".", "-") + "-" + rand.Text()
}

// ---- Runs (history + resume) ----

type StepRecord struct {
	ID      string   `json:"id"`
	Title   string   `json:"title"`
	Item    string   `json:"item"`
	Status  string   `json:"status"` // pending running done failed skipped intervention cancelled
	Error   string   `json:"error,omitempty"`
	Note    string   `json:"note,omitempty"`
	Op      string   `json:"op,omitempty"`
	Problem *Problem `json:"problem,omitempty"`
}

type Problem struct {
	Code    string `json:"code"`
	Summary string `json:"summary"`
	Advice  string `json:"advice"`
}

type Run struct {
	ID         string                       `json:"id"`
	Started    string                       `json:"started"`
	Finished   string                       `json:"finished,omitempty"`
	Status     string                       `json:"status"` // running done failed intervention cancelled interrupted
	Selections map[string]catalog.Selection `json:"selections"`
	Reverts    []string                     `json:"reverts,omitempty"`
	Steps      []StepRecord                 `json:"steps"`
	LogFile    string                       `json:"logFile"`
	Simulated  bool                         `json:"simulated"`
}

func (s *Store) SaveRun(r Run) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return writeJSON(filepath.Join(s.State, "runs", r.ID+".json"), r)
}

func (s *Store) Runs() []Run {
	return s.runs(30)
}

// RecoveryRuns includes older unresolved verification records beyond the UI's
// recent-history window.
func (s *Store) RecoveryRuns() []Run { return s.runs(0) }

func (s *Store) runs(limit int) []Run {
	files, _ := filepath.Glob(filepath.Join(s.State, "runs", "*.json"))
	sort.Sort(sort.Reverse(sort.StringSlice(files)))
	var out []Run
	for i, f := range files {
		if limit > 0 && i >= limit {
			break
		}
		var r Run
		if readJSON(f, &r) == nil {
			out = append(out, r)
		}
	}
	return out
}

func (s *Store) LogPath(runID string) string { return filepath.Join(s.State, "logs", runID+".log") }
