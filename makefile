# Define variables
GOOS = $(shell go env GOOS)
GOARCH = $(shell go env GOARCH)
OUTDIR = build
OUTFILE = main

# Define targets
.PHONY: all clean

all: windows

windows:
ifeq ($(GOOS),windows)
 
  $(GOARCH)-$(GOOS)go build -o $(OUTDIR)/$(OUTFILE).exe main.go
  $(GOARCH)-$(GOOS)go tool rsrc -arch amd64 -ico icon.ico $(OUTDIR)/$(OUTFILE).exe
  $(GOARCH)-$(GOOS)go tool rsrc -arch 386 -ico icon.ico $(OUTDIR)/$(OUTFILE).exe
else
  $(GOARCH)-$(GOOS)go build -o $(OUTDIR)/$(OUTFILE) .
endif

clean:
  rm -rf $(OUTDIR)

