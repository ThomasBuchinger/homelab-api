
build: ui/out/
	mkdir -p ui/out/
	cd ui/ && npm run build
	go build .

test:
	go test ./...

build-container-image:
	podman build --tag homelab-api .

run-container:
	podman run -it --publish 8080:8080 homelab-api:latest
