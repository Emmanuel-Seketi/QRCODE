// Code generated by ent, DO NOT EDIT.

package migrate

import (
	"entgo.io/ent/dialect/sql/schema"
	"entgo.io/ent/schema/field"
)

var (
	// FileReferencesColumns holds the columns for the "file_references" table.
	FileReferencesColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "filename", Type: field.TypeString},
		{Name: "url", Type: field.TypeString},
		{Name: "size", Type: field.TypeInt64},
		{Name: "type", Type: field.TypeString},
		{Name: "qr_code_file_refs", Type: field.TypeInt, Nullable: true},
	}
	// FileReferencesTable holds the schema information for the "file_references" table.
	FileReferencesTable = &schema.Table{
		Name:       "file_references",
		Columns:    FileReferencesColumns,
		PrimaryKey: []*schema.Column{FileReferencesColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "file_references_qr_codes_file_refs",
				Columns:    []*schema.Column{FileReferencesColumns[5]},
				RefColumns: []*schema.Column{QrCodesColumns[0]},
				OnDelete:   schema.SetNull,
			},
		},
	}
	// QrCodesColumns holds the columns for the "qr_codes" table.
	QrCodesColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "type", Type: field.TypeString},
		{Name: "title", Type: field.TypeString},
		{Name: "description", Type: field.TypeString, Nullable: true},
		{Name: "redirect_url", Type: field.TypeString, Nullable: true},
		{Name: "short_url", Type: field.TypeString, Nullable: true},
		{Name: "content", Type: field.TypeJSON},
		{Name: "created_at", Type: field.TypeTime},
		{Name: "updated_at", Type: field.TypeTime},
		{Name: "expires_at", Type: field.TypeTime, Nullable: true},
		{Name: "analytics", Type: field.TypeBool, Default: false},
		{Name: "active", Type: field.TypeBool, Default: true},
		{Name: "tags", Type: field.TypeJSON, Nullable: true},
		{Name: "design", Type: field.TypeJSON, Nullable: true},
		{Name: "group_id", Type: field.TypeInt, Nullable: true},
	}
	// QrCodesTable holds the schema information for the "qr_codes" table.
	QrCodesTable = &schema.Table{
		Name:       "qr_codes",
		Columns:    QrCodesColumns,
		PrimaryKey: []*schema.Column{QrCodesColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "qr_codes_qr_code_groups_qrcodes",
				Columns:    []*schema.Column{QrCodesColumns[14]},
				RefColumns: []*schema.Column{QrCodeGroupsColumns[0]},
				OnDelete:   schema.SetNull,
			},
		},
	}
	// QrCodeAnalyticsColumns holds the columns for the "qr_code_analytics" table.
	QrCodeAnalyticsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "ip_address", Type: field.TypeString},
		{Name: "user_agent", Type: field.TypeString},
		{Name: "location", Type: field.TypeString, Nullable: true},
		{Name: "device", Type: field.TypeString, Nullable: true},
		{Name: "scanned_at", Type: field.TypeTime},
		{Name: "qr_code_analytics_records", Type: field.TypeInt, Nullable: true},
	}
	// QrCodeAnalyticsTable holds the schema information for the "qr_code_analytics" table.
	QrCodeAnalyticsTable = &schema.Table{
		Name:       "qr_code_analytics",
		Columns:    QrCodeAnalyticsColumns,
		PrimaryKey: []*schema.Column{QrCodeAnalyticsColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "qr_code_analytics_qr_codes_analytics_records",
				Columns:    []*schema.Column{QrCodeAnalyticsColumns[6]},
				RefColumns: []*schema.Column{QrCodesColumns[0]},
				OnDelete:   schema.SetNull,
			},
		},
	}
	// QrCodeGroupsColumns holds the columns for the "qr_code_groups" table.
	QrCodeGroupsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "name", Type: field.TypeString},
		{Name: "description", Type: field.TypeString, Nullable: true},
		{Name: "created_at", Type: field.TypeTime},
		{Name: "updated_at", Type: field.TypeTime},
	}
	// QrCodeGroupsTable holds the schema information for the "qr_code_groups" table.
	QrCodeGroupsTable = &schema.Table{
		Name:       "qr_code_groups",
		Columns:    QrCodeGroupsColumns,
		PrimaryKey: []*schema.Column{QrCodeGroupsColumns[0]},
	}
	// Tables holds all the tables in the schema.
	Tables = []*schema.Table{
		FileReferencesTable,
		QrCodesTable,
		QrCodeAnalyticsTable,
		QrCodeGroupsTable,
	}
)

func init() {
	FileReferencesTable.ForeignKeys[0].RefTable = QrCodesTable
	QrCodesTable.ForeignKeys[0].RefTable = QrCodeGroupsTable
	QrCodeAnalyticsTable.ForeignKeys[0].RefTable = QrCodesTable
}
