#!/bin/sh

echo "build app..."

Server_PATH=D:\\GO\\src\\collector\\cmd\\server
Client_PATH=D:\\GO\\src\\collector\\cmd\\client

cd $Server_PATH
go build -ldflags="-w -s -H windowsgui" -o server-collect.exe
upx -9 server-collect.exe

cd $Client_PATH
go build -ldflags="-w -s -H windowsgui" -o client-collect.exe
upx -9 client-collect.exe