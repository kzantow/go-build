.PHONY: format
format:
	@go run -C build . format
.PHONY: lint-fix
lint-fix:
	@go run -C build . lint-fix
.PHONY: static-analysis
static-analysis:
	@go run -C build . static-analysis
.PHONY: unit
unit:
	@go run -C build . unit
.PHONY: test
test:
	@go run -C build . test
.PHONY: help
help:
	@go run -C build . help
.PHONY: makefile
makefile:
	@go run -C build . makefile
.PHONY: *
.DEFAULT:
	@go run -C build . $@
