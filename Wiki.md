#The epl programming language

##Contents:
1. introduction
2. language semantics
3. runtime
4. sub-languages
5. optimization
6. aliases
7. avm
8. packages and libraries
9. What next?

#Introduction:

Epl is multi paradigmatic . That focuses on easy complex software systems
production. The language is easy to learn and code.

#Language semantics:
1. Variables and types
2. Functions
3. If –else
4. Loops
5. Cfalgs
6. Memory
7. Functions groups
8. Air
9. Next updates

##Variables and types:
###types:

Types are names for differ data types. For example if we want to represent an
Integer we will use the data type 'int' and if we want to represent a floating
point number we will use the type 'float'.

###List of types and their meaning:

* int – system dependent sized integer
* int8 – 8 bit integer
* int16 – 16 bit integer
* int32 – 32 bit integer
* int64 – 64 bit integer
* float16 – 16 bit floating point numebr
* float32 – 32 bit floating point numebr
* float64 – 64 bit floating point number
* cmx – complex number
* long – 128 bits integer
* float – system dependent sized floating point number

###Variables:
Variable is a name that have a certain value of certain type.
To declare variable we use the saved word 'decl' followed by the variable name
and is type. For example to declare a variable of the type int we write:

```
#!epl

decl a int;
```

to initialize an already declared variable we write:

```
#!epl

a = value;
```
were value is an integer.
to declare variable and initialize him we write:

```
#!epl

a type = value;
```

were type is the type of the data that the variable is holding and value is the data
that the variable is holding

####arithmatics:
we can add subtract multiply and devide variable for example:

```
#!epl

a int = 9;
b int = (a*2)-2/2+3
```

here we initialize the variable a to the value of 9
after that we initialize b to have the value of (9*2) - 1 +3 i.e 20.

Epl is also supports the modulo (aka reminder) operator (%). Which can be used like
that (a is variable with the type int and the value 8) 
```
#!epl

b int = a % 2;
```
(the value of b is 0
because a is devided by 2)

##Functions:
A function is a grouped sequence of operations that can be called over and over.
To declare a function we use the reserved word fnc.
Example:

```
#!epl

fnc hello_world();
```

here we declare a function that’s called hello_world(). To use ( aka implement ) the
function after we declared her we write:

```
#!epl


fnc hello_world(): void {
    out.put("Hello World");
}
```

This function wont compile, first we need to "import" the group (will be explaind
later) out. 
To import a group we will use the reserved word import.
Example:

```
#!epl

Import out;

```
The complete code:

```
#!epl

import out;

@MainFunc
fnc hello_world() {
   out.put("Hello world");
}
```


To compile save this code as hello.epl and write *eplc –run hello.epl*.
Function also can have an argumets. Argument is some variable that’s the function
gets and do some processing on him. To declare function with arguments a and b of
the type int we write 
```
#!epl

fnc func(a int, b int);
```

We can use the variable a and inside the function func.

###Functions return types:
As in variables function can also have values, in this case the value of the function is
the value that the function returns. For example suppose we want to create a
function that’s computes the sum of to numbers x and y. To create such a function
we write:
```
#!epl


fnc sum(x int, y int): int;
```

which means declare a function with 2 arguments x and y of the type int and return
value of the type int.
the sum of x and y is x+y. To return a value we use the keyword return.
The complete code:

```
#!epl

import out;

fnc sum (x int, y int): int {
   return x+y;
}
```

###Main function:
In order for epl to recognize the start function of a program. We call the start
function main . But some times we don’t wont to do that. So epl implements
something that called compiler flags ( this topic will be coverd in much more detail in
the rest of this article) to call a compiler flag we use the symbol @ flowed by the
name of the flag. For example @MainFunc which will tell the compiler that the
function is the main function.

##If-else:
If and else statements checks if Boolean operation is returning true or false

Example:

```
#!epl

if (232 > 90)
{
   out.put("Wowww 232 is bigger then 90");
}
else
{
   Out.put("WTH?!?!? 232 is bigger then 90. The world is broken");
}
```


The stmt prints Wowww 232 is bigger then 90.

##Loops:
 ```
to be continue
```