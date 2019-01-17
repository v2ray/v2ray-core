```
 _____ _     ____        _        _
|_   _| |   / ___|      | |_ _ __(_)___
  | | | |   \___ \ _____| __| '__| / __|
  | | | |___ ___) |_____| |_| |  | \__ \
  |_| |_____|____/       \__|_|  |_|___/

```

crypto/tls, now with 100% more 1.3.

THE API IS NOT STABLE AND DOCUMENTATION IS NOT GUARANTEED.

[![Build Status](https://travis-ci.org/cloudflare/tls-tris.svg?branch=master)](https://travis-ci.org/cloudflare/tls-tris)

## Usage

Since `crypto/tls` is very deeply (and not that elegantly) coupled with the Go stdlib,
tls-tris shouldn't be used as an external package.  It is also impossible to vendor it
as `crypto/tls` because stdlib packages would import the standard one and mismatch.

So, to build with tls-tris, you need to use a custom GOROOT.

A script is provided that will take care of it for you: `./_dev/go.sh`.
Just use that instead of the `go` tool.

The script also transparently fetches the custom Cloudflare Go 1.10 compiler with the required backports.

## Development

### Dependencies

Copy paste line bellow to install all required dependencies:

* ArchLinux:
```
pacman -S go docker gcc git make patch python2 python-docker rsync
```

* Debian:
```
apt-get install build-essential docker go patch python python-pip rsync
pip install setuptools
pip install docker
```

* Ubuntu (18.04) :
```
apt-get update
apt-get install build-essential docker docker.io golang patch python python-pip rsync sudo
pip install setuptools
pip install docker
sudo usermod -a -G docker $USER
```

Similar dependencies can be found on any UNIX based system/distribution.

### Building

There are number of things that need to be setup before running tests. Most important step is to copy ``go env GOROOT`` directory to ``_dev`` and swap TLS implementation and recompile GO. Then for testing we use go implementation from ``_dev/GOROOT``.

```
git clone https://v2ray.com/core/external/github.com/cloudflare/tls-tris.git
cd tls-tris; cp _dev/utils/pre-commit .git/hooks/ 
make -f _dev/Makefile build-all
```

### Testing

We run 3 kinds of test:.

* Unit testing: <br/>``make -f _dev/Makefile test-unit``
* Testing against BoringSSL test suite: <br/>``make -f _dev/Makefile test-bogo``
* Compatibility testing (see below):<br/>``make -f _dev/Makefile test-interop``

To run all the tests in one go use:
```
make -f _dev/Makefile test
```

### Testing interoperability with 3rd party libraries

In order to ensure compatibility we are testing our implementation against BoringSSL, NSS and PicoTLS.

Makefile has a specific target for testing interoperability with external libraries. Following command can be used in order to run such test:

```
make -f _dev/Makefile test-interop
```

The makefile target is just a wrapper and it executes ``_dev/interop_test_runner`` script written in python. The script implements interoperability tests using ``python unittest`` framework. 

Script can be started from command line directly. For example:

```
> ./interop_test_runner -v InteropServer_NSS.test_zero_rtt
test_zero_rtt (__main__.InteropServer_NSS) ... ok

----------------------------------------------------------------------
Ran 1 test in 8.765s

OK
```

### Debugging

When the environment variable `TLSDEBUG` is set to `error`, Tris will print a hexdump of the Client Hello and a stack trace if an handshake error occurs. If the value is `short`, only the error and the first meaningful stack frame are printed.
