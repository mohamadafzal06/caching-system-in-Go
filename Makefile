build:
	go build -o bin/dcache
run: build
	./bin/dcache

runfollower:
	./bin/dcache --listenaddr :3000 --leaderaddr :4000

test:
	go test -v ./...
