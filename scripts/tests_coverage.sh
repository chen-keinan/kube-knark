#!/bin/sh
go install github.com/golang/mock/mockgen@latest
go install -v github.com/golang/mock/mockgen && export PATH=$GOPATH/bin:$PATH;
go generate ./...
go test ./cmd... ./internal... ./pkg... -coverprofile coverage.md fmt
go tool cover -html=coverage.md -o coverage.html
go tool cover  -func coverage.md
go install github.com/jpoles1/gopherbadger
gopherbadger -md="README.md,coverage.md" -manualcov=$(go tool cover -func coverage.md | grep total | awk '{print substr($3, 1, length($3)-1)}') -prefix=""
mv coverage_badge.png ./pkg/images/coverage_badge.png
