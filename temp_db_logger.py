#!/usr/bin/env python

import htu21df
import sqlite3
import time

def main():
  while 1:
    temp = None
    humidity = None
  
    with htu21df.Htu21df() as sensor:
      temp = sensor.temperature
      humidity = sensor.humidity
  
    if not temp and not humidity:
      print 'Failed to read sensor'
      return
  
    with sqlite3.connect('temp.db') as conn:
      now = int(time.time() * 1e3)
      cur = conn.cursor()
      cur.execute('''
        INSERT INTO temp_and_humidity VALUES (%d, %.3f, %.3f, 'htu21df', '%s');
        ''' % (now, temp, humidity, __file__))
    time.sleep(60 * 5)


if __name__ == '__main__':
  main()