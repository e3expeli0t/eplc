version=0.0.1
branch=devel
msg=bug fixing
dep_args=
exec=test.eplc.nvis
check_defined = \
    $(strip $(foreach 1,$1, \
        $(call __check_defined,$1,$(strip $(value 2)))))
__check_defined = \
    $(if $(value $1),, \
      $(error Undefined $1$(if $2, ($2))))

build:
	@#make clean
	@command -v go >/dev/null 2>&1 || { echo >&2 "Please install go. Aborting."; exit 1; }
	@command -v dep >/dev/null 2>&1 || { echo >&2 "Please install dep. Aborting."; exit 1; } # for the future
	
	@echo Checking dependencies...
	dep ensure $(dep_args)
	
	@echo Running tests...
	go test -v eplc/src/libepl/epllex -cover 
	
	@echo Building eplc...
	go build -i -v -o eplc-$(version) src/eplc.go
	@mkdir target
	@mkdir target/bin
	@mv eplc-$(version) target/bin

	@echo Building support tools
	cd tools/Support/epldbg/; cmake .
	cd tools/Support/epldbg/; make

rebuild:
	@echo Running tests...
	go test -v eplc/src/libepl/epllex -cover
	
	@echo Rebuilding eplc...
	go build -i -v -o eplc-$(version) src/eplc.go
	@rm -rf target/bin/eplc-$(version)
	@mv eplc-$(version) target/bin/

install:
	make build
	sudo install target/bin/eplc-$(version) /bin/
run: 
	make rebuild
	./target/bin/eplc-$(version) $(exec)
clean:
	@echo Removing eplc targets...
	@rm -rf target
	@sudo rm -rf /bin/eplc-$(version)	 

	@echo Removing Support targets...
	@cd tools/Support/epldbg/; make clean	
devel_tests:
	dep ensure $(dep_args)
	go test -v eplc/src/libepl/epllex -covermode=count -coverprofile=count.out fmt
	go tool cover -html=count.out
list:
	@ls  /bin/epl*

update:clean
	@rm -rf vendor/*
	@git add .
	@git commit -a -m "$(msg)"
	@git push
sync:
	@git pull

switch:
	git checkout $(branch) 
