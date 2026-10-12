import { useEffect, useState } from "react";
import {
  AlertTriangle, CheckCircle2, XCircle, Loader2, Circle, MinusCircle, Hand, Play, Square, RotateCcw, ShieldCheck,
  LogOut, Power, ChevronDown, ChevronRight, KeyRound, BookOpen,
} from "lucide-react";
import type { Ctx } from "./App";
import type { Plan, PlanStep, StepEvent, LogEvent } from "./api";
import { api } from "./api";
import { sensitiveReason } from "./App";
import { StatusBadge } from "./Catalog";
import { SaveProfileDialog } from "./Modals";
import { settingLabel } from "./settingLabels";
import { AppIcon } from "./AppIcon";
import { ManualGuide } from "./ManualGuide";

export type RunState = {
  id: string;
  plan: Plan;
  steps: Record<string, StepEvent>;
  logs: Record<string, LogEvent[]>;
  started: string;
  last?: string;
  lastLine?: string;
  finished?: string;
  error?: string;
};

const INSTALL_OPS = ["apt-install", "flatpak-install", "extension", "user-script", "npm-user"];
const CONFIG_OPS = ["gsettings", "gitconfig", "power-profile", "template", "admin", "post", "extension-enable"];
const REVERT_OPS = ["undo", "revert"];
function itemSteps(steps: PlanStep[]) {
  const result = new Map<string, PlanStep[]>();
  steps.forEach((s) => { if (s.item) result.set(s.item, [...(result.get(s.item) || []), s]); });
  return [...result.entries()];
}
function actionLabel(steps: PlanStep[]) {
  const labels: string[] = [];
  if (steps.some((s) => INSTALL_OPS.includes(s.op))) labels.push("Instalar");
  if (steps.some((s) => CONFIG_OPS.includes(s.op))) labels.push("Configurar");
  if (steps.some((s) => REVERT_OPS.includes(s.op))) labels.push("Restaurar o estado anterior");
  if (steps.some((s) => s.op === "repo-add")) labels.push("Adicionar o repositório do fornecedor");
  return labels.join(" · ") || "Preparar";
}
const itemCount = (n: number) => `${n} ${n === 1 ? "item" : "itens"}`;

const STEP_LABEL: Record<string, string> = {
  pending: "Pendente", running: "Em curso", done: "Concluído", failed: "Falhou", skipped: "Ignorado",
  intervention: "Requer intervenção", cancelled: "Cancelado",
};

function StepIcon({ s }: { s: string }) {
  const p = { size: 17, className: "st-icon st-" + s };
  switch (s) {
    case "done": return <CheckCircle2 {...p} />;
    case "failed": return <XCircle {...p} />;
    case "running": return <Loader2 {...p} className={p.className + " spin"} />;
    case "intervention": return <Hand {...p} />;
    case "cancelled": case "skipped": return <MinusCircle {...p} />;
  }
  return <Circle {...p} />;
}

function ago(iso?: string, now = Date.now()) {
  if (!iso) return "—";
  const s = Math.max(0, Math.round((now - new Date(iso).getTime()) / 1000));
  if (s < 60) return `há ${s}s`;
  return `há ${Math.floor(s / 60)} min ${s % 60}s`;
}

export function PlanPage({ ctx, plan, busy, executing, run, onExecute, onRetry, onClose }: {
  ctx: Ctx; plan: Plan | null; busy: boolean; executing: boolean; run: RunState | null;
  onExecute: (confirmed: Record<string, string>) => void; onRetry: () => void; onClose: () => void;
}) {
  const [confirmed, setConfirmed] = useState<Record<string, string>>({});
  const [script, setScript] = useState<{ item: string; text: string } | null>(null);
  useEffect(() => { setConfirmed({}); setScript(null); }, [plan]);
  if (run) return <RunView ctx={ctx} run={run} onRetry={onRetry} onClose={onClose} />;
  const interrupted = ctx.data.interrupted;
  if (busy) return <div className="empty"><Loader2 className="spin" /> <p>A detetar o estado real e a gerar o plano…</p></div>;
  if (!plan) {
    return (
      <div className="empty">
        {interrupted && <ResumeNote ctx={ctx} />}
        <p>Selecione itens e clique em «Rever plano». Nada é executado sem esta revisão.</p>
      </div>
    );
  }
  const n = (id: string) => ctx.items[id]?.name || id;
  const admin = plan.steps.filter((s) => s.kind === "admin").length;
  const needConfirm = plan.steps.filter((s) => s.confirm);
  const allConfirmed = needConfirm.every((s) => confirmed[s.item] === (s.expect || "yes"));
  const sensitiveSteps = plan.steps.filter((s) => !s.confirm && sensitiveReason(ctx.items[s.item]));
  const countItems = (ops: string[]) => new Set(plan.steps.filter((s) => s.item && ops.includes(s.op)).map((s) => s.item)).size;
  const ready = (plan.skipped || []).filter((s) => s.status === "installed" || s.status === "provided");
  const unavailable = (plan.skipped || []).filter((s) => s.status !== "installed" && s.status !== "provided");
  return (
    <>
      <div className="page-head">
        <div style={{ flex: 1 }}>
          <h1>Revisão do plano</h1>
          <p>Confira o que vai mudar. Os itens já como pretende não precisam de ser aplicados novamente.</p>
        </div>
        <button className="btn" disabled={executing} onClick={() => ctx.reviewPlan(plan.selection, plan.reverts || [])}><RotateCcw size={15} /> Atualizar plano</button>
        <button className="btn primary" disabled={executing || plan.blocked || plan.steps.length === 0 || !allConfirmed} onClick={() => onExecute(confirmed)} title={allConfirmed ? "" : "Faltam confirmações explícitas"}>
          <Play size={15} /> {executing ? "A verificar e iniciar…" : "Aplicar plano"}
        </button>
      </div>

      <p className="muted">Desmarcar só retira o item das escolhas; não pede desinstalação nem reversão.</p>

      <div className="facts" style={{ marginBottom: 16 }} aria-label="O que o plano vai fazer">
        <div className="fact"><small>Vai instalar</small><div>{itemCount(countItems(INSTALL_OPS))}</div></div>
        <div className="fact"><small>Vai configurar</small><div>{itemCount(countItems(CONFIG_OPS))}</div></div>
        <div className="fact"><small>Vai reverter</small><div>{itemCount(countItems(REVERT_OPS))}</div></div>
        <div className="fact"><small>Já como pretende</small><div>{itemCount(ready.length)}</div></div>
      </div>

      {ctx.data.system.simulated && <div className="note warn"><AlertTriangle size={16} /> Modo simulação: os passos usam um executor controlado e não alteram o sistema.</div>}
      {!plan.blocked && plan.steps.length === 0 && ready.length > 0 && unavailable.length === 0 && !(plan.manual || []).length && !(plan.kept || []).length && <div className="note ok"><CheckCircle2 size={16} /><span>As escolhas já estão como pretende. Não é necessário aplicar novamente.</span></div>}
      {plan.blocked && <div className="note bad"><XCircle size={16} /> Há conflitos que bloqueiam a execução. Remova um dos itens em conflito.</div>}
      {(plan.conflicts || []).map((c, i) => (
        <div key={i} className={"note " + (c.hard ? "bad" : "warn")}><AlertTriangle size={16} /><div><b>{n(c.a)} × {n(c.b)}</b> — {c.reason}<div className="btnrow" style={{ marginTop: 8 }}>
          <button className="btn sm" onClick={() => ctx.editPlanItem(c.a, plan)}>Rever {n(c.a)}</button>
          <button className="btn sm" onClick={() => ctx.editPlanItem(c.b, plan)}>Rever {n(c.b)}</button>
        </div></div></div>
      ))}
      {(plan.warnings || []).map((w, i) => <div key={i} className="note warn"><AlertTriangle size={16} /> {w}</div>)}

      {needConfirm.length > 0 && (
        <div className="panel" style={{ borderColor: "var(--warn)" }}>
          <h2><Hand size={17} /> Confirmações necessárias</h2>
          {needConfirm.map((s) => (
            <div key={s.id} style={{ marginBottom: 10 }}>
              <label className="confirm-row">
                <input type="checkbox" checked={confirmed[s.item] === (s.expect || "yes")} disabled={s.op === "user-script" && !s.expect}
                  onChange={(e) => setConfirmed({ ...confirmed, [s.item]: e.target.checked ? (s.expect || "yes") : "" })} />
                <span><b>{s.title}</b><br />{s.confirm}{s.expect && <><br /><span className="mono muted">SHA-256 {s.expect}</span></>}
                  {s.op === "user-script" && !s.expect && <><br /><span className="st-failed">Não foi possível transferir o script para revisão; este passo não pode ser confirmado.</span></>}</span>
              </label>
              {s.op === "user-script" && s.expect && (
                <button className="btn sm" style={{ marginLeft: 26 }} onClick={async () => {
                  try { setScript({ item: s.item, text: await api.scriptPreview(s.item, s.options, s.expect!) }); } catch (e) { ctx.toast(String(e), true); }
                }}>Ver script</button>
              )}
              {script?.item === s.item && <pre className="log" style={{ marginLeft: 26 }}>{script.text}</pre>}
            </div>
          ))}
          {needConfirm.some((s) => s.op === "user-script") && <p className="muted" style={{ marginBottom: 0 }}>Reveja o script antes de confirmar. Se mudar até à execução, a aplicação pede uma nova revisão.</p>}
        </div>
      )}
      {sensitiveSteps.length > 0 && (
        <div className="note warn"><AlertTriangle size={16} /><div><b>Itens sensíveis neste plano:</b> {sensitiveSteps.map((s) => n(s.item) + " (" + sensitiveReason(ctx.items[s.item]) + ")").filter((v, i, a) => a.indexOf(v) === i).join("; ")}</div></div>
      )}
      <div className="cols" style={{ marginBottom: 14 }}>
        <div className="panel" style={{ marginBottom: 0 }}>
          <h2><ShieldCheck size={17} /> Para aplicar este plano</h2>
          <div className="counter"><b>{plan.steps.length}</b> {plan.steps.length === 1 ? "passo" : "passos"} · <b>{admin}</b> administrativos · <b>{(plan.manual || []).length}</b> itens com passos manuais</div>
          {plan.needsRoot && <p className="muted" style={{ marginBottom: 0 }}><KeyRound size={13} /> O sistema vai pedir a sua senha uma vez (polkit/pkexec). A senha não é guardada e a interface continua sem root.</p>}
          {plan.relogin && <p className="muted" style={{ margin: "4px 0 0" }}><LogOut size={13} /> Algumas alterações exigem terminar sessão.</p>}
          {plan.reboot && <p className="muted" style={{ margin: "4px 0 0" }}><Power size={13} /> Algumas alterações exigem reiniciar.</p>}
          {plan.added && plan.added.length > 0 && <p className="muted" style={{ margin: "4px 0 0" }}>Dependências adicionadas: {plan.added.map(n).join(", ")}</p>}
        </div>
        {(plan.kept || []).length > 0 && (
          <div className="panel" style={{ marginBottom: 0 }}>
            <h2>Reversões não automáticas</h2>
            <p className="muted" style={{ marginTop: 0 }}>Pacotes e aplicações nunca são removidos pelo Meu Ubuntu.</p>
            <p>{(plan.kept || []).map((k) => k.name).join(", ")}</p>
          </div>
        )}
        {ready.length > 0 && (
          <div className="panel" style={{ marginBottom: 0 }}>
            <h2>Já estão como pretende</h2>
            <table className="t"><tbody>
              {ready.map((s) => <tr key={s.item}><td>{s.name}</td><td><StatusBadge st={{ status: s.status, detail: s.reason }} method={ctx.items[s.item]?.method} adminOp={ctx.items[s.item]?.adminOp} /></td><td><button className="btn sm" onClick={() => ctx.editPlanItem(s.item, plan)}>Alterar escolha</button></td></tr>)}
            </tbody></table>
          </div>
        )}
      </div>

      {unavailable.length > 0 && <div className="panel">
        <h2>Itens que precisam de atenção</h2>
        {unavailable.map((s) => <div className="review-item" key={s.item}><div className="item-description"><b>{s.name}</b><p className="muted">{s.reason}</p></div><button className="btn sm" onClick={() => ctx.editPlanItem(s.item, plan)}>Alterar escolha</button></div>)}
      </div>}

      {plan.steps.some((s) => s.item) && <div className="panel">
        <h2>Alterações por item</h2>
        {itemSteps(plan.steps).map(([id, steps]) => <div className="review-item" key={id}>
          {ctx.items[id] && <AppIcon item={ctx.items[id]} />}
          <div className="item-description"><b>{n(id)}</b><p>{actionLabel(steps)}</p>
            {steps.flatMap((s) => s.changes || []).filter((d) => d.current !== d.proposed).map((d, i) => <p className="muted" style={{ whiteSpace: "pre-wrap" }} key={i} title={d.schema + " " + d.key}>
              {settingLabel(ctx.items[id], d, d.current)} → {settingLabel(ctx.items[id], d, d.proposed)}
            </p>)}
          </div>
          <button className="btn sm" aria-label={"Alterar escolha de " + n(id)} onClick={() => ctx.editPlanItem(id, plan)}>Alterar escolha</button>
        </div>)}
      </div>}

      <div className="panel">
        <details className="technical-details"><summary>Detalhes técnicos ({plan.steps.length} {plan.steps.length === 1 ? "passo" : "passos"})</summary>
        {plan.steps.length === 0 && <p className="muted">Nada a executar.</p>}
        <div className="steps">
          {plan.steps.map((s, i) => <PlanStepRow key={s.id} idx={i + 1} title={s.title} kind={s.kind} detail={s.detail} />)}
        </div>
        </details>
      </div>

      {(plan.manual || []).length > 0 && (
        <div className="panel">
          <h2><BookOpen size={17} /> Passos manuais (não automatizados)</h2>
          {(plan.manual || []).map((m) => (
            <div key={m.item} style={{ marginBottom: 10 }}>
              <b>{m.name}</b>
              <ManualGuide ctx={ctx} item={m.item} steps={m.steps} url={m.url} />
            </div>
          ))}
        </div>
      )}
    </>
  );
}

function PlanStepRow({ idx, title, kind, detail }: { idx: number; title: string; kind: string; detail: string[] }) {
  const [open, setOpen] = useState(false);
  return (
    <div className="step">
      <button type="button" className="step-h" onClick={() => setOpen(!open)} aria-expanded={open}>
        {open ? <ChevronDown size={15} /> : <ChevronRight size={15} />}
        <span className="muted mono">{String(idx).padStart(2, "0")}</span>
        <span className="t">{title}</span>
        <span className="right">{kind === "admin" ? <span className="tag"><KeyRound size={11} /> administrativo</span> : <span className="tag">utilizador</span>}</span>
      </button>
      {open && detail?.length > 0 && <div className="step-body"><ul>{detail.map((d, i) => <li key={i} className="mono" style={{ whiteSpace: "pre-wrap" }}>{d}</li>)}</ul></div>}
    </div>
  );
}

export function ResumeNote({ ctx }: { ctx: Ctx }) {
  const r = ctx.data.interrupted;
  if (!r) return null;
  return (
    <div className="note warn" style={{ textAlign: "left" }}>
      <AlertTriangle size={16} />
      <div style={{ flex: 1 }}>
        Uma execução de {new Date(r.started).toLocaleString("pt-PT")} foi interrompida ({r.steps.filter((s) => s.status === "done").length}/{r.steps.length} passos concluídos).
        Retomar volta a verificar o estado real de cada item — não confia no registo de «concluído».
      </div>
      <button className="btn sm" onClick={() => { ctx.setSel(r.selections); ctx.reviewPlan(r.selections, r.reverts || []); }}>Retomar</button>
    </div>
  );
}

function RunView({ ctx, run, onRetry, onClose }: { ctx: Ctx; run: RunState; onRetry: () => void; onClose: () => void }) {
  const [now, setNow] = useState(Date.now());
  const [open, setOpen] = useState<Record<string, boolean>>({});
  const [cancelling, setCancelling] = useState(false);
  const [saving, setSaving] = useState(false);
  useEffect(() => {
    if (run.finished) return;
    const t = window.setInterval(() => setNow(Date.now()), 1000);
    return () => window.clearInterval(t);
  }, [run.finished]);
  const status = (id: string) => run.steps[id]?.status || "pending";
  const counts = run.plan.steps.reduce<Record<string, number>>((m, s) => { const k = status(s.id); m[k] = (m[k] || 0) + 1; return m; }, {});
  const finished = run.plan.steps.length - (counts.pending || 0) - (counts.running || 0);
  const current = run.plan.steps.find((s) => status(s.id) === "running");
  const failed = (counts.failed || 0) + (counts.intervention || 0);
  const manual = run.plan.manual || [];
  const outcomes = [
    ...itemSteps(run.plan.steps).map(([id, steps]) => ({ id, status: steps.some((s) => status(s.id) === "failed") ? "failed" : steps.some((s) => status(s.id) === "intervention") ? "intervention" : steps.every((s) => status(s.id) === "done") ? "done" : "pending" })),
    ...manual.map((m) => ({ id: m.item, status: "manual" })),
  ];
  const successful = run.plan.steps.filter((s) => ["done", "intervention"].includes(status(s.id)));
  const needsReboot = successful.some((s) => ctx.items[s.item]?.reboot);
  const needsRelogin = successful.some((s) => ctx.items[s.item]?.relogin);

  return (
    <>
      <div className="page-head">
        <div style={{ flex: 1 }}>
          <h1>{run.finished ? ({ done: manual.length ? "Parte automática concluída" : "Execução concluída", intervention: "Execução requer intervenção", failed: "Execução terminou com falhas", cancelled: "Execução cancelada" } as Record<string, string>)[run.finished] || "Execução terminada" : "A executar o plano"}</h1>
          <p className="counter">
            <b>{finished}</b> de <b>{run.plan.steps.length}</b> passos terminados
            {counts.done ? <> · <span className="st-done">{counts.done} concluídos</span></> : null}
            {counts.failed ? <> · <span className="st-failed">{counts.failed} falharam</span></> : null}
            {counts.intervention ? <> · <span className="st-intervention">{counts.intervention} requerem intervenção</span></> : null}
            {counts.cancelled ? <> · {counts.cancelled} cancelados</> : null}
          </p>
        </div>
        {!run.finished && (
          <button className="btn danger" disabled={cancelling} onClick={() => { setCancelling(true); api.cancel(); }}>
            <Square size={14} /> {cancelling ? "A cancelar após o passo atual…" : "Cancelar"}
          </button>
        )}
        {run.finished && (failed > 0 || run.error || outcomes.some((o) => o.status !== "done")) && <button className="btn primary" onClick={onRetry}><RotateCcw size={15} /> Rever o que falta</button>}
        {run.finished && <button className="btn" onClick={onClose}>Fechar</button>}
      </div>
      {run.finished && <div className="panel">
        <h2>Resultado desta execução</h2>
        {run.error && <div className="note bad" role="alert">{run.error}</div>}
        <div className="facts">
          <div className="fact"><small>Itens confirmados</small><div>{outcomes.filter((o) => o.status === "done").length}</div></div>
          <div className="fact"><small>Com falhas</small><div>{outcomes.filter((o) => o.status === "failed").length}</div></div>
          <div className="fact"><small>A confirmar</small><div>{outcomes.filter((o) => o.status === "intervention").length}</div></div>
          <div className="fact"><small>Por concluir</small><div>{outcomes.filter((o) => o.status === "pending" || o.status === "manual").length}</div></div>
        </div>
        {outcomes.some((o) => o.status !== "done") && <p className="muted">Itens por concluir: {outcomes.filter((o) => o.status !== "done").map((o) => ctx.items[o.id]?.name || o.id).join(", ")}. «Rever o que falta» consulta o sistema novamente antes de propor outra execução.</p>}
        <div className="btnrow" style={{ marginTop: 12 }}>
          {Object.keys(run.plan.selection).length > 0 && <button className="btn" onClick={() => setSaving(true)}>Guardar o perfil escolhido</button>}
          <button className="btn" onClick={() => ctx.go("history")}>Ver histórico e reversão</button>
        </div>
      </div>}
      {!run.finished && (
        <div className="note info">
          <Loader2 size={16} className="spin" />
          <div>
            <div><b>Passo atual:</b> {current?.title || "a preparar…"}</div>
            <div className="muted">Última atividade {ago(run.last, now)}{run.lastLine ? <> — <span className="mono">{run.lastLine.slice(0, 140)}</span></> : null}</div>
            {cancelling && <div>O cancelamento acontece num limite seguro: uma operação APT/dpkg em curso termina primeiro.</div>}
          </div>
        </div>
      )}
      <div className="steps">
        {run.plan.steps.map((s, i) => {
          const st = status(s.id);
          const ev = run.steps[s.id];
          const logs = run.logs[s.id] || [];
          const isOpen = open[s.id] ?? (st === "running" || st === "failed" || st === "intervention");
          return (
            <div className="step" key={s.id}>
              <button type="button" className="step-h" onClick={() => setOpen({ ...open, [s.id]: !isOpen })} aria-expanded={isOpen}>
                <StepIcon s={st} />
                <span className="muted mono">{String(i + 1).padStart(2, "0")}</span>
                <span className="t">{s.title}</span>
                <span className="right">
                  {ev?.note && <span className="k">{ev.note}</span>}
                  <span className={"st-" + st}>{STEP_LABEL[st]}</span>
                  {logs.length > 0 && <span className="k">{logs.length} linhas</span>}
                </span>
              </button>
              {isOpen && (
                <div className="step-body">
                  {ev?.problem && <div className="note bad" style={{ margin: "0 0 6px" }} role="alert"><XCircle size={15} /><div><b>{ev.problem.summary}</b><div>{ev.problem.advice}</div></div></div>}
                  {!ev?.problem && ev?.error && <div className="note bad" style={{ margin: "0 0 6px" }}><XCircle size={15} /> {ev.error}</div>}
                  {st === "intervention" && ev?.note && <div className="note warn">{ev.note}</div>}
                  <details className="technical-details"><summary>Erro e registo técnico</summary>
                    {ev?.error && <p className="mono">{ev.error}</p>}
                    {logs.length === 0 ? <span className="muted">Sem registo ainda.</span> : <pre className="log">{logs.map((l) => new Date(l.time).toLocaleTimeString("pt-PT") + "  " + l.line).join("\n")}</pre>}
                  </details>
                </div>
              )}
            </div>
          );
        })}
      </div>
      {run.finished && (needsReboot || needsRelogin) && (
        <div className="panel" style={{ marginTop: 12 }}>
          <h2><Power size={17} /> Concluir</h2>
          {needsReboot ? <p className="muted" style={{ marginTop: 0 }}>Reinicie para ativar as alterações realizadas que precisam de reinício.</p>
            : <p className="muted" style={{ marginTop: 0 }}>Termine sessão para ativar as alterações realizadas que precisam de uma nova sessão.</p>}
          <div style={{ display: "flex", gap: 8 }}>
            {needsRelogin && !needsReboot && <button className="btn primary" onClick={() => api.endSession("logout").catch((e) => ctx.toast(String(e), true))}><LogOut size={15} /> Terminar sessão</button>}
            {needsReboot && <button className="btn primary" onClick={() => api.endSession("reboot").catch((e) => ctx.toast(String(e), true))}><Power size={15} /> Reiniciar</button>}
          </div>
          <p className="muted" style={{ marginBottom: 0 }}>O GNOME pede confirmação antes; guarde o seu trabalho noutras aplicações.</p>
        </div>
      )}
      {run.finished && manual.length > 0 && <div className="panel" style={{ marginTop: 12 }}>
        <h2><BookOpen size={17} /> Passos manuais por concluir</h2>
        <p className="muted">Estes passos não foram executados pela aplicação. Conclua-os e use «Rever o que falta» para voltar a consultar o sistema.</p>
        {manual.map((m) => <div key={m.item}>
          <h3>{m.name}</h3>
          <ManualGuide ctx={ctx} item={m.item} steps={m.steps} url={m.url} />
        </div>)}
      </div>}
      {saving && <SaveProfileDialog ctx={ctx} selections={run.plan.selection} onClose={() => setSaving(false)} />}
      <p className="muted" style={{ marginTop: 10 }}>Registo completo guardado em {ctx.data.paths.state}/logs/{run.id}.log</p>
    </>
  );
}
