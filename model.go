package i18n

type Text struct {
	// ID is a globally unique id
	ID string
	// Description describes the message to give additional
	// context to translators that may be relevant for translation.
	Description string
	Zero        string
	One         string
	Two         string
	Few         string
	Many        string
	Other       string
}
