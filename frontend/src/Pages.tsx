import { useEffect, useState } from "react";
import {
  Bookmark, Download, Upload, Trash2, Save, FileText, Stethoscope, HardDrive, Lightbulb, AlertTriangle, CheckCircle2,
  Undo2, RefreshCw, Archive, ArrowRightLeft, Flame, KeyRound, ListChecks, Info,
} from "lucide-react";
import type { Ctx } from "./App";
import type { Profile, DiagReport, InvEntry, FFRoot, Run, Change } from "./api";
import { api, errText } from "./api";
import { ResumeNote, type RunState } from "./PlanRun";
import { ExportDialog, RevertDialog, SaveProfileDialog } from "./Modals";
import { StatusBadge } from "./Catalog";

/* ------------------------------ Início ------------------------------ */
export function HomePage({ ctx, run }: { ctx: Ctx; run: RunState | null }) {
  const [saving, setSaving] = useState(false);
  const selected = Object.keys(ctx.sel);
  const ready = selected.filter((id) => ["installed", "provided"].includes(ctx.states[id]?.status)).length;
  const manual = selected.filter((id) => ctx.states[id]?.status === "manual").length;
  const attention = selected.filter((id) => !ctx.states[id] || ["unknown", "incompatible"].includes(ctx.states[id].status)).length;
  const pending = selected.length - ready - manual - attention;
  const running = !!run && !run.finished;
  return <>
    <div className="home-intro">
      <span className="eyebrow">{ctx.data.system.pretty} · {ctx.data.system.arch}</span>
      <h1>Prepare este computador</h1>
      <p>Escolha as aplicações e os ajustes de que precisa. Veja o que vai mudar antes de aplicar.</p>
    </div>
    <ResumeNote ctx={ctx} />
    {!ctx.data.system.simulated && ctx.data.system.id !== "ubuntu" && <div className="note warn"><AlertTriangle size={16} /> Esta versão foi preparada para Ubuntu. As receitas Ubuntu estão indisponíveis neste sistema.</div>}
    {running && <div className="note info"><Info size={16} /> Há uma execução em curso. <button className="btn sm" onClick={() => ctx.go("plan")}>Ver progresso</button></div>}
    {selected.length || ctx.reverts.length ? <div className="panel profile-readiness">
      <div className="readiness-head"><div><h2><ListChecks size={18} /> Minhas escolhas</h2><p className="muted">{selected.length} aplicações e ajustes escolhidos para este computador.</p></div><div className="readiness-total"><b>{ready}<small> / {selected.length}</small></b><span>já como pretende</span></div></div>
      {selected.length > 0 && <><progress className="readiness-bar" value={ready} max={selected.length} aria-label="Escolhas confirmadas neste computador" /><div className="readiness-counts"><span><CheckCircle2 size={14} /> {ready} já prontos</span><span><Download size={14} /> {pending} para aplicar</span>{manual > 0 && <span>{manual} com passos manuais</span>}{attention > 0 && <span><AlertTriangle size={14} /> {attention} a verificar</span>}</div></>}
      {ctx.reverts.length > 0 && <p className="muted">{ctx.reverts.length} ajustes escolhidos para restaurar.</p>}
      <div className="btnrow"><button className="btn primary" onClick={() => ctx.go("selection")}>Ver minhas escolhas</button><button className="btn" disabled={running} onClick={() => ctx.reviewPlan()}>Rever e aplicar</button>{selected.length > 0 && <button className="btn" disabled={running || !!ctx.editing} onClick={() => setSaving(true)}><Save size={15} /> Guardar como perfil…</button>}</div>
      <p className="readiness-help">As escolhas ficam guardadas neste computador. Guarde um perfil para as reutilizar ou partilhar.</p>
    </div> : <div className="start-options">
      <button className="start-option" onClick={() => ctx.showCatalogProfile("")}><span className="start-option-icon"><Download size={24} /></span><h2>Escolher aplicações</h2><p>Explore o catálogo e escolha apenas o que precisa.</p><span className="start-option-action">Explorar catálogo →</span></button>
      <button className="start-option" onClick={() => ctx.go("profiles")}><span className="start-option-icon"><Bookmark size={24} /></span><h2>Usar um perfil pronto</h2><p>Comece com um conjunto de escolhas e personalize.</p><span className="start-option-action">Ver perfis →</span></button>
    </div>}
    <div className="workflow-guide" aria-label="Como preparar este computador">
      <div><span>1</span><h2>Escolher</h2><p>Adicione aplicações e ajustes às suas escolhas.</p></div>
      <div><span>2</span><h2>Rever</h2><p>Confira instalações, alterações e confirmações necessárias.</p></div>
      <div><span>3</span><h2>Aplicar</h2><p>Confirme o plano e acompanhe cada resultado.</p></div>
    </div>
    <div className="home-tools"><button className="btn" onClick={() => ctx.go("diag")}><Stethoscope size={16} /> Ver diagnóstico do computador</button><button className="btn ghost" onClick={() => ctx.openHelp()}>Ajuda para começar</button></div>
    <details className="catalog-secondary"><summary><KeyRound size={15} /> Como são protegidas as alterações</summary><ul className="muted"><li>Escolher e guardar um perfil não instala aplicações.</li><li>As operações administrativas pedem autorização ao sistema.</li><li>Ajustes registados podem ser revertidos individualmente; instalações não têm reversão automática.</li></ul></details>
    {saving && <SaveProfileDialog ctx={ctx} selections={ctx.sel} onClose={() => setSaving(false)} />}
  </>;
}

/* ------------------------------ Perfis ------------------------------ */
export function ProfilesPage({ ctx }: { ctx: Ctx }) {
  const [profiles, setProfiles] = useState<Profile[]>(ctx.data.profiles);
  const [tab, setTab] = useState("builtin");
  const [saving, setSaving] = useState(false);
  const [exporting, setExporting] = useState<Profile | null>(null);
  const reload = async () => { try { const next = await api.profiles(); setProfiles(next); ctx.setProfiles(next); } catch (e) { ctx.toast(errText(e), true); } };
  useEffect(() => { reload(); }, []);
  useEffect(() => { setProfiles(ctx.data.profiles); }, [ctx.data.profiles]);
  const saveCopy = async (profile: Profile) => {
    const copy: Profile = { ...profile, id: (profile.id || "perfil") + "-copia", name: profile.name + " (cópia)", source: undefined, builtin: false, created: undefined };
    try { await api.saveProfile(copy, false); await reload(); setTab("mine"); ctx.toast("Cópia guardada em «Meus perfis»."); }
    catch (e) { ctx.toast(errText(e), true); }
  };
  const included = profiles.filter(p => p.source === "builtin" && Object.keys(p.selections || {}).length > 0);
  const mine = profiles.filter(p => p.source === "mine");
  const imported = profiles.filter(p => p.source === "imported");
  const sections: [string, Profile[]][] = tab === "builtin" ? [["Pontos de partida", included.filter(p => p.kind !== "complement")], ["Complementos opcionais", included.filter(p => p.kind === "complement")]] : [[tab === "mine" ? "Meus perfis" : "Perfis importados", tab === "mine" ? mine : imported]];
  return <>
    <div className="page-head"><div style={{flex:1}}><h1>Perfis prontos</h1><p>Conjuntos de aplicações e ajustes para começar. Pode adicionar, substituir e personalizar as escolhas antes de aplicar.</p></div><button className="btn" disabled={ctx.running || !Object.keys(ctx.sel).length || !!ctx.editing} onClick={() => setSaving(true)}><Save size={15} /> Guardar minhas escolhas…</button></div>
    <div className="profile-tabs" role="group" aria-label="Origem dos perfis">{[["builtin", "Incluídos", included.length], ["mine", "Meus perfis", mine.length], ["imported", "Importados", imported.length]].map(([key, label, count]) => <button key={key} className={"profile-tab" + (tab === key ? " on" : "")} aria-pressed={tab === key} onClick={() => setTab(String(key))}>{label} <span>{count}</span></button>)}<button className="btn sm" onClick={async () => {try {const review = await api.importProfile(); if (review) ctx.openReview(review);} catch(e) {ctx.toast("Perfil recusado: " + errText(e), true);}}}><Upload size={14} /> Importar perfil…</button></div>
    {sections.map(([title, list]) => <section key={title}><div className="cat-title">{title}<span className="line" /></div>{list.length === 0 ? <div className="panel empty-profiles"><h2>{tab === "mine" ? "Ainda não guardou um perfil" : "Ainda não importou perfis"}</h2><p className="muted">{tab === "mine" ? "Faça as suas escolhas e guarde-as como perfil para voltar a usá-las." : "Importe um perfil partilhado consigo para o rever antes de usar."}</p><button className="btn" onClick={() => ctx.go("selection")}>Ver minhas escolhas</button></div> : <div className="grid profile-grid">{list.map(profile => {
      const ids = Object.keys(profile.selections || {});
      return <article key={profile.source + ":" + profile.id} className="card profile-card" aria-label={profile.name}>
        <div className="card-h"><span className="profile-icon"><Bookmark size={20} /></span><span className="badge b-unknown">{ids.length} {ids.length === 1 ? "item" : "itens"}</span></div>
        <h2 className="card-title">{profile.name}</h2><p className="profile-description">{profile.audience || profile.description}</p>
        <p className="profile-sample">{ids.slice(0,4).map(id => ctx.items[id]?.name || id).join(" · ")}{ids.length > 4 && ` · +${ids.length - 4}`}</p>
        <button className="btn primary profile-use" disabled={ctx.running || !!ctx.editing} aria-label={"Ver e usar " + profile.name} onClick={async () => {try {ctx.openReview(await api.reviewProfile(profile));} catch(e) {ctx.toast(errText(e), true);}}}>Ver e usar perfil…</button>
        <details className="profile-secondary"><summary>Mais opções</summary><div className="btnrow"><button className="btn sm" onClick={() => ctx.showCatalogProfile(`${profile.source}:${profile.id || profile.name}`)}>Explorar os itens</button><button className="btn sm" onClick={() => saveCopy(profile)}>Guardar uma cópia</button>{profile.source !== "builtin" && <button className="btn sm" disabled={!!ctx.editing || ctx.running} onClick={() => ctx.startEdit(profile)}>Editar escolhas</button>}<button className="btn sm" onClick={() => setExporting(profile)} aria-label={"Exportar " + profile.name}><Download size={13} /> Exportar</button>{profile.source !== "builtin" && <button className="btn sm danger" aria-label={"Apagar " + profile.name} onClick={async () => {if (!window.confirm("Apagar o perfil «" + profile.name + "»?")) return;try {await api.deleteProfile(profile.source!,profile.id!);await reload();} catch(e) {ctx.toast(errText(e),true);}}}><Trash2 size={13} /> Apagar</button>}</div></details>
      </article>;
    })}</div>}</section>)}
    <details className="catalog-secondary"><summary><FileText size={15} /> Importar a lista de aplicações de outro computador</summary><InventoryImport ctx={ctx} /></details>
    {saving && <SaveProfileDialog ctx={ctx} selections={ctx.sel} onClose={() => setSaving(false)} />}
    {exporting && <ExportDialog ctx={ctx} profile={{...exporting,id:exporting.id || slug(exporting.name)}} onClose={() => setExporting(null)} />}
  </>;
}

function slug(n: string) {
  const s = n.toLowerCase().normalize("NFD").replace(/[\u0300-\u036f]/g, "").replace(/[^a-z0-9]+/g, "-").replace(/^-+|-+$/g, "").slice(0, 60);
  return s.length >= 2 ? s : "perfil-" + s;
}

const CLASS_LABEL: Record<string, string> = { app: "Aplicação no catálogo", unknown: "Desconhecido", legacy: "Legado", dependency: "Dependência", system: "Sistema" };

function InventoryImport({ ctx }: { ctx: Ctx }) {
  const [text, setText] = useState("");
  const [rows, setRows] = useState<InvEntry[] | null>(null);
  const [pick, setPick] = useState<Record<string, boolean>>({});
  const [filter, setFilter] = useState("app");
  const parse = async (t: string) => { setRows((await api.parseInventory(t)) || []); setPick({}); };
  const counts = (rows || []).reduce<Record<string, number>>((m, r) => { m[r.class] = (m[r.class] || 0) + 1; return m; }, {});
  return (
    <div className="panel" style={{ marginTop: 16 }}>
      <h2><FileText size={17} /> Importar inventário do PC antigo</h2>
      <p className="muted" style={{ marginTop: 0 }}>Cole a saída de <span className="mono">apt-mark showmanual</span>, <span className="mono">snap list</span> e/ou <span className="mono">flatpak list --app --columns=application</span>. A lista é uma fonte de candidatos: nada fica selecionado automaticamente; kernels, GRUB, shim, libc e ICU antigas nunca são transportados.</p>
      <textarea className="inv" value={text} onChange={(e) => setText(e.target.value)} placeholder={"git\nlinux-image-6.8.0-45-generic\nfirefox  131.0  4793  latest/stable  mozilla✓  -\norg.freefilesync.FreeFileSync"} aria-label="Inventário" />
      <div style={{ display: "flex", gap: 8, marginTop: 8 }}>
        <button className="btn" onClick={async () => { try { const t = await api.openInventoryFile(); if (t) { setText(t); parse(t); } } catch (e) { ctx.toast(errText(e), true); } }}><Upload size={15} /> Abrir ficheiro</button>
        <button className="btn primary" disabled={!text.trim()} onClick={() => parse(text)}>Classificar</button>
      </div>
      {rows && (
        <>
          <div className="chips" style={{ marginTop: 12 }}>
            {Object.keys(CLASS_LABEL).map((k) => <button key={k} className={"chip" + (filter === k ? " on" : "")} onClick={() => setFilter(k)}>{CLASS_LABEL[k]}<span className="n">{counts[k] || 0}</span></button>)}
          </div>
          <table className="t">
            <thead><tr><th></th><th>Nome</th><th>Origem</th><th>Classificação</th><th>Item do catálogo</th></tr></thead>
            <tbody>
              {rows.filter((r) => r.class === filter).slice(0, 400).map((r) => (
                <tr key={r.source + r.name}>
                  <td>{r.item ? <input type="checkbox" checked={!!pick[r.item]} onChange={(e) => setPick({ ...pick, [r.item!]: e.target.checked })} aria-label={"Escolher " + r.name} /> : null}</td>
                  <td className="mono">{r.name}</td><td>{r.source}</td><td className="muted">{r.reason}</td>
                  <td>{r.item ? <>{ctx.items[r.item]?.name} <StatusBadge st={ctx.states[r.item]} method={ctx.items[r.item]?.method} adminOp={ctx.items[r.item]?.adminOp} /></> : "—"}</td>
                </tr>
              ))}
            </tbody>
          </table>
          <button className="btn primary" style={{ marginTop: 10 }} disabled={!Object.values(pick).some(Boolean)} onClick={() => {
            const next = { ...ctx.sel }; const org = { ...ctx.origins };
            Object.entries(pick).forEach(([id, v]) => { if (v && !next[id]) { next[id] = {}; org[id] = ["perfil:inventário do PC antigo"]; } });
            ctx.setSel(next, org);
            ctx.toast("Itens escolhidos adicionados à seleção.");
          }}>Adicionar escolhidos à seleção</button>
        </>
      )}
    </div>
  );
}

/* ------------------------------ Diagnóstico ------------------------------ */
export function DiagnosticsPage({ ctx }: { ctx: Ctx }) {
  const [rep, setRep] = useState<DiagReport | null>(null);
  const [exportOpen, setExportOpen] = useState(false);
  const [format, setFormat] = useState<"pdf" | "csv" | "xlsx">("pdf");
  const [sections, setSections] = useState<Record<string, boolean>>({facts: true, disks: true, startup: true, suggestions: true});
  const [exporting, setExporting] = useState(false);
  const [exportError, setExportError] = useState("");
  const exportReport = async () => {
    if (!rep || exporting) return;
    setExporting(true); setExportError("");
    try {
      const path = await api.exportDiagnostics(rep.snapshotId, format, Object.keys(sections).filter(s => sections[s]));
      if (path) ctx.toast("Diagnóstico exportado para " + path);
    } catch (e) { setExportError(errText(e)); } finally { setExporting(false); }
  };
  const load = () => { setRep(null); api.diagnostics().then(setRep).catch((e) => ctx.toast(errText(e), true)); };
  useEffect(load, []);
  return (
    <>
      <div className="page-head">
        <div style={{ flex: 1 }}><h1>Diagnóstico</h1><p>Leitura apenas — nada é alterado. As sugestões mostram a evidência que as motivou; nenhuma é aplicada sem a selecionar e rever o plano.</p></div>
        <button className="btn" disabled={!rep || exporting} onClick={() => setExportOpen(!exportOpen)}><Download size={15} /> Exportar relatório…</button>
        <button className="btn" disabled={exporting} onClick={load}><RefreshCw size={15} /> Atualizar</button>
      </div>
      {!rep ? <div className="empty">A ler o sistema…</div> : (
        <>
          {exportOpen && <section className="panel" aria-label="Exportar diagnóstico">
            <h2>Exportar diagnóstico</h2>
            <p>Recolhido em {rep.collectedAt ? new Date(rep.collectedAt).toLocaleString("pt-PT") : "agora"}. Guarda uma cópia dos dados apresentados, com versão, sistema e modo. Não é um backup do computador.</p>
            <label className="field">Formato<select disabled={exporting} value={format} onChange={e => setFormat(e.target.value as "pdf" | "csv" | "xlsx")}><option value="pdf">PDF — relatório para leitura</option><option value="csv">CSV — tabela de dados</option><option value="xlsx">Excel (.xlsx) — tabela com filtros</option></select></label>
            <fieldset className="diag-export-sections"><legend>Secções incluídas</legend>{([
              ["facts", "Sistema", rep.facts.length], ["disks", "Discos", rep.disks?.length || 0], ["startup", "Entradas de arranque", rep.startup?.length || 0], ["suggestions", "Sugestões", rep.suggestions?.length || 0],
            ] as [string, string, number][]).map(([id, label, count]) => <label className="confirm-row" key={id}><input type="checkbox" disabled={exporting} checked={sections[id]} onChange={e => setSections({...sections, [id]: e.target.checked})} /> {label} ({count})</label>)}</fieldset>
            {sections.startup && <p className="muted">Entradas de arranque incluem os caminhos de origem e o executável. Reveja os dados abaixo antes de partilhar o ficheiro.</p>}
            {exportError && <div className="note bad" role="alert">{exportError}</div>}
            <button className="btn primary" disabled={exporting || !Object.values(sections).some(Boolean)} onClick={exportReport}><Download size={15} /> {exporting ? "A guardar…" : "Escolher ficheiro e guardar"}</button>
          </section>}
          <div className="panel">
            <h2><Lightbulb size={17} /> Sugestões contextualizadas</h2>
            {(rep.suggestions || []).length === 0 && <p className="muted">Sem sugestões: nada de anómalo detetado.</p>}
            {(rep.suggestions || []).map((s, i) => (
              <div key={i} className={"note " + (s.level === "warn" ? "warn" : s.level === "suggest" ? "info" : "ok")}>
                {s.level === "warn" ? <AlertTriangle size={16} /> : <Lightbulb size={16} />}
                <div style={{ flex: 1 }}><b>{s.title}</b> — {s.detail}<div className="muted">Evidência: {s.evidence}</div></div>
                {s.item && ctx.items[s.item] && (
                  ctx.sel[s.item] ? <span className="tag">selecionado</span> :
                    <button className="btn sm" onClick={() => { ctx.toggle(s.item!); ctx.toast(ctx.items[s.item!].name + " selecionado"); }}>Selecionar</button>
                )}
              </div>
            ))}
            <p className="muted" style={{ marginBottom: 0 }}>Swappiness, zram, escalonador de I/O, sysctl, DNS e limites globais não são alterados sem evidência específica. Nenhum ganho de desempenho é afirmado sem medição.</p>
          </div>
          <div className="panel">
            <h2><Stethoscope size={17} /> Sistema</h2>
            <div className="facts">{rep.facts.map((f) => <div className="fact" key={f.label}><small>{f.label}</small><div>{f.value}</div></div>)}</div>
          </div>
          <div className="cols">
            <div className="panel">
              <h2><HardDrive size={17} /> Discos</h2>
              <table className="t"><thead><tr><th>Disco</th><th>Tipo</th><th>TRIM</th><th>Livre</th></tr></thead><tbody>
                {(rep.disks || []).map((d) => <tr key={d.name}><td>{d.name} <span className="muted">{d.model} · {d.size}</span></td><td>{d.ssd ? "SSD" : "HDD"}</td><td>{d.discard ? "suporta" : "não"}</td><td>{d.mount ? `${d.freePct}% em ${d.mount} (${d.free})` : "—"}</td></tr>)}
              </tbody></table>
            </div>
            <div className="panel">
              <h2>Entradas de arranque</h2>
              <p className="muted">São ficheiros de configuração. Uma entrada pode existir mesmo depois de a aplicação ter sido removida.</p>
              {(rep.startup || []).length === 0 ? <p className="muted">Nenhuma entrada encontrada.</p> : <ul className="startup-list">{(rep.startup || []).map((s) => <li key={s.path}>
                <div><b>{s.name}</b> <span className="muted">({s.scope})</span></div>
                <span className={"badge " + (s.status === "available" ? "b-installed" : s.status === "missing" ? "b-config" : "b-unknown")}>{({available: "Comando disponível", missing: "Executável não encontrado", disabled: "Desativada nesta sessão", unknown: "Por verificar"})[s.status]}</span>
                <p className="muted">{s.detail}</p>
                <details className="technical-details"><summary>Ver origem da entrada</summary><p className="mono">{s.path}</p>{s.command && <p className="mono">Executável: {s.command}</p>}</details>
              </li>)}</ul>}
              <p className="muted">O diagnóstico não inicia nem remove estas entradas. Reveja a origem antes de desativar algo que não reconhece.</p>
            </div>
          </div>
        </>
      )}
    </>
  );
}

/* ------------------------------ Migração ------------------------------ */
const CHECKLIST: [string, string, string][] = [
  ["keyring", "Chaveiro (GNOME Keyring / Senhas e chaves)", "Exporte do PC antigo com «Senhas e chaves» ou copie ~/.local/share/keyrings manualmente com a sessão terminada. O Meu Ubuntu não lê nem copia segredos."],
  ["browsers", "Navegadores", "Ative a sincronização (Firefox Account / Brave Sync) ou use a migração de perfil abaixo. Senhas: exporte/importe pelo próprio navegador."],
  ["ssh", "Chaves SSH", "Copie ~/.ssh por meio seguro (pen cifrada/scp) e aplique chmod 700 ~/.ssh e 600 às chaves privadas. As chaves existentes nunca são substituídas."],
  ["gpg", "Chaves GPG", "gpg --export-secret-keys --armor > ficheiro (guarde cifrado) e gpg --import no novo PC; confirme a confiança."],
  ["vpn", "VPN e certificados", "Reimporte os .ovpn no NetworkManager e os certificados de cliente; não são copiados automaticamente."],
  ["docker", "Volumes Docker", "Faça dump das bases de dados (mysqldump) ou docker run --rm -v vol:/v -v $PWD:/b alpine tar czf /b/vol.tgz -C /v . — não copie /var/lib/docker inteiro."],
  ["projects", "Projetos e .env", "Copie os projetos; os .env contêm segredos: transfira-os por meio seguro e nunca para perfis."],
  ["home", "Pasta pessoal", "Restaure documentos seletivamente. Não restaure a home inteira por cima das configurações novas (~/.config)."],
];

export function MigrationPage({ ctx }: { ctx: Ctx }) {
  const [ff, setFf] = useState<{ roots: FFRoot[] | null; running: boolean } | null>(null);
  const [busy, setBusy] = useState("");
  const [checks, setChecks] = useState<Record<string, boolean>>(() => { try { return JSON.parse(localStorage.getItem("mu-checklist") || "{}"); } catch { return {}; } });
  const load = () => api.firefox().then(setFf).catch((e) => ctx.toast(errText(e), true));
  useEffect(() => { load(); }, []);
  useEffect(() => { try { localStorage.setItem("mu-checklist", JSON.stringify(checks)); } catch { /* ignore */ } }, [checks]);
  const roots = ff?.roots || [];
  const flatpakHasProfile = roots.some((r) => r.kind === "flatpak" && (r.profiles || []).length > 0);
  const act = async (kind: "backup" | "migrate", path: string) => {
    setBusy(kind + path);
    try {
      const out = kind === "backup" ? await api.firefoxBackup(path) : await api.firefoxMigrate(path);
      ctx.toast((kind === "backup" ? "Backup criado: " : "Perfil copiado para: ") + out);
      load();
    } catch (e) { ctx.toast(errText(e), true); } finally { setBusy(""); }
  };
  return (
    <>
      <div className="page-head"><div style={{ flex: 1 }}><h1>Migração</h1><p>Perfis do Firefox, checklist de dados pessoais e remoção do Firefox Snap — passos separados e reversíveis, nunca automáticos.</p></div>
        <button className="btn" onClick={load}><RefreshCw size={15} /> Detetar</button></div>
      <div className="panel">
        <h2><Flame size={17} /> Firefox: perfis detetados</h2>
        {ff?.running && <div className="note warn"><AlertTriangle size={16} /> O Firefox está aberto. Feche todas as janelas (Snap e Flatpak) antes de backup ou migração.</div>}
        {roots.length === 0 && <p className="muted">Nenhum perfil encontrado em ~/snap/firefox, ~/.mozilla, ~/.config/mozilla ou ~/.var/app/org.mozilla.firefox.</p>}
        {roots.map((r) => (
          <div key={r.root} style={{ marginBottom: 10 }}>
            <h3>{r.label} <span className="muted mono" style={{ fontWeight: 400 }}>{r.root}</span></h3>
            <table className="t"><tbody>
              {(r.profiles || []).map((p) => (
                <tr key={p.path}><td><b>{p.name}</b><div className="muted mono">{p.path}</div></td><td>{p.size}</td>
                  <td style={{ textAlign: "right", whiteSpace: "nowrap" }}>
                    <button className="btn sm" disabled={!!busy || ff?.running} onClick={() => act("backup", p.path)}><Archive size={13} /> Backup</button>{" "}
                    {r.kind !== "flatpak" && <button className="btn sm" disabled={!!busy || ff?.running} onClick={() => act("migrate", p.path)}><ArrowRightLeft size={13} /> Copiar para Flatpak</button>}
                  </td></tr>
              ))}
            </tbody></table>
          </div>
        ))}
        <p className="muted">O backup é um .tar.gz em ~/MeuUbuntu-backups (só leitura do perfil). A cópia cria um <b>novo</b> perfil «Migrado» no Firefox Flatpak — nunca sobrescreve um existente — e faz backup do profiles.ini. Escolha-o em about:profiles. O Brave não importa automaticamente perfis do Firefox Flatpak/Snap: use Sync ou exporte marcadores (HTML) e senhas pelo próprio navegador.</p>
        <div className="note info"><Info size={16} /><div style={{ flex: 1 }}>Remover o Firefox Snap é uma ação independente: faça-o só depois de validar o novo perfil. O snapd não é removido.
          {flatpakHasProfile ? "" : " (Ainda não há perfil no Firefox Flatpak.)"}</div>
          {ctx.items["firefox-snap-remove"] && <button className="btn sm" disabled={!flatpakHasProfile || !!ctx.sel["firefox-snap-remove"]} onClick={() => { ctx.toggle("firefox-snap-remove"); ctx.toast("«Remover Firefox Snap» adicionado; o plano pede confirmação explícita de que validou o perfil migrado."); }}>{ctx.sel["firefox-snap-remove"] ? "Selecionado" : "Selecionar remoção"}</button>}
        </div>
      </div>
      <div className="panel">
        <h2><ListChecks size={17} /> Checklist de migração</h2>
        <div className="checklist">
          {CHECKLIST.map(([id, t, d]) => (
            <label key={id}><input type="checkbox" checked={!!checks[id]} onChange={(e) => setChecks({ ...checks, [id]: e.target.checked })} /><div><b>{t}</b><small>{d}</small></div></label>
          ))}
        </div>
      </div>
    </>
  );
}

/* ------------------------------ Histórico ------------------------------ */
const RUN_STATUS: Record<string, string> = { done: "Concluída", intervention: "Requer intervenção", failed: "Com falhas", cancelled: "Cancelada", interrupted: "Interrompida", running: "Em curso" };

export function HistoryPage({ ctx }: { ctx: Ctx }) {
  const [reviewing, setReviewing] = useState<Change | null>(null);
  const [runs, setRuns] = useState<Run[]>([]);
  const [changes, setChanges] = useState<Change[]>([]);
  const [log, setLog] = useState<{ id: string; text: string } | null>(null);
  const load = () => {
    api.runs().then((r) => setRuns(r || [])).catch((e) => ctx.toast(errText(e), true));
    api.changes().then((c) => setChanges((c || []).slice().reverse())).catch((e) => ctx.toast("Não foi possível ler o histórico: " + errText(e), true));
  };
  useEffect(load, []);
  return (
    <>
      <div className="page-head"><div style={{ flex: 1 }}><h1>Histórico e reversão</h1><p>Cada preferência alterada fica registada com o valor anterior. A reversão é específica por alteração e recusa-se se o valor mudou depois. Instalações de pacotes não têm rollback automático.</p></div>
        <button className="btn" onClick={load}><RefreshCw size={15} /> Atualizar</button></div>
      <div className="panel">
        <h2><Undo2 size={17} /> Alterações registadas</h2>
        {changes.length === 0 ? <p className="muted">Ainda não há alterações.</p> : (
          <table className="t"><thead><tr><th>Quando</th><th>Ajuste</th><th>Antes → depois</th><th></th></tr></thead><tbody>
            {changes.map((c) => (
              <tr key={c.id}>
                <td className="muted">{new Date(c.time).toLocaleString("pt-PT")}</td>
                <td><b>{c.title}</b><div className="muted mono">{c.schema ? c.schema + " " : ""}{c.key}</div>{c.status === "pending" && <small className="st-intervention">Aplicação sem confirmação final; verifique o estado atual antes de reverter.</small>}</td>
                <td className="mono">{c.kind === "template" ? "(criado)" : c.kind === "admin" ? (c.after || "aplicado") : <>{c.before || "(sem valor)"} → {c.after}</>}</td>
                <td style={{ textAlign: "right" }}>{c.reverted ? <span className="badge b-installed"><CheckCircle2 size={11} /> revertido</span> : <button className="btn sm" onClick={() => setReviewing(c)}><Undo2 size={13} /> Reverter…</button>}</td>
              </tr>
            ))}
          </tbody></table>
        )}
      </div>
      <div className="panel">
        <h2>Execuções</h2>
        {runs.length === 0 ? <p className="muted">Ainda não há execuções.</p> : (
          <table className="t"><thead><tr><th>Início</th><th>Estado</th><th>Passos</th><th></th></tr></thead><tbody>
            {runs.map((r) => (
              <tr key={r.id}>
                <td>{new Date(r.started).toLocaleString("pt-PT")}{r.simulated && <span className="tag" style={{ marginLeft: 6 }}>simulação</span>}</td>
                <td>{RUN_STATUS[r.status] || r.status}</td>
                <td className="muted">{r.steps.filter((s) => s.status === "done").length} concluídos · {r.steps.filter((s) => s.status === "failed").length} falharam · {r.steps.filter((s) => s.status === "intervention").length} requerem intervenção · {r.steps.length} total{r.steps.filter((s) => s.note || s.error).map(s => <div key={s.id}>{s.title}: {s.note || s.error}</div>)}</td>
                <td style={{ textAlign: "right" }}><button className="btn sm" onClick={async () => setLog({ id: r.id, text: await api.runLog(r.id) })}>Ver registo</button></td>
              </tr>
            ))}
          </tbody></table>
        )}
        {log && <><h3>Registo {log.id}</h3><pre className="log" style={{ maxHeight: 420 }}>{log.text || "(vazio)"}</pre></>}
      </div>
      {reviewing && <RevertDialog ctx={ctx} change={reviewing} onClose={() => setReviewing(null)} onDone={() => { load(); ctx.refresh(); }} />}
    </>
  );
}
