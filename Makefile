exe = timeat

build:
	go build

test:
	go test -v

archives:
	GOARCH=386 go build
	xz -c $(exe) > $(exe)-386.xz
	go build
	xz -c $(exe) > $(exe)-amd64.xz

clean:
	rm -f *.xz
	rm -f $(exe)

github:
	hg bookmark -r default master
	hg push git+ssh://git@github.com/tebeka/timeat.git

.PHONY: build test archives clean github
