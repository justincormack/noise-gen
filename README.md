This is a small package that generates the
[Noise protocol interactive patterns](http://noiseprotocol.org/noise.html#interactive-patterns).

It is based on Revision 34 with eth interactive patterns.

See [mailing list discussion for rationale](https://moderncrypto.org/mail-archive/noise/2018/001706.html).

Currently it adds `ss` tokens where possible, unlike the standard that omits them in most cases. Otherwise
the patterns are identical.
