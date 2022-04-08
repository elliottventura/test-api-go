[![Open in Gitpod](https://gitpod.io/button/open-in-gitpod.svg)](https://gitpod.io/projects/test-api-go)

# Install dependencies
go get .
# Launch the application
go run .

# To launch with docker
## Build the image
docker build --tag test-api-go:v1 --file=".\Dockerfile" .

Note: Since the docker file is named Dockefile, the parameter -- file is not mandatory. Do not forget the . at the end

## Launch the container
docker run --publish 3083:8083 test-api-go:v1

Note: 3083 is the port used to access the container, and it is mapped to the port 8083 (port of the api defined in the code). This parameter must be provided beofre the image.