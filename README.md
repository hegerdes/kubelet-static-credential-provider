# Kubelet-static-credential-provider
---
[![Test Build](https://github.com/hegerdes/kubelet-static-credential-provider/actions/workflows/test.yml/badge.svg)](https://github.com/hegerdes/kubelet-static-credential-provider/actions/workflows/test.yml)
![GitHub Release](https://img.shields.io/github/v/release/hegerdes/kubelet-static-credential-provider)


This plugin implements the [kubernetes credential-provider-api](https://kubernetes.io/docs/reference/config-api/kubelet-credentialprovider.v1/). This allows users to pull images from private image registries without having to create and reference a `image-pull-secret` in every namespace and deployment.

AWS, Azure and Google all use their own version of this plugin protocol to allow user friendly image pulls from their hosted image registries. Now you can also use this developer friendly approach to pull images from your private registries or password protected pull-through caches of DockerHub.

## Quickstart
Download the latest version from the release page:

```bash
# Download
wget https://<RELEASE_URL> -O static-credential-provider.tar.gz
mkdir -p /var/lib/kubelet-plugins/
tar xzf static-credential-provider.tar.gz -C /var/lib/kubelet-plugins/
```
And create a config:
```yaml
# /srv/kscp-conf.yaml
apiVersion: kubelet.config.k8s.io/v1
kind: CredentialProviderConfig
providers:
  # The name must match the name of the credential-provider binary
  - name: static-credential-provider
    # You can also use a config file instead of envs
    # args:[--config, <path-to-password-conf>]
    matchImages: [docker.io]
    defaultCacheDuration: "12h"
    apiVersion: credentialprovider.kubelet.k8s.io/v1
    env:
      - name: KSCP_REGISTRY_USERNAME
        value: <my-user>
      - name: KSCP_REGISTRY_PASSWORD
        value: <my-password>
  # Optional: More credential providers
```

Add the following args to the kubelet startup args and restart the kubelet with `service kubelet restart`.
 * `--image-credential-provider-bin-dir=/var/lib/kubelet-plugins/`
 * `--image-credential-provider-config=/srv/kscp-conf.yaml`

If you create a new cluster you can set these as your `kubeletExtraArgs` in your [Kubeadm NodeRegistrationOptions  under InitConfiguration](https://kubernetes.io/docs/reference/config-api/kubeadm-config.v1beta4/#kubeadm-k8s-io-v1beta4-NodeRegistrationOptions) or manually edit the `/var/lib/kubelet/kubeadm-flags.env` file and add the above args.

**Note:** This credential-provider has to be present on every node in the cluster where you want to pull private images.

## Features
 * Authenticate private registry mirrors
 * Use password protected pull-through caches
 * Increase docker-pull rate-limit by always providing credentials
 * Smaller deployment config for developers
 * Use core components form protected/trusted registries - like `registry.k8s.io/kube-apiserver`
 * Choose between a native go binary or a generic bash script.

## Example
```bash
export KSCP_REGISTRY_USERNAME=my-user
export KSCP_REGISTRY_PASSWORD=my-password
echo '
{
  "apiVersion": "credentialprovider.kubelet.k8s.io/v1",
  "kind": "CredentialProviderRequest",
  "image": "your.registry.example.org/org/image:version"
}' | jq -c | static-credential-provider [--config scp-conf.yaml] | jq
```

## Alternatives
 * There is a pure bash version in the [hack](/hack) folder that should work on every linux system with jq installed. It is also included in every release.
 * There is a python implementation from [JonTheNiceGuy](https://github.com/JonTheNiceGuy) see [generic-credential-provider](https://github.com/JonTheNiceGuy/generic-credential-provider)
 * There is the [AWS ECR](https://cloud-provider-aws.sigs.k8s.io/credential_provider/) variant that is used by default for ever EKS cluster

## Contributing
Found a bug? Or want to add an feature? Feel free to open an [issue](https://github.com/hegerdes/kubelet-static-credential-provider/issues) or [PR](https://github.com/hegerdes/kubelet-static-credential-provider/pulls).

## Docs
For more info see:
 * [Credential-provider-spec](https://kubernetes.io/docs/reference/config-api/kubelet-credentialprovider.v1/)
 * [kubelet-provider-setup](https://kubernetes.io/docs/tasks/administer-cluster/kubelet-credential-provider/)
 * [kubelet-args](https://kubernetes.io/docs/reference/command-line-tools-reference/kubelet/)
