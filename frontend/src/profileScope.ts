import type { Item, Profile } from "./api";

export function profileKey(profile: Profile) {
  return `${profile.source || (profile.builtin ? "builtin" : "mine")}:${profile.id || profile.name}`;
}

// A view filter only. Include dependencies and declared alternatives so a
// profile can be explored without hiding the choices within a group.
export function profileScope(profile: Profile | null, items: Item[]): Set<string> | null {
  if (!profile) return null;
  const byID = new Map(items.map((item) => [item.id, item]));
  const scope = new Set<string>();
  const include = (id: string) => {
    const item = byID.get(id);
    if (!item || scope.has(id)) return;
    scope.add(id);
    (item.requires || []).forEach(include);
    if (item.group) items.filter((other) => other.group === item.group).forEach((other) => include(other.id));
  };
  Object.keys(profile.selections || {}).forEach(include);
  return scope;
}
