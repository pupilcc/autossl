// Code generated by ent, DO NOT EDIT.

package ent

import (
	"autossl/infrastructure/ent/cert"
	"autossl/infrastructure/ent/predicate"
	"context"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
)

// CertUpdate is the builder for updating Cert entities.
type CertUpdate struct {
	config
	hooks    []Hook
	mutation *CertMutation
}

// Where appends a list predicates to the CertUpdate builder.
func (cu *CertUpdate) Where(ps ...predicate.Cert) *CertUpdate {
	cu.mutation.Where(ps...)
	return cu
}

// SetCode sets the "code" field.
func (cu *CertUpdate) SetCode(s string) *CertUpdate {
	cu.mutation.SetCode(s)
	return cu
}

// SetNillableCode sets the "code" field if the given value is not nil.
func (cu *CertUpdate) SetNillableCode(s *string) *CertUpdate {
	if s != nil {
		cu.SetCode(*s)
	}
	return cu
}

// SetDomain sets the "domain" field.
func (cu *CertUpdate) SetDomain(s string) *CertUpdate {
	cu.mutation.SetDomain(s)
	return cu
}

// SetNillableDomain sets the "domain" field if the given value is not nil.
func (cu *CertUpdate) SetNillableDomain(s *string) *CertUpdate {
	if s != nil {
		cu.SetDomain(*s)
	}
	return cu
}

// SetCreatedAt sets the "created_at" field.
func (cu *CertUpdate) SetCreatedAt(t time.Time) *CertUpdate {
	cu.mutation.SetCreatedAt(t)
	return cu
}

// SetNillableCreatedAt sets the "created_at" field if the given value is not nil.
func (cu *CertUpdate) SetNillableCreatedAt(t *time.Time) *CertUpdate {
	if t != nil {
		cu.SetCreatedAt(*t)
	}
	return cu
}

// SetUpdatedAt sets the "updated_at" field.
func (cu *CertUpdate) SetUpdatedAt(t time.Time) *CertUpdate {
	cu.mutation.SetUpdatedAt(t)
	return cu
}

// Mutation returns the CertMutation object of the builder.
func (cu *CertUpdate) Mutation() *CertMutation {
	return cu.mutation
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (cu *CertUpdate) Save(ctx context.Context) (int, error) {
	cu.defaults()
	return withHooks(ctx, cu.sqlSave, cu.mutation, cu.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (cu *CertUpdate) SaveX(ctx context.Context) int {
	affected, err := cu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (cu *CertUpdate) Exec(ctx context.Context) error {
	_, err := cu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (cu *CertUpdate) ExecX(ctx context.Context) {
	if err := cu.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (cu *CertUpdate) defaults() {
	if _, ok := cu.mutation.UpdatedAt(); !ok {
		v := cert.UpdateDefaultUpdatedAt()
		cu.mutation.SetUpdatedAt(v)
	}
}

func (cu *CertUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := sqlgraph.NewUpdateSpec(cert.Table, cert.Columns, sqlgraph.NewFieldSpec(cert.FieldID, field.TypeInt))
	if ps := cu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := cu.mutation.Code(); ok {
		_spec.SetField(cert.FieldCode, field.TypeString, value)
	}
	if value, ok := cu.mutation.Domain(); ok {
		_spec.SetField(cert.FieldDomain, field.TypeString, value)
	}
	if value, ok := cu.mutation.CreatedAt(); ok {
		_spec.SetField(cert.FieldCreatedAt, field.TypeTime, value)
	}
	if value, ok := cu.mutation.UpdatedAt(); ok {
		_spec.SetField(cert.FieldUpdatedAt, field.TypeTime, value)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, cu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{cert.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	cu.mutation.done = true
	return n, nil
}

// CertUpdateOne is the builder for updating a single Cert entity.
type CertUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *CertMutation
}

// SetCode sets the "code" field.
func (cuo *CertUpdateOne) SetCode(s string) *CertUpdateOne {
	cuo.mutation.SetCode(s)
	return cuo
}

// SetNillableCode sets the "code" field if the given value is not nil.
func (cuo *CertUpdateOne) SetNillableCode(s *string) *CertUpdateOne {
	if s != nil {
		cuo.SetCode(*s)
	}
	return cuo
}

// SetDomain sets the "domain" field.
func (cuo *CertUpdateOne) SetDomain(s string) *CertUpdateOne {
	cuo.mutation.SetDomain(s)
	return cuo
}

// SetNillableDomain sets the "domain" field if the given value is not nil.
func (cuo *CertUpdateOne) SetNillableDomain(s *string) *CertUpdateOne {
	if s != nil {
		cuo.SetDomain(*s)
	}
	return cuo
}

// SetCreatedAt sets the "created_at" field.
func (cuo *CertUpdateOne) SetCreatedAt(t time.Time) *CertUpdateOne {
	cuo.mutation.SetCreatedAt(t)
	return cuo
}

// SetNillableCreatedAt sets the "created_at" field if the given value is not nil.
func (cuo *CertUpdateOne) SetNillableCreatedAt(t *time.Time) *CertUpdateOne {
	if t != nil {
		cuo.SetCreatedAt(*t)
	}
	return cuo
}

// SetUpdatedAt sets the "updated_at" field.
func (cuo *CertUpdateOne) SetUpdatedAt(t time.Time) *CertUpdateOne {
	cuo.mutation.SetUpdatedAt(t)
	return cuo
}

// Mutation returns the CertMutation object of the builder.
func (cuo *CertUpdateOne) Mutation() *CertMutation {
	return cuo.mutation
}

// Where appends a list predicates to the CertUpdate builder.
func (cuo *CertUpdateOne) Where(ps ...predicate.Cert) *CertUpdateOne {
	cuo.mutation.Where(ps...)
	return cuo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (cuo *CertUpdateOne) Select(field string, fields ...string) *CertUpdateOne {
	cuo.fields = append([]string{field}, fields...)
	return cuo
}

// Save executes the query and returns the updated Cert entity.
func (cuo *CertUpdateOne) Save(ctx context.Context) (*Cert, error) {
	cuo.defaults()
	return withHooks(ctx, cuo.sqlSave, cuo.mutation, cuo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (cuo *CertUpdateOne) SaveX(ctx context.Context) *Cert {
	node, err := cuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (cuo *CertUpdateOne) Exec(ctx context.Context) error {
	_, err := cuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (cuo *CertUpdateOne) ExecX(ctx context.Context) {
	if err := cuo.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (cuo *CertUpdateOne) defaults() {
	if _, ok := cuo.mutation.UpdatedAt(); !ok {
		v := cert.UpdateDefaultUpdatedAt()
		cuo.mutation.SetUpdatedAt(v)
	}
}

func (cuo *CertUpdateOne) sqlSave(ctx context.Context) (_node *Cert, err error) {
	_spec := sqlgraph.NewUpdateSpec(cert.Table, cert.Columns, sqlgraph.NewFieldSpec(cert.FieldID, field.TypeInt))
	id, ok := cuo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "Cert.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := cuo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, cert.FieldID)
		for _, f := range fields {
			if !cert.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != cert.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := cuo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := cuo.mutation.Code(); ok {
		_spec.SetField(cert.FieldCode, field.TypeString, value)
	}
	if value, ok := cuo.mutation.Domain(); ok {
		_spec.SetField(cert.FieldDomain, field.TypeString, value)
	}
	if value, ok := cuo.mutation.CreatedAt(); ok {
		_spec.SetField(cert.FieldCreatedAt, field.TypeTime, value)
	}
	if value, ok := cuo.mutation.UpdatedAt(); ok {
		_spec.SetField(cert.FieldUpdatedAt, field.TypeTime, value)
	}
	_node = &Cert{config: cuo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, cuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{cert.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	cuo.mutation.done = true
	return _node, nil
}