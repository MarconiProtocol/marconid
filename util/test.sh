#!/bin/bash

usage() {
  cat <<EOF
Usage: ./test.sh [options]

Options:
	-h             Print Help (this message) and exit
	-v             Verbose output: log all tests as they are run
EOF
  exit 1
}

if [[ "$#" -gt 1 ]]; then
  usage
fi

VERBOSE=false

if [[ "$#" == 1 ]]; then
  if [[ "-h" == $1 ]]; then
    usage
  elif [[ "-v" == $1 ]]; then
    VERBOSE=true
  else
    usage
  fi
fi

# check if GOROOT and GOPATH are properly set
if [[ -z "$GOROOT" ]]; then
  echo "GOROOT is not set"
  exit 1
fi

if [[ -z "$GOPATH" ]]; then
  echo "GOPATH is not set"
  exit 1
fi

if [[ ! -d $GOROOT/src ]]; then
  echo "Invalid GOROOT"
  exit 1
fi

if [[ ! -d $GOPATH/src ]]; then
  mkdir -p $GOPATH/src
fi

echo "GOROOT: ${GOROOT}"
echo "GOPATH: ${GOPATH}"
echo ""
read -p "If GOROOT and GOPATH look good to you, press <Enter> to continue: "

# get all the required packages
sudo env GOPATH=$GOPATH GOROOT=$GOROOT PATH=$PATH:$GOROOT/bin go get -d ../...

# create test key and directories
mkdir -p ../core/crypto/build
touch ../core/crypto/build/test_key
sudo mkdir -p /opt/marconi/etc/marconid/keys

# run unit tests
if [[ "$VERBOSE" = true ]]; then
  sudo env GOPATH=$GOPATH GOROOT=$GOROOT PATH=$PATH:$GOROOT/bin go test -v -cover ../...
else
  sudo env GOPATH=$GOPATH GOROOT=$GOROOT PATH=$PATH:$GOROOT/bin go test -cover ../...
fi

# clean up
rm -rf ../core/crypto/build
sudo rm -rf /opt/marconi/etc/marconid/keys
