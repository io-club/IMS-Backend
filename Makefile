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
	MODE=${MODE} ./target/nms >> | tee -a ./log/nms/nms.log
	@echo "Finish service"

run-fresh: build stop run

build:
	go build -o ./target/nms ./internal/nms.go
	go build -o ./target/user ./internal/user/cmd/user.go
	go build -o ./target/device ./internal/device/cmd/device.go

stop:
	pkill nms || exit 0
	pkill user || exit 0
	pkill device || exit 0

clean:
	rm ./log/*.log