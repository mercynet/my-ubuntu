// Test-only Wails adapter. Production frontend assets are copied unchanged.
// This adapter models data and operations; it does not validate the Go backend.
(() => {
  const { catalog, profiles } = window.__fixtureData;
  const listeners = {};
  const items = catalog.items;
  const byId = new Map(items.map(item => [item.id, item]));
  const installed = new Set(['git', 'curl', 'firefox', 'ubuntu-dock', 'appindicator', 'tiling-assistant']);
  const configured = new Set();
  const records = [];
  const calls = [];
  let prefs = JSON.parse(localStorage.getItem('ux-beta-prefs') || 'null') || {
    theme: 'light', selections: {}, showAdvanced: false, showAllAreas: false,
    origins: {}, onboarded: false, reverts: [], manualProgress: {}, catalogProfile: '',
  };
  let running = false;
  const recordCall = (name, args) => calls.push({ name, args, time: new Date().toISOString() });
  const domain = async (Method, parameters) => {
    const response = await fetch('/ux-domain', {method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({Method,...parameters})});
    const value = await response.json();
    if (!response.ok) throw new Error(value.error || 'Domain test bridge failed');
    return value;
  };
  const state = (item) => ({
    status: installed.has(item.id) ? 'installed' : configured.has(item.id) ? 'applied' : 'available',
    detail: installed.has(item.id) ? 'Presente no sistema simulado' : configured.has(item.id) ? 'Configurado no sistema simulado' : 'Disponível',
  });
  const review = (profile) => ({ profile, warnings: [], requirement: [], existsMine: profiles.some(p => p.id === profile.id && p.source === 'mine'), existsImported: false,
    items: Object.keys(profile.selections || {}).map(id => ({ id, name: byId.get(id)?.name || id, area: byId.get(id)?.area || '', ...state(byId.get(id) || {id}) })) });
  function buildPlan(selections, reverts = []) {
    const selection = structuredClone(selections), added = [];
    const include = id => {
      const item = byId.get(id);
      if (!item) return;
      for (const dep of item.requires || []) if (!selection[dep]) { selection[dep] = {}; added.push(dep); include(dep); }
    };
    Object.keys(selection).forEach(include);
    const steps = [], manual = [], kept = [], conflicts = [];
    for (const [id, selected] of Object.entries(selection)) {
      const item = byId.get(id);
      if (!item) continue;
      if (installed.has(id) || configured.has(id)) { kept.push({item:id,name:item.name,reason:'Já presente; sem alteração',...state(item)}); continue; }
      if (item.method === 'manual') { manual.push({item:id,name:item.name,steps:item.manual?.steps || [],url:item.manual?.url || ''}); continue; }
      steps.push({id:`step-${id}`,item:id,title:item.name,kind:['apt','apt-repo','admin'].includes(item.method)?'admin':'user',op:item.method,
        detail:[...(item.packages || []), ...(item.settings || []).map(s=>`${s.schema} ${s.key}`)],category:item.category,options:selected.options||{},confirm:item.confirm||'',expect:''});
      for (const conflict of item.conflicts || []) if (selection[conflict.id]) conflicts.push({a:id,b:conflict.id,hard:!!conflict.hard,reason:conflict.reason});
    }
    return {reviewToken:'test-fixture-token',steps,manual,kept,conflicts,blocked:conflicts.some(c=>c.hard),added,skipped:[],warnings:[],selection,reverts,
      relogin:Object.keys(selection).some(id=>byId.get(id)?.relogin),reboot:Object.keys(selection).some(id=>byId.get(id)?.reboot),needsRoot:steps.some(s=>s.kind==='admin')};
  }
  const methods = {
    Init: () => ({version:'0.6.0-beta.1',system:{id:'ubuntu',release:'26.04',codename:'resolute',pretty:'Ubuntu 26.04 LTS',arch:'amd64',gnome:'50',gnomeMajor:50,session:'wayland',simulated:true},
      prefs,profiles,interrupted:null,repos:catalog.repos,groups:catalog.groups,items,paths:{config:'/test/config',data:'/test/data',state:'/test/state'},pendingImport:null,pendingError:''}),
    States: selections => domain('States', {Selections:selections}),
    ItemState: id => state(byId.get(id)),
    SavePrefs: value => { prefs=structuredClone(value);localStorage.setItem('ux-beta-prefs',JSON.stringify(prefs)); },
    Profiles: () => profiles,
    ReviewProfile: review,
    BuildPlan: (selections,reverts) => domain('BuildPlan', {Selections:selections,Reverts:reverts}),
    MergeSelections: (current,origins,incoming,label) => {
      const selections=structuredClone(current), resultOrigins=structuredClone(origins||{}), added=[], conflicts=[];
      for(const [id, value] of Object.entries(incoming)) {
        const item=byId.get(id), other=Object.keys(selections).find(key => item?.group && key !== id && byId.get(key)?.group === item.group);
        if(other) { conflicts.push({kind:'group',item:other,name:byId.get(other).name,current:selections[other],currentOrigin:resultOrigins[other]||[],incomingItem:id,incomingName:item.name,incoming:value,incomingLabel:label,reason:'Alternativas exclusivas'});continue; }
        if(selections[id] && JSON.stringify(selections[id])!==JSON.stringify(value)) { conflicts.push({kind:'options',item:id,name:item?.name||id,current:selections[id],currentOrigin:resultOrigins[id]||[],incomingItem:id,incomingName:item?.name||id,incoming:value,incomingLabel:label,reason:'Opções diferentes'});continue; }
        if(!selections[id]) added.push(id);
        selections[id]=value;resultOrigins[id]=[...new Set([...(resultOrigins[id]||[]),label])];
      }
      return {selections,origins:resultOrigins,added,conflicts};
    },
    SaveProfile: (profile,overwrite) => {
      if(window.__ux.failSave) throw new Error('Não foi possível escrever o perfil (falha injetada)');
      const index=profiles.findIndex(p=>p.source==='mine'&&p.id===profile.id);
      if(index>=0&&!overwrite) throw new Error('Perfil já existe');
      const p={...structuredClone(profile),source:'mine',builtin:false};
      if(index>=0) profiles[index]=p;else profiles.push(p);
    },
    SaveImported: profile => profiles.push({...structuredClone(profile),source:'imported',builtin:false}),
    DeleteProfile: (source,id) => {const index=profiles.findIndex(p=>p.source===source&&p.id===id);if(index>=0)profiles.splice(index,1);},
    ExportProfile: () => '/test/perfil.json', ImportProfile: () => null,
    Runs: () => records, Changes: () => [], Running: () => running,
    RunLog: () => '[ponte de teste] Registo simulado',
    Execute: (selection,reverts) => {
      const plan=buildPlan(selection,reverts), id=`test-run-${records.length+1}`;
      const run={id,started:new Date().toISOString(),status:'running',selections:selection,reverts,steps:[],logFile:'/test/run.log',simulated:true};records.unshift(run);running=true;
      plan.steps.forEach((step,index)=>setTimeout(()=>{
        const failure=step.item===window.__ux.failItem;
        if(!failure) {if(['apt','apt-repo','flatpak','extension'].includes(step.op))installed.add(step.item);else configured.add(step.item);}
        const event={runId:id,stepId:step.id,status:failure?'failed':'done',error:failure?'exit status 100':'',time:new Date().toISOString(),
          problem:failure?{code:'apt',summary:'Não foi possível instalar o pacote',advice:'Verifique as fontes e tente novamente.'}:undefined};
        run.steps.push({...event,id:step.id,item:step.item,title:step.title});
        (listeners['run:step']||[]).forEach(fn=>fn(event));
        (listeners['run:log']||[]).forEach(fn=>fn({runId:id,stepId:step.id,time:event.time,line:failure?'E: pacote indisponível (falha injetada)':'Concluído (ponte de teste)'}));
        if(index===plan.steps.length-1) {running=false;run.status=run.steps.some(s=>s.status==='failed')?'failed':'done';run.finished=new Date().toISOString();(listeners['run:done']||[]).forEach(fn=>fn(run));}
      },(index+1)*700));
      return id;
    },
    Cancel: () => {running=false;}, EndSession: () => {},
    Diagnostics: () => ({facts:[{label:'Sistema',value:'Ubuntu 26.04 (ponte de teste)'},{label:'Sessão',value:'Wayland'}],disks:[],startup:[],suggestions:[],battery:false}),
    Firefox: () => ({roots:[],running:false}),SSHKeys: () => [],
    ParseInventory: () => [], OpenInventoryFile: () => '',
    About: () => ({version:'0.6.0-beta.1',author:'Meu Ubuntu',goVersion:'ponte de teste',wailsVersion:'v2.16.0',attributions:[]}),
  };
  window.__ux={calls,reset:()=>{localStorage.clear();location.reload();},failItem:'',failSave:false};
  window.go={main:{App:Object.fromEntries(Object.entries(methods).map(([name,fn])=>[name,async(...args)=>{recordCall(name,args);return structuredClone(await fn(...args));}]))}};
  window.runtime={EventsOn:(event,fn)=>{(listeners[event]||=[]).push(fn);return ()=>{listeners[event]=listeners[event].filter(f=>f!==fn);};},BrowserOpenURL:()=>{},ClipboardSetText:async text=>{window.__ux.clipboard=text;return true;}};
})();
