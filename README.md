#Semantic Code Search Engine

Q: What is the Semantic Code Search Engine (SemQuery)?  
A: SemQuery is a search engine intended for querying codebases.

###Full Explanation:
Code is a system used to represent logical concepts such as looping, conditions, and other computational actions in textual form. Many programming languages represent similar concepts using different textual forms. For example, this is how a function that adds two numbers, a and b, is declared in various languages:

######C:
```C
int add(int a, int b) { 
  return a + b;
}
```
######JavaScript:
```javascript
function add(a, b) { 
  return a + b 
}
```
######Ruby:
```ruby
def add(a, b)
  a + b
end
```
######Go:
```go
func add(a int, b int) int { 
  return a + b 
}
```

Although all these declarations represent roughly the same concept, they have different textual forms. In fact, a concept can have multiple textual forms within a single language. This C snippet below defines an identifier, foo, in three different ways.
```C
// first way
int foo = 3;

// second way
#define foo 3

// third way
void foo() {}
```
These three declarations have differences, such as type, and how they are processed by the compiler; nonetheless, they are all declarations that define foo in some way.

Say we are inspecting a C codebase/program we are unfamiliar with and we come across this:
```C
printf(“%p”, foo);
```
What is foo (where was it defined)?  Perhaps if it was defined in the line above it is not hard to figure out. But what if it was defined hundreds of lines above? Or in a different file? Search for foo, of course. But any occurrence of foo is not necessarily the definition. Ok, search for a definition of foo, in that case… right? Not so easy, definitions can have multiple forms (as demonstrated above!). Even if you searched for the correct type of declaration, different coding styles, whitespace, etc. could mess up your search. Why is this so hard? All we want is to find a definition of foo. If a compiler can do it, we should be able to do it easily as well. 

Now, imagine a search engine that can search for the idea you’re thinking of, that is, it understands the semantics of your code. We could enter a query such as:
```
definition(name = “foo”)
```
and instead of searching for the text `definition(name = “foo”)` directly within the code, it would search the code based on the code’s meaning. 
