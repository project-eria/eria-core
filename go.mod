module github.com/project-eria/eria-core

go 1.17

require (
	github.com/project-eria/go-wot v1.0.1
	github.com/rs/zerolog v1.26.1
	gopkg.in/yaml.v3 v3.0.0-20200313102051-9f266ea9e77c
)

require (
	github.com/andybalholm/brotli v1.0.4 // indirect
	github.com/fasthttp/websocket v1.5.0 // indirect
	github.com/gofiber/fiber/v2 v2.31.0 // indirect
	github.com/gofiber/websocket/v2 v2.0.20 // indirect
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/klauspost/compress v1.15.1 // indirect
	github.com/savsgio/gotils v0.0.0-20220401102855-e56b59f40436 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasthttp v1.35.0 // indirect
	github.com/valyala/tcplisten v1.0.0 // indirect
	golang.org/x/sys v0.0.0-20220412211240-33da011f77ad // indirect
)

// replace github.com/project-eria/go-wot => ../go-wot

retract (
	v0.3.2
	v0.3.1
	v0.3.0
	v0.2.2
	v0.2.1
	v0.2.0
	v0.1.0
	v0.0.2
	v0.0.1
	v0.0.0
)
