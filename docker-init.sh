#!/bin/bash

docker run --rm=true -v $(pwd):/web-compiler -w /web-compiler golang /web-compiler/init

