"""Exercise local configuration without touching the checkout's .env."""
import base64
from pathlib import Path
import subprocess
import sys
import tempfile
import unittest

ROOT = Path(__file__).resolve().parents[1]


class InitEnvironmentTest(unittest.TestCase):
    def run_initializer(self, directory):
        return subprocess.run(
            [sys.executable, str(ROOT / 'scripts/init-env.py')],
            cwd=directory, capture_output=True, text=True,
        )

    def test_key_and_file_permissions(self):
        with tempfile.TemporaryDirectory() as directory:
            root = Path(directory)
            (root / '.env.example').write_text('POSTGRES_DB=example\n')
            result = self.run_initializer(directory)
            self.assertEqual(result.returncode, 0, result.stderr)
            config = (root / '.env').read_text()
            key = config.split('APP_KEY=base64:', 1)[1].strip()
            self.assertEqual(len(base64.b64decode(key, validate=True)), 32)
            self.assertEqual((root / '.env').stat().st_mode & 0o777, 0o600)
            self.assertNotIn(key, result.stdout + result.stderr)

    def test_existing_configuration_is_preserved(self):
        with tempfile.TemporaryDirectory() as directory:
            path = Path(directory) / '.env'
            path.write_text('EXISTING=user-configuration\n')
            result = self.run_initializer(directory)
            self.assertNotEqual(result.returncode, 0)
            self.assertEqual(path.read_text(), 'EXISTING=user-configuration\n')
