module github.com/guregodevo/loopchain_tools

go 1.23.2

replace github.com/guregodevo/mario-llm => ../mario-llm

require (
	github.com/chromedp/chromedp v0.11.1
	github.com/tmc/langchaingo v0.1.12
)

require (
	github.com/chromedp/cdproto v0.0.0-20241022234722-4d5d5faf59fb // indirect
	github.com/chromedp/sysutil v1.1.0 // indirect
	github.com/dlclark/regexp2 v1.10.0 // indirect
	github.com/gobwas/httphead v0.1.0 // indirect
	github.com/gobwas/pool v0.2.1 // indirect
	github.com/gobwas/ws v1.4.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/pkoukk/tiktoken-go v0.1.7 // indirect
    go.starlark.net v0.0.0-20240925182052-1207426daebd // Explicitly added
	golang.org/x/sys v0.26.0 // indirect
)
