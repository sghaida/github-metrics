
install-deps:
	go mod download

lint:
	golangci-lint --version &> /dev/null
	@if [ $? -ne 0 ];then\
  		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin;\
  	if\
  	golangci-lint run


run: install-deps
	@if [ ! -w $(OUTPUT_PATH) ];then\
		echo "$(OUTPUT_PATH) is not writable, please update the permission";\
		exit -1;\
	fi\

	go build -v -o github-metrics
	./github-metrics -out $(OUTPUT_PATH)

docker-build:
	docker build -t github-metrics .

docker-run: docker-build
	mkdir -p /tmp/github-metrics &> /dev/null

	docker run \
	  -e OUTPUT_PATH=/tmp/github-metrics \
	  -v /tmp/github-metrics:/tmp/github-metrics \
	  github-metrics





