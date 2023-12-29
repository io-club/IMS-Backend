MODE? = debug

prepare:
	mkdir -p ./log/nms
	mkdir -p ./log/user
	mkdir -p ./log/device
	mkdir -p ./target

run: prepare build
	@echo "Start service"
	MODE=${MODE} nohup ./target/device &
	MODE=${MODE} nohup ./target/user &
	MODE=${MODE} nohup ./target/work &
	#MODE=${MODE} nohup ./target/nms &
	# MODE=${MODE} ./target/user | tee -a ./log/user/user.log
	# MODE=${MODE} ./target/nms
	@echo "Service has started"

run-fresh: build stop run

build:
	go build -o ./target/nms ./main.go
	go build -o ./target/user ./microservices/user/main.go
	go build -o ./target/device ./microservices/device/main.go
	go build -o ./target/work ./microservices/work/main.go

stop:
	pkill nms || exit 0
	pkill user || exit 0
	pkill device || exit 0

clean:
	rm ./log/*.log