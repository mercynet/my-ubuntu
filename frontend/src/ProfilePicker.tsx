import { useRef, useState } from "react";
import { Bookmark, Loader2 } from "lucide-react";
import type { Ctx } from "./App";
import type { Profile } from "./api";
import { api, errText } from "./api";
import { profileKey } from "./profileScope";

export function ProfilePicker({ ctx, disabled }: { ctx: Ctx; disabled: boolean }) {
  const [busy, setBusy] = useState(false);
  const request = useRef(0);
  const profile = ctx.catalogProfile;
  const sections: [string, Profile[]][] = [
    ["Perfis incluídos", ctx.data.profiles.filter((p) => p.source === "builtin" && p.kind !== "complement" && Object.keys(p.selections || {}).length > 0)],
    ["Complementos", ctx.data.profiles.filter((p) => p.source === "builtin" && p.kind === "complement")],
    ["Meus perfis", ctx.data.profiles.filter((p) => p.source === "mine")],
    ["Perfis importados", ctx.data.profiles.filter((p) => p.source === "imported")],
  ];
  const review = async () => {
    if (!profile || busy || disabled) return;
    const token = ++request.current;
    setBusy(true);
    try {
      const result = await api.reviewProfile(profile);
      if (request.current === token) ctx.openReview(result);
    } catch (e) { if (request.current === token) ctx.toast(errText(e), true); }
    finally { if (request.current === token) setBusy(false); }
  };
  return <div className="profile-picker">
    <label className="profile-select"><span>Mostrar catálogo de</span>
      <select aria-label="Mostrar catálogo de" value={profile ? profileKey(profile) : ""} disabled={busy} onChange={(e) => ctx.showCatalogProfile(e.target.value)}>
        <option value="">Todo o catálogo</option>
        {sections.map(([title, profiles]) => profiles.length > 0 && <optgroup key={title} label={title}>{profiles.map((p) => <option key={profileKey(p)} value={profileKey(p)}>{p.name}</option>)}</optgroup>)}
      </select>
    </label>
    {profile && <button className="btn" disabled={disabled || busy || !!ctx.editing} onClick={review}>{busy ? <Loader2 size={14} className="spin" /> : <Bookmark size={14} />} Adicionar este perfil…</button>}
  </div>;
}
