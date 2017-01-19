echo %~dp0
set GOPATH=%GOPATH%;%~dp0
go build -o dist/win/buildit.exe src/run/main.go
