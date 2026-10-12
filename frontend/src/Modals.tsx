import { useEffect, useRef, useState, type ReactNode } from "react";
import { AlertTriangle, ShieldAlert, X, Download, GitMerge, Replace, Save } from "lucide-react";
import type { Ctx } from "./App";
import { AREAS, originLabel, sensitiveReason } from "./App";
import type { Change, MergeConflict, Origins, Profile, ProfileReview, RevertPreview, Selections, Sel } from "./api";
import { api, errText } from "./api";
import { statusLabel } from "./Catalog";

// Accessible modal: labelled dialog, focus moved inside, Escape closes, focus restored.
export function Modal({ title, onClose, children, footer, wide, closeDisabled = false }: { title: string; onClose: () => void; children: ReactNode; footer?: ReactNode; wide?: boolean; closeDisabled?: boolean }) {
  const ref = useRef<HTMLDivElement>(null);
  const closeRef = useRef(onClose); closeRef.current = onClose;
  const disabledRef = useRef(closeDisabled); disabledRef.current = closeDisabled;
  useEffect(() => {
    const prev = document.activeElement as HTMLElement | null;
    const first = ref.current?.querySelector<HTMLElement>("[data-autofocus]") || ref.current?.querySelector<HTMLElement>("button:not(:disabled), [href], input:not(:disabled), select:not(:disabled), textarea:not(:disabled)");
    first?.focus();
    const onKey = (e: KeyboardEvent) => {
      if (e.key === "Escape" && !disabledRef.current) closeRef.current();
      if (e.key === "Tab" && ref.current) {
        const f = ref.current.querySelectorAll<HTMLElement>("button:not(:disabled), [href], input:not(:disabled), select:not(:disabled), textarea:not(:disabled)");
        if (!f.length) return;
        if (!ref.current.contains(document.activeElement)) { e.preventDefault(); (e.shiftKey ? f[f.length - 1] : f[0]).focus(); }
        else if (e.shiftKey && document.activeElement === f[0]) { e.preventDefault(); f[f.length - 1].focus(); }
        else if (!e.shiftKey && document.activeElement === f[f.length - 1]) { e.preventDefault(); f[0].focus(); }
      }
    };
    document.addEventListener("keydown", onKey);
    return () => { document.removeEventListener("keydown", onKey); if (prev?.isConnected) prev.focus(); };
  }, []);
  return (
    <div className="modal-back" onMouseDown={(e) => e.target === e.currentTarget && !closeDisabled && onClose()}>
      <div className={"modal" + (wide ? " wide" : "")} role="dialog" aria-modal="true" aria-label={title} ref={ref}>
        <div className="modal-h"><h2>{title}</h2><button className="iconbtn" disabled={closeDisabled} onClick={onClose} aria-label="Fechar"><X size={16} /></button></div>
        <div className="modal-b">{children}</div>
        {footer && <div className="modal-f">{footer}</div>}
      </div>
    </div>
  );
}

function optText(ctx: Ctx, id: string, s: Sel) {
  const it = ctx.items[id];
  const parts = (it?.options || []).map((o) => {
    const vals = s.options?.[o.key] || o.default;
    return o.label + ": " + vals.map((v) => o.values.find((x) => x.value === v)?.label || v).join(", ");
  });
  return parts.join(" · ") || "sem opções";
}

/* ---------- Conflict resolution: one explicit choice per conflict ---------- */
export function ConflictDialog({ ctx, merge, onDone, onCancel }: {
  ctx: Ctx; merge: { selections: Selections; origins: Origins; conflicts: MergeConflict[] };
  onDone: (s: Selections, o: Origins) => void; onCancel: () => void;
}) {
  const [choice, setChoice] = useState<Record<number, "keep" | "incoming">>({});
  const all = merge.conflicts.every((_, i) => choice[i]);
  const apply = () => {
    const s = { ...merge.selections }; const o = { ...merge.origins };
    merge.conflicts.forEach((c, i) => {
      if (choice[i] !== "incoming") return;
      if (c.kind !== "options") { delete s[c.item]; delete o[c.item]; }
      s[c.incomingItem] = c.incoming;
      o[c.incomingItem] = [c.incomingLabel];
    });
    onDone(s, o);
  };
  return (
    <Modal title={`Resolver ${merge.conflicts.length} conflito(s)`} onClose={onCancel} wide
      footer={<><span className="muted" style={{ flex: 1 }}>{all ? "Pronto para aplicar." : "Escolha uma opção em cada conflito."}</span>
        <button className="btn" onClick={onCancel}>Cancelar</button>
        <button className="btn primary" disabled={!all} onClick={apply}>Aplicar escolhas</button></>}>
      <p className="muted" style={{ marginTop: 0 }}>Os perfis combinados pedem coisas incompatíveis. Nada é decidido pela ordem: escolha o que fica. O resto da combinação já foi deduplicado.</p>
      {merge.conflicts.map((c, i) => (
        <fieldset key={i} className="conflict">
          <legend><AlertTriangle size={14} /> {c.kind === "options" ? c.name : `${c.name} × ${c.incomingName}`} <span className="muted">— {c.reason}</span></legend>
          <label className={"radio-card" + (choice[i] === "keep" ? " on" : "")}>
            <input type="radio" name={"c" + i} checked={choice[i] === "keep"} onChange={() => setChoice({ ...choice, [i]: "keep" })} />
            <div><b>Manter atual:</b> {c.kind === "options" ? optText(ctx, c.item, c.current) : c.name}
              <small>Origem: {(c.currentOrigin || ["manual"]).map(originLabel).join(", ")}</small></div>
          </label>
          <label className={"radio-card" + (choice[i] === "incoming" ? " on" : "")}>
            <input type="radio" name={"c" + i} checked={choice[i] === "incoming"} onChange={() => setChoice({ ...choice, [i]: "incoming" })} />
            <div><b>Usar do {originLabel(c.incomingLabel)}:</b> {c.kind === "options" ? optText(ctx, c.incomingItem, c.incoming) : c.incomingName}
              {c.kind !== "options" && <small>Substitui {c.name} na seleção.</small>}</div>
          </label>
        </fieldset>
      ))}
    </Modal>
  );
}

/* ---------- Profile review (import, CLI, compare) ---------- */
export function ReviewDialog({ ctx, review, onClose }: { ctx: Ctx; review: ProfileReview; onClose: () => void }) {
  const p = review.profile;
  const items = review.items || [];
  const sensitive = items.filter((i) => i.sensitive);
  const group = (f: (s: string) => boolean) => items.filter((i) => f(i.status));
  const change = group((s) => s === "available" || s === "config");
  const ok = group((s) => s === "installed" || s === "provided");
  const manual = group((s) => s === "manual");
  const bad = group((s) => s === "incompatible" || s === "unknown");
  const imported = p.source === "imported";
  const [replaceConfirm, setReplaceConfirm] = useState(false);
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState("");
  const working = useRef(false);

  const saveImported = async () => {
    if (working.current) return;
    working.current = true;
    setBusy(true); setError("");
    try {
      let target: Profile = { ...p };
      if (review.existsImported || review.existsMine) {
        const over = review.existsImported && window.confirm(`Já existe um perfil importado com o identificador «${p.id}». Substituí-lo?\n\nCancelar guarda como cópia.`);
        if (!over) target = { ...p, id: (p.id || "perfil") + "-" + Date.now().toString(36).slice(-4), name: p.name + " (cópia)" };
        await api.saveImported(target, over);
      } else await api.saveImported(target, false);
      ctx.setProfiles(await api.profiles());
      ctx.toast("Guardado em «Perfis importados». Nada foi executado.");
    } catch (e) { setError(errText(e)); } finally { working.current = false; setBusy(false); }
  };
  const apply = async (replace: boolean) => {
    if (working.current) return;
    working.current = true;
    setBusy(true); setError("");
    try {
      if (await ctx.applyProfiles([p], replace)) { onClose(); ctx.go("selection"); }
      else setError("Não foi possível aplicar o perfil. A seleção foi preservada. Pode tentar novamente.");
    } catch (e) { setError(errText(e)); }
    finally { working.current = false; setBusy(false); }
  };

  const List = ({ title, list, cls }: { title: string; list: typeof items; cls?: string }) => list.length ? (
    <div className="review-sec"><h3>{title} ({list.length})</h3>
      <ul className={"review-list " + (cls || "")}>{list.map((i) => <li key={i.id}><b>{i.name}</b> <span className="muted">· {AREAS[i.area]?.title || i.area} · {statusLabel(i.status, ctx.items[i.id]?.method, ctx.items[i.id]?.adminOp)}</span>{i.detail && <small>{i.detail}</small>}</li>)}</ul>
    </div>) : null;

  return (
    <Modal title={(imported ? "Rever perfil importado: " : "Usar perfil: ") + p.name} onClose={onClose} wide closeDisabled={busy}
      footer={<>
        <span className="muted" style={{ flex: 1 }}>Pode personalizar estas escolhas antes de aplicar.</span>
        {imported && <button className="btn" disabled={busy} onClick={saveImported}><Save size={15} /> Guardar em importados</button>}
        {replaceConfirm ? <><button className="btn" disabled={busy} onClick={() => setReplaceConfirm(false)}>Voltar</button><button className="btn danger" disabled={busy} onClick={() => apply(true)}>Confirmar substituição</button></> : <><button className="btn primary" disabled={busy} onClick={() => apply(false)}><GitMerge size={15} /> Adicionar às escolhas</button><button className="btn" disabled={busy} onClick={() => Object.keys(ctx.sel).length ? setReplaceConfirm(true) : apply(true)}><Replace size={15} /> Substituir escolhas…</button></>}
      </>}>
      {busy && <p className="muted" role="status">A preparar o perfil…</p>}
      {error && <div className="note bad" role="alert">{error}</div>}
      {replaceConfirm && <div className="note warn" role="alert"><AlertTriangle size={16} /><div>Vai substituir as escolhas pendentes pelas de <b>{p.name}</b>. Os itens já instalados continuam marcados, salvo exclusões feitas por si. Os perfis guardados mantêm-se. Use «Voltar» para cancelar.</div></div>}
      <dl className="kv meta">
        <dt>Identificador</dt><dd className="mono">{p.id}</dd>
        {p.audience && <><dt>Para quem</dt><dd>{p.audience}</dd></>}
        {p.description && <><dt>Descrição</dt><dd>{p.description}</dd></>}
        <dt>Autoria</dt><dd>{p.author || "não indicada"} <span className="muted">(informativo — não é assinatura nem prova de confiança)</span></dd>
        <dt>Requisitos</dt><dd>{p.requires?.minRelease ? "Ubuntu " + p.requires.minRelease + "+" : "nenhum"}{p.requires?.arch?.length ? " · " + p.requires.arch.join("/") : ""}</dd>
        {review.path && <><dt>Ficheiro</dt><dd className="mono">{review.path}</dd></>}
      </dl>
      {(review.requirement || []).map((r, i) => <div key={i} className="note warn"><AlertTriangle size={15} /> {r}</div>)}
      {(review.warnings || []).map((w, i) => <div key={i} className="note info"><AlertTriangle size={15} /> {w}</div>)}
      {sensitive.length > 0 && (
        <div className="note bad"><ShieldAlert size={16} /><div><b>Itens que merecem atenção ({sensitive.length})</b>
          <ul style={{ margin: "4px 0 0", paddingLeft: 18 }}>{sensitive.map((i) => <li key={i.id}><b>{i.name}</b>: {i.sensitive}</li>)}</ul></div></div>
      )}
      <div className="cols">
        <List title="Vai instalar ou alterar neste PC" list={change} />
        <List title="Já está como pretendido" list={ok} cls="muted" />
      </div>
      <div className="cols">
        <List title="Passos manuais" list={manual} />
        <List title="Incompatível ou desconhecido neste PC" list={bad} />
      </div>
    </Modal>
  );
}

/* ---------- Export summary ---------- */
export function ExportDialog({ ctx, profile, onClose }: { ctx: Ctx; profile: Profile; onClose: () => void }) {
  const ids = Object.keys(profile.selections || {});
  const byArea: Record<string, string[]> = {};
  ids.forEach((id) => { const a = ctx.items[id]?.area || "?"; (byArea[a] ||= []).push(ctx.items[id]?.name || id); });
  const sens = ids.filter((id) => sensitiveReason(ctx.items[id]));
  return (
    <Modal title={"Exportar «" + profile.name + "»"} onClose={onClose} wide
      footer={<><button className="btn" onClick={onClose}>Cancelar</button>
        <button className="btn primary" onClick={async () => {
          try { const path = await api.exportProfile(profile); if (path) { ctx.toast("Exportado para " + path); onClose(); } }
          catch (e) { ctx.toast(errText(e), true); }
        }}><Download size={15} /> Escolher ficheiro e exportar</button></>}>
      <p className="muted" style={{ marginTop: 0 }}>O JSON contém só identificadores de itens e opções validadas. Não contém comandos, scripts, URLs de instalação, senhas, tokens nem certificados. Quem o importar revê tudo antes de qualquer execução.</p>
      <dl className="kv meta">
        <dt>Identificador</dt><dd className="mono">{profile.id}</dd>
        <dt>Itens</dt><dd>{ids.length}</dd>
        {profile.author && <><dt>Autoria</dt><dd>{profile.author}</dd></>}
      </dl>
      {Object.entries(byArea).map(([a, names]) => <p key={a}><b>{AREAS[a]?.title || a}:</b> <span className="muted">{names.join(", ")}</span></p>)}
      {sens.length > 0 && <div className="note warn"><ShieldAlert size={15} /> Inclui itens sensíveis: {sens.map((id) => ctx.items[id]?.name).join(", ")}. Quem importar verá este aviso.</div>}
    </Modal>
  );
}

export function SaveProfileDialog({ ctx, selections, onClose }: { ctx: Ctx; selections: Selections; onClose: () => void }) {
  const [name, setName] = useState("");
  const [description, setDescription] = useState("");
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState("");
  const [overwrite, setOverwrite] = useState(false);
  const save = async () => {
    if (!name.trim() || busy) return;
    setBusy(true); setError("");
    try {
      await api.saveProfile({ name: name.trim(), description, selections, requires: { minRelease: ctx.data.system.release || undefined, arch: [ctx.data.system.arch] } }, overwrite);
      ctx.setProfiles(await api.profiles());
      ctx.toast("Perfil «" + name.trim() + "» guardado.");
      onClose();
    } catch (e) {
      const message = errText(e);
      setError(message);
      if (message.includes("já existe")) setOverwrite(true);
    } finally { setBusy(false); }
  };
  return <Modal title="Guardar escolhas como perfil" onClose={onClose} closeDisabled={busy} footer={<>
    <button className="btn" disabled={busy} onClick={onClose}>Cancelar</button>
    <button className="btn primary" disabled={busy || !name.trim()} onClick={save}>{busy ? "A guardar…" : overwrite ? "Substituir perfil existente" : "Guardar perfil"}</button>
  </>}>
    <p className="muted">Guarda {Object.keys(selections).length === 1 ? "o item escolhido" : `os ${Object.keys(selections).length} itens escolhidos`} nesta execução, incluindo os que ainda estejam pendentes. Pode reutilizar este perfil numa reinstalação.</p>
    <label className="field">Nome<input data-autofocus value={name} disabled={busy} onChange={(e) => { setName(e.target.value); setOverwrite(false); setError(""); }} placeholder="O meu Ubuntu" /></label>
    <label className="field">Descrição (opcional)<textarea value={description} disabled={busy} onChange={(e) => setDescription(e.target.value)} /></label>
    {error && <div className="note bad" role="alert">{error}</div>}
    {overwrite && <p className="muted">Confirme no botão abaixo se pretende substituir as escolhas guardadas nesse perfil.</p>}
  </Modal>;
}

export function RevertDialog({ ctx, change, onClose, onDone }: { ctx: Ctx; change: Change; onClose: () => void; onDone: () => void }) {
  const [preview, setPreview] = useState<RevertPreview | null>(null);
  const [loading, setLoading] = useState(false);
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState("");
  const [force, setForce] = useState(false);
  const [remove, setRemove] = useState(false);
  const request = useRef(0);
  const load = async () => {
    const version = ++request.current;
    setLoading(true); setPreview(null); setForce(false); setError("");
    try { const next = await api.previewRevert(change.id); if (version === request.current) setPreview(next); }
    catch (e) { if (version === request.current) setError(errText(e)); }
    finally { if (version === request.current) setLoading(false); }
  };
  useEffect(() => { load(); return () => { request.current++; }; }, [change.id]);
  const apply = async () => {
    if (!preview || busy) return;
    setBusy(true); setError("");
    try {
      await api.revertReviewed(change.id, force, preview.reviewToken);
      if (remove && ctx.sel[change.item]) ctx.toggle(change.item);
      ctx.toast("Revertido: " + change.title);
      onDone(); onClose();
    } catch (e) { setError(errText(e)); setPreview(null); setForce(false); }
    finally { setBusy(false); }
  };
  return <Modal title={"Rever reversão: " + change.title} onClose={onClose} closeDisabled={busy} wide footer={<>
    <button className="btn" disabled={busy || loading} onClick={load}>Voltar a verificar</button>
    <button className="btn" disabled={busy} onClick={onClose}>Cancelar</button>
    <button className="btn primary" disabled={busy || loading || !preview || (!preview.canRevert && !(preview.forceAllowed && force))} onClick={apply}>{busy ? "A restaurar…" : "Confirmar reversão"}</button>
  </>}>
    {loading && <p role="status">A ler o valor atual e a configuração anterior…{change.kind === "admin" ? " O sistema pode pedir autorização para consultar a cópia administrativa." : ""}</p>}
    {error && <div className="note bad" role="alert">{error}</div>}
    {preview && <>
      <div className="cols">
        <div><b>Estado atual</b><pre className="value-preview">{preview.current || "(sem valor)"}</pre></div>
        <div><b>Será restaurado</b><pre className="value-preview">{preview.target || "(sem valor)"}</pre></div>
      </div>
      {preview.warning && <div className="note warn" style={{ marginTop: 12 }} role="status">{preview.warning}</div>}
      {!preview.canRevert && preview.forceAllowed && <label className="confirm-row"><input type="checkbox" checked={force} disabled={busy} onChange={(e) => setForce(e.target.checked)} /> Quero substituir a alteração posterior pelo valor anterior apresentado.</label>}
      {ctx.sel[change.item] && <label className="confirm-row"><input type="checkbox" checked={remove} disabled={busy} onChange={(e) => setRemove(e.target.checked)} /> Retirar também este item do perfil, para que um novo plano não volte a pedi-lo.</label>}
    </>}
  </Modal>;
}
