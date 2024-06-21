.PHONY: format
format:
	@go run -C cmd . format
.PHONY: lint-fix
lint-fix:
	@go run -C cmd . lint-fix
.PHONY: help
help:
	@go run -C cmd . help
.PHONY: makefile
makefile:
	@go run -C cmd . makefile
