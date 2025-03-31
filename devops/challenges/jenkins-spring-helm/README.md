
### Install kind

```shell
# For AMD64 / x86_64
[ $(uname -m) = x86_64 ] && curl -Lo ./kind https://kind.sigs.k8s.io/dl/v0.27.0/kind-linux-amd64
# For ARM64
[ $(uname -m) = aarch64 ] && curl -Lo ./kind https://kind.sigs.k8s.io/dl/v0.27.0/kind-linux-arm64
chmod +x ./kind
sudo mv ./kind /usr/local/bin/kind
```

### Create cluster
```shell
kind create cluster
```

### Install helm
```shell
curl https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash
```


helm repo add harbor https://helm.goharbor.io
helm install --set expose.type=nodePort --set expose.tls.enabled=false harbor harbor/harbor

kubectl port-forward service/harbor 30005:80
