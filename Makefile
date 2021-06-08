passhider.exe: .
	go build -ldflags="-s -w" -o=$@ $<