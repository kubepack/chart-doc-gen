# Stash

[Stash by AppsCode](https://github.com/stashed/stash) - Backup your Kubernetes Volumes

## TL;DR;

```console
$ helm repo add appscode https://charts.appscode.com/stable/
$ helm repo update
$ helm install stash-operator appscode/stash -n kube-system --version=v0.9.0-rc.0
```

## Introduction

This chart deploys a Stash operator on a [Kubernetes](http://kubernetes.io) cluster using the [Helm](https://helm.sh) package manager.

## Prerequisites

- Kubernetes v1.14+
- `--allow-privileged` flag must be set to true for both the API server and the kubelet
- (If you use Docker) The Docker daemon of the cluster nodes must allow shared mounts
- Pre-installed HashiCorp Vault server.

## Installing the Chart

To install the chart with the release name `stash-operator`:

```console
$ helm install stash-operator appscode/stash -n kube-system --version=v0.9.0-rc.0
```

The command deploys a Stash operator on the Kubernetes cluster in the default configuration. The [configuration](#configuration) section lists the parameters that can be configured during installation.

> **Tip**: List all releases using `helm list`

## Uninstalling the Chart

To uninstall/delete the `stash-operator`:

```console
$ helm delete stash-operator -n kube-system
```

The command removes all the Kubernetes components associated with the chart and deletes the release.

## Configuration

The following table lists the configurable parameters of the `stash` chart and their default values.

|               Parameter               |                                                                                                                    Description                                                                                                                    |                                Default                                |
|---------------------------------------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|-----------------------------------------------------------------------|
| nameOverride                          | Overrides name template                                                                                                                                                                                                                           | `""`                                                                  |
| fullnameOverride                      | Overrides fullname template                                                                                                                                                                                                                       | `""`                                                                  |
| replicaCount                          | Number of stash operator replicas to create (only 1 is supported)                                                                                                                                                                                 | `1`                                                                   |
| operator.registry                     | Docker registry used to pull operator image                                                                                                                                                                                                       | `appscode`                                                            |
| operator.repository                   | Name of operator container image                                                                                                                                                                                                                  | `stash`                                                               |
| operator.tag                          | Operator container image tag                                                                                                                                                                                                                      | `v0.9.0-rc.6`                                                         |
| operator.resources                    | Compute Resources required by the operator container                                                                                                                                                                                              | `{"requests":{"cpu":"100m"}}`                                         |
| operator.securityContext              | Security options the operator container should run with                                                                                                                                                                                           | `{}`                                                                  |
| pushgateway.registry                  | Docker registry used to pull Prometheus pushgateway image                                                                                                                                                                                         | `prom`                                                                |
| pushgateway.repository                | Prometheus pushgateway container image                                                                                                                                                                                                            | `pushgateway`                                                         |
| pushgateway.tag                       | Prometheus pushgateway container image tag                                                                                                                                                                                                        | `v0.5.2`                                                              |
| pushgateway.resources                 | Compute Resources required by the Prometheus pushgateway container                                                                                                                                                                                | `{}`                                                                  |
| pushgateway.securityContext           | Security options the Prometheus pushgateway container should run with                                                                                                                                                                             | `{}`                                                                  |
| cleaner.registry                      | Docker registry used to pull Webhook cleaner image                                                                                                                                                                                                | `appscode`                                                            |
| cleaner.repository                    | Webhook cleaner container image                                                                                                                                                                                                                   | `kubectl`                                                             |
| cleaner.tag                           | Webhook cleaner container image tag                                                                                                                                                                                                               | `v1.16`                                                               |
| imagePullSecrets                      | Specify an array of imagePullSecrets. Secrets must be manually created in the namespace. <br> Example: <br> `helm template charts/stash \` <br> `--set imagePullSecrets[0].name=sec0 \` <br> `--set imagePullSecrets[1].name=sec1`                | `[]`                                                                  |
| imagePullPolicy                       | Container image pull policy                                                                                                                                                                                                                       | `IfNotPresent`                                                        |
| criticalAddon                         | If true, installs Stash operator as critical addon                                                                                                                                                                                                | `false`                                                               |
| logLevel                              | Log level for operator                                                                                                                                                                                                                            | `3`                                                                   |
| annotations                           | Annotations applied to operator deployment                                                                                                                                                                                                        | `{}`                                                                  |
| podAnnotations                        | Annotations passed to operator pod(s).                                                                                                                                                                                                            | `{}`                                                                  |
| nodeSelector                          | Node labels for pod assignment                                                                                                                                                                                                                    | `{"beta.kubernetes.io/arch":"amd64","beta.kubernetes.io/os":"linux"}` |
| tolerations                           | Tolerations for pod assignment                                                                                                                                                                                                                    | `[]`                                                                  |
| affinity                              | Affinity rules for pod assignment                                                                                                                                                                                                                 | `{}`                                                                  |
| podSecurityContext                    | Security options the operator pod should run with.                                                                                                                                                                                                | `{"fsGroup":65535}`                                                   |
| serviceAccount.create                 | Specifies whether a service account should be created                                                                                                                                                                                             | `true`                                                                |
| serviceAccount.annotations            | Annotations to add to the service account                                                                                                                                                                                                         | `{}`                                                                  |
| serviceAccount.name                   | The name of the service account to use. If not set and create is true, a name is generated using the fullname template                                                                                                                            | ``                                                                    |
| apiserver.groupPriorityMinimum        | The minimum priority the webhook api group should have at least. Please see https://github.com/kubernetes/kube-aggregator/blob/release-1.9/pkg/apis/apiregistration/v1beta1/types.go#L58-L64 for more information on proper values of this field. | `10000`                                                               |
| apiserver.versionPriority             | The ordering of the webhook api inside of the group. Please see https://github.com/kubernetes/kube-aggregator/blob/release-1.9/pkg/apis/apiregistration/v1beta1/types.go#L66-L70 for more information on proper values of this field              | `15`                                                                  |
| apiserver.enableMutatingWebhook       | If true, mutating webhook is configured for Kubernetes workloads                                                                                                                                                                                  | `true`                                                                |
| apiserver.enableValidatingWebhook     | If true, validating webhook is configured for Stash CRDss                                                                                                                                                                                         | `true`                                                                |
| apiserver.ca                          | CA certificate used by the Kubernetes api server. This field is automatically assigned by the operator.                                                                                                                                           | `not-ca-cert`                                                         |
| apiserver.bypassValidatingWebhookXray | If true, bypasses checks that validating webhook is actually enabled in the Kubernetes cluster.                                                                                                                                                   | `false`                                                               |
| apiserver.useKubeapiserverFqdnForAks  | If true, uses kube-apiserver FQDN for AKS cluster to workaround https://github.com/Azure/AKS/issues/522 (default true)                                                                                                                            | `true`                                                                |
| apiserver.healthcheck.enabled         | If true, enables the readiness and liveliness probes for the operator pod.                                                                                                                                                                        | `false`                                                               |
| apiserver.servingCerts.generate       | If true, generates on install/upgrade the certs that allow the kube-apiserver (and potentially ServiceMonitor) to authenticate operators pods. Otherwise specify certs in `apiserver.servingCerts.{caCrt, serverCrt, serverKey}`.                 | `true`                                                                |
| apiserver.servingCerts.caCrt          | CA certficate used by serving certificate of webhook server.                                                                                                                                                                                      | `""`                                                                  |
| apiserver.servingCerts.serverCrt      | Serving certficate used by webhook server.                                                                                                                                                                                                        | `""`                                                                  |
| apiserver.servingCerts.serverKey      | Private key for the serving certificate used by webhook server.                                                                                                                                                                                   | `""`                                                                  |
| enableAnalytics                       | If true, sends usage analytics                                                                                                                                                                                                                    | `true`                                                                |
| monitoring.agent                      | Name of monitoring agent (either "prometheus.io/operator" or "prometheus.io/builtin")                                                                                                                                                             | `"none"`                                                              |
| monitoring.backup                     | Specify whether to monitor Stash backup and recovery                                                                                                                                                                                              | `false`                                                               |
| monitoring.operator                   | Specify whether to monitor Stash operator                                                                                                                                                                                                         | `false`                                                               |
| monitoring.prometheus.namespace       | Specify the namespace where Prometheus server is running or will be deployed.                                                                                                                                                                     | `""`                                                                  |
| monitoring.serviceMonitor.labels      | Specify the labels for ServiceMonitor. Prometheus crd will select ServiceMonitor using these labels. Only usable when monitoring agent is `prometheus.io/operator`.                                                                               | `{}`                                                                  |
| additionalPodSecurityPolicies         | Additional psp names passed to operator <br> Example: <br> `helm install appscode/stash \` <br> `--set additionalPodSecurityPolicies[0]=abc \` <br> `--set additionalPodSecurityPolicies[1]=xyz`                                                  | `[]`                                                                  |
| platform.openshift                    | Set true, if installed in OpenShift                                                                                                                                                                                                               | `false`                                                               |
| empty                                 | This top-level empty map should be documented                                                                                                                                                                                                     | `{}`                                                                  |


Specify each parameter using the `--set key=value[,key=value]` argument to `helm install`. For example:

```console
$ helm install stash-operator appscode/stash -n kube-system --version=v0.9.0-rc.0 --set replicaCount=1
```

Alternatively, a YAML file that specifies the values for the parameters can be provided while
installing the chart. For example:

```console
$ helm install stash-operator appscode/stash -n kube-system --version=v0.9.0-rc.0 --values values.yaml
```
