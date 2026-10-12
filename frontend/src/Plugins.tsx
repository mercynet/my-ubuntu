import { useState } from "react";
import { ChevronDown, ListPlus, RefreshCw } from "lucide-react";
import type { Ctx } from "./App";
import type { Item } from "./api";
import { AppIcon } from "./AppIcon";
import { Card, StatusBadge } from "./Catalog";
import { ManualGuide } from "./ManualGuide";

export function PluginsPage({ ctx, tool, onTool }: { ctx: Ctx; tool: string; onTool: (id: string) => void }) {
  const targets = new Set(ctx.data.items.flatMap((item) => item.addonsFor || []));
  const tools = ctx.data.items.filter((item) => targets.has(item.id));
  const current = tools.some((item) => item.id === tool) ? tool : "all";
  return <>
    <div className="page-head"><div><h1>Plugins e complementos</h1><p>Escolha primeiro a ferramenta. Encontre os seus plugins, integrações e passos de ativação no mesmo grupo.</p></div>
      <button className="btn" onClick={() => ctx.refresh()}><RefreshCw size={15} /> Atualizar estado</button></div>
    <div className="catalog-filters"><label>Ferramenta<select value={current} onChange={(e) => onTool(e.target.value)} aria-label="Agrupar complementos por ferramenta">
      <option value="all">Todas as ferramentas ({tools.length})</option>
      {tools.map((item) => <option key={item.id} value={item.id}>{item.name}</option>)}
    </select></label><span className="muted">Pacotes instaláveis e guias de configuração identificados em cada grupo.</span></div>
    {tools.filter((item) => current === "all" || item.id === current).map((item) => <PluginGroup key={current + ":" + item.id} ctx={ctx} tool={item} expanded={current !== "all" || item.id === "zsh"} />)}
  </>;
}

function PluginGroup({ ctx, tool, expanded }: { ctx: Ctx; tool: Item; expanded: boolean }) {
  const [open, setOpen] = useState(expanded);
  const items = ctx.data.items.filter((item) => item.addonsFor?.includes(tool.id));
  const packages = items.filter((item) => item.method !== "manual");
  const guides = items.filter((item) => item.method === "manual");
  const selected = packages.filter((item) => ctx.sel[item.id]).length;
  const addPackages = () => {
    const selections = { ...ctx.sel }, origins = { ...ctx.origins };
    packages.forEach((item) => {
      if (ctx.states[item.id]?.status === "incompatible") return;
      if (!selections[item.id]) { selections[item.id] = {}; origins[item.id] = ["manual"]; }
    });
    ctx.setSel(selections, origins);
    ctx.setReverts(ctx.reverts.filter((id) => !selections[id]));
    ctx.toast("Complementos incluídos nas escolhas. Reveja o plano para instalar o que falta.");
  };
  return <section className="plugin-group">
    <button className="plugin-group-heading" onClick={() => setOpen(!open)} aria-expanded={open} aria-controls={"plugins-" + tool.id}>
      <AppIcon item={tool} /><span className="plugin-group-label"><b>{tool.name}</b><small>{packages.length ? `${packages.length} complementos · ${selected} marcados` : "Guia de plugins e extensões"}</small></span>
      <StatusBadge st={ctx.states[tool.id]} method={tool.method} /><ChevronDown size={18} className={open ? "rotated" : ""} />
    </button>
    {open && <div id={"plugins-" + tool.id} className="plugin-group-body">
      <div className="plugin-actions"><p className="muted">{packages.length ? "Instalar os pacotes e ativá-los na shell são etapas separadas. Consulte as notas de cada complemento." : "Configure os complementos dentro desta ferramenta. Este guia não instala plugins automaticamente."}</p>
        {!ctx.sel[tool.id] && !["installed", "provided"].includes(ctx.states[tool.id]?.status) && <button className="btn sm" disabled={ctx.running} onClick={() => ctx.toggle(tool.id)}>Incluir {tool.name}</button>}
        {packages.length > 0 && <button className="btn sm" disabled={ctx.running || packages.every((item) => !!ctx.sel[item.id] || ctx.states[item.id]?.status === "incompatible")} onClick={addPackages}><ListPlus size={14} /> Incluir complementos do grupo</button>}
      </div>
      {packages.length > 0 && <div className="grid">{packages.map((item) => <Card key={item.id} ctx={ctx} item={item} />)}</div>}
      {guides.map((item) => <details className="plugin-guide" key={item.id}><summary>{item.name} <span className="badge b-manual">Guia manual</span></summary>
        <p>{item.summary}</p>{item.manual && <ManualGuide ctx={ctx} item={item.id} steps={item.manual.steps} url={item.manual.url} />}
      </details>)}
    </div>}
  </section>;
}
