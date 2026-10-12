import { useEffect, useState } from "react";
import { Check, Copy, ExternalLink } from "lucide-react";
import type { Ctx } from "./App";
import { errText } from "./api";

function CodeBlock({ text, language }: { text: string; language: string }) {
  const [copied, setCopied] = useState(false);
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState("");
  useEffect(() => {
    if (!copied) return;
    const timer = window.setTimeout(() => setCopied(false), 2200);
    return () => window.clearTimeout(timer);
  }, [copied]);
  const copy = async () => {
    if (busy) return;
    setBusy(true); setError(""); setCopied(false);
    try {
      if (window.runtime?.ClipboardSetText) {
        if (!await window.runtime.ClipboardSetText(text)) throw new Error("A área de transferência não aceitou o texto.");
      } else if (navigator.clipboard?.writeText) await navigator.clipboard.writeText(text);
      else throw new Error("A cópia não está disponível nesta sessão.");
      setCopied(true);
    } catch (e) { setError("Não foi possível copiar. Pode selecionar o texto abaixo. " + errText(e)); }
    finally { setBusy(false); }
  };
  const label = language === "zsh" ? "Configuração para o .zshrc" : language === "sh" || language === "bash" ? "Comandos para o terminal" : "Configuração";
  return <div className="guide-code">
    <div className="guide-code-head"><span>{label}</span>
      <button type="button" className="btn sm" disabled={busy} onClick={copy} aria-label={"Copiar " + label.toLocaleLowerCase("pt-PT")}>
        {copied ? <Check size={14} /> : <Copy size={14} />}<span aria-live="polite">{copied ? "Copiado" : "Copiar"}</span>
      </button>
    </div>
    <pre><code>{text}</code></pre>
    {error && <p className="st-failed" role="alert">{error}</p>}
  </div>;
}

// Only explicitly fenced, curated examples receive a copy button. Prose is
// never guessed to be a command; all content stays escaped React text.
function StepContent({ text }: { text: string }) {
  const blocks = [...text.matchAll(/```([a-z][a-z0-9-]*)?\n([\s\S]*?)\n```/g)];
  if (!blocks.length) return <p className="guide-text">{text}</p>;
  let cursor = 0;
  const parts = blocks.map((match, i) => {
    const before = text.slice(cursor, match.index).trim();
    cursor = match.index! + match[0].length;
    return <div key={i}>{before && <p className="guide-text">{before}</p>}<CodeBlock language={match[1] || ""} text={match[2]} /></div>;
  });
  const after = text.slice(cursor).trim();
  return <>{parts}{after && <p className="guide-text">{after}</p>}</>;
}

export function ManualGuide({ ctx, item, steps, url }: { ctx: Ctx; item: string; steps: string[]; url?: string }) {
  const marked = ctx.manualProgress[item] || [];
  const done = steps.filter((step) => marked.includes(step)).length;
  return <div className="manual-guide">
    {steps.length > 0 && <>
      <div className="guide-progress"><span aria-live="polite">{done} de {steps.length} passos assinalados por si</span><progress value={done} max={steps.length} aria-label="Progresso do checklist pessoal" /></div>
      <p className="muted guide-help">Marque os passos à medida que os conclui. As marcas ficam guardadas; a instalação continua a ser verificada pelo estado do computador.</p>
    </>}
    <ol className="guide-steps">{steps.map((step, i) => <li key={i} className={marked.includes(step) ? "marked" : ""}>
      <label className="guide-check"><input type="checkbox" checked={marked.includes(step)} onChange={() => ctx.toggleManualStep(item, step)} aria-label={`Assinalar passo ${i + 1} de ${ctx.items[item]?.name || item}`} />
        <span>Passo {i + 1}{marked.includes(step) ? " · assinalado por si" : ""}</span>
      </label>
      <StepContent text={step} />
    </li>)}</ol>
    {url && <a href="#" onClick={(e) => { e.preventDefault(); window.runtime?.BrowserOpenURL(url); }}>Documentação oficial <ExternalLink size={11} /></a>}
  </div>;
}
