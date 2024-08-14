@echo off
go env -w CGO_ENABLED=1
go env -w GOOS=windows
go env -w GOARCH=386
go env -w CC=C:\TDM-GCC-32\bin\mingw32-gcc.exe
go env -w CXX=C:\TDM-GCC-32\bin\mingw32-c++.exe
go build -ldflags "-w -s -H=windowsgui" -tags 'release'
pause