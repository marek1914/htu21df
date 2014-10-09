#include <stdio.h>
#include <stdlib.h>
#include <stdint.h>
#include <errno.h>
#include <unistd.h>
#include <wiringPi.h>
#include <wiringPiI2C.h>

#include "htu21df.h"

#define HTU21DF_I2CADDR 0x40

#define HTU21DF_READTEMP 0xE3
#define HTU21DF_READHUM 0xE5
#define HTU21DF_WRITEREG 0xE6
#define HTU21DF_READREG 0xE7

#define HTU21DF_RESET 0xFE
#define HTU21DF_OK 0x02

#define TWO_TO_THE_16TH 65536.0

void die(const char *message) {
  if(errno) {
    perror(message);
  } else {
    printf("ERROR: %s\n", message);
  }
  exit(1);
}

int connect_or_die() {
  int i2c_fd = wiringPiI2CSetup(HTU21DF_I2CADDR);
  if (i2c_fd < 0) {
    die("wiringPiI2CSetup failed");
  }
  return i2c_fd;
}

void reset(int fd) {
  wiringPiI2CWrite(fd, HTU21DF_RESET);
  delay(15);
}

int open_connection() {
  int fd = connect_or_die();
  reset(fd);
  uint8_t status = wiringPiI2CReadReg8(fd, HTU21DF_READREG);
  if (HTU21DF_OK != status) {
    printf("Status not ok byte: %02X", status);
    die("Status not ok");
  }
  return fd;
}

void close_connection(int fd) {
  close(fd);
}

double sensor_convert(double offset, double factor, int sensor_out) {
  double scaled_output = (sensor_out / TWO_TO_THE_16TH);
  return offset + (factor * scaled_output);
}

double degrees_c(int sensor_out) {
  return sensor_convert(-46.85, 175.72, sensor_out);
}

uint16_t read_and_check_crc(fd) {
  unsigned int bytes_read;
  uint8_t buf[4] = {0, 0, 0, 0};
  bytes_read = read(fd, buf, 3);
  if (bytes_read != 3) {
    printf ("%d: %02X %02X %02X\n", bytes_read, buf[0], buf[1], buf[2]);
    die("More than 3 bytes returned\n");
  }
  // TODO: check CRC.
  return (buf [0] << 8 | buf [1]) & 0xFFFC;
}

double read_temperature(int fd) {
  wiringPiI2CWrite(fd, HTU21DF_READTEMP);
  delay(50);
  uint16_t raw_temp = read_and_check_crc(fd);
  return degrees_c(raw_temp);
}

double percent_relative_humidity(int sensor_out) {
  return sensor_convert(-6.0, 125.0, sensor_out);
}

double read_humidity(int fd) {
  wiringPiI2CWrite(fd, HTU21DF_READHUM);
  delay(50);
  uint16_t raw_humidity = read_and_check_crc(fd);
  return percent_relative_humidity(raw_humidity);
}
