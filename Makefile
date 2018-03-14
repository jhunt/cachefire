deploy:
	GOOS=linux GOARCH=amd64 go build
	cf push --no-start
