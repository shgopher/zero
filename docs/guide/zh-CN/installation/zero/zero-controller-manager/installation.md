# zero-controller-manager 部署指南（容器化）

1. 创建 ConfigMap

```bash
$ sed "s/127.0.0.1/$HOSTIP/g" configs/zero-controller-manager.yaml > $HOME/.zero/zero-controller-manager.yaml
$ kubectl -n zero create configmap zero-controller-manager --from-file $HOME/.zero/zero-controller-manager.yaml
```

> 注意：创建前，记得修改 `zero-controller-manager.yaml` 中相关配置，例如：访问地址、密码等。

2. 创建 Workload

```bash
$ kubectl -n zero apply -f deployments/zero/zero-controller-manager
```

3. 测试是否部署成功

```bash
$ curl -H "Host: zero.controllermanager.superproj.com" http://127.0.0.1:18080/healthz
```
