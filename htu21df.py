#!/usr/bin/env python

import ctypes
from ctypes import cdll
import os
import sys


HTU21DF_C_SO = './libhtu21df.so.1.0.1'


if not os.path.isfile(HTU21DF_C_SO):
  print 'Error, could not find %s' % HTU21DF_C_SO
  sys.exit(1)


libhtu21df = cdll.LoadLibrary(HTU21DF_C_SO)
libhtu21df.read_temperature.restype = ctypes.c_double
libhtu21df.read_humidity.restype = ctypes.c_double


class Htu21df(object):
  def __init__(self):
    self._fd = libhtu21df.open_connection()

  def __enter__(self):
    return self
  
  def __exit__(self, *args):
    self.close()

  def close(self):
    libhtu21df.close_connection(self._fd)
    
  @property
  def temperature(self):
    return libhtu21df.read_temperature(self._fd)
  
  @property
  def humidity(self):
    return libhtu21df.read_humidity(self._fd)


def main():
  with Htu21df() as sensor:
    print '%.2f C' % sensor.temperature
    print '%.2f%% rh' % sensor.humidity


if __name__ == '__main__':
  main()