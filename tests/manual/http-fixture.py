"""Serve a local DIPO HTTP 2 fixture for the documentation examples."""
import json
from http.server import BaseHTTPRequestHandler, ThreadingHTTPServer


class Handler(BaseHTTPRequestHandler):
    protocol_version = "HTTP/1.1"

    def reply(self, status, body, kind="application/json"):
        self.send_response(status)
        self.send_header("Content-Type", kind)
        self.send_header("Content-Length", str(len(body)))
        self.end_headers()
        if self.command != "HEAD":
            self.wfile.write(body)

    def do_GET(self):
        if self.path == "/health":
            self.reply(200, b'{"ok":true}')
        else:
            self.reply(404, b'{"error":"missing"}')

    do_HEAD = do_GET

    def do_OPTIONS(self):
        self.reply(204, b"")

    do_DELETE = do_OPTIONS

    def do_POST(self):
        length = int(self.headers.get("Content-Length", "0"))
        if not 0 < length <= 65535 or self.path != "/echo":
            self.reply(400, b'{"error":"invalid_request"}')
            return
        payload = self.rfile.read(length)
        try:
            json.loads(payload)
        except ValueError:
            self.reply(400, b'{"error":"invalid_json"}')
            return
        self.reply(201 if self.command == "POST" else 200, payload)

    do_PUT = do_POST
    do_PATCH = do_POST


if __name__ == "__main__":
    ThreadingHTTPServer(("127.0.0.1", 8765), Handler).serve_forever()
