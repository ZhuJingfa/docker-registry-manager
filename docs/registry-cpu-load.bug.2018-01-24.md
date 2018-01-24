# BUG cpu占用偏高

平均占用cpu 15%左右，一天运行后，机器就负载10几了。

初步分析应该是models.AddRegistry.Refresh go routine函数导致的，配置的五分钟，网页添加的刷新10秒钟，获取所有数据，可怕。

1. 所有repo列表接口： https://registry.alishui.com/v2/_catalog
    
1. registry.tags url=https://registry.alishui.com:443/v2/kong/tags/list repository=kong

    {"name":"kong","tags":["0.12.1-alpine"]}     
    
1. registry.manifest.get url=https://registry.alishui.com:443/v2/jenkins/jenkins/manifests/lts-alpine repository=jenkins/jenkins reference=lts-alpine 

    method DELETE: https://registry.alishui.com:443/v2/jenkins/jenkins/manifests/lts-alpine

1. registry.layer.check url=https://registry.alishui.com:443/v2/kong/blobs/sha256:b54c874d32ceed59274451ee7e0cbb4a04e6d968eec92345a06eacf06aeaf34e repository=kong digest=sha256:b54c874d32ceed59274451ee7e0cbb4a04e6d968eec92345a06eacf06aeaf34e



```
root@b0bb9826271c:/usr/src$ kill -3 30
root@b0bb9826271c:/usr/src$ SIGQUIT: quit
PC=0x45c3d1 m=0 sigcode=0

goroutine 0 [idle]:
runtime.futex(0xd79710, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x7fff27521d60, 0x40f85b, ...)
	/server/golang/src/runtime/sys_linux_amd64.s:438 +0x21
runtime.futexsleep(0xd79710, 0x7fff00000000, 0xffffffffffffffff)
	/server/golang/src/runtime/os_linux.go:45 +0x62
runtime.notesleep(0xd79710)
	/server/golang/src/runtime/lock_futex.go:151 +0x9b
runtime.stopm()
	/server/golang/src/runtime/proc.go:1680 +0xe5
runtime.gcstopm()
	/server/golang/src/runtime/proc.go:1884 +0xb7
runtime.schedule()
	/server/golang/src/runtime/proc.go:2219 +0x2c0
runtime.goschedImpl(0xc420000180)
	/server/golang/src/runtime/proc.go:2333 +0xff
runtime.gopreempt_m(0xc420000180)
	/server/golang/src/runtime/proc.go:2361 +0x36
runtime.newstack(0x0)
	/server/golang/src/runtime/stack.go:1042 +0x2d1
runtime.morestack()
	/server/golang/src/runtime/asm_amd64.s:415 +0x86

goroutine 1 [running]:
	goroutine running on other thread; stack unavailable

goroutine 18 [syscall]:
os/signal.signal_recv(0x0)
	/server/golang/src/runtime/sigqueue.go:131 +0xa6
os/signal.loop()
	/server/golang/src/os/signal/signal_unix.go:22 +0x22
created by os/signal.init.0
	/server/golang/src/os/signal/signal_unix.go:28 +0x41

goroutine 13 [chan receive]:
github.com/snagles/docker-registry-manager/app/models.AddRegistry.func1(0xc4201080e0)
	/data/go/src/github.com/snagles/docker-registry-manager/app/models/registry.go:247 +0x90
created by github.com/snagles/docker-registry-manager/app/models.AddRegistry
	/data/go/src/github.com/snagles/docker-registry-manager/app/models/registry.go:246 +0x49c

goroutine 10 [IO wait]:
internal/poll.runtime_pollWait(0x7fd49cd66eb0, 0x72, 0x0)
	/server/golang/src/runtime/netpoll.go:173 +0x57
internal/poll.(*pollDesc).wait(0xc4203d4418, 0x72, 0xffffffffffffff00, 0xd32d00, 0xd2e420)
	/server/golang/src/internal/poll/fd_poll_runtime.go:85 +0xae
internal/poll.(*pollDesc).waitRead(0xc4203d4418, 0xc4203ea000, 0x2000, 0x2000)
	/server/golang/src/internal/poll/fd_poll_runtime.go:90 +0x3d
internal/poll.(*FD).Read(0xc4203d4400, 0xc4203ea000, 0x2000, 0x2000, 0x0, 0x0, 0x0)
	/server/golang/src/internal/poll/fd_unix.go:126 +0x18a
net.(*netFD).Read(0xc4203d4400, 0xc4203ea000, 0x2000, 0x2000, 0x1, 0x4, 0xc4204a5788)
	/server/golang/src/net/fd_unix.go:202 +0x52
net.(*conn).Read(0xc420072030, 0xc4203ea000, 0x2000, 0x2000, 0x0, 0x0, 0x0)
	/server/golang/src/net/net.go:176 +0x6d
crypto/tls.(*block).readFromUntil(0xc420268000, 0x7fd49cbc2028, 0xc420072030, 0x5, 0xc420072030, 0xc4201562d0)
	/server/golang/src/crypto/tls/conn.go:488 +0x95
crypto/tls.(*Conn).readRecord(0xc420098380, 0xa56117, 0xc4200984a0, 0x42d88b)
	/server/golang/src/crypto/tls/conn.go:590 +0xe0
crypto/tls.(*Conn).Read(0xc420098380, 0xc42025c000, 0x1000, 0x1000, 0x0, 0x0, 0x0)
	/server/golang/src/crypto/tls/conn.go:1134 +0x110
net/http.(*persistConn).Read(0xc420348240, 0xc42025c000, 0x1000, 0x1000, 0xc4201e6220, 0xc420094ef8, 0x455090)
	/server/golang/src/net/http/transport.go:1391 +0x140
bufio.(*Reader).fill(0xc4201046c0)
	/server/golang/src/bufio/bufio.go:97 +0x11a
bufio.(*Reader).Peek(0xc4201046c0, 0x1, 0x0, 0x0, 0x0, 0xc4200168a0, 0x0)
	/server/golang/src/bufio/bufio.go:129 +0x3a
net/http.(*persistConn).readLoop(0xc420348240)
	/server/golang/src/net/http/transport.go:1539 +0x185
created by net/http.(*Transport).dialConn
	/server/golang/src/net/http/transport.go:1186 +0xa2e

goroutine 11 [select]:
net/http.(*persistConn).writeLoop(0xc420348240)
	/server/golang/src/net/http/transport.go:1759 +0x165
created by net/http.(*Transport).dialConn
	/server/golang/src/net/http/transport.go:1187 +0xa53

goroutine 40 [select]:
net/http.(*persistConn).writeLoop(0xc420096120)
	/server/golang/src/net/http/transport.go:1759 +0x165
created by net/http.(*Transport).dialConn
	/server/golang/src/net/http/transport.go:1187 +0xa53

goroutine 39 [IO wait]:
internal/poll.runtime_pollWait(0x7fd49cd66f70, 0x72, 0x0)
	/server/golang/src/runtime/netpoll.go:173 +0x57
internal/poll.(*pollDesc).wait(0xc4203d2a18, 0x72, 0xffffffffffffff00, 0xd32d00, 0xd2e420)
	/server/golang/src/internal/poll/fd_poll_runtime.go:85 +0xae
internal/poll.(*pollDesc).waitRead(0xc4203d2a18, 0xc42050e000, 0x4000, 0x4000)
	/server/golang/src/internal/poll/fd_poll_runtime.go:90 +0x3d
internal/poll.(*FD).Read(0xc4203d2a00, 0xc42050e000, 0x4000, 0x4000, 0x0, 0x0, 0x0)
	/server/golang/src/internal/poll/fd_unix.go:126 +0x18a
net.(*netFD).Read(0xc4203d2a00, 0xc42050e000, 0x4000, 0x4000, 0x1, 0x5, 0xc4204a4788)
	/server/golang/src/net/fd_unix.go:202 +0x52
net.(*conn).Read(0xc420072118, 0xc42050e000, 0x4000, 0x4000, 0x0, 0x0, 0x0)
	/server/golang/src/net/net.go:176 +0x6d
crypto/tls.(*block).readFromUntil(0xc4203da930, 0x7fd49cbc2028, 0xc420072118, 0x5, 0xc420072118, 0xc42000e750)
	/server/golang/src/crypto/tls/conn.go:488 +0x95
crypto/tls.(*Conn).readRecord(0xc420028e00, 0xa56117, 0xc420028f20, 0x42d88b)
	/server/golang/src/crypto/tls/conn.go:590 +0xe0
crypto/tls.(*Conn).Read(0xc420028e00, 0xc420177000, 0x1000, 0x1000, 0x0, 0x0, 0x0)
	/server/golang/src/crypto/tls/conn.go:1134 +0x110
net/http.(*persistConn).Read(0xc420096120, 0xc420177000, 0x1000, 0x1000, 0xc4201e70c0, 0xc420017558, 0x455090)
	/server/golang/src/net/http/transport.go:1391 +0x140
bufio.(*Reader).fill(0xc4203e09c0)
	/server/golang/src/bufio/bufio.go:97 +0x11a
bufio.(*Reader).Peek(0xc4203e09c0, 0x1, 0x0, 0x0, 0x0, 0xc4200949c0, 0x0)
	/server/golang/src/bufio/bufio.go:129 +0x3a
net/http.(*persistConn).readLoop(0xc420096120)
	/server/golang/src/net/http/transport.go:1539 +0x185
created by net/http.(*Transport).dialConn
	/server/golang/src/net/http/transport.go:1186 +0xa2e

rax    0xca
rbx    0xd79600
rcx    0x45c3d3
rdx    0x0
rdi    0xd79710
rsi    0x0
rbp    0x7fff27521d28
rsp    0x7fff27521ce0
r8     0x0
r9     0x0
r10    0x0
r11    0x286
r12    0xc
r13    0xc4203d46e0
r14    0x1
r15    0x29
rip    0x45c3d1
rflags 0x286
cs     0x33
fs     0x0
gs     0x0

[1]+  Exit 2                  /app -c /conf.yml
root@b0bb9826271c:/usr/src$
```


```
DEBU[0330] registry.layer.check url=https://registry.alishui.com:443/v2/fanghui/system/python3/blobs/sha256:b89fd87b23636cf742fa893e8da55747952bf360820925c3287fe4b128dcf11d repository=fanghui/system/python3 digest=sha256:b89fd87b23636cf742fa893e8da55747952bf360820925c3287fe4b128dcf11d 
DEBU[0330] registry.manifest.get url=https://registry.alishui.com:443/v2/fanghui/system/python3/manifests/latest repository=fanghui/system/python3 reference=latest 
DEBU[0330] registry.tags url=https://registry.alishui.com:443/v2/fanghui/system/python3-dev/tags/list repository=fanghui/system/python3-dev 
DEBU[0330] registry.manifest.get url=https://registry.alishui.com:443/v2/fanghui/system/python3-dev/manifests/3.6 repository=fanghui/system/python3-dev reference=3.6 
DEBU[0330] registry.layer.check url=https://registry.alishui.com:443/v2/fanghui/system/python3-dev/blobs/sha256:4ff5d061c205bb14df5291e158ec3db53645a94fe31eb05ab1399a1dd659ee43 repository=fanghui/system/python3-dev digest=sha256:4ff5d061c205bb14df5291e158ec3db53645a94fe31eb05ab1399a1dd659ee43 
DEBU[0330] registry.manifest.get url=https://registry.alishui.com:443/v2/fanghui/system/python3-dev/manifests/3.6 repository=fanghui/system/python3-dev reference=3.6 
DEBU[0330] registry.manifest.get url=https://registry.alishui.com:443/v2/fanghui/system/python3-dev/manifests/latest repository=fanghui/system/python3-dev reference=latest 
DEBU[0330] registry.layer.check url=https://registry.alishui.com:443/v2/fanghui/system/python3-dev/blobs/sha256:4ff5d061c205bb14df5291e158ec3db53645a94fe31eb05ab1399a1dd659ee43 repository=fanghui/system/python3-dev digest=sha256:4ff5d061c205bb14df5291e158ec3db53645a94fe31eb05ab1399a1dd659ee43 
DEBU[0330] registry.manifest.get url=https://registry.alishui.com:443/v2/fanghui/system/python3-dev/manifests/latest repository=fanghui/system/python3-dev reference=latest 
DEBU[0330] registry.tags url=https://registry.alishui.com:443/v2/fanghui/system/redis/tags/list repository=fanghui/system/redis 
DEBU[0331] registry.manifest.get url=https://registry.alishui.com:443/v2/fanghui/system/redis/manifests/latest repository=fanghui/system/redis reference=latest 
DEBU[0331] registry.layer.check url=https://registry.alishui.com:443/v2/fanghui/system/redis/blobs/sha256:b7fddd0da2bb7c83de38409ff754b53022aabd5e366e19ad2402854c82d0a8d0 repository=fanghui/system/redis digest=sha256:b7fddd0da2bb7c83de38409ff754b53022aabd5e366e19ad2402854c82d0a8d0 
DEBU[0331] registry.manifest.get url=https://registry.alishui.com:443/v2/fanghui/system/redis/manifests/latest repository=fanghui/system/redis reference=latest 
DEBU[0331] registry.manifest.get url=https://registry.alishui.com:443/v2/fanghui/system/redis/manifests/3.2 repository=fanghui/system/redis reference=3.2 
DEBU[0331] registry.layer.check url=https://registry.alishui.com:443/v2/fanghui/system/redis/blobs/sha256:b7fddd0da2bb7c83de38409ff754b53022aabd5e366e19ad2402854c82d0a8d0 repository=fanghui/system/redis digest=sha256:b7fddd0da2bb7c83de38409ff754b53022aabd5e366e19ad2402854c82d0a8d0 
DEBU[0331] registry.manifest.get url=https://registry.alishui.com:443/v2/fanghui/system/redis/manifests/3.2 repository=fanghui/system/redis reference=3.2 
DEBU[0331] registry.tags url=https://registry.alishui.com:443/v2/fanghui/system/registry-manager/tags/list repository=fanghui/system/registry-manager 
DEBU[0331] registry.manifest.get url=https://registry.alishui.com:443/v2/fanghui/system/registry-manager/manifests/latest repository=fanghui/system/registry-manager reference=latest 
DEBU[0331] registry.layer.check url=https://registry.alishui.com:443/v2/fanghui/system/registry-manager/blobs/sha256:1bd56d4da1f58f0807f5217809a72c042dd7e8eb77263335b67d711968fbdcfe repository=fanghui/system/registry-manager digest=sha256:1bd56d4da1f58f0807f5217809a72c042dd7e8eb77263335b67d711968fbdcfe 
DEBU[0331] registry.manifest.get url=https://registry.alishui.com:443/v2/fanghui/system/registry-manager/manifests/latest repository=fanghui/system/registry-manager reference=latest 
DEBU[0331] registry.tags url=https://registry.alishui.com:443/v2/fanghui/system/scratch/tags/list repository=fanghui/system/scratch 
DEBU[0331] registry.manifest.get url=https://registry.alishui.com:443/v2/fanghui/system/scratch/manifests/1.0 repository=fanghui/system/scratch reference=1.0 
DEBU[0331] registry.layer.check url=https://registry.alishui.com:443/v2/fanghui/system/scratch/blobs/sha256:6989d8ab72dba4d70d113d10123d9a9ecc98a2ab4bf9a18732e403867ac7188f repository=fanghui/system/scratch digest=sha256:6989d8ab72dba4d70d113d10123d9a9ecc98a2ab4bf9a18732e403867ac7188f 
DEBU[0331] registry.manifest.get url=https://registry.alishui.com:443/v2/fanghui/system/scratch/manifests/1.0 repository=fanghui/system/scratch reference=1.0 
DEBU[0331] registry.manifest.get url=https://registry.alishui.com:443/v2/fanghui/system/scratch/manifests/latest repository=fanghui/system/scratch reference=latest 
DEBU[0331] registry.layer.check url=https://registry.alishui.com:443/v2/fanghui/system/scratch/blobs/sha256:6989d8ab72dba4d70d113d10123d9a9ecc98a2ab4bf9a18732e403867ac7188f repository=fanghui/system/scratch digest=sha256:6989d8ab72dba4d70d113d10123d9a9ecc98a2ab4bf9a18732e403867ac7188f 
DEBU[0331] registry.manifest.get url=https://registry.alishui.com:443/v2/fanghui/system/scratch/manifests/latest repository=fanghui/system/scratch reference=latest 
DEBU[0331] registry.tags url=https://registry.alishui.com:443/v2/fanghui/system/sshd/tags/list repository=fanghui/system/sshd 
DEBU[0331] registry.manifest.get url=https://registry.alishui.com:443/v2/fanghui/system/sshd/manifests/7.5 repository=fanghui/system/sshd reference=7.5 
DEBU[0331] registry.layer.check url=https://registry.alishui.com:443/v2/fanghui/system/sshd/blobs/sha256:6b42608684e13e4a0b3b91d654ff7e6e2044fddbbfef2f5284f17b41b8894a6a repository=fanghui/system/sshd digest=sha256:6b42608684e13e4a0b3b91d654ff7e6e2044fddbbfef2f5284f17b41b8894a6a 
DEBU[0331] registry.manifest.get url=https://registry.alishui.com:443/v2/fanghui/system/sshd/manifests/7.5 repository=fanghui/system/sshd reference=7.5 
DEBU[0331] registry.manifest.get url=https://registry.alishui.com:443/v2/fanghui/system/sshd/manifests/latest repository=fanghui/system/sshd reference=latest 
DEBU[0331] registry.layer.check url=https://registry.alishui.com:443/v2/fanghui/system/sshd/blobs/sha256:6b42608684e13e4a0b3b91d654ff7e6e2044fddbbfef2f5284f17b41b8894a6a repository=fanghui/system/sshd digest=sha256:6b42608684e13e4a0b3b91d654ff7e6e2044fddbbfef2f5284f17b41b8894a6a 
DEBU[0331] registry.manifest.get url=https://registry.alishui.com:443/v2/fanghui/system/sshd/manifests/latest repository=fanghui/system/sshd reference=latest 
DEBU[0331] registry.tags url=https://registry.alishui.com:443/v2/fanghui/system/strongswan/tags/list repository=fanghui/system/strongswan 
DEBU[0331] registry.manifest.get url=https://registry.alishui.com:443/v2/fanghui/system/strongswan/manifests/latest repository=fanghui/system/strongswan reference=latest 
DEBU[0331] registry.layer.check url=https://registry.alishui.com:443/v2/fanghui/system/strongswan/blobs/sha256:52ae80b80f35a22bf82f58cdc290d0a3b6652b97a69311ae4117ea1b2faed5b4 repository=fanghui/system/strongswan digest=sha256:52ae80b80f35a22bf82f58cdc290d0a3b6652b97a69311ae4117ea1b2faed5b4 
DEBU[0331] registry.manifest.get url=https://registry.alishui.com:443/v2/fanghui/system/strongswan/manifests/latest repository=fanghui/system/strongswan reference=latest 
DEBU[0331] registry.tags url=https://registry.alishui.com:443/v2/fanghui/system/trafficserver/tags/list repository=fanghui/system/trafficserver 
DEBU[0331] registry.manifest.get url=https://registry.alishui.com:443/v2/fanghui/system/trafficserver/manifests/7.1 repository=fanghui/system/trafficserver reference=7.1 
DEBU[0331] registry.layer.check url=https://registry.alishui.com:443/v2/fanghui/system/trafficserver/blobs/sha256:15b73ae5509c324cf352b2e1e262d50bd7314561abfb1e400cc2ddf82fc81a30 repository=fanghui/system/trafficserver digest=sha256:15b73ae5509c324cf352b2e1e262d50bd7314561abfb1e400cc2ddf82fc81a30 
DEBU[0331] registry.manifest.get url=https://registry.alishui.com:443/v2/fanghui/system/trafficserver/manifests/7.1 repository=fanghui/system/trafficserver reference=7.1 
DEBU[0331] registry.manifest.get url=https://registry.alishui.com:443/v2/fanghui/system/trafficserver/manifests/latest repository=fanghui/system/trafficserver reference=latest 
DEBU[0331] registry.layer.check url=https://registry.alishui.com:443/v2/fanghui/system/trafficserver/blobs/sha256:15b73ae5509c324cf352b2e1e262d50bd7314561abfb1e400cc2ddf82fc81a30 repository=fanghui/system/trafficserver digest=sha256:15b73ae5509c324cf352b2e1e262d50bd7314561abfb1e400cc2ddf82fc81a30 
DEBU[0332] registry.manifest.get url=https://registry.alishui.com:443/v2/fanghui/system/trafficserver/manifests/latest repository=fanghui/system/trafficserver reference=latest 
DEBU[0332] registry.tags url=https://registry.alishui.com:443/v2/fanghui/system/ubuntu/tags/list repository=fanghui/system/ubuntu 
DEBU[0332] registry.manifest.get url=https://registry.alishui.com:443/v2/fanghui/system/ubuntu/manifests/latest repository=fanghui/system/ubuntu reference=latest 
DEBU[0332] registry.layer.check url=https://registry.alishui.com:443/v2/fanghui/system/ubuntu/blobs/sha256:e0efdf5ea7fc2171cac23ab5924fcfda2c3b3367bbd2df4bc2580983355118b2 repository=fanghui/system/ubuntu digest=sha256:e0efdf5ea7fc2171cac23ab5924fcfda2c3b3367bbd2df4bc2580983355118b2 
DEBU[0332] registry.manifest.get url=https://registry.alishui.com:443/v2/fanghui/system/ubuntu/manifests/latest repository=fanghui/system/ubuntu reference=latest 
DEBU[0332] registry.manifest.get url=https://registry.alishui.com:443/v2/fanghui/system/ubuntu/manifests/16.04 repository=fanghui/system/ubuntu reference=16.04 
DEBU[0332] registry.layer.check url=https://registry.alishui.com:443/v2/fanghui/system/ubuntu/blobs/sha256:e0efdf5ea7fc2171cac23ab5924fcfda2c3b3367bbd2df4bc2580983355118b2 repository=fanghui/system/ubuntu digest=sha256:e0efdf5ea7fc2171cac23ab5924fcfda2c3b3367bbd2df4bc2580983355118b2 
DEBU[0332] registry.manifest.get url=https://registry.alishui.com:443/v2/fanghui/system/ubuntu/manifests/16.04 repository=fanghui/system/ubuntu reference=16.04 
DEBU[0332] registry.tags url=https://registry.alishui.com:443/v2/fanghui/system/zentao/tags/list repository=fanghui/system/zentao 
DEBU[0332] registry.manifest.get url=https://registry.alishui.com:443/v2/fanghui/system/zentao/manifests/latest repository=fanghui/system/zentao reference=latest 
DEBU[0332] registry.layer.check url=https://registry.alishui.com:443/v2/fanghui/system/zentao/blobs/sha256:8ba75e096fc7940a0199ac913575d4a421220d299d056ace5e9fa4ec89df9c8e repository=fanghui/system/zentao digest=sha256:8ba75e096fc7940a0199ac913575d4a421220d299d056ace5e9fa4ec89df9c8e 
DEBU[0332] registry.manifest.get url=https://registry.alishui.com:443/v2/fanghui/system/zentao/manifests/latest repository=fanghui/system/zentao reference=latest 
DEBU[0332] registry.manifest.get url=https://registry.alishui.com:443/v2/fanghui/system/zentao/manifests/9.5 repository=fanghui/system/zentao reference=9.5 
DEBU[0332] registry.layer.check url=https://registry.alishui.com:443/v2/fanghui/system/zentao/blobs/sha256:8ba75e096fc7940a0199ac913575d4a421220d299d056ace5e9fa4ec89df9c8e repository=fanghui/system/zentao digest=sha256:8ba75e096fc7940a0199ac913575d4a421220d299d056ace5e9fa4ec89df9c8e 
DEBU[0332] registry.manifest.get url=https://registry.alishui.com:443/v2/fanghui/system/zentao/manifests/9.5 repository=fanghui/system/zentao reference=9.5 
DEBU[0332] registry.tags url=https://registry.alishui.com:443/v2/fanghui/system/zookeeper/tags/list repository=fanghui/system/zookeeper 
DEBU[0332] registry.manifest.get url=https://registry.alishui.com:443/v2/fanghui/system/zookeeper/manifests/latest repository=fanghui/system/zookeeper reference=latest 
DEBU[0332] registry.layer.check url=https://registry.alishui.com:443/v2/fanghui/system/zookeeper/blobs/sha256:c1b59253ae0377453ecdd61d895c11e3ca8091d7f617b048f11447058a40c631 repository=fanghui/system/zookeeper digest=sha256:c1b59253ae0377453ecdd61d895c11e3ca8091d7f617b048f11447058a40c631 
DEBU[0332] registry.manifest.get url=https://registry.alishui.com:443/v2/fanghui/system/zookeeper/manifests/latest repository=fanghui/system/zookeeper reference=latest 
DEBU[0332] registry.manifest.get url=https://registry.alishui.com:443/v2/fanghui/system/zookeeper/manifests/3.5 repository=fanghui/system/zookeeper reference=3.5 
DEBU[0332] registry.layer.check url=https://registry.alishui.com:443/v2/fanghui/system/zookeeper/blobs/sha256:c1b59253ae0377453ecdd61d895c11e3ca8091d7f617b048f11447058a40c631 repository=fanghui/system/zookeeper digest=sha256:c1b59253ae0377453ecdd61d895c11e3ca8091d7f617b048f11447058a40c631 
DEBU[0333] registry.manifest.get url=https://registry.alishui.com:443/v2/fanghui/system/zookeeper/manifests/3.5 repository=fanghui/system/zookeeper reference=3.5 
DEBU[0333] registry.tags url=https://registry.alishui.com:443/v2/fanghui/test/lua/tags/list repository=fanghui/test/lua 
DEBU[0333] registry.manifest.get url=https://registry.alishui.com:443/v2/fanghui/test/lua/manifests/5.2 repository=fanghui/test/lua reference=5.2 
DEBU[0333] registry.layer.check url=https://registry.alishui.com:443/v2/fanghui/test/lua/blobs/sha256:256776419ea258b21ff69322b2fe4cdf66e22a40bc2ec2c29386f490eb3fc87c repository=fanghui/test/lua digest=sha256:256776419ea258b21ff69322b2fe4cdf66e22a40bc2ec2c29386f490eb3fc87c 
DEBU[0333] registry.manifest.get url=https://registry.alishui.com:443/v2/fanghui/test/lua/manifests/5.2 repository=fanghui/test/lua reference=5.2 
DEBU[0333] registry.manifest.get url=https://registry.alishui.com:443/v2/fanghui/test/lua/manifests/lts repository=fanghui/test/lua reference=lts 
DEBU[0333] registry.layer.check url=https://registry.alishui.com:443/v2/fanghui/test/lua/blobs/sha256:256776419ea258b21ff69322b2fe4cdf66e22a40bc2ec2c29386f490eb3fc87c repository=fanghui/test/lua digest=sha256:256776419ea258b21ff69322b2fe4cdf66e22a40bc2ec2c29386f490eb3fc87c 
DEBU[0333] registry.manifest.get url=https://registry.alishui.com:443/v2/fanghui/test/lua/manifests/lts repository=fanghui/test/lua reference=lts 
DEBU[0333] registry.tags url=https://registry.alishui.com:443/v2/haimi/system/php7-phalcon/tags/list repository=haimi/system/php7-phalcon 
DEBU[0333] registry.manifest.get url=https://registry.alishui.com:443/v2/haimi/system/php7-phalcon/manifests/latest repository=haimi/system/php7-phalcon reference=latest 
DEBU[0333] registry.layer.check url=https://registry.alishui.com:443/v2/haimi/system/php7-phalcon/blobs/sha256:9f159da739544acca04da6325ceaaffa89cb9948ff8073c1f7704e8c08f2031c repository=haimi/system/php7-phalcon digest=sha256:9f159da739544acca04da6325ceaaffa89cb9948ff8073c1f7704e8c08f2031c 
DEBU[0333] registry.manifest.get url=https://registry.alishui.com:443/v2/haimi/system/php7-phalcon/manifests/latest repository=haimi/system/php7-phalcon reference=latest 
DEBU[0333] registry.tags url=https://registry.alishui.com:443/v2/jenkins/jenkins/tags/list repository=jenkins/jenkins 
DEBU[0333] registry.manifest.get url=https://registry.alishui.com:443/v2/jenkins/jenkins/manifests/lts-alpine repository=jenkins/jenkins reference=lts-alpine 
DEBU[0333] registry.layer.check url=https://registry.alishui.com:443/v2/jenkins/jenkins/blobs/sha256:86c5dddf5c404e0838010644316f1dd44e0fd2e4b08e1aaaa690878033e5b12b repository=jenkins/jenkins digest=sha256:86c5dddf5c404e0838010644316f1dd44e0fd2e4b08e1aaaa690878033e5b12b 
DEBU[0333] registry.manifest.get url=https://registry.alishui.com:443/v2/jenkins/jenkins/manifests/lts-alpine repository=jenkins/jenkins reference=lts-alpine 
DEBU[0333] registry.tags url=https://registry.alishui.com:443/v2/kong/tags/list repository=kong 
DEBU[0333] registry.manifest.get url=https://registry.alishui.com:443/v2/kong/manifests/0.12.1-alpine repository=kong reference=0.12.1-alpine 
DEBU[0333] registry.layer.check url=https://registry.alishui.com:443/v2/kong/blobs/sha256:b54c874d32ceed59274451ee7e0cbb4a04e6d968eec92345a06eacf06aeaf34e repository=kong digest=sha256:b54c874d32ceed59274451ee7e0cbb4a04e6d968eec92345a06eacf06aeaf34e 
DEBU[0333] registry.manifest.get url=https://registry.alishui.com:443/v2/kong/manifests/0.12.1-alpine repository=kong reference=0.12.1-alpine 

```