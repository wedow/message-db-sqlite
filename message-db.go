package main

import (
	"go.riyazali.net/sqlite"
)

func init() {
	funcs := map[string]sqlite.Function{
		"identity": &identity{},
	}

	for name, fn := range funcs {
		sqlite.Register(func(api *sqlite.ExtensionApi) (sqlite.ErrorCode, error) {
			if err := api.CreateFunction(name, fn); err != nil {
				return sqlite.SQLITE_ERROR, err
			}
			return sqlite.SQLITE_OK, nil
		})
	}
}

func main() {}

type identity struct{}

func (m *identity) Args() int           { return 1 }
func (m *identity) Deterministic() bool { return true }
func (m *identity) Apply(ctx *sqlite.Context, values ...sqlite.Value) {
	ctx.ResultValue(values[0])
}
