package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// QRCode holds the schema definition for the QRCode entity.
type QRCode struct {
	ent.Schema
}

// Fields of the QRCode.
func (QRCode) Fields() []ent.Field {
	return []ent.Field{
		field.String("type").NotEmpty(),
		field.String("title").NotEmpty(),
		field.String("description").Optional(),
		field.String("redirect_url").Optional(),
		field.String("short_url").Optional(),
		field.JSON("content", map[string]interface{}{}),
		field.Time("created_at").Default(time.Now),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
		field.Time("expires_at").Optional().Nillable(),
		field.Bool("analytics").Default(false),
		field.Bool("active").Default(true),
		field.JSON("tags", []string{}).Optional(),
		field.JSON("design", map[string]interface{}{}).Optional(),
		field.Int("group_id").Optional().Nillable(),
	}
}

// Edges of the QRCode.
func (QRCode) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("file_refs", FileReference.Type),
		edge.From("group", QRCodeGroup.Type).Ref("qrcodes").Unique().Field("group_id"),
		edge.To("analytics_records", QRCodeAnalytics.Type),
	}
}
