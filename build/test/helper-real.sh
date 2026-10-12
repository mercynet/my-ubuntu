#!/usr/bin/env bash
# Real installs in a disposable Ubuntu 26.04 container, through the shipped binary's restricted helper.
set -u
export PKEXEC_UID=1000
req() { echo "$1"; }
{
  req '{"op":"apt-update"}'
  req '{"op":"apt-install","item":"htop"}'
  req '{"op":"apt-install","item":"php","options":{"versions":["8.5"]}}'
  req '{"op":"apt-install","item":"xdebug","options":{"versions":["8.5"]}}'
  req '{"op":"post","item":"xdebug","options":{"versions":["8.5"]}}'
  req '{"op":"repo-add","item":"docker-engine"}'
  req '{"op":"apt-install","item":"docker-engine"}'
  req '{"op":"repo-add","item":"gh"}'
  req '{"op":"apt-install","item":"gh"}'
  req '{"op":"repo-add","item":"stripe-cli"}'
  req '{"op":"apt-install","item":"stripe-cli"}'
  req '{"op":"repo-add","item":"vscode"}'
  req '{"op":"repo-add","item":"brave"}'
  req '{"op":"repo-add","item":"warp"}'
  req '{"op":"repo-add","item":"php-ppa"}'
  req '{"op":"repo-add","item":"docker-engine"}'
  req '{"op":"admin","item":"journald-limit","options":{"value":["1G"]}}'
  req '{"op":"revert","item":"journald-limit"}'
  req '{"op":"admin","item":"docker-group"}'
  req '{"op":"shell","item":"htop"}'
  req '{"op":"apt-install","item":"php","options":{"versions":["8.5; reboot"]}}'
} | meu-ubuntu --helper 2>/dev/null | grep '"type":"done"' > /tmp/done.txt
paste -d' ' <(printf '%s\n' apt-update htop php8.5 xdebug xdebug-off repo:docker docker-ce repo:gh gh repo:stripe stripe repo:vscode repo:brave repo:warp repo:ondrej repo:docker-again journald revert-journald docker-group INVALID:shell INVALID:option) /tmp/done.txt
echo "--- verificação real"
dpkg -l htop php8.5-cli php8.5-xdebug docker-ce gh stripe 2>/dev/null | awk '/^ii/{print $2, $3}'
php -m | grep -i xdebug || echo "xdebug NÃO carregado (desativado como pretendido)"
ls /etc/php/8.5/cli/conf.d/ | grep -i xdebug || true
docker --version; gh --version | head -1; stripe version
ls /etc/apt/sources.list.d/ /etc/apt/keyrings/
cat /etc/apt/sources.list.d/meu-ubuntu-docker.sources
id -nG utilizador
ls /etc/systemd/journald.conf.d/ 2>/dev/null || echo "journald drop-in removido (revertido)"
apt-get -o DPkg::Lock::Timeout=5 update -qq 2>&1 | grep -iE "^(E|W):" || echo "apt update limpo (todas as fontes válidas)"
