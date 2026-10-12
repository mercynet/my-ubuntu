import { useState } from "react";
import { AppWindow, Code2, Settings2, ShieldCheck, Terminal, BatteryCharging, Server, Puzzle, FileText } from "lucide-react";
import type { Item } from "./api";
import mapping from "./appLogos.json";

const assets = import.meta.glob<string>("./assets/apps/*.svg", { eager: true, query: "?url", import: "default" });
const files: Record<string, string> = mapping;

export function AppIcon({ item }: { item: Item }) {
  const [failed, setFailed] = useState(false);
  const file = files[item.id];
  const src = file && assets["./assets/apps/" + file];
  const Icon = item.method === "extension" ? Puzzle
    : item.method === "template" ? FileText
    : ["gsettings", "admin"].includes(item.method) ? Settings2
    : item.method === "power-profile" || item.area === "energy" ? BatteryCharging
    : item.area === "infra" || item.category.startsWith("Serviços locais ·") ? Server
    : /seguran|privacidade/i.test(item.category) ? ShieldCheck
    : item.category === "Terminal" ? Terminal
    : item.area === "dev" ? Code2 : AppWindow;
  return <span className={"app-icon" + (src && !failed ? " branded" : "")} aria-hidden="true">
    {src && !failed ? <img src={src} width={32} height={32} alt="" draggable={false} onError={() => setFailed(true)} /> : <Icon size={26} />}
  </span>;
}
