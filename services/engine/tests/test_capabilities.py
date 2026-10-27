import json
import subprocess
import sys
import unittest

from observatorio_engine.__main__ import capabilities


class CapabilitiesTest(unittest.TestCase):
    def test_no_models_available(self):
        self.assertIs(capabilities()["electoral_models"], False)
        self.assertEqual(capabilities()["status"], "scaffold")

    def test_cli_reports_capabilities(self):
        result = subprocess.run(
            [sys.executable, "-m", "observatorio_engine"],
            check=True, capture_output=True, text=True,
        )
        self.assertEqual(json.loads(result.stdout), capabilities())
