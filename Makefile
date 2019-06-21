STATIC := static

all: version build clean

build:
	packr2
	go build

version:
	[[ -d ${STATIC} ]] || mkdir ${STATIC}
	git describe > static/version

clean:
	rm -rf ${STATIC}
	packr2 clean
