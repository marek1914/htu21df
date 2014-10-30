all: htu21df/temp_and_humidity.pb.go

htu21df/temp_and_humidity.pb.go: ../../common/temp_and_humidity.proto
	@echo 'Compiling temp_and_humidity.proto '
	protoc --go_out=./htu21df/ --proto_path=../../common ../../common/temp_and_humidity.proto

.PHONY:	clean
clean:
	rm -f htu21df/temp_and_humidity.pb.go