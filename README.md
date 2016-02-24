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
* Windows/Linux/Mac on x86, x86_64, arm
