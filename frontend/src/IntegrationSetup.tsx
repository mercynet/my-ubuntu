import { useEffect, useState } from "react";
import { Loader2, Pause, Play, X } from "lucide-react";
import { api, errText, type SetupJob } from "./api";

export const setupLabels: Record<string, string> = {
  pending: "Na fila", running: "Em execução", pausing: "Vai pausar após a etapa atual",
  paused: "Em pausa", cancelling: "Vai parar após a etapa atual", done: "Concluído",
  failed: "Falhou", intervention: "Requer atenção", cancelled: "Cancelado", interrupted: "Interrompido",
};

export function SetupProgress({ job, onOpen, toast }: { job: SetupJob | null; onOpen: () => void; toast: (text: string, error?: boolean) => void }) {
  const [now, setNow] = useState(Date.now());
  const [busy, setBusy] = useState(false);
  useEffect(() => {
    if (!job || job.finished) return;
    setNow(Date.now());
    const timer = window.setInterval(() => setNow(Date.now()), 1000);
    return () => window.clearInterval(timer);
  }, [job?.id, job?.finished]);
  if (!job || job.finished) return null;
  const done = job.steps.filter(s => ["done", "intervention", "failed", "cancelled"].includes(s.status)).length;
  const current = job.steps.find(s => s.status === "running");
  const elapsed = Math.max(0, Math.floor((now - Date.parse(job.started)) / 1000));
  const paused = job.state === "paused" || job.state === "pausing";
  const control = async (action: "cancel" | "pause" | "resume") => {
    setBusy(true);
    try { await api[action](); } catch (e) { toast(errText(e), true); } finally { setBusy(false); }
  };
  return <div className="execution-banner" role="region" aria-label="Preparação do ambiente em andamento">
    {job.state === "paused" ? <Pause size={19} /> : <Loader2 className="spin" size={19} />}
    <div className="execution-summary"><b>Ambiente {job.package} · {setupLabels[job.state] || job.state}</b>
      <div>{current?.title || "Aguarda próxima etapa"}</div>
      <small>{done}/{job.steps.length} etapas · {Math.floor(elapsed / 60)} min {elapsed % 60} s decorridos</small>
    </div>
    <button className="btn sm" disabled={busy || job.state === "cancelling"} onClick={() => control(paused ? "resume" : "pause")}>{paused ? <Play size={14} /> : <Pause size={14} />}{paused ? "Retomar" : "Pausar"}</button>
    <button className="btn sm" disabled={busy || job.state === "cancelling"} onClick={() => control("cancel")}><X size={14} /> Cancelar</button>
    <button className="btn sm primary" onClick={onOpen}>Ver preparação</button>
  </div>;
}
