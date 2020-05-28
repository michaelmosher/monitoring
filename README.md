# Monitoring

A collection of Go executables and custom libraries that they depend on.
Pre-compiled darwin/amd executables will be attached to each release (see below to build for other platforms).

## Commands

- [cdc_status](cmd/cdc_status)

Each command may have specific setup steps required in addition to the general steps below;
see the command's README for more details.

## Building and Installing

### Prerequisites

- [install Go](https://golang.org/doc/install)
- clone this repo

### Installing

Installing a package will create an executable in **$HOME/go/bin**
(or %USERPROFILE%\go\bin on Windows).
To modify this behavior, use the `GOPATH` and `GOBIN` env variables:

```shell
export GOPATH="$HOME/my_new_go_path" GOBIN="my_goodie_bag"
# now executables are installed in "$HOME/my_new_go_path/my_goodie_bag"
```

From the root of this repo, you can install a command using `go install`:

```shell
$ go install github.com/michaelmosher/monitoring/cmd/cdc_status
```

**Note**: Don't forget to add `$GOPATH/$GOBIN` to your `$PATH`:

```shell
# in (for example) ~/.bash_profile
export PATH="$PATH:~/go/bin"
```

### Building

Installing a Go package is useful for running it locally, but if you want to run it on a different computer, you may need to "build" instead.
The Go compiler has first-class cross-compiling support, which is controlled by the `GOOS` and `GOARCH` env variables:

```shell
$ export GOOS=linux GOARCH=arm64
$ go build github.com/michaelmosher/monitoring/cmd/cdc_status
```

When cross-compiling, it can also be useful to name the generated executable appropriately:

```shell
$ go build -o "cdc_status-$GOOS-$GOARCH" github.com/michaelmosher/monitoring/cmd/cdc_status
```

See `go help build` for more information.
