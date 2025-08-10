.PHONY: build run test clean install

build:
	go build -o bin/akasha ./cmd/akasha

run: build
	./bin/akasha run

test:
	go test ./...

clean:
	rm -rf bin/

install: build
	@echo "安装 akasha 到 /usr/local/bin"
	@sudo cp bin/akasha /usr/local/bin/
	@echo "安装完成！现在可以直接使用 'akasha' 命令"

