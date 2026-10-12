import test from "node:test";
import assert from "node:assert/strict";
import { readFile } from "node:fs/promises";
import ts from "typescript";

const source = await readFile(new URL("../src/catalogSearch.ts", import.meta.url), "utf8");
const compiled = ts.transpileModule(source, { compilerOptions: { target: ts.ScriptTarget.ES2022, module: ts.ModuleKind.ES2022 } }).outputText;
const { matches, searchScore } = await import("data:text/javascript;base64," + Buffer.from(compiled).toString("base64"));
const editor = { id: "vscode", name: "Visual Studio Code", summary: "Editor da Microsoft", why: "Programação", category: "Editores", packages: ["code"] };

test("common editor names find the editor before an unrelated recipe that mentions it", () => {
  const watchers = { id: "inotify-watches", name: "Mais watchers inotify", summary: "Limite de ficheiros", why: "Se o VS Code/Vite avisar", category: "Sistema" };
  for (const query of ["VS Code", "vscode", "Visual Studio Code", " CODE "]) {
    assert.equal(matches(editor, query), true, query);
    assert.ok(searchScore(editor, query) > searchScore(watchers, query), query);
  }
});

test("search supports accents, IDs, package identifiers and all query words", () => {
  const item = { id: "lang-pt", name: "Português (PT/BR)", summary: "Traduções e dicionários", why: "Idioma", category: "Idiomas", packages: ["language-pack-pt"] };
  assert.equal(matches(item, "portugues"), true);
  assert.equal(matches(item, "lang-pt"), true);
  assert.equal(matches(item, "language-pack-pt"), true);
  assert.equal(matches(item, "traducoes dicionarios"), true);
  assert.equal(matches(item, "traducoes docker"), false);
  assert.equal(matches(item, "   "), true);
});
