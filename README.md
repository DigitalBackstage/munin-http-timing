# munin-http-timing [![Build Status](https://travis-ci.org/DigitalBackstage/munin-http-timing.svg?branch=master)](https://travis-ci.org/DigitalBackstage/munin-http-timing)
munin-node plugin to get detailed HTTP response timing information from
requesting an URI.

![rendered graph example](example.png)

## Usage
Build using `make release`, link the executable from `releases/` in
`/etc/munin/plugins/`, configure it in `/etc/munin/plugin-conf.d/` and restart
the `munin-node` service.  
Two binaries are provided, one for _ARMv6_ (Raspberry-Pi compatible) and one
for AMD64.

## Configuration
URIs must be registered in the environment variables using variables named
`TARGET_<name>`.

Example:
```
[http-timing]
env.TARGET_EXAMPLE https://example.com/
env.TARGET_GITHUB https://github.com/L-P
```

Other options:
    - `env.RANDOM_DELAY` (default 0) when set to `1` requests will be delayed
      by a random amount. This is useful when you test many URIs on the same
      server and don't want to have them arrive at the same time.

## Tests
```bash
# run test suite
go test

# get code coverage and display it in browser
make cover
```

## License
[MIT](LICENSE)
