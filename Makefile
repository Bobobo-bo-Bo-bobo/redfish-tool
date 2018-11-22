GOPATH	= $(CURDIR)
BINDIR	= $(CURDIR)/bin

PROGRAMS = redfish-tool

build:
	env GOPATH=$(GOPATH) go install $(PROGRAMS)

destdirs:
	mkdir -p -m 0755 $(DESTDIR)/usr/bin

strip: build
	strip --strip-all $(BINDIR)/redfish-tool

install: strip destdirs install-bin

install-bin:
	install -m 0755 $(BINDIR)/redfish-tool $(DESTDIR)/usr/bin

clean:
	/bin/rm -f bin/redfish-tool

uninstall:
	/bin/rm -f $(DESTDIR)/usr/bin

all: build strip install

