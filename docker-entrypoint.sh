#!/bin/sh

cd /go/src/

# install the graphql-server application that we copied into the container
go install graphql-server

cd /go/bin/

if [ "$hotswapenv" -ge 1 ]; then
  # install code hotswap library
  go get github.com/pilu/fresh

  # run fresh to listen for any code changes and rebuild the application
  fresh -c runner.conf &
else
  if [ "$debugenv" -ge 1 ]; then
    # install remote go debugging package
    go get github.com/derekparker/delve/cmd/dlv

    # execute the graphql-server project in debug mode using delve debugger when spinning up the container
    dlv --listen=:2345 --headless=true --api-version=2 --log exec ./graphql-server
  else
    ./graphql-server
  fi
fi
