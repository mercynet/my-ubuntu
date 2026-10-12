import test from "node:test";
import assert from "node:assert/strict";
import { readFile } from "node:fs/promises";
import ts from "typescript";

const source = await readFile(new URL("../src/installedSelection.ts", import.meta.url), "utf8");
const compiled = ts.transpileModule(source, { compilerOptions: { target: ts.ScriptTarget.ES2022, module: ts.ModuleKind.ES2022 } }).outputText;
const { includeInstalled, matchingStates } = await import("data:text/javascript;base64," + Buffer.from(compiled).toString("base64"));
const state = (status) => ({ status, detail: "fixture" });

test("all confirmed installed items are selected on startup, including manual and Ubuntu-provided items", () => {
  const items = [{ id: "backup", method: "apt" }, { id: "vpn", method: "apt" }, { id: "omz", method: "manual" }, { id: "dock", method: "extension" }, { id: "remove", method: "admin" }, { id: "missing" }, { id: "unknown" }];
  const next = includeInstalled(items, { backup: state("installed"), vpn: state("installed"), omz: state("installed"), dock: state("provided"), remove: state("installed"), missing: state("available"), unknown: state("unknown") }, {}, {}, [], []);
  assert.deepEqual(next.added, ["backup", "vpn", "omz", "dock"]);
  assert.deepEqual(Object.keys(next.selections), next.added);
  assert.deepEqual(next.origins.backup, ["estado"]);
  const repeat = includeInstalled(items, { backup: state("installed") }, next.selections, next.origins, [], []);
  assert.equal(repeat.selections, next.selections);
});

test("refresh includes newly discovered installations without undoing exclusions, reverts or pending options", () => {
  const selection = { php: { options: { version: ["8.5"] } }, performance: {} };
  const origins = { php: ["manual"], performance: ["perfil:work"] };
  const items = [{ id: "php" }, { id: "backup" }, { id: "revert" }, { id: "new" }, { id: "performance", group: "energy" }, { id: "balanced", group: "energy" }];
  const result = includeInstalled(items, Object.fromEntries(items.map((item) => [item.id, state("installed")])), selection, origins, ["backup"], ["revert"]);
  assert.deepEqual(result.added, ["new"]);
  assert.deepEqual(result.selections.php, selection.php);
  assert.deepEqual(result.origins.php, ["manual"]);
  assert.equal(selection.new, undefined);
});

test("two installed exclusive alternatives never become two conflicting automatic choices", () => {
  const result = includeInstalled([{ id: "balanced", group: "energy" }, { id: "saver", group: "energy" }], { balanced: state("installed"), saver: state("installed") }, {}, {}, [], []);
  assert.equal(Object.keys(result.selections).length, 1);
});

test("stale option detection cannot auto-select a configuration for superseded options", () => {
  const items = { php: { id: "php", options: [{ key: "version", default: ["8.4"] }] }, backup: { id: "backup" } };
  const result = matchingStates(items, { php: state("installed"), backup: state("installed"), gone: state("installed") }, {}, { php: { options: { version: ["8.5"] } } });
  assert.deepEqual(Object.keys(result), ["backup"]);
});

test("the catalog groups complement items under real tools without treating the collection as exclusive", async () => {
  const catalog = JSON.parse(await readFile(new URL("../../internal/catalog/catalog.json", import.meta.url), "utf8"));
  const items = new Map(catalog.items.map((item) => [item.id, item]));
  const tools = new Set(catalog.items.flatMap((item) => item.addonsFor || []));
  for (const id of ["zsh", "codex-cli", "claude-code", "opencode", "pi-coding-agent", "gemini-cli"]) assert.ok(tools.has(id));
  for (const item of catalog.items.filter((item) => item.addonsFor?.length)) {
    assert.equal(item.group, undefined);
    item.addonsFor.forEach((id) => assert.ok(items.has(id)));
  }
});
