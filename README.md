# i18n ![wip](https://img.shields.io/badge/-work%20in%20progress-red) ![draft](https://img.shields.io/badge/-draft-red) [![Travis-CI](https://travis-ci.com/worldiety/i18n.svg?branch=master)](https://travis-ci.com/worldiety/i18n) [![Go Report Card](https://goreportcard.com/badge/github.com/worldiety/i18n)](https://goreportcard.com/report/github.com/worldiety/i18n) [![GoDoc](https://godoc.org/github.com/worldiety/i18n?status.svg)](http://godoc.org/github.com/worldiety/i18n)


A library (ready) and a go (golang) generator (wip) which creates code based and type safe translation units.

## milestones and road map

- [x] Android xml support
- [x] CLDR plural support
- [x] CLDR language tag support
- [x] support priority matching of wanted locales and available locales
- [ ] dynamic fallthrough resources, if strings are missing
- [ ] compile time checker for kind of value and placeholders
- [ ] runtime checker for kind of value and placeholders
- [ ] runtime checker for consistent placeholders across translations
- [ ] type safe generator for accessor facade

## library usage

1. use the [Android XML Format](https://developer.android.com/guide/topics/resources/string-resource).
1. import the i18n dependency `go get github.com/worldiety/i18n` in your module.
1. configure and usage is as easy as this
    ```go
        package main
   
        import "github.com/worldiety/i18n"
   
        func main(){
           err := i18n.ImportFile(i18n.AndroidImporter{}, "en-US", "usecase/strings.xml")
           if err != nil {
              panic(err)
           }
   
           res := i18n.From("de-DE")   
   
           str, err := res.QuantityText("x_has_y_cats2", 1, "nick", 1)
           if err != nil {
              panic(err)
           }
           
           fmt.Println(str)
        }   
    ```

## related work
Popular existing libraries are [go-18n](https://github.com/nicksnyder/go-i18n) or 
[i18n4go](https://github.com/maximilien/i18n4go). There is also a pending localization 
[proposal](https://go.googlesource.com/proposal/+/master/design/12750-localization.md). 

## goals and design decisions
The known tools are at their core simple libraries and fall entirely short when it comes to type safety. 
This can only be avoided by one of the following two approaches:

  1. create a linter which runs before any compilation and proofs that whatever text based solution you 
  use, you have consistent translations (e.g. a translation for each key, equal placeholders and plurals 
  for each key) and that you use the keys and formatting methods correctly and consistently (e.g. correct 
  sprintf directives for correct types) OR
  1. create a generator which creates source code from your text based translation configuration and 
  solve all the hassle using simply the type system of your programming language. Even if your language 
  does not provide type safety, the generator can also provide the role of a linter.

The following decisions have been discussed
  1. A new tool should support go modules and go packages. Instead of writing the code first, we assume that it is 
  equally fine or better to write a default translation first, to ensure that you have always a valid text at your hand.
  1. Every access should only be made by type safe accessors, which provides type safe parameters for ordered 
  placeholders and pluralization.
  1. A good encapsulation strategy requires to put related things together, sometimes just on module level but in larger
  projects also per package level. So this applies also to translations, which may be scattered across packages to fit
  best to your divide and conquer strategy.
  1. However, scattering translation files wildly across a module, or even worse, across modules of modules, is probably 
  not desirable for your translation (agency) process and perhaps not feasible at all, because you may be out of control of
  some modules. At best, you have to provide a single file per language in a common format and get the translated 
  languages also back the same way.
  1. The conclusion is to have a single state of truth at the top of your root module, which aggregates and merges
  all translations together and is also the truth for the generated type safe accessors.
  1. A statically proofed translation cannot be guaranteed, if the values can be overridden after generation
  time. So there should be also a runtime checker at startup, because the trade of for a slower start is better than
  a malfunction or crash of your productive service.
  1. The value of introducing a central dependency to a translation dictionary is better than to expect that a developer
  is aware of registering each translatable package from unknown modules by hand. This can only be accomplished 
  with a singleton.
  1. The supported file format must be a well known format, so that common translation software used by agencies
  can simply import and export them (see also for example available SDL Trados
  [file formats](https://docs.sdl.com/LiveContent/content/en-US/SDL%20Passolo-v1/GUID-93FC4141-8209-40A0-B2D6-6E2B8B471D1F#addHistory=true&filename=GUID-AE8DADC4-AE34-4E32-BEAC-F23586BA1DAE.xml&docid=GUID-B2D20814-5CFC-464E-9696-2A19261C0589&inner_id=&tid=&query=&scope=&resource=&toc=false&eventType=lcContent.loadDocGUID-B2D20814-5CFC-464E-9696-2A19261C0589)). Obviously a custom JSON or even TOML format is usually a bad choice.
  
  
  


## go generate usage

1. use the [Android XML Format](https://developer.android.com/guide/topics/resources/string-resource).
 In contrast to the specification, the file name is important and must be prefixed with *strings-* and postfixed with
 the locale, e.g. `mymodule/myusecase/strings-en-US.xml`. For the default fallback language the name *strings.xml*
 is sufficient. 
    ```xml
    <resources>
        <string name="app_name" translatable="false">EasyApp</string>
        <string name="hello_world">Hello World</string>
        <string name="hello_x">Hello %s</string>
        <string name="x_runs_around_Y_and_sings_z">%1s runs around the %2$s and sings %3$s</string>
        <plurals name="x_has_y_cats">
            <item quantity="one">%1$s has %2$d cat</item>
            <item quantity="other">the owner of %2$d cats is %1$s</item>
        </plurals>
        <string-array name="selector_details_array">
            <item>first line</item>
            <item>second line</item>
            <item>third line</item>
            <item>fourth line</item>
        </string-array>
    
        <plurals name="x_has_y_cats2">
            <item quantity="one">%1$s has %2$d cat2</item>
            <item quantity="other">the owner of %2$d cats2 is %1$s</item>
        </plurals>
    
        <string-array name="selector_details_array2">
            <item>a</item>
            <item>b</item>
            <item>c</item>
            <item>d</item>
        </string-array>
    
        <string name="bad_0">\@ \? &lt; &amp; &apos; &quot; \" \'</string>
        <string name="bad_1">"hello '"</string>
    </resources>
    ```
1. import the i18n dependency `go get github.com/worldiety/i18n` in your module.
1. create a generator file, e.g. `mymodule/gen/i18n.go`
    ```go
   package main
   
   import "github.com/worldiety/i18n" 
   
   func main(){
       // invoke the generator in your current project. It will process the entire module.
       i18n.Bundle()
   }
    ```
1. create a file in the root of your module, e.g. in `myproject/gen.go` 
   ```go
   package myproject
   
   //go:generate go run gen/i18n.go
   ```
1. invoke `go generate` and you are done. For each file set within a package you have now a `strings_gen.go`
   file, which contains a *Strings* struct and an according constructor. 

The example output for this example would be `mymodule/myusecase/strings.go`:

```go
// Code generated by go generate; DO NOT EDIT.
package myusecase

import "github.com/worldiety/i18n" 


func init(){
    // add default values, may be overloaded any time later. This is a singleton, see discussion above.
    i18n.Add(i18n.Text{ID:"hello_world", Locale: "en-US", Comment: "HelloWorld returns the text for saying hello world"})
    // ...
}
// Strings is a type safe wrapper around i8n resources
type Strings struct {
    db *i18n.Resources
}

// NewStrings returns a type safe wrapper around an i8n database
func NewStrings(locale string)Strings{
    // tbd validation
    return Strings{i18n.For(locale)}
} 

 // HelloWorld returns the text for saying hello world
func (s Strings) HelloWorld() string {
    return s.db.Text("hello_world")
}

   // HelloX returns a string where X has been replaced by a value.
func (s Strings) HelloX(x string)string{
    return s.db.Text("hello_x", x)
}

   // XRunsAroundYAndSingsZ returns an interpolated and positional string
func (s Strings) XRunsAroundYAndSingsZ(x,y,z string) string{
    return s.db.Text("x_runs_around_Y_and_sings_z", x, y, z)
}

   // XHasYCats returns an interpolated and pluralized sentence.
func (s Strings) XHasYCats(quantity int, x string, y int) string{
       return s.db.PluralText(quantity, "x_has_y_cats", x, y)
}
   
   
func (s Strings) SelectorDetailsArray() []string{
    return s.db.Array("selector_details_array", x, y, z)
}

```


# releases

No code has been written yet.
