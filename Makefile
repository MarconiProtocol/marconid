OUT_DIR=out/bin
MARCONID_BUILD_TARGET=agent/service/*.go
ARCH=amd64

#BUILD set to debug to include -race flag
BUILD=release
flags.debug=-race
flags.release=
FLAGS=${flags.${BUILD}}

#VERBOSE set to true to include -x flag
VERBOSE=false
verbosity.true=-x
verbosity.false=
VERBOSITY=${verbosity.${VERBOSE}}

.PHONY: all phony_explicit

default: all

all: marconid_linux

marconid_%: phony_explicit
	env GOOS=$* GOARCH=${ARCH} go build ${FLAGS} ${VERBOSITY} -o ${OUT_DIR}/$*/$@_${ARCH} ${MARCONID_BUILD_TARGET}

clean:
	rm -rf ${OUT_DIR}
