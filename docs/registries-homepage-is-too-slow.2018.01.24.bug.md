# INTRO

请求http://n77.node.ifanghui.cn:18080/registries页面特别慢。但是本地有好像正常。

凭空多了两次ping，特别慢。

* Beego语法模板使用的是官方text/template包，{{.Name}}，可能调用的是对象的Key或者Method都有可能。

下面这段代码调用了两次status，导致了两次请求。

改用goroutine在后台更新Status，然后函数返回缓存状态值，速度快了很多。

01.26 /registries页面请求/registries/all/repositories/count，请求所有库目录的，还是发现多了两个请求
删除了即时统计的需求，改为了使用TTL更新的缓存。解析IP导致的。Go程序异常block，通过引起block的时候，再kill -3 去调试，获取到当时执行的协程信息。

```
goroutine 1580 [chan receive]:
net.(*Resolver).goLookupIPCNAMEOrder(0xd771b0, 0xd38000, 0xc4200140c0, 0xc4203b6038, 0x17, 0x2, 0xd31340, 0xc4203ba990, 0x1, 0xc4201fd838, ...)
	/server/golang/src/net/dnsclient_unix.go:489 +0x374
net.(*Resolver).goLookupHostOrder(0xd771b0, 0xd38000, 0xc4200140c0, 0xc4203b6038, 0x17, 0x2, 0x0, 0x0, 0x59c40, 0xd992e0, ...)
	/server/golang/src/net/dnsclient_unix.go:425 +0x126
net.(*Resolver).lookupHost(0xd771b0, 0xd38000, 0xc4200140c0, 0xc4203b6038, 0x17, 0xa11400, 0x7fac00059c40, 0x953c40, 0xc4201fd918, 0x43be33)
	/server/golang/src/net/lookup_unix.go:86 +0x96
net.(*Resolver).LookupHost(0xd771b0, 0xd38000, 0xc4200140c0, 0xc4203b6038, 0x17, 0x7faccf39c238, 0x0, 0x8, 0xd79600, 0x7faccf3ed000)
	/server/golang/src/net/lookup.go:152 +0x241
net.LookupHost(0xc4203b6038, 0x17, 0xb80000c4201fdc18, 0xc4201fd9d8, 0x4110a8, 0x20, 0xc4206ab980)
	/server/golang/src/net/lookup.go:138 +0x5d
github.com/zhujingfa/docker-registry-manager/app/models.(*Registry).IP(0xc4203b1780, 0x0, 0x0)
	/data/go/src/github.com/zhujingfa/docker-registry-manager/app/models/registry.go:49 +0x38
```



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

01.26 /registries页面请求/registries/all/repositories/count，请求所有库目录的，还是发现多了两个请求

```
DEBU[0064] |      127.0.0.1| 200 |  16.170527ms|   match| GET      /registries   r:/registries  file=log.go line=615 source=beego
DEBU[0064] |      127.0.0.1| 200 |     64.484µs|   match| GET      /registries/all/count   r:/registries/all/count  file=log.go line=615 source=beego
DEBU[0064] registry.repositories url=https://aliyun.node.ifanghui.cn:7788/v2/_catalog  file=log.go line=615 source=beego
DEBU[0064] registry.repositories url=https://registry.alishui.com/v2/_catalog  file=log.go line=615 source=beego
DEBU[0064] |      127.0.0.1| 200 |  67.123618ms|   match| GET      /registries/all/repositories/count   r:/registries/all/repositories/count  file=log.go line=615 source=beego
```

```
DEBU[0494] |      127.0.0.1| 200 |    26.8587ms|   match| GET      /registries   r:/registries  file=log.go line=615 source=beego
DEBU[0495] |      127.0.0.1| 200 |  20.344322ms|   match| GET      /registries   r:/registries  file=log.go line=615 source=beego
DEBU[0495] |      127.0.0.1| 200 |  45.792349ms|   match| GET      /registries   r:/registries  file=log.go line=615 source=beego
SIGQUIT: quit
PC=0x45c3d1 m=0 sigcode=0

goroutine 0 [idle]:
runtime.futex(0xd79710, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x7ffc456c1a48, 0x40f85b, ...)
	/server/golang/src/runtime/sys_linux_amd64.s:438 +0x21
runtime.futexsleep(0xd79710, 0x0, 0xffffffffffffffff)
	/server/golang/src/runtime/os_linux.go:45 +0x62
runtime.notesleep(0xd79710)
	/server/golang/src/runtime/lock_futex.go:151 +0x9b
runtime.stopm()
	/server/golang/src/runtime/proc.go:1680 +0xe5
runtime.findrunnable(0xc42001b300, 0x0)
	/server/golang/src/runtime/proc.go:2135 +0x4d2
runtime.schedule()
	/server/golang/src/runtime/proc.go:2255 +0x12c
runtime.goexit0(0xc4205c2a80)
	/server/golang/src/runtime/proc.go:2406 +0x236
runtime.mcall(0x7ffc456c1c00)
	/server/golang/src/runtime/asm_amd64.s:286 +0x5b

goroutine 1 [chan receive, 8 minutes]:
github.com/zhujingfa/docker-registry-manager/app/vendor/github.com/astaxie/beego.(*App).Run(0xc420043ee0)
	/data/go/src/github.com/zhujingfa/docker-registry-manager/app/vendor/github.com/astaxie/beego/app.go:186 +0x4af
github.com/zhujingfa/docker-registry-manager/app/vendor/github.com/astaxie/beego.Run(0x0, 0x0, 0x0)
	/data/go/src/github.com/zhujingfa/docker-registry-manager/app/vendor/github.com/astaxie/beego/beego.go:67 +0x51
main.main.func1(0xc4200909a0)
	/data/go/src/github.com/zhujingfa/docker-registry-manager/app/main.go:125 +0xc86
github.com/zhujingfa/docker-registry-manager/app/vendor/github.com/urfave/cli.HandleAction(0x9543c0, 0xc420382560, 0xc4200909a0, 0xc42037c4e0, 0x0)
	/data/go/src/github.com/zhujingfa/docker-registry-manager/app/vendor/github.com/urfave/cli/app.go:492 +0x7c
github.com/zhujingfa/docker-registry-manager/app/vendor/github.com/urfave/cli.(*App).Run(0xc42005bba0, 0xc42000e090, 0x3, 0x3, 0x0, 0x0)
	/data/go/src/github.com/zhujingfa/docker-registry-manager/app/vendor/github.com/urfave/cli/app.go:264 +0x635
main.main()
	/data/go/src/github.com/zhujingfa/docker-registry-manager/app/main.go:127 +0x31c

goroutine 5 [syscall, 8 minutes]:
os/signal.signal_recv(0x0)
	/server/golang/src/runtime/sigqueue.go:131 +0xa6
os/signal.loop()
	/server/golang/src/os/signal/signal_unix.go:22 +0x22
created by os/signal.init.0
	/server/golang/src/os/signal/signal_unix.go:28 +0x41

goroutine 39 [chan receive, 3 minutes]:
github.com/zhujingfa/docker-registry-manager/app/models.AddRegistry.func1(0xc4203c0100)
	/data/go/src/github.com/zhujingfa/docker-registry-manager/app/models/registry.go:271 +0x90
created by github.com/zhujingfa/docker-registry-manager/app/models.AddRegistry
	/data/go/src/github.com/zhujingfa/docker-registry-manager/app/models/registry.go:270 +0x42a

goroutine 1378 [select, 3 minutes]:
net/http.(*persistConn).writeLoop(0xc42065b0e0)
	/server/golang/src/net/http/transport.go:1759 +0x165
created by net/http.(*Transport).dialConn
	/server/golang/src/net/http/transport.go:1187 +0xa53

goroutine 1377 [IO wait, 3 minutes]:
internal/poll.runtime_pollWait(0x7faccf391c70, 0x72, 0x0)
	/server/golang/src/runtime/netpoll.go:173 +0x57
internal/poll.(*pollDesc).wait(0xc4207e8d18, 0x72, 0xffffffffffffff00, 0xd32d00, 0xd2e420)
	/server/golang/src/internal/poll/fd_poll_runtime.go:85 +0xae
internal/poll.(*pollDesc).waitRead(0xc4207e8d18, 0xc42041a000, 0x2000, 0x2000)
	/server/golang/src/internal/poll/fd_poll_runtime.go:90 +0x3d
internal/poll.(*FD).Read(0xc4207e8d00, 0xc42041a000, 0x2000, 0x2000, 0x0, 0x0, 0x0)
	/server/golang/src/internal/poll/fd_unix.go:126 +0x18a
net.(*netFD).Read(0xc4207e8d00, 0xc42041a000, 0x2000, 0x2000, 0x1, 0x4, 0xc4204a4788)
	/server/golang/src/net/fd_unix.go:202 +0x52
net.(*conn).Read(0xc4205d6418, 0xc42041a000, 0x2000, 0x2000, 0x0, 0x0, 0x0)
	/server/golang/src/net/net.go:176 +0x6d
crypto/tls.(*block).readFromUntil(0xc4209a6a50, 0x7faccf34d0e0, 0xc4205d6418, 0x5, 0xc4205d6418, 0xc420b49500)
	/server/golang/src/crypto/tls/conn.go:488 +0x95
crypto/tls.(*Conn).readRecord(0xc420463500, 0xa54b17, 0xc420463620, 0x42d88b)
	/server/golang/src/crypto/tls/conn.go:590 +0xe0
crypto/tls.(*Conn).Read(0xc420463500, 0xc4207c0000, 0x1000, 0x1000, 0x0, 0x0, 0x0)
	/server/golang/src/crypto/tls/conn.go:1134 +0x110
net/http.(*persistConn).Read(0xc42065b0e0, 0xc4207c0000, 0x1000, 0x1000, 0xc42047d780, 0xc4203d1678, 0x455090)
	/server/golang/src/net/http/transport.go:1391 +0x140
bufio.(*Reader).fill(0xc420b94180)
	/server/golang/src/bufio/bufio.go:97 +0x11a
bufio.(*Reader).Peek(0xc420b94180, 0x1, 0x0, 0x0, 0x0, 0xc42026f560, 0x0)
	/server/golang/src/bufio/bufio.go:129 +0x3a
net/http.(*persistConn).readLoop(0xc42065b0e0)
	/server/golang/src/net/http/transport.go:1539 +0x185
created by net/http.(*Transport).dialConn
	/server/golang/src/net/http/transport.go:1186 +0xa2e

goroutine 65 [chan receive, 3 minutes]:
github.com/zhujingfa/docker-registry-manager/app/models.AddRegistry.func1(0xc4203c1800)
	/data/go/src/github.com/zhujingfa/docker-registry-manager/app/models/registry.go:271 +0x90
created by github.com/zhujingfa/docker-registry-manager/app/models.AddRegistry
	/data/go/src/github.com/zhujingfa/docker-registry-manager/app/models/registry.go:270 +0x42a

goroutine 1366 [IO wait, 3 minutes]:
internal/poll.runtime_pollWait(0x7faccf391eb0, 0x72, 0x0)
	/server/golang/src/runtime/netpoll.go:173 +0x57
internal/poll.(*pollDesc).wait(0xc4207e9d18, 0x72, 0xffffffffffffff00, 0xd32d00, 0xd2e420)
	/server/golang/src/internal/poll/fd_poll_runtime.go:85 +0xae
internal/poll.(*pollDesc).waitRead(0xc4207e9d18, 0xc42099c000, 0x4000, 0x4000)
	/server/golang/src/internal/poll/fd_poll_runtime.go:90 +0x3d
internal/poll.(*FD).Read(0xc4207e9d00, 0xc42099c000, 0x4000, 0x4000, 0x0, 0x0, 0x0)
	/server/golang/src/internal/poll/fd_unix.go:126 +0x18a
net.(*netFD).Read(0xc4207e9d00, 0xc42099c000, 0x4000, 0x4000, 0x1, 0x2, 0xc42037b788)
	/server/golang/src/net/fd_unix.go:202 +0x52
net.(*conn).Read(0xc4205d6700, 0xc42099c000, 0x4000, 0x4000, 0x0, 0x0, 0x0)
	/server/golang/src/net/net.go:176 +0x6d
crypto/tls.(*block).readFromUntil(0xc420b49dd0, 0x7faccf34d0e0, 0xc4205d6700, 0x5, 0xc4205d6700, 0xc420e24e10)
	/server/golang/src/crypto/tls/conn.go:488 +0x95
crypto/tls.(*Conn).readRecord(0xc420b4e380, 0xa54b17, 0xc420b4e4a0, 0x42d88b)
	/server/golang/src/crypto/tls/conn.go:590 +0xe0
crypto/tls.(*Conn).Read(0xc420b4e380, 0xc4208a5000, 0x1000, 0x1000, 0x0, 0x0, 0x0)
	/server/golang/src/crypto/tls/conn.go:1134 +0x110
net/http.(*persistConn).Read(0xc4207ad8c0, 0xc4208a5000, 0x1000, 0x1000, 0xc420a8e0e0, 0xc42026fe58, 0x455090)
	/server/golang/src/net/http/transport.go:1391 +0x140
bufio.(*Reader).fill(0xc420780900)
	/server/golang/src/bufio/bufio.go:97 +0x11a
bufio.(*Reader).Peek(0xc420780900, 0x1, 0x0, 0x0, 0x0, 0xc420d9b140, 0x0)
	/server/golang/src/bufio/bufio.go:129 +0x3a
net/http.(*persistConn).readLoop(0xc4207ad8c0)
	/server/golang/src/net/http/transport.go:1539 +0x185
created by net/http.(*Transport).dialConn
	/server/golang/src/net/http/transport.go:1186 +0xa2e

goroutine 103 [IO wait, 8 minutes]:
internal/poll.runtime_pollWait(0x7faccf391d30, 0x72, 0xffffffffffffffff)
	/server/golang/src/runtime/netpoll.go:173 +0x57
internal/poll.(*pollDesc).wait(0xc4207e8118, 0x72, 0xc420a13b00, 0x0, 0x0)
	/server/golang/src/internal/poll/fd_poll_runtime.go:85 +0xae
internal/poll.(*pollDesc).waitRead(0xc4207e8118, 0xffffffffffffff00, 0x0, 0x0)
	/server/golang/src/internal/poll/fd_poll_runtime.go:90 +0x3d
internal/poll.(*FD).Accept(0xc4207e8100, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0)
	/server/golang/src/internal/poll/fd_unix.go:335 +0x1e2
net.(*netFD).accept(0xc4207e8100, 0x7faccf3ed6c8, 0x0, 0xa54be8)
	/server/golang/src/net/fd_unix.go:238 +0x42
net.(*TCPListener).accept(0xc420192010, 0xc420a13d08, 0x4110a8, 0x30)
	/server/golang/src/net/tcpsock_posix.go:136 +0x2e
net.(*TCPListener).AcceptTCP(0xc420192010, 0xc42089c270, 0xc42089c270, 0x979620)
	/server/golang/src/net/tcpsock.go:234 +0x49
net/http.tcpKeepAliveListener.Accept(0xc420192010, 0xc4200140c0, 0x979620, 0xd66c70, 0xa11620)
	/server/golang/src/net/http/server.go:3120 +0x2f
net/http.(*Server).Serve(0xc4202f6000, 0xd37a40, 0xc420192010, 0x0, 0x0)
	/server/golang/src/net/http/server.go:2695 +0x1b2
net/http.(*Server).ListenAndServe(0xc4202f6000, 0xc4202f6000, 0x0)
	/server/golang/src/net/http/server.go:2636 +0xa9
net/http.ListenAndServe(0xc4201f40b0, 0x5, 0x0, 0x0, 0x1, 0xc4201f40b0)
	/server/golang/src/net/http/server.go:2882 +0x7f
github.com/zhujingfa/docker-registry-manager/app/vendor/github.com/astaxie/beego.(*adminApp).Run(0xc42000c600)
	/data/go/src/github.com/zhujingfa/docker-registry-manager/app/vendor/github.com/astaxie/beego/admin.go:399 +0x3af
created by github.com/zhujingfa/docker-registry-manager/app/vendor/github.com/astaxie/beego.registerAdmin
	/data/go/src/github.com/zhujingfa/docker-registry-manager/app/vendor/github.com/astaxie/beego/hooks.go:89 +0x6f

goroutine 104 [IO wait]:
internal/poll.runtime_pollWait(0x7faccf391df0, 0x72, 0xffffffffffffffff)
	/server/golang/src/runtime/netpoll.go:173 +0x57
internal/poll.(*pollDesc).wait(0xc4203b0398, 0x72, 0xc4204a7c00, 0x0, 0x0)
	/server/golang/src/internal/poll/fd_poll_runtime.go:85 +0xae
internal/poll.(*pollDesc).waitRead(0xc4203b0398, 0xffffffffffffff00, 0x0, 0x0)
	/server/golang/src/internal/poll/fd_poll_runtime.go:90 +0x3d
internal/poll.(*FD).Accept(0xc4203b0380, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0)
	/server/golang/src/internal/poll/fd_unix.go:335 +0x1e2
net.(*netFD).accept(0xc4203b0380, 0xa54be8, 0xc4204a7db8, 0x40239b)
	/server/golang/src/net/fd_unix.go:238 +0x42
net.(*TCPListener).accept(0xc42000cef0, 0x974d60, 0xc4204a7de8, 0x401137)
	/server/golang/src/net/tcpsock_posix.go:136 +0x2e
net.(*TCPListener).AcceptTCP(0xc42000cef0, 0xc4204a7e30, 0xc4204a7e38, 0xc4204a7e28)
	/server/golang/src/net/tcpsock.go:234 +0x49
net/http.tcpKeepAliveListener.Accept(0xc42000cef0, 0xa54438, 0xc420a56140, 0xd38080, 0xc42084c6f0)
	/server/golang/src/net/http/server.go:3120 +0x2f
net/http.(*Server).Serve(0xc420104680, 0xd37a40, 0xc42000cef0, 0x0, 0x0)
	/server/golang/src/net/http/server.go:2695 +0x1b2
net/http.(*Server).ListenAndServe(0xc420104680, 0xa9c5b0, 0xc4209b1f80)
	/server/golang/src/net/http/server.go:2636 +0xa9
github.com/zhujingfa/docker-registry-manager/app/vendor/github.com/astaxie/beego.(*App).Run.func4(0xc420043ee0, 0xc4206435c0, 0x5, 0xc4206fc3f0)
	/data/go/src/github.com/zhujingfa/docker-registry-manager/app/vendor/github.com/astaxie/beego/app.go:178 +0x2ca
created by github.com/zhujingfa/docker-registry-manager/app/vendor/github.com/astaxie/beego.(*App).Run
	/data/go/src/github.com/zhujingfa/docker-registry-manager/app/vendor/github.com/astaxie/beego/app.go:160 +0x4fa

goroutine 1367 [select, 3 minutes]:
net/http.(*persistConn).writeLoop(0xc4207ad8c0)
	/server/golang/src/net/http/transport.go:1759 +0x165
created by net/http.(*Transport).dialConn
	/server/golang/src/net/http/transport.go:1187 +0xa53

goroutine 1610 [IO wait]:
internal/poll.runtime_pollWait(0x7faccf391af0, 0x72, 0x0)
	/server/golang/src/runtime/netpoll.go:173 +0x57
internal/poll.(*pollDesc).wait(0xc420478c18, 0x72, 0xffffffffffffff00, 0xd32d00, 0xd2e420)
	/server/golang/src/internal/poll/fd_poll_runtime.go:85 +0xae
internal/poll.(*pollDesc).waitRead(0xc420478c18, 0xc420d1f000, 0x200, 0x200)
	/server/golang/src/internal/poll/fd_poll_runtime.go:90 +0x3d
internal/poll.(*FD).Read(0xc420478c00, 0xc420d1f000, 0x200, 0x200, 0x0, 0x0, 0x0)
	/server/golang/src/internal/poll/fd_unix.go:126 +0x18a
net.(*netFD).Read(0xc420478c00, 0xc420d1f000, 0x200, 0x200, 0x29, 0xc420374bf0, 0x443897)
	/server/golang/src/net/fd_unix.go:202 +0x52
net.(*conn).Read(0xc4205d6818, 0xc420d1f000, 0x200, 0x200, 0x0, 0x0, 0x0)
	/server/golang/src/net/net.go:176 +0x6d
net.(*dnsPacketConn).dnsRoundTrip(0xc42046d280, 0xc420478b00, 0xd78800, 0xd78800, 0x0)
	/server/golang/src/net/dnsclient_unix.go:57 +0xde
net.(*Resolver).exchange(0xd771b0, 0xd38000, 0xc4200140c0, 0xc4203aa700, 0xc, 0xc4209e1740, 0x18, 0x1, 0x12a05f200, 0x0, ...)
	/server/golang/src/net/dnsclient_unix.go:138 +0x48e
net.(*Resolver).tryOneName(0xd771b0, 0xd38000, 0xc4200140c0, 0xc4203a06e0, 0xc4209e1740, 0x18, 0x1, 0x50, 0x1, 0x6, ...)
	/server/golang/src/net/dnsclient_unix.go:161 +0x148
net.(*Resolver).goLookupIPCNAMEOrder.func1(0xd771b0, 0xd38000, 0xc4200140c0, 0xc4203a06e0, 0xc42046d1f0, 0xc42037d3e0, 0xc400000001)
	/server/golang/src/net/dnsclient_unix.go:483 +0x8d
created by net.(*Resolver).goLookupIPCNAMEOrder
	/server/golang/src/net/dnsclient_unix.go:482 +0x21c

goroutine 1604 [IO wait]:
internal/poll.runtime_pollWait(0x7faccf391f70, 0x72, 0x0)
	/server/golang/src/runtime/netpoll.go:173 +0x57
internal/poll.(*pollDesc).wait(0xc42052ac98, 0x72, 0xffffffffffffff00, 0xd32d00, 0xd2e420)
	/server/golang/src/internal/poll/fd_poll_runtime.go:85 +0xae
internal/poll.(*pollDesc).waitRead(0xc42052ac98, 0xc420d07d00, 0x1, 0x1)
	/server/golang/src/internal/poll/fd_poll_runtime.go:90 +0x3d
internal/poll.(*FD).Read(0xc42052ac80, 0xc420d07d21, 0x1, 0x1, 0x0, 0x0, 0x0)
	/server/golang/src/internal/poll/fd_unix.go:126 +0x18a
net.(*netFD).Read(0xc42052ac80, 0xc420d07d21, 0x1, 0x1, 0xc4208100c0, 0x205c2328, 0xc420023748)
	/server/golang/src/net/fd_unix.go:202 +0x52
net.(*conn).Read(0xc420192608, 0xc420d07d21, 0x1, 0x1, 0x0, 0x0, 0x0)
	/server/golang/src/net/net.go:176 +0x6d
net/http.(*connReader).backgroundRead(0xc420d07d10)
	/server/golang/src/net/http/server.go:660 +0x62
created by net/http.(*connReader).startBackgroundRead
	/server/golang/src/net/http/server.go:656 +0xd8

goroutine 1580 [chan receive]:
net.(*Resolver).goLookupIPCNAMEOrder(0xd771b0, 0xd38000, 0xc4200140c0, 0xc4203b6038, 0x17, 0x2, 0xd31340, 0xc4203ba990, 0x1, 0xc4201fd838, ...)
	/server/golang/src/net/dnsclient_unix.go:489 +0x374
net.(*Resolver).goLookupHostOrder(0xd771b0, 0xd38000, 0xc4200140c0, 0xc4203b6038, 0x17, 0x2, 0x0, 0x0, 0x59c40, 0xd992e0, ...)
	/server/golang/src/net/dnsclient_unix.go:425 +0x126
net.(*Resolver).lookupHost(0xd771b0, 0xd38000, 0xc4200140c0, 0xc4203b6038, 0x17, 0xa11400, 0x7fac00059c40, 0x953c40, 0xc4201fd918, 0x43be33)
	/server/golang/src/net/lookup_unix.go:86 +0x96
net.(*Resolver).LookupHost(0xd771b0, 0xd38000, 0xc4200140c0, 0xc4203b6038, 0x17, 0x7faccf39c238, 0x0, 0x8, 0xd79600, 0x7faccf3ed000)
	/server/golang/src/net/lookup.go:152 +0x241
net.LookupHost(0xc4203b6038, 0x17, 0xb80000c4201fdc18, 0xc4201fd9d8, 0x4110a8, 0x20, 0xc4206ab980)
	/server/golang/src/net/lookup.go:138 +0x5d
github.com/zhujingfa/docker-registry-manager/app/models.(*Registry).IP(0xc4203b1780, 0x0, 0x0)
	/data/go/src/github.com/zhujingfa/docker-registry-manager/app/models/registry.go:49 +0x38
reflect.Value.call(0xa11400, 0xc4203b1780, 0x1213, 0xa26e11, 0x4, 0xd98460, 0x0, 0x0, 0xa23d20, 0x1, ...)
	/server/golang/src/reflect/value.go:434 +0x905
reflect.Value.Call(0xa11400, 0xc4203b1780, 0x1213, 0xd98460, 0x0, 0x0, 0xd3e120, 0xc42037c360, 0xc42037c360)
	/server/golang/src/reflect/value.go:302 +0xa4
text/template.(*state).evalCall(0xc4201feed8, 0xa11400, 0xc4203b1780, 0x16, 0xa11400, 0xc4203b1780, 0x1213, 0xd39aa0, 0xc4202aa6f0, 0xc420acf95a, ...)
	/server/golang/src/text/template/exec.go:670 +0x580
text/template.(*state).evalField(0xc4201feed8, 0xa11400, 0xc4203b1780, 0x16, 0xc420acf95a, 0x2, 0xd39aa0, 0xc4202aa6f0, 0x0, 0x0, ...)
	/server/golang/src/text/template/exec.go:560 +0xd38
text/template.(*state).evalFieldChain(0xc4201feed8, 0xa11400, 0xc4203b1780, 0x16, 0xa11400, 0xc4203b1780, 0x16, 0xd39aa0, 0xc4202aa6f0, 0xc42027e950, ...)
	/server/golang/src/text/template/exec.go:528 +0x22a
text/template.(*state).evalVariableNode(0xc4201feed8, 0xa11400, 0xc4203b1780, 0x16, 0xc4202aa6f0, 0x0, 0x0, 0x0, 0x0, 0x0, ...)
	/server/golang/src/text/template/exec.go:516 +0x192
text/template.(*state).evalArg(0xc4201feed8, 0xa11400, 0xc4203b1780, 0x16, 0xd3e120, 0x974c40, 0xd39aa0, 0xc4202aa6f0, 0xc42046cde0, 0xc42046d140, ...)
	/server/golang/src/text/template/exec.go:746 +0xb09
text/template.(*state).evalCall(0xc4201feed8, 0xa11400, 0xc4203b1780, 0x16, 0x97e120, 0xa53ce0, 0x13, 0xd395c0, 0xc4202aa660, 0xc420543759, ...)
	/server/golang/src/text/template/exec.go:645 +0x281
text/template.(*state).evalFunction(0xc4201feed8, 0xa11400, 0xc4203b1780, 0x16, 0xc4202aa690, 0xd395c0, 0xc4202aa660, 0xc42037e640, 0x3, 0x4, ...)
	/server/golang/src/text/template/exec.go:538 +0x176
text/template.(*state).evalCommand(0xc4201feed8, 0xa11400, 0xc4203b1780, 0x16, 0xc4202aa660, 0x0, 0x0, 0x0, 0x951400, 0x3, ...)
	/server/golang/src/text/template/exec.go:435 +0x53a
text/template.(*state).evalPipeline(0xc4201feed8, 0xa11400, 0xc4203b1780, 0x16, 0xc420372e60, 0xc4202aa390, 0xc42037f728, 0x81)
	/server/golang/src/text/template/exec.go:408 +0x115
text/template.(*state).walkIfOrWith(0xc4201feed8, 0xa, 0xa11400, 0xc4203b1780, 0x16, 0xc420372e60, 0xc4202aa720, 0x0)
	/server/golang/src/text/template/exec.go:263 +0xbc
text/template.(*state).walk(0xc4201feed8, 0xa11400, 0xc4203b1780, 0x16, 0xd39740, 0xc42037e780)
	/server/golang/src/text/template/exec.go:239 +0x318
text/template.(*state).walk(0xc4201feed8, 0xa11400, 0xc4203b1780, 0x16, 0xd397a0, 0xc4204773e0)
	/server/golang/src/text/template/exec.go:242 +0x11d
text/template.(*state).walkRange.func1(0x956a40, 0xc42046cf80, 0x98, 0xa11400, 0xc4203b1780, 0x16)
	/server/golang/src/text/template/exec.go:329 +0x12f
text/template.(*state).walkRange(0xc4201feed8, 0x981360, 0xc420476e70, 0x15, 0xc42037ecc0)
	/server/golang/src/text/template/exec.go:346 +0x592
text/template.(*state).walk(0xc4201feed8, 0x981360, 0xc420476e70, 0x15, 0xd39920, 0xc42037ecc0)
	/server/golang/src/text/template/exec.go:245 +0x451
text/template.(*state).walk(0xc4201feed8, 0x981360, 0xc420476e70, 0x15, 0xd397a0, 0xc420477260)
	/server/golang/src/text/template/exec.go:242 +0x11d
text/template.(*state).walkTemplate(0xc4201ff0a8, 0x981360, 0xc420476e70, 0x15, 0xc42037ef40)
	/server/golang/src/text/template/exec.go:391 +0x25d
text/template.(*state).walk(0xc4201ff0a8, 0x981360, 0xc420476e70, 0x15, 0xd399e0, 0xc42037ef40)
	/server/golang/src/text/template/exec.go:247 +0x1f6
text/template.(*state).walk(0xc4201ff0a8, 0x981360, 0xc420476e70, 0x15, 0xd397a0, 0xc4202ab290)
	/server/golang/src/text/template/exec.go:242 +0x11d
text/template.(*state).walkTemplate(0xc4201ff278, 0x981360, 0xc420476e70, 0x15, 0xc42048bd40)
	/server/golang/src/text/template/exec.go:391 +0x25d
text/template.(*state).walk(0xc4201ff278, 0x981360, 0xc420476e70, 0x15, 0xd399e0, 0xc42048bd40)
	/server/golang/src/text/template/exec.go:247 +0x1f6
text/template.(*state).walk(0xc4201ff278, 0x981360, 0xc420476e70, 0x15, 0xd397a0, 0xc4204771d0)
	/server/golang/src/text/template/exec.go:242 +0x11d
text/template.(*Template).execute(0xc42048bd00, 0xd2fc80, 0xc42022b570, 0x981360, 0xc420476e70, 0x0, 0x0)
	/server/golang/src/text/template/exec.go:197 +0x1f9
text/template.(*Template).Execute(0xc42048bd00, 0xd2fc80, 0xc42022b570, 0x981360, 0xc420476e70, 0x0, 0xc42030e8a8)
	/server/golang/src/text/template/exec.go:180 +0x53
html/template.(*Template).ExecuteTemplate(0xc420477170, 0xd2fc80, 0xc42022b570, 0xa2f5c0, 0xe, 0x981360, 0xc420476e70, 0xd79600, 0x0)
	/server/golang/src/html/template/template.go:137 +0xa6
github.com/zhujingfa/docker-registry-manager/app/vendor/github.com/astaxie/beego.ExecuteViewPathTemplate(0xd2fc80, 0xc42022b570, 0xa2f5c0, 0xe, 0xa28444, 0x5, 0x981360, 0xc420476e70, 0x0, 0x0)
	/data/go/src/github.com/zhujingfa/docker-registry-manager/app/vendor/github.com/astaxie/beego/template.go:64 +0x193
github.com/zhujingfa/docker-registry-manager/app/vendor/github.com/astaxie/beego.(*Controller).renderTemplate(0xc420bad1d0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, ...)
	/data/go/src/github.com/zhujingfa/docker-registry-manager/app/vendor/github.com/astaxie/beego/controller.go:263 +0x431
github.com/zhujingfa/docker-registry-manager/app/vendor/github.com/astaxie/beego.(*Controller).RenderBytes(0xc420bad1d0, 0xc420e037b0, 0x83db71, 0x981360, 0xc420476e70, 0xc420e037a0)
	/data/go/src/github.com/zhujingfa/docker-registry-manager/app/vendor/github.com/astaxie/beego/controller.go:214 +0x8e
github.com/zhujingfa/docker-registry-manager/app/vendor/github.com/astaxie/beego.(*Controller).Render(0xc420bad1d0, 0x0, 0x0)
	/data/go/src/github.com/zhujingfa/docker-registry-manager/app/vendor/github.com/astaxie/beego/controller.go:194 +0x3e
github.com/zhujingfa/docker-registry-manager/app/vendor/github.com/astaxie/beego.(*ControllerRegister).ServeHTTP(0xc4200106e0, 0xd375c0, 0xc4206ac1c0, 0xc4201b0200)
	/data/go/src/github.com/zhujingfa/docker-registry-manager/app/vendor/github.com/astaxie/beego/router.go:827 +0x1dbe
net/http.serverHandler.ServeHTTP(0xc420104680, 0xd375c0, 0xc4206ac1c0, 0xc4201b0200)
	/server/golang/src/net/http/server.go:2619 +0xb4
net/http.(*conn).serve(0xc420a56140, 0xd37fc0, 0xc420dc5640)
	/server/golang/src/net/http/server.go:1801 +0x71d
created by net/http.(*Server).Serve
	/server/golang/src/net/http/server.go:2720 +0x288

rax    0xca
rbx    0xd79600
rcx    0x45c3d3
rdx    0x0
rdi    0xd79710
rsi    0x0
rbp    0x7ffc456c1a10
rsp    0x7ffc456c19c8
r8     0x0
r9     0x0
r10    0x0
r11    0x286
r12    0x0
r13    0xf1
r14    0x11
r15    0x0
rip    0x45c3d1
rflags 0x286
cs     0x33
fs     0x0
gs     0x0
```