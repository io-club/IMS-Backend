MODE? = release

prepare:
	mkdir -p ./log/nms
	mkdir -p ./log/user
	mkdir -p ./log/device
	mkdir -p ./log/work
	mkdir -p ./target

run: prepare build
	@echo "Start service"
	# 启动 user
	MODE=${MODE}  ./target/user &
	# 启动 device
	MODE=${MODE}  ./target/mqtt &
	MODE=${MODE}  ./target/device &
	# 启动 work
	MODE=${MODE}  ./target/weather &
	MODE=${MODE}  ./target/work &
	# 启动 nms
	MODE=${MODE} ./target/nms
	@echo "Service has started"

run-fresh: build stop run

build:
	go build -o ./target/nms ./main.go
	# 编译 user
	go build -o ./target/user ./microservices/user/main.go
	# 编译 device
	go build -o ./target/mqtt ./microservices/device/internal/cmd/mqtt/mqtt.go
	go build -o ./target/device ./microservices/device/main.go
	# 编译 work
	go build -o ./target/weather ./microservices/work/internal/cmd/weather/weather.go
	go build -o ./target/work ./microservices/work/main.go

stop:
	#终止 nms
	pkill nms || exit 0
	#终止 user
	pkill user || exit 0
	#终止 device
	pkill mqtt || exit 0
	pkill device || exit 0
	#终止 work
	pkill weather || exit 0
	pkill work || exit 0

clean:
	rm ./log/*.log