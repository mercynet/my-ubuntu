// Typed wrapper over the Wails bindings (window.go.main.App) and runtime events.
// The UI never sends commands — only item IDs and validated option values.

export type Sel = { options?: Record<string, string[]> };
export type Selections = Record<string, Sel>;

export interface OptionValue { value: string; label: string }
export interface Option { key: string; label: string; multi?: boolean; values: OptionValue[]; default: string[]; help?: string }
export interface Conflict { id: string; hard: boolean; reason: string }
export interface Item {
  id: string; name: string; area: string; category: string; summary: string; why: string; notes?: string;
  source: string; method: string; packages?: string[]; repo?: string; flatpak?: string; extension?: string;
  providedBy?: string[]; requires?: string[]; conflicts?: Conflict[]; group?: string; addonsFor?: string[]; options?: Option[];
  advanced?: boolean; service?: boolean; relogin?: boolean; reboot?: boolean; verified?: string;
  manual?: { steps: string[]; url?: string }; compat?: { arch?: string[]; notes?: string };
  settings?: { schema: string; key: string; value?: string; option?: string; map?: Record<string, string> }[];
  git?: { key: string; value?: string; option?: string }[];
  script?: { url: string; allowHosts: string[] };
  sensitive?: string; confirm?: string; adminOp?: string;
  npm?: { package: string; version: string; binary: string; ignoreScripts?: boolean };
}
export interface SettingState { schema: string; key: string; current: string; proposed: string; missing: boolean }
export interface ItemState { status: string; detail: string; settings?: SettingState[] }
export interface Repo { id: string; name: string; fingerprints: string[]; url: string; docs: string }
export interface System { id: string; release: string; codename: string; pretty: string; arch: string; gnome: string; gnomeMajor: number; session: string; simulated: boolean }
export interface Profile {
  schema?: string; version?: number; id?: string; name: string; description?: string; audience?: string; author?: string;
  kind?: "usage" | "complement" | ""; requires?: { minRelease?: string; arch?: string[] }; created?: string; updated?: string;
  selections: Selections; source?: "builtin" | "mine" | "imported"; builtin?: boolean;
}
export type Origins = Record<string, string[]>;
export interface Prefs { theme: string; selections: Selections; showAdvanced: boolean; showAllAreas?: boolean; reverts?: string[] | null; origins?: Origins | null; onboarded?: boolean; manualProgress?: Record<string, string[]> | null; catalogProfile?: string; excludedInstalled?: string[] }
export interface ReviewItem { id: string; name: string; area: string; status: string; detail: string; sensitive?: string }
export interface ProfileReview { profile: Profile; warnings: string[] | null; items: ReviewItem[] | null; requirement: string[] | null; existsMine: boolean; existsImported: boolean; path?: string }
export interface MergeConflict { kind: "options" | "group" | "exclusive"; item: string; name: string; current: Sel; currentOrigin: string[] | null; incomingItem: string; incomingName: string; incoming: Sel; incomingLabel: string; reason: string }
export interface MergeResult { selections: Selections; origins: Origins; added: string[] | null; conflicts: MergeConflict[] | null }
export interface Group { id: string; name: string; area: string; summary: string; why: string }
export interface Problem { code: string; summary: string; advice: string }
export interface StepRecord { id: string; title: string; item: string; status: string; error?: string; note?: string; op?: string; problem?: Problem }
export interface Run { id: string; started: string; finished?: string; status: string; selections: Selections; reverts?: string[] | null; steps: StepRecord[]; logFile: string; simulated: boolean }
export interface InitData { version: string; system: System; prefs: Prefs; profiles: Profile[]; interrupted: Run | null; repos: Repo[]; groups: Group[] | null; items: Item[]; paths: Record<string, string>; pendingImport: ProfileReview | null; pendingError: string }
export interface AboutInfo { version: string; author: string; goVersion: string; wailsVersion: string; attributions: { name: string; license: string; url: string; text: string }[] }
export interface PlanStep { id: string; item: string; title: string; kind: string; op: string; detail: string[]; changes?: SettingState[]; category: string; options?: Record<string, string[]>; confirm?: string; expect?: string }
export interface Plan {
	 reviewToken: string;
  steps: PlanStep[]; conflicts: { a: string; b: string; hard: boolean; reason: string }[]; blocked: boolean; added: string[] | null;
  skipped: { item: string; name: string; reason: string; status: string }[] | null; manual: { item: string; name: string; steps: string[]; url: string }[] | null;
  relogin: boolean; reboot: boolean; needsRoot: boolean; warnings: string[] | null; selection: Selections;
  reverts: string[] | null; kept: { item: string; name: string; reason: string; status: string }[] | null;
}
export interface Change { id: string; time: string; runId: string; item: string; title: string; kind: string; schema?: string; key?: string; before: string; after: string; reverted?: string; status?: string }
export interface RevertPreview { id: string; title: string; current: string; target: string; warning?: string; canRevert: boolean; forceAllowed: boolean; needsRoot: boolean; reviewToken: string }
export interface Fact { label: string; value: string }
export interface Suggestion { title: string; detail: string; evidence: string; item?: string; level: string }
export interface StartupEntry { name: string; scope: string; path: string; status: "available" | "missing" | "disabled" | "unknown"; detail: string; command?: string }
export interface DiagReport { snapshotId: string; collectedAt: string; facts: Fact[]; disks: { name: string; model: string; size: string; ssd: boolean; discard: boolean; mount: string; freePct: number; free: string }[] | null; startup: StartupEntry[] | null; suggestions: Suggestion[] | null; battery: boolean }
export interface InvEntry { name: string; source: string; class: string; item?: string; reason: string }
export interface FFRoot { kind: string; label: string; root: string; profiles: { name: string; path: string; size: string }[] | null }
export interface StepEvent { runId: string; stepId: string; status: string; error?: string; note?: string; time: string; problem?: Problem }
export interface LogEvent { runId: string; stepId: string; time: string; line: string }

declare global {
  interface Window {
    go?: { main: { App: Record<string, (...args: unknown[]) => Promise<unknown>> } };
    runtime?: { EventsOn: (ev: string, cb: (d: unknown) => void) => () => void; BrowserOpenURL: (u: string) => void; ClipboardSetText?: (text: string) => Promise<boolean> };
  }
}

function call<T>(name: string, ...args: unknown[]): Promise<T> {
  const fn = window.go?.main?.App?.[name];
  if (!fn) return Promise.reject(new Error("Ponte Wails indisponível (" + name + ")"));
  return fn(...args) as Promise<T>;
}

export const api = {
  init: () => call<InitData>("Init"),
  about: () => call<AboutInfo>("About"),
  states: (s: Selections) => call<Record<string, ItemState>>("States", s),
  itemState: (id: string, s: Sel) => call<ItemState>("ItemState", id, s),
  savePrefs: (p: Prefs) => call<void>("SavePrefs", p),
  buildPlan: (s: Selections, r: string[]) => call<Plan>("BuildPlan", s, r),
  execute: (s: Selections, r: string[], confirmed: Record<string, string>, token: string) => call<string>("Execute", s, r, confirmed, token),
  scriptPreview: (item: string, options: Record<string, string[]> | undefined, expected: string) => call<string>("ScriptPreview", item, options || {}, expected),
  endSession: (action: "logout" | "reboot") => call<void>("EndSession", action),
  cancel: () => call<void>("Cancel"),
  running: () => call<boolean>("Running"),
  runs: () => call<Run[] | null>("Runs"),
  runLog: (id: string) => call<string>("RunLog", id),
  changes: () => call<Change[] | null>("Changes"),
  previewRevert: (id: string) => call<RevertPreview>("PreviewRevert", id),
  revertReviewed: (id: string, force: boolean, token: string) => call<void>("RevertReviewed", id, force, token),
  profiles: () => call<Profile[]>("Profiles"),
  saveProfile: (p: Profile, overwrite: boolean) => call<void>("SaveProfile", p, overwrite),
  saveImported: (p: Profile, overwrite: boolean) => call<void>("SaveImported", p, overwrite),
  deleteProfile: (source: string, id: string) => call<void>("DeleteProfile", source, id),
  exportProfile: (p: Profile) => call<string>("ExportProfile", p),
  importProfile: () => call<ProfileReview | null>("ImportProfile"),
  reviewProfile: (p: Profile) => call<ProfileReview>("ReviewProfile", p),
  mergeSelections: (cur: Selections, origins: Origins, incoming: Selections, label: string) => call<MergeResult>("MergeSelections", cur, origins, incoming, label),
  parseInventory: (t: string) => call<InvEntry[] | null>("ParseInventory", t),
  openInventoryFile: () => call<string>("OpenInventoryFile"),
  diagnostics: () => call<DiagReport>("Diagnostics"),
  exportDiagnostics: (snapshot: string, format: "pdf" | "csv" | "xlsx", sections: string[]) => call<string>("ExportDiagnostics", snapshot, format, sections),
  firefox: () => call<{ roots: FFRoot[] | null; running: boolean }>("Firefox"),
  firefoxBackup: (p: string) => call<string>("FirefoxBackup", p),
  firefoxMigrate: (p: string) => call<string>("FirefoxMigrate", p),
  sshKeys: () => call<{ file: string; fingerprint: string }[] | null>("SSHKeys"),
};

export function on<T>(ev: string, cb: (d: T) => void): () => void {
  return window.runtime?.EventsOn(ev, (d) => cb(d as T)) ?? (() => {});
}

export function errText(e: unknown): string {
  return e instanceof Error ? e.message : String(e);
}
