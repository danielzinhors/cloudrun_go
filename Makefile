generate-mocks:
	mockgen -source=./internal/usecases/get_temp.go -destination ./internal/usecases/mocks/get_temp.go -package mocks
	mockgen -source=./internal/services/viacep.go -destination ./internal/services/mocks/viacep.go -package mocks
	mockgen -source=./internal/services/weatherapi.go -destination ./internal/services/mocks/weatherapi.go -package mocks

test:
	go test ./... -v
