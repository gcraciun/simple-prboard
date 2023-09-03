# simple-prboard

A simple web board written in Go, which displays information about Github repos you want, including number of open PRs.


## Inputs

* Reads an env var called GIT_TOKEN. Specify the token used to access the Github API.
* Reads a config file for the rest of the settings. It has to be in the conf/config.yaml path relative to the binary file.

```yaml
collector:
  owner: gcraciun   # your org owner
  interval: 30      # number of seconds between github polls
  refresh: true     # unused
repos:
  category01:       # user supplied name for category (add as many as you need)
    - repo01        # name of a repository you want to poll
    - repo02
  category02:
    - repo03
    - repo04
  category03:
    - repo05
    - repo06
```

## TODO

* Set TCP Listening port as variable, currently staticly binded to 8084.


## Notes

* Weekend project, probably not well maintained.
* Dockerfile has commented lines for building on ubuntu:22.04 instead of scratch (if you need to debug)

