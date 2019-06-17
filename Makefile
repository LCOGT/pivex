STATIC := static

all: version build clean

build:
	packr build

version:
	mkdir ${STATIC}
	git describe > static/version

clean:
	rm -rf ${STATIC}
