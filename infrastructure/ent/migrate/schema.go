// Code generated by ent, DO NOT EDIT.

package migrate

import (
	"entgo.io/ent/dialect/sql/schema"
	"entgo.io/ent/schema/field"
)

var (
	// CertsColumns holds the columns for the "certs" table.
	CertsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "code", Type: field.TypeString},
		{Name: "domain", Type: field.TypeString},
		{Name: "created_at", Type: field.TypeTime},
		{Name: "updated_at", Type: field.TypeTime},
	}
	// CertsTable holds the schema information for the "certs" table.
	CertsTable = &schema.Table{
		Name:       "certs",
		Columns:    CertsColumns,
		PrimaryKey: []*schema.Column{CertsColumns[0]},
	}
	// Tables holds all the tables in the schema.
	Tables = []*schema.Table{
		CertsTable,
	}
)

func init() {
}
