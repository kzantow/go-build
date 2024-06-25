.PHONY: $(MAKECMDGOALS)
$(MAKECMDGOALS):
	@go run -C build . $@

.PHONY: build
.PHONY: *
.DEFAULT:
	@go run -C build . $@
