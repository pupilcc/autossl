package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"time"
)

// Cert holds the schema definition for the Cert entity.
type Cert struct {
	ent.Schema
}

// Fields of the Cert.
func (Cert) Fields() []ent.Field {
	return []ent.Field{
		field.String("code"),
		field.String("domain"),
		field.Time("created_at").
			Default(time.Now).
			Nillable(),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now).
			Nillable(),
	}
}

// Edges of the Cert.
func (Cert) Edges() []ent.Edge {
	return nil
}
