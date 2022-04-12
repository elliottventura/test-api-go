[![Open in Gitpod](https://gitpod.io/button/open-in-gitpod.svg)](https://gitpod.io/projects/test-api-go)

# Install dependencies
go get .
# Launch the application
go run .

# To add mongo dependencies
go get go.mongodb.org/mongo-driver

# To launch with docker
## Build the image
docker build --tag test-api-go:v1 --file=".\Dockerfile" .

Note: Since the docker file is named Dockefile, the parameter -- file is not mandatory. Do not forget the . at the end

## Launch the container
docker run --publish 3083:8083 test-api-go:v1

Note: 3083 is the port used to access the container, and it is mapped to the port 8083 (port of the api defined in the code). This parameter must be provided before the image.

# To allow communication beetween containers : 

## Solution 1 (bad solution)
Find the internal ip of the container to access : docker inspect <container_id> | grep IPAddress
Use this ip instead of localhost in the url. The port is the port of the api (and not the container exposed port)
Here we use the bridge. Each container has its ip and we use the ip to communicate.

## Solution 2
Create your own network in which your applications will run.

- docker network create test-api-net
- docker run -i --rm --net test-api-net --publish 3083:8083 --name test-api-go test-api-go:v1

- Use docker network rm test-api-net to delete the network


# Summary
docker build --tag test-api-go:v1 --file=".\Dockerfile" .
docker network create test-api-net
docker run -i --rm --net test-api-net --publish 3083:8083 --name test-api-go test-api-go:v1