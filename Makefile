check_defined = \
    $(strip $(foreach 1,$1, \
        $(call __check_defined,$1,$(strip $(value 2)))))
__check_defined = \
    $(if $(value $1),, \
      $(error Undefined $1$(if $2, ($2))))
install:

	@:$(call check_defined, GOPATH, please set it up and then run make)
	@:$(call check_defined, GOBIN, please set it up and then run make)
	@:$(call check_defined, GOROOT, please set it up and then run make)
	@command -v go >/dev/null 2>&1 || { echo >&2 "Please install go. Aborting."; exit 1; }
	@command -v dep >/dev/null 2>&1 || { echo >&2 "Please install dep. Aborting."; exit 1; } # for the future
	@go build src/libapl/aplc.go
	@sudo cp aplc /bin/aplc
