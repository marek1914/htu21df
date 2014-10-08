#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>

#include "htu21df.h"

int main (int argc, char const *argv[]) {
  int fd = open_connection();
  double temp = read_temperature(fd);
  double humidity = read_humidity(fd);
  printf("%f %f", temp, humidity);
  printf("%5.1fC ", temp);
  printf("%5.2f%% rh\n", humidity);
  close_connection(fd);
  return 0;
}