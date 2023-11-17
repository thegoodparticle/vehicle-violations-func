## Steps to Run

### Install Go (v1.21+)

https://go.dev/doc/install

### Download dependent libraries/modules

Run
```
cd vehicle-violations-func/

go mod tidy
```

### To run service locally

```
go run main.go
```

### Build docker file

I've used buildx and mentioned the platform as x86_64 (as my system builds images with amd64)

```
docker buildx build -f Dockerfile -t 237291162833.dkr.ecr.us-east-1.amazonaws.com/vehicle-violations-func:v1.0.3 --platform linux/x86_64 .
```

### Push Docker image to ECR

To push the docker images to ECR, the local system needs aws CLI, AWS credentials and config also added in local. 

Installing AWS CLI - https://docs.aws.amazon.com/cli/latest/userguide/getting-started-install.html

Configuring CLI - https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-configure.html

```
aws ecr get-login-password --region us-east-1 | docker login --username AWS --password-stdin 237291162833.dkr.ecr.us-east-1.amazonaws.com

docker push 237291162833.dkr.ecr.us-east-1.amazonaws.com/vehicle-violations-func:v1.0.x
```
