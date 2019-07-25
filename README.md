# TOTP

It's a time-based one-time password (TOTP) code generator.

[![Build Status](https://travis-ci.com/arcanericky/totp.svg?branch=master)](https://travis-ci.com/arcanericky/totp)
[![codecov](https://codecov.io/gh/arcanericky/totp/branch/master/graph/badge.svg)](https://codecov.io/gh/arcanericky/totp)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](http://makeapullrequest.com)

## What it Does

It generates TOTP codes used for two-factor authentication at sites such as Google, GitHub, Dropbox, and AWS.

**Warning**
Every copy of your two-factor credentials increases your risk profile. Using this utility is no exception. This utility will store your TOTP secrets unencrypted on your filesystem. The only protection offered is to store these secrets in a file readable by only your user and protected by the operating system only.

## How to Use

**Add TOTP secrets** to the TOTP configuration file with the `config add` option, specifying the name and secret value. Note the secret names are **case sensitive**.

```
$ totp config add mysecretname seed
```

**Generate TOTP codes** using the `totp` command to specify the secret name. Note that because `totp` reserves the use of the words `config` and `version`, don't use them to name a secret.

```
$ totp mysecretname
```

**List the secret entries** with the `config list` command.

```
$ totp config list
```

**Update secret entries** using the `config update` command. Note that `config update` and `config add` are actually the same command and can be used interchangeably.

```
$ totp config update mysecretname newseed
```

**Rename the secret entries** with the `config rename` command

```
$ totp config rename mysecretname mynewname
```

**Delete secret entries** with the `config delete` command

```
$ totp config delete mynewname
```

**Remove all the secrets** and start over using the `config reset` command

```
$ totp config reset
```

**Use an ad-hoc secret** to generate a code by using the `--secret` option

```
$ totp --secret seed
```

**Continuous code output** can be generated with the `--follow` option.

```
$ totp --follow mysecretname
```

**For help** on any of the above, use the `--help` option. Examples are

```
$ totp --help
$ totp config --help
```

**Bash completion** can be enabled by using `config completion`.

```
$ . <(totp config completion)
```

## Using the Time Machine

`totp` has the `--time`, `--forward`, and `--backward` options that are used to manipulate the time for which the TOTP code is generated. This is useful if `totp` is being used on a machine with the incorrect time.

The `--time` option takes an [RFC3339 formatted time string](https://tools.ietf.org/html/rfc3339) as its argument and uses it to generate the TOTP code. Note that the `--forward` and `--backward` options will modify this option value.

Examples with `--time`:

```
$ date '+%FT%T%:z'
2019-06-01T19:58:47-05:00
$ totp --time $(date '+%FT%T%:z') --secret seed
931665
$ totp --time 2019-06-01T20:00:00-05:00 --secret seed
526171
```

The `--forward` and `--backward` options move the current time forward and backward by their duration formatted arguments. See [Go's `time.ParseDuration()`](https://golang.org/pkg/time/#ParseDuration) documentation for more details on this format.

Examples with `--forward` and `--backward`

```
$ totp --time 2019-06-01T20:00:00-05:00 --backward 3m --secret seed
222296
$ totp --time 2019-06-01T20:00:00-05:00 --forward 30s --secret seed
820148
```

The `--follow` option is also compatible with the time machine.

```
$ totp --time 2001-10-31T20:00:00-05:00 --follow --secret seed
```

## Using the Stdio Option

If storing secrets in the clear isn't ideal for you, `totp` supports streaming the shared secret collection through stdin and stdout with the `--stdio` option. This allows you to roll your own encryption or support other methods of maintaining shared secrets.

The `totp <secret name>` and `totp config list` commands support loading the collection via standard input. The 
`totp config update`, `totp config delete`, and `totp config rename` commands support loading via standard input and sending the modified collection to standard output. Experiment with the `--stdio` option to observe how this works.

**Learning with Cleartext Data**

Note the `--file` option can achieve the same results as this example. This is meant to teach how stdio works with `totp`.

Create a collection

```
totp config add --stdio secretname myvalue < /dev/null > totp.json
```

View the collection

```
totp config list --stdio < totp.json
```

Generate a TOTP code

```
totp secretname --stdio < totp.json
```

**Encrypting Shared Secret Collection**

Using what was learned above, a contrived example for encrypting data with [GnuPG](https://gnupg.org/) follows.

Create an encrypted collection
```
totp config add --stdio secretname myvalue < /dev/null | \
  gpg --batch --yes --passphrase mypassphrase --output totp-collection.gpg --symmetric
```

View the collection

```
gpg --quiet --batch --passphrase mypassphrase --decrypt totp-collection.gpg | \
  totp config list --stdio
```

Add another secret

```
gpg --quiet --batch --passphrase mypassphrase --decrypt totp-collection.gpg | \
  totp config add  --stdio newname newvalue | \
  gpg --batch --yes --passphrase mypassphrase --output totp-collection.gpg --symmetric
```

View the modified collection

```
gpg --quiet --batch --passphrase mypassphrase --decrypt totp-collection.gpg | \
  totp config list --stdio
```

Generate a TOTP code

```
gpg --quiet --batch --passphrase mypassphrase --decrypt totp-collection.gpg | totp --stdio secretname
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

The coverage is output to `coverage.html`. Load it in browser for review. For example:

```
$ /opt/google/chrome/chrome file://$PWD/coverage.html
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

My [ga-cmd project](https://github.com/arcanericky/ga-cmd) is more popular than I expected. It's basically the same as `totp` with a much smaller executable, but the list of secrets must be edited manually and there aren't as many command line options. This `totp` project allows the user to maintain the secret collection through the `totp` command line interface, run on a variety of operating systems, and gives me a platform to practice my Go coding.

## Credits

This utility uses the [otp package by pquerna](https://github.com/pquerna/otp). Without this library, I probably woudn't have bothered creating this.