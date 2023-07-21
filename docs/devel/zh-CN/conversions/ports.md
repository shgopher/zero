## zero 端口列表

Linux端口使用规则：

- 0 不使用
- 1–1023 系统保留,只能由root用户使用
- 1024—4999 由客户端程序自由分配
- 5000—65535 由服务器端程序自由分配


zero端口分配规则：

wwxyz: 
- ww: 服务分类，从50开始
- x: 程序自己分配模块类别，默认为0
- y: 0 http/https, 1 grpc
- z: 0 非health类/health+业务,1 only health


zero-usercenter: 
  http/https/rpc: 50000/50010
zero-gateway: 
  http/https/rpc: 51000/51010
zero-apiserver: 
   http/https/rpc: 52000/52010
zero-osscenter: 
  http/https/rpc: 53000/53010
zero-nightwatch: 
  health: 54001
zero-pump: 
  health: 55001
zero-toyblc: 
  http: 58000
zero-tyk-server: 
  http/https/rpc: 56000/56010/56020
zero-event-monitor: 
  health: 57001
swagger server:
  http: 65534
