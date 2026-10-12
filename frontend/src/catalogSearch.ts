import type { Item } from "./api";

const aliases: Record<string, string[]> = {
  vscode: ["VS Code", "vscode"],
  "gnome-tweaks": ["ajustes", "tweaks"],
  "extension-manager": ["gestor de extensões"],
};

const normalize = (text: string) => text.normalize("NFD").replace(/[\u0300-\u036f]/g, "").toLowerCase().trim();

// Names and common aliases rank ahead of incidental mentions in descriptions.
export function searchScore(item: Item, query: string): number {
  const q = normalize(query);
  if (!q) return 1;
  const names = [item.name, item.id, ...(aliases[item.id] || [])].map(normalize);
  if (names.includes(q)) return 100;
  const words = q.split(/\s+/);
  if (names.some(name => words.every(word => name.includes(word)))) return 80;
  const identifiers = normalize([...(item.packages || []), item.flatpak || "", item.npm?.package || ""].join(" "));
  if (words.every(word => identifiers.includes(word))) return 60;
  const description = normalize([item.summary, item.why, item.category, identifiers].join(" "));
  return words.every(word => description.includes(word)) ? 20 : 0;
}

export const matches = (item: Item, query: string) => searchScore(item, query) > 0;
