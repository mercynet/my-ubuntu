import type { Item, SettingState } from "./api";

export function settingLabel(item: Item | undefined, setting: SettingState, value: string) {
  if (!value) return "(sem valor)";
  const settingRecipe = item?.settings?.find((s) => s.key === setting.key && (s.schema === setting.schema || s.schema.includes("{")));
  const recipe = settingRecipe || item?.git?.find((s) => s.key === setting.key);
  const option = item?.options?.find((o) => o.key === recipe?.option)
    || (item?.method === "admin" ? item.options?.find((o) => o.key === "value") : undefined);
  if (option) {
    const mapping = settingRecipe?.map;
    const found = option.values.find((v) => (mapping?.[v.value] || v.value) === value);
    if (found) return found.label;
  }
  if (value === "true") return "Ativado";
  if (value === "false") return "Desativado";
  return value;
}
