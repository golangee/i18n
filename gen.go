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

import (
	"fmt"
	. "github.com/dave/jennifer/jen"
	"github.com/golangee/i18n/internal"
	"github.com/iancoleman/strcase"
	"golang.org/x/text/language"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

const stringsPrefix = "strings"
const stringsPostfix = ".xml"

type resourceFile struct {
	filename string
	values   *Resources
}

type packageTranslation struct {
	pkg   *internal.Package
	files []resourceFile
}

func (t *packageTranslation) Emit() error {
	var tmp []*Resources
	for _, res := range t.files {
		tmp = append(tmp, res.values)
	}
	err := validate(tmp)
	if err != nil {
		return err
	}

	file := NewFile(t.pkg.Name)
	file.HeaderComment("Code generated by go generate; DO NOT EDIT.")
	file.HeaderComment("This file was generated by github.com/golangee/i18n")

	// the value import
	file.Func().Id("init").Params().BlockFunc(func(group *Group) {
		group.Var().Id("tag").String()
		for _, resFile := range t.files {
			group.Line()
			group.Comment("from " + filepath.Base(resFile.filename))
			group.Id("tag").Op("=").Lit(resFile.values.tag.String())
			group.Line()
			for _, k := range resFile.values.Keys() {
				val := resFile.values.values[k]
				val.goEmitImportValue(group)
			}
			group.Id("_").Op("=").Id("tag")
		}
		group.Line()

	})

	// typesafe accessors
	file.Comment("Resources wraps the package strings to get invoked safely.")
	file.Type().Id("Resources").Struct(Id("res").Op("*").Qual("github.com/golangee/i18n", "Resources"))
	file.Comment("NewResources creates a new localized resource instance.")
	file.Func().Id("NewResources").Params(Id("locale").String()).Id("Resources").BlockFunc(func(group *Group) {
		group.Return(Id("Resources").Op("{").Qual("github.com/golangee/i18n", "From").Call(Id("locale"))).Op("}")
	})
	for _, value := range t.collectValues() {
		file.Comment(strcase.ToCamel(value.ID()) + " returns a translated text for \"" + value.exampleText() + "\"")
		file.Custom(Options{}, value.goEmitGetter())
	}

	// funcmap for templates
	file.Comment("FuncMap returns the named functions to be used with a template")
	file.Func().Params(Id("r").Id("Resources")).Id("FuncMap").Params().Map(Id("string")).Id("interface{}").BlockFunc(func(group *Group) {
		group.Id("m").Op(":=").Make(Map(Id("string")).Id("interface{}"))
		for _, value := range t.collectValues() {
			methodName := strcase.ToCamel(value.ID())
			group.Id("m").Op("[").Lit(methodName).Op("]=").Id("r").Dot(methodName)
		}
		group.Return(Id("m"))
	})

	dstFname := filepath.Join(t.pkg.Dir, "strings_gen.go")
	f, err := os.Create(dstFname)
	if err != nil {
		return fmt.Errorf("cannot write to %s: %w", dstFname, err)
	}
	defer func() {
		_ = f.Close()
	}()
	err = file.Render(f)
	if err != nil {
		_ = os.Remove(dstFname)
		return fmt.Errorf("generated invalid go code %s: %w", dstFname, err)
	}
	return nil
}

// Returns all available values, aggregated across all translations. It does not perform a validation and uses
// a random value per key. However the returned values are sorted by their id.
func (t *packageTranslation) collectValues() []Value {
	tmp := make(map[string]Value)
	for _, file := range t.files {
		for _, value := range file.values.values {
			tmp[value.ID()] = value
		}
	}
	keys := make([]string, 0, len(tmp))
	for k := range tmp {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	res := make([]Value, 0, len(tmp))
	for _, k := range keys {
		res = append(res, tmp[k])
	}
	return res
}

type goGenerator struct {
	dir          string
	pgk          *internal.Package
	translations []*packageTranslation
}

func newGoGenerator(dir string) *goGenerator {
	return &goGenerator{dir: dir}
}

// Scan identifies all available package translations
func (g *goGenerator) Scan() error {
	pkg, err := internal.GoList(g.dir, true)
	if err != nil {
		return fmt.Errorf("%s is not a module: %w", g.dir, err)
	}
	g.pgk = pkg

	return g.scanCandidates(pkg)
}

func (g *goGenerator) scanCandidates(root *internal.Package) error {
	importer := AndroidImporter{}

	var androidTranslationFiles []resourceFile
	for _, file := range root.ListFiles() {
		fname := filepath.Base(file)
		if strings.HasPrefix(fname, stringsPrefix) && strings.HasSuffix(fname, stringsPostfix) {
			localeName := fname[len(stringsPrefix) : len(fname)-len(stringsPostfix)]
			tag := language.Make(localeName)
			res := newResources(tag)
			reader, err := os.Open(file)
			if err != nil {
				return fmt.Errorf("unable to open %s: %w", file, err)
			}
			err = importer.Import(res, reader)
			_ = reader.Close()
			if err != nil {
				return fmt.Errorf("unable to import %s: %w", file, err)
			}
			androidTranslationFiles = append(androidTranslationFiles, resourceFile{
				filename: file,
				values:   res,
			})
		}
	}
	if len(androidTranslationFiles) > 0 {
		g.translations = append(g.translations, &packageTranslation{
			pkg:   root,
			files: androidTranslationFiles,
		})
	}

	for _, child := range root.Packages {
		err := g.scanCandidates(child)
		if err != nil {
			return err
		}
	}
	return nil
}

func (g *goGenerator) Emit() error {
	for _, translation := range g.translations {
		err := translation.Emit()
		if err != nil {
			return err
		}
	}
	return nil
}

func (p pluralValue) goEmitGetter() *Statement {
	params := ParsePrintf(p.other)
	return Func().Params(Id("r").Id("Resources")).Id(strcase.ToCamel(p.ID())).ParamsFunc(func(group *Group) {
		group.Id("quantity").Int()
		emitParams(params, group)
	}).String().BlockFunc(func(group *Group) {
		group.Id("str").Op(",").Id("err").Op(":=").Id("r").Dot("res").Dot("QuantityText").ParamsFunc(func(group *Group) {
			group.Lit(p.ID())
			group.Id("quantity")
			emitCallParams(params, group)
		})

		emitCheckReturn(p, group)
	})
}

func (p pluralValue) exampleText() string {
	return p.other
}

func (p pluralValue) goEmitImportValue(group *Group) {
	call := Qual("github.com/golangee/i18n", "NewQuantityText").Params(Id("tag"), Lit(p.Id))
	if len(p.zero) > 0 {
		call = call.Dot("Zero").Params(Lit(p.zero))
	}

	if len(p.one) > 0 {
		call = call.Dot("One").Params(Lit(p.one))
	}
	if len(p.two) > 0 {
		call = call.Dot("Two").Params(Lit(p.two))
	}
	if len(p.few) > 0 {
		call = call.Dot("Few").Params(Lit(p.few))
	}
	if len(p.many) > 0 {
		call = call.Dot("Many").Params(Lit(p.many))
	}
	if len(p.other) > 0 {
		call = call.Dot("Other").Params(Lit(p.other))
	}

	group.Qual("github.com/golangee/i18n", "ImportValue").Params(call)
}

func emitParams(params []PrintfFormatSpecifier, group *Group) {
	for i, p := range params {
		switch p.Verb() {
		case 'd':
			group.Id("num" + strconv.Itoa(i)).Int()
		case 'f':
			group.Id("fl" + strconv.Itoa(i)).Float64()
		case 's':
			group.Id("str" + strconv.Itoa(i)).String()
		default:
			group.Id("val" + strconv.Itoa(i)).Interface()
		}
	}
}

func emitCallParams(params []PrintfFormatSpecifier, group *Group) {
	for i, p := range params {
		switch p.Verb() {
		case 'd':
			group.Id("num" + strconv.Itoa(i))
		case 'f':
			group.Id("fl" + strconv.Itoa(i))
		case 's':
			group.Id("str" + strconv.Itoa(i))
		default:
			group.Id("val" + strconv.Itoa(i))
		}
	}
}

func emitCheckReturn(s Value, group *Group) {
	group.If(Id("err").Op("!=").Nil()).Block(Return(Qual("fmt", "Errorf").Call(Lit("MISS!" + s.ID() + ": %w").Op(",").Id("err")).Dot("Error").Call()))
	group.Return(Id("str"))
}

func (s simpleValue) goEmitGetter() *Statement {
	params := ParsePrintf(s.String)
	return Func().Params(Id("r").Id("Resources")).Id(strcase.ToCamel(s.ID())).ParamsFunc(func(group *Group) {
		emitParams(params, group)
	}).String().BlockFunc(func(group *Group) {
		group.Id("str").Op(",").Id("err").Op(":=").Id("r").Dot("res").Dot("Text").ParamsFunc(func(group *Group) {
			group.Lit(s.ID())
			emitCallParams(params, group)
		})
		emitCheckReturn(s, group)

	})
}

func (s simpleValue) exampleText() string {
	return s.String
}

func (s simpleValue) goEmitImportValue(group *Group) {
	call := Qual("github.com/golangee/i18n", "NewText").Params(Id("tag"), Lit(s.Id), Lit(s.String))
	group.Qual("github.com/golangee/i18n", "ImportValue").Params(call)
}

func (a arrayValue) goEmitGetter() *Statement {
	return Func().Params(Id("r").Id("Resources")).Id(strcase.ToCamel(a.ID())).Params().Op("[]").String().BlockFunc(func(group *Group) {
		group.Id("str").Op(",").Id("err").Op(":=").Id("r").Dot("res").Dot("TextArray").Params(Lit(a.ID()))

		group.If(Id("err").Op("!=").Nil()).Block(Return(Op("[]string{").Qual("fmt", "Errorf").Call(Lit("MISS!" + a.ID() + ": %w").Op(",").Id("err")).Dot("Error").Call().Op("}")))
		group.Return(Id("str"))
	})
}

func (a arrayValue) exampleText() string {
	str, _ := a.Text()
	return str
}

func (a arrayValue) goEmitImportValue(group *Group) {
	varArgs := ListFunc(func(group *Group) {
		for _, s := range a.Strings {
			group.Lit(s)
		}
	})
	call := Qual("github.com/golangee/i18n", "NewTextArray").Params(Id("tag"), Lit(a.Id), varArgs)
	group.Qual("github.com/golangee/i18n", "ImportValue").Params(call)
}
