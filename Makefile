exe = timeat
src = $(exe).go
arch386 = $(exe)-386.xz
arch64 = $(exe)-amd64.xz

exe: $(exe)

$(exe): $(src)
	go build

test: $(exe)
	go test -v

$(arch386): $(src)
	GOARCH=386 go build
	xz -c $(exe) > $@

$(arch64): $(src)
	go build
	xz -c $(exe) > $@

archives: $(arch386) $(arch64)


github:
	hg bookmark -r default master
	hg push git+ssh://git@github.com/tebeka/timeat.git

clean:
	rm *.xz

.PHONY: exe test clean
