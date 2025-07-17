package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// QRCodeAnalytics holds the schema definition for the QRCodeAnalytics entity.
type QRCodeAnalytics struct {
	ent.Schema
}

// Fields of the QRCodeAnalytics.
func (QRCodeAnalytics) Fields() []ent.Field {
	return []ent.Field{
		field.String("ip_address"),
		field.String("user_agent"),
		field.String("location").Optional(),
		field.String("device").Optional(),
		field.Time("scanned_at").Default(time.Now),
	}
}

// Edges of the QRCodeAnalytics.
func (QRCodeAnalytics) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("qr_code", QRCode.Type).Ref("analytics_records").Unique(),
	}
}
