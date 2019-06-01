# TOTP

It's a time-based one-time password (TOTP) code generator.

[![Build Status](https://travis-ci.com/arcanericky/totp.svg?branch=master)](https://travis-ci.com/arcanericky/totp)
[![codecov](https://codecov.io/gh/arcanericky/totp/branch/master/graph/badge.svg)](https://codecov.io/gh/arcanericky/totp)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](http://makeapullrequest.com)

## What it Does

It generates TOTP codes used for two-factor authentication at sites such as Google, GitHub, Dropbox, and AWS.

## How to Use

**Add TOTP seeds** to the TOTP configuration file with the `config add` option, specifying the name and seed value. Note the seed names are **case sensitive**.

```
$ totp config add google seed
```

**Generate TOTP codes** using the `totp` command to specify the seed name. Note that because `totp` reserves the use of the words `config` and `version`, don't use them to name a seed.

```
$ totp google
```

**List the seed entries** with the `config list` command.

```
$ totp config list
```

**Update seed entries** using the `config update` command. Note that `config update` and `config add` are actually the same command and can be used interchangeably.

```
$ totp config update google newseed
```

**Rename the seed entries** with the `config rename` command

```
$ totp config rename google google-main
```

**Delete seed entries** with the `config delete` command

```
$ totp config delete google-main
```

**Remove all the keys** and start over, use the `config reset` command

```
$ totp config reset
```

**Use an ad-hoc seed** to generate a code by using the `--seed` option

```
$ totp --seed seed
```

**For help** on any of the above, use the `--help` option. Examples are

```
$ totp --help
$ totp config --help
```

**Bash completion** can be enabled by using `config completion`

```
$ . <(totp config completion)
```

## Building

`totp` is mostly developed using Go 1.12.x on Debian based systems. Only `go` is required but to use the automated actions the `Makefile` provides, `make` must be installed.

To build everything:

```
$ git clone https://github.com/arcanericky/totp.git
$ cd totp
$ make
```

For unit tests and code coverage reports:

```
$ make test
```

To build for a single platform (see the `Makefile` for the different targets)

```
$ make linux-amd64
```

See the `Makefile` for how to use the `go` command natively.

## Contributing

Contributions and issues are welcome. These include bugs reports and fixes, code comments, spelling corrections, and new features. If adding a new feature, please file an issue so it can be discussed prior to implementation so your time isn't wasted.

Unit tests for new code are required. Use `make test` to verify coverage. Coverage will also be checked with Codecov when pull requests are made.

## Inspiration

My [ga-cmd project](https://github.com/arcanericky/ga-cmd) is more popular than I expected. It's basically the same as `totp` with a much smaller executable, but the list of seeds must be edited manually. This `totp` project allows the user to maintain the seed collection through the `totp` command line interface, run on a variety of operating systems, and gives me a platform to practice my Go coding.
