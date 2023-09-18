HOST_DEV="192.168.82.186"
USER_DEV="root"
BINARY=bin/linux/stunel
linux:
	GOOS=linux GOARCH=amd64 go build -o ${BINARY}
	bin/linux/stunel
macosx:
	GOOS=darwin GOARCH=amd64 go build -o bin/macosx/stunel
	bin/macosx/stunel

windows:
	GOOS=windows GOARCH=amd64 go build -o bin/windows/stunel
	bin/windows/stunel

run:
	go run main.go
td: 
	go mod tidy

# test:
# 	curl 'http://127.0.0.1:80/api?api=123&msg=$(COMMIT)'

dp:
	make linux
	ssh ${USER_DEV}@${HOST_DEV} "cd /home/stunnel && rm stunnel"
	scp ${BINARY} ${USER_DEV}@${HOST_DEV}:/home/stunnel/stunnel
	ssh ${USER_DEV}@${HOST_DEV} "cd /home/stunnel && chmod +x stunnel"
	ssh ${USER_DEV}@${HOST_DEV} "cd /home/stunnel && sh start.sh"