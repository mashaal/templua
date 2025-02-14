module templua

go 1.21

require (
	github.com/fsnotify/fsnotify v1.8.0
	github.com/gorilla/websocket v1.5.3
	github.com/labstack/echo/v4 v4.11.4
	github.com/mashaal/tempura v0.0.0
)

replace github.com/mashaal/tempura => ./tempura

require (
	github.com/labstack/gommon v0.4.2 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasttemplate v1.2.2 // indirect
	github.com/yuin/gopher-lua v1.1.1 // indirect
	golang.org/x/crypto v0.17.0 // indirect
	golang.org/x/net v0.19.0 // indirect
	golang.org/x/sys v0.15.0 // indirect
	golang.org/x/text v0.14.0 // indirect
)
