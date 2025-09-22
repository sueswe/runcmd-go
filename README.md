# runcmd-go

![build workflow](https://github.com/sueswe/runcmd-go/actions/workflows/go.yml/badge.svg?event=push)

runcmd, written in Go.

Used for logging in automation and batch processing, etc.


## Usage

`runcmd command parameter1 parameter2 ...`

* Example:

~~~
runcmd ls -l -t -r *.c
~~~


## Optional: configuration

You can (but you don't have to) create a config-file: `$HOME/.runcmd.toml`
with following content:

~~~toml
[default]
RUNCMD_BASE = "HOME"   
RUNCMD_PATH = "logging/runcmd_logging"
~~~

where `RUNCMD_BASE` has to be an Environment-variable,
and `RUNCMD_PATH` is a directory name that will be
automatically created if it does no exist.


## Dry-runs

export `RUNCMD_DRY` with any value you want, and nothing will happen, just a log-entry.

.
