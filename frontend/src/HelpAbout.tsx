import { useEffect, useRef, useState } from "react";
import { BookOpen, Copy, Search, ChevronLeft, ChevronRight, ExternalLink, RefreshCw } from "lucide-react";
import type { Ctx } from "./App";
import { api, errText, type AboutInfo } from "./api";
import manual from "../../docs/manual.json";
import icon from "./icon.png";

type Block = { kind: string; text?: string; items?: string[]; headers?: string[]; rows?: string[][]; tone?: string; language?: string };
const normalize = (s: string) => s.normalize("NFD").replace(/[\u0300-\u036f]/g, "").toLowerCase();
const substitute = (s: string, version: string) => s.replaceAll("{{version}}", version);

function ManualBlock({ block: b, version }: { block: Block; version: string }) {
  const text = (s: string) => substitute(s, version);
  if (b.kind === "steps") return <ol>{b.items?.map((s, i) => <li key={i}>{text(s)}</li>)}</ol>;
  if (b.kind === "list") return <ul>{b.items?.map((s, i) => <li key={i}>{text(s)}</li>)}</ul>;
  if (b.kind === "code") return <pre className="manual-code"><code>{text(b.text || "")}</code></pre>;
  if (b.kind === "note") return <div className={"note " + (b.tone === "warn" ? "warn" : "info")}>{text(b.text || "")}</div>;
  if (b.kind === "table") return <div className="manual-table"><table className="t"><thead><tr>{b.headers?.map((h, i) => <th key={i} scope="col">{h}</th>)}</tr></thead><tbody>{b.rows?.map((row, i) => <tr key={i}>{row.map((c, j) => <td key={j}>{text(c)}</td>)}</tr>)}</tbody></table></div>;
  return <p>{text(b.text || "")}</p>;
}

export function HelpPage({ ctx, topic, onTopic }: { ctx: Ctx; topic: string; onTopic: (id: string) => void }) {
  const [query, setQuery] = useState("");
  const article = useRef<HTMLElement>(null);
  const active = manual.sections.find((s) => s.id === topic) || manual.sections[0];
  const index = manual.sections.indexOf(active);
  const words = normalize(query.trim()).split(/\s+/).filter(Boolean);
  const chapters = manual.sections.filter((s) => words.every((word) => normalize(JSON.stringify(s)).includes(word)));
  const navigate = (id: string) => {
    onTopic(id);
    requestAnimationFrame(() => article.current?.scrollIntoView({ block: "start", behavior: "smooth" }));
  };
  return <>
    <div className="page-head"><div><h1>Ajuda e manual</h1><p>Todos os capítulos estão incluídos na aplicação. Pode consultá-los sem internet.</p></div></div>
    <div className="help-layout">
      <nav className="help-toc" aria-label="Capítulos do manual">
        <label className="search"><Search size={15} aria-hidden="true" /><input aria-label="Pesquisar no manual" placeholder="Pesquisar no manual…" value={query} onChange={(e) => setQuery(e.target.value)} /></label>
        <p className="muted help-count" role="status">{chapters.length} de {manual.sections.length} capítulos</p>
        {chapters.length === 0 && <p className="muted">Nenhum capítulo encontrado. Experimente outras palavras.</p>}
        {chapters.map((s) => <button className={"help-chapter" + (s.id === active.id ? " on" : "")} key={s.id} onClick={() => navigate(s.id)} aria-current={s.id === active.id ? "page" : undefined}>{s.title}</button>)}
      </nav>
      <article className="panel manual-article" ref={article} aria-labelledby="manual-title">
        <div className="manual-meta">Meu Ubuntu {ctx.data.version} · {index + 1}/{manual.sections.length}</div>
        <h2 id="manual-title" tabIndex={-1}>{active.title}</h2>
        <p className="manual-summary">{active.summary}</p>
        {active.blocks.map((b, i) => <ManualBlock key={active.id + i} block={b} version={ctx.data.version} />)}
        <div className="manual-paging">
          <button className="btn sm" disabled={index === 0} onClick={() => navigate(manual.sections[index - 1].id)}><ChevronLeft size={14} /> Anterior</button>
          <button className="btn sm" disabled={index === manual.sections.length - 1} onClick={() => navigate(manual.sections[index + 1].id)}>Seguinte <ChevronRight size={14} /></button>
        </div>
      </article>
    </div>
  </>;
}

export function AboutPage({ ctx }: { ctx: Ctx }) {
  const [info, setInfo] = useState<AboutInfo | null>(null);
  const [error, setError] = useState("");
  const [copying, setCopying] = useState(false);
  const [copyError, setCopyError] = useState("");
  const mounted = useRef(true);
  const load = async () => {
    setError("");
    try { const result = await api.about(); if (mounted.current) setInfo(result); }
    catch (e) { if (mounted.current) setError(errText(e)); }
  };
  useEffect(() => { mounted.current = true; load(); return () => { mounted.current = false; }; }, []);
  const s = ctx.data.system;
  const support = [
    `Meu Ubuntu ${info?.version || ctx.data.version}`,
    `Sistema: ${s.pretty || s.id} (${s.release}, ${s.codename})`,
    `Arquitetura: ${s.arch}`, `GNOME: ${s.gnome || "não detetado"}`, `Sessão: ${s.session || "não detetada"}`,
    `Modo: ${s.simulated ? "simulação" : "real"}`,
    ...(info ? [`Go: ${info.goVersion}`, `Wails: ${info.wailsVersion || "não disponível"}`] : []),
    `Catálogo: ${ctx.data.items.length} itens`,
  ].join("\n");
  const copy = async () => {
    if (copying) return;
    setCopying(true); setCopyError("");
    try {
      if (window.runtime?.ClipboardSetText) {
        if (!await window.runtime.ClipboardSetText(support)) throw new Error("O sistema recusou a cópia.");
      } else if (navigator.clipboard?.writeText) await navigator.clipboard.writeText(support);
      else throw new Error("A área de transferência não está disponível.");
      if (mounted.current) ctx.toast("Informações técnicas copiadas.");
    } catch (e) { if (mounted.current) setCopyError(errText(e)); }
    finally { if (mounted.current) setCopying(false); }
  };
  return <>
    <div className="page-head"><div><h1>Sobre</h1><p>Versão, documentação, informações técnicas e atribuições.</p></div></div>
    <div className="panel about-intro">
      <img src={icon} alt="" />
      <div><h2>Meu Ubuntu <span className="badge">{info?.version || ctx.data.version}</span></h2><p>Aplicação gratuita e open source para instalar ferramentas, personalizar o Ubuntu e guardar perfis, com revisão antes de aplicar e verificação do resultado.</p><p className="muted">Versão em execução: <b>{info?.version || ctx.data.version}</b>{(info?.version || ctx.data.version).includes("-beta") && " · Beta para testes"}{(info?.version || ctx.data.version).includes("-rc") && " · Candidata ao lançamento"}</p><p className="muted">Criado por {info?.author || "Paulo Silva"} · Alvo: Ubuntu 26.04 LTS, GNOME, amd64</p></div>
    </div>
    <div className="about-grid">
      <section className="panel"><h2><BookOpen size={18} /> Documentação</h2><p>Aprenda os fluxos, consulte permissões e dados locais ou encontre uma solução para um problema.</p><div className="about-actions"><button className="btn primary" onClick={() => ctx.openHelp()}>Abrir manual completo</button><button className="btn" onClick={() => ctx.openHelp("problemas-frequentes")}>Problemas frequentes</button><button className="btn" onClick={() => ctx.openHelp("dados-e-privacidade")}>Dados e privacidade</button></div><p className="muted">A documentação técnica e a referência integral do catálogo acompanham o código-fonte e o pacote .deb.</p></section>
      <section className="panel"><h2>Informações para suporte</h2>{error && <div className="note bad" role="alert">{error}<button className="btn sm" onClick={load}><RefreshCw size={14} /> Tentar novamente</button></div>}<pre className="support-info">{support}</pre><button className="btn" onClick={copy} disabled={copying || !info}><Copy size={15} /> {copying ? "A copiar…" : "Copiar informações técnicas"}</button>{!info && !error && <p className="muted" role="status">A ler informações do executável…</p>}{copyError && <p className="st-failed" role="alert">{copyError} Pode selecionar e copiar o resumo acima.</p>}<p className="muted">Este resumo não inclui escolhas, caminhos pessoais nem registos de execução.</p></section>
    </div>
    <section className="panel"><h2>Projeto e licenças</h2><p>O código do Meu Ubuntu usa a GNU GPLv3. A aplicação é gratuita; suporte e implantação podem ser serviços pagos. A licença acompanha o código-fonte e o pacote.</p><p>As dependências, a fonte dos relatórios e os logos têm avisos próprios. As marcas pertencem aos seus titulares e não implicam patrocínio.</p>{info?.attributions.map((a) => <details className="attribution" key={a.name}><summary>{a.name} · {a.license}</summary><button className="btn sm" onClick={() => window.runtime?.BrowserOpenURL(a.url)}><ExternalLink size={13} /> Fonte e licença</button><pre className="license-text">{a.text}</pre></details>)}<button className="btn sm" onClick={() => ctx.openHelp("contribuir-e-licencas")}>Contribuir e conhecer os limites</button></section>
  </>;
}
