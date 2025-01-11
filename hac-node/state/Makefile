
all: types.pb.go

%.pb.go: %.proto
	protoc --go_out . $<

clean:
	rm -rf types.pb.go
