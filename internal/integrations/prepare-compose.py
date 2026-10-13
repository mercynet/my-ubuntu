#!/usr/bin/env python3
"""Compiled development-environment operations; manifests supply data only."""
import argparse
import copy
import fcntl
import hashlib
import json
import os
from pathlib import Path
import re
import secrets
import shlex
import socket
import subprocess
import tempfile
import time

LABEL = 'io.meuubuntu.workspace'
TERMINAL = ('done', 'failed', 'intervention', 'cancelled')


def emit(**values):
    print(json.dumps(values, ensure_ascii=False), flush=True)


def run(argv, timeout=1800, capture=True):
    # Raw Docker/Composer output can contain credentials. Keep it out of IPC/logs.
    result = subprocess.run(argv, stdout=subprocess.PIPE if capture else subprocess.DEVNULL, stderr=subprocess.DEVNULL,
                            timeout=timeout, check=False)
    if result.returncode:
        raise RuntimeError(f'{argv[0]} {argv[1]} terminou com código {result.returncode}')
    return result.stdout.decode('utf-8', errors='replace') if capture else ''


def inside(root, path):
    value = Path(path)
    if not value.is_absolute():
        value = root / value
    if not value.resolve().is_relative_to(root.resolve()):
        raise ValueError('Caminho do projeto aponta para fora da pasta escolhida')
    return value


def write(path, text, mode=0o600):
    path = Path(path)
    if path.is_symlink() or any(p.is_symlink() for p in path.parents):
        raise ValueError('Ficheiro de preparação não pode ser um symlink')
    path.parent.mkdir(parents=True, exist_ok=True, mode=0o700)
    handle, temporary_name = tempfile.mkstemp(dir=path.parent, prefix='.'+path.name)
    temporary = Path(temporary_name)
    with os.fdopen(handle, 'w') as stream:
        stream.write(text)
    temporary.chmod(mode)
    temporary.replace(path)


def acquire_lock(path):
    path = Path(path)
    if any(p.is_symlink() for p in path.parents):
        raise ValueError('Pasta pessoal de preparação não pode ser symlink')
    path.parent.mkdir(parents=True, exist_ok=True, mode=0o700)
    fd = os.open(path, os.O_CREAT | os.O_RDWR | os.O_NOFOLLOW, 0o600)
    try:
        fcntl.flock(fd, fcntl.LOCK_EX | fcntl.LOCK_NB)
    except OSError:
        os.close(fd)
        raise ValueError('Já existe uma preparação ativa; aguarde o limite seguro')
    return os.fdopen(fd, 'r+')


def fingerprint(project, spec):
    digest = hashlib.sha256(json.dumps(spec, sort_keys=True).encode())
    digest.update(str(project).encode())
    paths = [project / spec['compose'], project / 'composer.json',
             project / 'composer.lock', project / '.env']
    for directory in ('docker', '_configs'):
        root = project / directory
        if root.exists():
            paths.extend(p for p in root.rglob('*') if p.is_file()
                         and 'logs' not in p.relative_to(root).parts)
    for path in sorted(set(paths)):
        if not path.exists():
            continue
        inside(project, path)
        if path.stat().st_size > 4 * 1024 * 1024:
            raise ValueError('Configuração demasiado grande para revisão')
        digest.update(str(path.relative_to(project)).encode())
        digest.update(hashlib.sha256(path.read_bytes()).digest())
    return digest.hexdigest()


def compose(project, file, namespace):
    return ['docker', 'compose', '--parallel', '1', '--project-directory', str(project),
            '-f', str(file), '--project-name', namespace]


def source_config(project, spec):
    if os.environ.get('DOCKER_HOST') and not os.environ['DOCKER_HOST'].startswith('unix://'):
        raise ValueError('Preparação exige um Docker local, com socket Unix')
    for name in ('composer.json', 'composer.lock', '.env', spec['template']):
        if not inside(project, name).is_file():
            raise ValueError('Falta configuração de desenvolvimento: peça .env, modelo de tenant e dependências à equipa')
    version = run(['docker', 'compose', 'version', '--short'], 15).strip()
    match = re.search(r'(\d+)\.(\d+)', version)
    if not match or tuple(map(int, match.groups())) < (2, 20):
        raise ValueError('Atualize Docker Compose para 2.20 ou superior')
    context = json.loads(run(['docker', 'context', 'inspect'], 15))[0]
    if not context['Endpoints']['docker']['Host'].startswith('unix://'):
        raise ValueError('Preparação exige um Docker local, com socket Unix')
    cfg = json.loads(run(compose(project, inside(project, spec['compose']), 'review') +
                         ['config', '--format', 'json'], 30))
    services = cfg.get('services', {})
    if not 1 <= len(services) <= 32:
        raise ValueError('Número de serviços não suportado')
    required = [spec['application'], spec['initializer']] + [x['service'] for x in spec['databases']] + [x['service'] for x in spec['redis']]
    if spec.get('worker'):
        required.append(spec['worker'])
    if any(name not in services for name in required):
        raise ValueError('O Compose não contém os serviços declarados pelo pacote')
    for service in services.values():
        if service.get('privileged') or service.get('network_mode') or service.get('pid') or service.get('ipc') == 'host' or service.get('devices') or service.get('cap_add') or service.get('volumes_from'):
            raise ValueError('Serviço pede privilégios não suportados pela preparação')
        if service.get('configs') or service.get('secrets') or service.get('external_links'):
            raise ValueError('Configurações externas exigem preparação manual')
        build = service.get('build')
        if build:
            inside(project, build['context'])
            inside(Path(build['context']), build.get('dockerfile', 'Dockerfile'))
        for volume in service.get('volumes', []):
            if volume['type'] == 'bind':
                inside(project, volume['source'])
    for value in cfg.get('volumes', {}).values():
        if value and (value.get('driver_opts') or value.get('driver', 'local') != 'local'):
            raise ValueError('Volumes especiais exigem preparação manual')
    for value in cfg.get('networks', {}).values():
        if value and (value.get('ipam') or value.get('driver_opts') or value.get('driver', 'bridge') != 'bridge'):
            raise ValueError('Redes especiais exigem preparação manual')
    return cfg


def free_port(preferred):
    with socket.socket() as listener:
        listener.bind(('127.0.0.1', preferred))
        return listener.getsockname()[1]


def choose_port(preferred):
    try:
        return free_port(preferred)
    except OSError:
        return free_port(0)


def preflight(cfg, namespace):
    # A deterministic namespace is reusable only if resources carry our marker.
    run(['docker', 'info', '--format', '{{.ServerVersion}}'], 15)
    checks = [('container', [namespace + '-' + name for name in cfg['services']]),
              ('volume', [namespace + '_' + name for name in cfg.get('volumes', {})]),
              ('network', [namespace + '_' + name for name in cfg.get('networks', {})])]
    owned = False
    for kind, names in checks:
        for name in names:
            result = subprocess.run(['docker', kind, 'inspect', name],
                                    stdout=subprocess.PIPE, stderr=subprocess.DEVNULL, timeout=15)
            if result.returncode:
                continue
            value = json.loads(result.stdout)[0]
            labels = value.get('Config', {}).get('Labels', {}) if kind == 'container' else value.get('Labels', {})
            if (labels or {}).get(LABEL) != namespace:
                raise ValueError('Já existe um recurso alheio com o nome reservado')
            owned = True
    return owned


def review(project, spec, workspace, selected_ports=None):
    cfg = source_config(project, spec)
    namespace = 'mu-og-' + hashlib.sha256(str(project).encode()).hexdigest()[:12]
    owned = preflight(cfg, namespace)
    saved = workspace / 'settings.json'
    if workspace.is_symlink() or saved.is_symlink():
        raise ValueError('Pasta pessoal de preparação não pode ser symlink')
    settings = json.loads(saved.read_text()) if saved.exists() else {}
    if owned and not settings:
        raise ValueError('Existem dados preparados mas falta a configuração pessoal; recupere settings.json antes de continuar')
    if settings and settings.get('namespace') != namespace:
        raise ValueError('Preparação anterior pertence a outro projeto')
    ports = dict(selected_ports or settings.get('ports', {}))
    if 'http' not in ports: ports['http'] = choose_port(8089)
    if 'https' not in ports: ports['https'] = choose_port(8449)
    for name, service in cfg['services'].items():
        for number, port in enumerate(service.get('ports', [])):
            key = name + ':' + str(number)
            if name == 'nginx' and port.get('target') in (80, 443):
                ports[key] = ports['http' if port['target'] == 80 else 'https']
            elif key not in ports:
                ports[key] = choose_port(0)
    steps = ['Preparar configuração pessoal', 'Construir imagens do projeto',
             'Iniciar containers', 'Criar/conferir schemas e preservar dados',
             'Conferir Redis Cluster', 'Instalar dependências Composer',
             'Conferir serviços e preparar instruções finais']
    hosts = ['localhost'] + [t['domain'] for t in spec['tenants']]
    manual = [
        'Importe os dumps de desenvolvimento nos schemas apresentados; esta preparação não importa dados.',
        'Nomes locais: confira /etc/hosts antes de copiar a linha seguinte para esse ficheiro com sudoedit /etc/hosts:\n127.0.0.1 ' + ' '.join(hosts[1:]),
        'Certificado local (opcional), execute:\nsudo install -m 644 ' + shlex.quote(str(workspace / 'ssl' / 'localhost.crt')) + ' /usr/local/share/ca-certificates/' + namespace + '.crt && sudo update-ca-certificates',
        'Configurações de tenants já existentes são preservadas; confira nelas os acessos às novas bases antes de importar os dumps.',
    ]
    binding = fingerprint(project, spec) + json.dumps(ports, sort_keys=True) + json.dumps(cfg, sort_keys=True)
    for name in ('settings.json', 'compose.json'):
        path = workspace / name
        if path.is_symlink(): raise ValueError('Configuração pessoal não pode ser symlink')
        if path.exists(): binding += hashlib.sha256(path.read_bytes()).hexdigest()
    token = hashlib.sha256(binding.encode()).hexdigest()
    return dict(token=token, project=str(project), namespace=namespace, services=list(cfg['services']),
                databases=spec['databases'], steps=steps, manual=manual, ports=ports,
                address=f'https://{spec["tenants"][0]["domain"] if spec["tenants"] else "localhost"}:{ports["https"]}', cfg=cfg, settings=settings)


def cluster_script(spec):
    nodes = [f'{n["service"]}:{n["port"]}' for n in spec['redis']]
    first = spec['redis'][0]
    endpoints = ' '.join(nodes)
    auth = 'export REDISCLI_AUTH="$REDIS_PASSWORD"\n'
    return '#!/bin/sh\nset -eu\n' + auth + f'''
for address in {endpoints}; do
  host=$(echo "$address" | cut -d: -f1); port=$(echo "$address" | cut -d: -f2)
  n=0; until redis-cli -h "$host" -p "$port" ping 2>/dev/null | grep -q PONG; do
    n=$((n+1)); test "$n" -lt 120 || exit 1; sleep 1
  done
done
if redis-cli -h {first['service']} -p {first['port']} cluster info | grep -q 'cluster_state:ok'; then exit 0; fi
for address in {endpoints}; do
  host=$(echo "$address" | cut -d: -f1); port=$(echo "$address" | cut -d: -f2)
  info=$(redis-cli -h "$host" -p "$port" cluster info)
  echo "$info" | grep -q 'cluster_slots_assigned:0' || exit 42
  echo "$info" | grep -q 'cluster_known_nodes:1' || exit 42
done
redis-cli --cluster create {endpoints} --cluster-replicas {1 if len(nodes)==6 else 0} --cluster-yes
'''


def php_string(value):
    return value.replace('\\', '\\\\').replace("'", "\\'")


def prepare(project, spec, workspace, details):
    if workspace.is_symlink():
        raise ValueError('Pasta pessoal de preparação não pode ser symlink')
    workspace.mkdir(parents=True, exist_ok=True, mode=0o700)
    settings = details['settings'] or dict(namespace=details['namespace'], password=secrets.token_hex(24))
    settings['ports'] = details['ports']
    # Preserve the same credentials when a failed preparation is reviewed again.
    write(workspace/'settings.json', json.dumps(settings, indent=2))
    password = settings['password']
    cfg = copy.deepcopy(details['cfg'])
    namespace = details['namespace']
    cfg['name'] = namespace
    for name, service in cfg['services'].items():
        old_name = service.get('container_name')
        service['container_name'] = namespace + '-' + name
        service.setdefault('labels', {})[LABEL] = namespace
        # Do not leave a newly prepared development stack running after reboot.
        service['restart'] = 'no'
        nets = service.setdefault('networks', {})
        for key, net in nets.items():
            if net is None: net = nets[key] = {}
            if old_name:
                net.setdefault('aliases', []).append(old_name)
        # Keep ports local and avoid fixed names/ports shared with other projects.
        for number, port in enumerate(service.get('ports', [])):
            key = name + ':' + str(number)
            if name == 'nginx' and port.get('target') in (80, 443):
                published = details['ports']['http' if port['target']==80 else 'https']
            else:
                published = settings['ports'][key]
            port['published'], port['host_ip'] = str(published), '127.0.0.1'
        if 'build' in service:
            service['image'] = namespace + '-' + name + ':dev'
    for name in (spec['application'], spec.get('worker')):
        if not name:
            continue
        service = cfg['services'][name]
        if name == spec.get('worker'):
            service['image'] = cfg['services'][spec['application']]['image']
        service.setdefault('environment', {}).update(
            REDIS_MODE='cluster', REDIS_CLUSTER_HOSTS=','.join(f'{n["service"]}:{n["port"]}' for n in spec['redis']),
            OG_REDIS_HOSTS=','.join(f'{n["service"]}:{n["port"]}' for n in spec['redis']), OG_REDIS_PASSW=password,
            WWWUSER=str(os.getuid()), WWWGROUP=str(os.getgid()), HTTP_HOST='localhost:'+str(details['ports']['http']))
        build = service.get('build')
        if build and isinstance(build, dict):
            build.setdefault('args', {}).update(WWWUSER=str(os.getuid()), WWWGROUP=str(os.getgid()))
    for kind in ('networks', 'volumes'):
        for name, value in cfg.get(kind, {}).items():
            if value is None: value = cfg[kind][name] = {}
            value.pop('external', None)
            value['name'] = namespace + '_' + name
            value.setdefault('labels', {})[LABEL] = namespace
    for database in spec['databases']:
        service = cfg['services'][database['service']]
        service.setdefault('environment', {}).update(MYSQL_ROOT_PASSWORD=password, MYSQL_USER='ogdevel', MYSQL_PASSWORD=password, MYSQL_DATABASE=database['schemas'][0])
        sql = "CREATE USER IF NOT EXISTS 'ogdevel'@'%' IDENTIFIED BY '" + password + "';\n"
        for schema in database['schemas']:
            quoted = chr(96) + schema + chr(96)
            sql += f'CREATE DATABASE IF NOT EXISTS {quoted} CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;\n'
            sql += f"GRANT ALL PRIVILEGES ON {quoted}.* TO 'ogdevel'@'%';\n"
        directory = workspace / database['service']
        write(directory / '00-init.sql', sql + 'FLUSH PRIVILEGES;\n', 0o644)
        service['volumes'] = [v for v in service.get('volumes', []) if v['target'] != '/docker-entrypoint-initdb.d']
        service['volumes'].append(dict(type='bind', source=str(directory/'00-init.sql'), target='/docker-entrypoint-initdb.d/00-init.sql', read_only=True))
        service['healthcheck'] = dict(test=['CMD-SHELL', 'MYSQL_PWD="$MYSQL_ROOT_PASSWORD" mariadb-admin -uroot ping --silent'], interval='5s', timeout='5s', retries=30)
    for node in spec['redis']:
        service = cfg['services'][node['service']]
        service.setdefault('environment', {})['REDIS_PASSWORD'] = password
        redis_argv = ['redis-server', '--cluster-enabled', 'yes', '--cluster-config-file', 'nodes.conf',
            '--cluster-node-timeout', '5000', '--appendonly', 'yes', '--port', str(node['port']),
            '--requirepass', password, '--masterauth', password,
            '--cluster-announce-port', str(node['port']), '--cluster-announce-bus-port', str(node['port']+10000)]
        service['command'] = ['sh', '-c', 'exec ' + shlex.join(redis_argv) + ' --cluster-announce-ip "$(hostname -i | cut -d\" \" -f1)"']
    script = workspace / 'init-cluster.sh'
    write(script, cluster_script(spec), 0o700)
    initializer = cfg['services'][spec['initializer']]
    initializer['volumes'] = [dict(type='bind', source=str(script), target='/init-cluster.sh', read_only=True)]
    initializer['entrypoint'] = ['sh', '/init-cluster.sh']
    initializer.setdefault('environment', {})['REDIS_PASSWORD'] = password
    initializer['restart'] = 'no'
    # Supply fresh local certificates to the repository's Nginx configuration.
    ssl = workspace / 'ssl'
    ssl.mkdir(exist_ok=True, mode=0o700)
    if not (ssl / 'localhost.crt').exists():
        hosts = ['localhost'] + [t['domain'] for t in spec['tenants']]
        run(['openssl', 'req', '-x509', '-nodes', '-newkey', 'rsa:2048',
             '-keyout', str(ssl/'localhost.key'), '-out', str(ssl/'localhost.crt'),
             '-days', '365', '-subj', '/CN=localhost', '-addext', 'subjectAltName='+','.join('DNS:'+h for h in hosts)], 30)
        (ssl/'localhost.key').chmod(0o600)
    for target, original in [('localhost.crt','localhost.crt'), ('localhost.key','localhost.key'),
                             ('_wildcard.test+2.pem','localhost.crt'), ('_wildcard.test+2-key.pem','localhost.key')]:
        if target != original:
            write(ssl/target, (ssl/original).read_text())
    if 'nginx' in cfg['services']:
        nginx = cfg['services']['nginx']
        nginx['volumes'] = [v for v in nginx.get('volumes', []) if v['target'].rstrip('/') != '/etc/nginx/ssl']
        nginx['volumes'].append(dict(type='bind', source=str(ssl), target='/etc/nginx/ssl/', read_only=True))
    template_path = inside(project, spec['template'])
    source = template_path.read_text()
    marker = 'config.php" <<EOF\n'
    if marker not in source:
        raise ValueError('Modelo de tenant não reconhecido; peça uma configuração à equipa')
    template = source.split(marker, 1)[1].split('\nEOF', 1)[0]
    for tenant in spec['tenants']:
        text = template
        values = dict(SELECTED_CONTAINER=tenant['service'], DB_USER='ogdevel', DB_PASS=password, NEW_SCHEMA=tenant['schema'])
        for key, value in values.items():
            text = text.replace('$'+'{'+key+'}', php_string(value))
        text = re.sub(r"\$\(echo [^)]+\)", tenant['schema'].replace('_', ''), text)
        text = text.replace('redis:6379', ','.join(f'{n["service"]}:{n["port"]}' for n in spec['redis']))
        if re.search(r'\$\{|\$\(', text):
            raise ValueError('Modelo de tenant tem campos não suportados; configuração manual necessária')
        for port in (details['ports']['http'], details['ports']['https']):
            path = project / '_configs' / (tenant['domain'] + ':' + str(port)) / 'config.php'
            inside(project, path)
            if not path.exists(): write(path, text + '\n')
    for service in cfg['services'].values():
        for volume in service.get('volumes', []):
            if volume['type'] == 'bind' and '/logs/' in volume['source']:
                inside(project, volume['source']).mkdir(parents=True, exist_ok=True)
    write(workspace/'settings.json', json.dumps(settings, indent=2))
    write(workspace/'compose.json', json.dumps(cfg, indent=2))
    launcher = '#!/bin/sh\nexec docker compose --parallel 1 --project-directory ' + shlex.quote(str(project)) + ' -f ' + shlex.quote(str(workspace/'compose.json')) + ' --project-name ' + namespace + ' "$@"\n'
    write(workspace/'dock', launcher, 0o700)
    return cfg


def boundary(control, project, spec, workspace=None, namespace=None):
    paused = False
    def snapshot():
        value = fingerprint(project, spec)
        if workspace:
            for name in ('settings.json', 'compose.json'):
                path = workspace/name
                if path.is_symlink(): raise ValueError('Configuração pessoal não pode ser symlink')
                if path.exists(): value += hashlib.sha256(path.read_bytes()).hexdigest()
        return value
    expected = snapshot()
    while Path(control + '.pause').exists():
        if not paused:
            emit(state='paused')
            paused = True
        if Path(control + '.cancel').exists():
            raise InterruptedError()
        time.sleep(0.1)
    if Path(control + '.cancel').exists():
        raise InterruptedError()
    if paused:
        if snapshot() != expected:
            raise ValueError('O projeto mudou durante a pausa; reveja a preparação')
        if workspace and (workspace/'compose.json').exists():
            preflight(json.loads((workspace/'compose.json').read_text()), namespace)
        emit(state='running')


def verify_services(prefix, cfg, spec):
    raw = run(prefix+['ps','--all','--format','json'],30)
    try:
        values = json.loads(raw)
        if not isinstance(values, list): values = [values]
    except json.JSONDecodeError:
        values = [json.loads(line) for line in raw.splitlines() if line.strip()]
    states = {v['Service']: v for v in values}
    for name in cfg['services']:
        value = states.get(name, {})
        if name == spec['initializer']:
            if value.get('State') != 'exited' or value.get('ExitCode') != 0:
                raise ValueError('Inicialização do cluster exige atenção; nenhum reset foi aplicado')
        elif value.get('State') != 'running' or value.get('Health', '') not in ('', 'healthy'):
            raise ValueError('Há serviços que não ficaram ativos; confira o terminal pessoal do projeto')


def apply(project, spec, workspace, details, control):
    cfg = None
    prefix = compose(project, workspace/'compose.json', details['namespace'])
    def database():
        for value in spec['databases']:
            argv = prefix + ['exec', '-T', value['service'], 'sh', '-c',
                'MYSQL_PWD="$MYSQL_ROOT_PASSWORD" mariadb -uroot < /docker-entrypoint-initdb.d/00-init.sql']
            for attempt in range(120):
                try:
                    run(argv, 15)
                    break
                except RuntimeError:
                    if attempt == 119:
                        raise
                    time.sleep(1)
            listed = run(prefix+['exec','-T',value['service'],'sh','-c',
                'MYSQL_PWD="$MYSQL_ROOT_PASSWORD" mariadb -uroot -N -e "SHOW DATABASES"'], 15)
            if any(name not in listed.splitlines() for name in value['schemas']):
                raise ValueError('Nem todas as bases foram criadas')
            run(prefix+['exec','-T',value['service'],'sh','-c',
                'MYSQL_PWD="$MYSQL_PASSWORD" mariadb -uogdevel -h127.0.0.1 -N -e "SELECT 1"'], 15)
    def redis():
        first = spec['redis'][0]
        argv = prefix + ['exec','-T',first['service'],'sh','-c',
            'REDISCLI_AUTH="$REDIS_PASSWORD" redis-cli -p '+str(first['port'])+' cluster info']
        for attempt in range(120):
            result = run(argv, 15)
            if 'cluster_state:ok' in result and 'cluster_slots_assigned:16384' in result:
                return
            time.sleep(1)
        raise ValueError('Redis Cluster não ficou pronto; dados existentes foram preservados')
    def configuration():
        nonlocal cfg
        cfg = prepare(project, spec, workspace, details)
        details['settings'] = json.loads((workspace/'settings.json').read_text())
    actions = [configuration,
               lambda: run(prefix+['build'], 3600, capture=False),
               lambda: run(prefix+['up','-d','--no-build'], 1800, capture=False),
               database, redis,
               lambda: run(prefix+['exec','-T','--user',f'{os.getuid()}:{os.getgid()}',spec['application'],
                                   'composer','install','--no-interaction','--prefer-dist','--no-progress'],1800,capture=False),
               lambda: verify_services(prefix, cfg, spec)]
    intervention = False
    for index, action in enumerate(actions):
        boundary(control, project, spec, workspace, details['namespace'])
        emit(index=index, status='running', title=details['steps'][index])
        try:
            action()
        except Exception:
            if index == 5:
                emit(index=index, status='intervention', message='Composer exige atenção. Use o terminal do projeto e a autenticação da equipa; não copie credenciais de produção.')
                intervention = True
                continue
            emit(index=index, status='failed', message='Não foi possível concluir este passo. Confira Docker, configuração e acesso à rede; os dados foram preservados.')
            raise
        emit(index=index, status='done', title=details['steps'][index])
    emit(status='intervention' if intervention else 'done', manual=details['manual'],
         address=details['address'], terminal=str(workspace/'dock'))


def main():
    os.umask(0o077)
    parser = argparse.ArgumentParser()
    parser.add_argument('--project', required=True)
    parser.add_argument('--spec', required=True)
    parser.add_argument('--workspace', required=True)
    parser.add_argument('--lock')
    parser.add_argument('--review', action='store_true')
    parser.add_argument('--token')
    parser.add_argument('--ports')
    parser.add_argument('--control', default='/nonexistent')
    args = parser.parse_args()
    project = Path(args.project).resolve()
    spec = json.loads(Path(args.spec).read_text())
    workspace = Path(args.workspace)
    lock = acquire_lock(args.lock or workspace/'setup.lock')
    details = review(project, spec, workspace, json.loads(args.ports) if args.ports else None)
    if args.review:
        details.pop('cfg')
        details.pop('settings')
        emit(**details)
        return
    if args.token != details['token']:
        raise ValueError('O projeto mudou desde a revisão; reveja antes de preparar')
    apply(project, spec, workspace, details, args.control)


if __name__ == '__main__':
    try:
        main()
    except InterruptedError:
        emit(status='cancelled', message='Cancelado entre operações; passos concluídos e dados preservados.')
        raise SystemExit(130)
    except Exception as error:
        # Do not disclose raw exceptions/config output containing project secrets.
        text = str(error) if isinstance(error, ValueError) else 'Preparação falhou; consulte a configuração e os requisitos.'
        emit(status='failed', message=text)
        raise SystemExit(1)
