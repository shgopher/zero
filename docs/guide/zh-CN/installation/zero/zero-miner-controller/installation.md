# zero-miner-controller 部署指南（容器化）

1. 创建 ConfigMap

```bash
$ sed "s/127.0.0.1/$HOSTIP/g" configs/zero-miner-controller.yaml > $HOME/.zero/zero-miner-controller.yaml
$ kubectl -n zero create configmap zero-miner-controller --from-file $HOME/.zero/zero-miner-controller.yaml --from-file config.kind=$HOME/.kube/config
```

> 注意：创建前，记得修改 `zero-miner-controller.yaml` 中相关配置，例如：访问地址、密码等。

2. 创建 Workload

```bash
$ kubectl -n zero apply -f deployments/zero/zero-miner-controller
```

3. 测试是否部署成功

```bash
$ curl -H "Host: zero.miner.superproj.com" http://127.0.0.1:18080/healthz
```
