#!/bin/sh
chmod 777 /etc/kubernetes/manifests/kube-apiserver.yaml
sleep 1s
chown vagrant:vagrant /etc/kubernetes/manifests/kube-apiserver.yaml
sleep 1s
chmod 777 /etc/kubernetes/manifests/kube-scheduler.yaml
sleep 1s
chown root:root /etc/kubernetes/manifests/kube-scheduler.yaml
sleep 1s
chmod 777 /etc/kubernetes/manifests/etcd.yaml
sleep 1s
chown vagrant:vagrant /etc/kubernetes/manifests/etcd.yaml 
sleep 1s
curl http://127.0.0.1:8080/api/v1/nodes
sleep 1s
curl http://127.0.0.1:8080/api/v1/nodes
sleep 1s
curl http://127.0.0.1:8080/api/v1/nodes
sleep 1s
curl http://127.0.0.1:8080/api/v1/nodes
sleep 1s
curl http://127.0.0.1:8080/api/v1/namespaces/default/resourcequotas

