import { useEffect, useMemo, useRef, useState } from "react";
import { Check, ArrowRight, Sparkles, ShieldAlert } from "lucide-react";
import type { Ctx } from "./App";
import { AREAS, sensitiveReason } from "./App";
import type { Profile, Selections } from "./api";
import { StatusBadge } from "./Catalog";
import { Modal } from "./Modals";
import icon from "./icon.png";

// First run: "Como vai usar este computador?" — choosing only fills the selection.
export function Wizard({ ctx, onDone }: { ctx: Ctx; onDone: () => void }) {
  const [mode, setMode] = useState("empty");
  const [replacing, setReplacing] = useState(false);
  const [existing, setExisting] = useState<string[]>([]);
  const changedExisting = useRef(false);
  const [busy, setBusy] = useState(false);
  useEffect(() => { if (mode === "current") ctx.refresh(); }, [mode]);
  const found = ctx.data.items.filter((it) => it.method !== "admin" && ["installed", "provided"].includes(ctx.states[it.id]?.status));
  useEffect(() => {
    if (mode === "current" && !changedExisting.current) setExisting(found.filter((it) => !!ctx.sel[it.id]).map((it) => it.id));
  }, [mode, ctx.states, ctx.sel]);
  const chooseExisting = (id: string) => setExisting((old) => {
    changedExisting.current = true;
    if (old.includes(id)) return old.filter((x) => x !== id);
    const group = ctx.items[id]?.group;
    return [...old.filter((x) => !group || ctx.items[x]?.group !== group), id];
  });
  const usage = ctx.data.profiles.filter((p) => p.builtin && p.kind === "usage");
  const comps = ctx.data.profiles.filter((p) => p.builtin && p.kind === "complement");
  const energy = comps.filter((p) => Object.keys(p.selections).some((id) => ctx.items[id]?.group === "power-profile"));
  const other = comps.filter((p) => !energy.includes(p));
  const [picked, setPicked] = useState<string[]>([]);
  const [energyPick, setEnergyPick] = useState<string>("");
  const toggle = (id: string) => setPicked((p) => (p.includes(id) ? p.filter((x) => x !== id) : [...p, id]));
  const chosen = useMemo(() => [...usage, ...other, ...energy].filter((p) => picked.includes(p.id!) || p.id === energyPick), [picked, energyPick, usage, other, energy]);
  const total = new Set(chosen.flatMap((p) => Object.keys(p.selections))).size;

  const finish = async (ps: Profile[]) => {
    if (busy) return;
    setBusy(true);
    try {
      if (mode === "profiles") {
        if (!await ctx.applyProfiles(ps, true)) return;
      } else {
        const next: Selections = mode === "empty" ? Object.fromEntries(Object.entries(ctx.sel).filter(([id]) => ["installed", "provided"].includes(ctx.states[id]?.status))) : {};
        if (mode === "current") existing.filter((id) => found.some((it) => it.id === id)).forEach((id) => { next[id] = structuredClone(ctx.sel[id] || {}); });
        ctx.setSel(next, Object.fromEntries(Object.keys(next).map((id) => [id, ["estado"]])));
      }
      ctx.setReverts([]);
      onDone();
    } finally { setBusy(false); }
  };

  return (
    <div className="wizard">
      <div className="wizard-in">
        <div className="hero" style={{ marginBottom: 6 }}>
          <img src={icon} alt="" width={52} height={52} />
          <div>
            <h1 style={{ margin: 0 }}>Como quer preparar este computador?</h1>
            <p className="muted" style={{ margin: "4px 0 0" }}>Escolha um ponto de partida. Vai poder personalizar tudo antes de aplicar.</p>
          </div>
        </div>
        <fieldset className="start-choices">
          <legend className="sr-only">Ponto de partida</legend>
          {[
            ["empty", "Escolher no catálogo", "Os já instalados ficam marcados; escolha o que quer acrescentar."],
            ["profiles", "Usar um perfil pronto", "Comece com escolhas para o seu tipo de utilização."],
            ["current", "Partir deste computador", "Reveja os itens encontrados, já marcados para si."],
          ].map(([value, title, desc]) => <label key={value} className={"radio-card" + (mode === value ? " on" : "")}>
            <input type="radio" name="start" value={value} checked={mode === value} onChange={() => setMode(value)} disabled={busy} />
            <div><b>{title}</b><small>{desc}</small></div>
          </label>)}
        </fieldset>
        {mode === "profiles" && <>
        <h2 className="wizard-sub">Como vai usar este computador?</h2>
        <div className="wizard-grid" role="group" aria-label="Perfis de utilização">
          {usage.filter((p) => Object.keys(p.selections).length > 0).map((p) => <ProfileChoice key={p.id} ctx={ctx} p={p} on={picked.includes(p.id!)} toggle={() => toggle(p.id!)} />)}
        </div>
        <details className="wizard-extras"><summary>Adicionar complementos (opcional)</summary>
        <div className="wizard-grid small" role="group" aria-label="Complementos">
          {other.map((p) => <ProfileChoice key={p.id} ctx={ctx} p={p} on={picked.includes(p.id!)} toggle={() => toggle(p.id!)} compact />)}
        </div>
        <fieldset className="energy-pick">
          <legend>Energia (escolha no máximo um)</legend>
          <label className={"radio-card" + (energyPick === "" ? " on" : "")}><input type="radio" name="energy" checked={energyPick === ""} onChange={() => setEnergyPick("")} /><div><b>Não alterar</b><small>Mantém o perfil de energia atual.</small></div></label>
          {energy.map((p) => (
            <label key={p.id} className={"radio-card" + (energyPick === p.id ? " on" : "")}><input type="radio" name="energy" checked={energyPick === p.id} onChange={() => setEnergyPick(p.id!)} /><div><b>{p.name}</b><small>{p.audience}</small></div></label>
          ))}
        </fieldset>
        </details>
        </>}
        {mode === "empty" && <div className="panel"><h2>Preparar a partir do catálogo</h2><p>Os itens já instalados ficam marcados automaticamente. Na próxima etapa, acrescente aplicações e ajustes. Outras escolhas pendentes anteriores serão retiradas.</p></div>}
        {mode === "current" && <div className="panel">
          <h2>Escolher itens encontrados neste computador</h2>
          <p className="muted">{found.length} itens confirmados no catálogo. As caixas escolhem o que vai para o perfil. Os ajustes apresentados já correspondem às opções atuais; configurações fora do catálogo não são importadas.</p>
          {found.length === 0 && <p className="muted">{Object.keys(ctx.states).length ? "Nenhum item confirmado agora." : "A detetar o sistema…"}</p>}
          <div className="inventory-choices">{found.map((it) => <label className="existing-choice" key={it.id}>
            <input type="checkbox" checked={existing.includes(it.id)} onChange={() => chooseExisting(it.id)} disabled={busy} />
            <span><b>{it.name}</b><small>{AREAS[it.area]?.title}</small></span>
            <StatusBadge st={ctx.states[it.id]} method={it.method} adminOp={it.adminOp} />
          </label>)}</div>
        </div>}
        <div className="wizard-foot">
          <span className="muted" style={{ flex: 1 }}>{mode === "profiles" ? (chosen.length ? `${total} itens escolhidos. Pode ajustá-los no catálogo.` : "Escolha um perfil para continuar.") : mode === "current" ? `${existing.filter((id) => found.some((it) => it.id === id)).length} itens escolhidos` : "Nada é instalado até confirmar um plano."}</span>
          <button className="btn ghost" disabled={busy} onClick={onDone}>Agora não</button>
          <button className="btn primary" disabled={busy || (mode === "profiles" && !chosen.length) || (mode === "current" && !existing.some((id) => found.some((it) => it.id === id)))} onClick={() => Object.keys(ctx.sel).some((id) => !(ctx.origins[id] || []).includes("estado")) || ctx.reverts.length ? setReplacing(true) : finish(chosen)}>{busy ? "A preparar…" : "Continuar para o catálogo"} <ArrowRight size={16} /></button>
        </div>
      </div>
      {replacing && <Modal title="Substituir as suas escolhas?" onClose={() => setReplacing(false)} footer={<><button className="btn" data-autofocus onClick={() => setReplacing(false)}>Manter escolhas</button><button className="btn primary" onClick={() => {setReplacing(false);finish(chosen);}}>Substituir e continuar</button></>}><p>O novo ponto de partida vai substituir as escolhas atuais e os pedidos de reversão. Os perfis guardados mantêm-se.</p></Modal>}
    </div>
  );
}

function ProfileChoice({ ctx, p, on, toggle, compact }: { ctx: Ctx; p: Profile; on: boolean; toggle: () => void; compact?: boolean }) {
  const ids = Object.keys(p.selections);
  const sample = ids.map((id) => ctx.items[id]?.name).filter(Boolean).slice(0, compact ? 3 : 4);
  const sens = ids.filter((id) => sensitiveReason(ctx.items[id])).length;
  return (
    <button className={"choice" + (on ? " on" : "")} role="checkbox" aria-checked={on} onClick={toggle}>
      <span className="choice-check" aria-hidden>{on ? <Check size={15} strokeWidth={3} /> : null}</span>
      <span className="choice-title">{p.name}</span>
      <span className="choice-aud">{p.audience}</span>
      <span className="choice-meta"><span className="tag">{ids.length} itens</span></span>
      <span className="choice-sample"><Sparkles size={12} /> {sample.join(", ")}{ids.length > sample.length ? ` e mais ${ids.length - sample.length}` : ""}</span>
      {sens > 0 && <span className="choice-sens"><ShieldAlert size={12} /> {sens} item(ns) sensível(is) — verá os detalhes no plano</span>}
      <span className="sr-only">{on ? "selecionado" : "não selecionado"}</span>
    </button>
  );
}
