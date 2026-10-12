from pathlib import Path
import hashlib
import os
import subprocess

root=Path('/src')
os.environ['VERSION']='1.0.0-rc.1'
subprocess.run(['python3',str(root/'build/scripts/source-archive.py')],check=True)
names=['meu-ubuntu-1.0.0-rc.1-linux-amd64','meu-ubuntu_1.0.0-rc.1_amd64.deb','meu-ubuntu-1.0.0-rc.1-source.tar.gz']
(root/'dist/SHA256SUMS').write_text(''.join(hashlib.sha256((root/'dist'/name).read_bytes()).hexdigest()+'  '+name+'\n' for name in names))
subprocess.run(['python3',str(root/'build/scripts/release-bundle.py')],check=True)
subprocess.run(['chown','-R','1000:1000',str(root/'dist'),str(root/'dogfood-output/rc.1')],check=True)
