# munin-http-timing
munin-node plugin to get detailed HTTP response timing information from
requesting an URI.

![rendered graph example](example.png)

## Usage
Build using `make`, link the executable in `/etc/munin/plugins/`, configure it
in `/etc/munin/plugin-conf.d/` and restart the `munin-node` service.

## Configuration
URIs must be registered in the environment variables using variables named
`TARGET_<name>`.

Example:
```
[http-timing]
env.TARGET_EXAMPLE https://example.com/
env.TARGET_GITHUB https://github.com/L-P
```

## Tests
```bash
# run test suite
go test

# get code coverage and display it in browser
make cover
```

## License
[MIT](LICENSE)
