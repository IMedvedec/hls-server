test-hls:
	go test -c ./hls && ./hls.test -test.v && rm hls.test