#!/bin/sh

echo "build app..."

Server_PATH=D:\\GO\\src\\collector\\cmd\\server
Client_PATH=D:\\GO\\src\\collector\\cmd\\client

set GOARCH=386
cd $Server_PATH

go build -o server-collect.exe -ldflags="-w -s -H windowsgui"  main.go
upx -9 server-collect.exe

cd $Client_PATH

go build -o client-collect.exe -ldflags="-w -s -H windowsgui"  main.go
upx -9 client-collect.exe

