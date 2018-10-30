[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

# gobackgammond

A backgammon player (artificial intelligence) built atop
https://github.com/chandler37/gobackgammon

The user interface is a web server (HTTP).

TODO(chandler37): Productionize the web server; stop using ListenAndServe().

## How do I use it?

`make srv`

Open http://localhost:8000/ in your web browser. Click the link to begin a new
game. Select your move from the list. To play against the PlayerConservative
AI, just choose the topmost choice.

## Copyright

Copyright 2018 David L. Chandler

See the LICENSE file in this directory.

## FAQ

Q: Is it any good?
A: Oh yes.

## Scalable Vector Graphics (SVG) Playground

To play with SVG, run `make build` to download the dependencies and then run
the following:

`GOPATH= go install github.com/ajstarks/svgo/svgplay/...`

`~/go/bin/svgplay`

Visit `http://127.0.0.1:1999/` in your favorite web browser. Copy and paste the
contents of `./bg.go`
