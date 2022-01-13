#!/usr/bin/env bash

if [ ! -d "./istio-1.12.1" ]; then
    echo "istio-1.12.1 folder not exsit, downloading..."
    curl -L https://istio.io/downloadIstio | ISTIO_VERSION=1.12.1 TARGET_ARCH=x86_64 sh -
fi

cd istio-1.12.1
export PATH=$PWD/bin:$PATH

istioctl install --set profile=demo -y
kubectl label namespace default istio-injection=enabled
kubectl apply -f samples/bookinfo/platform/kube/bookinfo.yaml
kubectl apply -f samples/bookinfo/networking/bookinfo-gateway.yaml
kubectl apply -f samples/addons
kubectl rollout status deployment/kiali -n istio-system
istioctl analyze

kubectl exec "$(kubectl get pod -l app=ratings -o jsonpath='{.items[0].metadata.name}')" -c ratings -- curl -s productpage:9080/productpage | grep -o "<title>.*</title>"
