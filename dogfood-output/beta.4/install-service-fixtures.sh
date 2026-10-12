#!/bin/sh
# Only in the disposable Ubuntu 26.04 QA container.
set -eu
export DEBIAN_FRONTEND=noninteractive
printf '#!/bin/sh\nexit 101\n' > /usr/sbin/policy-rc.d
chmod 755 /usr/sbin/policy-rc.d
apt-get install -y --no-install-recommends apache2 nginx redis-server rabbitmq-server postgresql memcached > /tmp/services-install.log 2>&1
dpkg-query -W -f='${Package}\t${Version}\t${db:Status-Status}\n' apache2 nginx redis-server rabbitmq-server postgresql memcached
/usr/sbin/apache2 -v
/usr/sbin/nginx -v
/usr/bin/redis-server --version
/usr/bin/psql --version
/usr/bin/memcached -V
nginx -t
apache2ctl configtest
