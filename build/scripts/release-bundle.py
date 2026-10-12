#!/usr/bin/env python3
"""Package the candidate and its corresponding sources for direct testing."""
import hashlib
import json
import os
from pathlib import Path
import subprocess
import time
import zipfile

ROOT = Path(__file__).resolve().parents[2]
version = os.environ.get('VERSION', '1.0.0-rc.1')
dist = ROOT/'dist'
names = [f'meu-ubuntu-{version}-linux-amd64', f'meu-ubuntu_{version}_amd64.deb',
         f'meu-ubuntu-{version}-source.tar.gz']
manifest = {'version':version, 'platform':'Ubuntu 26.04 LTS / GNOME / amd64',
            'license':'GPL-3.0-only', 'publication':'local candidate',
            'go':subprocess.check_output(['go','version'],text=True).strip(),
            'artifacts':{name:hashlib.sha256((dist/name).read_bytes()).hexdigest() for name in names}}
(dist/'release-manifest.json').write_text(json.dumps(manifest,ensure_ascii=False,indent=2)+'\n')
instructions = f'''Meu Ubuntu {version} — candidata para testes

Destino: Ubuntu 26.04 LTS, GNOME, amd64/x86_64.
Antes de executar, confirme na pasta extraída: sha256sum -c SHA256SUMS

Instalação recomendada:
  sudo apt install ./meu-ubuntu_{version}_amd64.deb
  meu-ubuntu --simular

Alternativa sem instalar o pacote (precisa de GTK/WebKitGTK/pkexec no sistema):
  chmod +x meu-ubuntu-{version}-linux-amd64
  ./meu-ubuntu-{version}-linux-amd64 --simular

Abra como utilizador normal. Depois de conhecer os fluxos, abra sem --simular
para consultar o estado real; operações só ocorrem após revisão e confirmação.
Use uma VM/máquina de teste com backup para experiências administrativas.

PDF/CSV/Excel: Diagnóstico > Exportar relatório; escolha secções e formato.
Consulte TESTAR.md e MANUAL.md. Reporte versão, ambiente, passos e resultados,
removendo segredos/dados pessoais. GNOME/Wayland, polkit, atualização e feedback
externo ainda precedem a 1.0 estável.

GPLv3 e código-fonte correspondente incluídos. Não houve publicação externa.
'''
bundle = dist/f'meu-ubuntu-{version}.zip'
epoch = max(315532800,int(os.environ.get('SOURCE_DATE_EPOCH','1790000000')))
date = time.gmtime(epoch)[:6]
files = {name:dist/name for name in names}
files.update({'SHA256SUMS':dist/'SHA256SUMS','release-manifest.json':dist/'release-manifest.json',
              'LICENSE':ROOT/'LICENSE','NOTICE.md':ROOT/'NOTICE.md',
              'MANUAL.md':ROOT/'docs/manual.md','TESTAR.md':ROOT/'docs/beta-usabilidade.md',
              'CHANGELOG.md':ROOT/'CHANGELOG.md','SECURITY.md':ROOT/'SECURITY.md'})
with zipfile.ZipFile(bundle,'w',compression=zipfile.ZIP_DEFLATED) as archive:
    for name,path in files.items():
        info=zipfile.ZipInfo('meu-ubuntu-'+version+'/'+name,date_time=date)
        info.compress_type=zipfile.ZIP_DEFLATED
        info.external_attr=(0o100755 if name.endswith('-linux-amd64') else 0o100644)<<16
        archive.writestr(info,path.read_bytes())
    info=zipfile.ZipInfo('meu-ubuntu-'+version+'/LEIA-ME.txt',date_time=date)
    info.compress_type=zipfile.ZIP_DEFLATED
    info.external_attr=0o100644<<16
    archive.writestr(info,instructions.encode('utf-8'))
with (dist/'SHA256SUMS').open('a') as checksum:
    checksum.write(hashlib.sha256(bundle.read_bytes()).hexdigest()+'  '+bundle.name+'\n')
print(f'Release bundle: {bundle.name} ({bundle.stat().st_size} bytes).')
