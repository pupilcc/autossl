// Code generated by ent, DO NOT EDIT.

package cert

import (
	"autossl/infrastructure/ent/predicate"
	"time"

	"entgo.io/ent/dialect/sql"
)

// ID filters vertices based on their ID field.
func ID(id int) predicate.Cert {
	return predicate.Cert(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id int) predicate.Cert {
	return predicate.Cert(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id int) predicate.Cert {
	return predicate.Cert(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...int) predicate.Cert {
	return predicate.Cert(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...int) predicate.Cert {
	return predicate.Cert(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id int) predicate.Cert {
	return predicate.Cert(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id int) predicate.Cert {
	return predicate.Cert(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id int) predicate.Cert {
	return predicate.Cert(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id int) predicate.Cert {
	return predicate.Cert(sql.FieldLTE(FieldID, id))
}

// Code applies equality check predicate on the "code" field. It's identical to CodeEQ.
func Code(v string) predicate.Cert {
	return predicate.Cert(sql.FieldEQ(FieldCode, v))
}

// Domain applies equality check predicate on the "domain" field. It's identical to DomainEQ.
func Domain(v string) predicate.Cert {
	return predicate.Cert(sql.FieldEQ(FieldDomain, v))
}

// CreatedAt applies equality check predicate on the "created_at" field. It's identical to CreatedAtEQ.
func CreatedAt(v time.Time) predicate.Cert {
	return predicate.Cert(sql.FieldEQ(FieldCreatedAt, v))
}

// UpdatedAt applies equality check predicate on the "updated_at" field. It's identical to UpdatedAtEQ.
func UpdatedAt(v time.Time) predicate.Cert {
	return predicate.Cert(sql.FieldEQ(FieldUpdatedAt, v))
}

// CodeEQ applies the EQ predicate on the "code" field.
func CodeEQ(v string) predicate.Cert {
	return predicate.Cert(sql.FieldEQ(FieldCode, v))
}

// CodeNEQ applies the NEQ predicate on the "code" field.
func CodeNEQ(v string) predicate.Cert {
	return predicate.Cert(sql.FieldNEQ(FieldCode, v))
}

// CodeIn applies the In predicate on the "code" field.
func CodeIn(vs ...string) predicate.Cert {
	return predicate.Cert(sql.FieldIn(FieldCode, vs...))
}

// CodeNotIn applies the NotIn predicate on the "code" field.
func CodeNotIn(vs ...string) predicate.Cert {
	return predicate.Cert(sql.FieldNotIn(FieldCode, vs...))
}

// CodeGT applies the GT predicate on the "code" field.
func CodeGT(v string) predicate.Cert {
	return predicate.Cert(sql.FieldGT(FieldCode, v))
}

// CodeGTE applies the GTE predicate on the "code" field.
func CodeGTE(v string) predicate.Cert {
	return predicate.Cert(sql.FieldGTE(FieldCode, v))
}

// CodeLT applies the LT predicate on the "code" field.
func CodeLT(v string) predicate.Cert {
	return predicate.Cert(sql.FieldLT(FieldCode, v))
}

// CodeLTE applies the LTE predicate on the "code" field.
func CodeLTE(v string) predicate.Cert {
	return predicate.Cert(sql.FieldLTE(FieldCode, v))
}

// CodeContains applies the Contains predicate on the "code" field.
func CodeContains(v string) predicate.Cert {
	return predicate.Cert(sql.FieldContains(FieldCode, v))
}

// CodeHasPrefix applies the HasPrefix predicate on the "code" field.
func CodeHasPrefix(v string) predicate.Cert {
	return predicate.Cert(sql.FieldHasPrefix(FieldCode, v))
}

// CodeHasSuffix applies the HasSuffix predicate on the "code" field.
func CodeHasSuffix(v string) predicate.Cert {
	return predicate.Cert(sql.FieldHasSuffix(FieldCode, v))
}

// CodeEqualFold applies the EqualFold predicate on the "code" field.
func CodeEqualFold(v string) predicate.Cert {
	return predicate.Cert(sql.FieldEqualFold(FieldCode, v))
}

// CodeContainsFold applies the ContainsFold predicate on the "code" field.
func CodeContainsFold(v string) predicate.Cert {
	return predicate.Cert(sql.FieldContainsFold(FieldCode, v))
}

// DomainEQ applies the EQ predicate on the "domain" field.
func DomainEQ(v string) predicate.Cert {
	return predicate.Cert(sql.FieldEQ(FieldDomain, v))
}

// DomainNEQ applies the NEQ predicate on the "domain" field.
func DomainNEQ(v string) predicate.Cert {
	return predicate.Cert(sql.FieldNEQ(FieldDomain, v))
}

// DomainIn applies the In predicate on the "domain" field.
func DomainIn(vs ...string) predicate.Cert {
	return predicate.Cert(sql.FieldIn(FieldDomain, vs...))
}

// DomainNotIn applies the NotIn predicate on the "domain" field.
func DomainNotIn(vs ...string) predicate.Cert {
	return predicate.Cert(sql.FieldNotIn(FieldDomain, vs...))
}

// DomainGT applies the GT predicate on the "domain" field.
func DomainGT(v string) predicate.Cert {
	return predicate.Cert(sql.FieldGT(FieldDomain, v))
}

// DomainGTE applies the GTE predicate on the "domain" field.
func DomainGTE(v string) predicate.Cert {
	return predicate.Cert(sql.FieldGTE(FieldDomain, v))
}

// DomainLT applies the LT predicate on the "domain" field.
func DomainLT(v string) predicate.Cert {
	return predicate.Cert(sql.FieldLT(FieldDomain, v))
}

// DomainLTE applies the LTE predicate on the "domain" field.
func DomainLTE(v string) predicate.Cert {
	return predicate.Cert(sql.FieldLTE(FieldDomain, v))
}

// DomainContains applies the Contains predicate on the "domain" field.
func DomainContains(v string) predicate.Cert {
	return predicate.Cert(sql.FieldContains(FieldDomain, v))
}

// DomainHasPrefix applies the HasPrefix predicate on the "domain" field.
func DomainHasPrefix(v string) predicate.Cert {
	return predicate.Cert(sql.FieldHasPrefix(FieldDomain, v))
}

// DomainHasSuffix applies the HasSuffix predicate on the "domain" field.
func DomainHasSuffix(v string) predicate.Cert {
	return predicate.Cert(sql.FieldHasSuffix(FieldDomain, v))
}

// DomainEqualFold applies the EqualFold predicate on the "domain" field.
func DomainEqualFold(v string) predicate.Cert {
	return predicate.Cert(sql.FieldEqualFold(FieldDomain, v))
}

// DomainContainsFold applies the ContainsFold predicate on the "domain" field.
func DomainContainsFold(v string) predicate.Cert {
	return predicate.Cert(sql.FieldContainsFold(FieldDomain, v))
}

// CreatedAtEQ applies the EQ predicate on the "created_at" field.
func CreatedAtEQ(v time.Time) predicate.Cert {
	return predicate.Cert(sql.FieldEQ(FieldCreatedAt, v))
}

// CreatedAtNEQ applies the NEQ predicate on the "created_at" field.
func CreatedAtNEQ(v time.Time) predicate.Cert {
	return predicate.Cert(sql.FieldNEQ(FieldCreatedAt, v))
}

// CreatedAtIn applies the In predicate on the "created_at" field.
func CreatedAtIn(vs ...time.Time) predicate.Cert {
	return predicate.Cert(sql.FieldIn(FieldCreatedAt, vs...))
}

// CreatedAtNotIn applies the NotIn predicate on the "created_at" field.
func CreatedAtNotIn(vs ...time.Time) predicate.Cert {
	return predicate.Cert(sql.FieldNotIn(FieldCreatedAt, vs...))
}

// CreatedAtGT applies the GT predicate on the "created_at" field.
func CreatedAtGT(v time.Time) predicate.Cert {
	return predicate.Cert(sql.FieldGT(FieldCreatedAt, v))
}

// CreatedAtGTE applies the GTE predicate on the "created_at" field.
func CreatedAtGTE(v time.Time) predicate.Cert {
	return predicate.Cert(sql.FieldGTE(FieldCreatedAt, v))
}

// CreatedAtLT applies the LT predicate on the "created_at" field.
func CreatedAtLT(v time.Time) predicate.Cert {
	return predicate.Cert(sql.FieldLT(FieldCreatedAt, v))
}

// CreatedAtLTE applies the LTE predicate on the "created_at" field.
func CreatedAtLTE(v time.Time) predicate.Cert {
	return predicate.Cert(sql.FieldLTE(FieldCreatedAt, v))
}

// UpdatedAtEQ applies the EQ predicate on the "updated_at" field.
func UpdatedAtEQ(v time.Time) predicate.Cert {
	return predicate.Cert(sql.FieldEQ(FieldUpdatedAt, v))
}

// UpdatedAtNEQ applies the NEQ predicate on the "updated_at" field.
func UpdatedAtNEQ(v time.Time) predicate.Cert {
	return predicate.Cert(sql.FieldNEQ(FieldUpdatedAt, v))
}

// UpdatedAtIn applies the In predicate on the "updated_at" field.
func UpdatedAtIn(vs ...time.Time) predicate.Cert {
	return predicate.Cert(sql.FieldIn(FieldUpdatedAt, vs...))
}

// UpdatedAtNotIn applies the NotIn predicate on the "updated_at" field.
func UpdatedAtNotIn(vs ...time.Time) predicate.Cert {
	return predicate.Cert(sql.FieldNotIn(FieldUpdatedAt, vs...))
}

// UpdatedAtGT applies the GT predicate on the "updated_at" field.
func UpdatedAtGT(v time.Time) predicate.Cert {
	return predicate.Cert(sql.FieldGT(FieldUpdatedAt, v))
}

// UpdatedAtGTE applies the GTE predicate on the "updated_at" field.
func UpdatedAtGTE(v time.Time) predicate.Cert {
	return predicate.Cert(sql.FieldGTE(FieldUpdatedAt, v))
}

// UpdatedAtLT applies the LT predicate on the "updated_at" field.
func UpdatedAtLT(v time.Time) predicate.Cert {
	return predicate.Cert(sql.FieldLT(FieldUpdatedAt, v))
}

// UpdatedAtLTE applies the LTE predicate on the "updated_at" field.
func UpdatedAtLTE(v time.Time) predicate.Cert {
	return predicate.Cert(sql.FieldLTE(FieldUpdatedAt, v))
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.Cert) predicate.Cert {
	return predicate.Cert(sql.AndPredicates(predicates...))
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.Cert) predicate.Cert {
	return predicate.Cert(sql.OrPredicates(predicates...))
}

// Not applies the not operator on the given predicate.
func Not(p predicate.Cert) predicate.Cert {
	return predicate.Cert(sql.NotPredicates(p))
}
