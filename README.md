# eplc compiler 
eplc is a compiler frontend for the epl programming language

# build
* git clone https://github.com/e3expeli0t/eplc
* cd eplc
* make build
which will create new eplc binary in target/bin/

# run 
* to compile into air code simply run eplc [options] filename

# make commands
|command|description|
|--|--|
|sync|syncs the local branch with the online branch|
|switch|switch branch|
|update|commit and push the changes in the local branch|
|build|builds the project|
|devel_tests|run tests and show how much of the code they cover|
|clean|clean build files|
|list|list all eplc binaries installed on the system|
|install|install eplc|

# disclaimer
The project code will not be at any stage a production code.
This project is simply a large amount of tests, ideas and experiments that I do
while I'm learning the basics of PLT
also the project writing will not progress fast mainly because I don't have a lot of time to work on it.

# The language architecture

The language tools (e.g compiler) are divided into 2 part:
* frontend
* backed

the frontend is handled by eplc (epl compiler) which includes the parser the lexer the type checker and the AIR generator.
the backend is handled by avm which includes the aeu (AIR execution unit) the optimizer the target language generator
and optimizer.

## backend
the backend is divided into two parts:
* air execution unit (aeu)
* target language toolchain (aka tlt)

### Air execution unit (AEU)
AEU is the main part of the avm
the aeu is responsible for handling all the things that relate to the AIR 
meaning is responsible to handle the parsing, the lexing and the optimizing of the AVM IR 
and is responsible tor handle the execution of the code (e.g interpret him) if needed. 
The aeu can be used as a wrapper to epl code meaning you can "inject" the epl code into the aeu 
the when you run the avm binary only the code that you injected will run (injecting code into the aey will create new avm binary)

### Target language toolchain (TLT)
The tlt is a library (written in go) that responsible for generating target working machine code, which includes:
* generating suitable RT 
* optimizing the target machine code 
* linking external libraries
* specific OS optimization