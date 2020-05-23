version=0.2.0
branch=devel
msg=bug fixes
exec=samples/test.epl

check_defined = \
    $(strip $(foreach 1,$1, \
        $(call __check_defined,$1,$(strip $(value 2)))))

__check_defined = \
    $(if $(value $1),, \
      $(error Undefined $1$(if $2, ($2))))

build:clean
	@command -v go >/dev/null 2>&1 || { echo >&2 "Please install go. Aborting."; exit 1; }
	@command -v dep >/dev/null 2>&1 || { echo >&2 "Please install dep. Aborting."; exit 1; }
	@echo Building eplc...
	go build -i -v -o eplc src/eplc.go
	@mkdir target
	@mkdir target/bin
	@mv eplc target/bin

buid_support:
	@echo Building support tools
	cd tools/Support/epldbg/; cmake .
	cd tools/Support/epldbg/; make

rebuild:
	@echo Running tests...
	go test -v eplc/tests -cover
	
	@echo Rebuilding eplc...
	go build -i -v -o eplc src/eplc.go
	@rm -rf target/bin/eplc
	@mv eplc target/bin/

install:build
	@echo "Installing..."
	@sudo mv target/bin/eplc /bin/
	@echo "Finished..."

run:rebuild
	./target/bin/eplc $(exec)

clean:
	@echo Removing eplc targets...
	@rm -rf target
	@echo Removing Support targets...
	@cd tools/Support/epldbg/; make clean
	@rm -rf samples/*.txt samples/syntax_errors/*.txt

devel_tests:
	dep ensure $(dep_args)
	go test -v eplc/tests -covermode=count -coverprofile=count.out fmt
	go tool cover -html=count.out

list:
	@ls  /bin/epl*

commit:clean
	@git add .
	@git commit -a -m "$(msg)"
	@git push

pull:
	@git pull

checkout:
	git checkout $(branch) 
