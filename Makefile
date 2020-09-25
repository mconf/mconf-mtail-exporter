VERSION=0.1
build:
	docker build -t mconf/mconf-api:mtail-exporter-v${VERSION} .