build:
	GOOS=linux GOARCH=arm GOARM=6 go build -ldflags '-s -w' -o teleinfo_exporter_armv6l *.go
	GOOS=linux GOARCH=amd64 go build -ldflags '-s -w' -o teleinfo_exporter_amd64 *.go

run:
	go run *.go

builddebug:
	GOOS=linux GOARCH=arm GOARM=6 go build -o teleinfo_exporter_armv6l *.go
	GOOS=linux GOARCH=amd64 go build -o teleinfo_exporter_amd64 *.go

debug: builddebug
	scp teleinfo_exporter_armv6l 192.168.29.250:/tmp/.
	ssh 192.168.29.250 /tmp/teleinfo_exporter_arm
