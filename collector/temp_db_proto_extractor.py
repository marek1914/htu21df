#!/usr/bin/env python

import argparse
import itertools
import os
import requests
import sqlite3
import sys

import temp_and_humidity_pb2 as pb


def pb_factory(cursor, row):
  rec_proto = pb.TempAndHumidity()
  rec_proto.recorded_timestamp_ms = row[1]
  rec_proto.temp_degrees_c = row[2]
  rec_proto.percent_relative_humidity = row[3]
  rec_proto.sensor_name = row[4]
  rec_proto.debug = row[5]
  return (row[0], rec_proto)


def extract_all(db_file):
  with sqlite3.connect(db_file) as conn:
    conn.row_factory = pb_factory
    protos = conn.execute('''
        SELECT
          rowid,
          timestamp,
          temp_degrees_c,
          relative_humidity,
          sensor,
          debug
        FROM
          temp_and_humidity
        WHERE
          NOT uploaded LIMIT 10;''').fetchall()
    return dict(protos)


def update_uploaded(db_file, ids):
  with sqlite3.connect(db_file) as conn:
    c = conn.execute('''
        UPDATE temp_and_humidity
        SET uploaded = 1
        WHERE rowid IN (%s)''' % ', '.join(map(str, ids)))
    return c.rowcount


def iter_chunk(iterable, chunk_size):
  it = iter(iterable)
  item = list(itertools.islice(it, chunk_size))
  while item:
    yield item
    item = list(itertools.islice(it, chunk_size))


def upload_records(record_items, remote_host, db_file):
  num_sent = len(record_items)
  request_proto = pb.UploadRequest()
  #request_proto.temp_and_humidty_data.extend(i[1] for i in record_items)
  for rowid, rec_proto in record_items:
    request_proto.temp_and_humidty_data.add().CopyFrom(rec_proto)
    proto_and_id = request_proto.data_and_client_id.add()
    proto_and_id.client_id = str(rowid)
    proto_and_id.temp_and_humidty_data.CopyFrom(rec_proto)

  url = 'http://%s/upload' % remote_host
  resp = requests.post(url, data=request_proto.SerializeToString())

  response_proto = pb.UploadResponse()
  response_proto.ParseFromString(resp.content)
  print ('Sent up %d records, response verified %d were received'
      % (num_sent, response_proto.num_saved))
  num_updated = update_uploaded(db_file, (r[0] for r in record_items))
  print '%d records marked uploaded' % response_proto.num_saved


def main():
  parser = argparse.ArgumentParser(
      description='Upload sensor data from a database to the cloud.')
  parser.add_argument('--remote_host', type=str, default='localhost:8080',
      help='Host of the upload service.')
  parser.add_argument('--db_file', type=str,
      default=os.path.join(os.path.dirname(__file__), 'temp.db'),
      help='sqlite database with records to upload.')
  parser.add_argument('--batch_size', type=int, default=-1,
      help='Upload records in batches of this size. Default (-1)'
           ' is to upload all in one batch.')
  args = parser.parse_args()

  records = extract_all(args.db_file)
  batch_size = len(records) if args.batch_size == -1 else args.batch_size
  for batch in iter_chunk(records.items(), batch_size):
    upload_records(batch, args.remote_host, args.db_file)


if __name__ == '__main__':
  main()