但 Spring Cloud 也有一些不可避免的缺点，如基于不同框架的不同组件带来的高应用门槛及学习成本、代码级别对诸多组件进行控制的需求与微服务多语言协作的目标背道而驰。

在我们内部，由于历史原因，不同小组所使用的 API 网关架构不统一，且存在多套 Spring Cloud，给统一管理造成了不便；Spring Cloud 无法实现灰度发布，也给公司业务发布带来了一定不便。更重要的是，作为一家周边游网站，我们经常会举行一些促销活动，面临在业务峰值期资源弹性扩缩容的需求，仅仅依靠 Spring Cloud 也无法实现资源调度来满足业务自动扩缩容的需求。



从 Spring Cloud 到 UK8S 的过程，也是内部服务模块再次梳理、统一的过程，在此过程中，我们对整体业务架构做了如下改动：

1.去掉原有的 Eureka，改用 Spring Cloud Kubernetes 项目下的 Discovery。Spring Cloud 官方推出的项目 Spring Cloud Kubernetes 提供了通用的接口来调用Kubernetes服务，让 Spring Cloud 和 Spring Boot 程序能够在 Kubernetes 环境中更好运行。在 Kubernetes 环境中，ETCD 已经拥有了服务发现所必要的信息，没有必要再使用 Eureka，通过 Discovery 就能够获取 Kubernetes ETCD 中注册的服务列表进行服务发现。

2.去掉 Feign 负载均衡，改用 Spring Cloud Kubernetes Ribbon。Ribbon 负载均衡模式有 Service / Pod 两种，在 Service 模式下，可以使用 Kubernetes 原生负载均衡，并通过 Istio 实现服务治理。

3. 网关边缘化。网关作为原来的入口，全部去除需要对原有代码进行大规模的改造，我们把原有的网关作为微服务部署在 Kubernetes 内，并利用 Istio 来管理流量入口。同时，我们还去掉了熔断器和智能路由，整体基于 Istio 实现服务治理。

4. 分布式配置 Config 统一为 Apollo。Apollo 能够集中管理应用在不同环境、不同集群的配置，修改后实时推送到应用端，并且具备规范的权限、流程治理等特性。

5. 增加 Prometheus 监控。特别是对 JVM 一些参数和一些定义指标的监控，并基于监控指标实现了 HPA 弹性伸缩。


在 Kubernetes中，HPA 通常通过 Pod 的 CPU、内存利用率等实现，但在 Java 中，内存控制通过 JVM 实现，当内存占用过高时，JVM 会进行内存回收，但 JVM 并不会返回给主机或容器，单纯基于 Pod / CPU 指标进行集群的扩缩容并不合理。我们通过 Prometheus 获取 Java 中 http_server_requests_seconds_count（请求数）参数，通过适配器将其转化成 Kubernetes API Server 能识别的参数，并基于这一指标实时动态调整 Pod 的数量。

Istio服务治理

基于应用程序安全性、可观察性、持续部署、弹性伸缩和性能、对开源工具的集成、开源控制平面的支持、方案成熟度等考虑，我们最终选择了 Istio 作为服务治理的方案，主要涉及以下几个部分：1. Istio-gateway 网关：Ingress Gateway 在逻辑上相当于网格边缘的一个负载均衡器，用于接收和处理网格边缘出站和入站的网络连接，其中包含开放端口和TLS的配置等内容，实现集群内部南北流量的治理。

2. Mesh 网关：Istio内部的虚拟Gateway，代表网格内部的所有Sidecar，实现所有网格内部服务之间的互相通信，即东西流量的治理。

3. 流量管理：在去除掉 Spring Cloud 原有的熔断、智能路由等组件后，我们通过对 Kubernetes 集群内部一系列的配置和管理，实现了 http 流量管理的功能。包括使用 Pod签对具体的服务进程进行分组（例如 V1/V2 版本应用）并实现流量调度，通过 Istio 内的 Destination Rule 单独定义服务负载均衡策略，根据来源服务、URL 进行重定向实现目标路由分流等，通过 MenQuota、RedisQuota 进行限流等。

4. 遥测：通过 Prometheus 获取遥测数据，实现灰度项目成功率、东西南北流量区分、服务峰值流量、服务动态拓扑的监控。

分类	SpringCloud	OpenShift ServiceMesh
服务配置	Spring Cloud Config Server	ConfigMap, Secret
服务注册与发现	Eureka	Etcd + Service + 集群内DNS
负载均衡	Ribbon	Service, Istio的Envoy数据平面
服务间调用	OpenFeign 或 RestTemplate	任意HTTP client
路由管理	Zuul 或 Spring Cloud Gateway	Istio的VirtualService和DetinationRule
对外API网关	Zuul 或 Spring Cloud Gateway	Route，Istio的Ingress gateway和Egress gateway
限流和熔断	Hystrix	Istio的Envoy数据平面
服务跟踪和调用链	Zipkin 或 OpenTracing + Jaeger	OpenTracing + Jager
服务网络拓扑	无	Kiali
灰度发布和蓝绿部署	无	Istio的Envoy数据平面


北京时间2月23日，在全球首届社区峰会IstioCon 2021中，华为云应用服务网格首席架构师张超盟发表了《Best practice:from Spring Cloud to Istio》主题演讲，分享了Istio在生产中使用的实际案例。

议题简介
https://zhuanlan.zhihu.com/p/358891699

微服务SDK曾经是一个常用的解决方案。
SDK形态中Spring cloud是最有影响力的代表项目。Spring cloud提供了构建分布式应用的开发工具集，如列表所示。其中被大部分开发者熟知的是微服务相关项目，如：服务注册发现eureka、配置管理 config、负载均衡ribbon、熔断容错Hystrix、调用链埋点sleuth、网关zuul或Spring cloud gateway等项目。在本次分享中提到的Spring cloud也特指Spring cloud的微服务开发套件。

以下是我们客户找到我们TOP的几个的问题，剖析下用户使用传统微服务框架碰到了哪些问题，这些大部分也是他们选择网格的最大动力。
1）多语言问题

2）将Spring cloud的微服务运行在K8s上会有很大的概率出现服务发现不及时

3）升级所有应用以应对服务管理需求变化

4）从单体式架构向微服务架构迁移



https://www.yisu.com/zixun/505602.html


https://www.jianshu.com/p/ee82c9f0965c



华为云SpringCloud微服务至Istio迁移指导


https://bbs.huaweicloud.com/forum/forum.php?mod=viewthread&tid=126766
https://support.huaweicloud.com/bestpractice-cce/istio_bestpractice_3012.html


你是否真的需要 Istio？
你可能参加过各种云原生、服务网格相关的 meetup，在社区里看到很多人在分享和讨论 Istio，但是对于自己是否真的需要 Istio 感到踌躇，甚至因为它的复杂性而对服务网格的前景感到怀疑。那么，在你继阅读 Istio SIG 后续文章之前，请先仔细阅读本文，审视一下自己公司的现状，看看你是否有必要使用服务网格，处于 Istio 应用的哪个阶段。

本文不是对应用服务网格的指导，而是根据社区里经常遇到的问题而整理。在使用 Istio 之前，请先考虑下以下因素：

你的团队里有多少人？
你的团队是否有使用 Kubernetes、Istio 的经验？
你有多少微服务？
这些微服务使用什么语言？
你的运维、SRE 团队是否可以支持服务网格管理？
你有采用开源项目的经验吗？
你的服务都运行在哪些平台上？
你的应用已经容器化并使用 Kubernetes 管理了吗？
你的服务有多少是部署在虚拟机、有多少是部署到 Kubernetes 集群上，比例如何？
你的团队有制定转移到云原生架构的计划吗？
你想使用 Istio 的什么功能？
Istio 的稳定性是否能够满足你的需求？
你是否可以忍受 Istio 带来的性能损耗？
请先思考一下上述问题，关于是否应该使用 Istio，及应用服务网格化的路径，欢迎到云原生社区 Istio SIG 中探讨。



来简单看一下他们的功能对比：

功能列表

Spring Cloud

Istio

服务注册与发现

支持，基于Eureka，consul等组件，提供server，和Client管理

支持，基于XDS接口获取服务信息，并依赖“虚拟服务路由表”实现服务发现

链路监控

支持，基于Zikpin或者Pinpoint或者Skywalking实现

支持，基于sideCar代理模型，记录网络请求信息实现

API网关

支持，基于zuul或者spring-cloud-gateway实现

支持，基于Ingress gateway以及egress实现

熔断器

支持，基于Hystrix实现

支持，基于声明配置文件，最终转化成路由规则实现

服务路由

支持，基于网关层实现路由转发

支持，基于iptables规则实现

安全策略

支持，基于spring-security组件实现，包括认证，鉴权等，支持通信加密

支持，基于RBAC的权限模型，依赖Kubernetes实现，同时支持通信加密

配置中心

支持，springcloud-config组件实现

不支持

性能监控

支持，基于Spring cloud提供的监控组件收集数据，对接第三方的监控数据存储

支持，基于SideCar代理，记录服务调用性能数据，并通过metrics adapter，导入第三方数据监控工具

日志收集

支持，提供client，对接第三方日志系统，例如ELK

支持，基于SideCar代理，记录日志信息，并通过log adapter，导入第三方日志系统

工具客户端集成

支持，提供消息，总线，部署管道，数据处理等多种工具客户端SDK

不支持

分布式事务

支持，支持不同的分布式事务模式：JTA，TCC，SAGA等，并且提供实现的SDK框架

不支持

其他

……

……

SpringCloud微服务Isito迁移指导
https://support.huaweicloud.com/bestpractice-cce/istio_bestpractice_3012.html
 

从上面表格中可以看到，如果从功能层面考虑，Spring Cloud与Service Mesh在服务治理场景下，有相当大量的重叠功能，从这个层面而言，为Spring Cloud向Service Mesh迁移提供了一种潜在的可能性。


https://developer.jdcloud.com/article/816

http://ideajava.com/2019/11/27/SpringCloud-to-K8S/


阿里巴巴 Service Mesh 落地的架构与挑战
https://zhuanlan.zhihu.com/p/97830665


Spring Cloud迁移Service Mesh（Istio）
https://www.jianshu.com/p/ee82c9f0965c
Spring Cloud向Service Mesh迁移，从我们考虑而言大体分为七个步骤，如图所示：
1 Spring Cloud架构解析

Spring Cloud架构解析的目的在于确定需要从当前的服务中去除与Service Mesh重叠的功能，为后续服务替换做准备。我们来看一个典型的Spring Cloud架构体系，如图所示：

2 服务改造 3 服务容器化 4 容器环境构建
5 ServiceMesh环境构建
6 服务侏儒
7 服务管理控制台

这里面哪些内容是我们可以拿掉或者说基于Service Mesh（以Istio为例）能力去做的？分析下来，可以替换的组件包括网关（gateway或者Zuul，由Ingress gateway或者egress替换），熔断器（hystrix，由SideCar替换），注册中心（Eureka及Eureka client，由Polit，SideCar替换），负责均衡（Ribbon，由SideCar替换），链路跟踪及其客户端（Pinpoint及Pinpoint client，由SideCar及Mixer替换）。这是我们在Spring Cloud解析中需要完成的目标：即确定需要删除或者替换的支撑模块。


服务迁移之路|Spring Cloud向Service Mesh转变
https://developer.jdcloud.com/article/816



https://zhuanlan.zhihu.com/p/358891699
我们的主要是思路是解耦和卸载。卸载原有SDK中非开发的功能，SDK只提供代码框架、应用协议等开发功能。涉及到微服务治理的内容都卸载到基础设施去做。



从图上可以看到开发人员接触到开发框架变薄了，开发人员的学习、使用和维护成本也相应的降低了。而基础设施变得厚重了，除了完成之前需要做的服务运行的基础能力外，还包括非侵入的服务治理能力。即将越来越多的之前认为是业务的能力提炼成通用能力，交给基础设施去做，交给云厂商去做，客户摆脱这些繁琐的非业务的事务，更多的时间和精力投入到业务的创新和开发上。在这种分工下，SDK才真的回归到开发框架的根本职能。


要使用网格的能力，前提是微服务出来的流量能走到网格的数据面来。主要的迁移工作在微服务的服务调用方。我们推荐3个步骤：

第一步：废弃原有的微服务注册中心，使用K8S的Service。

第二步：短路SDK中服务发现和负载均衡等逻辑，直接使用k8s的服务名和端口访问目标服务；

第三步：结合自身项目需要，逐步使用网格中的治理能力替换原有SDK中提供的对应功能，当然这步是可选的，如调用链埋点等原有功能使用满足要求，也可以作为应用自身功能保留。


为了达成以上迁移，我们有两种方式，供不同的用户场景采用。



一种是只修改配置的方式：Spring cloud本身除了支持基于Eureka的服务端的服务发现外，还可以给Ribbon配置静态服务实例地址。利用这种机制给对应微服务的后端实例地址中配置服务的Kubernetes服务名和端口。



当Spring cloud框架中还是访问原有的服务端微服务名时，会将请求转发到k8s的服务和端口上。这样访问会被网格数据面拦截，从而走到网格数据面上来。服务发现负载均衡和各种流量规则都可以应用网格的能力。



这种方式其实是用最小的修改将SDK的访问链路和网格数据面的访问链路串接起来。在平台中使用时，可以借助流水线工具辅助，减少直接修改配置文件的工作量和操作错误。可以看到我这个实际项目中，只是修改了项目的application.yaml配置文件，其他代码都是0修改。当然对于基于annotation的方式的配置也是同样的语义修改即可。


前面一种方式对原有项目的修改比较少，但是Spring cloud的项目依赖都还在。



我们有些客户选择了另外一种更简单直接的方式，既然原有SDK中服务发现负载均衡包括各种服务治理能力都不需要了，干脆这些依赖也全部干掉。从最终的镜像大小看，整个项目的体量也得到了极大的瘦身。这种方式客户根据自己的实际使用方式，进行各种裁剪，最终大部分是把Spring cloud退化成Spring boot。



迁移中还有另外一部分比较特殊，就是微服务外部访问的Gateway。



Spring cloud 有两种功能类似的Gateway，Zuul和Spring cloud Gateway。基于Eureka的服务发现，将内部微服务映射成外部服务，并且在入口处提供安全、分流等能力。在切换到k8s和Istio上来时，和内部服务一样，将入口各个服务的服务发现迁移到k8s上来。



差别在于对于用户如果在Gateway上开发了很多私有的业务强相关的filter时，这时候Gateway其实是微服务的门面服务，为了业务延续性，方案上可以直接将其当成普通的微服务部署在网格中进行管理。

https://zhuanlan.zhihu.com/p/358891699

应用服务网格（Application Service Mesh，简称ASM）是华为云基于开源Istio推出的服务网格平台，它深度、无缝对接了华为云的企业级Kubernetes集群服务CCE，可为客户提供开箱即用的上手体验。

计费项
ASM计费单位为实例数（CCE集群中Pod数量）。
ASM套餐价格中不包含APM服务费用，推荐购买APM服务享受全方位的服务性能监控：APM计费模式介绍。
ASM套餐价格不包含用户使用华为云上的资源费用（弹性云服务器、CCE集群管理费、ELB费用等），相关链接如下：
计费模式
ASM分为“按需计费”与“包年/包月”两种计费模式。

按需计费模式

按需计费模式是根据当前集群中服务网格治理的pod实例数量，按每小时扣费。

托管网格治理不足20实例时，按20实例数收取费用，超过20实例时，按实际实例数收取费用。
托管网格仅可选择5000实例治理。
专有网格可选择5000实例治理，也可选择体验20实例。
包周期模式

包年/包月计费模式：为您提供实例套餐包，相比于按需模式更优惠；包年更优惠，仅需支付10个月的套餐包费用即可包一整年。
表1 包周期收费表
网格可管理的最大实例数（Pod个数）

配置费用（元/月）

20

2,700

50

6,750

100

10,800

200

17,280

500

27,600

1,000

45,000

2,000

70,200
华为asm专业版 社区版对比
https://support.huaweicloud.com/productdesc-asm/asm_productdesc_0018.html

规格推荐列表
https://support.huaweicloud.com/productdesc-asm/asm_productdesc_0006.html

1）微服务和容器都有轻量和敏捷的共同特点，容器是微服务非常适合的一个运行环境；

2）在云原生场景下，在微服务场景下，容器从来都不是独立存在的，使用k8s来编排容器已经是一个事实标准；

3）Istio和k8s在架构和应用场景上的紧密结合，一起提供了一个端到端的微服务运行和治理的平台。

4）也是我们推荐的方案，使用Istio进行微服务治理正在成为越来越多用户的技术选择。



以上四个关系顺时针结合在一起为我们的解决方案构造一个完整的闭环。


## istio版本支持状态

要求的k8s版本范围窄
每个版本6个月支持期
1.10 支持 k8s 1.16-。122

Kubernetes 1.22 删除了一些已弃用的 API，因此 1.10.0 之前的 Istio 版本将不再工作。