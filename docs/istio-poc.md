# Isito

Istio Gateway 资源本身只能配置L4到L6的功能，例如暴露的端口、TLS 设置等；但 Gateway 可与 VirtualService 绑定，在VirtualService 中可以配置七层路由规则，例如按比例和版本的流量路由，故障注入，HTTP 重定向，HTTP 重写等所有Mesh内部支持的路由规则。

