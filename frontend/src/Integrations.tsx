import { useEffect, useState } from "react";
import { Copy, FolderOpen, PackagePlus, RefreshCw, Loader2 } from "lucide-react";
import { api, errText, type IntegrationPackage, type IntegrationReport, type IntegrationReview, type SetupReview, type SetupJob } from "./api";
import { Modal } from "./Modals";
import { setupLabels } from "./IntegrationSetup";
import type { Ctx } from "./App";

export function IntegrationsPage({ ctx, selected, packages, onPackages, onSelect, job, onJob, toast }: { ctx: Ctx; selected: string; packages: IntegrationPackage[]; onPackages: (p: IntegrationPackage[]) => void; onSelect: (id: string) => void; job: SetupJob | null; onJob: (job: SetupJob) => void; toast: (message: string, error?: boolean) => void }) {
  const [review, setReview] = useState<IntegrationReview | null>(null);
  const [confirmed, setConfirmed] = useState(false);
  const [busy, setBusy] = useState("");
  const [report, setReport] = useState<IntegrationReport | null>(null);
  const [error, setError] = useState("");
  const [setupReview, setSetupReview] = useState<SetupReview | null>(null);
  const [setupConfirmed, setSetupConfirmed] = useState(false);
  const pack = packages.find(p => p.id === selected);
  const blocked = !!busy || ctx.running;
  const packJob = job?.package === pack?.id ? job : null;
  useEffect(() => { setReport(null); setError(""); setSetupReview(null); }, [selected]);
  const act = async (name: string, fn: () => Promise<void>) => {
    setBusy(name); setError("");
    try { await fn(); } catch (e) { setError(errText(e)); } finally { setBusy(""); }
  };
  const copy = async (text: string) => {
    try { if (!await window.runtime?.ClipboardSetText?.(text)) throw new Error("Não foi possível copiar"); toast("Copiado. Reveja antes de executar no terminal."); } catch (e) { toast(errText(e), true); }
  };
  const tools = () => {
    if (!pack?.setup) return;
    const selections = { ...ctx.sel };
    for (const id of pack.setup.tools) {
      if (id === "docker-engine" && ctx.states["docker-ubuntu"]?.status === "installed") continue;
      if (ctx.items[id]) selections[id] = selections[id] || {};
    }
    ctx.setSel(selections);
    ctx.reviewPlan(selections);
  };
  return <>
    <div className="page-head"><div style={{flex: 1}}><h1>{pack?.name || "Pacotes de integração"}</h1><p>{pack?.description || "Acrescente um separador e um guia de preparação fornecidos pela sua equipa."}</p></div>
      <button className="btn" disabled={blocked} onClick={() => act("import", async () => { setConfirmed(false); setReview(await api.previewIntegration()); })}><PackagePlus size={16} /> Instalar pacote…</button></div>
    {error && <div className="note bad" role="alert">{error}</div>}
    {!pack && <div className="panel"><h2>Integrações instaladas</h2>{packages.length ? packages.map(p => <button className="btn" key={p.id} onClick={() => onSelect(p.id)}>{p.name} · {p.version}</button>) : <p>Escolha o ZIP de integração recebido da equipa. O pacote acrescenta verificações e instruções; instalar o pacote não prepara o projeto automaticamente.</p>}</div>}
    {pack && <>
      <div className="panel"><h2>Projeto local</h2><p className="muted">Pacote {pack.version} · {pack.publisher}. Escolha a pasta onde já clonou o projeto.</p><p className="mono">{pack.project || "Pasta ainda não escolhida"}</p>
        <button className="btn" disabled={blocked} onClick={() => act("folder", async () => { await api.integrationProject(pack.id); onPackages(await api.integrations()); setReport(null); })}><FolderOpen size={16} /> Escolher pasta…</button>{" "}
        <button className="btn" disabled={blocked || !pack.project} onClick={() => act("diagnose", async () => setReport(await api.integrationDiagnostics(pack.id)))}><RefreshCw size={16} /> {busy === "diagnose" ? "A verificar…" : "Verificar ambiente"}</button>
        <p className="muted">Verificar ambiente só consulta. {pack.setup ? "Preparar ambiente mostra uma revisão e pede confirmação antes de criar ficheiros ou iniciar serviços." : "Os comandos do guia são copiados para execução manual."}</p>
      </div>
      {pack.setup && <div className="panel"><h2>Preparar ambiente</h2><p>Prepara containers isolados, schemas de desenvolvimento, Redis Cluster e dependências. Os dumps são importados por si; passos administrativos aparecem no final.</p>
        <button className="btn" disabled={blocked} onClick={tools}>1. Rever ferramentas necessárias</button>{" "}
        <button className="btn primary" disabled={blocked || !pack.project} onClick={() => act("setup-review", async () => { setSetupConfirmed(false); setSetupReview(await api.integrationSetupReview(pack.id)); })}>{busy === "setup-review" && <Loader2 size={16} className="spin" />}2. Rever preparação do projeto</button>
        <p className="muted">Se instalou Docker ou mudou as permissões, termine sessão e entre novamente antes da preparação. Use um checkout e configuração de desenvolvimento fornecidos pela equipa.</p>
      </div>}
      {packJob && <div className="panel"><h2>Última preparação · {setupLabels[packJob.state] || packJob.state}{packJob.simulated ? " · Simulação" : ""}</h2>
        <ol className="setup-steps">{packJob.steps.map((s, i) => <li key={i} className={s.status}>{s.status === "running" && <Loader2 className="spin" size={15} />}<b>{s.title}</b> · {setupLabels[s.status] || s.status}{s.message && <p>{s.message}</p>}</li>)}</ol>
        {packJob.message && <div className={"note " + (packJob.state === "done" ? "ok" : "warn")}>{packJob.message}</div>}
        {packJob.finished && <><p>Para repetir após uma falha ou interrupção, reveja novamente. Ficheiros existentes e dados são conferidos; nenhuma base é reposta.</p>
          {packJob.address && <p>Endereço após importar os dumps: <span className="mono">{packJob.address}</span></p>}
          {packJob.terminal && <><p>Use este comando para gerir apenas o ambiente preparado:</p><pre className="code">{packJob.terminal}</pre><button className="btn sm" onClick={() => copy("'" + packJob.terminal.replace(/'/g, "'\\''") + "' ps")}><Copy size={14} /> Copiar comando de estado</button><p>Credenciais locais em <span className="mono">settings.json</span>, junto deste comando. Não partilhe esse ficheiro.</p></>}
          {packJob.manual.map((m, i) => <div className="panel" key={i}><p>{m}</p><button className="btn sm" onClick={() => copy(m)}><Copy size={14} /> Copiar instrução</button></div>)}
        </>}
      </div>}
      {report && <div className="panel"><h2>Diagnóstico {report.simulated ? "de exemplo" : "do ambiente"}</h2><p className="muted">{new Date(report.collectedAt).toLocaleString("pt-PT")}</p>{report.results.map((r, i) => <div key={i} className={"note " + (r.status === "ok" ? "ok" : r.status === "example" ? "info" : "warn")}><div><b>{r.title}</b><div>{r.detail}</div></div></div>)}
        <button className="btn" onClick={() => copy(`${pack.name} ${pack.version}\n${report.collectedAt}\n${report.simulated ? "Simulação" : "Leitura real"}\n` + report.results.map(r => `${r.title}: ${r.status} — ${r.detail}`).join("\n"))}><Copy size={15} /> Copiar diagnóstico</button></div>}
      <div className="panel"><h2>Guia da equipa</h2><p>{pack.setup ? "Consulte estas instruções para requisitos e passos que exigem intervenção." : "Conclua os passos pela ordem indicada; este pacote não prepara o projeto automaticamente."}</p>{pack.guide.map((g, i) => <div className="panel" key={i}><h3>{i + 1}. {g.title}</h3><p>{g.text}</p>{g.command && <><pre className="code">{g.command}</pre><button className="btn sm" onClick={() => copy(g.command!)}><Copy size={14} /> Copiar comando</button></>}</div>)}</div>
    </>}
    {review && <Modal title="Rever pacote de integração" onClose={() => setReview(null)} closeDisabled={!!busy} footer={<><button className="btn" disabled={!!busy} onClick={() => setReview(null)}>Cancelar</button><button className="btn primary" disabled={!confirmed || !!busy} onClick={() => act("install", async () => { await api.installIntegration(review.token); onPackages(await api.integrations()); onSelect(review.manifest.id); setReview(null); toast("Pacote instalado; o novo separador está disponível."); })}>{busy ? "A instalar…" : "Instalar pacote"}</button></>}>
      <h3>{review.manifest.name} · {review.manifest.version}</h3><p>{review.manifest.description}</p><p>Fornecido por: {review.manifest.publisher}</p>{review.existing && <p>Atualiza o pacote {review.existing}; a pasta do projeto é preservada.</p>}
      <p>Inclui {review.manifest.checks.length} verificações e {review.manifest.guide.length} passos de guia. O pacote contém dados; os comandos do guia são executados apenas quando os copia e usa no terminal.</p>
      <label className="confirm-row"><input type="checkbox" checked={confirmed} onChange={e => setConfirmed(e.target.checked)} /> Confirmei que este pacote veio da minha equipa.</label>
    </Modal>}
    {setupReview && pack && <Modal title="Rever preparação do ambiente" onClose={() => setSetupReview(null)} closeDisabled={!!busy} footer={<><button className="btn" disabled={!!busy} onClick={() => setSetupReview(null)}>Cancelar</button><button className="btn primary" disabled={!setupConfirmed || blocked} onClick={() => act("setup-start", async () => { onJob(await api.integrationSetupStart(pack.id, setupReview.token)); setSetupReview(null); })}>{busy === "setup-start" ? "A iniciar…" : "Preparar ambiente"}</button></>}>
      <p className="mono">{setupReview.project}</p>{setupReview.simulated && <div className="note info">Simulação: nenhum ficheiro do projeto, container ou base será alterado.</div>}
      <ol>{setupReview.steps.map(s => <li key={s}>{s}</li>)}</ol>
      <p>Serviços: {setupReview.services.join(", ")}</p><p>Endereço: {setupReview.address || "Exemplo controlado"}</p>
      <p>Publicação de portas apenas neste computador. Ficheiros Compose originais e serviços existentes são preservados. Novas configurações de tenants são criadas nas pastas indicadas pelo projeto.</p>
      {setupReview.databases.map(d => <p key={d.service}><b>{d.service}:</b> {d.schemas.join(", ")}</p>)}
      <p>O build e Composer executam código do checkout escolhido. Bases e Redis existentes não são apagados nem repostos. Os dumps e comandos administrativos continuam manuais.</p>
      <label className="confirm-row"><input type="checkbox" checked={setupConfirmed} onChange={e => setSetupConfirmed(e.target.checked)} /> Confio neste checkout de desenvolvimento e revi os serviços e alterações.</label>
    </Modal>}
  </>;
}
