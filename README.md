## lcx_go
lcx基于go的实现

基本是从[NATBypass](https://github.com/cw1997/NATBypass)重写过来的，改了main的部分逻辑，添加了域名解析，去除了日志输出(但大小只小了0.4m。。。C语言版的50k，yyds

不过go写这种工具属实方便很多～

用法与lcx一样：
```
usage: "-listen port1 port2"
       "-tran port1 ip:port2"
       "-slave ip1:port1 ip2:port2"
```