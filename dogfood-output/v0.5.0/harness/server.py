"""Local-only test server; exposes original Go domain logic with a Sim executor."""
import json
import subprocess
from http.server import SimpleHTTPRequestHandler, ThreadingHTTPServer


class Handler(SimpleHTTPRequestHandler):
    def do_POST(self):
        if self.path != '/ux-domain':
            self.send_error(404)
            return
        data = self.rfile.read(int(self.headers.get('Content-Length', '0')))
        process = subprocess.run(['/domain/ux-domain'], input=data, capture_output=True, timeout=60)
        self.send_response(200 if process.returncode == 0 else 400)
        self.send_header('Content-Type', 'application/json; charset=utf-8')
        self.end_headers()
        self.wfile.write(process.stdout if process.returncode == 0 else json.dumps({'error': process.stderr.decode()}).encode())


ThreadingHTTPServer(('0.0.0.0', 8765), Handler).serve_forever()
