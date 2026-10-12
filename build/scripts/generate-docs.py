#!/usr/bin/env python3
"""Generate repository/package docs from the same sources as the offline help."""
import argparse
import json
import re
import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parents[2]
AREAS = {"apps": "Aplicações", "dev": "Desenvolvimento", "infra": "Infraestrutura", "security": "Segurança", "style": "Estilo", "desktop": "Desktop", "extensions": "Extensões", "energy": "Energia e hardware"}


def read(path):
    return json.loads((ROOT / path).read_text())


def table(headers, rows):
    def cell(value):
        return str(value).replace("|", "\\|").replace("\n", "<br>")
    return "\n".join(["| " + " | ".join(map(cell, headers)) + " |", "| " + " | ".join(["---"] * len(headers)) + " |"] + ["| " + " | ".join(map(cell, row)) + " |" for row in rows])


def manual_doc(manual, version):
    ids = [section["id"] for section in manual["sections"]]
    if len(ids) != len(set(ids)) or not all(re.fullmatch(r"[a-z0-9-]+", id) for id in ids):
        raise ValueError("IDs dos capítulos inválidos ou duplicados")
    out = ["<!-- Gerado de docs/manual.json; execute python3 build/scripts/generate-docs.py. -->", "# " + manual["title"], manual["intro"], "## Índice", "\n".join(f'- [{s["title"]}](#{s["id"]})' for s in manual["sections"])]
    for s in manual["sections"]:
        out.extend([f'<a id="{s["id"]}"></a>', "## " + s["title"], s["summary"]])
        for b in s["blocks"]:
            kind = b["kind"]
            if kind == "paragraph":
                out.append(b["text"])
            elif kind in ("steps", "list"):
                out.append("\n".join(f'{i + 1}. {text}' if kind == "steps" else "- " + text for i, text in enumerate(b["items"])))
            elif kind == "note":
                if b["tone"] not in ("info", "warn"):
                    raise ValueError("Tom de nota desconhecido")
                out.append("> " + b["text"])
            elif kind == "code":
                out.append(f'```{b["language"]}\n{b["text"]}\n```')
            elif kind == "table":
                if not all(len(row) == len(b["headers"]) for row in b["rows"]):
                    raise ValueError("Tabela com colunas inconsistentes")
                out.append(table(b["headers"], b["rows"]))
            else:
                raise ValueError("Bloco desconhecido: " + kind)
    return ("\n\n".join(out) + "\n").replace("{{version}}", version)


def catalog_doc(cat, profiles):
    items = {it["id"]: it for it in cat["items"]}
    def link(id):
        if id not in items:
            raise ValueError("Referência desconhecida: " + id)
        return f'[{items[id]["name"]}](#item-{id})'
    out = ["<!-- Gerado do catálogo e dos perfis; não editar diretamente. -->", "# Referência do catálogo e dos perfis", f'{len(items)} itens e {len(profiles)} perfis incluídos. Esta referência descreve as receitas embutidas; o estado e a compatibilidade efetivos são detetados no computador. Os comandos de verificação são evidência usada pelo software, não instruções para aplicar a receita manualmente.', "## Índice por área"]
    for area, title in AREAS.items():
        subset = [i for i in cat["items"] if i["area"] == area]
        out += ["### " + title, table(["Item", "Categoria", "Método", "Avançado"], [[link(i["id"]), i["category"], i["method"], "Sim" if i.get("advanced") else "Não"] for i in subset])]
    out.append("## Perfis incluídos")
    for p in profiles:
        out += [f'### {p["name"]}', p.get("description", ""), f'ID: `{p["id"]}`. Tipo: `{p.get("kind", "usage")}`. ' + p.get("audience", ""), "Requisitos: " + json.dumps(p.get("requires", {}), ensure_ascii=False), table(["Item", "Opções explícitas"], [[link(id), json.dumps(sel.get("options", {}), ensure_ascii=False)] for id, sel in p["selections"].items()]) if p["selections"] else "Começa com escolhas vazias."]
    out.append("## Repositórios APT")
    for repo in cat["repos"]:
        out += [f'### {repo["name"]}', "Receita completa, incluindo URL da chave, fingerprints permitidos e documentação do fornecedor:", "```json\n" + json.dumps(repo, ensure_ascii=False, indent=2) + "\n```"]
    out.append("## Grupos exclusivos")
    for group in cat["groups"]:
        out += ["### " + group["name"], group["summary"] + " " + group["why"], ", ".join(link(i["id"]) for i in cat["items"] if i.get("group") == group["id"])]
    out.append("## Receitas dos itens")
    for it in cat["items"]:
        out += [f'<a id="item-{it["id"]}"></a>', f'### {it["name"]}', it["summary"], it["why"]]
        rows = [["ID", it["id"]], ["Área / categoria", AREAS[it["area"]] + " / " + it["category"]], ["Método", it["method"]], ["Fonte", it["source"]]]
        for key, label in [("packages", "Pacotes"), ("flatpak", "Flatpak"), ("extension", "UUID"), ("repo", "Repositório"), ("verified", "Evidência da receita"), ("sensitive", "Atenção"), ("confirm", "Confirmação")]:
            if it.get(key):
                rows.append([label, ", ".join(it[key]) if isinstance(it[key], list) else it[key]])
        if it.get("requires"):
            rows.append(["Dependências", ", ".join(map(link, it["requires"]))])
        if it.get("npm"):
            npm = it["npm"]
            rows += [["Pacote npm fixo", npm["package"] + "@" + npm["version"]], ["Comando", "~/.local/bin/" + npm["binary"]], ["Destino", "~/.local (utilizador, sem sudo)"]]
        if it.get("conflicts"):
            rows.append(["Conflitos", "; ".join(link(c["id"]) + (" (bloqueia): " if c["hard"] else " (aviso): ") + c["reason"] for c in it["conflicts"])])
        rows += [[label, "Sim"] for key, label in [("advanced", "Avançado"), ("service", "Serviço local"), ("relogin", "Terminar sessão"), ("reboot", "Reiniciar")] if it.get(key)]
        out.append(table(["Campo", "Valor"], rows))
        for option in it.get("options", []):
            out += [f'Opção **{option["label"]}** (`{option["key"]}`): ' + ("múltiplos valores" if option.get("multi") else "um valor") + ". Padrão: " + ", ".join(option["default"]) + ". " + option.get("help", ""), table(["Valor", "Descrição"], [[v["value"], v["label"]] for v in option["values"]])]
        if it.get("notes"):
            out.append(it["notes"])
        if it.get("manual"):
            out += ["#### Guia manual", "\n".join(f'{i + 1}. {text}' for i, text in enumerate(it["manual"]["steps"]))]
            if it["manual"].get("url"):
                out.append(f'[Documentação do fornecedor]({it["manual"]["url"]})')
        # Preserve every field, including less common detection and recipe data.
        out += ["<details>\n<summary>Receita JSON integral</summary>\n\n```json\n" + json.dumps(it, ensure_ascii=False, indent=2) + "\n```\n\n</details>"]
    return "\n\n".join(out) + "\n"


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--check", action="store_true", help="recusa documentação desatualizada sem alterar ficheiros")
    args = parser.parse_args()
    version = read("frontend/package.json")["version"]
    profiles = [json.loads(p.read_text()) for p in sorted((ROOT / "internal/profiles").glob("*.json"))]
    outputs = {"docs/manual.md": manual_doc(read("docs/manual.json"), version), "docs/catalogo.md": catalog_doc(read("internal/catalog/catalog.json"), profiles)}
    stale = []
    for name, content in outputs.items():
        path = ROOT / name
        if args.check:
            if not path.exists() or path.read_text() != content:
                stale.append(name)
        else:
            path.write_text(content)
            print(name)
    if stale:
        print("Documentação desatualizada: " + ", ".join(stale) + ". Execute python3 build/scripts/generate-docs.py.", file=sys.stderr)
        return 1
    return 0


if __name__ == "__main__":
    sys.exit(main())
