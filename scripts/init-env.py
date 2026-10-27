"""Create local configuration once, without overwriting existing user values."""
import base64
import secrets
from pathlib import Path

path = Path(".env")
if path.exists():
    raise SystemExit(".env already exists; kept unchanged")
key = base64.b64encode(secrets.token_bytes(32)).decode()
with path.open("x") as output:
    output.write(Path(".env.example").read_text())
    output.write(f"APP_KEY=base64:{key}\n")
path.chmod(0o600)
print("Local .env created (values omitted).")
