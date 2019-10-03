.PHONY: docs, test
docs:
	@docker run -v $$PWD/:/docs pandoc/latex -f markdown /docs/README.md -o /docs/build/output/README.pdf

test:
	@docker run --rm \
	-v $$PWD/:/go/src/github.com/rhymond/interview-accountapi \
	-w /go/src/github.com/rhymond/interview-accountapi \
	golang:1.13.1-stretch \
	go get -d -v -t ./... && go test --cover -v ./...

lint:
