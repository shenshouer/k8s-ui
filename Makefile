all: push

# 0.0 shouldn't clobber any release builds
TAG = 0.1
PREFIX = dhub.yunpro.cn/shenshouer/k8s-ui

controller: main.go
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-w' -o k8s-ui ./main.go

container: controller
	docker build -t $(PREFIX):$(TAG) .

push: container
	docker push $(PREFIX):$(TAG)

clean:
	rm -f k8s-ui