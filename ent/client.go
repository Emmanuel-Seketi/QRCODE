// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"log"
	"reflect"

	"qr_backend/ent/migrate"

	"qr_backend/ent/filereference"
	"qr_backend/ent/qrcode"
	"qr_backend/ent/qrcodeanalytics"
	"qr_backend/ent/qrcodegroup"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
)

// Client is the client that holds all ent builders.
type Client struct {
	config
	// Schema is the client for creating, migrating and dropping schema.
	Schema *migrate.Schema
	// FileReference is the client for interacting with the FileReference builders.
	FileReference *FileReferenceClient
	// QRCode is the client for interacting with the QRCode builders.
	QRCode *QRCodeClient
	// QRCodeAnalytics is the client for interacting with the QRCodeAnalytics builders.
	QRCodeAnalytics *QRCodeAnalyticsClient
	// QRCodeGroup is the client for interacting with the QRCodeGroup builders.
	QRCodeGroup *QRCodeGroupClient
}

// NewClient creates a new client configured with the given options.
func NewClient(opts ...Option) *Client {
	client := &Client{config: newConfig(opts...)}
	client.init()
	return client
}

func (c *Client) init() {
	c.Schema = migrate.NewSchema(c.driver)
	c.FileReference = NewFileReferenceClient(c.config)
	c.QRCode = NewQRCodeClient(c.config)
	c.QRCodeAnalytics = NewQRCodeAnalyticsClient(c.config)
	c.QRCodeGroup = NewQRCodeGroupClient(c.config)
}

type (
	// config is the configuration for the client and its builder.
	config struct {
		// driver used for executing database requests.
		driver dialect.Driver
		// debug enable a debug logging.
		debug bool
		// log used for logging on debug mode.
		log func(...any)
		// hooks to execute on mutations.
		hooks *hooks
		// interceptors to execute on queries.
		inters *inters
	}
	// Option function to configure the client.
	Option func(*config)
)

// newConfig creates a new config for the client.
func newConfig(opts ...Option) config {
	cfg := config{log: log.Println, hooks: &hooks{}, inters: &inters{}}
	cfg.options(opts...)
	return cfg
}

// options applies the options on the config object.
func (c *config) options(opts ...Option) {
	for _, opt := range opts {
		opt(c)
	}
	if c.debug {
		c.driver = dialect.Debug(c.driver, c.log)
	}
}

// Debug enables debug logging on the ent.Driver.
func Debug() Option {
	return func(c *config) {
		c.debug = true
	}
}

// Log sets the logging function for debug mode.
func Log(fn func(...any)) Option {
	return func(c *config) {
		c.log = fn
	}
}

// Driver configures the client driver.
func Driver(driver dialect.Driver) Option {
	return func(c *config) {
		c.driver = driver
	}
}

// Open opens a database/sql.DB specified by the driver name and
// the data source name, and returns a new client attached to it.
// Optional parameters can be added for configuring the client.
func Open(driverName, dataSourceName string, options ...Option) (*Client, error) {
	switch driverName {
	case dialect.MySQL, dialect.Postgres, dialect.SQLite:
		drv, err := sql.Open(driverName, dataSourceName)
		if err != nil {
			return nil, err
		}
		return NewClient(append(options, Driver(drv))...), nil
	default:
		return nil, fmt.Errorf("unsupported driver: %q", driverName)
	}
}

// ErrTxStarted is returned when trying to start a new transaction from a transactional client.
var ErrTxStarted = errors.New("ent: cannot start a transaction within a transaction")

// Tx returns a new transactional client. The provided context
// is used until the transaction is committed or rolled back.
func (c *Client) Tx(ctx context.Context) (*Tx, error) {
	if _, ok := c.driver.(*txDriver); ok {
		return nil, ErrTxStarted
	}
	tx, err := newTx(ctx, c.driver)
	if err != nil {
		return nil, fmt.Errorf("ent: starting a transaction: %w", err)
	}
	cfg := c.config
	cfg.driver = tx
	return &Tx{
		ctx:             ctx,
		config:          cfg,
		FileReference:   NewFileReferenceClient(cfg),
		QRCode:          NewQRCodeClient(cfg),
		QRCodeAnalytics: NewQRCodeAnalyticsClient(cfg),
		QRCodeGroup:     NewQRCodeGroupClient(cfg),
	}, nil
}

// BeginTx returns a transactional client with specified options.
func (c *Client) BeginTx(ctx context.Context, opts *sql.TxOptions) (*Tx, error) {
	if _, ok := c.driver.(*txDriver); ok {
		return nil, errors.New("ent: cannot start a transaction within a transaction")
	}
	tx, err := c.driver.(interface {
		BeginTx(context.Context, *sql.TxOptions) (dialect.Tx, error)
	}).BeginTx(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("ent: starting a transaction: %w", err)
	}
	cfg := c.config
	cfg.driver = &txDriver{tx: tx, drv: c.driver}
	return &Tx{
		ctx:             ctx,
		config:          cfg,
		FileReference:   NewFileReferenceClient(cfg),
		QRCode:          NewQRCodeClient(cfg),
		QRCodeAnalytics: NewQRCodeAnalyticsClient(cfg),
		QRCodeGroup:     NewQRCodeGroupClient(cfg),
	}, nil
}

// Debug returns a new debug-client. It's used to get verbose logging on specific operations.
//
//	client.Debug().
//		FileReference.
//		Query().
//		Count(ctx)
func (c *Client) Debug() *Client {
	if c.debug {
		return c
	}
	cfg := c.config
	cfg.driver = dialect.Debug(c.driver, c.log)
	client := &Client{config: cfg}
	client.init()
	return client
}

// Close closes the database connection and prevents new queries from starting.
func (c *Client) Close() error {
	return c.driver.Close()
}

// Use adds the mutation hooks to all the entity clients.
// In order to add hooks to a specific client, call: `client.Node.Use(...)`.
func (c *Client) Use(hooks ...Hook) {
	c.FileReference.Use(hooks...)
	c.QRCode.Use(hooks...)
	c.QRCodeAnalytics.Use(hooks...)
	c.QRCodeGroup.Use(hooks...)
}

// Intercept adds the query interceptors to all the entity clients.
// In order to add interceptors to a specific client, call: `client.Node.Intercept(...)`.
func (c *Client) Intercept(interceptors ...Interceptor) {
	c.FileReference.Intercept(interceptors...)
	c.QRCode.Intercept(interceptors...)
	c.QRCodeAnalytics.Intercept(interceptors...)
	c.QRCodeGroup.Intercept(interceptors...)
}

// Mutate implements the ent.Mutator interface.
func (c *Client) Mutate(ctx context.Context, m Mutation) (Value, error) {
	switch m := m.(type) {
	case *FileReferenceMutation:
		return c.FileReference.mutate(ctx, m)
	case *QRCodeMutation:
		return c.QRCode.mutate(ctx, m)
	case *QRCodeAnalyticsMutation:
		return c.QRCodeAnalytics.mutate(ctx, m)
	case *QRCodeGroupMutation:
		return c.QRCodeGroup.mutate(ctx, m)
	default:
		return nil, fmt.Errorf("ent: unknown mutation type %T", m)
	}
}

// FileReferenceClient is a client for the FileReference schema.
type FileReferenceClient struct {
	config
}

// NewFileReferenceClient returns a client for the FileReference from the given config.
func NewFileReferenceClient(c config) *FileReferenceClient {
	return &FileReferenceClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `filereference.Hooks(f(g(h())))`.
func (c *FileReferenceClient) Use(hooks ...Hook) {
	c.hooks.FileReference = append(c.hooks.FileReference, hooks...)
}

// Intercept adds a list of query interceptors to the interceptors stack.
// A call to `Intercept(f, g, h)` equals to `filereference.Intercept(f(g(h())))`.
func (c *FileReferenceClient) Intercept(interceptors ...Interceptor) {
	c.inters.FileReference = append(c.inters.FileReference, interceptors...)
}

// Create returns a builder for creating a FileReference entity.
func (c *FileReferenceClient) Create() *FileReferenceCreate {
	mutation := newFileReferenceMutation(c.config, OpCreate)
	return &FileReferenceCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// CreateBulk returns a builder for creating a bulk of FileReference entities.
func (c *FileReferenceClient) CreateBulk(builders ...*FileReferenceCreate) *FileReferenceCreateBulk {
	return &FileReferenceCreateBulk{config: c.config, builders: builders}
}

// MapCreateBulk creates a bulk creation builder from the given slice. For each item in the slice, the function creates
// a builder and applies setFunc on it.
func (c *FileReferenceClient) MapCreateBulk(slice any, setFunc func(*FileReferenceCreate, int)) *FileReferenceCreateBulk {
	rv := reflect.ValueOf(slice)
	if rv.Kind() != reflect.Slice {
		return &FileReferenceCreateBulk{err: fmt.Errorf("calling to FileReferenceClient.MapCreateBulk with wrong type %T, need slice", slice)}
	}
	builders := make([]*FileReferenceCreate, rv.Len())
	for i := 0; i < rv.Len(); i++ {
		builders[i] = c.Create()
		setFunc(builders[i], i)
	}
	return &FileReferenceCreateBulk{config: c.config, builders: builders}
}

// Update returns an update builder for FileReference.
func (c *FileReferenceClient) Update() *FileReferenceUpdate {
	mutation := newFileReferenceMutation(c.config, OpUpdate)
	return &FileReferenceUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *FileReferenceClient) UpdateOne(fr *FileReference) *FileReferenceUpdateOne {
	mutation := newFileReferenceMutation(c.config, OpUpdateOne, withFileReference(fr))
	return &FileReferenceUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOneID returns an update builder for the given id.
func (c *FileReferenceClient) UpdateOneID(id int) *FileReferenceUpdateOne {
	mutation := newFileReferenceMutation(c.config, OpUpdateOne, withFileReferenceID(id))
	return &FileReferenceUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for FileReference.
func (c *FileReferenceClient) Delete() *FileReferenceDelete {
	mutation := newFileReferenceMutation(c.config, OpDelete)
	return &FileReferenceDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a builder for deleting the given entity.
func (c *FileReferenceClient) DeleteOne(fr *FileReference) *FileReferenceDeleteOne {
	return c.DeleteOneID(fr.ID)
}

// DeleteOneID returns a builder for deleting the given entity by its id.
func (c *FileReferenceClient) DeleteOneID(id int) *FileReferenceDeleteOne {
	builder := c.Delete().Where(filereference.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &FileReferenceDeleteOne{builder}
}

// Query returns a query builder for FileReference.
func (c *FileReferenceClient) Query() *FileReferenceQuery {
	return &FileReferenceQuery{
		config: c.config,
		ctx:    &QueryContext{Type: TypeFileReference},
		inters: c.Interceptors(),
	}
}

// Get returns a FileReference entity by its id.
func (c *FileReferenceClient) Get(ctx context.Context, id int) (*FileReference, error) {
	return c.Query().Where(filereference.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *FileReferenceClient) GetX(ctx context.Context, id int) *FileReference {
	obj, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return obj
}

// QueryQrCode queries the qr_code edge of a FileReference.
func (c *FileReferenceClient) QueryQrCode(fr *FileReference) *QRCodeQuery {
	query := (&QRCodeClient{config: c.config}).Query()
	query.path = func(context.Context) (fromV *sql.Selector, _ error) {
		id := fr.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(filereference.Table, filereference.FieldID, id),
			sqlgraph.To(qrcode.Table, qrcode.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, filereference.QrCodeTable, filereference.QrCodeColumn),
		)
		fromV = sqlgraph.Neighbors(fr.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// Hooks returns the client hooks.
func (c *FileReferenceClient) Hooks() []Hook {
	return c.hooks.FileReference
}

// Interceptors returns the client interceptors.
func (c *FileReferenceClient) Interceptors() []Interceptor {
	return c.inters.FileReference
}

func (c *FileReferenceClient) mutate(ctx context.Context, m *FileReferenceMutation) (Value, error) {
	switch m.Op() {
	case OpCreate:
		return (&FileReferenceCreate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdate:
		return (&FileReferenceUpdate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdateOne:
		return (&FileReferenceUpdateOne{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpDelete, OpDeleteOne:
		return (&FileReferenceDelete{config: c.config, hooks: c.Hooks(), mutation: m}).Exec(ctx)
	default:
		return nil, fmt.Errorf("ent: unknown FileReference mutation op: %q", m.Op())
	}
}

// QRCodeClient is a client for the QRCode schema.
type QRCodeClient struct {
	config
}

// NewQRCodeClient returns a client for the QRCode from the given config.
func NewQRCodeClient(c config) *QRCodeClient {
	return &QRCodeClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `qrcode.Hooks(f(g(h())))`.
func (c *QRCodeClient) Use(hooks ...Hook) {
	c.hooks.QRCode = append(c.hooks.QRCode, hooks...)
}

// Intercept adds a list of query interceptors to the interceptors stack.
// A call to `Intercept(f, g, h)` equals to `qrcode.Intercept(f(g(h())))`.
func (c *QRCodeClient) Intercept(interceptors ...Interceptor) {
	c.inters.QRCode = append(c.inters.QRCode, interceptors...)
}

// Create returns a builder for creating a QRCode entity.
func (c *QRCodeClient) Create() *QRCodeCreate {
	mutation := newQRCodeMutation(c.config, OpCreate)
	return &QRCodeCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// CreateBulk returns a builder for creating a bulk of QRCode entities.
func (c *QRCodeClient) CreateBulk(builders ...*QRCodeCreate) *QRCodeCreateBulk {
	return &QRCodeCreateBulk{config: c.config, builders: builders}
}

// MapCreateBulk creates a bulk creation builder from the given slice. For each item in the slice, the function creates
// a builder and applies setFunc on it.
func (c *QRCodeClient) MapCreateBulk(slice any, setFunc func(*QRCodeCreate, int)) *QRCodeCreateBulk {
	rv := reflect.ValueOf(slice)
	if rv.Kind() != reflect.Slice {
		return &QRCodeCreateBulk{err: fmt.Errorf("calling to QRCodeClient.MapCreateBulk with wrong type %T, need slice", slice)}
	}
	builders := make([]*QRCodeCreate, rv.Len())
	for i := 0; i < rv.Len(); i++ {
		builders[i] = c.Create()
		setFunc(builders[i], i)
	}
	return &QRCodeCreateBulk{config: c.config, builders: builders}
}

// Update returns an update builder for QRCode.
func (c *QRCodeClient) Update() *QRCodeUpdate {
	mutation := newQRCodeMutation(c.config, OpUpdate)
	return &QRCodeUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *QRCodeClient) UpdateOne(qc *QRCode) *QRCodeUpdateOne {
	mutation := newQRCodeMutation(c.config, OpUpdateOne, withQRCode(qc))
	return &QRCodeUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOneID returns an update builder for the given id.
func (c *QRCodeClient) UpdateOneID(id int) *QRCodeUpdateOne {
	mutation := newQRCodeMutation(c.config, OpUpdateOne, withQRCodeID(id))
	return &QRCodeUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for QRCode.
func (c *QRCodeClient) Delete() *QRCodeDelete {
	mutation := newQRCodeMutation(c.config, OpDelete)
	return &QRCodeDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a builder for deleting the given entity.
func (c *QRCodeClient) DeleteOne(qc *QRCode) *QRCodeDeleteOne {
	return c.DeleteOneID(qc.ID)
}

// DeleteOneID returns a builder for deleting the given entity by its id.
func (c *QRCodeClient) DeleteOneID(id int) *QRCodeDeleteOne {
	builder := c.Delete().Where(qrcode.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &QRCodeDeleteOne{builder}
}

// Query returns a query builder for QRCode.
func (c *QRCodeClient) Query() *QRCodeQuery {
	return &QRCodeQuery{
		config: c.config,
		ctx:    &QueryContext{Type: TypeQRCode},
		inters: c.Interceptors(),
	}
}

// Get returns a QRCode entity by its id.
func (c *QRCodeClient) Get(ctx context.Context, id int) (*QRCode, error) {
	return c.Query().Where(qrcode.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *QRCodeClient) GetX(ctx context.Context, id int) *QRCode {
	obj, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return obj
}

// QueryFileRefs queries the file_refs edge of a QRCode.
func (c *QRCodeClient) QueryFileRefs(qc *QRCode) *FileReferenceQuery {
	query := (&FileReferenceClient{config: c.config}).Query()
	query.path = func(context.Context) (fromV *sql.Selector, _ error) {
		id := qc.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(qrcode.Table, qrcode.FieldID, id),
			sqlgraph.To(filereference.Table, filereference.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, qrcode.FileRefsTable, qrcode.FileRefsColumn),
		)
		fromV = sqlgraph.Neighbors(qc.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryGroup queries the group edge of a QRCode.
func (c *QRCodeClient) QueryGroup(qc *QRCode) *QRCodeGroupQuery {
	query := (&QRCodeGroupClient{config: c.config}).Query()
	query.path = func(context.Context) (fromV *sql.Selector, _ error) {
		id := qc.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(qrcode.Table, qrcode.FieldID, id),
			sqlgraph.To(qrcodegroup.Table, qrcodegroup.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, qrcode.GroupTable, qrcode.GroupColumn),
		)
		fromV = sqlgraph.Neighbors(qc.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryAnalyticsRecords queries the analytics_records edge of a QRCode.
func (c *QRCodeClient) QueryAnalyticsRecords(qc *QRCode) *QRCodeAnalyticsQuery {
	query := (&QRCodeAnalyticsClient{config: c.config}).Query()
	query.path = func(context.Context) (fromV *sql.Selector, _ error) {
		id := qc.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(qrcode.Table, qrcode.FieldID, id),
			sqlgraph.To(qrcodeanalytics.Table, qrcodeanalytics.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, qrcode.AnalyticsRecordsTable, qrcode.AnalyticsRecordsColumn),
		)
		fromV = sqlgraph.Neighbors(qc.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// Hooks returns the client hooks.
func (c *QRCodeClient) Hooks() []Hook {
	return c.hooks.QRCode
}

// Interceptors returns the client interceptors.
func (c *QRCodeClient) Interceptors() []Interceptor {
	return c.inters.QRCode
}

func (c *QRCodeClient) mutate(ctx context.Context, m *QRCodeMutation) (Value, error) {
	switch m.Op() {
	case OpCreate:
		return (&QRCodeCreate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdate:
		return (&QRCodeUpdate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdateOne:
		return (&QRCodeUpdateOne{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpDelete, OpDeleteOne:
		return (&QRCodeDelete{config: c.config, hooks: c.Hooks(), mutation: m}).Exec(ctx)
	default:
		return nil, fmt.Errorf("ent: unknown QRCode mutation op: %q", m.Op())
	}
}

// QRCodeAnalyticsClient is a client for the QRCodeAnalytics schema.
type QRCodeAnalyticsClient struct {
	config
}

// NewQRCodeAnalyticsClient returns a client for the QRCodeAnalytics from the given config.
func NewQRCodeAnalyticsClient(c config) *QRCodeAnalyticsClient {
	return &QRCodeAnalyticsClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `qrcodeanalytics.Hooks(f(g(h())))`.
func (c *QRCodeAnalyticsClient) Use(hooks ...Hook) {
	c.hooks.QRCodeAnalytics = append(c.hooks.QRCodeAnalytics, hooks...)
}

// Intercept adds a list of query interceptors to the interceptors stack.
// A call to `Intercept(f, g, h)` equals to `qrcodeanalytics.Intercept(f(g(h())))`.
func (c *QRCodeAnalyticsClient) Intercept(interceptors ...Interceptor) {
	c.inters.QRCodeAnalytics = append(c.inters.QRCodeAnalytics, interceptors...)
}

// Create returns a builder for creating a QRCodeAnalytics entity.
func (c *QRCodeAnalyticsClient) Create() *QRCodeAnalyticsCreate {
	mutation := newQRCodeAnalyticsMutation(c.config, OpCreate)
	return &QRCodeAnalyticsCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// CreateBulk returns a builder for creating a bulk of QRCodeAnalytics entities.
func (c *QRCodeAnalyticsClient) CreateBulk(builders ...*QRCodeAnalyticsCreate) *QRCodeAnalyticsCreateBulk {
	return &QRCodeAnalyticsCreateBulk{config: c.config, builders: builders}
}

// MapCreateBulk creates a bulk creation builder from the given slice. For each item in the slice, the function creates
// a builder and applies setFunc on it.
func (c *QRCodeAnalyticsClient) MapCreateBulk(slice any, setFunc func(*QRCodeAnalyticsCreate, int)) *QRCodeAnalyticsCreateBulk {
	rv := reflect.ValueOf(slice)
	if rv.Kind() != reflect.Slice {
		return &QRCodeAnalyticsCreateBulk{err: fmt.Errorf("calling to QRCodeAnalyticsClient.MapCreateBulk with wrong type %T, need slice", slice)}
	}
	builders := make([]*QRCodeAnalyticsCreate, rv.Len())
	for i := 0; i < rv.Len(); i++ {
		builders[i] = c.Create()
		setFunc(builders[i], i)
	}
	return &QRCodeAnalyticsCreateBulk{config: c.config, builders: builders}
}

// Update returns an update builder for QRCodeAnalytics.
func (c *QRCodeAnalyticsClient) Update() *QRCodeAnalyticsUpdate {
	mutation := newQRCodeAnalyticsMutation(c.config, OpUpdate)
	return &QRCodeAnalyticsUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *QRCodeAnalyticsClient) UpdateOne(qca *QRCodeAnalytics) *QRCodeAnalyticsUpdateOne {
	mutation := newQRCodeAnalyticsMutation(c.config, OpUpdateOne, withQRCodeAnalytics(qca))
	return &QRCodeAnalyticsUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOneID returns an update builder for the given id.
func (c *QRCodeAnalyticsClient) UpdateOneID(id int) *QRCodeAnalyticsUpdateOne {
	mutation := newQRCodeAnalyticsMutation(c.config, OpUpdateOne, withQRCodeAnalyticsID(id))
	return &QRCodeAnalyticsUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for QRCodeAnalytics.
func (c *QRCodeAnalyticsClient) Delete() *QRCodeAnalyticsDelete {
	mutation := newQRCodeAnalyticsMutation(c.config, OpDelete)
	return &QRCodeAnalyticsDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a builder for deleting the given entity.
func (c *QRCodeAnalyticsClient) DeleteOne(qca *QRCodeAnalytics) *QRCodeAnalyticsDeleteOne {
	return c.DeleteOneID(qca.ID)
}

// DeleteOneID returns a builder for deleting the given entity by its id.
func (c *QRCodeAnalyticsClient) DeleteOneID(id int) *QRCodeAnalyticsDeleteOne {
	builder := c.Delete().Where(qrcodeanalytics.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &QRCodeAnalyticsDeleteOne{builder}
}

// Query returns a query builder for QRCodeAnalytics.
func (c *QRCodeAnalyticsClient) Query() *QRCodeAnalyticsQuery {
	return &QRCodeAnalyticsQuery{
		config: c.config,
		ctx:    &QueryContext{Type: TypeQRCodeAnalytics},
		inters: c.Interceptors(),
	}
}

// Get returns a QRCodeAnalytics entity by its id.
func (c *QRCodeAnalyticsClient) Get(ctx context.Context, id int) (*QRCodeAnalytics, error) {
	return c.Query().Where(qrcodeanalytics.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *QRCodeAnalyticsClient) GetX(ctx context.Context, id int) *QRCodeAnalytics {
	obj, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return obj
}

// QueryQrCode queries the qr_code edge of a QRCodeAnalytics.
func (c *QRCodeAnalyticsClient) QueryQrCode(qca *QRCodeAnalytics) *QRCodeQuery {
	query := (&QRCodeClient{config: c.config}).Query()
	query.path = func(context.Context) (fromV *sql.Selector, _ error) {
		id := qca.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(qrcodeanalytics.Table, qrcodeanalytics.FieldID, id),
			sqlgraph.To(qrcode.Table, qrcode.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, qrcodeanalytics.QrCodeTable, qrcodeanalytics.QrCodeColumn),
		)
		fromV = sqlgraph.Neighbors(qca.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// Hooks returns the client hooks.
func (c *QRCodeAnalyticsClient) Hooks() []Hook {
	return c.hooks.QRCodeAnalytics
}

// Interceptors returns the client interceptors.
func (c *QRCodeAnalyticsClient) Interceptors() []Interceptor {
	return c.inters.QRCodeAnalytics
}

func (c *QRCodeAnalyticsClient) mutate(ctx context.Context, m *QRCodeAnalyticsMutation) (Value, error) {
	switch m.Op() {
	case OpCreate:
		return (&QRCodeAnalyticsCreate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdate:
		return (&QRCodeAnalyticsUpdate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdateOne:
		return (&QRCodeAnalyticsUpdateOne{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpDelete, OpDeleteOne:
		return (&QRCodeAnalyticsDelete{config: c.config, hooks: c.Hooks(), mutation: m}).Exec(ctx)
	default:
		return nil, fmt.Errorf("ent: unknown QRCodeAnalytics mutation op: %q", m.Op())
	}
}

// QRCodeGroupClient is a client for the QRCodeGroup schema.
type QRCodeGroupClient struct {
	config
}

// NewQRCodeGroupClient returns a client for the QRCodeGroup from the given config.
func NewQRCodeGroupClient(c config) *QRCodeGroupClient {
	return &QRCodeGroupClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `qrcodegroup.Hooks(f(g(h())))`.
func (c *QRCodeGroupClient) Use(hooks ...Hook) {
	c.hooks.QRCodeGroup = append(c.hooks.QRCodeGroup, hooks...)
}

// Intercept adds a list of query interceptors to the interceptors stack.
// A call to `Intercept(f, g, h)` equals to `qrcodegroup.Intercept(f(g(h())))`.
func (c *QRCodeGroupClient) Intercept(interceptors ...Interceptor) {
	c.inters.QRCodeGroup = append(c.inters.QRCodeGroup, interceptors...)
}

// Create returns a builder for creating a QRCodeGroup entity.
func (c *QRCodeGroupClient) Create() *QRCodeGroupCreate {
	mutation := newQRCodeGroupMutation(c.config, OpCreate)
	return &QRCodeGroupCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// CreateBulk returns a builder for creating a bulk of QRCodeGroup entities.
func (c *QRCodeGroupClient) CreateBulk(builders ...*QRCodeGroupCreate) *QRCodeGroupCreateBulk {
	return &QRCodeGroupCreateBulk{config: c.config, builders: builders}
}

// MapCreateBulk creates a bulk creation builder from the given slice. For each item in the slice, the function creates
// a builder and applies setFunc on it.
func (c *QRCodeGroupClient) MapCreateBulk(slice any, setFunc func(*QRCodeGroupCreate, int)) *QRCodeGroupCreateBulk {
	rv := reflect.ValueOf(slice)
	if rv.Kind() != reflect.Slice {
		return &QRCodeGroupCreateBulk{err: fmt.Errorf("calling to QRCodeGroupClient.MapCreateBulk with wrong type %T, need slice", slice)}
	}
	builders := make([]*QRCodeGroupCreate, rv.Len())
	for i := 0; i < rv.Len(); i++ {
		builders[i] = c.Create()
		setFunc(builders[i], i)
	}
	return &QRCodeGroupCreateBulk{config: c.config, builders: builders}
}

// Update returns an update builder for QRCodeGroup.
func (c *QRCodeGroupClient) Update() *QRCodeGroupUpdate {
	mutation := newQRCodeGroupMutation(c.config, OpUpdate)
	return &QRCodeGroupUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *QRCodeGroupClient) UpdateOne(qcg *QRCodeGroup) *QRCodeGroupUpdateOne {
	mutation := newQRCodeGroupMutation(c.config, OpUpdateOne, withQRCodeGroup(qcg))
	return &QRCodeGroupUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOneID returns an update builder for the given id.
func (c *QRCodeGroupClient) UpdateOneID(id int) *QRCodeGroupUpdateOne {
	mutation := newQRCodeGroupMutation(c.config, OpUpdateOne, withQRCodeGroupID(id))
	return &QRCodeGroupUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for QRCodeGroup.
func (c *QRCodeGroupClient) Delete() *QRCodeGroupDelete {
	mutation := newQRCodeGroupMutation(c.config, OpDelete)
	return &QRCodeGroupDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a builder for deleting the given entity.
func (c *QRCodeGroupClient) DeleteOne(qcg *QRCodeGroup) *QRCodeGroupDeleteOne {
	return c.DeleteOneID(qcg.ID)
}

// DeleteOneID returns a builder for deleting the given entity by its id.
func (c *QRCodeGroupClient) DeleteOneID(id int) *QRCodeGroupDeleteOne {
	builder := c.Delete().Where(qrcodegroup.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &QRCodeGroupDeleteOne{builder}
}

// Query returns a query builder for QRCodeGroup.
func (c *QRCodeGroupClient) Query() *QRCodeGroupQuery {
	return &QRCodeGroupQuery{
		config: c.config,
		ctx:    &QueryContext{Type: TypeQRCodeGroup},
		inters: c.Interceptors(),
	}
}

// Get returns a QRCodeGroup entity by its id.
func (c *QRCodeGroupClient) Get(ctx context.Context, id int) (*QRCodeGroup, error) {
	return c.Query().Where(qrcodegroup.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *QRCodeGroupClient) GetX(ctx context.Context, id int) *QRCodeGroup {
	obj, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return obj
}

// QueryQrcodes queries the qrcodes edge of a QRCodeGroup.
func (c *QRCodeGroupClient) QueryQrcodes(qcg *QRCodeGroup) *QRCodeQuery {
	query := (&QRCodeClient{config: c.config}).Query()
	query.path = func(context.Context) (fromV *sql.Selector, _ error) {
		id := qcg.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(qrcodegroup.Table, qrcodegroup.FieldID, id),
			sqlgraph.To(qrcode.Table, qrcode.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, qrcodegroup.QrcodesTable, qrcodegroup.QrcodesColumn),
		)
		fromV = sqlgraph.Neighbors(qcg.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// Hooks returns the client hooks.
func (c *QRCodeGroupClient) Hooks() []Hook {
	return c.hooks.QRCodeGroup
}

// Interceptors returns the client interceptors.
func (c *QRCodeGroupClient) Interceptors() []Interceptor {
	return c.inters.QRCodeGroup
}

func (c *QRCodeGroupClient) mutate(ctx context.Context, m *QRCodeGroupMutation) (Value, error) {
	switch m.Op() {
	case OpCreate:
		return (&QRCodeGroupCreate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdate:
		return (&QRCodeGroupUpdate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdateOne:
		return (&QRCodeGroupUpdateOne{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpDelete, OpDeleteOne:
		return (&QRCodeGroupDelete{config: c.config, hooks: c.Hooks(), mutation: m}).Exec(ctx)
	default:
		return nil, fmt.Errorf("ent: unknown QRCodeGroup mutation op: %q", m.Op())
	}
}

// hooks and interceptors per client, for fast access.
type (
	hooks struct {
		FileReference, QRCode, QRCodeAnalytics, QRCodeGroup []ent.Hook
	}
	inters struct {
		FileReference, QRCode, QRCodeAnalytics, QRCodeGroup []ent.Interceptor
	}
)
