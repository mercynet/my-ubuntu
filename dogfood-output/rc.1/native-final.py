import os
from pathlib import Path
import subprocess
import time

os.environ['DISPLAY']=':99'

def x(*args):
    subprocess.run(['xdotool',*map(str,args)],check=True)

def click(px,py):
    x('mousemove',px,py,'click',1)
    time.sleep(0.25)

window=subprocess.check_output(['xdotool','search','--onlyvisible','--class','.'],text=True).splitlines()[0]
x('windowmove',window,0,0)
click(379,648)
click(1026,155)
for index,extension in enumerate(('pdf','csv','xlsx')):
    if index:
        click(700,333)
        x('key','Down','Return')
        time.sleep(0.2)
    click(380,554)
    x('key','ctrl+a')
    x('type','--clearmodifiers','/tmp/rc-final-diagnostico.'+extension)
    click(1139,808)
    path=Path('/tmp/rc-final-diagnostico.'+extension)
    for attempt in range(20):
        if path.is_file(): break
        time.sleep(0.1)
    assert path.is_file(),extension
    Path('/tmp/rc-diagnostico.'+extension).write_bytes(path.read_bytes())
subprocess.run(['python3','/src/dogfood-output/rc.1/verify-native-exports.py'],check=True)
x('windowsize',window,900,600)
time.sleep(0.3)
click(62,499)
x('mousemove',880,25)
time.sleep(0.4)
subprocess.run(['import','-window',window,'/src/dogfood-output/rc.1/screenshots/native-sobre-900.png'],check=True)
print('Final candidate exported PDF/CSV/XLSX through GTK and opened About at 900x600.')
