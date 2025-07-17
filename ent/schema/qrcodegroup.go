package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// QRCodeGroup holds the schema definition for the QRCodeGroup entity.
type QRCodeGroup struct {
	ent.Schema
}

// Fields of the QRCodeGroup.
func (QRCodeGroup) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").NotEmpty(),
		field.String("description").Optional(),
		field.Time("created_at").Default(time.Now),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}

// Edges of the QRCodeGroup.
func (QRCodeGroup) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("qrcodes", QRCode.Type),
	}
}
