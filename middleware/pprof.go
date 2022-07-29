package middleware

// third part provide service, this just give some pprof usage

/**
提供服务的第三方包：https://github.com/gin-contrib/pprof


Pprof 简单使用：

1、采样pprof
go tool pprof http://127.0.0.1:9000/debug/pprof/profile -seconds 30

2、查看采样效果，火焰图+逻辑调用图
go tool pprof -http=:8080 "~\pprof\pprof.samples.cpu.004.pb.gz"
*/
