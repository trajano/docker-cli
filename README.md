# Wrapped Docker CLI

This wraps the Docker CLI so that it establishes saner defaults for my own needs.

## What this does

- [x] `docker run` will set to logging to none and `-it` by default

- [x] `docker rm` is forced
- [x] `docker ls` maps to `docker ps`
- [x] `docker ps` will actually call `docker inspect` then render the data using github.com/jedib0t/go-pretty.
  - Primarily because I want to know WHEN the bloody container started,
  - how long it took to start rather than about a minute ago.
  - [x] Show ports only if I am not running in service mode
- [x] `docker service restart` maps to `docker service update --force`
- [x] `docker service ls` shows only relavent columns
  - [x] drop `:latest` if that's the image tag
  - [x] use github.com/jedib0t/go-pretty to render the table
  - [x] `--down` to list services that are not fully up
- [x] `docker service env` shows the environment variables
- [x] `docker service ps` without the service list will do all services
  - [x] if primary is running and is desired to be running, don't bother showing the others
- [x] `docker ptag` to tag and push image in one command

  - [x] `-p` to add support for generating a patch tag

- [ ] `docker service push <service> <image>` replaces the image of the service if image is not provided it pulls and then does the update to ensure it is the latest copy. It will also add `--with-registry-auth` as appropriate

- [x] `docker use <context>` maps to `docker context use`
- [ ] `docker context use` allows `docker context use <target>` which uses the target rather than the context name.
  - [ ] automatically create the context if it does not exist
- [ ] `docker service --help` should call `docker service --help` but append the extra commands
- [x] `docker service inspect` shows data
  - [x] gets rid of the previous spec
  - [x] use network names rather than the IDs
- [x] `docker du` shows disk usage stats in JSON
- [ ] `docker stat` runs `docker stats --no-stream` with sane columns (i.e. no ID)
- [ ] `docker stats` runs with sane columns (i.e. no ID) and hopefully less flashing
- [x] `docker context create <name> <dockerhost>` maps to `docker context create <name> --docker "host=<dockerhost>"`
- [ ] `docker context create <dockerhost>` maps to `docker context create <host portion of dockerhost> --docker "host=<dockerhost>"`
- [ ] Special handling for ridiculously long image names, specifically "ghcr.io/open-telemetry/opentelemetry-collector-releases/opentelemetry-collector-contrib"

## Architecture

- Use Cobra to manage the CLI
- Use the CLI for the most part until they have a package that would handle `DOCKER_HOST`
- output will be in JSON like the CLI because there's no yaml mappings in the Docker types
