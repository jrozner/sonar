# Sonar [![Build Status](https://travis-ci.org/jrozner/sonar.svg?branch=master)](https://travis-ci.org/jrozner/sonar)
Sonar is a reconnaissance tool for enumerating sub domains. It was modeled after [Knock](https://github.com/guelfoweb/knock) and [DNSRecon](https://github.com/darkoperator/dnsrecon) though explicitly not written in Python to avoid the limitations of threading and dependencies. Sonar is statically compiled meaning it has no dependencies and even dynamically builds the default wordlist in at compile time to ensure it is portable. It has native support for most modern operating systems and most modern architectures using Go's extremely simple and fast standard cross compilation toolchain.

## Features
* Zone Transfers
* Wordlist based brute force
* Multiple output formats (json, xml, nmap list)
* Wildcard Detection and bypass
* Threading
* Static compilation
* No external dependencies
* Windows/Linux/Mac/FreeBSD on x86, x86_64, arm

## Building
Pre-built binaries will be distributed in the "Releases" tab on GitHub. If you wish you compile yourself you first need to get the Go compiler either from [https://www.golang.org](https://www.golang.org) or through your operating system's package manager. Once setup and installed follow these steps from within the cloned repository to compile:

```sh
cd cmd/sonar
go build
```

This will produce an executable called sonar (sonar.exe on windows) for the platform you are currently on. If you would like to cross compile for another platform follow the instructions [here](http://dave.cheney.net/2015/08/22/cross-compilation-with-go-1-5) for configuring the Go compiler for cross compilation.

Note: By default for your native platform cgo will likely be enabled for the compiler and thus link against the system's domain resovler instead of using the pure Go one. If you would like to stop that use:
```sh
export CGO_ENABLED=0
```
before compiling

## Custom Wordlists
Sonar is designed to be totally self contained and thus compiles in the wordlist to the executable so that it doesn't have to find it on disk. A default wordlist is provided as part of the source but there is also a utility provided, the wordlist_generator, for generating your own from a newline delimited wordlist. You can find the utility in cmd/wordlist_generator and use that to generate a new source file with the custom wordlist before compiling sonar.
