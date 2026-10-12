import type { Item, ItemState, Origins, Sel, Selections } from "./api";

export const optionSignature = (item: Item, selection?: Sel) =>
  JSON.stringify((item.options || []).map((option) => selection?.options?.[option.key] || option.default));

// Detection is asynchronous: a result for old options cannot select an item.
export function matchingStates(items: Record<string, Item>, next: Record<string, ItemState>, snapshot: Selections, current: Selections) {
  return Object.fromEntries(Object.entries(next).filter(([id]) => items[id] &&
    optionSignature(items[id], snapshot[id]) === optionSignature(items[id], current[id])));
}

export function includeInstalled(items: Item[], states: Record<string, ItemState>, selections: Selections, origins: Origins, excluded: string[], reverts: string[]) {
  const omitted = new Set([...excluded, ...reverts]);
  const occupied = new Set(items.filter((item) => item.group && selections[item.id]).map((item) => item.group));
  const added: string[] = [];
  for (const item of items) {
    if (item.method === "admin") continue; // Completed administrative actions are not installed applications.
    if (selections[item.id] || omitted.has(item.id) || !["installed", "provided"].includes(states[item.id]?.status)) continue;
    if (item.group && occupied.has(item.group)) continue;
    added.push(item.id);
    if (item.group) occupied.add(item.group);
  }
  if (!added.length) return { selections, origins, added };
  const next = { ...selections };
  const nextOrigins = { ...origins };
  added.forEach((id) => { next[id] = {}; nextOrigins[id] = ["estado"]; });
  return { selections: next, origins: nextOrigins, added };
}
