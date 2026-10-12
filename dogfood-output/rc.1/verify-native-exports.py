from pathlib import Path
import csv
import json
import shutil
import subprocess
import openpyxl

out = Path('/src/dogfood-output/rc.1/exports/native')
out.mkdir(parents=True, exist_ok=True)
for extension in ('pdf', 'csv', 'xlsx'):
    source = Path('/tmp/rc-diagnostico.' + extension)
    assert source.is_file(), source
    assert source.stat().st_mode & 0o777 == 0o600
    shutil.copyfile(source, out / source.name)
with (out/'rc-diagnostico.csv').open(encoding='utf-8-sig', newline='') as stream:
    rows = list(csv.reader(stream, delimiter=';'))
workbook = openpyxl.load_workbook(out/'rc-diagnostico.xlsx')
sheet = workbook.active
xrows = [list(row) for row in sheet.values]
assert rows == xrows, (len(rows),len(xrows))
assert all(cell.data_type != 'f' for row in sheet for cell in row)
pdf = subprocess.check_output(['pdftotext', str(out/'rc-diagnostico.pdf'), '-'], text=True)
assert '1.0.0-rc.1' in pdf and 'simulação' in pdf and 'Bateria' in pdf
assert any('1.0.0-rc.1' in str(row) for row in rows)
assert any('simulação' in str(row) for row in rows)
result = {'nativeSaveDialog': True, 'mode':'simulação', 'csvRows':len(rows),
          'xlsxRows':len(xrows), 'csvMatchesXlsx':True, 'xlsxContainsFormulas':False,
          'permissions':'0600', 'pdfVersionAndMode':True}
(out/'validation.json').write_text(json.dumps(result,ensure_ascii=False,indent=2)+'\n')
subprocess.run(['pdftoppm','-f','1','-singlefile','-scale-to','1500','-png',str(out/'rc-diagnostico.pdf'),str(out/'pdf-page')],check=True)
print(json.dumps(result,ensure_ascii=False))
