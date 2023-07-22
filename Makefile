.PHONY: install tidy clean

install:
	go mod download

tidy:
	go mod tidy
	go mod vendor

clean:
	go clean -cache -modcache
