// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package query

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"

	"gorm.io/gen"
	"gorm.io/gen/field"

	"gorm.io/plugin/dbresolver"

	"github.com/fivetentaylor/pointy/pkg/models"
)

func newStripeWebhookEvent(db *gorm.DB, opts ...gen.DOOption) stripeWebhookEvent {
	_stripeWebhookEvent := stripeWebhookEvent{}

	_stripeWebhookEvent.stripeWebhookEventDo.UseDB(db, opts...)
	_stripeWebhookEvent.stripeWebhookEventDo.UseModel(&models.StripeWebhookEvent{})

	tableName := _stripeWebhookEvent.stripeWebhookEventDo.TableName()
	_stripeWebhookEvent.ALL = field.NewAsterisk(tableName)
	_stripeWebhookEvent.ID = field.NewString(tableName, "id")
	_stripeWebhookEvent.EventID = field.NewString(tableName, "event_id")
	_stripeWebhookEvent.EventType = field.NewString(tableName, "event_type")
	_stripeWebhookEvent.Payload = field.NewString(tableName, "payload")
	_stripeWebhookEvent.ReceivedAt = field.NewTime(tableName, "received_at")
	_stripeWebhookEvent.Processed = field.NewBool(tableName, "processed")

	_stripeWebhookEvent.fillFieldMap()

	return _stripeWebhookEvent
}

type stripeWebhookEvent struct {
	stripeWebhookEventDo

	ALL        field.Asterisk
	ID         field.String
	EventID    field.String
	EventType  field.String
	Payload    field.String
	ReceivedAt field.Time
	Processed  field.Bool

	fieldMap map[string]field.Expr
}

func (s stripeWebhookEvent) Table(newTableName string) *stripeWebhookEvent {
	s.stripeWebhookEventDo.UseTable(newTableName)
	return s.updateTableName(newTableName)
}

func (s stripeWebhookEvent) As(alias string) *stripeWebhookEvent {
	s.stripeWebhookEventDo.DO = *(s.stripeWebhookEventDo.As(alias).(*gen.DO))
	return s.updateTableName(alias)
}

func (s *stripeWebhookEvent) updateTableName(table string) *stripeWebhookEvent {
	s.ALL = field.NewAsterisk(table)
	s.ID = field.NewString(table, "id")
	s.EventID = field.NewString(table, "event_id")
	s.EventType = field.NewString(table, "event_type")
	s.Payload = field.NewString(table, "payload")
	s.ReceivedAt = field.NewTime(table, "received_at")
	s.Processed = field.NewBool(table, "processed")

	s.fillFieldMap()

	return s
}

func (s *stripeWebhookEvent) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := s.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (s *stripeWebhookEvent) fillFieldMap() {
	s.fieldMap = make(map[string]field.Expr, 6)
	s.fieldMap["id"] = s.ID
	s.fieldMap["event_id"] = s.EventID
	s.fieldMap["event_type"] = s.EventType
	s.fieldMap["payload"] = s.Payload
	s.fieldMap["received_at"] = s.ReceivedAt
	s.fieldMap["processed"] = s.Processed
}

func (s stripeWebhookEvent) clone(db *gorm.DB) stripeWebhookEvent {
	s.stripeWebhookEventDo.ReplaceConnPool(db.Statement.ConnPool)
	return s
}

func (s stripeWebhookEvent) replaceDB(db *gorm.DB) stripeWebhookEvent {
	s.stripeWebhookEventDo.ReplaceDB(db)
	return s
}

type stripeWebhookEventDo struct{ gen.DO }

type IStripeWebhookEventDo interface {
	gen.SubQuery
	Debug() IStripeWebhookEventDo
	WithContext(ctx context.Context) IStripeWebhookEventDo
	WithResult(fc func(tx gen.Dao)) gen.ResultInfo
	ReplaceDB(db *gorm.DB)
	ReadDB() IStripeWebhookEventDo
	WriteDB() IStripeWebhookEventDo
	As(alias string) gen.Dao
	Session(config *gorm.Session) IStripeWebhookEventDo
	Columns(cols ...field.Expr) gen.Columns
	Clauses(conds ...clause.Expression) IStripeWebhookEventDo
	Not(conds ...gen.Condition) IStripeWebhookEventDo
	Or(conds ...gen.Condition) IStripeWebhookEventDo
	Select(conds ...field.Expr) IStripeWebhookEventDo
	Where(conds ...gen.Condition) IStripeWebhookEventDo
	Order(conds ...field.Expr) IStripeWebhookEventDo
	Distinct(cols ...field.Expr) IStripeWebhookEventDo
	Omit(cols ...field.Expr) IStripeWebhookEventDo
	Join(table schema.Tabler, on ...field.Expr) IStripeWebhookEventDo
	LeftJoin(table schema.Tabler, on ...field.Expr) IStripeWebhookEventDo
	RightJoin(table schema.Tabler, on ...field.Expr) IStripeWebhookEventDo
	Group(cols ...field.Expr) IStripeWebhookEventDo
	Having(conds ...gen.Condition) IStripeWebhookEventDo
	Limit(limit int) IStripeWebhookEventDo
	Offset(offset int) IStripeWebhookEventDo
	Count() (count int64, err error)
	Scopes(funcs ...func(gen.Dao) gen.Dao) IStripeWebhookEventDo
	Unscoped() IStripeWebhookEventDo
	Create(values ...*models.StripeWebhookEvent) error
	CreateInBatches(values []*models.StripeWebhookEvent, batchSize int) error
	Save(values ...*models.StripeWebhookEvent) error
	First() (*models.StripeWebhookEvent, error)
	Take() (*models.StripeWebhookEvent, error)
	Last() (*models.StripeWebhookEvent, error)
	Find() ([]*models.StripeWebhookEvent, error)
	FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*models.StripeWebhookEvent, err error)
	FindInBatches(result *[]*models.StripeWebhookEvent, batchSize int, fc func(tx gen.Dao, batch int) error) error
	Pluck(column field.Expr, dest interface{}) error
	Delete(...*models.StripeWebhookEvent) (info gen.ResultInfo, err error)
	Update(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	Updates(value interface{}) (info gen.ResultInfo, err error)
	UpdateColumn(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateColumnSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	UpdateColumns(value interface{}) (info gen.ResultInfo, err error)
	UpdateFrom(q gen.SubQuery) gen.Dao
	Attrs(attrs ...field.AssignExpr) IStripeWebhookEventDo
	Assign(attrs ...field.AssignExpr) IStripeWebhookEventDo
	Joins(fields ...field.RelationField) IStripeWebhookEventDo
	Preload(fields ...field.RelationField) IStripeWebhookEventDo
	FirstOrInit() (*models.StripeWebhookEvent, error)
	FirstOrCreate() (*models.StripeWebhookEvent, error)
	FindByPage(offset int, limit int) (result []*models.StripeWebhookEvent, count int64, err error)
	ScanByPage(result interface{}, offset int, limit int) (count int64, err error)
	Scan(result interface{}) (err error)
	Returning(value interface{}, columns ...string) IStripeWebhookEventDo
	UnderlyingDB() *gorm.DB
	schema.Tabler
}

func (s stripeWebhookEventDo) Debug() IStripeWebhookEventDo {
	return s.withDO(s.DO.Debug())
}

func (s stripeWebhookEventDo) WithContext(ctx context.Context) IStripeWebhookEventDo {
	return s.withDO(s.DO.WithContext(ctx))
}

func (s stripeWebhookEventDo) ReadDB() IStripeWebhookEventDo {
	return s.Clauses(dbresolver.Read)
}

func (s stripeWebhookEventDo) WriteDB() IStripeWebhookEventDo {
	return s.Clauses(dbresolver.Write)
}

func (s stripeWebhookEventDo) Session(config *gorm.Session) IStripeWebhookEventDo {
	return s.withDO(s.DO.Session(config))
}

func (s stripeWebhookEventDo) Clauses(conds ...clause.Expression) IStripeWebhookEventDo {
	return s.withDO(s.DO.Clauses(conds...))
}

func (s stripeWebhookEventDo) Returning(value interface{}, columns ...string) IStripeWebhookEventDo {
	return s.withDO(s.DO.Returning(value, columns...))
}

func (s stripeWebhookEventDo) Not(conds ...gen.Condition) IStripeWebhookEventDo {
	return s.withDO(s.DO.Not(conds...))
}

func (s stripeWebhookEventDo) Or(conds ...gen.Condition) IStripeWebhookEventDo {
	return s.withDO(s.DO.Or(conds...))
}

func (s stripeWebhookEventDo) Select(conds ...field.Expr) IStripeWebhookEventDo {
	return s.withDO(s.DO.Select(conds...))
}

func (s stripeWebhookEventDo) Where(conds ...gen.Condition) IStripeWebhookEventDo {
	return s.withDO(s.DO.Where(conds...))
}

func (s stripeWebhookEventDo) Order(conds ...field.Expr) IStripeWebhookEventDo {
	return s.withDO(s.DO.Order(conds...))
}

func (s stripeWebhookEventDo) Distinct(cols ...field.Expr) IStripeWebhookEventDo {
	return s.withDO(s.DO.Distinct(cols...))
}

func (s stripeWebhookEventDo) Omit(cols ...field.Expr) IStripeWebhookEventDo {
	return s.withDO(s.DO.Omit(cols...))
}

func (s stripeWebhookEventDo) Join(table schema.Tabler, on ...field.Expr) IStripeWebhookEventDo {
	return s.withDO(s.DO.Join(table, on...))
}

func (s stripeWebhookEventDo) LeftJoin(table schema.Tabler, on ...field.Expr) IStripeWebhookEventDo {
	return s.withDO(s.DO.LeftJoin(table, on...))
}

func (s stripeWebhookEventDo) RightJoin(table schema.Tabler, on ...field.Expr) IStripeWebhookEventDo {
	return s.withDO(s.DO.RightJoin(table, on...))
}

func (s stripeWebhookEventDo) Group(cols ...field.Expr) IStripeWebhookEventDo {
	return s.withDO(s.DO.Group(cols...))
}

func (s stripeWebhookEventDo) Having(conds ...gen.Condition) IStripeWebhookEventDo {
	return s.withDO(s.DO.Having(conds...))
}

func (s stripeWebhookEventDo) Limit(limit int) IStripeWebhookEventDo {
	return s.withDO(s.DO.Limit(limit))
}

func (s stripeWebhookEventDo) Offset(offset int) IStripeWebhookEventDo {
	return s.withDO(s.DO.Offset(offset))
}

func (s stripeWebhookEventDo) Scopes(funcs ...func(gen.Dao) gen.Dao) IStripeWebhookEventDo {
	return s.withDO(s.DO.Scopes(funcs...))
}

func (s stripeWebhookEventDo) Unscoped() IStripeWebhookEventDo {
	return s.withDO(s.DO.Unscoped())
}

func (s stripeWebhookEventDo) Create(values ...*models.StripeWebhookEvent) error {
	if len(values) == 0 {
		return nil
	}
	return s.DO.Create(values)
}

func (s stripeWebhookEventDo) CreateInBatches(values []*models.StripeWebhookEvent, batchSize int) error {
	return s.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (s stripeWebhookEventDo) Save(values ...*models.StripeWebhookEvent) error {
	if len(values) == 0 {
		return nil
	}
	return s.DO.Save(values)
}

func (s stripeWebhookEventDo) First() (*models.StripeWebhookEvent, error) {
	if result, err := s.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*models.StripeWebhookEvent), nil
	}
}

func (s stripeWebhookEventDo) Take() (*models.StripeWebhookEvent, error) {
	if result, err := s.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*models.StripeWebhookEvent), nil
	}
}

func (s stripeWebhookEventDo) Last() (*models.StripeWebhookEvent, error) {
	if result, err := s.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*models.StripeWebhookEvent), nil
	}
}

func (s stripeWebhookEventDo) Find() ([]*models.StripeWebhookEvent, error) {
	result, err := s.DO.Find()
	return result.([]*models.StripeWebhookEvent), err
}

func (s stripeWebhookEventDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*models.StripeWebhookEvent, err error) {
	buf := make([]*models.StripeWebhookEvent, 0, batchSize)
	err = s.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (s stripeWebhookEventDo) FindInBatches(result *[]*models.StripeWebhookEvent, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return s.DO.FindInBatches(result, batchSize, fc)
}

func (s stripeWebhookEventDo) Attrs(attrs ...field.AssignExpr) IStripeWebhookEventDo {
	return s.withDO(s.DO.Attrs(attrs...))
}

func (s stripeWebhookEventDo) Assign(attrs ...field.AssignExpr) IStripeWebhookEventDo {
	return s.withDO(s.DO.Assign(attrs...))
}

func (s stripeWebhookEventDo) Joins(fields ...field.RelationField) IStripeWebhookEventDo {
	for _, _f := range fields {
		s = *s.withDO(s.DO.Joins(_f))
	}
	return &s
}

func (s stripeWebhookEventDo) Preload(fields ...field.RelationField) IStripeWebhookEventDo {
	for _, _f := range fields {
		s = *s.withDO(s.DO.Preload(_f))
	}
	return &s
}

func (s stripeWebhookEventDo) FirstOrInit() (*models.StripeWebhookEvent, error) {
	if result, err := s.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*models.StripeWebhookEvent), nil
	}
}

func (s stripeWebhookEventDo) FirstOrCreate() (*models.StripeWebhookEvent, error) {
	if result, err := s.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*models.StripeWebhookEvent), nil
	}
}

func (s stripeWebhookEventDo) FindByPage(offset int, limit int) (result []*models.StripeWebhookEvent, count int64, err error) {
	result, err = s.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = s.Offset(-1).Limit(-1).Count()
	return
}

func (s stripeWebhookEventDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = s.Count()
	if err != nil {
		return
	}

	err = s.Offset(offset).Limit(limit).Scan(result)
	return
}

func (s stripeWebhookEventDo) Scan(result interface{}) (err error) {
	return s.DO.Scan(result)
}

func (s stripeWebhookEventDo) Delete(models ...*models.StripeWebhookEvent) (result gen.ResultInfo, err error) {
	return s.DO.Delete(models)
}

func (s *stripeWebhookEventDo) withDO(do gen.Dao) *stripeWebhookEventDo {
	s.DO = *do.(*gen.DO)
	return s
}
