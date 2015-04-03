all: htu21df/temp_and_humidity.pb.go

htu21df/temp_and_humidity.pb.go: ../../common/temp_and_humidity.proto
	@echo 'Compiling temp_and_humidity.proto '
	mkdir -p ${GOPATH}/src/htu21df
	protoc --go_out=${GOPATH}/src/htu21df/ --proto_path=../../common ../../common/temp_and_humidity.proto

.PHONY:	clean
clean:
	rm -f ${GOPATH}/src/htu21df/temp_and_humidity.pb.go