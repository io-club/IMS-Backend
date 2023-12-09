MODE? = debug

prepare:
	mkdir -p ./log/nms
	mkdir -p ./log/user
	mkdir -p ./log/device
	mkdir -p ./target

run: prepare build
	@echo "Running service"
    MODE=${MODE} nohup ./target/device >> ./log/device/device.log 2>&1 &
    MODE=${MODE} nohup ./target/user >> ./log/user/user.log 2>&1 &
    #MODE=${MODE} nohup ./target/nms >> ./log/nms/nms.log 2>&1 &
    #MODE=${MODE} nohup ./target/user | tee -a ./log/user/user.log
    MODE=${MODE} ./target/nms | tee -a ./log/nms/nms.log

run-fresh: build stop run

build:
	go build -o ./target/nms ./main.go
	go build -o ./target/user ./microservices/user/main.go
	go build -o ./target/device ./microservices/device/main.go

stop:
	pkill nms || exit 0
	pkill user || exit 0
	pkill device || exit 0

clean:
	rm ./log/*.log