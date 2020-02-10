package i18n

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
