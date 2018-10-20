from beat.beat import TestCase
import os
import sys
sys.path.append('../../vendor/github.com/elastic/beats/libbeat/tests/system')


class BaseTest(TestCase):

    @classmethod
    def setUpClass(self):
        self.beat_name = "s3logsbeat"
        self.beat_path = os.path.abspath(os.path.join(os.path.dirname(__file__), "../../"))
        super(BaseTest, self).setUpClass()
