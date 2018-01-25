# 线上调试bug指南

新建一个同环境同主机容器，然后修改配置文件，打开调试按钮。

```
docker run -ti registry.alishui.com:443/fanghui/system/registry-manager bash -l

docker run -ti -v /data/docker/registry-manager/config/registry-manager-self.yml:/conf.yml registry.alishui.com:443/fanghui/system/registry-manager bash -l

然后可用修改/conf.yml到本地目录编辑打开debug选项就可以了。

最后运行：/app -c conf.yml

不是/app -c /conf.yml根目录是不待调试的。

```