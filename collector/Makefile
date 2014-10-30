all: temp.db sensor/htu21df_demo sensor/libhtu21df.so.1.0.1 temp_and_humidity_pb2.py

temp.db:
	@echo 'Ensuring DB exists'
	sqlite3 temp.db < create_temp_db.sql

sensor/htu21df.o:
	@echo 'Compiling htu21df.c'
	gcc -Wall -c -fPIC sensor/htu21df.c -o sensor/htu21df.o -lwiringPi

sensor/htu21df_demo: sensor/htu21df.o
	@echo 'Compiling htu21df_demo'
	gcc -Wall sensor/htu21df_demo.c sensor/htu21df.o -o sensor/htu21df_demo -lwiringPi

sensor/libhtu21df.so.1.0.1: sensor/htu21df.o
	@echo 'Building libhtu21df.so'
	gcc -shared -Wl,-soname,libhtu21df.so -o sensor/libhtu21df.so.1.0.1 sensor/htu21df.o -lwiringPi

temp_and_humidity_pb2.py: ../common/temp_and_humidity.proto
	@echo 'Compiling temp_and_humidity.proto '
	protoc --python_out=. --proto_path=../common ../common/temp_and_humidity.proto

test: all
	@echo 'Testing C demo'
	@./sensor/htu21df_demo
	@echo 'Testing python wrapper'
	@./sensor/htu21df.py

.PHONY:	clean
clean:
	rm -f sensor/htu21df.o
	rm -f sensor/libhtu21df.so.1.0.1
	rm -f sensor/htu21df_demo
	rm -f sensor/htu21df.pyc
	rm -f temp_and_humidity_pb2.py
	rm -f temp_and_humidity_pb2.pyc
