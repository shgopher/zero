# 项目笔记

> 注意：所有操作均在 `${ZROOT}` 目录下进行。

本项目所有密码均为：`zero(#)666`。

## 组件名、数据库、密码

为了能够追踪请求源，方便以后排障，这里针对每一个组件申请一个账号：

| 组件名 | 用户名 | 密码 | 备注 |
| ----------- | ----------- | ----------- | ----------- |
| - | zero | zero(#)666 | root account |
| zero-usercenter | usercenter | zero(#)666 | |
| zero-gateway | gateway | zero(#)666 | |
| zero-apiserver | apiserver | zero(#)666 | |
| zero-nightwatch | nightwatch | zero(#)666 | |

## 启动 GitBook

```bash
$ cd docs
$ gitbook serve # 访问 http://127.0.0.1:4000/
```

## 组件启动命令

1. 命令行直接启动

```bash
_output/platforms/linux/amd64/zero-usercenter --db.host=127.0.0.1 --db.username=zero --db.password='zero(#)666' --db.database=zero --redis.addr=127.0.0.1:6379 --redis.password='zero(#)666' --etcd.endpoints=127.0.0.1:2379 --kafka.brokers=localhost:9092 --http.addr=0.0.0.0:50000 --grpc.addr=0.0.0.0:50010

_output/platforms/linux/amd64/zero-usercenter --db.host=127.0.0.1 --db.username=zero --db.password='zero(#)666' --db.database=zero --redis.addr=127.0.0.1:6379 --redis.password='zero(#)666' --etcd.endpoints=127.0.0.1:2379 --kafka.brokers=localhost:9092 --http.addr=0.0.0.0:50000 --grpc.addr=0.0.0.0:50010 --tls.use-tls=true --tls.cert=/home/colin/.zero/cert/zero-usercenter.pem --tls.key=/home/colin/.zero/cert/zero-usercenter-key.pem

_output/platforms/linux/amd64/zero-gateway --db.host=127.0.0.1 --db.username=zero --db.password='zero(#)666' --db.database=zero --etcd.endpoints=127.0.0.1:2379 --insecure.addr=0.0.0.0:51000 --secure.bind-address=0.0.0.0 --secure.bind-port=51010 --grpc.addr=0.0.0.0:51020 --kubeconfig=$HOME/.zero/config --usercenter.server=127.0.0.1:50000

_output/platforms/linux/amd64/zero-apiserver --etcd-servers=127.0.0.1:2379 --secure-port=52010 --bind-address=0.0.0.0 --client-ca-file=/home/colin/.zero/cert/ca.pem --tls-cert-file=/home/colin/.zero/cert/zero-apiserver.pem --tls-private-key-file=/home/colin/.zero/cert/zero-apiserver-key.pem

 _output/platforms/linux/amd64/zero-controller-manager --kubeconfig /home/colin/.zero/config --mysql-database=zero --mysql-host=127.0.0.1:3306 --mysql-username=zero --mysql-password='zero(#)666'

_output/platforms/linux/amd64/zero-controller-manager --kubeconfig ~/.zero/config --config configs/zero-controller-manager.yaml

_output/platforms/linux/amd64/zero-nightwatch --kubeconfig /home/colin/.zero/config --db.host=127.0.0.1 --db.username=zero --db.password='zero(#)666' --db.database=zero --redis.addr=127.0.0.1:6379 --redis.password='zero(#)666' --redis.database=1
_output/platforms/linux/amd64/zero-nightwatch --config ~/.zero/zero-nightwatch.yaml

```

2. Kubernetes 部署


## 其他

1. 转发

```bash
kubectl port-forward -n kube-system --address 0.0.0.0 $(kubectl get pods -n kube-system --selector "app.kubernetes.io/name=traefik" --output=name) 8000:9000
kubectl port-forward -n zero --address 0.0.0.0 services/zero-apiserver 52010:https
```

机器学习资料：
https://github.com/microsoft/AI-System
https://openmlsys.github.io/index.html
