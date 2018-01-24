# Build and PPROF

[profiling-go-programs](https://blog.golang.org/profiling-go-programs)

```
go get -v github.com/zhujingfa/docker-registry-manager/app

/Users/bruce/project/godev/bin/app -c /Users/bruce/project/godev/src/github.com/zhujingfa/docker-registry-manager/docs/test.yml -p cpu-test.log

go tool pprof /Users/bruce/project/godev/bin/app cpu-test.log 

pprof cmd: 

  - top 20 
  - top 20 -cum
  - web
  
web生成的svg图片，特别有用，一下子cpu占用就很清楚了。


嵌入代码，检测cpu：

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func main() {
  flag.Parse()
  if *cpuprofile != "" {
      f, err := os.Create(*cpuprofile)
      if err != nil {
          log.Fatal(err)
      }
      pprof.StartCPUProfile(f)
      defer pprof.StopCPUProfile()
  }
  ...
    
检测内存：

var memprofile = flag.String("memprofile", "", "write memory profile to this file")
...

FindHavlakLoops(cfgraph, lsgraph)
if *memprofile != "" {
    f, err := os.Create(*memprofile)
    if err != nil {
        log.Fatal(err)
    }
    pprof.WriteHeapProfile(f)
    f.Close()
    return
}    
    

  
```

