import { useEffect, useState } from "react";
import { Loader2, Pause, Play, CheckCircle2, AlertTriangle, XCircle } from "lucide-react";
import type { RunState } from "./PlanRun";
import { api, errText } from "./api";

const ended = ["done", "failed", "intervention", "skipped", "cancelled"];

export function ExecutionProgress({ run, onOpen, toast }: { run: RunState | null; onOpen: () => void; toast: (text: string, error?: boolean) => void }) {
  const [now, setNow] = useState(Date.now());
  const [busy, setBusy] = useState(false);
  useEffect(() => {
    if (!run || run.finished) return;
    const timer = window.setInterval(() => setNow(Date.now()), 1000);
    return () => window.clearInterval(timer);
  }, [run?.id, run?.finished]);
  if (!run || run.finished) return null;
  const steps = run.plan.steps;
  const done = steps.filter(s => ended.includes(run.steps[s.id]?.status)).length;
  const current = steps.find(s => run.steps[s.id]?.status === "running");
  const elapsed = Math.max(0, Math.floor((now - Date.parse(run.started)) / 1000));
  const state = run.state || "running";
  const label = state === "paused" ? "Execução em pausa" : state === "pausing" ? "Vai pausar após o passo atual" : state === "resuming" ? "A conferir antes de retomar" : state === "cancelling" ? "Vai parar após o passo atual" : "Instalação em andamento";
  const toggle = async () => {
    setBusy(true);
    try { if (state === "paused" || state === "pausing") await api.resume(); else await api.pause(); }
    catch (error) { toast(errText(error), true); }
    finally { setBusy(false); }
  };
  return <div className="execution-banner" role="region" aria-label="Execução em andamento">
    {state === "paused" ? <Pause size={19} /> : <Loader2 className="spin" size={19} />}
    <div className="execution-summary"><b>{label}</b><div>{current?.title || (state === "paused" ? "Aguarda retoma; nenhum passo novo será iniciado." : "A preparar o próximo passo…")}</div>
      <small>{done} de {steps.length} passos terminados · {Math.floor(elapsed / 60)} min {elapsed % 60} s decorridos</small>
      {run.lastLine && state !== "paused" && <small className="execution-last">{run.lastLine.slice(0, 160)}</small>}
    </div>
    <button className="btn sm" disabled={busy || state === "resuming" || state === "cancelling"} onClick={toggle}>{state === "paused" || state === "pausing" ? <Play size={14} /> : <Pause size={14} />}{state === "paused" || state === "pausing" ? "Retomar" : "Pausar"}</button>
    <button className="btn sm primary" onClick={onOpen}>Ver execução</button>
  </div>;
}

export function ItemExecution({ run, item }: { run: RunState | null; item: string }) {
  if (!run) return null;
  const steps = run.plan.steps.filter(step => step.item === item);
  if (!steps.length) return null;
  const done = steps.filter(step => ended.includes(run.steps[step.id]?.status)).length;
  const current = steps.find(step => run.steps[step.id]?.status === "running");
  const failed = steps.some(step => ["failed", "intervention"].includes(run.steps[step.id]?.status));
  const cancelled = steps.some(step => run.steps[step.id]?.status === "cancelled");
  const paused = run.state === "paused" && !run.finished;
  const label = failed ? "Requer atenção" : current ? "Em execução" : cancelled ? "Cancelado" : done === steps.length ? "Passos terminados" : paused ? "Em pausa" : run.finished ? "Por concluir" : "Na fila";
  return <div className="item-execution" role="status">
    {current ? <Loader2 size={14} className="spin" /> : failed ? <AlertTriangle size={14} /> : cancelled ? <XCircle size={14} /> : done === steps.length ? <CheckCircle2 size={14} /> : paused ? <Pause size={14} /> : null}
    <span><b>{label}</b> · {done}/{steps.length} passos{current && <small>{current.title}</small>}</span>
  </div>;
}
