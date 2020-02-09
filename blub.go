package geni18n

type Text struct {
	// ID is a globally unique id
	ID string

	// Description describes the message to give additional
	// context to translators that may be relevant for translation.
	Description string

	// Zero is the content of the text for the CLDR plural form "zero".
	Zero string

	// One is the content of the text for the CLDR plural form "one".
	One string

	// Two is the content of the text for the CLDR plural form "two".
	Two string

	// Few is the content of the text for the CLDR plural form "few".
	Few string

	// Many is the content of the text for the CLDR plural form "many".
	Many string

	// Other is the content of the text for the CLDR plural form "other".
	Other string
}

type Translation struct {
}

func (t *Translation) String(key, val string) *Translation {
	return nil
}

func (t *Translation) Add(txt Text)*Translation{
	return nil
}

func Translate() *Translation {
	return nil
}
