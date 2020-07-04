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

depend:
	env GOPATH=$(GOPATH) go get -u github.com/sirupsen/logrus/
	env GOPATH=$(GOPATH) go get -u git.ypbind.de/repository/go-redfish.git/
	env GOPATH=$(GOPATH) go get -u golang.org/x/crypto/ssh/terminal

distclean:
	/bin/rm -rf src/github.com/
	/bin/rm -rf src/git.ypbind.de/
	/bin/rm -rf src/golang.org/

all: depend build strip install

