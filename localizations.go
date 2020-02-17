package i18n

// nolint: goimports // the linter is broken
import (
	"golang.org/x/text/language"
	"sync"
)

// maxLocaleCacheSize limits how many locales we cache. We want to avoid garbage and safe parsing time. However
// an attacker may flood the cache and may cause an oom easily without a limit.
const maxLocaleCacheSize = 1000

// localizations is our internal implementation using CLDR rules for locales and plurals from
// golang.org/x/text/* implementations.
// All public methods are safe for concurrent usage.
type localizations struct {
	translationPriority []language.Tag

	translations      map[language.Tag]*Resources
	translationsMutex sync.RWMutex

	matcher language.Matcher
	// we expect, that the amount of locales are a small number of finite tags
	parsedTags      map[string]language.Tag
	parsedTagsMutex sync.RWMutex
}

func newLocalizations() *localizations {
	return &localizations{
		translations: make(map[language.Tag]*Resources),
		parsedTags:   make(map[string]language.Tag),
	}
}

// SetTranslationPriority updates the resolution order and removes unwanted translations
func (l *localizations) SetTranslationPriority(order []string) {
	l.translationsMutex.Lock()
	defer l.translationsMutex.Unlock()

	var wanted []language.Tag
	for _, ordered := range order {
		wantedTag := language.Make(ordered)
		for _, existingTag := range l.translationPriority {
			if existingTag == wantedTag {
				wanted = append(wanted, existingTag)
				break
			}
		}
	}

	// free memory for unwanted translations
	for _, t := range l.translationPriority {
		found := false
		for _, w := range wanted {
			if t == w {
				found = true
				break
			}
		}
		if !found {
			delete(l.translations, t)
		}
	}
	l.translationPriority = wanted
}

// Locale parses the given unsafe locale into a valid BCP 47 tag.
func (l *localizations) Locale(locale string) language.Tag {
	l.parsedTagsMutex.RLock()
	tag, exists := l.parsedTags[locale]
	l.parsedTagsMutex.RUnlock()

	if exists {
		return tag
	}

	l.parsedTagsMutex.Lock()
	defer l.parsedTagsMutex.Unlock()

	if len(l.parsedTags) > maxLocaleCacheSize {
		l.parsedTags = make(map[string]language.Tag)
	}

	res := language.Make(locale)
	l.parsedTags[locale] = res

	return res
}

// Configure allocates, if required, a new resource and returns it
func (l *localizations) Configure(locale string) *Resources {
	tag := l.Locale(locale)
	l.translationsMutex.RLock()
	res := l.translations[tag]
	l.translationsMutex.RUnlock()

	if res != nil {
		return res
	}

	l.translationsMutex.Lock()
	defer l.translationsMutex.Unlock()

	res = newResources(tag)
	l.translations[tag] = res
	l.translationPriority = append(l.translationPriority, tag)
	l.matcher = language.NewMatcher(l.translationPriority)

	return res
}

// Returns the best matching resource. If no resources are available, panics because it is a programming error
// to call Match without configuring.
func (l *localizations) Match(locales ...string) *Resources {
	l.translationsMutex.RLock()
	defer l.translationsMutex.RUnlock()

	if len(l.translationPriority) == 0 {
		panic("illegal state: not yet configured")
	}

	tmp := make([]language.Tag, 0, len(locales))
	for _, locale := range locales {
		tmp = append(tmp, l.Locale(locale))
	}

	bestTag, _, _ := l.matcher.Match(tmp...)

	res := l.translations[bestTag]
	if res == nil {
		panic("assert: may not be nil")
	}

	return res
}
