# webfinger-enough

just enough webfinger that [Tailscale](https://tailscale.com/kb/1240/sso-custom-oidc) to let me use my own OIDC provider.

usage:

```
go run . -idp=https://idp.nsood.in -domain=nsood.in
```

will start an HTTP server on port 9090 that responds to `acct:<username>@nsood.in`.

## license

MIT license; see LICENSE.md.
