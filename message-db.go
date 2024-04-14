package main

import (
	"go.riyazali.net/sqlite"
)

func init() {
	sqlite.Register(func(api *sqlite.ExtensionApi) (sqlite.ErrorCode, error) {
		for name, fn := range funcs(api.Connection()) {
			if err := api.CreateFunction(name, fn); err != nil {
				return sqlite.SQLITE_ERROR, err
			}
		}
		return sqlite.SQLITE_OK, nil
	})
}

func funcs(c *sqlite.Conn) map[string]sqlite.Function {
	return map[string]sqlite.Function{
		"create_table": &createTable{c},
	}
}

func main() {}

type createTable struct{ c *sqlite.Conn }

func (c *createTable) Args() int           { return 0 }
func (c *createTable) Deterministic() bool { return false }
func (c *createTable) Apply(ctx *sqlite.Context, values ...sqlite.Value) {
	const ddl = `
		CREATE TABLE IF NOT EXISTS messages (
			global_position INTEGER PRIMARY KEY,
			position INTEGER NOT NULL,
			time DATETIME NOT NULL,
			stream_name TEXT NOT NULL,
			type TEXT NOT NULL,
			data JSONB,
			metadata JSONB,
			id UUID NOT NULL
		)
	`

	if err := c.c.Exec(ddl, nil); err != nil {
		ctx.ResultError(err)
	}
}
