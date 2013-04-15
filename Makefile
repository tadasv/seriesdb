BUILDDIR=build
FORMAT_DIRS=seriesdbd util series
TEST_DIRS=series
SERIESDB_SOURCES = $(filter-out %_test.go, $(wildcard seriesdbd/*.go util/*.go series/*.go))
BINARIES=seriesdbd

all: $(BINARIES)

format:
	@for d in $(FORMAT_DIRS) ; do\
		 pushd $$d > /dev/null;\
		 go fmt;\
		 popd > /dev/null;\
	done

tests:
	@for d in $(TEST_DIRS) ; do\
		pushd $$d > /dev/null;\
		go test;\
		popd > /dev/null;\
	done

$(BUILDDIR)/%:
	mkdir -p $(dir $@)
	cd $* && go build -o $(abspath $@)

.PHONY: $(BINARIES)

$(BINARIES): %: $(BUILDDIR)/%

$(BUILDDIR)/seriesdbd: $(SERIESDB_SOURCES)

clean:
	$(RM) -r $(BUILDDIR)
