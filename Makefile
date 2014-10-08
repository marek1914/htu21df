all: temp.db htu21df_demo libhtu21df.so.1.0.1 temp_and_humidity_pb2.py

temp.db:
	@echo 'Ensuring DB exists'
	sqlite3 temp.db < create_temp_db.sql

htu21df.o:
	@echo 'Compiling htu21df.c'
	gcc -Wall -c -fPIC htu21df.c -o htu21df.o -lwiringPi

htu21df_demo: htu21df.o
	@echo 'Compiling htu21df_demo'
	gcc -Wall htu21df_demo.c htu21df.o -o htu21df_demo -lwiringPi

libhtu21df.so.1.0.1: htu21df.o
	@echo 'Building libhtu21df.so'
	gcc -shared -Wl,-soname,libhtu21df.so -o libhtu21df.so.1.0.1 htu21df.o -lwiringPi

temp_and_humidity_pb2.py: temp_and_humidity.proto
	@echo 'Compiling temp_and_humidity.proto '
	protoc --python_out=. temp_and_humidity.proto 

test: all
	@echo 'Testing C demo'
	@./htu21df_demo
	@echo 'Testing python wrapper'
	@./htu21df.py

.PHONY:	clean
clean:
	rm -f htu21df.o
	rm -f libhtu21df.so.1.0.1
	rm -f htu21df_demo
	rm -f htu21df.pyc
	rm -f temp_and_humidity_pb2.py
	rm -f temp_and_humidity_pb2.pyc
	