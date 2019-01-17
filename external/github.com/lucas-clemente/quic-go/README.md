# A QUIC implementation in pure Go

<img src="docs/quic.png" width=303 height=124>

[![Godoc Reference](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/v2ray.com/core/external/github.com/lucas-clemente/quic-go)
[![Travis Build Status](https://img.shields.io/travis/lucas-clemente/quic-go/master.svg?style=flat-square&label=Travis+build)](https://travis-ci.org/lucas-clemente/quic-go)
[![CircleCI Build Status](https://img.shields.io/circleci/project/github/lucas-clemente/quic-go.svg?style=flat-square&label=CircleCI+build)](https://circleci.com/gh/lucas-clemente/quic-go)
[![Windows Build Status](https://img.shields.io/appveyor/ci/lucas-clemente/quic-go/master.svg?style=flat-square&label=windows+build)](https://ci.appveyor.com/project/lucas-clemente/quic-go/branch/master)
[![Code Coverage](https://img.shields.io/codecov/c/github/lucas-clemente/quic-go/master.svg?style=flat-square)](https://codecov.io/gh/lucas-clemente/quic-go/)

quic-go is an implementation of the [QUIC](https://en.wikipedia.org/wiki/QUIC) protocol in Go. It roughly implements the [IETF QUIC draft](https://v2ray.com/core/external/github.com/quicwg/base-drafts), although we don't fully support any of the draft versions at the moment.

## Version compatibility

Since quic-go is under active development, there's no guarantee that two builds of different commits are interoperable. The QUIC version used in the *master* branch is just a placeholder, and should not be considered stable.

If you want to use quic-go as a library in other projects, please consider using a [tagged release](https://v2ray.com/core/external/github.com/lucas-clemente/quic-go/releases). These releases expose [experimental QUIC versions](https://v2ray.com/core/external/github.com/quicwg/base-drafts/wiki/QUIC-Versions), which are guaranteed to be stable.

## Google QUIC

quic-go used to support both the QUIC versions supported by Google Chrome and QUIC as deployed on Google's servers, as well as IETF QUIC. Due to the divergence of the two protocols, we decided to not support both versions any more.

The *master* branch **only** supports IETF QUIC. For Google QUIC support, please refer to the [gquic branch](https://v2ray.com/core/external/github.com/lucas-clemente/quic-go/tree/gquic). 

## Guides

We currently support Go 1.9+.

Installing and updating dependencies:

    go get -t -u ./...

Running tests:

    go test ./...

### HTTP mapping

We're currently not implementing the HTTP mapping as described in the [QUIC over HTTP draft](https://quicwg.org/base-drafts/draft-ietf-quic-http.html). The HTTP mapping here is a leftover from Google QUIC.

### QUIC without HTTP/2

Take a look at [this echo example](example/echo/echo.go).

## Usage

### As a server

See the [example server](example/main.go). Starting a QUIC server is very similar to the standard lib http in go:

```go
http.Handle("/", http.FileServer(http.Dir(wwwDir)))
h2quic.ListenAndServeQUIC("localhost:4242", "/path/to/cert/chain.pem", "/path/to/privkey.pem", nil)
```

### As a client

See the [example client](example/client/main.go). Use a `h2quic.RoundTripper` as a `Transport` in a `http.Client`.

```go
http.Client{
  Transport: &h2quic.RoundTripper{},
}
```

## Contributing

We are always happy to welcome new contributors! We have a number of self-contained issues that are suitable for first-time contributors, they are tagged with [help wanted](https://v2ray.com/core/external/github.com/lucas-clemente/quic-go/issues?q=is%3Aissue+is%3Aopen+label%3A%22help+wanted%22). If you have any questions, please feel free to reach out by opening an issue or leaving a comment.
