module github.com/project-eria/eria-core

go 1.17

require (
	github.com/project-eria/go-wot v0.1.1
	github.com/rs/zerolog v1.26.1
	gopkg.in/yaml.v3 v3.0.0-20200313102051-9f266ea9e77c
)

require (
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/julienschmidt/httprouter v1.3.0 // indirect
)

// replace github.com/project-eria/go-wot => ../go-wot
// replace github.com/project-eria/eria-core => ../eria-core