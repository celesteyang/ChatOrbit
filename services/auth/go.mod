module auth

go 1.22

require github.com/celesteyang/ChatOrbit/shared/logger v0.1.0

require (
	go.uber.org/multierr v1.10.0 // indirect
	go.uber.org/zap v1.27.0 // indirect
)

replace github.com/celesteyang/ChatOrbit/shared/logger => ../../shared/logger
