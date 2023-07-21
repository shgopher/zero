# zero-toyblc 部署指南（容器化）

1. 创建 ConfigMap

```bash
$ cp configs/zero-toyblc.yaml $HOME/.zero/zero-toyblc.yaml
$ kubectl -n zero create configmap zero-toyblc --from-file $HOME/.zero/zero-toyblc.yaml
```

> 注意：创建前，记得修改 `zero-toyblc.yaml` 中相关配置，例如：访问地址、密码等。

2. 创建 Workload

```bash
$ kubectl -n zero apply -f deployments/zero/zero-toyblc
```

3. 测试是否部署成功

```bash
$ curl -H "Host: zero.toyblc.superproj.com" http://127.0.0.1:18080/healthz
```

## 使用

### 1. 查询 peers

```bash
$ curl -H "Host: zero.toyblc.superproj.com" http://127.0.0.1:18080/v1/peers
```

### 2. 查询 blocks

```bash
$ curl -H "Host: zero.toyblc.superproj.com" http://127.0.0.1:18080/v1/blocks
```

> curl http://genesis.kube-system.svc.zero.io:8080/v1/blocks

### 3. 挖矿 

```bash
$ curl -XPOST -H "Host: zero.toyblc.superproj.com" -H"Content-type: application/json" -d'{"data": "Some data to the first block"}' http://127.0.0.1:18080/v1/blocks
```
