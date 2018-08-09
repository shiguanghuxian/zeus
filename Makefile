all: 
	go build -o ./apps/bin/zues-dispatchd ./apps/zues-dispatchd
	go build -o ./apps/bin/zues-portal ./apps/zues-portal
	go build -o ./apps/bin/zues-serverd ./apps/zues-serverd
	go build -o ./apps/bin/zues-statisd ./apps/zues-statisd
clean: 
	rm -f ./apps/bin/zues-dispatchd
	rm -f ./apps/bin/zues-portal
	rm -f ./apps/bin/zues-serverd
	rm -f ./apps/bin/zues-statisd
	rm -f ./apps/bin/logs/*.log
	