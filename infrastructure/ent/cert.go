// Code generated by ent, DO NOT EDIT.

package ent

import (
	"autossl/infrastructure/ent/cert"
	"fmt"
	"strings"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
)

// Cert is the model entity for the Cert schema.
type Cert struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// Code holds the value of the "code" field.
	Code string `json:"code,omitempty"`
	// Domain holds the value of the "domain" field.
	Domain string `json:"domain,omitempty"`
	// CreatedAt holds the value of the "created_at" field.
	CreatedAt *time.Time `json:"created_at,omitempty"`
	// UpdatedAt holds the value of the "updated_at" field.
	UpdatedAt    *time.Time `json:"updated_at,omitempty"`
	selectValues sql.SelectValues
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Cert) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case cert.FieldID:
			values[i] = new(sql.NullInt64)
		case cert.FieldCode, cert.FieldDomain:
			values[i] = new(sql.NullString)
		case cert.FieldCreatedAt, cert.FieldUpdatedAt:
			values[i] = new(sql.NullTime)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Cert fields.
func (c *Cert) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case cert.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			c.ID = int(value.Int64)
		case cert.FieldCode:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field code", values[i])
			} else if value.Valid {
				c.Code = value.String
			}
		case cert.FieldDomain:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field domain", values[i])
			} else if value.Valid {
				c.Domain = value.String
			}
		case cert.FieldCreatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field created_at", values[i])
			} else if value.Valid {
				c.CreatedAt = new(time.Time)
				*c.CreatedAt = value.Time
			}
		case cert.FieldUpdatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field updated_at", values[i])
			} else if value.Valid {
				c.UpdatedAt = new(time.Time)
				*c.UpdatedAt = value.Time
			}
		default:
			c.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the Cert.
// This includes values selected through modifiers, order, etc.
func (c *Cert) Value(name string) (ent.Value, error) {
	return c.selectValues.Get(name)
}

// Update returns a builder for updating this Cert.
// Note that you need to call Cert.Unwrap() before calling this method if this Cert
// was returned from a transaction, and the transaction was committed or rolled back.
func (c *Cert) Update() *CertUpdateOne {
	return NewCertClient(c.config).UpdateOne(c)
}

// Unwrap unwraps the Cert entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (c *Cert) Unwrap() *Cert {
	_tx, ok := c.config.driver.(*txDriver)
	if !ok {
		panic("ent: Cert is not a transactional entity")
	}
	c.config.driver = _tx.drv
	return c
}

// String implements the fmt.Stringer.
func (c *Cert) String() string {
	var builder strings.Builder
	builder.WriteString("Cert(")
	builder.WriteString(fmt.Sprintf("id=%v, ", c.ID))
	builder.WriteString("code=")
	builder.WriteString(c.Code)
	builder.WriteString(", ")
	builder.WriteString("domain=")
	builder.WriteString(c.Domain)
	builder.WriteString(", ")
	if v := c.CreatedAt; v != nil {
		builder.WriteString("created_at=")
		builder.WriteString(v.Format(time.ANSIC))
	}
	builder.WriteString(", ")
	if v := c.UpdatedAt; v != nil {
		builder.WriteString("updated_at=")
		builder.WriteString(v.Format(time.ANSIC))
	}
	builder.WriteByte(')')
	return builder.String()
}

// Certs is a parsable slice of Cert.
type Certs []*Cert
