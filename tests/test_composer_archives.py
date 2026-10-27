import copy
import unittest

from scripts.normalize_composer_dist import normalize


class ArchiveIntegrityTest(unittest.TestCase):
    def package(self, checksum=''):
        commit = 'a' * 40
        return {'packages': [{'name': 'example/library', 'version': '1.0.0', 'dist': {
            'url': f'https://api.github.com/repos/example/library/zipball/{commit}',
            'reference': commit, 'shasum': checksum,
        }}]}

    def test_existing_checksum_and_archive_are_preserved(self):
        lock = self.package(checksum='existing-checksum')
        original = copy.deepcopy(lock)
        self.assertEqual(normalize(lock), 0)
        self.assertEqual(lock, original)

    def test_archive_stays_pinned_and_normalization_is_idempotent(self):
        lock = self.package()
        original = copy.deepcopy(lock['packages'][0])
        self.assertEqual(normalize(lock), 1)
        package = lock['packages'][0]
        self.assertEqual(package['version'], original['version'])
        self.assertEqual(package['dist']['reference'], original['dist']['reference'])
        self.assertTrue(package['dist']['url'].startswith('https://codeload.github.com/'))
        self.assertEqual(normalize(lock), 0)

    def test_inconsistent_reference_is_rejected(self):
        lock = self.package()
        lock['packages'][0]['dist']['reference'] = 'b' * 40
        with self.assertRaises(ValueError):
            normalize(lock)
