// Copyright 2020 Torben Schinke
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package i18n

// nolint: goimports // the linter is broken
import (
	"golang.org/x/text/language"
	"sync"
)

// Resources is a type for accessing an applications text resources. It is safe to use concurrently.
type Resources struct {
	tag    language.Tag
	values map[string]Value
	mutex  sync.RWMutex
}

func newResources(tag language.Tag) *Resources {
	return &Resources{
		tag:    tag,
		values: make(map[string]Value),
	}
}


// TextArray returns a defensive copy of the according string array
// or ErrTextNotFound.
func (l *Resources) TextArray(id string) ([]string, error) {
	l.mutex.RLock()
	defer l.mutex.RUnlock()

	value := l.values[id]
	if value == nil {
		return nil, ErrTextNotFound
	}

	return value.TextArray()
}

// Text returns a translated string or ErrTextNotFound
func (l *Resources) Text(id string, args ...interface{}) (string, error) {
	l.mutex.RLock()
	defer l.mutex.RUnlock()

	value := l.values[id]
	if value == nil {
		return "", ErrTextNotFound
	}

	return value.Text(args...)
}

// QuantityText returns a translated and grammatically correct pluralization string or ErrTextNotFound
func (l *Resources) QuantityText(id string, quantity int, args ...interface{}) (string, error) {
	l.mutex.RLock()
	defer l.mutex.RUnlock()

	value := l.values[id]
	if value == nil {
		return "", ErrTextNotFound
	}

	return value.QuantityText(quantity, args...)
}
