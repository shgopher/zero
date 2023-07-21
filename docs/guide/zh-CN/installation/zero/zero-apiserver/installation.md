# zero-apiserver 部署指南（容器化）

1. 创建 Workload

```bash
$ sed "s/127.0.0.1/$HOSTIP/g" deployments/zero/zero-apiserver/*|kubectl -n zero apply -f -
```

2. 测试是否部署成功

```bash
$ kubectl -s https://zero.apiserver.superproj.com:18443 --kubeconfig=$HOME/.zero/config get ms
```
