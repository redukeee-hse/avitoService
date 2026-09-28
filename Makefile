.PHONY: generate
generate:
	go tool oapi-codegen \
		-generate types,chi-server \
		-package api \
		-include-operation-ids createTrip,getTrip,finishTrip,health,ready \
		-o internal/generated/api.gen.go \
		contracts/openapi/trip-service.openapi.yaml
.PHONY: run
run:
	go build -o bin/trip-service ./cmd/trip-service/ && ./bin/trip-service
	