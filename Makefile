init:
	@wire

run:
	@go run main.go wire_gen.go

push-build:
	@docker buildx build --push -t registry.fly.io/partita:latest -f Dockerfile .

deploy:
	flyctl m run -a partita \
	--memory 256 \
	--cpus 1 \
	--region sjc \
	-p 443:3000/tcp:http:tls \
	registry.fly.io/partita:latest 

roll:
	flyctl m restart -a partita