import { useMemo, useState } from "react";
import { Check, ChevronDown, ChevronUp, RefreshCw, ShieldAlert, LogOut, Server, Undo2, ListPlus } from "lucide-react";
import type { Ctx } from "./App";
import { AREAS, UNDOABLE, originLabel, sensitiveReason } from "./App";
import type { Item, ItemState } from "./api";
import { settingLabel } from "./settingLabels";
import { AppIcon } from "./AppIcon";
import { ManualGuide } from "./ManualGuide";
import { ProfilePicker } from "./ProfilePicker";
import { SaveProfileDialog } from "./Modals";
import { matches, searchScore } from "./catalogSearch";

export const STATUS: Record<string, string> = {
  installed: "Já instalado / configurado", available: "Disponível", config: "Requer configuração", incompatible: "Incompatível",
  unknown: "Desconhecido", manual: "Manual", provided: "Já incluído no Ubuntu",
};

const METHOD: Record<string, string> = {
  apt: "APT (arquivo Ubuntu)", "apt-repo": "APT (repositório oficial)", flatpak: "Flatpak (utilizador)", gsettings: "Preferência GNOME",
  gitconfig: "git config", template: "Ficheiro modelo", extension: "Extensão GNOME", "power-profile": "Perfil de energia",
  admin: "Operação administrativa", manual: "Manual", "user-script": "Instalador do fornecedor (sem root)",
  "npm-user": "Instalação por utilizador",
};

export { matches } from "./catalogSearch";

export function ItemsPage({ ctx, area }: { ctx: Ctx; area: string }) {
  const [cat, setCat] = useState("Todas");
  const [statusF, setStatusF] = useState("all");
  const [saving, setSaving] = useState(false);
  const global = area === "search";
  const catalog = area === "catalog";
  const selection = area === "selection";
  const agents = area === "ai";
  const services = area === "services";
  const inScope = (i: Item) => selection ? !!ctx.sel[i.id] || ctx.reverts.includes(i.id) : agents ? i.category === "Agentes de IA" : services ? i.category.startsWith("Serviços locais ·") || i.category === "Contentores" : (global || catalog || i.area === area) && (global || !ctx.catalogScope || ctx.catalogScope.has(i.id));
  const all = ctx.data.items.filter((i) => inScope(i) && (ctx.showAdv || !i.advanced || ctx.sel[i.id]));
  const filtered = all.filter((i) => matches(i, ctx.query) && (statusF === "all" || (statusF === "sel" ? !!ctx.sel[i.id] : ctx.states[i.id]?.status === statusF)));
  const cats = useMemo(() => {
    const m = new Map<string, number>();
    filtered.forEach((i) => m.set(i.category, (m.get(i.category) || 0) + 1));
    return [...m.entries()];
  }, [filtered]);
  const category = cats.some(([name]) => name === cat) ? cat : "Todas";
  const shown = filtered.filter((i) => category === "Todas" || i.category === category).sort((a, b) => searchScore(b, ctx.query) - searchScore(a, ctx.query));
  const groups = new Map<string, Item[]>();
  shown.forEach((i) => { const category = i.addonsFor?.length ? "Complementos de " + i.addonsFor.map((id) => ctx.items[id]?.name || id).join(", ") : i.category; const key = global ? "Resultados" : catalog || selection ? (AREAS[i.area]?.title || i.area) + " · " + category : category; groups.set(key, [...(groups.get(key) || []), i]); });
  const grouped = [...groups.entries()];
  if (services) {
    const order = ["Serviços locais · Web", "Serviços locais · Dados e cache", "Serviços locais · Mensageria", "Contentores"];
    grouped.sort(([a], [b]) => order.indexOf(a) - order.indexOf(b));
  }
  const hiddenAdv = ctx.data.items.filter((i) => inScope(i) && matches(i,ctx.query) && i.advanced && !ctx.showAdv && !ctx.sel[i.id]).length;
  const outside = ctx.catalogScope ? Object.keys(ctx.sel).filter((id) => !ctx.catalogScope!.has(id)).length : 0;

  return (
    <>
      <div className="page-head">
        <div style={{ flex: 1 }}>
          <h1>{selection ? "Minhas escolhas" : global ? "Resultados da pesquisa" : catalog ? "Catálogo" : AREAS[area].title}</h1>
          <p>{selection ? "Aqui está o que escolheu para este computador. Reveja as opções antes de aplicar." : global ? `${shown.length} ${shown.length === 1 ? "resultado" : "resultados"} para «${ctx.query.trim()}» em todo o catálogo.` : catalog ? "Escolha aplicações e ajustes. Nada é instalado nesta etapa." : AREAS[area].desc}</p>
        </div>
        {selection ? <><button className="btn" disabled={!Object.keys(ctx.sel).length || ctx.running || !!ctx.editing} onClick={() => setSaving(true)}>Guardar como perfil…</button><button className="btn primary" disabled={ctx.running || (!Object.keys(ctx.sel).length && !ctx.reverts.length)} onClick={() => ctx.reviewPlan()}>Rever e aplicar</button></> : <button className="btn" onClick={() => ctx.refresh()} title="Voltar a detetar o estado real do sistema"><RefreshCw size={15} /> Atualizar estado</button>}
      </div>
      {agents && <div className="note info"><span>O plano prepara as dependências e instala a ferramenta para o seu utilizador. Depois, siga o guia do agente para entrar na sua conta ou ligar um fornecedor. O acesso aos modelos depende desse fornecedor.</span><button className="btn sm" onClick={() => ctx.openPlugins("all")}>Plugins dos agentes</button></div>}
      {(services || area === "dev") && <div className="note info"><span>{services ? "Os pacotes de servidor podem ativar serviços e arrancar com o sistema. Já instalado confirma o pacote; consulte o guia para verificar o serviço e configurar o acesso. Kafka tem um guia manual em Docker." : "Apache, Nginx, Redis, RabbitMQ, PostgreSQL, MariaDB e Kafka estão reunidos na vista Serviços locais."}</span>{!services && <button className="btn sm" onClick={() => ctx.go("services")}>Ver serviços locais</button>}</div>}
      {!selection && !global && !agents && !services && <div className="catalog-view"><ProfilePicker ctx={ctx} disabled={ctx.running} /></div>}
      <div className="catalog-filters" role="group" aria-label="Filtros do catálogo">
        <label>Categoria<select value={category} onChange={(e) => setCat(e.target.value)}>
          <option value="Todas">Todas as categorias ({filtered.length})</option>
          {cats.map(([c, n]) => <option key={c} value={c}>{c} ({n})</option>)}
        </select></label>
        <label>Estado
        <select className="btn" value={statusF} onChange={(e) => setStatusF(e.target.value)} aria-label="Filtrar por estado">
          <option value="all">Todos os estados</option>
          <option value="sel">Selecionados</option>
          {Object.entries(STATUS).map(([k, v]) => <option key={k} value={k}>{v}</option>)}
        </select>
        </label>
        <span className="filter-count" role="status">{shown.length} {shown.length === 1 ? "item" : "itens"}</span>
        {(category !== "Todas" || statusF !== "all") && <button className="btn sm ghost" onClick={() => {setCat("Todas");setStatusF("all");}}>Limpar filtros</button>}
      </div>
      {!selection && !global && !agents && !services && ctx.catalogProfile && <div className="catalog-context">
        <span>Está a explorar <b>{ctx.catalogProfile.name}</b>. Este filtro não muda as suas escolhas.{outside > 0 && <> Tem {outside} escolhas noutras partes do catálogo.</>}</span>
        <button className="btn sm" onClick={() => ctx.go("selection")}>Ver minhas escolhas</button>
      </div>}
      {hiddenAdv > 0 && <p className="muted">{hiddenAdv} {hiddenAdv === 1 ? "opção avançada está oculta" : "opções avançadas estão ocultas"}. Pode mostrá-las no controlo acima.</p>}
      {shown.length === 0 && <div className="empty"><h2>{selection && !Object.keys(ctx.sel).length && !ctx.reverts.length ? "Ainda não escolheu nada" : "Não encontramos itens com estes filtros"}</h2><p>{selection ? "Explore o catálogo ou use um perfil pronto para começar." : "Tente outro nome ou retire os filtros de categoria e estado."}</p><div className="btnrow"><button className="btn primary" onClick={() => ctx.showCatalogProfile("")}>Explorar catálogo completo</button>{selection && <button className="btn" onClick={() => ctx.go("profiles")}>Escolher um perfil</button>}</div></div>}
      {grouped.map(([c, list]) => (
        <section key={c}>
          <div className="cat-title">{c}<span className="line" /><span style={{ textTransform: "none", letterSpacing: 0, fontWeight: 400 }}>{list.filter((i) => ctx.sel[i.id]).length}/{list.length} selecionados</span></div>
          <div className="grid">
            {list.map((i) => {
              if (!i.group) return <Card key={i.id} ctx={ctx} item={i} />;
              const first = list.find((x) => x.group === i.group);
              return first?.id === i.id ? <GroupCard key={"g-" + i.group} ctx={ctx} group={i.group} /> : null;
            })}
          </div>
        </section>
      ))}
      {selection ? <div className="selection-actions"><button className="btn" onClick={() => ctx.go("catalog")}>Adicionar mais itens</button><button className="btn ghost danger" disabled={ctx.running || (!Object.keys(ctx.sel).length && !ctx.reverts.length)} onClick={ctx.clearSelection}>Limpar escolhas…</button></div> : !global && <details className="catalog-secondary"><summary>Mais ações do catálogo</summary><button className="btn" onClick={() => ctx.includeApplied(undefined, shown.map((i) => i.id))}><ListPlus size={15} /> Voltar a marcar os já instalados nesta vista</button></details>}
      {saving && <SaveProfileDialog ctx={ctx} selections={ctx.sel} onClose={() => setSaving(false)} />}
    </>
  );
}

const CONFIG_METHODS = ["gsettings", "gitconfig", "power-profile"];

export function statusLabel(s: string, method?: string, adminOp?: string) {
  if (s === "installed") {
    if (method && CONFIG_METHODS.includes(method)) return "Já configurado";
    if (method === "template") return "Já existe";
    if (method === "admin") return adminOp === "snap-remove-firefox" ? "Snap ausente" : "Já configurado";
    return "Já instalado";
  }
  return STATUS[s] || s;
}

export function StatusBadge({ st, method, adminOp }: { st?: ItemState; method?: string; adminOp?: string }) {
  const s = st?.status || "unknown";
  return <span className={"badge b-" + s} title={st?.detail ? "Estado do sistema: " + st.detail : undefined}>{statusLabel(s, method, adminOp)}</span>;
}

// One card per group of exclusive alternatives, with a select instead of checkboxes.
export function GroupCard({ ctx, group }: { ctx: Ctx; group: string }) {
  const g = ctx.data.groups?.find((x) => x.id === group);
  const members = ctx.data.items.filter((i) => i.group === group);
  const chosen = members.find((m) => ctx.sel[m.id]);
  const current = members.find((m) => ctx.states[m.id]?.status === "installed");
  const shown = chosen || current;
  const st = shown ? ctx.states[shown.id] : undefined;
  const pending = !!chosen && ctx.states[chosen.id]?.status !== "installed";
  return (
    <article className={"card" + (chosen ? " sel" : "")} aria-label={g?.name || group}>
      <div className="card-h">
        <div style={{ minWidth: 0 }}>
          <div className="card-title">{g?.name || group}</div>
          <div className={"card-profile" + (chosen ? " included" : "")}>{chosen ? "Incluído nas escolhas" : "Sem alteração pedida"}</div>
        </div>
        {current ? <span className="badge b-installed" title={ctx.states[current.id]?.detail}>Atual: {current.name.replace(/^Energia: /, "")}</span> : <StatusBadge st={st} method={shown?.method} adminOp={shown?.adminOp} />}
      </div>
      {g && <p>{g.summary}</p>}
      {g && <div className="why">{g.why}</div>}
      <div className="opts">
        <label>Opção:</label>
        <select value={chosen?.id || ""} onChange={(e) => ctx.pickGroup(group, e.target.value || null)} aria-label={g?.name || group} style={{ flex: 1 }}>
          <option value="">Não alterar</option>
          {members.map((m) => (
            <option key={m.id} value={m.id} disabled={ctx.states[m.id]?.status === "incompatible"}>
              {m.name.replace(/^Energia: /, "")}{ctx.states[m.id]?.status === "installed" ? " (atual)" : ""}{ctx.states[m.id]?.status === "incompatible" ? " — incompatível" : ""}
            </option>
          ))}
        </select>
      </div>
      {chosen && <div className="why">{chosen.summary} {chosen.why}</div>}
      {pending && <div className="card-sub">A alternativa escolhida será aplicada ao executar o plano.</div>}
      {chosen && (ctx.states[chosen.id]?.settings || []).filter((d) => d.current !== d.proposed).map((d) => (
        <div className="diff" key={d.key}><div className="row"><span className="mono muted">{d.key}</span><span className="mono cur">{d.current || "(sem valor)"}</span>→<span className="mono new">{d.proposed}</span></div></div>
      ))}
    </article>
  );
}

export function Card({ ctx, item }: { ctx: Ctx; item: Item }) {
  const [open, setOpen] = useState(false);
  const st = ctx.states[item.id];
  const selected = !!ctx.sel[item.id];
  const reverting = ctx.reverts.includes(item.id);
  const present = st?.status === "installed" || st?.status === "provided";
  const canRevert = st?.status === "installed" &&
    (UNDOABLE.includes(item.method) && item.method !== "admin" || ["docker-group", "inotify-watches", "journald-limit"].includes(item.id));
  const sens = sensitiveReason(item);
  const origin = ctx.origins[item.id] || [];
  const depOf = !selected ? Object.keys(ctx.sel).filter((id) => (ctx.items[id]?.requires || []).includes(item.id)) : [];
  const disabled = st?.status === "incompatible" && !selected;
  const diffs = (st?.settings || []).filter((s) => s.current !== s.proposed);
  const conflictsSel = (item.conflicts || []).filter((c) => ctx.sel[c.id]);
  const groupSel = item.group ? ctx.data.items.filter((o) => o.group === item.group && o.id !== item.id && ctx.sel[o.id]) : [];

  return (
    <article className={"card" + (selected ? " sel" : "") + (reverting ? " unsel" : "")} aria-label={item.name}>
      <div className="card-h has-logo">
        <button className={"check" + (selected ? " on" : "")} role="checkbox" aria-checked={selected} disabled={disabled}
          onClick={() => ctx.toggle(item.id)} aria-label={(selected ? "Retirar " : "Incluir ") + item.name + (selected ? " das escolhas" : " nas escolhas")}
          title={selected ? "Retirar apenas das escolhas; mantém o estado do computador" : "Incluir nas escolhas e no perfil; o plano verifica o que falta"}>
          {selected && <Check size={14} strokeWidth={3} />}
        </button>
        <AppIcon item={item} />
        <div style={{ minWidth: 0 }}>
          <div className="card-title">{item.name}</div>
          <div className="card-sub">{METHOD[item.method] || item.method}</div>
          <div className={"card-profile" + (selected ? " included" : "")}>{selected ? (present && origin.includes("estado") ? "Marcado pelo estado atual" : "Incluído nas escolhas") : present ? (ctx.editing ? "Fora do perfil em edição" : "Desmarcado por si") : "Por escolher"}</div>
          {depOf.length > 0 && <div className="origin">Dependência de {depOf.map((id) => ctx.items[id]?.name).join(", ")}</div>}
        </div>
        <StatusBadge st={st} method={item.method} adminOp={item.adminOp} />
      </div>
      <p>{item.summary}</p>
      {!selected && present && !reverting && <div className="selection-result kept">Retirado das escolhas; permanece instalado.</div>}
      {selected && <div className={"selection-result" + (["incompatible", "unknown"].includes(st?.status || "unknown") ? " attention" : "")}>
        {st?.status === "installed" || st?.status === "provided" ? "Já está como pretende; nenhuma alteração necessária."
          : st?.status === "manual" ? "O plano apresentará os passos manuais."
          : st?.status === "incompatible" ? "Este item não poderá ser aplicado neste sistema."
          : st?.status === "unknown" || !st ? "É preciso confirmar o estado antes de decidir o que aplicar."
          : ["apt", "apt-repo", "flatpak", "user-script", "extension", "npm-user"].includes(item.method) && st?.status !== "config" ? "Instalação pendente; reveja o plano para aplicar."
          : "Configuração pendente; reveja o plano para aplicar."}
      </div>}
      {selected && item.options && item.options.length > 0 && (
        <div className="opts">
          {item.options.map((o) => {
            const cur = ctx.sel[item.id]?.options?.[o.key] || o.default;
            return (
              <div key={o.key} style={{ display: "flex", gap: 6, alignItems: "center", flexWrap: "wrap" }}>
                <label>{o.label}:</label>
                {o.multi ? o.values.map((v) => {
                  const on = cur.includes(v.value);
                  return <button key={v.value} className={"optchip" + (on ? " on" : "")} aria-pressed={on}
                    onClick={() => {
                      const next = on ? cur.filter((x) => x !== v.value) : [...cur, v.value];
                      if (next.length) ctx.setOpt(item.id, o.key, next);
                    }}>{v.label}</button>;
                }) : (
                  <select value={cur[0]} onChange={(e) => ctx.setOpt(item.id, o.key, [e.target.value])} aria-label={o.label}>
                    {o.values.map((v) => <option key={v.value} value={v.value}>{v.label}</option>)}
                  </select>
                )}
              </div>
            );
          })}
        </div>
      )}
      {item.npm && item.manual && <details className="agent-start"><summary>Como começar depois de instalar</summary><ManualGuide ctx={ctx} item={item.id} steps={item.manual.steps} url={item.manual.url} /></details>}
      {item.addonsFor?.length && item.method !== "manual" && item.manual && <details className="agent-start"><summary>Como ativar este complemento</summary><ManualGuide ctx={ctx} item={item.id} steps={item.manual.steps} url={item.manual.url} /></details>}
      {item.category.startsWith("Serviços locais ·") && item.manual && <details className="agent-start"><summary>{item.method === "manual" ? "Abrir guia de laboratório" : "Verificar e configurar o serviço"}</summary><ManualGuide ctx={ctx} item={item.id} steps={item.manual.steps} url={item.manual.url} /></details>}
      {ctx.data.items.some((addon) => addon.addonsFor?.includes(item.id)) && <button className="btn sm" style={{ alignSelf: "flex-start" }} onClick={() => ctx.openPlugins(item.id)}>Ver plugins e complementos</button>}
      {selected && diffs.length > 0 && (
        <div className="diff" aria-label="Atual versus proposto">
          {diffs.map((d) => (
            <div className="row" key={d.schema + d.key}>
              <span className="mono muted">{d.key}</span>
              <span className="cur">{settingLabel(item, d, d.current)}</span>→<span className="new">{settingLabel(item, d, d.proposed)}</span>
            </div>
          ))}
        </div>
      )}
      {reverting && (
        <div className="note warn" style={{ margin: 0 }}>
          <Undo2 size={15} />
          <div>Reversão pedida: o plano vai repor o valor anterior registado (ou o padrão do sistema). Será revisto antes de executar.</div>
        </div>
      )}
      {(conflictsSel.length > 0 || groupSel.length > 0) && (
        <div className={"note " + (conflictsSel.some((c) => c.hard) || groupSel.length ? "bad" : "warn")} style={{ margin: 0 }}>
          <ShieldAlert size={15} />
          <div>
            {conflictsSel.map((c) => <div key={c.id}>{c.hard ? "Conflito" : "Atenção"} com {ctx.items[c.id]?.name}: {c.reason}</div>)}
            {groupSel.map((g) => <div key={g.id}>Alternativa a {g.name}: escolha apenas uma.</div>)}
          </div>
        </div>
      )}
      <div className="tags">
        {item.relogin && <span className="tag"><LogOut size={11} /> terminar sessão</span>}
        {item.reboot && <span className="tag">reiniciar</span>}
        {item.service && <span className="tag"><Server size={11} /> serviço local</span>}
        {item.advanced && <span className="tag">avançado</span>}
        {sens && <span className="tag sens" title={sens}><ShieldAlert size={11} /> sensível</span>}
        {(item.requires || []).length > 0 && <span className="tag">requer {(item.requires || []).map((r) => ctx.items[r]?.name || r).join(", ")}</span>}
      </div>
      {canRevert && !selected && (
        <button className={"btn sm" + (reverting ? " danger" : "")} style={{ alignSelf: "flex-start" }} onClick={() => ctx.toggleRevert(item.id)} aria-pressed={reverting}>
          <Undo2 size={13} /> {reverting ? "Cancelar reversão" : "Reverter…"}
        </button>
      )}
      <button className="more" onClick={() => setOpen(!open)} aria-expanded={open}>
        {open ? <ChevronUp size={13} /> : <ChevronDown size={13} />} {open ? "Menos detalhes" : "Detalhes, origem e verificação"}
      </button>
      {open && (
        <dl className="details kv">
          <dt>Estado</dt><dd>{st?.detail || "—"}</dd>
          <dt>Porque incluir</dt><dd>{item.why}</dd>
          {selected && origin.length > 0 && <><dt>Escolhido a partir de</dt><dd>{origin.map(originLabel).join(", ")}</dd></>}
          <dt>Fonte</dt><dd>{item.source}</dd>
          {sens && <><dt>Atenção</dt><dd>{sens}</dd></>}
          {item.packages && <><dt>Pacotes</dt><dd className="mono">{item.packages.join(" ")}</dd></>}
          {item.flatpak && <><dt>Flatpak</dt><dd className="mono">{item.flatpak}</dd></>}
          {item.extension && <><dt>UUID</dt><dd className="mono">{item.extension}</dd></>}
          {item.repo && <><dt>Repositório</dt><dd>{ctx.data.repos.find((r) => r.id === item.repo)?.url}</dd></>}
          {item.script && <><dt>Script</dt><dd className="mono">{item.script.url}</dd></>}
          {item.settings && <><dt>Chaves</dt><dd className="mono">{item.settings.map((s) => s.schema + " " + s.key).join("\n")}</dd></>}
          {item.conflicts && item.conflicts.length > 0 && <><dt>Conflitos</dt><dd>{item.conflicts.map((c) => (ctx.items[c.id]?.name || c.id) + (c.hard ? " (bloqueia)" : " (aviso)")).join(", ")}</dd></>}
          {item.notes && <><dt>Notas</dt><dd>{item.notes}</dd></>}
          {item.verified && <><dt>Verificação</dt><dd>{item.verified}</dd></>}
          {item.manual && <><dt>Passos</dt><dd><ManualGuide ctx={ctx} item={item.id} steps={item.manual.steps} url={item.manual.url} /></dd></>}
        </dl>
      )}
    </article>
  );
}
