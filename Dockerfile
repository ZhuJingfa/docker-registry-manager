#构建命令：docker build -t fanghui/system/registry-manager .
#使用docker多层次构建新特性
# step 1
FROM fanghui/app/dev-golang as builder
RUN wget https://codeload.github.com/zhujingfa/docker-registry-manager/zip/master && unzip master \
    && mkdir -p src/github.com/zhujingfa/docker-registry-manager \
    && mv docker-registry-manager-master/* src/github.com/zhujingfa/docker-registry-manager/ \
    && CGO_ENABLED=0 go get -a -ldflags '-s' -v github.com/zhujingfa/docker-registry-manager/app

# step 2
FROM fanghui/system/alpine

# make app base dir
ENV REGISTRY_MANAGER_BASE_DIR=/root/go/src/github.com/zhujingfa/docker-registry-manager/app

# ban error: FATA[0000] mkdir /root/go/src/github.com/zhujingfa/docker-registry-manager/app/logs/: no such file or directory
RUN mkdir -p ${REGISTRY_MANAGER_BASE_DIR}/logs

COPY --from=builder /data/go/bin/app /app

# make sure workdir, fetch source.
WORKDIR ${FH_BUILD_DIR}
COPY --from=builder /data/go/src/github.com/zhujingfa/docker-registry-manager/app/static/ ${FH_BUILD_DIR}/static/
COPY --from=builder /data/go/src/github.com/zhujingfa/docker-registry-manager/app/views/ ${FH_BUILD_DIR}/views/

CMD ["/app"]