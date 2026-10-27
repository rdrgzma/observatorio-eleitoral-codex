"""Use GitHub's official archive endpoint without changing pinned revisions.

Run after composer update only when api.github.com is unavailable. Archives with
an existing checksum remain untouched: alternate ZIP bytes may differ.
"""
import json
from pathlib import Path
import re


def normalize(lock):
    changed = 0
    for package in lock.get('packages', []) + lock.get('packages-dev', []):
        dist = package.get('dist', {})
        match = re.fullmatch(
            r'https://api\.github\.com/repos/([^/]+)/([^/]+)/zipball/([0-9a-f]{40})',
            dist.get('url', ''),
        )
        if match and not dist.get('shasum'):
            owner, repository, commit = match.groups()
            if dist.get('reference') != commit:
                raise ValueError('Archive reference differs from pinned commit')
            dist['url'] = f'https://codeload.github.com/{owner}/{repository}/zip/{commit}'
            changed += 1
    return changed


if __name__ == '__main__':
    path = Path(__file__).resolve().parents[1] / 'painel/composer.lock'
    lock = json.loads(path.read_text())
    count = normalize(lock)
    temporary = path.with_suffix('.lock.tmp')
    temporary.write_text(json.dumps(lock, indent=4, ensure_ascii=False) + '\n')
    temporary.replace(path)
    print(f'Updated {count} archive URLs; versions and checksums preserved.')
