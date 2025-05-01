# runcmd-go

runcmd, written in Go

## Usage


**runcmd** command parameter1 parameter2 ...


* Example:

~~~
runcmd ls -l -t -r *.c
~~~

## Configuration

You can create a config-file: `$HOME/.runcmd.toml`
with following content:

~~~toml
[default]
RUNCMD_BASE = "HOME"   
RUNCMD_PATH = "runcmd_logging"
~~~

where RUNCMD_BASE has to be an Environment-variable,
and RUNCMD_PATH is a directory name that will be
automatically created if it does no exist.

