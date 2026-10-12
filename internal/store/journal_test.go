package store

import (
	"os"
	"path/filepath"
	"sync"
	"testing"
)

func TestCorruptJournalIsNeverOverwritten(t *testing.T) {
	s := &Store{State: t.TempDir()}
	b := []byte("{broken journal")
	if err := os.WriteFile(s.changesPath(), b, 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := s.BeginChange(Change{Item: "dock-position"}); err == nil {
		t.Fatal("corruption must block changes")
	}
	if err := s.MarkReverted("missing"); err == nil {
		t.Fatal("corruption must block reversal marking")
	}
	got, _ := os.ReadFile(s.changesPath())
	if string(got) != string(b) {
		t.Fatal("corrupt journal was overwritten")
	}
}

func TestPendingJournalSurvivesReopen(t *testing.T) {
	s := &Store{State: t.TempDir()}
	id, err := s.BeginChange(Change{Before: "LEFT", After: "BOTTOM", Item: "dock-position"})
	if err != nil {
		t.Fatal(err)
	}
	other := &Store{State: s.State}
	got, err := other.ReadChanges()
	if err != nil || len(got) != 1 || got[0].ID != id || got[0].Status != "pending" {
		t.Fatalf("pending intent lost: %+v, %v", got, err)
	}
	if err := other.FinishChange(id, "RIGHT"); err != nil {
		t.Fatal(err)
	}
	got, _ = s.ReadChanges()
	if got[0].Status != "applied" || got[0].After != "RIGHT" {
		t.Fatal("completion was not saved")
	}
}

func TestConcurrentPrefsRemainValid(t *testing.T) {
	s := &Store{Config: t.TempDir()}
	var wg sync.WaitGroup
	for i := 0; i < 24; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := s.SavePrefs(Prefs{Theme: "dark"}); err != nil {
				t.Error(err)
			}
		}()
	}
	wg.Wait()
	var p Prefs
	if err := readJSON(filepath.Join(s.Config, "prefs.json"), &p); err != nil || p.Theme != "dark" {
		t.Fatalf("invalid preferences: %v", err)
	}
}
