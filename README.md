# One-Time-Padlock (OTPLock)

[![Yum](https://img.shields.io/badge/-Buy%20me%20a%20cookie-blue?style=for-the-badge&logo=cookiecutter)](https://www.buymeacoffee.com/mjwhitta)

[![Go Report Card](https://goreportcard.com/badge/github.com/mjwhitta/otplock)](https://goreportcard.com/report/github.com/mjwhitta/otplock)

## What is this?

This go package provides a utility for Red Teamers to host dynamic OTP
codes for their payloads.

## How to install

Open a terminal and run the following:

```
$ go install --ldflags "-s -w" --trimpath \
    github.com/mjwhitta/otplock/cmd/otplock@latest
```

## Usage

Simply run `otplock [--unsafe]` in a terminal, and open the URL it
prints out.

### Simple

To start, enter the endpoint, the length of time the OTP key should be
valid, and the payload in hex. After hitting submit, you will be given
the URL for the decryption key and the encrypted payload in hex. Copy
and paste those to your source code and compile.

### Advanced

**Warning:** This usage can be unsafe. This will allow anyone with the
link to run arbitrary commands on your box. It is suggested to only
run this on a fresh VM with minimal network connections (separate
VLAN).

Change to the Advanced config level to get a little more
functionality. The Advanced config will let you upload your source
code and will attempt to compile it for you as you submit payloads.

To start, enter the endpoint, the source filename, the command to
compile, and the name of the compiled binary to return to the user
upon payload submission. After hitting submit, you will be given a URL
to save for submitting payloads.

From here it is quite similar to the Simple config, except a binary is
returned instead of the URL for the decryption key and the encrypted
payload.

## Links

- [Source](https://github.com/mjwhitta/otplock)

## TODO

- Better README
