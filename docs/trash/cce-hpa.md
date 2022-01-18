# 自动扩容 (Horizontal Pod Autoscaler)

## 介绍

Pod 水平自动扩缩（Horizontal Pod Autoscaler） 可以基于 CPU 利用率自动扩缩 ReplicationController、Deployment、ReplicaSet 和 StatefulSet 中的 Pod 数量。 除了 CPU 利用率，也可以基于其他应程序提供的 自定义度量指标 来执行自动扩缩。 Pod 自动扩缩不适用于无法扩缩的对象，比如 DaemonSet。

Pod 水平自动扩缩器由`controller-manager`的 --horizontal-pod-autoscaler-sync-period 参数指定周期（默认值为 15 秒）。

## 示例

与其他 API 资源类似，`kubectl` 以标准方式支持 HPA。

我们可以通过 `kubectl create` 命令创建一个 HPA 对象， 通过 `kubectl get hpa` 命令来获取所有 HPA 对象， 通过 `kubectl describe hpa` 命令来查看 HPA 对象的详细信息。 最后，可以使用 `kubectl delete hpa` 命令删除对象。

此外，还有个简便的命令 `kubectl autoscale` 来创建 HPA 对象。 例如，命令 `kubectl autoscale rs foo --min=2 --max=5 --cpu-percent=80` 将会为名 为 `foo` 的 ReplicationSet 创建一个 HPA 对象， 目标 CPU 使用率为 80%，副本数量配置为 2 到 5 之间。

HPA对象示例：

```yaml
apiVersion: autoscaling/v2beta2
kind: HorizontalPodAutoscaler
metadata:
  name: php-apache
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: php-apache
  minReplicas: 1
  maxReplicas: 10
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 50
status:
  observedGeneration: 1
  lastScaleTime: <some-time>
  currentReplicas: 1
  desiredReplicas: 1
  currentMetrics:
  - type: Resource
    resource:
      name: cpu
      current:
        averageUtilization: 0
        averageValue: 0
```

`scaleTargetRef` 为HPA监视的工作负载对象

`metrics` 为具体的度量指标

### 对资源指标的支持

HPA 的任何目标资源都可以基于其中的 Pods 的资源用量来实现扩缩。 在定义 Pod 规约时，类似 cpu 和 memory 这类资源请求必须被设定。 这些设定值被用来确定资源利用量并被 HPA 控制器用来对目标资源完成扩缩操作。 要使用基于资源利用率的扩缩，可以像下面这样指定一个指标源：

```yaml
type: Resource
resource:
  name: cpu
  target:
    type: Utilization
    averageUtilization: 60
```

基于这一指标设定，HPA 控制器会维持扩缩目标中的 Pods 的平均资源利用率在 60%。

有多种可以选用的检测项，具体文档参考

https://kubernetes.io/zh/docs/tasks/run-application/horizontal-pod-autoscale/

https://kubernetes.io/zh/docs/tasks/run-application/horizontal-pod-autoscale-walkthrough/

## 在华为云CCE上使用HPA

### 创建工作负载弹性伸缩（HPA）

- 前提条件

使用HPA需要安装metrics-server插件，metrics-server负责采集kubernetes集群中kubelet的公开指标项，包含CPU利用率、内存利用率

### 操作步骤

1. 在CCE控制台中，单击左侧导航栏的“弹性伸缩”，在“工作负载伸缩”页签下，单击“创建HPA策略”。默认每 30 秒 轮询一次

2. 进入创建工作负载HPA策略页面，在“插件检测”步骤中：

- 若插件名称后方显示橙色的`未安装`，请单击插件后方的“现在安装”，根据业务需求配置插件参数后单击“立即安装”，等待插件安装完成。

- 若插件名称后方显示绿色的`已安装`，则说明插件已安装成功。

3. 确认插件已安装成功后，单击“下一步：策略配置”。

4. 在“策略配置”步骤中，策略参数配置。

实例范围： 输入最小实例数和最大实例数。策略触发时，工作负载实例将在此范围内伸缩。

冷却时间： 输入缩容和扩容的冷却时间，单位为分钟，缩容扩容冷却时间不能小于1分钟。

策略成功触发后，在此缩容/扩容冷却时间内，不会再次触发缩容/扩容，目的是等待伸缩动作完成后在系统稳定且集群正常的情况下进行下一次策略匹配。

5. 设置完成后，单击“创建”，在“完成”步骤中若显示“创建工作负载策略***提交成功”，可单击“返回工作负载伸缩策略”。

6. 在“工作负载伸缩”页签下，可以看到刚刚创建的HPA策略。
