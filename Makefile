.PHONY: install
install:
	go get
	go install

.PHONY: test
test:
	go test ./...
