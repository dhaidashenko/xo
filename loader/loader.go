// Package loader contains database schema, type and query loaders.
package loader

import (
	"context"
	"fmt"
	"math/rand"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/kenshaw/snaker"
	"github.com/xo/xo/models"
	"github.com/xo/xo/templates"
)

// loaders are registered database loaders.
var loaders = make(map[string]*Loader)

// Register registers a database loader.
func Register(loader *Loader) {
	loaders[loader.Driver] = loader
}

// Get retrieves a database loader for the provided driver name.
func Get(driver string) *Loader {
	return loaders[driver]
}

// Flags returns the additional driver flags for the loaders.
//
// These should be added to the invocation context for any call to a loader
// func.
func Flags() []FlagSet {
	var drivers []string
	for driver := range loaders {
		drivers = append(drivers, driver)
	}
	sort.Strings(drivers)
	var flags []FlagSet
	for _, driver := range drivers {
		l := loaders[driver]
		if l.Flags == nil {
			continue
		}
		for _, flag := range l.Flags() {
			flags = append(flags, FlagSet{
				Driver: driver,
				Name:   string(flag.ContextKey),
				Flag:   flag,
			})
		}
	}
	return flags
}

// FlagSet is a set of flags for a driver.
type FlagSet struct {
	Driver string
	Name   string
	Flag   Flag
}

// Flag is a option flag.
type Flag struct {
	ContextKey  ContextKey
	Desc        string
	PlaceHolder string
	Default     string
	Short       rune
	Value       interface{}
	Enums       []string
}

// ContextKey is a context key.
type ContextKey string

// Loader loads type information from a database.
type Loader struct {
	Driver           string
	Mask             string
	ViewPrefix       string
	Flags            func() []Flag
	GoType           func(context.Context, string, bool) (string, string, int, error)
	Schema           func(context.Context, models.DB) (string, error)
	Enums            func(context.Context, models.DB, string) ([]*models.Enum, error)
	EnumValues       func(context.Context, models.DB, string, string) ([]*models.EnumValue, error)
	Procs            func(context.Context, models.DB, string) ([]*models.Proc, error)
	ProcParams       func(context.Context, models.DB, string, string) ([]*models.ProcParam, error)
	Tables           func(context.Context, models.DB, string, string) ([]*models.Table, error)
	TableColumns     func(context.Context, models.DB, string, string) ([]*models.Column, error)
	TableSequences   func(context.Context, models.DB, string) ([]*models.Sequence, error)
	TableForeignKeys func(context.Context, models.DB, string, string) ([]*models.ForeignKey, error)
	TableIndexes     func(context.Context, models.DB, string, string) ([]*models.Index, error)
	IndexColumns     func(context.Context, models.DB, string, string, string) ([]*models.IndexColumn, error)
	QueryStrip       func([]string, []string)
	QueryColumns     func(context.Context, models.DB, string, []string) ([]*models.Column, error)
	ViewCreate       interface{}
	ViewDrop         interface{}
}

// NthParam returns the 0-based Nth param for the Loader.
func (l *Loader) NthParam(i int) string {
	mask := l.Mask
	if mask == "" {
		return "?"
	}
	if strings.Contains(mask, "%d") {
		return fmt.Sprintf(mask, i+1)
	}
	return mask
}

// SchemaName loads the active schema name for a database.
func (l *Loader) SchemaName(ctx context.Context, db models.DB) (string, error) {
	if l.Schema != nil {
		return l.Schema(ctx, db)
	}
	return "", nil
}

// intRE matches Go int types.
var intRE = regexp.MustCompile(`^int(32|64)?$`)

// letters for GenRandomID.
const letters = "abcdefghijklmnopqrstuvwxyz0123456789"

// rnd is a random source.
var rnd = rand.New(rand.NewSource(time.Now().UTC().UnixNano()))

// randomID generates a 8 character random string.
func randomID() string {
	buf := make([]byte, 8)
	for i := range buf {
		buf[i] = letters[rnd.Intn(len(letters))]
	}
	return string(buf)
}

// parsePrec parses "type[ (precision[,scale])]" strings returning the parsed
// precision and scale.
func parsePrec(typ string) (string, int, int, error) {
	typ, prec, scale := strings.ToLower(typ), -1, -1
	if m := precRE.FindStringIndex(typ); m != nil {
		s := typ[m[0]+1 : m[1]-1]
		if i := strings.LastIndex(s, ","); i != -1 {
			var err error
			if scale, err = strconv.Atoi(strings.TrimSpace(s[i+1:])); err != nil {
				return "", 0, 0, fmt.Errorf("could not parse scale: %w", err)
			}
			s = s[:i]
		}
		// extract precision
		var err error
		if prec, err = strconv.Atoi(strings.TrimSpace(s)); err != nil {
			return "", 0, 0, fmt.Errorf("could not parse precision: %w", err)
		}
		typ = typ[:m[0]]
	}
	return strings.TrimSpace(typ), prec, scale, nil
}

// precRE is the regexp that matches "(precision[,scale])" definitions in a
// database.
var precRE = regexp.MustCompile(`\(([0-9]+)(\s*,\s*[0-9]+\s*)?\)$`)

// schemaGoType returns Go type and zero for a type, removing a "<schema>."
// prefix when the type is determined to be in the same package.
func schemaGoType(ctx context.Context, typ string) (string, string) {
	if schema := templates.Schema(ctx); strings.HasPrefix(typ, schema+".") {
		// in the same schema, so chop off
		typ = typ[len(schema)+1:]
	}
	s := snaker.SnakeToCamelIdentifier(typ)
	return s, s + "{}"
}
