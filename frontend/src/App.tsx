import { useCallback, useEffect, useMemo, useRef, useState, type ReactNode } from "react";
import {
  AppWindow, Code2, Palette, Puzzle, BatteryCharging, Stethoscope, FolderInput, Bookmark, ListChecks, History, Home,
  Search, Moon, Sun, Monitor, X, AlertTriangle, Brush, Server, ShieldCheck, Info, BookOpen, ChevronDown, Undo2, Sparkles,
} from "lucide-react";
import { api, on, errText, InitData, ItemState, Selections, Plan, StepEvent, LogEvent, Item, Origins, ProfileReview, MergeConflict, Profile, IntegrationPackage } from "./api";
import { ItemsPage } from "./Catalog";
import { PlanPage, RunState } from "./PlanRun";
import { HomePage, ProfilesPage, DiagnosticsPage, MigrationPage, HistoryPage } from "./Pages";
import { Wizard } from "./Wizard";
import { ConflictDialog, ReviewDialog, Modal } from "./Modals";
import { profileKey, profileScope } from "./profileScope";
import { includeInstalled, matchingStates, optionSignature } from "./installedSelection";
import { PluginsPage } from "./Plugins";
import { AboutPage, HelpPage } from "./HelpAbout";
import { IntegrationsPage } from "./Integrations";
import icon from "./icon.png";

export const AREAS: Record<string, { title: string; desc: string }> = {
  services: { title: "Serviços locais", desc: "Servidores web, bases de dados, cache, mensageria e contentores para desenvolvimento. Reveja os serviços e as portas antes de instalar." },
  plugins: { title: "Plugins e complementos", desc: "Organize os complementos por ferramenta. Instale só o que utiliza e siga o guia de ativação de cada grupo." },
  ai: { title: "Agentes de IA", desc: "Ferramentas para programar no terminal. Escolha o agente que quer instalar e depois configure o acesso dentro dele." },
  apps: { title: "Aplicações", desc: "Ferramentas, navegadores, escritório e utilitários. Nada é instalado sem rever o plano." },
  dev: { title: "Desenvolvimento", desc: "PHP/Laravel, Node, Docker, Git, bases de dados e editores — com versões explícitas e conflitos resolvidos." },
  infra: { title: "Infraestrutura", desc: "Ferramentas de cliente para acesso remoto, rede, diagnóstico, virtualização e discos. Não configura servidores, firewall, DNS, VPN nem serviços." },
  security: { title: "Segurança", desc: "Gestores de senhas, 2FA, cifragem, privacidade e auditoria. A firewall é só instalada, nunca ativada automaticamente." },
  style: { title: "Estilo", desc: "Cores, wallpapers, transparência do terminal e da dock, cursor e fontes. Cada ajuste mostra o valor atual e o proposto." },
  desktop: { title: "Desktop", desc: "Dock, janelas, atalhos, ficheiros e acessibilidade. Cada ajuste mostra o valor atual e o proposto." },
  extensions: { title: "Extensões", desc: "Um conjunto pequeno e coerente. Equivalentes já presentes no sistema são detetados e não duplicados." },
  energy: { title: "Energia e hardware", desc: "Perfis de energia (alternativos entre si), suspensão, disco e drivers — só com evidência do Diagnóstico." },
};

export const UNDOABLE = ["gsettings", "gitconfig", "template", "admin"];

export function sensitiveReason(it?: Item): string {
  if (!it) return "";
  if (it.sensitive) return it.sensitive;
  if (it.method === "user-script") return "Executa um instalador do fornecedor (confiança no fornecedor)";
  if (it.method === "npm-user") return "Instala código do fornecedor como o seu utilizador";
  if (it.method === "apt-repo") return "Adiciona um repositório de terceiros ao APT";
  if (it.service) return "Instala um serviço local que pode arrancar com o sistema";
  if (it.method === "admin") return "Operação administrativa no sistema";
  return "";
}

export function originLabel(o: string) {
  if (o === "manual") return "escolha manual";
  if (o === "estado") return "incluído do estado atual";
  return o.replace(/^perfil:/, "perfil ");
}

export type Ctx = {
  data: InitData;
  states: Record<string, ItemState>;
  sel: Selections;
  origins: Origins;
  reverts: string[];
  items: Record<string, Item>;
  toggle: (id: string) => void;
  toggleRevert: (id: string) => void;
  setOpt: (id: string, key: string, values: string[]) => void;
  setSel: (s: Selections, origins?: Origins) => void;
  pickGroup: (group: string, id: string | null) => void;
  includeApplied: (area?: string, candidates?: string[]) => void;
  applyProfiles: (ps: Profile[], replace: boolean) => Promise<boolean>;
  openReview: (r: ProfileReview) => void;
  openWizard: () => void;
  go: (page: string) => void;
  openPlugins: (tool: string) => void;
  toast: (msg: string, err?: boolean) => void;
  refresh: () => Promise<void>;
  showAdv: boolean;
  query: string;
  reviewPlan: (s?: Selections, r?: string[]) => void;
  setProfiles: (p: Profile[]) => void;
  editing: Profile | null;
  startEdit: (p: Profile) => void;
  setReverts: (ids: string[]) => void;
  editPlanItem: (id: string, plan: Plan) => void;
  manualProgress: Record<string, string[]>;
  toggleManualStep: (id: string, step: string) => void;
  catalogProfile: Profile | null;
  catalogScope: Set<string> | null;
  showCatalogProfile: (key: string) => void;
  openHelp: (topic?: string) => void;
  clearSelection: () => void;
  running: boolean;
};

type Theme = "system" | "light" | "dark";
type PendingMerge = { selections: Selections; origins: Origins; conflicts: MergeConflict[] };

export default function App() {
  const [data, setData] = useState<InitData | null>(null);
  const [states, setStates] = useState<Record<string, ItemState>>({});
  const [sel, setSelState] = useState<Selections>({});
  const [origins, setOrigins] = useState<Origins>({});
  const [excludedInstalled, setExcludedInstalled] = useState<string[]>([]);
  const excludedRef = useRef<string[]>([]);
  excludedRef.current = excludedInstalled;
  const editingRef = useRef(false);
  const [reverts, setReverts] = useState<string[]>([]);
  const [onboarded, setOnboarded] = useState(true);
  const [page, setPage] = useState("home");
  const [pluginTool, setPluginTool] = useState("all");
  const [integrations, setIntegrations] = useState<IntegrationPackage[]>([]);
  const [query, setQuery] = useState("");
  const [theme, setTheme] = useState<Theme>("system");
  const [showAdv, setShowAdv] = useState(false);
  const [showAllAreas, setShowAllAreas] = useState(false);
  const [manualProgress, setManualProgress] = useState<Record<string, string[]>>({});
  const [catalogProfileKey, setCatalogProfileKey] = useState("");
  const [helpTopic, setHelpTopic] = useState("primeiros-passos");
  const [plan, setPlan] = useState<Plan | null>(null);
  const [planBusy, setPlanBusy] = useState(false);
  const planRequest = useRef(0);
  const [executeBusy, setExecuteBusy] = useState(false);
  const [run, setRun] = useState<RunState | null>(null);
  const [review, setReview] = useState<ProfileReview | null>(null);
  const [merge, setMerge] = useState<PendingMerge | null>(null);
  const [edit, setEdit] = useState<{ profile: Profile; selection: Selections; origins: Origins; reverts: string[] } | null>(null);
  const [savingEdit, setSavingEdit] = useState(false);
  const [planItem, setPlanItem] = useState<string | null>(null);
  const contentRef = useRef<HTMLElement>(null);
  const [toastMsg, setToastMsg] = useState<{ msg: string; err: boolean } | null>(null);
  const [loadErr, setLoadErr] = useState("");
  const [clearConfirm, setClearConfirm] = useState(false);
  const [selectionBackup, setSelectionBackup] = useState<{ selections: Selections; origins: Origins; reverts: string[] } | null>(null);
  const loaded = useRef(false);
  const selRef = useRef<Selections>({});
  selRef.current = sel;
  const stateItems = useRef<Record<string, Item>>({});
  const originsRef = useRef<Origins>({});
  originsRef.current = origins;
  const revRef = useRef<string[]>([]);
  revRef.current = reverts;
  // Events can arrive before Execute's promise returns the run ID.
  const starting = useRef(false);
  const earlySteps = useRef<StepEvent[]>([]);
  const earlyLogs = useRef<LogEvent[]>([]);
  const earlyDone = useRef<{ runId: string; status: string; error?: string } | null>(null);
  const latestSteps = useRef<Record<string, StepEvent>>({});
  const executionPlan = useRef<Plan | null>(null);
  const runRef = useRef(run);
  runRef.current = run;
  const statesRef = useRef(states);
  statesRef.current = states;
  editingRef.current = !!edit;

  const toast = useCallback((msg: string, err = false) => {
    setToastMsg({ msg, err });
    window.setTimeout(() => setToastMsg((t) => (t?.msg === msg ? null : t)), err ? 9000 : 4500);
  }, []);

  const acceptStates = useCallback((next: Record<string, ItemState>, snapshot: Selections) => {
    const accepted = matchingStates(stateItems.current, next, snapshot, selRef.current);
    statesRef.current = { ...statesRef.current, ...accepted };
    setStates((old) => ({ ...old, ...accepted }));
    if (editingRef.current) return;
    const merged = includeInstalled(Object.values(stateItems.current), accepted, selRef.current, originsRef.current, excludedRef.current, revRef.current);
    if (!merged.added.length) return;
    selRef.current = merged.selections;
    originsRef.current = merged.origins;
    planRequest.current++;
    setSelState(merged.selections);
    setOrigins(merged.origins);
  }, []);
  const refresh = useCallback(async () => {
    const snapshot = selRef.current;
    try { acceptStates(await api.states(snapshot), snapshot); } catch (e) { toast(errText(e), true); }
  }, [acceptStates, toast]);

  useEffect(() => {
    api.init().then((d) => {
      setData(d);
      api.integrations().then(setIntegrations).catch((e) => toast("Não foi possível ler os pacotes: " + errText(e), true));
      stateItems.current = Object.fromEntries(d.items.map((it) => [it.id, it]));
      const s = d.prefs.selections || {};
      setSelState(s);
      selRef.current = s;
      setOrigins(d.prefs.origins || {});
      originsRef.current = d.prefs.origins || {};
      setReverts(d.prefs.reverts || []);
      revRef.current = d.prefs.reverts || [];
      setExcludedInstalled(d.prefs.excludedInstalled || []);
      excludedRef.current = d.prefs.excludedInstalled || [];
      setOnboarded(!!d.prefs.onboarded || Object.keys(s).length > 0 || !!d.pendingImport);
      setTheme((d.prefs.theme as Theme) || "system");
      setShowAdv(!!d.prefs.showAdvanced);
      setShowAllAreas(!!d.prefs.showAllAreas);
      setManualProgress(d.prefs.manualProgress || {});
      setCatalogProfileKey(d.prefs.catalogProfile || "");
      loaded.current = true;
      api.states(s).then((next) => acceptStates(next, s)).catch((e) => toast(errText(e), true));
      if (d.pendingImport) setReview(d.pendingImport);
      if (d.pendingError) toast("Perfil recusado: " + d.pendingError, true);
    }).catch((e) => setLoadErr(errText(e)));
  }, [acceptStates, toast]);

  useEffect(() => {
    const mq = window.matchMedia("(prefers-color-scheme: dark)");
    const apply = () => {
      document.documentElement.dataset.theme = theme === "dark" || (theme === "system" && mq.matches) ? "dark" : "light";
    };
    apply();
    mq.addEventListener("change", apply);
    return () => mq.removeEventListener("change", apply);
  }, [theme]);

  useEffect(() => {
    if (!loaded.current || !data) return;
    const t = window.setTimeout(() => {
      api.savePrefs({ theme, selections: sel, showAdvanced: showAdv, showAllAreas, reverts, origins, onboarded, manualProgress, catalogProfile: catalogProfileKey, excludedInstalled }).catch((e) => toast("Não foi possível guardar as preferências: " + errText(e), true));
    }, 400);
    return () => window.clearTimeout(t);
  }, [sel, theme, showAdv, showAllAreas, data, reverts, origins, onboarded, manualProgress, catalogProfileKey, excludedInstalled, toast]);

  useEffect(() => {
    const offStep = on<StepEvent>("run:step", (e) => {
      latestSteps.current[e.stepId] = e;
      if (starting.current) earlySteps.current.push(e);
      setRun((r) => (r && r.id === e.runId ? { ...r, steps: { ...r.steps, [e.stepId]: e }, last: e.time } : r));
    });
    const offLog = on<LogEvent>("run:log", (e) => {
      if (starting.current) earlyLogs.current.push(e);
      setRun((r) => {
      if (!r || r.id !== e.runId) return r;
      const prev = r.logs[e.stepId] || [];
      return { ...r, logs: { ...r.logs, [e.stepId]: [...prev.slice(-1500), e] }, last: e.time, lastLine: e.line };
      });
    });
    const offDone = on<{ runId: string; status: string; error?: string }>("run:done", (e) => {
      if (starting.current) earlyDone.current = e;
      setRun((r) => (r && r.id === e.runId ? { ...r, finished: e.status, error: e.error } : r));
      const active = executionPlan.current;
      if (active) {
        const successful = active.steps.filter((s) => s.op === "undo" || s.op === "revert").filter((s) => latestSteps.current[s.id]?.runId === e.runId && latestSteps.current[s.id]?.status === "done").map((s) => s.item);
        setReverts((r) => r.filter((id) => !successful.includes(id)));
      }
      if (e.error) toast(e.error, true);
      refresh();
    });
    const offReview = on<ProfileReview>("profile:review", (r) => setReview(r));
    const offReviewErr = on<string>("profile:review-error", (m) => toast("Perfil recusado: " + m, true));
    return () => { offStep(); offLog(); offDone(); offReview(); offReviewErr(); };
  }, [refresh, toast]);

  const items = useMemo(() => Object.fromEntries((data?.items || []).map((i) => [i.id, i])), [data]);
  const catalogProfile = useMemo(() => data?.profiles.find((p) => profileKey(p) === catalogProfileKey) || null, [data?.profiles, catalogProfileKey]);
  const catalogScope = useMemo(() => profileScope(catalogProfile, data?.items || []), [catalogProfile, data?.items]);
  useEffect(() => {
    if (data && catalogProfileKey && (!catalogProfile || !Object.keys(catalogProfile.selections || {}).length)) setCatalogProfileKey("");
  }, [data, catalogProfileKey, catalogProfile]);
  useEffect(() => {
    planRequest.current++;
    setPlan(null);
    setPlanBusy(false);
    if (Object.keys(sel).length || reverts.length) setSelectionBackup(null);
  }, [sel, reverts]);

  const toggleManualStep = (id: string, step: string) => {
    const allowed = items[id]?.manual?.steps || [];
    if (!allowed.includes(step)) return;
    setManualProgress((previous) => {
      const completed = (previous[id] || []).filter((s) => allowed.includes(s));
      const next = { ...previous };
      next[id] = completed.includes(step) ? completed.filter((s) => s !== step) : [...completed, step];
      if (!next[id].length) delete next[id];
      return next;
    });
  };

  const setSel = useCallback((s: Selections, o?: Origins) => {
    if (!editingRef.current) {
      const removed = Object.keys(selRef.current).filter((id) => !s[id]);
      const installed = Object.keys(statesRef.current).filter((id) => !s[id] && stateItems.current[id]?.method !== "admin" && ["installed", "provided"].includes(statesRef.current[id]?.status));
      const excluded = [...new Set([...excludedRef.current, ...removed, ...installed])].filter((id) => !s[id]);
      excludedRef.current = excluded;
      setExcludedInstalled(excluded);
    }
    const nextOrigins = o ?? Object.fromEntries(Object.keys(s).map((id) => [id, originsRef.current[id] || ["manual"]]));
    selRef.current = s;
    originsRef.current = nextOrigins;
    setSelState(s);
    setOrigins(nextOrigins);
  }, []);

  // Checkbox = choice kept in the profile. It never changes the system by itself.
  const toggle = useCallback((id: string) => {
    const was = !!selRef.current[id];
    const s = { ...selRef.current }, o = { ...originsRef.current };
    if (was) { delete s[id]; delete o[id]; } else { s[id] = {}; o[id] = ["manual"]; }
    setSel(s, o);
    if (!was) { revRef.current = revRef.current.filter((x) => x !== id); setReverts(revRef.current); }
  }, [setSel]);

  // Explicit revert of an applied adjustment (reviewed in the plan before running).
  const toggleRevert = useCallback((id: string) => {
    revRef.current = revRef.current.includes(id) ? revRef.current.filter((x) => x !== id) : [...revRef.current, id];
    setReverts(revRef.current);
    if (selRef.current[id]) {
      const s = { ...selRef.current }, o = { ...originsRef.current };
      delete s[id]; delete o[id]; setSel(s, o);
    }
  }, [setSel]);

  const optTimers = useRef<Record<string, number>>({});
  const optVersions = useRef<Record<string, number>>({});
  const previousSelection = useRef<Selections | null>(null);
  useEffect(() => {
    if (!data) return;
    const previous = previousSelection.current;
    previousSelection.current = sel;
    if (!previous) return;
    data.items.forEach((it) => {
      if (!it.options?.length || optionSignature(it, previous[it.id]) === optionSignature(it, sel[it.id])) return;
      const id = it.id;
      window.clearTimeout(optTimers.current[id]);
      const version = (optVersions.current[id] || 0) + 1;
      optVersions.current[id] = version;
      setStates((old) => ({ ...old, [id]: { status: "unknown", detail: "A detetar o estado para as opções atuais…" } }));
      optTimers.current[id] = window.setTimeout(() => {
        const snapshot = selRef.current[id] || {};
        api.itemState(id, snapshot).then((st) => {
          if (version === optVersions.current[id] && optionSignature(it, snapshot) === optionSignature(it, selRef.current[id])) {
            setStates((old) => ({ ...old, [id]: st }));
          }
        }).catch((e) => { if (version === optVersions.current[id]) toast(errText(e), true); });
      }, 300);
    });
  }, [data, sel, toast]);
  useEffect(() => () => { Object.values(optTimers.current).forEach(window.clearTimeout); }, []);
  const setOpt = useCallback((id: string, key: string, values: string[]) => {
    setSel({ ...selRef.current, [id]: { options: { ...(selRef.current[id]?.options || {}), [key]: values } } },
      { ...originsRef.current, [id]: [...new Set([...(originsRef.current[id] || []), "manual"])] });
  }, [setSel]);

  const pickGroup = useCallback((group: string, id: string | null) => {
    const members = (data?.items || []).filter((i) => i.group === group).map((i) => i.id);
    const s = { ...selRef.current }, o = { ...originsRef.current };
    members.forEach((m) => { delete s[m]; delete o[m]; });
    if (id) { s[id] = selRef.current[id] || {}; o[id] = ["manual"]; }
    setSel(s, o);
  }, [data, setSel]);

  // Restore installed choices that were explicitly excluded by the user.
  const includeApplied = useCallback((area?: string, candidates?: string[]) => {
    const shown = (data?.items || []).filter((it) => (!area || it.area === area) && (!candidates || candidates.includes(it.id)));
    const visibleStates = Object.fromEntries(shown.map((it) => [it.id, states[it.id]]));
    const merged = includeInstalled(data?.items || [], visibleStates, selRef.current, originsRef.current, [], revRef.current);
    if (!merged.added.length) { toast("Os itens instalados desta vista já estão marcados."); return; }
    setSel(merged.selections, merged.origins);
    toast(`${merged.added.length} item(ns) instalados novamente marcados.`);
  }, [data, states, setSel, toast]);

  // Combine profiles: duplicates merge, conflicts go to the dialog (never decided by order).
  const applyProfiles = useCallback(async (ps: Profile[], replace: boolean) => {
    let cur: Selections = replace ? {} : selRef.current;
    let org: Origins = replace ? {} : originsRef.current;
    const conflicts: MergeConflict[] = [];
    try {
      for (const p of ps) {
        const r = await api.mergeSelections(cur, org, p.selections || {}, "perfil:" + p.name);
        cur = r.selections; org = r.origins;
        conflicts.push(...(r.conflicts || []));
      }
    } catch (e) { toast(errText(e), true); return false; }
    if (replace && !editingRef.current) {
      const merged = includeInstalled(Object.values(stateItems.current), statesRef.current, cur, org, excludedRef.current, revRef.current);
      cur = merged.selections; org = merged.origins;
    }
    if (conflicts.length) { setMerge({ selections: cur, origins: org, conflicts }); return true; }
    setSel(cur, org);
    setReverts((r) => r.filter((x) => !cur[x]));
    toast(`Perfil(s) ${ps.map((p) => "«" + p.name + "»").join(", ")} ${replace ? "aplicado(s)" : "juntado(s)"} à seleção. Nada foi executado.`);
    return true;
  }, [setSel, toast]);

  const reviewPlan = useCallback(async (s?: Selections, r?: string[]) => {
    const request = ++planRequest.current;
    setRun((previous) => previous?.finished ? null : previous);
    setPlanItem(null);
    setQuery("");
    setPage("plan");
    setPlanBusy(true);
    try { const next = await api.buildPlan(s || selRef.current, r || revRef.current); if (request === planRequest.current) setPlan(next); }
    catch (e) { if (request === planRequest.current) {setPlan(null); toast(errText(e), true);} }
    finally { if (request === planRequest.current) setPlanBusy(false); }
  }, [toast]);

  const go = (next: string) => {
    if (next === "plan" && !runRef.current?.id) { reviewPlan(); return; }
    setQuery(""); setPage(next);
  };
  const openHelp = (topic = "primeiros-passos") => { setHelpTopic(topic); go("help"); };
  const showCatalogProfile = (key: string) => { setCatalogProfileKey(key); go("catalog"); setPlanItem(null); };
  useEffect(() => { contentRef.current?.scrollTo({ top: 0 }); }, [page]);
  const editPlanItem = (id: string, reviewed: Plan) => {
    const item = items[id];
    if (!item) return;
    setSel(reviewed.selection);
    setReverts(reviewed.reverts || []);
    setPlanItem(id); setPage(item.area); setQuery(item.name);
  };
  const startEdit = (p: Profile) => {
    if (edit) { toast("Guarde ou descarte o perfil em edição antes de abrir outro.", true); return; }
    editingRef.current = true;
    setEdit({ profile: structuredClone(p), selection: structuredClone(sel), origins: structuredClone(origins), reverts: [...reverts] });
    setSel(structuredClone(p.selections), Object.fromEntries(Object.keys(p.selections).map((id) => [id, ["perfil:" + p.name]])));
    setReverts([]);
  };
  const saveEdit = async () => {
    if (!edit || savingEdit) return;
    setSavingEdit(true);
    try {
      const p = { ...edit.profile, selections: selRef.current };
      if (p.source === "imported") await api.saveImported(p, true); else await api.saveProfile(p, true);
      const profiles = await api.profiles();
      setData((d) => d ? { ...d, profiles } : d);
      setEdit(null);
      editingRef.current = false;
      setSel(selRef.current);
      toast("Perfil «" + p.name + "» atualizado.");
    } catch (e) { toast(errText(e), true); }
    finally { setSavingEdit(false); }
  };
  const discardEdit = () => {
    if (!edit) return;
    setSel(edit.selection, edit.origins); setReverts(edit.reverts); setEdit(null);
    editingRef.current = false;
    toast("Edição descartada; seleção anterior restaurada.");
  };

  if (loadErr) return <div className="empty"><AlertTriangle /> <p>Não foi possível iniciar: {loadErr}</p></div>;
  if (!data) return <div className="empty">A carregar…</div>;

  const ctx: Ctx = {
    data, states, sel, origins, reverts, items, toggle, toggleRevert, setOpt, setSel, pickGroup, includeApplied, applyProfiles,
    openReview: setReview, openWizard: () => setOnboarded(false), go, openPlugins: (tool) => { setPluginTool(tool); go("plugins"); }, toast, refresh, showAdv, query, reviewPlan, editing: edit?.profile || null, startEdit,
    setProfiles: (p) => setData((d) => (d ? { ...d, profiles: p } : d)),
    setReverts, editPlanItem, manualProgress, toggleManualStep, catalogProfile, catalogScope, showCatalogProfile, openHelp,
    clearSelection: () => { if (!runRef.current || runRef.current.finished) setClearConfirm(true); },
    running: !!run && !run.finished,
  };
  const pending = (id: string) => states[id]?.status !== "installed" && states[id]?.status !== "provided";
  const inProfile = Object.keys(sel).length;
  const inArea = (id: string, a: string) => a === "services" ? items[id]?.category.startsWith("Serviços locais ·") || items[id]?.category === "Contentores" : items[id]?.area === a;
  const areaCount = (a: string) => Object.keys(sel).filter((id) => inArea(id, a) && pending(id)).length + reverts.filter((id) => inArea(id, a)).length;
  const running = !!run && !run.finished;

  if (!onboarded) {
    return (
      <>
        <Wizard ctx={ctx} onDone={() => { setOnboarded(true); setCatalogProfileKey(""); go("catalog"); }} />
        {merge && <ConflictDialog ctx={ctx} merge={merge} onDone={(s, o) => { setSel(s, o); setMerge(null); }} onCancel={() => setMerge(null)} />}
        {toastMsg && <Toast t={toastMsg} close={() => setToastMsg(null)} />}
      </>
    );
  }

  const NAV: [string, string, ReactNode][] = [
    ["apps", "Aplicações", <AppWindow size={17} />],
    ["dev", "Desenvolvimento", <Code2 size={17} />],
    ["services", "Serviços locais", <Server size={17} />],
    ["infra", "Infraestrutura", <Server size={17} />],
    ["security", "Segurança", <ShieldCheck size={17} />],
    ["style", "Estilo", <Brush size={17} />],
    ["desktop", "Desktop", <Palette size={17} />],
    ["extensions", "Extensões", <Puzzle size={17} />],
    ["energy", "Energia", <BatteryCharging size={17} />],
  ];
  const NAV2: [string, string, ReactNode][] = [
    ["diag", "Diagnóstico", <Stethoscope size={17} />],
    ["migration", "Migração", <FolderInput size={17} />],
    ["history", "Histórico e reversão", <History size={17} />],
    ["help", "Ajuda e manual", <BookOpen size={17} />],
    ["integrations", "Pacotes de integração", <Puzzle size={17} />],
  ];
  const cycleTheme = () => setTheme(theme === "system" ? "light" : theme === "light" ? "dark" : "system");
  const themeName = { system: "sistema", light: "claro", dark: "escuro" }[theme];

  return (
    <div className="app">
      <aside className="side" aria-label="Navegação">
        <div className="brand">
          <img src={icon} alt="" />
          <div><b>Meu Ubuntu</b><small>{data.system.pretty || "Ubuntu"}</small></div>
        </div>
        <nav className="side-nav" aria-label="Secções da aplicação">
        <button className={"nav" + (page === "home" ? " on" : "")} onClick={() => go("home")} aria-current={page === "home" ? "page" : undefined}><Home size={17} /> Início</button>
        <div className="nav-sec">Preparar</div>
        <button className={"nav" + (page === "catalog" || page === "search" ? " on" : "")} onClick={() => { setCatalogProfileKey(""); go("catalog"); }} aria-current={page === "catalog" ? "page" : undefined}><AppWindow size={17} /> Catálogo</button>
        <button className={"nav" + (page === "profiles" ? " on" : "")} onClick={() => go("profiles")} aria-current={page === "profiles" ? "page" : undefined}><Bookmark size={17} /> Perfis prontos</button>
        <button className={"nav" + (page === "ai" ? " on" : "")} onClick={() => go("ai")} aria-current={page === "ai" ? "page" : undefined}><Sparkles size={17} /> Agentes de IA</button>
        <button className={"nav" + (page === "plugins" ? " on" : "")} onClick={() => { setPluginTool("all"); go("plugins"); }} aria-current={page === "plugins" ? "page" : undefined}><Puzzle size={17} /> Plugins e complementos</button>
        <button className={"nav" + (page === "selection" ? " on" : "")} onClick={() => go("selection")} aria-current={page === "selection" ? "page" : undefined}><ListChecks size={17} /> Minhas escolhas{inProfile > 0 && <span className="count">{inProfile}</span>}</button>
        <button className="nav nav-disclosure" aria-expanded={showAllAreas} aria-controls="catalog-areas" onClick={() => setShowAllAreas(!showAllAreas)}><ChevronDown size={15} className={showAllAreas ? "rotated" : ""} /> Explorar por área</button>
        {showAllAreas && <div id="catalog-areas" className="nav-areas">{NAV.map(([id, label, ic]) => (
          <button key={id} className={"nav" + (page === id ? " on" : "")} onClick={() => go(id)} aria-current={page === id ? "page" : undefined}>
            {ic}{label}{AREAS[id] && areaCount(id) ? <span className="count" title="alterações pendentes">{areaCount(id)}</span> : null}
          </button>
        ))}</div>}
        <div className="nav-sec">Aplicar</div>
        <button className={"nav" + (page === "plan" ? " on" : "")} onClick={() => running ? go("plan") : reviewPlan()} aria-current={page === "plan" ? "page" : undefined}><ListChecks size={17} /> {running ? "Acompanhar execução" : "Rever e aplicar"}</button>
        <div className="nav-sec">Ferramentas</div>
        {NAV2.map(([id, label, ic]) => (
          <button key={id} className={"nav" + (page === id ? " on" : "")} onClick={() => go(id)} aria-current={page === id ? "page" : undefined}>
            {ic}{label}{id === "plan" && running ? <span className="count">a correr</span> : null}
          </button>
        ))}
        {integrations.length > 0 && <div className="nav-sec">Equipa</div>}
        {integrations.map(p => <button key={p.id} className={"nav" + (page === "integration:" + p.id ? " on" : "")} onClick={() => go("integration:" + p.id)} aria-current={page === "integration:" + p.id ? "page" : undefined}><Server size={17} /> {p.name}</button>)}
        </nav>
        <div className="side-foot">
          <button className={"nav side-about" + (page === "about" ? " on" : "")} onClick={() => go("about")} aria-current={page === "about" ? "page" : undefined} title={"Sobre o Meu Ubuntu · versão em execução " + data.version}><Info size={17} /><span>Sobre</span><span className="side-version">{data.version}</span></button>
          <div className="side-system">GNOME {data.system.gnome || "?"} · {data.system.session || "?"} · {data.system.arch}</div>
          {data.system.simulated && <div style={{ color: "var(--warn)", fontWeight: 600, marginTop: 4 }}>Modo simulação: nada é instalado</div>}
        </div>
      </aside>
      <div className="main">
        <header className="top">
          <label className="search">
            <Search size={16} className="muted" />
            <input placeholder="Pesquisar aplicações e ajustes…" value={query} onChange={(e) => { setQuery(e.target.value); if (e.target.value.trim()) setPage("search"); else if (page === "search") setPage("catalog"); }} aria-label="Pesquisar em todo o catálogo" />
            {query && <button className="iconbtn" style={{ border: 0, padding: 2 }} onClick={() => go("catalog")} aria-label="Limpar pesquisa"><X size={14} /></button>}
          </label>
          <div className="spacer" />
          <button className="iconbtn" onClick={cycleTheme} title={"Tema: " + themeName} aria-label={"Mudar tema (atual: " + themeName + ")"}>
            {theme === "system" ? <Monitor size={17} /> : theme === "light" ? <Sun size={17} /> : <Moon size={17} />}
          </button>
          <div className="selbar">
            <button className="btn ghost choices-link" onClick={() => go("selection")} aria-label={`Minhas escolhas: ${inProfile} ${inProfile === 1 ? "item" : "itens"}`}><Bookmark size={16} /><b>{inProfile}</b> {inProfile === 1 ? "escolha" : "escolhas"}</button>
            <button className="btn primary" disabled={(!inProfile && !reverts.length) || running || planBusy || executeBusy} onClick={() => reviewPlan()}>
              <ListChecks size={16} /> Rever plano
            </button>
          </div>
        </header>
        {(AREAS[page] || ["catalog", "search", "selection"].includes(page)) && <div className="catalog-preferences" role="region" aria-label="Preferências do catálogo"><label><input type="checkbox" checked={showAdv} onChange={(e) => setShowAdv(e.target.checked)} /> Mostrar opções avançadas</label><span>Os já instalados ficam marcados automaticamente. Desmarcar não desinstala nem reverte.</span></div>}
        {selectionBackup && <div className="editbar" role="status"><span>As escolhas foram limpas.</span><button className="btn sm" onClick={() => { setSel(selectionBackup.selections, selectionBackup.origins); setReverts(selectionBackup.reverts); setSelectionBackup(null); toast("Escolhas anteriores recuperadas."); }}><Undo2 size={14} /> Desfazer limpeza</button></div>}
        {edit && <div className="editbar" role="status">
          <span>A editar <b>{edit.profile.name}</b> · {inProfile} itens · {JSON.stringify(Object.keys(sel).sort().map((id) => [id, items[id] && optionSignature(items[id], sel[id])])) === JSON.stringify(Object.keys(edit.profile.selections).sort().map((id) => [id, items[id] && optionSignature(items[id], edit.profile.selections[id])])) ? "Sem alterações" : "Alterações por guardar"}</span>
          <button className="btn sm primary" disabled={savingEdit} onClick={saveEdit}>{savingEdit ? "A guardar…" : "Guardar perfil"}</button>
          <button className="btn sm" disabled={savingEdit} onClick={discardEdit}>Descartar edição</button>
        </div>}
        {planItem && <div className="editbar">
          <span>A ajustar o plano: <b>{items[planItem]?.name}</b></span>
          <button className="btn sm primary" onClick={() => reviewPlan()}>Voltar à revisão</button>
        </div>}
        <main className="content" ref={contentRef}>
          {page === "home" && <HomePage ctx={ctx} run={run} />}
          {page === "plugins" && <PluginsPage ctx={ctx} tool={pluginTool} onTool={setPluginTool} />}
          {((AREAS[page] && page !== "plugins") || page === "search" || page === "catalog" || page === "selection") && <ItemsPage key={`${page}:${catalogProfileKey}`} ctx={ctx} area={page} />}
          {page === "plan" && (
            <PlanPage ctx={ctx} plan={plan} busy={planBusy} executing={executeBusy} run={run}
              onExecute={async (confirmed) => {
                if (!plan) return;
                if (starting.current) return;
                starting.current = true; earlySteps.current = []; earlyLogs.current = []; earlyDone.current = null;
                setExecuteBusy(true);
                latestSteps.current = {}; executionPlan.current = plan;
                try {
                  const id = await api.execute(plan.selection, plan.reverts || [], confirmed, plan.reviewToken);
                  const steps = Object.fromEntries(earlySteps.current.filter((e) => e.runId === id).map((e) => [e.stepId, e]));
                  const logs: Record<string, LogEvent[]> = {};
                  earlyLogs.current.filter((e) => e.runId === id).forEach((e) => (logs[e.stepId] ||= []).push(e));
                  const done = earlyDone.current as { runId: string; status: string; error?: string } | null;
                  const nextRun = { id, plan, steps, logs, started: new Date().toISOString(), finished: done?.runId === id ? done.status : undefined, error: done?.runId === id ? done.error : undefined };
                  runRef.current = nextRun; setRun(nextRun);
                  if (done?.runId === id) {
                    const successful = plan.steps.filter((s) => (s.op === "undo" || s.op === "revert") && steps[s.id]?.status === "done").map((s) => s.item);
                    setReverts((r) => r.filter((x) => !successful.includes(x)));
                  }
                } catch (e) { toast(errText(e), true); }
                finally { starting.current = false; setExecuteBusy(false); }
              }}
              onRetry={() => {
                const pending = (run?.plan.reverts || []).filter((id) => run?.plan.steps.some((s) => s.item === id && (s.op === "undo" || s.op === "revert") && run.steps[s.id]?.status !== "done"));
                setRun(null); reviewPlan(run?.plan.selection, pending);
              }}
              onClose={() => { setRun(null); setPlan(null); go("home"); }}
            />
          )}
          {page === "profiles" && <ProfilesPage ctx={ctx} />}
          {page === "diag" && <DiagnosticsPage ctx={ctx} />}
          {page === "migration" && <MigrationPage ctx={ctx} />}
          {page === "history" && <HistoryPage ctx={ctx} />}
          {page === "help" && <HelpPage ctx={ctx} topic={helpTopic} onTopic={setHelpTopic} />}
          {page === "about" && <AboutPage ctx={ctx} />}
          {(page === "integrations" || page.startsWith("integration:")) && <IntegrationsPage selected={page.startsWith("integration:") ? page.slice(12) : ""} packages={integrations} onPackages={setIntegrations} onSelect={id => go("integration:" + id)} toast={toast} />}
        </main>
      </div>
      {review && <ReviewDialog ctx={ctx} review={review} onClose={() => setReview(null)} />}
      {merge && <ConflictDialog ctx={ctx} merge={merge} onDone={(s, o) => { setSel(s, o); setMerge(null); toast("Conflitos resolvidos; seleção atualizada. Nada foi executado."); }} onCancel={() => setMerge(null)} />}
      {clearConfirm && <Modal title="Limpar as suas escolhas?" onClose={() => setClearConfirm(false)} footer={<><button className="btn" data-autofocus onClick={() => setClearConfirm(false)}>Manter escolhas</button><button className="btn danger" onClick={() => {
        setSelectionBackup({selections:structuredClone(sel),origins:structuredClone(origins),reverts:[...reverts]});
        setSel({}, {}); setReverts([]); setClearConfirm(false);
      }}>Limpar escolhas</button></>}><p>Vai retirar {inProfile} {inProfile === 1 ? "item" : "itens"} das escolhas{reverts.length ? ` e ${reverts.length} pedidos de reversão` : ""}.</p><p>As aplicações instaladas e os perfis guardados mantêm-se. Pode desfazer esta limpeza antes de fazer novas escolhas.</p></Modal>}
      {toastMsg && <Toast t={toastMsg} close={() => setToastMsg(null)} />}
    </div>
  );
}

function Toast({ t, close }: { t: { msg: string; err: boolean }; close: () => void }) {
  return (
    <div className={"toast" + (t.err ? " err" : "")} role={t.err ? "alert" : "status"}>
      <span style={{ flex: 1 }}>{t.msg}</span>
      <button className="iconbtn" style={{ border: 0, background: "none", padding: 0, color: "inherit" }} onClick={close} aria-label="Fechar aviso"><X size={15} /></button>
    </div>
  );
}
