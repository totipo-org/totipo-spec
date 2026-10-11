"""Exercise spec structure/history guards through the standalone checker CLI."""
from pathlib import Path
import shutil
import subprocess
import sys
import tempfile
import unittest

BASE = Path(__file__).resolve().parent.parent


class StructureTests(unittest.TestCase):
    def setUp(self):
        temp = tempfile.TemporaryDirectory()
        self.addCleanup(temp.cleanup)
        self.root = Path(temp.name)
        for name in ('spec/totipo-vault-format-v1.md',
                     'review/V1_PRE_R16_REVISION_HISTORY.md'):
            target = self.root / name
            target.parent.mkdir(parents=True, exist_ok=True)
            shutil.copyfile(BASE / name, target)

    def check(self):
        return subprocess.run(
            [sys.executable, str(BASE / 'tools/check_spec.py')],
            cwd=self.root, capture_output=True, text=True)

    def test_current_structure_passes(self):
        result = self.check()
        self.assertEqual(result.returncode, 0, result.stderr)

    def test_structure_mutations_fail(self):
        spec = self.root / 'spec/totipo-vault-format-v1.md'
        original = spec.read_text(encoding='utf-8')
        for old, new in (
            ('## 20.', '## 22.'),
            ('**Revision:** r19', '**Revision:** r99'),
            ('Section 12', 'Section 99'),
            ('### v1/r19\n', '### v1/r99\n'),
            ('../review/V1_PRE_R16_REVISION_HISTORY.md', '../review/missing.md'),
            ('0x0001', '0x00ff'),
            ('regular file only', 'any entry type'),
        ):
            with self.subTest(rule=old):
                self.assertIn(old, original)
                spec.write_text(original.replace(old, new, 1), encoding='utf-8')
                self.assertNotEqual(self.check().returncode, 0)
        spec.write_text(original, encoding='utf-8')
        archive = self.root / 'review/V1_PRE_R16_REVISION_HISTORY.md'
        archive.write_text(archive.read_text(encoding='utf-8').replace(
            '### v1/r1\n', '### v1/r99\n'), encoding='utf-8')
        self.assertNotEqual(self.check().returncode, 0)


if __name__ == '__main__':
    unittest.main()
