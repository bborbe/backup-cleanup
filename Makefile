install:
	GOBIN=$(GOPATH)/bin GO15VENDOREXPERIMENT=1 go install bin/backup_cleanup_cron/backup_cleanup_cron.go
test:
	GO15VENDOREXPERIMENT=1 go test `glide novendor`
check:
	golint ./...
	errcheck ./...
run:
	backup_cleanup_cron \
	-loglevel=DEBUG \
	-dir=/tmp \
	-lock=/tmp/backup_cleanup_cron.lock \
	-match=backup-to-devel
open:
	open http://localhost:8080/
format:
	find . -name "*.go" -exec gofmt -w "{}" \;
	goimports -w=true .
prepare:
	npm install
	go get -u golang.org/x/tools/cmd/goimports
	go get -u github.com/Masterminds/glide
	go get -u github.com/golang/lint/golint
	go get -u github.com/kisielk/errcheck
	glide install
update:
	glide up
clean:
	rm -rf vendor
