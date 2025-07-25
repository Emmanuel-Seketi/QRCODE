// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"qr_backend/ent/filereference"
	"qr_backend/ent/predicate"
	"qr_backend/ent/qrcode"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
)

// FileReferenceUpdate is the builder for updating FileReference entities.
type FileReferenceUpdate struct {
	config
	hooks    []Hook
	mutation *FileReferenceMutation
}

// Where appends a list predicates to the FileReferenceUpdate builder.
func (fru *FileReferenceUpdate) Where(ps ...predicate.FileReference) *FileReferenceUpdate {
	fru.mutation.Where(ps...)
	return fru
}

// SetFilename sets the "filename" field.
func (fru *FileReferenceUpdate) SetFilename(s string) *FileReferenceUpdate {
	fru.mutation.SetFilename(s)
	return fru
}

// SetNillableFilename sets the "filename" field if the given value is not nil.
func (fru *FileReferenceUpdate) SetNillableFilename(s *string) *FileReferenceUpdate {
	if s != nil {
		fru.SetFilename(*s)
	}
	return fru
}

// SetURL sets the "url" field.
func (fru *FileReferenceUpdate) SetURL(s string) *FileReferenceUpdate {
	fru.mutation.SetURL(s)
	return fru
}

// SetNillableURL sets the "url" field if the given value is not nil.
func (fru *FileReferenceUpdate) SetNillableURL(s *string) *FileReferenceUpdate {
	if s != nil {
		fru.SetURL(*s)
	}
	return fru
}

// SetSize sets the "size" field.
func (fru *FileReferenceUpdate) SetSize(i int64) *FileReferenceUpdate {
	fru.mutation.ResetSize()
	fru.mutation.SetSize(i)
	return fru
}

// SetNillableSize sets the "size" field if the given value is not nil.
func (fru *FileReferenceUpdate) SetNillableSize(i *int64) *FileReferenceUpdate {
	if i != nil {
		fru.SetSize(*i)
	}
	return fru
}

// AddSize adds i to the "size" field.
func (fru *FileReferenceUpdate) AddSize(i int64) *FileReferenceUpdate {
	fru.mutation.AddSize(i)
	return fru
}

// SetType sets the "type" field.
func (fru *FileReferenceUpdate) SetType(s string) *FileReferenceUpdate {
	fru.mutation.SetType(s)
	return fru
}

// SetNillableType sets the "type" field if the given value is not nil.
func (fru *FileReferenceUpdate) SetNillableType(s *string) *FileReferenceUpdate {
	if s != nil {
		fru.SetType(*s)
	}
	return fru
}

// SetQrCodeID sets the "qr_code" edge to the QRCode entity by ID.
func (fru *FileReferenceUpdate) SetQrCodeID(id int) *FileReferenceUpdate {
	fru.mutation.SetQrCodeID(id)
	return fru
}

// SetNillableQrCodeID sets the "qr_code" edge to the QRCode entity by ID if the given value is not nil.
func (fru *FileReferenceUpdate) SetNillableQrCodeID(id *int) *FileReferenceUpdate {
	if id != nil {
		fru = fru.SetQrCodeID(*id)
	}
	return fru
}

// SetQrCode sets the "qr_code" edge to the QRCode entity.
func (fru *FileReferenceUpdate) SetQrCode(q *QRCode) *FileReferenceUpdate {
	return fru.SetQrCodeID(q.ID)
}

// Mutation returns the FileReferenceMutation object of the builder.
func (fru *FileReferenceUpdate) Mutation() *FileReferenceMutation {
	return fru.mutation
}

// ClearQrCode clears the "qr_code" edge to the QRCode entity.
func (fru *FileReferenceUpdate) ClearQrCode() *FileReferenceUpdate {
	fru.mutation.ClearQrCode()
	return fru
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (fru *FileReferenceUpdate) Save(ctx context.Context) (int, error) {
	return withHooks(ctx, fru.sqlSave, fru.mutation, fru.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (fru *FileReferenceUpdate) SaveX(ctx context.Context) int {
	affected, err := fru.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (fru *FileReferenceUpdate) Exec(ctx context.Context) error {
	_, err := fru.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (fru *FileReferenceUpdate) ExecX(ctx context.Context) {
	if err := fru.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (fru *FileReferenceUpdate) check() error {
	if v, ok := fru.mutation.Filename(); ok {
		if err := filereference.FilenameValidator(v); err != nil {
			return &ValidationError{Name: "filename", err: fmt.Errorf(`ent: validator failed for field "FileReference.filename": %w`, err)}
		}
	}
	if v, ok := fru.mutation.URL(); ok {
		if err := filereference.URLValidator(v); err != nil {
			return &ValidationError{Name: "url", err: fmt.Errorf(`ent: validator failed for field "FileReference.url": %w`, err)}
		}
	}
	return nil
}

func (fru *FileReferenceUpdate) sqlSave(ctx context.Context) (n int, err error) {
	if err := fru.check(); err != nil {
		return n, err
	}
	_spec := sqlgraph.NewUpdateSpec(filereference.Table, filereference.Columns, sqlgraph.NewFieldSpec(filereference.FieldID, field.TypeInt))
	if ps := fru.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := fru.mutation.Filename(); ok {
		_spec.SetField(filereference.FieldFilename, field.TypeString, value)
	}
	if value, ok := fru.mutation.URL(); ok {
		_spec.SetField(filereference.FieldURL, field.TypeString, value)
	}
	if value, ok := fru.mutation.Size(); ok {
		_spec.SetField(filereference.FieldSize, field.TypeInt64, value)
	}
	if value, ok := fru.mutation.AddedSize(); ok {
		_spec.AddField(filereference.FieldSize, field.TypeInt64, value)
	}
	if value, ok := fru.mutation.GetType(); ok {
		_spec.SetField(filereference.FieldType, field.TypeString, value)
	}
	if fru.mutation.QrCodeCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   filereference.QrCodeTable,
			Columns: []string{filereference.QrCodeColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(qrcode.FieldID, field.TypeInt),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := fru.mutation.QrCodeIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   filereference.QrCodeTable,
			Columns: []string{filereference.QrCodeColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(qrcode.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, fru.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{filereference.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	fru.mutation.done = true
	return n, nil
}

// FileReferenceUpdateOne is the builder for updating a single FileReference entity.
type FileReferenceUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *FileReferenceMutation
}

// SetFilename sets the "filename" field.
func (fruo *FileReferenceUpdateOne) SetFilename(s string) *FileReferenceUpdateOne {
	fruo.mutation.SetFilename(s)
	return fruo
}

// SetNillableFilename sets the "filename" field if the given value is not nil.
func (fruo *FileReferenceUpdateOne) SetNillableFilename(s *string) *FileReferenceUpdateOne {
	if s != nil {
		fruo.SetFilename(*s)
	}
	return fruo
}

// SetURL sets the "url" field.
func (fruo *FileReferenceUpdateOne) SetURL(s string) *FileReferenceUpdateOne {
	fruo.mutation.SetURL(s)
	return fruo
}

// SetNillableURL sets the "url" field if the given value is not nil.
func (fruo *FileReferenceUpdateOne) SetNillableURL(s *string) *FileReferenceUpdateOne {
	if s != nil {
		fruo.SetURL(*s)
	}
	return fruo
}

// SetSize sets the "size" field.
func (fruo *FileReferenceUpdateOne) SetSize(i int64) *FileReferenceUpdateOne {
	fruo.mutation.ResetSize()
	fruo.mutation.SetSize(i)
	return fruo
}

// SetNillableSize sets the "size" field if the given value is not nil.
func (fruo *FileReferenceUpdateOne) SetNillableSize(i *int64) *FileReferenceUpdateOne {
	if i != nil {
		fruo.SetSize(*i)
	}
	return fruo
}

// AddSize adds i to the "size" field.
func (fruo *FileReferenceUpdateOne) AddSize(i int64) *FileReferenceUpdateOne {
	fruo.mutation.AddSize(i)
	return fruo
}

// SetType sets the "type" field.
func (fruo *FileReferenceUpdateOne) SetType(s string) *FileReferenceUpdateOne {
	fruo.mutation.SetType(s)
	return fruo
}

// SetNillableType sets the "type" field if the given value is not nil.
func (fruo *FileReferenceUpdateOne) SetNillableType(s *string) *FileReferenceUpdateOne {
	if s != nil {
		fruo.SetType(*s)
	}
	return fruo
}

// SetQrCodeID sets the "qr_code" edge to the QRCode entity by ID.
func (fruo *FileReferenceUpdateOne) SetQrCodeID(id int) *FileReferenceUpdateOne {
	fruo.mutation.SetQrCodeID(id)
	return fruo
}

// SetNillableQrCodeID sets the "qr_code" edge to the QRCode entity by ID if the given value is not nil.
func (fruo *FileReferenceUpdateOne) SetNillableQrCodeID(id *int) *FileReferenceUpdateOne {
	if id != nil {
		fruo = fruo.SetQrCodeID(*id)
	}
	return fruo
}

// SetQrCode sets the "qr_code" edge to the QRCode entity.
func (fruo *FileReferenceUpdateOne) SetQrCode(q *QRCode) *FileReferenceUpdateOne {
	return fruo.SetQrCodeID(q.ID)
}

// Mutation returns the FileReferenceMutation object of the builder.
func (fruo *FileReferenceUpdateOne) Mutation() *FileReferenceMutation {
	return fruo.mutation
}

// ClearQrCode clears the "qr_code" edge to the QRCode entity.
func (fruo *FileReferenceUpdateOne) ClearQrCode() *FileReferenceUpdateOne {
	fruo.mutation.ClearQrCode()
	return fruo
}

// Where appends a list predicates to the FileReferenceUpdate builder.
func (fruo *FileReferenceUpdateOne) Where(ps ...predicate.FileReference) *FileReferenceUpdateOne {
	fruo.mutation.Where(ps...)
	return fruo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (fruo *FileReferenceUpdateOne) Select(field string, fields ...string) *FileReferenceUpdateOne {
	fruo.fields = append([]string{field}, fields...)
	return fruo
}

// Save executes the query and returns the updated FileReference entity.
func (fruo *FileReferenceUpdateOne) Save(ctx context.Context) (*FileReference, error) {
	return withHooks(ctx, fruo.sqlSave, fruo.mutation, fruo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (fruo *FileReferenceUpdateOne) SaveX(ctx context.Context) *FileReference {
	node, err := fruo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (fruo *FileReferenceUpdateOne) Exec(ctx context.Context) error {
	_, err := fruo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (fruo *FileReferenceUpdateOne) ExecX(ctx context.Context) {
	if err := fruo.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (fruo *FileReferenceUpdateOne) check() error {
	if v, ok := fruo.mutation.Filename(); ok {
		if err := filereference.FilenameValidator(v); err != nil {
			return &ValidationError{Name: "filename", err: fmt.Errorf(`ent: validator failed for field "FileReference.filename": %w`, err)}
		}
	}
	if v, ok := fruo.mutation.URL(); ok {
		if err := filereference.URLValidator(v); err != nil {
			return &ValidationError{Name: "url", err: fmt.Errorf(`ent: validator failed for field "FileReference.url": %w`, err)}
		}
	}
	return nil
}

func (fruo *FileReferenceUpdateOne) sqlSave(ctx context.Context) (_node *FileReference, err error) {
	if err := fruo.check(); err != nil {
		return _node, err
	}
	_spec := sqlgraph.NewUpdateSpec(filereference.Table, filereference.Columns, sqlgraph.NewFieldSpec(filereference.FieldID, field.TypeInt))
	id, ok := fruo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "FileReference.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := fruo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, filereference.FieldID)
		for _, f := range fields {
			if !filereference.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != filereference.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := fruo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := fruo.mutation.Filename(); ok {
		_spec.SetField(filereference.FieldFilename, field.TypeString, value)
	}
	if value, ok := fruo.mutation.URL(); ok {
		_spec.SetField(filereference.FieldURL, field.TypeString, value)
	}
	if value, ok := fruo.mutation.Size(); ok {
		_spec.SetField(filereference.FieldSize, field.TypeInt64, value)
	}
	if value, ok := fruo.mutation.AddedSize(); ok {
		_spec.AddField(filereference.FieldSize, field.TypeInt64, value)
	}
	if value, ok := fruo.mutation.GetType(); ok {
		_spec.SetField(filereference.FieldType, field.TypeString, value)
	}
	if fruo.mutation.QrCodeCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   filereference.QrCodeTable,
			Columns: []string{filereference.QrCodeColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(qrcode.FieldID, field.TypeInt),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := fruo.mutation.QrCodeIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   filereference.QrCodeTable,
			Columns: []string{filereference.QrCodeColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(qrcode.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_node = &FileReference{config: fruo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, fruo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{filereference.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	fruo.mutation.done = true
	return _node, nil
}
