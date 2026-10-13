#!/usr/bin/env python3
"""Drive the native simulation through AT-SPI in an isolated X11 container."""
import json
import os
from pathlib import Path
import subprocess
import time
import pyatspi

OUT = Path('/out')
LOG = (OUT/'app.log').open('a')
app = subprocess.Popen(['/src/dist/meu-ubuntu','--simular'], stdout=LOG, stderr=LOG)


def nodes(obj=None, depth=0):
    if obj is None: obj = pyatspi.Registry.getDesktop(0)
    if depth > 32: return
    yield obj
    try:
        for child in obj:
            yield from nodes(child, depth+1)
    except Exception:
        pass


def wait(check, seconds=20):
    deadline = time.monotonic()+seconds
    while time.monotonic() < deadline:
        result = check()
        if result: return result
        time.sleep(.1)
    raise AssertionError('Expected native UI state was not reached')


def button(name):
    for obj in nodes():
        try:
            if obj.name == name and obj.getRoleName() in ('button','push button','check box'):
                if obj.getState().contains(pyatspi.STATE_ENABLED): return obj
        except Exception:
            continue


def contains_text(text):
    for obj in nodes():
        try:
            if text in (obj.name or ''):return True
            if text in obj.queryText().getText(0,-1):return True
        except Exception:
            continue
    return False


def click(name):
    obj=wait(lambda: button(name))
    action=obj.queryAction()
    assert action.doAction(0), name
    time.sleep(.2)


def file_dialog(path):
    time.sleep(.6)
    title='Escolher pasta do projeto' if Path(path).is_dir() else 'Instalar pacote de integração'
    windows=subprocess.check_output(['xdotool','search','--onlyvisible','--name',title],text=True).splitlines()
    subprocess.run(['xdotool','windowfocus','--sync',windows[-1]],check=True,timeout=5)
    subprocess.run(['xdotool','key','ctrl+l'],check=True)
    def focused_entry():
        for obj in nodes():
            if obj.getState().contains(pyatspi.STATE_FOCUSED) and obj.getState().contains(pyatspi.STATE_EDITABLE):
                return obj
    entry=wait(focused_entry)
    assert entry.queryEditableText().setTextContents(str(path)), 'file chooser path'
    if Path(path).is_dir():
        rect=entry.queryComponent().getExtents(pyatspi.DESKTOP_COORDS)
        subprocess.run(['xdotool','mousemove',str(rect.x+30),str(rect.y+rect.height//2),'click','1'],check=True)
        subprocess.run(['xdotool','key','Return'],check=True)
        time.sleep(.6)
    click('Open')
    time.sleep(.6)


def shot(name):
    subprocess.run(['import','-window','root',str(OUT/(name+'.png'))],check=True)


def job():
    path=OUT/'data/meu-ubuntu-sim/setup-last.json'
    try:return json.loads(path.read_text())
    except (FileNotFoundError,ValueError):return {}


try:
    wait(lambda: button('Pacotes de integração'),40)
    shot('01-home')
    click('Pacotes de integração')
    click('Instalar pacote…')
    file_dialog('/src/dist/officegest-0.2.0.zip')
    click('Confirmei que este pacote veio da minha equipa.')
    shot('02-package-review')
    click('Instalar pacote')
    # Project selection is a persisted synthetic fixture; the GTK directory
    # chooser is outside this smoke's acceptance scope.
    click('2. Rever preparação do projeto')
    click('Confio neste checkout de desenvolvimento e revi os serviços e alterações.')
    shot('03-setup-review')
    click('Preparar ambiente')
    click('Pausar')
    wait(lambda:job().get('state')=='paused')
    before=job()['steps']
    click('Aplicações')
    time.sleep(.8)
    assert job()['steps']==before, 'steps advanced while paused'
    wait(lambda: contains_text('Em pausa'))
    assert button('Retomar'), 'pause controls disappeared after navigation'
    shot('04-paused-other-page')
    click('Retomar')
    wait(lambda:job().get('finished'),20)
    assert job()['state']=='done', job()
    click('OfficeGest')
    shot('05-completed')
    assert sorted(p.name for p in (OUT/'project').iterdir())==['sentinel'], 'simulation changed project'
    print('PASS: native ZIP import, setup review, pause/resume, navigation and completion; project path seeded')
except Exception:
    shot('failure')
    print('Visible buttons:', sorted(set(obj.name for obj in nodes() if obj.getRoleName() in ('button','push button','check box'))))
    raise
finally:
    app.terminate()
    try:app.wait(timeout=5)
    except subprocess.TimeoutExpired:app.kill();app.wait()
