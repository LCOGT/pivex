STATIC := static

all: version build clean

build:
	packr build

version:
	rm -rf ${STATIC}
	mkdir ${STATIC}
	git describe > static/version

clean:
	rm -rf ${STATIC}
