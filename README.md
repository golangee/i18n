# i18n ![wip](https://img.shields.io/badge/-work%20in%20progress-red) ![draft](https://img.shields.io/badge/-draft-red)

A go (golang) generator which creates code based and type safe translation units.

## related work
Popular existing libraries are [go-18n](https://github.com/nicksnyder/go-i18n) or 
[i18n4go](https://github.com/maximilien/i18n4go). There is also a pending localization 
[proposal](https://go.googlesource.com/proposal/+/master/design/12750-localization.md). 

## goals and design decisions
The known tools are at their core simple libraries and fall entirely short when it comes to typesafety. 
This can only be avoided by one of the following two approaches:

  1. create a linter which runs before any compilation and proofs that whatever text based solution you 
  use, you have consistent translations (e.g. a translation for each key, equal placeholders and plurals 
  for each key) and that you use the keys and formatting methods corectly and consistently (e.g. correct 
  sprintf directives for correct types) OR
  1. create a generator which creates source code from your text based translation configuration and 
  solve all the hassle using simply the type system of your programming language. Even if your language 
  does not provide type safety, the generator can also provide the role of a linter.

The following descision are made
  1. A new tool should support go modules and go packages. Instead of writing the code first, we assume that it is 
  equally fine or better to write a default translation first, to ensure that you have always a valid text at your hand.
  1. Every access should only be made by type safe accessors, which provides type safe parameters for ordered 
  placeholders and pluralization.
  1. A good encapsulation strategy requires to put related things together, sometimes just on module level but in larger
  projects also per package level. So this applies also to translations, which may be scattered across packages to fit
  best to your divide and conquer strategy.
  1. However, scattering translation files wildly accross a module, or even worse, accross modules of modules, is probably not
  desirable for your translation (agency) process and perhaps not feasible at all, because you may be out of control of
  some modules.
  1. The conclusion is to have a single state of truth at the top of your root module, which aggregates and merges
  all translations together and is also the truth for the generated type safe accessors.
  1. A statically proofed translation cannot be guaranteed, if the values can be overriden after generation
  time. So there should be also a runtime checker at startup, because the trade of for a slower start is better than
  a malfunction or crash of your productive service.
  


# usage

Installation (*note: do not call this in your project with go.mod, to avoid inclusion as unneeded dependency*)  
```bash
cd /not/my/module/but/go/path/or/bin
go get github.com/worldiety/geni18n
# add resulting binary into your path
```

Apply generator (everything with automatic detection)  
```bash
cd /my/project/dir
geni18n
```

Apply generator (everything set manually)
```bash
geni18n -dir /my/project/dir -targetLang=go -targetArch=default -fallback=en-US
```

# generated code examples

you write *values-en-US.xml*
```xml
<resources>
    <string name="app_name" translatable="false">EasyApp</string>
    <string name="action_settings">Settings</string>
    <string-array name="selector_details_array">
       <item>first line</item>
       <item>second line</item>
       <item>third line</item>
       <item>fourth line</item>
    </string-array>
    <plurals name="test0">
        <item quantity="one">Test ok</item>
        <item quantity="other">Tests ok</item>
    </plurals>
    <plurals name="test1">
        <item quantity="one">%d test ok</item>
        <item quantity="other">%d tests ok</item>
    </plurals>
  
</resources>
```

after invoking geni18n a *values.go* is generated
```go

  // ValuesEnUS is the generated version of the file values-en-US.xml
  type ValuesEnUS struct{}
  
  // AppName returns the language independent value for app_name
  func (v ValuesEnUS) AppName() string{
     return "EasyApp"
  }
  
  // ActionSettings returns the value of action_settings for en-US
  func (v ValuesEnUS) ActionSettings() string{
     return "Settings"
  }
  
  // SelectorDetailsArray returns the value of selector_details_array for en-US
  func (v ValuesEnUS) SelectorDetailsArray() []string{
     return []string{"first line","second line","third line","fourth line"}
  }
  
  // Test0 returns the value of test0 for en-US
  func (v ValuesEnUS) Test0(quantity int) string{
    switch quantity{
      case 0:
        return return "Test ok"
      default:
        return return "Tests ok"
    }
  }
  
  // Test1 returns the value of test1 for en-US
  func (v ValuesEnUS) Test1(quantity int, param0 int) string{
    switch quantity{
      case 0:
        return return fmt.Sprintf("%d test ok", param0)
      default:
        return return fmt.Sprintf("%d tests ok", param1)
    }
  }
  
  //let the compiler check type safety
  var _ Values = (ValuesDeDE)(nil)
  
  // ValuesDeDE is the generated version of the file values-de-DE.xml
  type ValuesDeDE struct{}
  // ...
  
  // Values is the common contract valid for all translations
  type Values interface{
    // AppName returns the language independent value for app_name
    AppName() string
    
    // ActionSettings returns the localized value of action_settings
    ActionSettings() string
    
    // ...
  }
  
  // LocalizationOf returns the resolved and localized Values type or a fallback and is never nil.
  func LocalizationOf(locale string) Values {
    switch locale{
      case "de-DE":
        return ValuesDeDE{}
      case strings.Contains(locale,"de-"):
        return ValuesDeDE{}
      default:
        return ValuesEnUS{}
    }
  }
  
```

# roadmap

## Version 1.0.0
 * A working prototype which detects automatically translation units (Android style in e.g. values-en-AU.xml) and project type (*.go). Optional parameters are
   * default fallback language (-default)
   * target language (e.g. Java, Go*, ...) (-targetLang)
   * target architecture (e.g. default*, Spring, ...) (-targetArch)
   * workspace (-dir)
 * Automatic Go package detection
 * Generation of individual structs with method accessors, one for each language and a common Interface to ensure that each struct fulfills the required contract. Create a static method which returns an interface (backed by a specifc struct) for a given language string (e.g. en-AU)
 * Basic support for plurals
 * Basic support for string placeholders

# releases

No code has been written yet.
