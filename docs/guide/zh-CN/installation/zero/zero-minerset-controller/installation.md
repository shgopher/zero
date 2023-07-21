# zero-minerset-controller 部署指南（容器化）

1. 创建 ConfigMap

```bash
$ cp configs/zero-minerset-controller.yaml $HOME/.zero/zero-minerset-controller.yaml
$ kubectl -n zero create configmap zero-minerset-controller --from-file $HOME/.zero/zero-minerset-controller.yaml
```

> 注意：创建前，记得修改 `zero-minerset-controller.yaml` 中相关配置，例如：访问地址、密码等。

2. 创建 Workload

```bash
$ kubectl -n zero apply -f deployments/zero/zero-minerset-controller
```

3. 测试是否部署成功

```bash
$ curl -H "Host: zero.minerset.superproj.com" http://127.0.0.1:18080/healthz
```
