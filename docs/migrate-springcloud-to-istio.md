# 微服务网格化方案

## 目的

### 兼容性和替代方案

nacos 保留 升级nacos支持service，使用k8s service做服务名 
1. 新版nacos支持k8s svc
2. nacos sdk支持 python/nodejs/cpp/c#/golang
3. 新版nacos可以集成istio
nacos 定时每10s 同步service list 给istiod.
实时感知服务上下线问题?
istiod服务发现同步给envoy服务列表时是全量同步的. 服务很多时有同步的压力. envoy直接对接nacos吗?
同步会检测服务的checksum. 服务信息如果有变更再同步. 但是好像有点bug…
————————————————
版权声明：本文为CSDN博主「Muroidea」的原创文章，遵循CC 4.0 BY-SA版权协议，转载请附上原文出处链接及本声明。
原文链接：https://blog.csdn.net/u014087707/article/details/120416146

sentinel适配
https://blog.csdn.net/lp19861126/article/details/105055272/

spring api gateway  vs  istio gateway as api gateway

Skywalking适配
1. skywalking istio bypass adapter

经过华为和 skywalking 核心开发者的确认，版本对应关系如下：
https://www.modb.pro/db/135357
istio
1.3     不支持生产 skywalking
使用
istio
1.7以上  skywalking
链路拓扑可以商用
istio
1.8     skywalking
日志商用
istio
1.11    trace
商用

一、确定 华为云istio的版本 看起来是自研基于envoy的，跟istio版本关系不打


我们归纳了服务网格支撑企业落地需要具备的 “三要素” ：通信协议，注册中心，部署环境。

通信协议：服务网格能支持的服务通信协议，常见的如 HTTP、gRPC、Dubbo 等，另外也有具备行业属性的私有 RPC 协议；

注册中心：服务网格能纳管的注册中心，包括常见的 Eureka、Consul、Nacos、Zookeeper 以及 Kubernetes （ETCD）；

部署环境：服务网格能支持的业务部署环境，除了天然云原生的 Kubernetes + Docker 外，对于遗留系统所在的虚拟机、物理机，也需要同等对待。

在满足 “三要素” 后，服务网格才能达到业务落地的 “及格线”。




在初步完成服务网格认知后，企业用户往往会发出灵魂拷问：为什么要上服务网格？服务网格有什么价值？

一般来说，通识的服务网格核心价值 “标准答案” 是：

业务无需感知微服务组件：微服务架构支撑、网络通信、治理等相关能力下沉到基础设施层，业务部门无需投入专人开发与维护，可以有效降低微服务架构下研发与维护成本；

支持多开发语言、框架：服务网格天然不限制开发语言、开发框架，提供多语言服务治理能力；

框架升级零成本：支持框架热升级，降低中间件和技术框架客户端、SDK 升级成本；

微服务体系统一纳管、演进：将存量微服务集群、遗留系统、外购系统微服务体系统一管理、演进。


作为 Istio 服务网格中的一部分，Kubernetes 集群中的 Pod 和 Service 必须满足以下要求：

命名的服务端口: Service 的端口必须命名。端口名键值对必须按以下格式：name: <protocol>[-<suffix>]。更多说明请参看协议选择。

Service 关联: 每个 Pod 必须至少属于一个 Kubernetes Service，不管这个 Pod 是否对外暴露端口。如果一个 Pod 同时属于多个 Kubernetes Service， 那么这些 Service 不能同时在一个端口号上使用不同的协议（比如：HTTP 和 TCP）。

带有 app 和 version 标签（label）的 Deployment: 我们建议显式地给 Deployment 加上 app 和 version 标签。给使用 Kubernetes Deployment 部署的 Pod 部署配置中增加这些标签，可以给 Istio 收集的指标和遥测信息中增加上下文信息。

app 标签：每个部署配置应该有一个不同的 app 标签并且该标签的值应该有一定意义。app label 用于在分布式追踪中添加上下文信息。

version 标签：这个标签用于在特定方式部署的应用中表示版本。

应用 UID: 确保你的 Pod 不会以用户 ID（UID）为 1337 的用户运行应用。

NET_ADMIN 功能: 如果你的集群执行 Pod 安全策略，必须给 Pod 配置 NET_ADMIN 功能。如果你使用 Istio CNI 插件 可以不配置。要了解更多 NET_ADMIN 功能的知识，请查看所需的 Pod 功能。
https://istio.io/latest/zh/docs/ops/deployment/requirements/

