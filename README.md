# Vaultflow
Proposal workflow for managing secrets across a team.

# TODO
- Add CLI configrations
- Add YAML file configrations
- Add local encryption options

# build
    make clean
    make

# build docker container
    make docker

# Misc
```docker run -it -p8500:8500 progrium/consul -server -bootstrap-expect 1```

```docker run -it --net=host cgswong/vault:latest server -dev```

``` export VAULT_ADDR=http://localhost:8200```

``` export VAULT_TOKEN=855a998d-db99-dd8d-b664-f73e57f3b648```

```go run main.go --pull```
