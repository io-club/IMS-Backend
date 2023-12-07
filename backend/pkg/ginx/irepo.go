package ioginx

import (
	"context"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	iodb "ims-server/pkg/db"
)

type IRepo[T schema.Tabler] struct {
	applier func(tx *gorm.DB) *gorm.DB
}

func (r IRepo[T]) DB() *gorm.DB {
	db := iodb.NewDB()
	if r.applier != nil {
		db = r.applier(db)
	}
	return db
}

func (r IRepo[T]) Create(ctx context.Context, t *T) error {
	return r.CreateWithDB(ctx, r.DB(), t)
}

func (r IRepo[T]) CreateWithDB(ctx context.Context, tx *gorm.DB, t *T) error {
	err := tx.WithContext(ctx).Create(t).Error
	if err != nil {
		return err
	}
	return nil
}

func (r IRepo[T]) MCreate(ctx context.Context, ts []*T) error {
	return r.MCreateWithDB(ctx, r.DB(), ts)
}

func (r IRepo[T]) MCreateWithDB(ctx context.Context, tx *gorm.DB, ts []*T) error {
	// 避免因 ts 为空导致的 sql 报错
	if len(ts) == 0 {
		return nil
	}
	err := tx.WithContext(ctx).Create(ts).Error
	if err != nil {
		return err
	}
	return nil
}

func (r IRepo[T]) Last() *T {
	var t T
	r.DB().Last(&t)
	return &t
}

func (r IRepo[T]) Get(ctx context.Context, id uint, selectField ...string) (*T, error) {
	return r.GetWithDBRaw(ctx, r.DB(), id, selectField...)
}

func (r IRepo[T]) GetWithDB(ctx context.Context, tx *gorm.DB, id uint, selectField ...string) (*T, error) {
	return r.GetWithDBRaw(ctx, tx, id, selectField...)
}

func (r IRepo[T]) GetWithDBRaw(ctx context.Context, tx *gorm.DB, id uint, selectField ...string) (*T, error) {
	var res T
	conn := tx.WithContext(ctx).Model(&res).Select(selectField)
	err := conn.First(&res, id).Error
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// MGet 批量获取
func (r IRepo[T]) MGet(ctx context.Context, ids []uint, selectField ...string) ([]T, error) {
	return r.MGetWithDB(ctx, r.DB(), ids, selectField...)
}

func (r IRepo[T]) MGetWithDB(ctx context.Context, tx *gorm.DB, ids []uint, selectField ...string) ([]T, error) {
	res := []T{}
	err := tx.WithContext(ctx).Select(selectField).Where(" id IN ? ", ids).Find(&res).Error
	if err != nil {
		return nil, err
	}
	return res, nil
}

// List 简单分页，从第 0 页开始
func (r IRepo[T]) List(ctx context.Context, page uint, size uint, selectField ...string) ([]T, error) {
	res := []T{}
	if size > 50 {
		size = 50
	}
	err := r.DB().WithContext(ctx).Select(selectField).Offset(int(page) * int(size)).Limit(int(size)).Find(&res).Error
	return res, err
}

func (r IRepo[T]) Count(ctx context.Context, req *iodb.PageBuilder) (int64, error) {
	var t T
	var count int64
	err := req.ToFilterDB(r.DB().WithContext(ctx).Model(&t)).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r IRepo[T]) Pageable(ctx context.Context, req *iodb.PageBuilder, selectField ...string) ([]T, error) {
	res := []T{}
	err := req.ToPageDB(r.DB().WithContext(ctx).Select(selectField)).Find(&res).Error
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (r IRepo[T]) ListAll(ctx context.Context, selectField ...string) ([]T, error) {
	res := []T{}
	err := r.DB().WithContext(ctx).Select(selectField).Find(&res).Error
	return res, err
}

func (r IRepo[T]) Update(ctx context.Context, id uint, fields map[string]interface{}) (*T, error) {
	return r.UpdateWithDB(ctx, r.DB(), id, fields)
}

func (r IRepo[T]) UpdateWithDB(ctx context.Context, tx *gorm.DB, id uint, fields map[string]interface{}) (*T, error) {
	var t T
	err := tx.WithContext(ctx).Model(&t).Where("id = ?", id).Updates(fields).Error
	if err != nil {
		return nil, err
	}
	res, err := r.GetWithDBRaw(ctx, tx, id)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (r IRepo[T]) MUpdate(ctx context.Context, ids []uint, fields map[string]interface{}) ([]T, error) {
	return r.MUpdateWithDB(ctx, r.DB(), ids, fields)
}

func (r IRepo[T]) MUpdateWithDB(ctx context.Context, tx *gorm.DB, ids []uint, fields map[string]interface{}) ([]T, error) {
	var t []T
	err := tx.WithContext(ctx).Model(&t).Where("id IN ?", ids).Updates(fields).Error
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (r IRepo[T]) Delete(ctx context.Context, id uint) error {
	return r.DeleteWithDB(ctx, r.DB(), id)
}

func (r IRepo[T]) DeleteWithDB(ctx context.Context, tx *gorm.DB, id uint) error {
	var t T
	err := tx.WithContext(ctx).Delete(&t, id).Error
	if err != nil {
		return err
	}
	return nil
}

func (r IRepo[T]) MDelete(ctx context.Context, ids []uint) error {
	return r.MDeleteWithDB(ctx, r.DB(), ids)
}

func (r IRepo[T]) MDeleteWithDB(ctx context.Context, tx *gorm.DB, ids []uint) error {
	var t T
	err := tx.WithContext(ctx).Where("id IN ?", ids).Delete(&t).Error
	if err != nil {
		return err
	}
	return nil
}
