exe = timeat
winexe = $(exe).exe
version = $(shell go run *.go --version)


build:
	go build

test:
	go test -v ./...

archives:
	GOARCH=386 go build
	xz -c $(exe) > $(exe)-$(version)-386.xz
	GOOS=windows go build
	xz -c $(winexe) > $(winexe)-$(version).xz
	# Make last so native exe will stay
	go build
	xz -c $(exe) > $(exe)-$(version)-amd64.xz

clean:
	rm -f *.xz
	rm -f $(exe) $(winexe)

.PHONY: build test archives clean
