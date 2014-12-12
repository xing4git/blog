```
def max(x:Int, y:Int):Int = {
	if(x>y) x else y
}
```

Sometimes the Scala compiler will require you to specify the result type of a function. If the function is recursive, for example, you must explicitly specify the function's result type. In the case of `max` however, you may leave the result type off and the compiler will inter it. Also, if a function consists of just one statement, you can optionally leave off the curly braces.

```
def max(x:Int, y:Int) = if(x>y) x else y
```

Nevertheless, it is often a good idea to indicate function result types explicitly, even when the compiler doesnot require it. Such type annotations can make the code easier to read, because the reader need not study the function body to figure out the inferred result type.

----
Java's ++i and i++ donot work in Scala. To increment in Scala, you need to say either i=i+1 or i+=1.

----
One way to print each command line argument is:
	
	args.foreach(arg => println(arg))

In this code, you call the foreach method on args, and pass in a function. In this case, you are passing in a function literal that takes one parameter named `arg`. 

In the previous example, the Scala interpreter inters the type of `arg` to be String, since String is the element type of the array on which you are calling foreach. If you would prefer to be more explicit, you can mention the type name, but when you do you will need to wrap the argument portion in parentheses.

	args.foreach((arg: String) => println(arg))
	
If a function literal consists of one statement that takes a single argument, you need not explicitly name and specify the argument.
	
	args.foreach(println)
	
To summarize, the syntax for a function literal is a list of named parameters, in parentheses, a right arrow, and then the body of the function.

----
Another way to print each command line argument is:

	for(arg <- args) 
		println(arg)

To the right of the <- symbol is the familiar args array. To the left <- is "arg", the name of a val, not a var. (Because it is always a val, you just write "arg" by itself, not "val arg"). Although arg may seem to be a var, because it will get a new value on each iteration, it really is a val: arg cannot be reassigned inside the body of the for expression. Instead, for each element of the args array, a new arg val will be created and initialized to the element value, and the body of the for will be executed.

----
When you define a variable with val, the variable cannot be reassigned, but the object it refers could potentially still be changed.

	val greets = new Array[String](3)
	greets(0) = "Hello"
	greets(1) = ","
	greets(2) = "World!"
	
So in this case, you could not reassign greets to a different array, greets will always point to the same Array[String] instance with which it was initialized. But you can change the elements of the Array[String] over time, so the array itself is mutable.

----
If a method(not a function) takes only one parameter, you can call it without a dot or parentheses.

	0.to(2) // equals
	0 to 2
	
Note that this syntax only works if you explicitly specify the receiver of the method call. You can not write `println 10`, but you can write `Console println 10`.

---
Scala does not technically have operator overloading, because it does not actually have operators in the traditional sense. Instead, characters such as +,-,*,/ can be used in method names. Thus, when you typed `1+2`, you were actually `invoking a method named + on the Int object 1, passing in 2 as a parameter`. You could alternatively have written `1+2` using traditional method invocation syntax `(1).+(2)`.

----
In Scala, you can access an element by specifying an index in parentheses. So the first element in Scala array named args is args(0), not args[0], as in Java.

In Scala, when you apply parentheses surrouding one or more values to a variable, Scala will transform the code into an invocation of a method named `apply` on that variable. So `greets(0)` gets transformed into `greets.apply(0)`. Thus accessing an element of an array in Scala is simply a method call like any other.

This principle is not restricted to arrays: any application of an object to some arguments in parentheses will be transformed to an `apply` method call. Of course this will compile only if that type of object actually defines an `apply` method. So it is not a special case; it is a general rule.

Similarly, when an assignment is made to a variable to which parentheses and one or more arguments have been applied, the compiler will transform that into an invocation of an `update` method that takes that arguments in parentheses as well as the object to the right of equals sign. For example, `greets(0) = "hello"` will be transformed into: `greets.update(0, "hello")`.

----
As you have seen, a Scala array is a mutable sequence of objects that all share the same type. An Array[String] contains only strings, although you can not change the length of an array after it is instantiated, you can change its element values. Thus, arrays are mutable objects. 

For an immutable sequence of objects that share the same type you can use scala.List class. As with arrays, a List[String] contains only strings.

	val numbers = List(1, 2, 3) // List.apply(1, 2, 3)
	
You can use `::` operator to prepends a new element to the beginning of an existing list, and returns the resulting list.

	val numbers = List(2, 3)
	val oneTwoThree = 1::numbers // oneTwoThree = List(1, 2, 3)
	
You can also use `:::` for list concatenation:

	val numbers = List(1, 2) ::: List(3, 4) // List(1, 2, 3, 4)
	
Note that in the expression `1::numbers`, `::` is a method of its right operand, the `numbers`. There is a simple rule to remember: If a method is used in operator notation, such as `a*b`, the method is invoked on the left operand, unless the method name ends in a colon. If the method name ends in a colon, the method is invoked on the right operand. Therefore, in `1::numbers` the `::` method is invoked on `numbers`, passing `1`.

Given that a shorthand way to specify an empty list is Nil, one way to initialize new lists is to string together elements with the `::` operator, with Nil as the last element.

	val numbers = 1::2::3::Nil // List(1, 2, 3)

----
Like lists, tuples are immutable, but unlike lists, tuples can contain different types of elements.

	var pair = (99, "Hello")
	println(pair._1)
	println(pair._2)
	
To instantiate a new tuple that holds some objects, just place the objects in parentheses, separated by commons. Once you have a tuple instantiated, you can access its elements individually with a dot, underscore, and the `one-based` index of the element.

The actual type of a tuple depends on the number of elements it contains and the types of those elements. Thus, the type of (99,"jell") is Tuple2[Int,String]; The type of ('u', 'r', "the", 1, 4, "me") is Tuple6[Char, Char, String, Int, Int, String].

---
Scala API contains a base trait for sets, and also provides two subtraits, one for mutable sets and another for immutable sets. These three traits all share the same simple name, Set. There full qualified names differ, however, because each resides in a different package.

Concrete set classes in the Scala API, such as HashSet classes extend either the mutable or immutable Set trait. Thus, if you want to use a HashSet, you can choose between mutable and immutable varieties depending upon your needs.

![set继承体系](scala.note.set.png)

	var jetSet = Set("Boeing", "Airbus")
	jetSet += "Lear"
	println(jetSet.contains("Cessna"))
	
To add a new element to a set, you call `+` on the set, passing in the new element. Both mutable and immutable sets offer a `+` method, but their behavior differs. Whereas a mutable set will add the element to itself, an immutable set will create and return a new set with the element added.

---
As with sets, Scala provides mutable and immutable versions of Map, using a class hierarchy.

	import scala.collection.mutable.Map
	
	val treasureMap = Map[Int, String]()
	treasureMap += (1->"Go to island")
	treasureMap += (2->"Find big x on ground")
	treasureMap += (3->"Dig")
	println(treasureMap(2))
	
This `->` method, witch you can invoke on any object in a Scala program, returns a two-element tuple containing the key and value.

	val romanMap = Map(
		1->"I", 2->"II", 3->"III"
	)
	
	
----
In Scala, `public` is the default access level.

Method parameters in Scala are vals, not vars.

The Scala compiler treats a function defined in the procedure style, with curly braces but no equals sign, essentially the same as a function that explicitly declares its result type to be Unit:

	def g() { "this string gets lost" }

The function `g` returns Unit. This is true no matter what the body contains, because the Scala compiler can convert any type to Unit. For example, if the last result of a method is a String, but the method result type is declared to be Unit, the String will be converted to Unit and its value lost.

If you intend to return a non-Unit value without explicitly declare result type, you need to insert the equals sign before the function body:

	def h() = {"this String gets returned" }
	
The function `h` returns String.

<<<<<<< HEAD
-----
Prefer vals, immutable objects, and methods without side effects. Reach for them first. Use vars, mutable objects, and moethods with side effects when you have a specific need and justification for them.

---
The rules of semicolon inference

A line ending is treated as a semicolon unless one of the following conditions is ture:

* The line in question ends in a word that would not be legal as the end of a statement, such as a period or an infix operator.
* The next line begins with a word that cannot start a statement.
* The line ends while inside parentheses (…) or brackets […], because these cannot contain multiple statements anyway.

---
Singleton objects

A singleton object definition looks like a class definition, except instead of the keyword `class` you use the keyword `object`.

```
object ChecksumAccumulator {
	private val cache = Map[String, Int]()
	
	def calculate(s: String): Int = {
		if(cache.contains(s)) 
			cache(s)
		else {
			val acc = new ChecksumAccumulator
			val cs = acc.calculate(s)
			cache += (s->cs)
			cs
		}
	}
}
```

When a singleton object shares the same name with a class, it is called that class's `companion object`. The class is called the `companion class` of the singleton object. A class and its companion object can access each other's private members, and you must define both the class and its companion object in the same source file.

If you are a Java programmer, one way to think of singleton objects is as the home for static methods. You can invoke methods on singleton objects using a similar syntax.

```
ChecksumAccumulator.calculate
```

One difference between classes and singleton objects is that singleton objects cannot take parameters, whereas classes can. Because you can't instantiate a singleton object with the new keyword, you have no way to pass parameters to it. Each singleton object is implemented as an instance of a synthetic class referenced from a static variable, so they have the same initializtion semantics as Java statics. In particular, a singleton object is initialized the first time some code accesses it.

A singleton obeject that does not share the same name with a companion class is called a standalone object. You can use standalone objects for many purposes, including collecting related utility methods together, or defining an entry point to a Scala application.


---
A Scala application

To run a Scala program, you must supply the name of a standalone singleton object with a main method that takes one parameter, an Array[String], and has a result type of Unit. Any standalone object with a main method of the proper signature can be used as the entry point into an application.

```
object Summer {
	def main(args: Array[String]) {
		for(arg <- args) println(arg)
	}
}
```

Scala provides a trait, scala.Application:

```
object Summer extends Application {
	for(arg <- args) println(arg)
}
```

To use the trait, you first write `extends Application` after the name of your singleton object. Then instead of writing a main method, you place the code you would have put in the main method directly between the curly braces of the singleton object.

---
Scala implicitly imports members of packages `java.lang` and `scala`, as well as the members of a singleton object named Predef, into every Scala source file. Predef, which resides in package `scala`, contains many useful methods. For example, when you say println in a Scala source file, you are actually invoking println on Predef. Predef.println turns around and invokes Console.println, which does the real work. When you say assert, you are invoking Predef.assert.


---
If an integer literal ends in an L or l, it is a Long, otherwise it is an Int.

If an Int literal is assigned to a variable of type Short or Byte, the literal is treated as if it were a Short or Byte type so long as the literal value is within the valid range for that type.

```
val prog = 111L // Long
val tower = 11 // Int
val little: Byte = 38 // Byte
```

If a floating-point literal ends in an F or f, it is a Float, otherwise it is a Double.

Scala includes a special syntax for raw strings. You start and end a raw string with three double quotation marks. The interior of a raw string may contain any chars, including newlines, quotation marks, and special characters.

```
println("""|Welcome to Ultamix 3000.             |Type "HELP" for help.""")
```


---
Most operators are infix operators, which mean the method to invoke sits between, the object and parameter or parameters you wish to pass to the method. Scala also has two other operator notations: prefix and postfix. In prefix notation, you put the method name before the object, for example, `-7`. In postfix notation, you put the method after the object, for example `7 toLong`.

As with the infix operators, prefix operators are a shorthand way of invoking methods. In this case, however, the name of the method has "unary_" prepended to the operator character. For instance, Scala will transform the expression `-2.0` into the method invocation `(2.0).unary_-`. 

Postfix operators are methods that take no arguments, when they are invoked without a dot or parentheses. In Scala, you can leave off empty parentheses on method calls. The convention is that you include parentheses if the method has side effects, such as println(), but you can leave them off if the method has no side effects, such as toLowerCase invoked on a String. In this case of method that requires no arguments, you can alternatively leave off the dot and use postfix operator notation.
=======
---
Prefer vals, immutable objects, and methods without side effects. Reach for them first. Use vars, mutable objects, and moethods with side effects when you have a specific need and justification for them.

---
If you want to compare two objects for equality, you can use either ==, or !=. == has been carefully crafted so that you get just the equality comparison you want in most cases. First check the left side for null, and if it is not null, call the equals method. Since equals is a method, the precise comparison you get depends on the type of the left-hand argument. Since there is an automatic null check, you donot have to do the check yourself. The automatic does not look at the right-hand side, but any reasonable equals method should return false if its argument is null.
>>>>>>> fa004fdc3394401077c95432c15e7df8a3cab039

In Java, you can use == to compare both primitive and reference types. On primitive types, Java's == compares value equality, as in Scala. On reference types, however, Java's == compares reference equality, which means the two variables point to the same object on the JVM's heap. Scala provides a facility for comparing reference equality, as well, under the name `eq`, or `ne`.

<<<<<<< HEAD
=======
---
You can invoke many methods on Scala's basic types, these methods are available via implicit conversions. For each basic type, there is also a rich wrapper that provides several additional methods.

---
Implicit conversions

```
class Rational(n: Int, d: Int) {
    require(d!=0)
    private val g = gcd(n.abs, d.abs)
    val number = n/g
    val denom = d/g

    def this(n:Int) = this(n, 1)

    def +(that:Rational):Rational = new Rational(number*that.denom + that.number*denom, denom*that.denom)
    
    def +(i:Int):Rational = new Rational(number+denom*i, denom)
}
```

For the Rational object r, you can write r+2, but you cannot write 2+r. Because there is not an add method accepts an Rational object on the Int 2.

There is a way to solve this problem in Scala: You can create an implicit conversion that automatically converts integers to rational numbers when needed:

```
implicit def intToRational(x: Int) = new Rational(x)
```

The implicit modifier in front of the method tells the compiler to apply it automatically in a number of situations. With the conversion defined, you can now write 2+r.

Implicit conversions are a very powerful technique for making libraries more fleible and more convenient to use. But they can also be easily misused.

---
Almost all of Scala's control structures result in some value. Programmers can use these result value to simplify their code, just as they use return values of functions.

```
val filename = if(!args.isEmpty) args(0) else "default.txt"
```

This code is slightly shorter, but its real advantage is that it uses a val instead of a var. Using a val is the functional style, and it helps you in much the same way as a final variable in Java.

The while and do-while constructs are called "loops", not expressions, because they donot result in an interesting value. The type of the result is Unit.

One other construct that results in the unit value, is reassignment to vars.

```
var line = ""
while((line=readLine())!="") {
	//
}
```

The code above doesnot work in Scala. Whereas in Java, assignment results in the value assigned, but in Scala assignment always results in the unit value.

---
For expression

The simplest thing you can do with for is to iterate through all the elements of a collection.

```
val files = (new java.io.File("/Users")).listFiles
for(file <- files) println(file)
```

Sometimes you want to filter a collection down to some subset. You can do this with a for expression by adding a filter: an if clause inside the for's parentheses.

```
for(file <- files if file.getName.endsWith(".scala")) println(file)
```

You can include more filters if you want, just keep adding if clauses.

You use for to iterate values and then forgotten them so far, you can also generate a value to remember for each iteration. To do so, you prefix the body of the for expression by the keyword yield.

```
def hiddens = for{
    file <- files if file.getName.contains(".")
} yield file
```

When the for expression completes, the result will include all the yielded values contained in a single collection. The type of the resulting collection is based on the kind of collections processed in the iteration clauses. The syntax of a for-yield expression is like this: `for clauses yield body`. The yield goes before the entire body. Even if the body is a block surrounded by curly braces, put the yield before the first curly brace, not before the last expression of the block.

---
try-catch-finally

In Scala, try-catch-finally results in a value. The result is that of the try clause if no exception is thrown, or the relevant catch clause if an exception is thrown and caught. If an exception is thrown but not caught, the expression has no value at all. The value computed in the finally clause, if there is one, is dropped. Usually finally clauses do some kind of clean up such as closing a file; they should not normally change the value computed in the main body or a catch clause fo the try.

```
def urlFor(path: String) = 
	try {
		new URL(path)
	} catch {
		case e: MalformedURLException => new URL("http://www.scala-lang.org")
	}
```

---
match expression

Scala's match expression lets you select from a number of alternatives, just like switch statements in other languages.

```
val firstArg = if(args.length>0) args(0) else ""

firstArg match {
	case "salt" => println("papper")
	case "chips" => println("salsa")
	case _ => println("huh?")
}
```

There are a few important differences from Java's switch statement. One is that any kind of constant, as well as other things, can be used in cases in Scala, not just the integer-type and enum constants of Java's case statements. Another difference is that there are no breaks at the end of the each alternative. Instead the break is implicit, and there is no fall through from one alternative to the next.

The match expressions also result in a value.

```
val firstArg = if (!args.isEmpty) args(0) else ""
val friend =
	firstArg match {
		case "salt" => "pepper"
		case "chips" => "salsa"
		case "eggs" => "bacon"
		case _ => "huh?"
	}
```

---
A funciton literal is compiled into a class that when instantiated at runtime is a function value. Thus the distinction between function literals and values is that function literals exist in the source code, whereas funciton values exist as object at runtime.

```
(x: Int) => x+1
```

The => designates that this function converts the thing on the left to the thing on the right. So, this is a function mapping any integer x to x+1.

Function values are objects, so you can store them in variables if you like. They are functions too, so you can invoke them using the usual parentheses function-call notation:

```
var increase = (x: Int) => x+1
increase(10) // 11
increase = (x: Int) => {
	x+999
}
increase(1) // 1000
```

Scala provides a number of ways to leave out redundant information and write function literals more briefly. One way to make a function literal more brief is to leave off the parameter types.

```
val someNums = List(-11, -10, -5, 0, 5, 10)
someNums.filter((x) => x>0) // List(5, 10)
```

In the previous example, the parentheses around x are unnecessary:

```
someNums.filter(x => x>0)
```

To make a function literal even more concise, you can use underscores as placeholders for one or more parameters, so long as each parameter appears only one time with the function literal.

```
someNums.filter(_>0)
```

You can think of the underscores as a "blank" in the expression that needs to be "filled in". This blank will be filled in with an argument to the function each time the function is invoked.

---
Partially applied functions

Scala treats `someNums.foreach(println _)` as if `someNums.foreach(x=>println(x)`. Thus the underscore in this case is not a placeholder for a single parameter. It is a placeholder for an entire parameter list. Remember that you need to leave a space between the function name and the underscore. When you use an underscore in this way, you are writing a partially applied functions. In Scala, when you invoke a function, passing in any needed arguments, you apply that function to the arguments.

```
def sum(a:Int, b:Int, c:Int) = a+b+c
sum(1, 2, 3) // apply the function sum to 1, 2, 3 arguments
```

A partially applied function is an expression in which you donot supply all of the arguments needed by the function. Instead, you supply some, or none, of the needed arguments. 

```
val a = sum _
a(1, 2, 3) // 6
```

Here is what just happened: The variable named a refers to a function value object. This function value is an instance of a class generated automatically by the Scala compiler from `sum _`, the partially applied function expression. The class generated by the compiler has an apply method that takes three arguments. The generated class's apply method takes 3 parameters because 3 is the number of arguments missing in the `sum _` expression. The Scala compiler translates the expression `a(1,2,3)` into an invocation of the function value's apply method, passing in the three arguments 1,2,3. Thus `a(1,2,3)` is a short form for: `a.apply(1,2,3)`.

This apply method, defined in the class generated automatically by the compiler from expression `sum _`, simply forwards those thress missing parameters to `sum`, and returns the result. 

You can also express a partially applied function by supplying some but not all of the required arguments.

```
val b = sum(1, _:Int, 3)
b(5) // 9
```

In this case, the middle argument is missing. Since only one parament is missing, the Scala compiler generates a new function class whose apply method takes one argument.

If you are writing a partially applied function expression in which you leave off all parameters, such as `println _`, or `sum _`, you can express it more concisely by leaving off the underscore if a function is required at that point in the code: `someNums.foreach(println)`. In situations where a funciton is not required, attempting to use this form will cause a compilation error.

---
Special function call forms

Scala allows you to indicate that the last parameter to a function may be repeated. This allows clients to pass variable length argument lists to the function. To denote a repeated parameter, place an asterisk after the type of the parameter.

```
def echo(args: String*) = for(arg<-args) println(arg)
```

Inside the function, the type of the repeated parameter is an Array of the declared type of the parameter. Nevertheless, if you have an array of the appropriate type, and you attempt to pass it as a repeated parameter, you will get a compiler error:

```
val arr = Array("1", "2")
echo(arr) //ERROR
```

To accomplish this, you will need to append the array argument with a colon and an _* symbol: `echo(arr: _*)`. This notation tells the compiler to pass each element of array as its own argument to echo, rather than all of it as a single argument.

Named arguments allow you to pass arguments to a function in a different order.

```
def speed(distance: Float, time: Float): Float = distance/time
speed(time=10, distance=100)
```

It is also possible to mix positional and named arguments. In that case, the positional arguments come first. Named arguments are most frequentlu used in combination with default paramter values.

Scala lets you specify default values for function parameters. The argument for such a parameter can optionally be omitted from a function call, in which case the corresponding argument will be filled in with the default.

```
def printTime(out: java.io.PrintStream = Console.out, divisor: Int = 1) = out.println("time = "+ System.currentTimeMillis()/divisor)
printTime(out=Console.err)
printTime(divisor=1000)
```

---
Tail recursion

```
def approximate(guess: Double): Double = 
	if(isGoodEnough(guess)) guess
	else approximate(improve(guess))
```

Note that the recursive call is the last thing that happens in the evaluation of function approximate's body. Functions like approximate, which call themselves as their last action, are called tail recursive. The Scala compiler detects tail recursion and replaces it with a jump back to the beginning of the funciton, after updating the function parameters with the new values.

A tail-recursive function will not build a new stack frame for each call; all calls will execute in a single frame. The use of tail recursive in Scala is fairly limited. Scala only optimizes directly recursive calls back to the same function making the call. If the recursion is indirect, no optimization is possible.

---
Curring function

A curried function is applied to multiple argument lists, instead of just one.

```
def curriedSum(x:Int)(y:Int) = x+y
curriedSum(1)(2) // 3
```

The first invocation takes a single Int parameter named x, and returns a function value for the second function. This second function takes the Int parameter y.

```
val onePlus = curriedSum(1)_
onePlus(2) // 3
```

---
In any Scala method invocation in which you are passing in exactly one argument, you can opt to use curly braces to surround the argument instead of parentheses.

```
println("hello")
println { "hello" }
```

The purpose of this ablility to substitute curly braces for parentheses for passing in one argument is to enable client programmers to write function literals between curly braces. This can make a method call feel more like a control abstraction.

```
def withPrintWriter(file: File)(op: PrintWriter=>Unit) {
	val writer = new PrintWriter(file)
	try {
		op(writer)
	} finally {
		writer.close()
	}
}

val file = new File("data.txt")
withPrintWriter(file) {
	writer => writer.println(new java.util.Date)
}
```

---
By-name parameters

What if there is no value to pass into the code between the curly braces? To help with such situation, Scala provides by-name parameters.

```
// Without using by-name parameters
var assertionEnable = true;
def myAssert(predicate: ()=>Boolean) =
	if(assertionEnable && !predicate()) throw new AssertionError
```

The definition is fine, but using it is a little bit awkward:
	
	myAssert(()=>5>3)

You would really prefer to leave out the empty parameter list and => symbol in the function literal and write the code like this:

	myAssert(5>3) // Won't work, because missing ()=>

By-name parameters exist precisely so that you can do this. To make a by-name parameter, you give the parameter a type starting with => instead of ()=>.

	def byNameAssert(predicate: =>Boolean) = 
		if(assertionEnable && !predicate) throw new AssertError

Now you can leave out the empty parameter in the property you want to assert. The result is that using byNameAssert looks exactly like using a built-in control structure:

	byNameAssert(5>3)


---
Parameterless methods

```
abstract class Element {
	def contents: Array[String]

	def height: Int = contents.length

	def width: Int = if(height==0) return 0 else return contents(0).length
}
```

Note that none of Element's three method has a parameter list, not even an empty one. Instead of:

	def width(): Int

The method is defined without parameter:

	def width: Int

Such parameterlesss methods are quite common in Scala. The recommended convention is to use a parameterless method whenever there are no parameters and the method accesses mutable state only by reading fields of the containing object(it doesn't change mutable state). 

In principle, it is possible to leave out all empty parentheses in Scala function calls. However, it is recommended to still write the empty parentheses when the invoked method represents more than a property of its receiver object. For instance, empty parentheses are appropriate if the method performs I/O, or writes reassignable vars, or reading vars other than the receiver's fields, either directly or indirectly by using mutable objects.

To summarize, it is encouraged style in Scala to define methods that take no parameters and have no side effects as parameterless methods. On the other hand, you should never define a method that has side-effects without parentheses, because then invocations of the method would like a field selection.


---
In Scala, fields and methods are in the same namespace. This makes it possible for a field to override a parameterless method.

```
abstract class Element {
	def contents: Array[String]
}

class ArrayElement(conts: Array[String]) extends Element {
	def contents: Array[String] = conts
}
```

Field contents in this version of ArrayElement is a perfectly good implementation of the parameterless method contents in class Element.

On the other hand, in Scala it is forbidden to define a field and method with the same name in the same class, whereas it is allowed in Java.


---

```
class ArrayElement (val contents: Array[String]) extends Element
```

Note that now the contents parameter is prefixed by val. This is a shorthand that defines at the same time a parameter and field with the same name. You can also prefix a class parameter with var, in which case the corresponding field would be reassignable. Finally, it is possible to add modifiers such as private, protected, or override to these parametric fields, just as you can do for any other class member.

```
class Cat {
	val dangerous = false
}

class Tiger(
	override val dangerous: Boolean, 
	private var age: Int
) extends Cat
```

---
Invoking superclass constructors

```
class LineElement(s: String) extends ArrayElement(Array(s)) {
	override def width = s.length
	override def height = 1
}
```

Since LineElement extends ArrayElement, and ArrayElement's constructor takes a parameter(an Array[String]), LineElement needs to pass an argument to the primary constructor of its superclass. To invoke a superclass constructor, you simply place the arguments you want to pass in parentheses following the name of the superclass.


---
Using override modifiers

Scala requires `override` for all members that override a concrete memeber in a parent class. The modifier is optional if a memeber implements an abstract memeber with same name. The modifier is forbidden if a member does not override or implement some other memeber in a base class.

Sometimes when designing an inheritance hierarchy, you want to ensure that a member cannot be overridden by subclass. In Scala, as in Java, you can do this by adding a final modifier to the member.

You may also at times want to ensure that an entire class not subclassed. To do this you simply declare the entire class final by adding a final modifier to the class declaration.
>>>>>>> fa004fdc3394401077c95432c15e7df8a3cab039
