version=0.1.2
check_defined = \
    $(strip $(foreach 1,$1, \
        $(call __check_defined,$1,$(strip $(value 2)))))
__check_defined = \
    $(if $(value $1),, \
      $(error Undefined $1$(if $2, ($2))))

install:
	@command -v go >/dev/null 2>&1 || { echo >&2 "Please install go. Aborting."; exit 1; }
	@command -v dep >/dev/null 2>&1 || { echo >&2 "Please install dep. Aborting."; exit 1; } # for the future
	dep ensure
	go build -i -v -o eplc-$(version) src/eplc.go
	go test
	mkdir target
	mkdir target/bin
	mv eplc-$(version) target/bin
	@sudo mv target/bin/eplc-$(version) /bin/

clean:
	rm -rf target
	sudo rm -rf /bin/eplc-$(version)	 

list:
	ls /bin/epl-*
