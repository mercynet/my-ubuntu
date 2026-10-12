import { useEffect, useState } from "react";
import { Copy, FolderOpen, PackagePlus, RefreshCw } from "lucide-react";
import { api, errText, type IntegrationPackage, type IntegrationReport, type IntegrationReview } from "./api";
import { Modal } from "./Modals";

export function IntegrationsPage({ selected, packages, onPackages, onSelect, toast }: { selected: string; packages: IntegrationPackage[]; onPackages: (p: IntegrationPackage[]) => void; onSelect: (id: string) => void; toast: (message: string, error?: boolean) => void }) {
  const [review, setReview] = useState<IntegrationReview | null>(null);
  const [confirmed, setConfirmed] = useState(false);
  const [busy, setBusy] = useState("");
  const [report, setReport] = useState<IntegrationReport | null>(null);
  const [error, setError] = useState("");
  const pack = packages.find(p => p.id === selected);
  useEffect(() => { setReport(null); setError(""); }, [selected]);
  const act = async (name: string, fn: () => Promise<void>) => {
    setBusy(name); setError("");
    try { await fn(); } catch (e) { setError(errText(e)); } finally { setBusy(""); }
  };
  const copy = async (text: string) => {
    try { if (!await window.runtime?.ClipboardSetText?.(text)) throw new Error("Não foi possível copiar"); toast("Copiado. Reveja antes de executar no terminal."); } catch (e) { toast(errText(e), true); }
  };
  return <>
    <div className="page-head"><div style={{flex: 1}}><h1>{pack?.name || "Pacotes de integração"}</h1><p>{pack?.description || "Acrescente um separador e um guia de preparação fornecidos pela sua equipa."}</p></div>
      <button className="btn" disabled={!!busy} onClick={() => act("import", async () => { setConfirmed(false); setReview(await api.previewIntegration()); })}><PackagePlus size={16} /> Instalar pacote…</button></div>
    {error && <div className="note bad" role="alert">{error}</div>}
    {!pack && <div className="panel"><h2>Integrações instaladas</h2>{packages.length ? packages.map(p => <button className="btn" key={p.id} onClick={() => onSelect(p.id)}>{p.name} · {p.version}</button>) : <p>Escolha o ZIP de integração recebido da equipa. O pacote acrescenta verificações e instruções; instalar o pacote não prepara o projeto automaticamente.</p>}</div>}
    {pack && <>
      <div className="panel"><h2>Projeto local</h2><p className="muted">Pacote {pack.version} · {pack.publisher}. Escolha a pasta onde já clonou o projeto.</p><p className="mono">{pack.project || "Pasta ainda não escolhida"}</p>
        <button className="btn" disabled={!!busy} onClick={() => act("folder", async () => { await api.integrationProject(pack.id); onPackages(await api.integrations()); setReport(null); })}><FolderOpen size={16} /> Escolher pasta…</button>{" "}
        <button className="btn primary" disabled={!!busy || !pack.project} onClick={() => act("diagnose", async () => setReport(await api.integrationDiagnostics(pack.id)))}><RefreshCw size={16} /> {busy === "diagnose" ? "A verificar…" : "Verificar ambiente"}</button>
        <p className="muted">As verificações consultam o ambiente. Os comandos do guia são copiados para execução manual; nada é instalado ou iniciado por este separador.</p>
      </div>
      {report && <div className="panel"><h2>Diagnóstico {report.simulated ? "de exemplo" : "do ambiente"}</h2><p className="muted">{new Date(report.collectedAt).toLocaleString("pt-PT")}</p>{report.results.map((r, i) => <div key={i} className={"note " + (r.status === "ok" ? "ok" : r.status === "example" ? "info" : "warn")}><div><b>{r.title}</b><div>{r.detail}</div></div></div>)}
        <button className="btn" onClick={() => copy(`${pack.name} ${pack.version}\n${report.collectedAt}\n${report.simulated ? "Simulação" : "Leitura real"}\n` + report.results.map(r => `${r.title}: ${r.status} — ${r.detail}`).join("\n"))}><Copy size={15} /> Copiar diagnóstico</button></div>}
      <div className="panel"><h2>Preparação assistida</h2><p>Conclua os passos pela ordem indicada. Confirme os problemas do diagnóstico antes de iniciar os serviços. Este primeiro pacote não cria tenants nem importa bases de dados.</p>{pack.guide.map((g, i) => <div className="panel" key={i}><h3>{i + 1}. {g.title}</h3><p>{g.text}</p>{g.command && <><pre className="code">{g.command}</pre><button className="btn sm" onClick={() => copy(g.command!)}><Copy size={14} /> Copiar comando</button></>}</div>)}</div>
    </>}
    {review && <Modal title="Rever pacote de integração" onClose={() => setReview(null)} closeDisabled={!!busy} footer={<><button className="btn" disabled={!!busy} onClick={() => setReview(null)}>Cancelar</button><button className="btn primary" disabled={!confirmed || !!busy} onClick={() => act("install", async () => { await api.installIntegration(review.token); onPackages(await api.integrations()); onSelect(review.manifest.id); setReview(null); toast("Pacote instalado; o novo separador está disponível."); })}>{busy ? "A instalar…" : "Instalar pacote"}</button></>}>
      <h3>{review.manifest.name} · {review.manifest.version}</h3><p>{review.manifest.description}</p><p>Fornecido por: {review.manifest.publisher}</p>{review.existing && <p>Atualiza o pacote {review.existing}; a pasta do projeto é preservada.</p>}
      <p>Inclui {review.manifest.checks.length} verificações e {review.manifest.guide.length} passos de guia. O pacote contém dados; os comandos do guia são executados apenas quando os copia e usa no terminal.</p>
      <label className="confirm-row"><input type="checkbox" checked={confirmed} onChange={e => setConfirmed(e.target.checked)} /> Confirmei que este pacote veio da minha equipa.</label>
    </Modal>}
  </>;
}
