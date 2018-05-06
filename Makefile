env ?= devel


PACKR := $(GOPATH)/bin/packr

$(PACKR):
	go get -u github.com/gobuffalo/packr/...


build: $(PACKR)
ifeq ($(env),prod)
	@echo "Making Production build"

else
	@echo "Making Development build"
	packr
	GOOS=linux GOARCH=amd64 go build -tags devel -o output/flair-devel-linux-amd64
	packr clean
endif


run: build
ifeq ($(env),prof)
	echo "production runing"
else
	echo "dev running"
endif
