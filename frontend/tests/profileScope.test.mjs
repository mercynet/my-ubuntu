import test from "node:test";
import assert from "node:assert/strict";
import { readFile } from "node:fs/promises";
import ts from "typescript";

const source = await readFile(new URL("../src/profileScope.ts", import.meta.url), "utf8");
const compiled = ts.transpileModule(source, { compilerOptions: { target: ts.ScriptTarget.ES2022, module: ts.ModuleKind.ES2022 } }).outputText;
const { profileScope, profileKey } = await import("data:text/javascript;base64," + Buffer.from(compiled).toString("base64"));

test("exploring a profile includes recursive dependencies and group alternatives without changing choices", () => {
  const items = [
    { id: "editor", requires: ["git", "missing"] },
    { id: "git", requires: ["editor"] }, // cycle must terminate
    { id: "balanced", group: "energy" },
    { id: "performance", group: "energy", requires: ["daemon"] },
    { id: "daemon" },
    { id: "unrelated" },
  ];
  const profile = { selections: { editor: { options: { version: ["1"] } }, balanced: {}, missing: {} } };
  const before = JSON.stringify({ items, profile });
  assert.deepEqual([...profileScope(profile, items)].sort(), ["balanced", "daemon", "editor", "git", "performance"]);
  assert.equal(JSON.stringify({ items, profile }), before);
  assert.equal(profileScope(null, items), null);
  assert.equal(profileScope({ selections: {} }, items).size, 0);
});

test("same profile ID from different origins stays distinguishable", () => {
  assert.notEqual(profileKey({ id: "terminal", source: "builtin" }), profileKey({ id: "terminal", source: "mine" }));
  assert.notEqual(profileKey({ id: "terminal", source: "mine" }), profileKey({ id: "terminal", source: "imported" }));
});
