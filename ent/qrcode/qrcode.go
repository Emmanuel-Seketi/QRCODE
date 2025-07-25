// Code generated by ent, DO NOT EDIT.

package qrcode

import (
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
)

const (
	// Label holds the string label denoting the qrcode type in the database.
	Label = "qr_code"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldType holds the string denoting the type field in the database.
	FieldType = "type"
	// FieldTitle holds the string denoting the title field in the database.
	FieldTitle = "title"
	// FieldDescription holds the string denoting the description field in the database.
	FieldDescription = "description"
	// FieldRedirectURL holds the string denoting the redirect_url field in the database.
	FieldRedirectURL = "redirect_url"
	// FieldShortURL holds the string denoting the short_url field in the database.
	FieldShortURL = "short_url"
	// FieldContent holds the string denoting the content field in the database.
	FieldContent = "content"
	// FieldCreatedAt holds the string denoting the created_at field in the database.
	FieldCreatedAt = "created_at"
	// FieldUpdatedAt holds the string denoting the updated_at field in the database.
	FieldUpdatedAt = "updated_at"
	// FieldExpiresAt holds the string denoting the expires_at field in the database.
	FieldExpiresAt = "expires_at"
	// FieldAnalytics holds the string denoting the analytics field in the database.
	FieldAnalytics = "analytics"
	// FieldActive holds the string denoting the active field in the database.
	FieldActive = "active"
	// FieldTags holds the string denoting the tags field in the database.
	FieldTags = "tags"
	// FieldDesign holds the string denoting the design field in the database.
	FieldDesign = "design"
	// FieldGroupID holds the string denoting the group_id field in the database.
	FieldGroupID = "group_id"
	// EdgeFileRefs holds the string denoting the file_refs edge name in mutations.
	EdgeFileRefs = "file_refs"
	// EdgeGroup holds the string denoting the group edge name in mutations.
	EdgeGroup = "group"
	// EdgeAnalyticsRecords holds the string denoting the analytics_records edge name in mutations.
	EdgeAnalyticsRecords = "analytics_records"
	// Table holds the table name of the qrcode in the database.
	Table = "qr_codes"
	// FileRefsTable is the table that holds the file_refs relation/edge.
	FileRefsTable = "file_references"
	// FileRefsInverseTable is the table name for the FileReference entity.
	// It exists in this package in order to avoid circular dependency with the "filereference" package.
	FileRefsInverseTable = "file_references"
	// FileRefsColumn is the table column denoting the file_refs relation/edge.
	FileRefsColumn = "qr_code_file_refs"
	// GroupTable is the table that holds the group relation/edge.
	GroupTable = "qr_codes"
	// GroupInverseTable is the table name for the QRCodeGroup entity.
	// It exists in this package in order to avoid circular dependency with the "qrcodegroup" package.
	GroupInverseTable = "qr_code_groups"
	// GroupColumn is the table column denoting the group relation/edge.
	GroupColumn = "group_id"
	// AnalyticsRecordsTable is the table that holds the analytics_records relation/edge.
	AnalyticsRecordsTable = "qr_code_analytics"
	// AnalyticsRecordsInverseTable is the table name for the QRCodeAnalytics entity.
	// It exists in this package in order to avoid circular dependency with the "qrcodeanalytics" package.
	AnalyticsRecordsInverseTable = "qr_code_analytics"
	// AnalyticsRecordsColumn is the table column denoting the analytics_records relation/edge.
	AnalyticsRecordsColumn = "qr_code_analytics_records"
)

// Columns holds all SQL columns for qrcode fields.
var Columns = []string{
	FieldID,
	FieldType,
	FieldTitle,
	FieldDescription,
	FieldRedirectURL,
	FieldShortURL,
	FieldContent,
	FieldCreatedAt,
	FieldUpdatedAt,
	FieldExpiresAt,
	FieldAnalytics,
	FieldActive,
	FieldTags,
	FieldDesign,
	FieldGroupID,
}

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	return false
}

var (
	// TypeValidator is a validator for the "type" field. It is called by the builders before save.
	TypeValidator func(string) error
	// TitleValidator is a validator for the "title" field. It is called by the builders before save.
	TitleValidator func(string) error
	// DefaultCreatedAt holds the default value on creation for the "created_at" field.
	DefaultCreatedAt func() time.Time
	// DefaultUpdatedAt holds the default value on creation for the "updated_at" field.
	DefaultUpdatedAt func() time.Time
	// UpdateDefaultUpdatedAt holds the default value on update for the "updated_at" field.
	UpdateDefaultUpdatedAt func() time.Time
	// DefaultAnalytics holds the default value on creation for the "analytics" field.
	DefaultAnalytics bool
	// DefaultActive holds the default value on creation for the "active" field.
	DefaultActive bool
)

// OrderOption defines the ordering options for the QRCode queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByType orders the results by the type field.
func ByType(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldType, opts...).ToFunc()
}

// ByTitle orders the results by the title field.
func ByTitle(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldTitle, opts...).ToFunc()
}

// ByDescription orders the results by the description field.
func ByDescription(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldDescription, opts...).ToFunc()
}

// ByRedirectURL orders the results by the redirect_url field.
func ByRedirectURL(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldRedirectURL, opts...).ToFunc()
}

// ByShortURL orders the results by the short_url field.
func ByShortURL(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldShortURL, opts...).ToFunc()
}

// ByCreatedAt orders the results by the created_at field.
func ByCreatedAt(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldCreatedAt, opts...).ToFunc()
}

// ByUpdatedAt orders the results by the updated_at field.
func ByUpdatedAt(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldUpdatedAt, opts...).ToFunc()
}

// ByExpiresAt orders the results by the expires_at field.
func ByExpiresAt(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldExpiresAt, opts...).ToFunc()
}

// ByAnalytics orders the results by the analytics field.
func ByAnalytics(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldAnalytics, opts...).ToFunc()
}

// ByActive orders the results by the active field.
func ByActive(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldActive, opts...).ToFunc()
}

// ByGroupID orders the results by the group_id field.
func ByGroupID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldGroupID, opts...).ToFunc()
}

// ByFileRefsCount orders the results by file_refs count.
func ByFileRefsCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newFileRefsStep(), opts...)
	}
}

// ByFileRefs orders the results by file_refs terms.
func ByFileRefs(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newFileRefsStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}

// ByGroupField orders the results by group field.
func ByGroupField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newGroupStep(), sql.OrderByField(field, opts...))
	}
}

// ByAnalyticsRecordsCount orders the results by analytics_records count.
func ByAnalyticsRecordsCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newAnalyticsRecordsStep(), opts...)
	}
}

// ByAnalyticsRecords orders the results by analytics_records terms.
func ByAnalyticsRecords(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newAnalyticsRecordsStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}
func newFileRefsStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(FileRefsInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, FileRefsTable, FileRefsColumn),
	)
}
func newGroupStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(GroupInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, GroupTable, GroupColumn),
	)
}
func newAnalyticsRecordsStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(AnalyticsRecordsInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, AnalyticsRecordsTable, AnalyticsRecordsColumn),
	)
}
