up:
	docker-compose up -d

build:
	cd cmd/shortener/ && go build -o shortener

run:
	go run cmd/shortener/main.go

run-w-db:
	go run cmd/shortener/main.go -d "host=localhost port=5432 user=myuser password=mypassword dbname=mydatabase sslmode=disable"

test10:
	shortenertestbeta -test.v -test.run=^TestIteration10$ \
          -binary-path=cmd/shortener/shortener \
          -source-path=. \
          -database-dsn='host=localhost port=5432 user=myuser password=mypassword dbname=mydatabase sslmode=disable'

test11:
	shortenertestbeta -test.v -test.run=^TestIteration11$ \
          -binary-path=cmd/shortener/shortener \
          -database-dsn='host=localhost port=5432 user=myuser password=mypassword dbname=mydatabase sslmode=disable'

gotest:
	go test -v ./...