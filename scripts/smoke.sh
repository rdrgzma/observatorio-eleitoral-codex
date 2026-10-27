#!/bin/sh
set -eu
cd "$(dirname "$0")/.."
docker compose up -d --build --wait --wait-timeout 120
python3 - <<'PY'
import json
import os
import time
import urllib.request
from pathlib import Path

# Read only the two non-secret port settings. Do not print .env contents.
ports = {'API_PORT': '8080', 'PANEL_PORT': '8000'}
for line in Path('.env').read_text().splitlines():
    key, _, value = line.partition('=')
    if key in ports:
        ports[key] = value
ports.update({k: os.environ[k] for k in ports if k in os.environ})
checks = [
    (ports['API_PORT'], '/healthz', {'status': 'ok', 'electoral_models': False}),
    (ports['PANEL_PORT'], '/', {'status': 'scaffold', 'electoral_models': False}),
    (ports['PANEL_PORT'], '/readyz', {'status': 'ok', 'database': 'reachable'}),
]
for port, path, expected in checks:
    for attempt in range(30):
        try:
            with urllib.request.urlopen(f'http://127.0.0.1:{port}{path}', timeout=3) as response:
                body = json.load(response)
            assert all(body.get(k) == v for k, v in expected.items()), body
            print(f'PASS {path}')
            break
        except (OSError, ValueError, AssertionError):
            if attempt == 29:
                raise
            time.sleep(1)
PY
