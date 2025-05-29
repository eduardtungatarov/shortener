up:
	docker-compose up -d

build:
	cd cmd/shortener/ && go build -o shortener

test10:
	shortenertestbeta -test.v -test.run=^TestIteration10$ \
          -binary-path=cmd/shortener/shortener

test11:
	shortenertestbeta -test.v -test.run=^TestIteration11$ \
          -binary-path=cmd/shortener/shortener

gotest:
	go test -v ./...