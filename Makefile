buildServer:
	CGO_ENABLED=0 GOOS=linux go build -o serverMain ./cmd/server/

.DEFAULT_GOAL = buildServer