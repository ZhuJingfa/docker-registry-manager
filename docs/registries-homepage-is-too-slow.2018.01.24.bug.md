# INTRO

请求http://n77.node.ifanghui.cn:18080/registries页面特别慢。但是本地有好像正常。

凭空多了两次ping，特别慢。

* Beego语法模板使用的是官方text/template包，{{.Name}}，可能调用的是对象的Key或者Method都有可能。

下面这段代码调用了两次status，导致了两次请求。

改用goroutine在后台更新Status，然后函数返回缓存状态值，速度快了很多。

```
{{if eq $registry.Status "UP" }}
  <span class="label label-success text-capitalize">{{$registry.Status}}</span>
{{else}}
  <span class="label label-danger text-capitalize">{{$registry.Status}}</span>
{{ end }}
```


```
/mysql digest=sha256:d62936e2f1f888d2454684340c88c36397e8efb83a377adc062479b0f4732fb7
DEBU[0009] registry.tags url=https://aliyun.node.ifanghui.cn:7788/v2/fanghui/system/redis/tags/list repository=fanghui/system/redis  file=registry.go line=58 source=app
DEBU[0009] registry.manifest.get url=https://aliyun.node.ifanghui.cn:7788/v2/fanghui/system/redis/manifests/latest repository=fanghui/system/redis reference=latest
DEBU[0009] registry.layer.check url=https://aliyun.node.ifanghui.cn:7788/v2/fanghui/system/redis/blobs/sha256:9b639d5de9dac83ced16ed2d0f1b331bfa529dcbd7de1cb02cd6c5b6bb2a7b81 repository=fanghui/system/redis digest=sha256:9b639d5de9dac83ced16ed2d0f1b331bfa529dcbd7de1cb02cd6c5b6bb2a7b81  file=registry.go line=58 source=app
DEBU[0009] registry.tags url=https://aliyun.node.ifanghui.cn:7788/v2/fanghui/system/registry-manager/tags/list repository=fanghui/system/registry-manager  file=registry.go line=58 source=app
DEBU[0009] registry.manifest.get url=https://aliyun.node.ifanghui.cn:7788/v2/fanghui/system/registry-manager/manifests/latest repository=fanghui/system/registry-manager reference=latest
DEBU[0009] registry.layer.check url=https://aliyun.node.ifanghui.cn:7788/v2/fanghui/system/registry-manager/blobs/sha256:c432530c33f263bbe2c9d9764cf325497e5f3743ab22866b0fc189b5fbfff5d9 repository=fanghui/system/registry-manager digest=sha256:c432530c33f263bbe2c9d9764cf325497e5f3743ab22866b0fc189b5fbfff5d9  file=registry.go line=58 source=app
INFO[0009] http server Running on http://:8080           file=log.go line=610 source=beego
INFO[0009] Admin server Running on :8088                 file=log.go line=610 source=beego
DEBU[0023] registry.ping url=https://aliyun.node.ifanghui.cn:7788/v2/  file=log.go line=610 source=beego
DEBU[0023] registry.ping url=https://aliyun.node.ifanghui.cn:7788/v2/  file=log.go line=610 source=beego
DEBU[0023] registry.ping url=https://registry.alishui.com:443/v2/  file=log.go line=610 source=beego
DEBU[0023] registry.ping url=https://registry.alishui.com:443/v2/  file=log.go line=610 source=beego
DEBU[0023] |      127.0.0.1| 200 | 185.371819ms|   match| GET      /registries   r:/registries  file=log.go line=615 source=beego
DEBU[0024] registry.ping url=https://aliyun.node.ifanghui.cn:7788/v2/
DEBU[0024] registry.ping url=https://aliyun.node.ifanghui.cn:7788/v2/
DEBU[0024] registry.ping url=https://registry.alishui.com:443/v2/
DEBU[0024] registry.ping url=https://registry.alishui.com:443/v2/
DEBU[0029] |      127.0.0.1| 200 | 5.300475969s|   match| GET      /registries   r:/registries  file=log.go line=615 source=beego
```