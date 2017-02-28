version = '0.1.0'
arch = 'amd64'
name = Gworker

default: .glide.run test build

.glide.run: glide.yaml
	touch .glide.run
	glide install

test:
	go test $(glide novendor)
	go vet $(glide novendor)

build:
	go build -i -o ./bin/$(name) ./

clean:
	$(RM) -rf build/*
	$(RM) bin/*

.PHONY: clean test build
