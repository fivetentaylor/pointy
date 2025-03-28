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

func newDocumentVersion(db *gorm.DB, opts ...gen.DOOption) documentVersion {
	_documentVersion := documentVersion{}

	_documentVersion.documentVersionDo.UseDB(db, opts...)
	_documentVersion.documentVersionDo.UseModel(&models.DocumentVersion{})

	tableName := _documentVersion.documentVersionDo.TableName()
	_documentVersion.ALL = field.NewAsterisk(tableName)
	_documentVersion.ID = field.NewString(tableName, "id")
	_documentVersion.DocumentID = field.NewString(tableName, "document_id")
	_documentVersion.Name = field.NewString(tableName, "name")
	_documentVersion.ContentAddress = field.NewString(tableName, "content_address")
	_documentVersion.CreatedAt = field.NewTime(tableName, "created_at")
	_documentVersion.UpdatedAt = field.NewTime(tableName, "updated_at")
	_documentVersion.CreatedBy = field.NewString(tableName, "created_by")
	_documentVersion.UpdatedBy = field.NewString(tableName, "updated_by")

	_documentVersion.fillFieldMap()

	return _documentVersion
}

type documentVersion struct {
	documentVersionDo

	ALL            field.Asterisk
	ID             field.String
	DocumentID     field.String
	Name           field.String
	ContentAddress field.String
	CreatedAt      field.Time
	UpdatedAt      field.Time
	CreatedBy      field.String
	UpdatedBy      field.String

	fieldMap map[string]field.Expr
}

func (d documentVersion) Table(newTableName string) *documentVersion {
	d.documentVersionDo.UseTable(newTableName)
	return d.updateTableName(newTableName)
}

func (d documentVersion) As(alias string) *documentVersion {
	d.documentVersionDo.DO = *(d.documentVersionDo.As(alias).(*gen.DO))
	return d.updateTableName(alias)
}

func (d *documentVersion) updateTableName(table string) *documentVersion {
	d.ALL = field.NewAsterisk(table)
	d.ID = field.NewString(table, "id")
	d.DocumentID = field.NewString(table, "document_id")
	d.Name = field.NewString(table, "name")
	d.ContentAddress = field.NewString(table, "content_address")
	d.CreatedAt = field.NewTime(table, "created_at")
	d.UpdatedAt = field.NewTime(table, "updated_at")
	d.CreatedBy = field.NewString(table, "created_by")
	d.UpdatedBy = field.NewString(table, "updated_by")

	d.fillFieldMap()

	return d
}

func (d *documentVersion) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := d.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (d *documentVersion) fillFieldMap() {
	d.fieldMap = make(map[string]field.Expr, 8)
	d.fieldMap["id"] = d.ID
	d.fieldMap["document_id"] = d.DocumentID
	d.fieldMap["name"] = d.Name
	d.fieldMap["content_address"] = d.ContentAddress
	d.fieldMap["created_at"] = d.CreatedAt
	d.fieldMap["updated_at"] = d.UpdatedAt
	d.fieldMap["created_by"] = d.CreatedBy
	d.fieldMap["updated_by"] = d.UpdatedBy
}

func (d documentVersion) clone(db *gorm.DB) documentVersion {
	d.documentVersionDo.ReplaceConnPool(db.Statement.ConnPool)
	return d
}

func (d documentVersion) replaceDB(db *gorm.DB) documentVersion {
	d.documentVersionDo.ReplaceDB(db)
	return d
}

type documentVersionDo struct{ gen.DO }

type IDocumentVersionDo interface {
	gen.SubQuery
	Debug() IDocumentVersionDo
	WithContext(ctx context.Context) IDocumentVersionDo
	WithResult(fc func(tx gen.Dao)) gen.ResultInfo
	ReplaceDB(db *gorm.DB)
	ReadDB() IDocumentVersionDo
	WriteDB() IDocumentVersionDo
	As(alias string) gen.Dao
	Session(config *gorm.Session) IDocumentVersionDo
	Columns(cols ...field.Expr) gen.Columns
	Clauses(conds ...clause.Expression) IDocumentVersionDo
	Not(conds ...gen.Condition) IDocumentVersionDo
	Or(conds ...gen.Condition) IDocumentVersionDo
	Select(conds ...field.Expr) IDocumentVersionDo
	Where(conds ...gen.Condition) IDocumentVersionDo
	Order(conds ...field.Expr) IDocumentVersionDo
	Distinct(cols ...field.Expr) IDocumentVersionDo
	Omit(cols ...field.Expr) IDocumentVersionDo
	Join(table schema.Tabler, on ...field.Expr) IDocumentVersionDo
	LeftJoin(table schema.Tabler, on ...field.Expr) IDocumentVersionDo
	RightJoin(table schema.Tabler, on ...field.Expr) IDocumentVersionDo
	Group(cols ...field.Expr) IDocumentVersionDo
	Having(conds ...gen.Condition) IDocumentVersionDo
	Limit(limit int) IDocumentVersionDo
	Offset(offset int) IDocumentVersionDo
	Count() (count int64, err error)
	Scopes(funcs ...func(gen.Dao) gen.Dao) IDocumentVersionDo
	Unscoped() IDocumentVersionDo
	Create(values ...*models.DocumentVersion) error
	CreateInBatches(values []*models.DocumentVersion, batchSize int) error
	Save(values ...*models.DocumentVersion) error
	First() (*models.DocumentVersion, error)
	Take() (*models.DocumentVersion, error)
	Last() (*models.DocumentVersion, error)
	Find() ([]*models.DocumentVersion, error)
	FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*models.DocumentVersion, err error)
	FindInBatches(result *[]*models.DocumentVersion, batchSize int, fc func(tx gen.Dao, batch int) error) error
	Pluck(column field.Expr, dest interface{}) error
	Delete(...*models.DocumentVersion) (info gen.ResultInfo, err error)
	Update(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	Updates(value interface{}) (info gen.ResultInfo, err error)
	UpdateColumn(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateColumnSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	UpdateColumns(value interface{}) (info gen.ResultInfo, err error)
	UpdateFrom(q gen.SubQuery) gen.Dao
	Attrs(attrs ...field.AssignExpr) IDocumentVersionDo
	Assign(attrs ...field.AssignExpr) IDocumentVersionDo
	Joins(fields ...field.RelationField) IDocumentVersionDo
	Preload(fields ...field.RelationField) IDocumentVersionDo
	FirstOrInit() (*models.DocumentVersion, error)
	FirstOrCreate() (*models.DocumentVersion, error)
	FindByPage(offset int, limit int) (result []*models.DocumentVersion, count int64, err error)
	ScanByPage(result interface{}, offset int, limit int) (count int64, err error)
	Scan(result interface{}) (err error)
	Returning(value interface{}, columns ...string) IDocumentVersionDo
	UnderlyingDB() *gorm.DB
	schema.Tabler
}

func (d documentVersionDo) Debug() IDocumentVersionDo {
	return d.withDO(d.DO.Debug())
}

func (d documentVersionDo) WithContext(ctx context.Context) IDocumentVersionDo {
	return d.withDO(d.DO.WithContext(ctx))
}

func (d documentVersionDo) ReadDB() IDocumentVersionDo {
	return d.Clauses(dbresolver.Read)
}

func (d documentVersionDo) WriteDB() IDocumentVersionDo {
	return d.Clauses(dbresolver.Write)
}

func (d documentVersionDo) Session(config *gorm.Session) IDocumentVersionDo {
	return d.withDO(d.DO.Session(config))
}

func (d documentVersionDo) Clauses(conds ...clause.Expression) IDocumentVersionDo {
	return d.withDO(d.DO.Clauses(conds...))
}

func (d documentVersionDo) Returning(value interface{}, columns ...string) IDocumentVersionDo {
	return d.withDO(d.DO.Returning(value, columns...))
}

func (d documentVersionDo) Not(conds ...gen.Condition) IDocumentVersionDo {
	return d.withDO(d.DO.Not(conds...))
}

func (d documentVersionDo) Or(conds ...gen.Condition) IDocumentVersionDo {
	return d.withDO(d.DO.Or(conds...))
}

func (d documentVersionDo) Select(conds ...field.Expr) IDocumentVersionDo {
	return d.withDO(d.DO.Select(conds...))
}

func (d documentVersionDo) Where(conds ...gen.Condition) IDocumentVersionDo {
	return d.withDO(d.DO.Where(conds...))
}

func (d documentVersionDo) Order(conds ...field.Expr) IDocumentVersionDo {
	return d.withDO(d.DO.Order(conds...))
}

func (d documentVersionDo) Distinct(cols ...field.Expr) IDocumentVersionDo {
	return d.withDO(d.DO.Distinct(cols...))
}

func (d documentVersionDo) Omit(cols ...field.Expr) IDocumentVersionDo {
	return d.withDO(d.DO.Omit(cols...))
}

func (d documentVersionDo) Join(table schema.Tabler, on ...field.Expr) IDocumentVersionDo {
	return d.withDO(d.DO.Join(table, on...))
}

func (d documentVersionDo) LeftJoin(table schema.Tabler, on ...field.Expr) IDocumentVersionDo {
	return d.withDO(d.DO.LeftJoin(table, on...))
}

func (d documentVersionDo) RightJoin(table schema.Tabler, on ...field.Expr) IDocumentVersionDo {
	return d.withDO(d.DO.RightJoin(table, on...))
}

func (d documentVersionDo) Group(cols ...field.Expr) IDocumentVersionDo {
	return d.withDO(d.DO.Group(cols...))
}

func (d documentVersionDo) Having(conds ...gen.Condition) IDocumentVersionDo {
	return d.withDO(d.DO.Having(conds...))
}

func (d documentVersionDo) Limit(limit int) IDocumentVersionDo {
	return d.withDO(d.DO.Limit(limit))
}

func (d documentVersionDo) Offset(offset int) IDocumentVersionDo {
	return d.withDO(d.DO.Offset(offset))
}

func (d documentVersionDo) Scopes(funcs ...func(gen.Dao) gen.Dao) IDocumentVersionDo {
	return d.withDO(d.DO.Scopes(funcs...))
}

func (d documentVersionDo) Unscoped() IDocumentVersionDo {
	return d.withDO(d.DO.Unscoped())
}

func (d documentVersionDo) Create(values ...*models.DocumentVersion) error {
	if len(values) == 0 {
		return nil
	}
	return d.DO.Create(values)
}

func (d documentVersionDo) CreateInBatches(values []*models.DocumentVersion, batchSize int) error {
	return d.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (d documentVersionDo) Save(values ...*models.DocumentVersion) error {
	if len(values) == 0 {
		return nil
	}
	return d.DO.Save(values)
}

func (d documentVersionDo) First() (*models.DocumentVersion, error) {
	if result, err := d.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*models.DocumentVersion), nil
	}
}

func (d documentVersionDo) Take() (*models.DocumentVersion, error) {
	if result, err := d.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*models.DocumentVersion), nil
	}
}

func (d documentVersionDo) Last() (*models.DocumentVersion, error) {
	if result, err := d.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*models.DocumentVersion), nil
	}
}

func (d documentVersionDo) Find() ([]*models.DocumentVersion, error) {
	result, err := d.DO.Find()
	return result.([]*models.DocumentVersion), err
}

func (d documentVersionDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*models.DocumentVersion, err error) {
	buf := make([]*models.DocumentVersion, 0, batchSize)
	err = d.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (d documentVersionDo) FindInBatches(result *[]*models.DocumentVersion, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return d.DO.FindInBatches(result, batchSize, fc)
}

func (d documentVersionDo) Attrs(attrs ...field.AssignExpr) IDocumentVersionDo {
	return d.withDO(d.DO.Attrs(attrs...))
}

func (d documentVersionDo) Assign(attrs ...field.AssignExpr) IDocumentVersionDo {
	return d.withDO(d.DO.Assign(attrs...))
}

func (d documentVersionDo) Joins(fields ...field.RelationField) IDocumentVersionDo {
	for _, _f := range fields {
		d = *d.withDO(d.DO.Joins(_f))
	}
	return &d
}

func (d documentVersionDo) Preload(fields ...field.RelationField) IDocumentVersionDo {
	for _, _f := range fields {
		d = *d.withDO(d.DO.Preload(_f))
	}
	return &d
}

func (d documentVersionDo) FirstOrInit() (*models.DocumentVersion, error) {
	if result, err := d.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*models.DocumentVersion), nil
	}
}

func (d documentVersionDo) FirstOrCreate() (*models.DocumentVersion, error) {
	if result, err := d.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*models.DocumentVersion), nil
	}
}

func (d documentVersionDo) FindByPage(offset int, limit int) (result []*models.DocumentVersion, count int64, err error) {
	result, err = d.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = d.Offset(-1).Limit(-1).Count()
	return
}

func (d documentVersionDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = d.Count()
	if err != nil {
		return
	}

	err = d.Offset(offset).Limit(limit).Scan(result)
	return
}

func (d documentVersionDo) Scan(result interface{}) (err error) {
	return d.DO.Scan(result)
}

func (d documentVersionDo) Delete(models ...*models.DocumentVersion) (result gen.ResultInfo, err error) {
	return d.DO.Delete(models)
}

func (d *documentVersionDo) withDO(do gen.Dao) *documentVersionDo {
	d.DO = *do.(*gen.DO)
	return d
}
