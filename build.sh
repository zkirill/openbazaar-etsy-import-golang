env GOOS=windows GOARCH=amd64 go build -v -o import_etsy_ob_windows_amd64.exe
env GOOS=windows GOARCH=386 go build -v -o import_etsy_ob_windows.exe
env GOOS=darwin GOARCH=amd64 go build -v -o import_etsy_ob_darwin
