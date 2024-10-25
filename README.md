# Kubelet-static-credential-provider
---

This plugin implements the [kubernetes credential-provider-api](https://kubernetes.io/docs/reference/config-api/kubelet-credentialprovider.v1/). This allows users to pull images from private image registries without having to create and reference a `image-pull-secret` in every namespace and deployment.

AWS, Azure and Google all use their own version of this plugin protocol to allow user friendly image pulls from their hosted image registries. Now you can also use this developer friendly approach to pull images from your private registries or password protected pull-through caches of Dockerhub.

## Quickstart
Download with `wget https://<RELEASE_URL> -O /var/lib/kubelet/plugins/static-credential-provider` and create a Config:
```yaml
# /srv/scp-con.yaml
apiVersion: kubelet.config.k8s.io/v1
kind: CredentialProviderConfig
providers:
  - name: static-credential-provider
    # You can also use a config file instead of envs
    # args:[--config, <path-to-password-conf>]
    matchImages: [docker.io]
    defaultCacheDuration: "12h"
    apiVersion: credentialprovider.kubelet.k8s.io/v1
    env:
      - name: SCP_REGISTRY_USERNAME
        value: <my-user>
      - name: SCP_REGISTRY_PASSWORD
        value: <my-password>
```

Add `--image-credential-provider-bin-dir=/var/lib/kubelet/plugins/` and `--image-credential-provider-config=/srv/scp-con.yaml` args to the kubelet startup args and restart the kubelet with `service kubelet restart`.

**Note:** This credential-provider has to be present on every node in the cluster where you want to pull private images.

## Features
 * Authenticate private registry mirrors
 * Use password protected pull-through caches
 * Increase docker-pull rate-limit by always providing credentials
 * Smaller deployment config for developers
 * Use core components form protected/trusted registries - like `registry.k8s.io/kube-apiserver`

## Example
```bash
echo '
{
  "apiVersion": "credentialprovider.kubelet.k8s.io/v1",
  "kind": "CredentialProviderRequest",
  "image": "your.registry.example.org/org/image:version"
}' | jq -c | static-credential-provider [--config scp-conf.yaml] | jq
```

## Docs
For more info see:
 * [Credential-provider-spec](https://kubernetes.io/docs/reference/config-api/kubelet-credentialprovider.v1/)
 * [kubelet-provider-setup](https://kubernetes.io/docs/tasks/administer-cluster/kubelet-credential-provider/)
 * [kubelet-args](https://kubernetes.io/docs/reference/command-line-tools-reference/kubelet/)
