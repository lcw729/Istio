docker_id="ketidevit2"
controller_name="hcp-loadbalancer"

export GO111MODULE=on
go mod vendor

go build -o build/_output/bin/$controller_name -gcflags all=-trimpath=`pwd` -asmflags all=-trimpath=`pwd` -mod=vendor ./src/main && \

# dockerfile build & image 생성
docker build -t $docker_id/$controller_name:v0.0.2 build && \

# image를 dockerhub에 push
docker push $docker_id/$controller_name:v0.0.2