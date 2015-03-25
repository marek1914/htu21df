#!/usr/bin/env python

import os
import os.path
import sqlite3
import unittest

import temp_db_logger
import temp_db_proto_extractor

_DB_NAME = '%s.db' % __file__
# _DATA = [
#   (),
#   (),
# ]


class ProtoExtractorTest(unittest.TestCase):
  def setUp(self):
    self._db_name = self.init_test_db()
  
  def init_test_db(self):
    pass

  def test(self):
    pass


def tear_down_db():
  if os.path.isfile(_DB_NAME):
    os.remove(_DB_NAME)


if __name__ == '__main__':
  try:
    unittest.main()
  finally:
    tear_down_db()
