# It checks your values.yaml for a user-provided nameOverride. If found, it uses that.
# If not, it defaults to the name defined in Chart.yaml (.Chart.Name).
# Crucially: It applies the trunc 63 and trimSuffix "-" functions.
# This is essential because Kubernetes object names (like labels) cannot be longer than 63 characters and cannot end in a hyphen.
# This template ensures the base name is always compliant.

{{/*
Expand the name of the chart.
*/}}
{{- define "my-full-app-chart.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because of K8s name limits.
*/}}
{{- define "my-full-app-chart.fullname" -}}
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
Common labels used throughout the chart.
This is the template called by {{- include "my-full-app-chart.labels" . | nindent 4 }}.
*/}}
{{- define "my-full-app-chart.labels" -}}
helm.sh/chart: {{ include "my-full-app-chart.name" . }}-{{ .Chart.AppVersion | replace "+" "_" }}
app.kubernetes.io/name: {{ include "my-full-app-chart.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end -}}