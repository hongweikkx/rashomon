# rashomon
Gateway

做一个网关系统， 这个系统拥有的功能点
* 协议代理。 可以作为grpc/http/tcp 的协议代理
* 可扩展。 也就是说可以通过增加机器来使得此系统不会成为瓶颈
* 负载均衡。 
* 服务发现。 

为了使得系统拥有高可用，高性能， 高并发， 安全。系统需满足
* 降级
* 熔断
* 隔离
* 限流
* 身份验证

beego vs gin: 选择gin。 因为不需要MVC架构。 GIN 可以满足需求。 

框架流程大致为：
负载均衡 -> 代理层 -> 身份验证 -> middleware -> 服务发现 


v0.1 
代理支持tcp
负载均衡



## how
build: make build
run: make run
