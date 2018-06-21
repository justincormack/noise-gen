This is a small package that generates the
[Noise protocol interactive patterns](http://noiseprotocol.org/noise.html#interactive-patterns).

It is based on Revision 34 with the deferred patterns.

See [mailing list discussion for rationale](https://moderncrypto.org/mail-archive/noise/2018/001706.html).

Currently it adds `ss` tokens in the deferred patterns too, unlike the draft that currently omits them.
Otherwise the patterns are identical.

To run
```
go get github.com/justincormack/noise-gen
```
Then the `noise-gen` binary should be in `~/go/bin`; run it to generate the patterns.
