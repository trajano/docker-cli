# Wrapped Docker CLI
This wraps the Docker CLI so that it establishes saner defaults for my own needs.

## What this does

* `docker run` will set to logging to none and `-it` by default
* `docker rm` is forced
* `docker ps` will actually call `docker inspect` then render the data using github.com/jedib0t/go-pretty.
  * Primarily because I want to know WHEN the bloody container started,
  * how long it took to start rather than about a minute ago.
* `docker service restart` maps to `docker service update --force`
* `docker service push <service> <image>` replaces the image of the service if image is not provided it pulls and then does the update to ensure it is the latest copy.  It will also add `--with-registry-auth` as appropriate
* `docker service --help` should call `docker service --help` but append the extra commands

## Architecture

* Command design pattern
* The invoker uses a Composite design pattern to collect a list of commands which are in priority order.  Each command has a "canHandle" method which if true will make the invoker set the command as the command to execute
* There's only one receiver though it can be mocked for testing
