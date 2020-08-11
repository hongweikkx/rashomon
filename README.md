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

todo:
有点奇怪的：
1. 如果我用 hystrix.Go(A, B) 就会报错 "nic: http: wrote more than the declared Content-Length"
但是如果是我用 hystrix.Go(c.next, B) 就不会报错。。。
2. jwt 虽然能用了  但是有点奇怪在jwt的认证方式 和 jwt-go的参数。 
3. i need todo 知道实际使用和原理之间的区别 比如还有hystrix 感觉很多都是
用不同的参数来实现的。 究竟怎样的才是好的呢？
