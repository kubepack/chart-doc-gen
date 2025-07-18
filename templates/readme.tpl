# {{ .Project.ShortName }}

[{{ .Project.Name }}]({{ .Project.URL }}) - {{ .Project.Description }}

## TL;DR;

{{- if ne (slice .Repository.URL 0 6) "oci://" }}
```bash
$ helm repo add {{ .Repository.Name }} {{ .Repository.URL }}
$ helm repo update
$ helm search repo {{ .Repository.Name }}/{{ .Chart.Name }}{{ with .Chart.Version }} --version={{.}}{{ end }}
$ helm upgrade -i {{ .Release.Name }} {{ .Repository.Name }}/{{ .Chart.Name }} -n {{ .Release.Namespace }} --create-namespace{{ with .Chart.Version }} --version={{.}}{{ end }}
```
{{ else }}
```bash
$ helm upgrade -i {{ .Release.Name }} {{ .Repository.URL }}/{{ .Chart.Name }} -n {{ .Release.Namespace }} --create-namespace{{ with .Chart.Version }} --version={{.}}{{ end }}
```
{{- end }}

## Introduction

This chart deploys {{ .Project.App }} on a [Kubernetes](http://kubernetes.io) cluster using the [Helm](https://helm.sh) package manager.

## Prerequisites
{{ range .Prerequisites }}
- {{ . }}
{{- end }}

## Installing the Chart

To install/upgrade the chart with the release name `{{ .Release.Name }}`:

{{- if ne (slice .Repository.URL 0 6) "oci://" }}
```bash
$ helm upgrade -i {{ .Release.Name }} {{ .Repository.Name }}/{{ .Chart.Name }} -n {{ .Release.Namespace }} --create-namespace{{ with .Chart.Version }} --version={{.}}{{ end }}
```
{{ else }}
```bash
$ helm upgrade -i {{ .Release.Name }} {{ .Repository.URL }}/{{ .Chart.Name }} -n {{ .Release.Namespace }} --create-namespace{{ with .Chart.Version }} --version={{.}}{{ end }}
```
{{- end }}

The command deploys {{ .Project.App }} on the Kubernetes cluster in the default configuration. The [configuration](#configuration) section lists the parameters that can be configured during installation.

> **Tip**: List all releases using `helm list`

## Uninstalling the Chart

To uninstall the `{{ .Release.Name }}`:

```bash
$ helm uninstall {{ .Release.Name }} -n {{ .Release.Namespace }}
```

The command removes all the Kubernetes components associated with the chart and deletes the release.

{{ if .Chart.Values -}}
## Configuration

The following table lists the configurable parameters of the `{{ .Chart.Name }}` chart and their default values.

{{ .Chart.Values }}

Specify each parameter using the `--set key=value[,key=value]` argument to `helm upgrade -i`. For example:

{{- if ne (slice .Repository.URL 0 6) "oci://" }}
```bash
$ helm upgrade -i {{ .Release.Name }} {{ .Repository.Name }}/{{ .Chart.Name }} -n {{ .Release.Namespace }} --create-namespace{{ with .Chart.Version }} --version={{.}}{{ end }} --set {{ .Chart.ValuesExample }}
```
{{ else }}
```bash
$ helm upgrade -i {{ .Release.Name }} {{ .Repository.URL }}/{{ .Chart.Name }} -n {{ .Release.Namespace }} --create-namespace{{ with .Chart.Version }} --version={{.}}{{ end }} --set {{ .Chart.ValuesExample }}
```
{{- end }}

Alternatively, a YAML file that specifies the values for the parameters can be provided while
installing the chart. For example:

{{- if ne (slice .Repository.URL 0 6) "oci://" }}
```bash
$ helm upgrade -i {{ .Release.Name }} {{ .Repository.Name }}/{{ .Chart.Name }} -n {{ .Release.Namespace }} --create-namespace{{ with .Chart.Version }} --version={{.}}{{ end }} --values values.yaml
```
{{ else }}
```bash
$ helm upgrade -i {{ .Release.Name }} {{ .Repository.URL }}/{{ .Chart.Name }} -n {{ .Release.Namespace }} --create-namespace{{ with .Chart.Version }} --version={{.}}{{ end }} --values values.yaml
```
{{- end }}

{{- end }}
