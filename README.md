
# Docker Registry Manager 

Docker Registry Manager is a golang written, beego driven, web interface for interacting with multiple docker registries (one to many).

## Impove by Zhujingfa，更快速更稳定。

- Lower cpu usage, disable snagles/docker-registry-manager's regex tag analyse, it is a eat cpu lion.
- Improve performance, tagsize(), Status(), IP() some sync remote request all use async Reresh method.
- Improve docker image build,after build , only 21M after expand.
- Remove dockerhub support, it is very slow in China.

# 怎么支持tag删除？

registry.yml add enable delete.
```
# https://docs.docker.com/registry/configuration/
version: 0.1
log:
  fields:
    service: registry
  accesslog:
    disabled: true
  level: error
storage:
  cache:
    blobdescriptor: inmemory
  filesystem:
    rootdirectory: /var/lib/registry
  # 开启删除功能
  delete:
    enabled: true
http:
  addr: :5000
  headers:
    X-Content-Type-Options: [nosniff]
health:
  storagedriver:
    enabled: true
    interval: 10s
    threshold: 3
auth:
  htpasswd:
    realm: basic-realm
    path: /registry/config/htpasswd
```

## Quickstart
 The below steps assume you have a docker registry currently running (with delete mode enabled (https://docs.docker.com/registry/configuration/). To add a registry to manage, add via the interface... or via the config.yml file


### Docker-Compose (Recommended)
 Install compose (https://docs.docker.com/compose/install/)

```bash
 git clone https://github.com/snagles/docker-registry-manager.git && cd docker-registry-manager
 vim config.yml # add your registry
 docker-compose up -d
 firefox localhost:8080
  ```

### Go
 ```bash
    git clone https://github.com/snagles/docker-registry-manager.git && cd docker-registry-manager
    vim config.yml # add your registry
    cd app && go build . && ./app
    firefox localhost:8080
 ```

### Dockerfile
 ```bash
    vim config.yml # add your registry
    docker run --detach --name docker-registry-manager -p 8080:8080 docker-registry-manager
 ```

## Current Features
 1. Support for docker distribution registry v2 (https and http)
 2. Viewable image/tags stages, commands, and sizes.
 3. Bulk deletes of tags
 4. Registry activity logs
 5. Comparison of registry images to public Dockerhub images

## Planned Features
 1. Authentication for users with admin/read only rights using TLS
 2. Global search
 3. List image shared layers
 4. Event timeline
