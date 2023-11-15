# Wrapped Docker CLI

This wraps the Docker CLI so that it establishes saner defaults for my own needs.

## What this does

- [x] `docker run` will set to logging to none and `-it` by default

- [x] `docker rm` is forced
- [x] `docker ps` will actually call `docker inspect` then render the data using github.com/jedib0t/go-pretty.
  - Primarily because I want to know WHEN the bloody container started,
  - how long it took to start rather than about a minute ago.
- [x] `docker service restart` maps to `docker service update --force`
- [x] `docker ptag` to tag and push image in one command

  - [x] `-p` to add support for generating a patch tag

- [ ] `docker service push <service> <image>` replaces the image of the service if image is not provided it pulls and then does the update to ensure it is the latest copy. It will also add `--with-registry-auth` as appropriate

- [ ] `docker service --help` should call `docker service --help` but append the extra commands

## Architecture

- Use Cobra to manage the CLI
