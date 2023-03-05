.PHONY: test
test.unit:
	echo "=> Running Tests"
	go test -tags=unit -v ./...

.PHONY: build
build:
	echo "=> Building..."
	CGO_ENABLED=0 go build -a -ldflags '-w -s' -o bin/httpnotifier

.PHONY: run
run:
	./bin/httpnotifier

.PHONY: generate
generate:
	mockgen -source=./pkg/notification/pool.go --build_flags=--mod=mod -package mock -destination=./mock/pool_mock.go . Pool
	mockgen -source=./pkg/notification/task.go --build_flags=--mod=mod -package mock -destination=./mock/task_mock.go . Task