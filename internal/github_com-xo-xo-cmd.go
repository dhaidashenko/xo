// Code generated by 'yaegi extract github.com/xo/xo/cmd'. DO NOT EDIT.

package internal

import (
	"github.com/xo/xo/cmd"
	"reflect"
)

func init() {
	Symbols["github.com/xo/xo/cmd/cmd"] = map[string]reflect.Value{
		// function, constant and variable definitions
		"BuildQuery":           reflect.ValueOf(cmd.BuildQuery),
		"BuildSchema":          reflect.ValueOf(cmd.BuildSchema),
		"DbLoaderSchema":       reflect.ValueOf(cmd.DbLoaderSchema),
		"Introspect":           reflect.ValueOf(cmd.Introspect),
		"LoadColumns":          reflect.ValueOf(cmd.LoadColumns),
		"LoadEnumValues":       reflect.ValueOf(cmd.LoadEnumValues),
		"LoadEnums":            reflect.ValueOf(cmd.LoadEnums),
		"LoadIndexColumns":     reflect.ValueOf(cmd.LoadIndexColumns),
		"LoadProcParams":       reflect.ValueOf(cmd.LoadProcParams),
		"LoadProcs":            reflect.ValueOf(cmd.LoadProcs),
		"LoadTableForeignKeys": reflect.ValueOf(cmd.LoadTableForeignKeys),
		"LoadTableIndexes":     reflect.ValueOf(cmd.LoadTableIndexes),
		"LoadTables":           reflect.ValueOf(cmd.LoadTables),
		"LoadTypeFields":       reflect.ValueOf(cmd.LoadTypeFields),
		"NewArgs":              reflect.ValueOf(cmd.NewArgs),
		"Open":                 reflect.ValueOf(cmd.Open),
		"ParseQuery":           reflect.ValueOf(cmd.ParseQuery),
		"ParseQueryFields":     reflect.ValueOf(cmd.ParseQueryFields),
		"Run":                  reflect.ValueOf(cmd.Run),
		"SplitFields":          reflect.ValueOf(cmd.SplitFields),

		// type definitions
		"Args":           reflect.ValueOf((*cmd.Args)(nil)),
		"DbParams":       reflect.ValueOf((*cmd.DbParams)(nil)),
		"OutParams":      reflect.ValueOf((*cmd.OutParams)(nil)),
		"QueryParams":    reflect.ValueOf((*cmd.QueryParams)(nil)),
		"SchemaParams":   reflect.ValueOf((*cmd.SchemaParams)(nil)),
		"TemplateParams": reflect.ValueOf((*cmd.TemplateParams)(nil)),
	}
}
