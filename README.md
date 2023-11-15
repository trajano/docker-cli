# Wrapped Docker CLI

This wraps the Docker CLI so that it establishes saner defaults for my own needs.

## What this does

- [x] `docker run` will set to logging to none and `-it` by default

- [x] `docker rm` is forced
- [x] `docker ls` maps to `docker ps`
- [x] `docker ps` will actually call `docker inspect` then render the data using github.com/jedib0t/go-pretty.
  - Primarily because I want to know WHEN the bloody container started,
  - how long it took to start rather than about a minute ago.
- [x] `docker service restart` maps to `docker service update --force`
- [x] `docker service ls` shows only relavent columns
  - [x] drop `:latest` if that's the image tag
  - [x] use github.com/jedib0t/go-pretty to render the table
  - [x] `--down` to list services that are not fully up
- [x] `docker ptag` to tag and push image in one command

  - [x] `-p` to add support for generating a patch tag

- [ ] `docker service push <service> <image>` replaces the image of the service if image is not provided it pulls and then does the update to ensure it is the latest copy. It will also add `--with-registry-auth` as appropriate

- [ ] `docker service --help` should call `docker service --help` but append the extra commands
- [ ] `docker service inspect` shows data in YAML
  - [ ] gets rid of the previous spec
  - [ ] use network names rather than the IDs
- [x] `docker du` shows disk usage stats in YAML (generally the output for almost all will be in YAML rather than JSON so it's easier for human consumption)
- [ ] `docker stat` runs `docker stats --no-stream` with sane columns (i.e. no ID)
- [ ] `docker stats` runs with sane columns (i.e. no ID) and hopefully less flashing

## Architecture

- Use Cobra to manage the CLI
- Use the CLI for the most part until they have a package that would handle `DOCKER_HOST`
