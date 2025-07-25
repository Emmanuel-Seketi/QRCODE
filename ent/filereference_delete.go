// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"qr_backend/ent/filereference"
	"qr_backend/ent/predicate"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
)

// FileReferenceDelete is the builder for deleting a FileReference entity.
type FileReferenceDelete struct {
	config
	hooks    []Hook
	mutation *FileReferenceMutation
}

// Where appends a list predicates to the FileReferenceDelete builder.
func (frd *FileReferenceDelete) Where(ps ...predicate.FileReference) *FileReferenceDelete {
	frd.mutation.Where(ps...)
	return frd
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (frd *FileReferenceDelete) Exec(ctx context.Context) (int, error) {
	return withHooks(ctx, frd.sqlExec, frd.mutation, frd.hooks)
}

// ExecX is like Exec, but panics if an error occurs.
func (frd *FileReferenceDelete) ExecX(ctx context.Context) int {
	n, err := frd.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (frd *FileReferenceDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := sqlgraph.NewDeleteSpec(filereference.Table, sqlgraph.NewFieldSpec(filereference.FieldID, field.TypeInt))
	if ps := frd.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	affected, err := sqlgraph.DeleteNodes(ctx, frd.driver, _spec)
	if err != nil && sqlgraph.IsConstraintError(err) {
		err = &ConstraintError{msg: err.Error(), wrap: err}
	}
	frd.mutation.done = true
	return affected, err
}

// FileReferenceDeleteOne is the builder for deleting a single FileReference entity.
type FileReferenceDeleteOne struct {
	frd *FileReferenceDelete
}

// Where appends a list predicates to the FileReferenceDelete builder.
func (frdo *FileReferenceDeleteOne) Where(ps ...predicate.FileReference) *FileReferenceDeleteOne {
	frdo.frd.mutation.Where(ps...)
	return frdo
}

// Exec executes the deletion query.
func (frdo *FileReferenceDeleteOne) Exec(ctx context.Context) error {
	n, err := frdo.frd.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{filereference.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (frdo *FileReferenceDeleteOne) ExecX(ctx context.Context) {
	if err := frdo.Exec(ctx); err != nil {
		panic(err)
	}
}
