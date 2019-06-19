STATIC := static

all: version build clean

build:
	packr build

version:
	[[ -d ${STATIC} ]] || mkdir ${STATIC}
	git describe > static/version

clean:
	rm -rf ${STATIC}
