package geni18n

func main(){
	Translate().
		String("HelloWorld","hello world").
		Add(Text{
			ID:          "HelloX",
			Description: "",
			Zero:        "",
			One:         "",
			Two:         "",
			Few:         "",
			Many:        "",
			Other:       "",
	})
}
// HelloWorld returns the text for saying hello world
HelloWorld() string

// HelloX returns a string where X has been replaced by a value.
HelloX(x string)string

// XRunsAroundYAndSingsZ returns an interpolated and positional string
XRunsAroundYAndSingsZ(x,y,z string) string

// XHasYCats returns an interpolated and pluralized sentence.
XHasYCats(x string, y int) string