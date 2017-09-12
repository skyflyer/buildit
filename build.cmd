echo %~dp0
set GOPATH=%GOPATH%;%~dp0
go build -o dist/win/buildit-win.exe src/run/main.go
