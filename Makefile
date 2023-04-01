.PHONY: docker-build docker-run lint

IMAGE_NAME := wifi-control-server
CONTAINER_NAME := wifi-control-server
PORT := 8080

docker-build:
		docker build -t $(IMAGE_NAME) .

docker-run:
		docker run -it --rm -p $(PORT):$(PORT) --name $(CONTAINER_NAME) $(IMAGE_NAME)

# docker run -d --name wifi-control-server -p 8080:8080 -e ROUTER_URL="http://10.0.0.1" -e USERNAME="admin" -e PASSWORD="secret" wifi-control-server
