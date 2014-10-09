#ifndef HTU21DF_H_4ADB11AE
#define HTU21DF_H_4ADB11AE

int open_connection();
double read_temperature(int fd);
double read_humidity(int fd);
void close_connection(int fd);

#endif /* end of include guard: HTU21DF_H_4ADB11AE */
