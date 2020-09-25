## Como instalar (v0.1)

O exporter precisa de acesso à API do kubernetes para puxar os IPS
dos pods do mtail-controller. Para isso é preciso criar um `clusterrolebinding`
no clusterrole view rodando o seguinte comando.

```bash
kubectl create clusterrolebinding default-view --clusterrole=view --serviceaccount=default:default
```

Assim o exporter terá visão e a partir das variáveis de ambiente passadas (`MTAIL_EXPORTER_POD_NAME` e
`MTAIL_EXPORTER_NAMESPACE`) irá reconhecer os IPS automaticamente a cada acesso ao `/metrics`.