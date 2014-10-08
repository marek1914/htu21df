#!/usr/bin/env python

import sqlite3
import sys

from google.protobuf import text_format

import temp_and_humidity_pb2 as pb


def pb_factory(cursor, row):
  rec_proto = pb.TempAndHumidity()
  rec_proto.recorded_timestamp_ms = row[0]
  rec_proto.temp_degrees_c = row[1]
  rec_proto.percent_relative_humidity = row[2]
  rec_proto.sensor_name = row[3]
  rec_proto.debug = row[4]
  return rec_proto


def extract_all():
  with sqlite3.connect('temp.db') as conn:
    conn.row_factory = pb_factory
    
    protos = conn.execute('''
        SELECT 
          timestamp, 
          temp_degrees_c, 
          relative_humidity, 
          sensor,
          debug
        FROM 
          temp_and_humidity;''').fetchall()
    return protos


def main():
  for p in extract_all():
    print text_format.MessageToString(p,  as_one_line=True)


if __name__ == '__main__':
  main()