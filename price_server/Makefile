
version=$(shell head -n 1 VERSION)

warehouse:
	go build -o build/warehouse .


build-docker:
	docker build -t aresprotocollab/warehouse:latest -f docker/Dockerfile ../

push-docker:
	docker tag aresprotocollab/warehouse:latest aresprotocollab/warehouse:${version}
	docker push aresprotocollab/warehouse:latest
	docker push aresprotocollab/warehouse:${version}