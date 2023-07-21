# zero-nightwatch 部署指南（容器化）

1. 创建 ConfigMap

```bash
$ sed "s/127.0.0.1/$HOSTIP/g" configs/zero-nightwatch.yaml > $HOME/.zero/zero-nightwatch.yaml
$ kubectl -n zero create configmap zero-nightwatch --from-file $HOME/.zero/zero-nightwatch.yaml
```

> 注意：创建前，记得修改 `zero-nightwatch.yaml` 中相关配置，例如：访问地址、密码等。

2. 创建 Workload

```bash
$ kubectl -n zero apply -f deployments/zero/zero-nightwatch
```

3. 测试是否部署成功

```bash
$ curl -H "Host: zero.nightwatch.superproj.com" http://127.0.0.1:18080/healthz
```
