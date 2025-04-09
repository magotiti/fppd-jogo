.PHONY: all build

all: build

go.mod:
	go mod init jogo
	go get -u github.com/nsf/termbox-go

build: go.mod
	go build -o jogo
	
clean:
	rm -f jogo

distclean: clean
	rm -f go.mod go.sum
