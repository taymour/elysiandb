tcp_benchmark:
	./elysian_bench -addr 127.0.0.1:8088   -vus 500 -duration 20s -keys 20000 -payload 16 -pair
http_benchmark:
	./benchmark.sh


.PHONY: test test-e2e-http test-unit test-cover

test:
	@pkgs=$$(go list -f '{{if or (len .TestGoFiles) (len .XTestGoFiles)}}{{.ImportPath}}{{end}}' ./... | grep -v '^$$'); \
	if [ -z "$$pkgs" ]; then echo "no test packages"; exit 0; fi; \
	go test $$pkgs -v -race -count=1

test-e2e-http:
	go test ./test/e2e/http -v -race -count=1

test-unit:
	@pkgs=$$(go list ./... | grep -v '^github.com/[^/]*/[^/]*/test/'); \
	if [ -z "$$pkgs" ]; then echo "no unit test packages"; exit 0; fi; \
	go test $$pkgs -v -race -count=1
test-cover:
	@pkgs=$$(go list -f '{{if or (len .TestGoFiles) (len .XTestGoFiles)}}{{.ImportPath}}{{end}}' ./... | grep -v '^$$'); \
	if [ -z "$$pkgs" ]; then echo "no test packages"; exit 0; fi; \
	go test $$pkgs -race -coverprofile=coverage.out -count=1 && go tool cover -func=coverage.out
