{{- define "was.serviceAccountName" -}}
  {{- if and .Values.was.enabled .Values.was.serviceAccount.create }}
    {{- default (printf "%s-service-account" .Release.Name) .Values.was.serviceAccount.name }}
  {{- else }}
    {{- default "default" .Values.was.serviceAccount.name }}
  {{- end }}
{{- end }}


{{- define "security_context" -}}
runAsUser: {{ .Values.uid }}
fsGroup: {{ .Values.gid }}
{{- end }}


{{- define "was_image" -}}
  {{- $repository := .Values.images.was.repository | default .Values.images.defaultWasRepository -}}
  {{- $tag := .Values.images.was.tag | default .Values.images.defaultWasTag -}}
  {{- $digest := .Values.images.was.digest | default .Values.images.defaultWasDigest -}}
  {{- if $digest }}
    {{- printf "%s@%s" $repository $digest -}}
  {{- else }}
    {{- printf "%s:%s" $repository $tag -}}
  {{- end }}
{{- end }}


{{- define "registry_docker_config" }}
  {{- $host := .Values.images.registry.secrets.connection.host }}
  {{- $email := .Values.images.registry.secrets.connection.email }}
  {{- $user := .Values.images.registry.secrets.connection.user }}
  {{- $pass := .Values.images.registry.secrets.connection.pass }}

  {{- $config := dict "auths" }}
  {{- $auth := dict }}
  {{- $data := dict }}
  {{- $_ := set $data "username" $user }}
  {{- $_ := set $data "password" $pass }}
  {{- $_ := set $data "email" $email }}
  {{- $_ := set $data "auth" (printf "%v:%v" $user $pass | b64enc) }}
  {{- $_ := set $auth $host $data }}
  {{- $_ := set $config "auths" $auth }}
  {{ $config | toJson | print }}
{{- end }}


{{- define "was_env_from" }}
  {{- $Global := . }}
  {{- with .Values.was.envFrom }}
    {{- tpl . $Global | nindent 2 }}
  {{- end }}
{{- end }}


{{- define "was_env" }}
  {{- range $i, $env := .Values.was.env }}
- name: {{ $env.name }}
  value: {{ $env.value | quote }}
  {{- end }}

  {{- $releaseName := .Release.Name -}}
  {{- if (and .Values.was.secrets.localData (eq .Values.was.secrets.type "local")) }}
    {{- range $i, $config := .Values.was.secrets.localData }}
- name: {{ $config.name }}
  valueFrom: 
    secretKeyRef:
      name: {{ printf "%s-secrets" $releaseName }}
      key: {{ $config.name }}
    {{- end }}
  {{- else if (and .Values.was.secrets.spec.data (eq .Values.was.secrets.type "cloud")) }}
    {{- range $i, $config := .Values.was.secrets.spec.data }}
- name: {{ $config.property }}
  valueFrom: 
    secretKeyRef:
      name: {{ printf "%s-secrets" $releaseName }}
      key: {{ $config.name }}
    {{- end }}
  {{- end }}
{{- end }}

