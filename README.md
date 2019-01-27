# geni8n
A go (golang) generator which creates code based and type safe translation units.

# goals and design decisions
Most translation tools are simply libraries and fall entirely short when it comes to typesafety. This can only be avoided by one of the following two approaches:

  1. create a linter which runs before any compilation and proofs that whatever text based solution you use, you have consistent translations (e.g. a translation for each key, equal placeholders and plurals for each key) and that you use the keys and formatting methods corectly and consistently (e.g. correct sprintf directives for correct types) OR
  1. create a generator which creates source code from your text based translation configuration and solve all the hassle using simply the type system of your programming language. Even if your language does not provide type safety, the generator can also provide the role of a linter.
  
Even if the tool is written in Go and will firstly support only a generator for it, it is not limited to that ecosystem. The tool will support multiple input formats and intentionally generates the code side by side with your translation file. This kind of fragmentation is intentional, to support the developer in encapsulation and his divide and conquer strategy. If a developer has chosen to have a central app-wide translation (which is generally a bad practice when it comes to reusing modules AND translations), he can still do so. This will also ensure that a localization does not introduce an unwanted dependency.

# usage

Installation (*note: do not call this in your project with go.mod, to avoid inclusion as unneeded dependency*)  
```bash
cd /not/my/module/but/go/path/or/bin
go get github.com/worldiety/geni8n
# add resulting binary into your path
```

Apply generator (everything with automatic detection)  
```bash
cd /my/project/dir
geni8n
```

Apply generator (everything set manually)
```bash
geni8n -dir /my/project/dir -targetLang=go -targetArch=default -fallback=en-US
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

after invoking geni8n a *values.go* is generated
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
        return ValuesDeDE
      case strings.Contains(locale,"de-"):
        return ValuesDeDE
      default:
        return ValuesEnUS
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
