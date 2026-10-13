#!/usr/bin/env python3
"""Bundle a candidate and an optional integration for a colleague's test."""
import argparse
import hashlib
import json
import os
from pathlib import Path
import re
import time
import zipfile

ROOT = Path(__file__).resolve().parents[2]
parser = argparse.ArgumentParser(description=__doc__)
parser.add_argument('--manifest', type=Path, required=True)
parser.add_argument('--version', default=os.environ.get('VERSION', (ROOT/'VERSION').read_text().strip()))
parser.add_argument('--package-only', action='store_true')
args = parser.parse_args()
if not re.fullmatch(r'\d+\.\d+\.\d+-(?:rc|beta)\.\d+', args.version):
    parser.error('Use a candidate version, for example 1.0.0-rc.4')
data = args.manifest.read_bytes()
manifest = json.loads(data)
identifier, package_version = manifest['id'], manifest['version']
if (len(data) > 256 * 1024 or manifest.get('format') not in (1, 2)
        or (manifest.get('format') == 2) != bool(manifest.get('setup'))
        or not re.fullmatch(r'[a-z][a-z0-9-]{1,63}', identifier)
        or not re.fullmatch(r'\d+\.\d+\.\d+', package_version)):
    parser.error('Invalid integration identity or format')
dist = ROOT / 'dist'
dist.mkdir(exist_ok=True)
date = time.gmtime(max(315532800, int(os.environ.get('SOURCE_DATE_EPOCH', '1790000000'))))[:6]


def write_zip(path, contents, prefix=''):
    with zipfile.ZipFile(path, 'w', compression=zipfile.ZIP_DEFLATED) as archive:
        for name, payload in contents.items():
            info = zipfile.ZipInfo(prefix + name, date_time=date)
            info.compress_type = zipfile.ZIP_DEFLATED
            info.external_attr = 0o100644 << 16
            archive.writestr(info, payload)


package = dist / f'{identifier}-{package_version}.zip'
write_zip(package, {'manifest.json': data})
print(f'Integration: {package.name}')
if args.package_only:
    raise SystemExit(0)
deb = f'meu-ubuntu_{args.version}_amd64.deb'
source = f'meu-ubuntu-{args.version}-source.tar.gz'
contents = {name: (dist / name).read_bytes() for name in (deb, source, package.name)}
checksums = ''.join(f'{hashlib.sha256(payload).hexdigest()}  {name}\n'
                    for name, payload in contents.items())
instructions = f'''# Meu Ubuntu {args.version} + {manifest['name']} {package_version}

Candidata para teste em Ubuntu 26.04 LTS, GNOME, amd64. Use uma máquina de
teste ou VM. Este pacote inclui o instalador, fontes e integração opcional.

## Instalar e abrir

Extraia o ZIP, abra um terminal na pasta extraída e execute:

```sh
sha256sum -c SHA256SUMS
sudo apt install ./{deb}
meu-ubuntu --simular
```

Explore primeiro a simulação. Feche-a e abra `meu-ubuntu` para consultar o
sistema real e testar operações depois de rever o plano. Abra como utilizador
normal; o sistema pede autorização administrativa quando necessária.

## Instalar a integração {manifest['name']}

1. Abra **Pacotes de integração** e escolha **Instalar pacote…**.
2. Selecione `{package.name}`; confira nome, fornecedor e versão, marque a
   confirmação apresentada e confirme a instalação do pacote.
3. Abra o separador **{manifest['name']}**, escolha a pasta do projeto e execute
   o diagnóstico. Reabra a aplicação e confira que pacote e pasta persistem.
4. Para OfficeGest, tenha o checkout e uma configuração de desenvolvimento
   sanitizada fornecida pela equipa. O diagnóstico exige Docker/Compose 2.20+;
   Docker ausente ou configuração incompleta devem aparecer como pendências.
5. Se o pacote mostra **Preparar ambiente**, reveja primeiro as ferramentas,
   aplique o plano e volte ao separador. Reveja a preparação do projeto e
   confirme. Build, containers, schemas, Redis e Composer são etapas separadas.
   Pode pausar/retomar entre etapas; a operação ativa termina primeiro.
6. Importe os dumps sanitizados manualmente. Siga as instruções finais para
   /etc/hosts, certificado e autenticação privada do Composer, quando necessária.
   Use o comando dock apresentado no final para gerir só este ambiente.
   A preparação completa do OfficeGest num desktop limpo ainda exige validação.

## Roteiro e retorno

- Confira versão em **Sobre**, abertura, tema e janela pequena.
- Teste uma receita pequena: reveja, cancele a autorização, volte a rever e
  aplique autorizando. Confira o resultado no sistema; repita o plano concluído.
- Teste uma preferência reversível e restaure-a em **Histórico e reversão**.
- Em **Diagnóstico**, exporte PDF/CSV/Excel e teste cancelar o diálogo.
- Na integração, confira as pendências, seleção da pasta e cópia dos comandos.
- Durante execução, navegue entre separadores; confira etapa, contador e pausa.

Envie versão Ubuntu/GNOME, tipo de sessão, passos, resultado esperado/observado
e capturas úteis. Remova credenciais e dados pessoais. Relate se instalou do
zero ou atualizou outra versão. Esta candidata ainda não é a 1.0 estável.

Projeto: https://github.com/mercynet/my-ubuntu
Segurança: https://github.com/mercynet/my-ubuntu/security/advisories/new
GPLv3 e código-fonte correspondente acompanham este ZIP.
'''
contents['TESTAR.md'] = instructions.encode()
checksums += f'{hashlib.sha256(contents["TESTAR.md"]).hexdigest()}  TESTAR.md\n'
contents['SHA256SUMS'] = checksums.encode()
bundle = dist / f'meu-ubuntu-{args.version}-{identifier}-teste.zip'
write_zip(bundle, contents, bundle.stem + '/')
(dist / f'{bundle.name}.sha256').write_text(
    f'{hashlib.sha256(bundle.read_bytes()).hexdigest()}  {bundle.name}\n')
print(f'Test bundle: {bundle.name} ({bundle.stat().st_size} bytes)')
