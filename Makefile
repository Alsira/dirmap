SOURCES := $(wildcard cmd/dirmap/*)

.ONESHELL:
release: $(SOURCES)

	if [ ! -d build ]; then
		mkdir build
	fi
	go build -o ./build ./cmd/dirmap

clean:
	rm -r build
