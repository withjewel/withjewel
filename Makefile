fmt:
	go fmt ./...
	@echo Done
watchcss:
	sass --style compact --sourcemap=none --no-cache --watch static/sass:static/css
