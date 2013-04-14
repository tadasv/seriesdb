BUILDDIR=build
FORMAT_DIRS=seriesdbd util
SERIESDB_SOURCES = $(wildcard seriesdbd/*.go util/*.go)
BINARIES=seriesdbd

all: $(BINARIES)

format:
	for d in $(FORMAT_DIRS) ; do\
		 pushd $$d;\
		 go fmt;\
		 popd;\
	 done

$(BUILDDIR)/%:
	mkdir -p $(dir $@)
	cd $* && go build -o $(abspath $@)

.PHONY: $(BINARIES)

$(BINARIES): %: $(BUILDDIR)/%

$(BUILDDIR)/seriesdbd: $(SERIESDB_SOURCES)

clean:
	$(RM) -r $(BUILDDIR)
