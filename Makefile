release:
	GOOS=windows GOARCH=amd64 go build -o ./bin/awsping_windows_amd64
	GOOS=linux GOARCH=amd64 go build -o ./bin/awsping_linux_amd64
	GOOS=darwin GOARCH=amd64 go build -o ./bin/awsping_darwin_amd64

push:
	setup/publish.sh
