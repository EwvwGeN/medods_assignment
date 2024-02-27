buildServer:
	CGO_ENABLED=0 GOOS=linux go build -o serverMain ./cmd/server/
prepareEnv:
	go run config_to_env.go ./configs/config.yaml
.DEFAULT_GOAL = buildServer