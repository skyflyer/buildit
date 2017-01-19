echo %~dp0
set GOPATH=%GOPATH%;%~dp0
go build -o dist/buildit.exe src/run/main.go