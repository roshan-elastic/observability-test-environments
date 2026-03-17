{{/* vim: set filetype=mustache: */}}
{{/*
Expand the name of the chart.
*/}}
{{- define "elastic-agent.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "elastic-agent.fullname" -}}
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
{{- define "elastic-agent.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Common labels
*/}}
{{- define "elastic-agent.labels" -}}
helm.sh/chart: {{ include "elastic-agent.chart" . }}
{{ include "elastic-agent.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- with .Values.labels }}
{{ toYaml . }}
{{- end }}
{{- end -}}

{{/*
Selector labels
*/}}
{{- define "elastic-agent.selectorLabels" -}}
app.kubernetes.io/name: {{ include "elastic-agent.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end -}}

{{/*
Create the name of the service account to use
*/}}
{{- define "elastic-agent.serviceAccountName" -}}
{{- if .Values.serviceAccount.create -}}
    {{ default (include "elastic-agent.fullname" .) .Values.serviceAccount.name }}
{{- else -}}
    {{ default "default" .Values.serviceAccount.name }}
{{- end -}}
{{- end -}}

{{/*
Kibana credentials name
*/}}
{{- define "elastic-agent.kibanaCreds" -}}
{{- if ((.Values.kibana).secretName) -}}
    {{ ((.Values.kibana).secretName) }}
{{- else -}}
    {{ template "elastic-agent.serviceAccountName" .}}-kibana-creds
{{- end -}}
{{- end -}}

{{/*
Elasticsearch credentials name
*/}}
{{- define "elastic-agent.esCreds" -}}
{{- if ((.Values.elasticsearch).secretName) -}}
    {{ ((.Values.elasticsearch).secretName) }}
{{- else -}}
    {{ template "elastic-agent.serviceAccountName" .}}-es-creds
{{- end -}}
{{- end -}}

{{/*
Fleet credentials name
*/}}
{{- define "elastic-agent.fleetCreds" -}}
{{- if .Values.fleet.secretName -}}
    {{ .Values.fleet.secretName }}
{{- else -}}
    {{ template "elastic-agent.serviceAccountName" .}}-fleet-creds
{{- end -}}
{{- end -}}

{{/*
Fleet Server credentials name
*/}}
{{- define "elastic-agent.fleetServerCreds" -}}
{{- if ((.Values.fleetServer).secretName) -}}
    {{ ((.Values.fleetServer).secretName) }}
{{- else -}}
    {{ template "elastic-agent.serviceAccountName" .}}-fleet-server-creds
{{- end -}}
{{- end -}}
