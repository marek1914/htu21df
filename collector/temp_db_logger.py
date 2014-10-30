#!/usr/bin/env python

import argparse
import os
import sqlite3
import time

from sensor import htu21df

def main():
  parser = argparse.ArgumentParser(
      description='Upload sensor data from a database to the cloud.')
  parser.add_argument('--db_file', type=str,
      default=os.path.join(os.path.dirname(__file__), 'temp.db'),
      help='sqlite database with records to upload.')
  args = parser.parse_args()

  while 1:
    temp = None
    humidity = None

    with htu21df.Htu21df() as sensor:
      temp = sensor.temperature
      humidity = sensor.humidity

    if not temp and not humidity:
      print 'Failed to read sensor'
      return

    with sqlite3.connect(args.db_file) as conn:
      now = int(time.time() * 1e3)
      cur = conn.cursor()
      cur.execute('''
        INSERT INTO temp_and_humidity VALUES (%d, %.3f, %.3f, 'htu21df', '%s', 0);
        ''' % (now, temp, humidity, __file__))
    time.sleep(60 * 5)


if __name__ == '__main__':
  main()