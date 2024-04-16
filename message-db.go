package main

import (
	"crypto/md5"
	"encoding/hex"
	"strconv"
	"strings"

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
		"id":           &id{},
		"cardinal_id":  &cardinalID{},
		"category":     &category{},
		"is_category":  &isCategory{},
		"hash_64":      &hash64{},
	}
}

func main() {}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

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

type id struct{ c *sqlite.Conn }

func (*id) Args() int           { return 1 }
func (*id) Deterministic() bool { return true }
func (*id) Apply(ctx *sqlite.Context, values ...sqlite.Value) {
	_, id, found := strings.Cut(values[0].Text(), "-")
	if !found {
		ctx.ResultNull()
		return
	}

	ctx.ResultText(id)
}

type cardinalID struct{ c *sqlite.Conn }

func (*cardinalID) Args() int           { return 1 }
func (*cardinalID) Deterministic() bool { return true }
func (*cardinalID) Apply(ctx *sqlite.Context, values ...sqlite.Value) {
	_, id, found := strings.Cut(values[0].Text(), "-")
	if !found {
		ctx.ResultNull()
		return
	}

	cardinalID, _, _ := strings.Cut(id, "+")
	ctx.ResultText(cardinalID)
}

type category struct{ c *sqlite.Conn }

func (*category) Args() int           { return 1 }
func (*category) Deterministic() bool { return true }
func (*category) Apply(ctx *sqlite.Context, values ...sqlite.Value) {
	category, _, _ := strings.Cut(values[0].Text(), "-")
	ctx.ResultText(category)
}

type isCategory struct{ c *sqlite.Conn }

func (*isCategory) Args() int           { return 1 }
func (*isCategory) Deterministic() bool { return true }
func (*isCategory) Apply(ctx *sqlite.Context, values ...sqlite.Value) {
	isCategory := !strings.Contains(values[0].Text(), "-")
	ctx.ResultInt(boolToInt(isCategory))
}

func hash64Fn(value string) (int64, error) {
	hash := md5.Sum([]byte(value))
	hashStr := "0x" + hex.EncodeToString(hash[:])
	hashStr = hashStr[:17]
	hashInt, err := strconv.ParseInt(hashStr, 0, 64)
	if err != nil {
		return 0, err
	}
	return hashInt, nil
}

type hash64 struct{ c *sqlite.Conn }

func (*hash64) Args() int           { return 1 }
func (*hash64) Deterministic() bool { return true }
func (*hash64) Apply(ctx *sqlite.Context, values ...sqlite.Value) {
	hash, err := hash64Fn(values[0].Text())
	if err != nil {
		ctx.ResultError(err)
	}

	ctx.ResultInt64(hash)
}
