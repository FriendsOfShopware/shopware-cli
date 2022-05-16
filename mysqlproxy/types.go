package mysqlproxy

import (
	"github.com/dolthub/go-mysql-server/sql"
	"strings"
)

type Entity struct {
	Name       string                    `json:"entity"`
	Properties map[string]EntityProperty `json:"properties"`
}

type EntityProperty struct {
	Type     string      `json:"type"`
	Relation string      `json:"relation"`
	Entity   string      `json:"entity"`
	Flags    interface{} `json:"flags,omitempty"`
}

func (p EntityProperty) GetType() sql.Type {
	switch p.Type {
	case "uuid":
	case "string":
	case "date":
		return sql.Text
	case "json_object":
	case "association":
		return sql.JSON
	case "boolean":
		return sql.Boolean
	case "float":
		return sql.Float64
	case "int":
		return sql.Float64
	}

	return sql.Text
}

func (p EntityProperty) IsPrimary() bool {
	flags, ok := p.Flags.(map[string]interface{})

	if !ok {
		return false
	}

	_, ok = flags["primary_key"]

	return ok
}

func (p EntityProperty) Comment() string {
	if strings.HasPrefix(p.Type, "json") {
		return "json"
	}

	return ""
}
