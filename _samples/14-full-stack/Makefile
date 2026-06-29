# ============================================================================
# stargo 全栈示例 Makefile
# ============================================================================

.PHONY: all build run clean proto docker-build docker-up docker-down k8s-apply k8s-delete

# 构建所有服务 / Build all services
build:
	cd user-service && go build -o ../bin/user-service . &
	cd product-service && go build -o ../bin/product-service . &
	cd order-service && go build -o ../bin/order-service . &
	cd gateway && go build -o ../bin/gateway . &
	wait && echo "All services built to bin/"

# 本地运行 (需要先启动 MySQL/Redis/etcd/Jaeger)
# Local run (requires MySQL/Redis/etcd/Jaeger running first)
run:
	cd user-service && go run . &
	sleep 2
	cd product-service && go run . &
	sleep 2
	cd order-service && go run . &
	sleep 2
	cd gateway && go run . &
	wait

# 清理 / Clean
clean:
	rm -rf bin/

# 生成 proto 代码 / Generate proto code
proto:
	protoc --go_out=. --go_opt=paths=source_relative \
	       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
	       proto/user/user.proto
	protoc --go_out=. --go_opt=paths=source_relative \
	       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
	       proto/product/product.proto

# Docker 构建 / Docker build
docker-build:
	docker compose build

# Docker 启动 / Docker start
docker-up:
	docker compose up -d

# Docker 停止 / Docker stop
docker-down:
	docker compose down

# K8s 部署 / K8s deploy
k8s-apply:
	kubectl apply -f k8s/

# K8s 清理 / K8s cleanup
k8s-delete:
	kubectl delete namespace stargo-demo

# 查看 K8s 状态 / Check K8s status
k8s-status:
	kubectl -n stargo-demo get all

# 测试 API / Test API
test-api:
	@echo "=== Health Check ==="
	curl -s http://localhost:8080/healthz | jq .
	@echo "=== Metrics ==="
	curl -s http://localhost:8080/metrics | head -20
