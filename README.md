## Intro

`buildit` is a very simple CI automation tool for MacOS, Linux and Windows, that connects to a git repository, performs a cleam checkout (retrieves content) and executes configured commands (steps). If there are changes from a previous run, it executes the configured steps in the build directory set as the current working directory (CWD). If not, it does nothing (unles `force` flag is set). Each command is executed as a separate process sequentially. If one command (step) fails, other remaining steps are not executed. `buildit` supports the following options:

```
  -force
    	force runs the build steps even if there are no changes in the repository
  -period int
    	When run in watch mode, defines how often to check the repository. Defaults to 1 minute (default 1)
  -watch
    	watch starts the agent in a periodic watch mode. It checks every minute for changes. Period can be adjusted by using the -period flag.
```

## Download

You can download pre-built binaries for Windows, Mac and Linux from [releases page](https://github.com/skyflyer/buildit/releases).

## Configuration

`buildit` is configured with a `buildit.yml` YAML configuration file. An annotated sample:

```yaml
# define the repository
repo: git@github.com:skyflyer/buildit.git
# if you want to build non-master branch, define it:
branch: develop
# authentication info for accessing repository
# if not provided, defaults to using ssh-agent authentication
auth:
  # the key to use for ssh authentication.
  key: /path/to/my/ssh-key
  username: git
  # if you want to use a password, comment the `key`
  # password: somepassword
steps:
# all of the steps below will be run
- echo Hello Dolly!
- ls
- sh ./script.sh
- go build
```

## Usage

The simplest way to use it is to run `buildit` in a directory where the configuration file `buildit.yml` resides. It will clone/update configured repository and run the steps if there are new commits in the repository. You could also run it with `buildit -watch` to periodically check for changes, or add it to cron or scheduled tasks on windows.

## Known issues

* symlinks in git repository are ignored, since they're not supported on Windows properly
* currently supports authentication configuration for ssh-based repos

## TODO

* more testing
* support http(s) repositories
* parse quotes in commands?