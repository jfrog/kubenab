{{/* vim: set filetype=mustache: */}}
{{/*
Expand the name of the chart.
*/}}
{{- define "kubenab.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "kubenab.fullname" -}}
{{- if .Values.fullnameOverride -}}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- $name := default .Chart.Name .Values.nameOverride -}}
{{- if contains $name .Release.Name -}}
{{- .Release.Name | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- end -}}
{{- end -}}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "kubenab.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "kubenab.selfSignedIssuer" -}}
{{ printf "%s-selfsign" (include "kubenab.fullname" .) }}
{{- end -}}

{{- define "kubenab.rootCAIssuer" -}}
{{ printf "%s-ca" (include "kubenab.fullname" .) }}
{{- end -}}

{{- define "kubenab.rootCACertificate" -}}
{{ printf "%s-ca" (include "kubenab.fullname" .) }}
{{- end -}}

{{- define "kubenab.servingCertificate" -}}
{{ printf "%s-kubenab-tls" (include "kubenab.fullname" .) }}
{{- end -}}
