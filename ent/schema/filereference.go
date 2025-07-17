package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// FileReference holds the schema definition for the FileReference entity.
type FileReference struct {
	ent.Schema
}

// Fields of the FileReference.
func (FileReference) Fields() []ent.Field {
	return []ent.Field{
		field.String("filename").NotEmpty(),
		field.String("url").NotEmpty(),
		field.Int64("size"),
		field.String("type"),
	}
}

// Edges of the FileReference.
func (FileReference) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("qr_code", QRCode.Type).Ref("file_refs").Unique(),
	}
}
