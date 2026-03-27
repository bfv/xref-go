# go-cli-template
A template repo for a Go CLI project 

## installation

`github/workflows/release.yaml` contains there placeholders:
```
env:
  DOCKER_IMAGE: devbfvio/xref
  EXE_NAME: xref
  IMAGE_TITLE: xref
  IMAGE_DESCRIPTION: OpenEdge xref parser and query tool
  IMAGE_LICENSES: MIT
  IMAGE_AUTHORS: Bronco Oostermeyer <dev@bfv.io>
```

These can be set via the `init.sh` script. For Windows, use either WSL2 or Git bash.
The `xref`, `xref` and `OpenEdge xref parser and query tool` placeholder can be used in any file. 

## Docker hub integration
`release.yaml` is configured to push a Docker image upon tagging with `vx.y.z` (via `gorelease`)
Set the value in Github:
```
# dont forget to set these repository secrets and variables:
#   - DOCKERHUB_USERNAME
#   - DOCKERHUB_TOKEN
```

