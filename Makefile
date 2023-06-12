
install-deps:
	go mod download

lint:
	golangci-lint --version &> /dev/null
	@if [ $? -ne 0 ];then\
  		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin;\
  	if\
  	golangci-lint run

run: install-deps
	go run .



