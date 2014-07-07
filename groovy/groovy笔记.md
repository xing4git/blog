### switch statement

Switch supports the following kinds of comparisions:

* Class case values matches if the switchValue is an instanceof the class
* Regular expression case value matches if the string of the switchValue matches the regex
* Collection case value matches if the switchValue is contained in the collection. This also includes ranges too.
* if none of the above are used then the case value matches if the case value equals the switch value

The case statement performs a match on the case value using the isCase(switchValue) method, which defaults to call equals(switchValue) but has been overloaded for various types like Class or regex etc.

So you could create your own kind of matcher class and add an isCase(switchValue) method to provide your own kind of matching.

### each and eachWithIndex

You can use each() and eachWithIndex() in place of most loops.

```
def stringList = [ "java", "perl", "python", "ruby", "c#", "cobol"]
stringList.each {print "$it "}
stringList.eachWithIndex {obj, i -> println "$i: $obj"}

def stringMap = [ "Su" : "Sunday", "Mo" : "Monday", "Tu" : "Tuesday"]
stringMap.each {k,v -> println "$k => $v"}
stringMap.eachWithIndex {obj,i -> println "$i: $obj"}
```

### returning values from if-else and try-catch blocks

Since groovy 1.6, it is possible for if/else and try/catch/finally blocks to return a value when they are the last expression in a method or a closure. No need to explicitly use the return keyword inside these constructs, as long as they are the last expression in the block of code.

```
def method() {
	if(true) 1 else 0
}

assert method()==1
```

For try/catch/finally blocks, the last expression evaluated is the one being returned. If an exception is thrown in the try block, the last expression in the catch block is returned instead. Note that finally blocks donot return any value.

```
def method(bool) {
	try {
		if(bool) throw new Exception("foo")
		1
	} catch(e) {
		2
	} finally {
		3
	}
}

assert method(false)==1
assert method(true)==2
```

### operator overloading

Groovy supports operator overloading which makes working with Numbers, Collections, Maps and various other data structures easier to use.

Various operators in Groovy are mapped onto regular Java method calls on objects. This allows you to provide your own Java or Groovy objects which can take advantage of operator overloading.

Operator | Method
-------- | ------
a+b | a.plus(b)
a-b | a.minus(b)
a*b | a.multiply(b)
a**b | a.power(b)
a/b | a.div(b)
a%b | a.mod(b)
a|b | a.or(b)
a&b | a.and(b)
a^b | a.xor(b)
a++ or ++a | a.next()
a-- or --a | a.previous()
a[b] | a.getAt(b)
a[b]=c | a.putAt(b, c)
a<<b | a.leftShift(b)
a>>b | a.rightShift(b)
switch(a){case(b):} | b.isCase(a)
~a | a.bitwiseNegate()
-a | a.negative()
+a | a.positive()

Note that all the following comparison operators handle nulls gracefully avoiding the throwing of NullPointerException

Operator | Method
-------- | ------
a==b | a.equals(b) or a.compareTo(b)==0
a!=b | !a.equals(b)
a<=>b | a.compareTo(b)
a>b | a.compareTo(b)>0
a>=b | a.compareTo(b)>=0
a<b | a.compareTo(b)<0
a<=b | a.compareTo(b)<=0

** Note: the == operator doesnot always exactly match the .equals method.

### spread operator(*.)

The spread operator is used to invoke an action on all items of an aggregate object. It is equivalent to calling the collect method like so:

```
parent*.action // equivalent to: 
parent.collect {child -> child?.action}
```
The action may either be a method call or property access, and returns a list of the items returned from each child call.

### using invokeMethod and getProperty

In any Groovy class you can override `invokeMethod` which will essentially intercept all method calls (to intercept calls to existing methods, the class additionally has to implement the `GroovyInterceptable` interface). This makes it possible to construct some quite interesting DSLs and builders.

```
class XmlBuilder {
   def out
   XmlBuilder(out) { this.out = out }
   def invokeMethod(String name, args) {
       out << "<$name>"
       if(args[0] instanceof Closure) {
            args[0].delegate = this
            args[0].call()
       }
       else {
           out << args[0].toString()
       }
       out << "</$name>"
   }
}
def xml = new XmlBuilder(new StringBuffer())
xml.html {
    head {
        title "Hello World"
    }
    body {
        p "Welcome!"
    }
}
println xml.out
```
You can also override property access using the `getProperty` and `setProperty` property access hooks:

```
class Expandable {
    def storage = [:]
    def getProperty(String name) { storage[name] }
    void setProperty(String name, value) { storage[name] = value }
}

def e = new Expandable()
e.foo = "bar"
println e.foo
```

### Object-Related operators

* java field

Groovy dynamically creates getter method for all your fields that can be referenced as properties.

```
class X {
	def field
}

x = new X()
x.field = 1
println x.field
```
You can override these getters with your own implementations if you like.

The @ operator allows you to override this behavior and access the field directly:

```
class X {
	def field

	def getField() {
		field + 1
	}
}

x = new X()
x.field = 1
println x.field // 2
println x.@field // 1
```
* safe navigation operator(?.)

The safe navigation operator is used to avoid a NullPointerException. Typically when you have a reference to an object you might need to verify that it is not null before accessing methods or properties of the object. To avoid this, the safe navigation operator will simply return `null` instead of throwing an exception:

```
def user = User.find('admin')
def streetName = user?.address?.street
```

* regular expression operator

Groovy supports regular expressions natively using the ~/string/ expression, which creates a compiled Java Pattern object from the given pattern string. Groovy also supports the =~ (create Matcher) and ==~ (returns boolean, whether String matches the pattern) operators.

```
import java.util.regex.Matcher
import java.util.regex.Pattern

def pattern = ~/\d+/
assert pattern instanceof Pattern
assert "1234" ==~ pattern
```

### classes

Classes are defined in Groovy similar to Java. Methods can be class(static) or instance based and can be public, protected, private and support all the usual Java modifiers like synchronized. Package and class imports use the Java syntax. Groovy automatically imports the following:

* java.lang
* java.io
* java.math
* java.net
* java.util
* groovy.lang
* groovy.util

One difference between Java and Groovy is that by default methods are public unless you specify otherwise. Groovy also merges the idea of fields and properties together to make code simpler. Here is an example:

```
class Customer {
    // properties
    Integer id
    String name
    Date dob

    // sample code
    static void main(args) {
        def customer = new Customer(id:1, name:"Gromit", dob:new Date())
        println("Hello ${customer.name}")
    }
}
```
The Groovy code above is equivalent to the following Java code:

```
import java.util.Date;

public class Customer {
    // properties
    private Integer id;
    private String name;
    private Date dob;

    public Integer getId() {
        return this.id;
    }

    public String getName() {
        return this.name;
    }

    public Date getDob() {
        return this.dob;
    }

    public void setId(Integer id) {
        this.id = id;
    }

    public void setName(String name) {
        this.name = name;
    }

    public void setDob(Date dob) {
        this.dob = dob;
    }

    // sample code
    public static void main(String[] args) {
        Customer customer = new Customer();
        customer.setId(1);
        customer.setName("Gromit");
        customer.setDob(new Date());

        System.out.println("Hello " + customer.getName());
    }
}
```

When Groovy is compiled to bytecode, the following rules are used:

* if the name is declared with an access modifier then a field is generated.
* a name declared with no access modifier generates a private field with public getter and setter (i.e. a property)
* if a property is declared final the private field is created final and no setter is generated
* you can declare a property and also declare your own getter or setter
* you can declare a property and a field of the same name, the property will use that field then
* if you want a private or protected property you have to provide your own getter and setter which must be declared private or protected.
* if you access a property from within the class the property is defined in at compile time with implicit or explicit this, Groovy will access the field directly instead of going though the getter and setter
* if you access a property that does not exist using the explicit or implicit foo, then Groovy will access the property through the meta class, which may fail at runtime.

Each class in Groovy is a Java class at the bytecode/JVM level. Any methods declared will be available to Java and vice versa. You can specify the types of parameters or return types on methods so that they work nicely in normal Java code. If you omit the types of any methods or properties they will default to java.lang.Object at the bytecode/JVM level.

### multiple assignments

Groovy is able to define and assign several variables at once:

```
def geo(location) {
	[48.824068, 2.531733]
}

def (lat,lon) = geo('ShangHai')
```

And you can also define the types of variables: `def (double lat, double lon) = geo("ShangHai")`

If the list on the right-hand side contains more elements than the number of variables on the left-hand side, only the first elements will be assigned in order into the variables. Also, when there are less elements than variables, the extra variables will assigned null.

### closure

Closure parameters are listed before the `->` token, like so: `def printSum = {a,b -> print a+b}`. The `->` token is optional and may be omitted if your closure definition takes fewer than two parameters. A closure without `->` is a closure with one argument that is implicitly named as `it`. In some case, you need to construct a closure with zero parameters, you have to explicity define your closure as `{->}` instead of just `{}`.

Closures may refer to variables not listed in their parameter list.

```
def myConst = 5
def inc = {num -> num+myConst}
myConst = 10
println inc(20) // 30
```

Within a Groovy closure, several variables are defined that have special meaning. 

* it. If you have a closure that takes a single argument, you may omit the parameter definition of the closure: `def clos = {print it}; clos("hi there");`
* this. as in Java, `this` refers to the instance of the enclosing class where a closure is defined
* owner. the enclosing object(`this` or a surrounding closre)
* delegate. by default the same as `owner`, but changeable.

```
class Class1 {
  def closure = {
    println this.class.name
    println delegate.class.name
    def nestedClos = {
      println owner.class.name
    }
    nestedClos()
  }
}

def clos = new Class1().closure
clos.delegate = this
clos()

/*  prints:
 Class1
 Script1
 Class1$_closure1  
 */
```







