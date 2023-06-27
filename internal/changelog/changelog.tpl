{{if .AiSummarize }}
{{ .AiSummarize }}

Commits:
{{end}}
{{range .Commits}}- [{{ .Message }}]({{ $.Config.VCSURL }}/{{ .Hash }})
{{end}}