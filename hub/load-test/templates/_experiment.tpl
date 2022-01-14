{{ define "load-test.experiment" -}}
# task 1: generate HTTP requests for application URL
# collect Iter8's built-in latency and error-related metrics
- task: gen-load-and-collect-metrics
  with:
    {{- if .Values.numQueries }}
    numQueries: {{ .Values.numQueries}}
    {{- end }}
    {{- if .Values.duration }}
    duration: {{ .Values.duration}}
    {{- end }}
    {{- if .Values.qps }}
    qps: {{ .Values.qps}}
    {{- end }}
    {{- if .Values.connections }}
    connections: {{ .Values.connections}}
    {{- end }}
    {{- if .Values.payloadStr }}
    payloadStr: {{ .Values.payloadStr}}
    {{- end }}
    {{- if .Values.payloadURL }}
    payloadURL: {{ .Values.payloadURL}}
    {{- end }}
    {{- if .Values.contentType }}
    contentType: {{ .Values.contentType}}
    {{- end }}
    {{- if .Values.errorRanges }}
    errorRanges:
{{ toYaml .Values.errorRanges | indent 4 }}
    {{- end }}    
    {{- if .Values.percentiles }}
    percentiles:
{{ toYaml .Values.percentiles | indent 4 }}
    {{- end }}
    versionInfo:
    - url: {{ required "A valid url value is required!" .Values.url }}
    {{- if .Values.headers }}
    headers:
{{ toYaml .Values.headers | indent 6 }}
    {{- end }}
# task 2: validate service level objectives for app using
# the metrics collected in the above task
- task: assess-app-versions
  with:
  {{- if .Values.SLOs }}
    SLOs:
{{ toYaml .Values.SLOs | indent 4 }}
  {{- end }}
{{ end }}