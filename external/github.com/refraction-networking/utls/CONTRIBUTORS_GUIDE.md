# How this package works
### Chapter 1: [Making private things public](./u_public.go)
There are numerous handshake-related structs in crypto/tls, most of which are either private or have private fields.
One of them — `clientHandshakeState` — has private function `handshake()`,
which is called in the beginning of default handshake.  
Unfortunately, user will not be able to directly access this struct outside of tls package.
As a result, we decided to employ following workaround: declare public copies of private structs.
Now user is free to manipulate fields of public `ClientHandshakeState`.
Then, right before handshake, we can shallow-copy public state into private `clientHandshakeState`,
call `handshake()` on it and carry on with default Golang handshake process.
After handshake is done we shallow-copy private state back to public, allowing user to read results of handshake.

### Chapter 2: [TLSExtension](./u_tls_extensions.go)
The way we achieve reasonable flexibilty with extensions is inspired by
[ztls'](https://github.com/zmap/zcrypto/blob/master/tls/handshake_extensions.go) design.
However, our design has several differences, so we wrote it from scratch.
This design allows us to have an array of `TLSExtension` objects and then marshal them in order:
```Golang
type TLSExtension interface {
	writeToUConn(*UConn) error

	Len() int // includes header

	// Read reads up to len(p) bytes into p.
	// It returns the number of bytes read (0 <= n <= len(p)) and any error encountered.
	Read(p []byte) (n int, err error) // implements io.Reader
}
```
`writeToUConn()` applies appropriate per-extension changes to `UConn`.

`Len()` provides the size of marshaled extension, so we can allocate appropriate buffer beforehand,
catch out-of-bound errors easily and guide size-dependent extensions such as padding.

`Read(buffer []byte)` _writes(see: io.Reader interface)_ marshaled extensions into provided buffer.
This avoids extra allocations.

### Chapter 3: [UConn](./u_conn.go)
`UConn` extends standard `tls.Conn`. Most notably, it stores slice with `TLSExtension`s and public
`ClientHandshakeState`.  
Whenever `UConn.BuildHandshakeState()` gets called (happens automatically in `UConn.Handshake()`
or could be called manually), config will be applied according to chosen `ClientHelloID`.
From contributor's view there are 2 main behaviors:  
 * `HelloGolang` simply calls default Golang's [`makeClientHello()`](./handshake_client.go)
 and directly stores it into `HandshakeState.Hello`. utls-specific stuff is ignored.  
 * Other ClientHelloIDs fill `UConn.Hello.{Random, CipherSuites, CompressionMethods}` and `UConn.Extensions` with
per-parrot setup, which then gets applied to appropriate standard tls structs,
and then marshaled by utls into `HandshakeState.Hello`.

### Chapter 4: Tests

Tests exist, but coverage is very limited. What's covered is a conjunction of
 * TLS 1.2
 * Working parrots without any unsupported extensions (only Android 5.1 at this time)
 * Ciphersuites offered by parrot.
 * Ciphersuites supported by Golang
 * Simple conversation with reference implementation of OpenSSL.
(e.g. no automatic checks for renegotiations, parroting quality and such)

plus we test some other minor things.
Basically, current tests aim to provide a sanity check.

# Merging upstream
```Bash
git remote add -f golang git@github.com:golang/go.git
git checkout -b golang-upstream golang/master
git subtree split -P src/crypto/tls/ -b golang-tls-upstream
git checkout master
git merge --no-commit golang-tls-upstream
```
