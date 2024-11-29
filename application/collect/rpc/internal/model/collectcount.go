package model

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type CollectCount struct {
	ID         int64 `gorm:"primary_key"`
	BizID      string
	ObjID      int64
	CollectNum int
	CreateTime time.Time
	UpdateTime time.Time
}

func (m *CollectCount) TableName() string {
	return "collect_count"
}

type CollectCountModel struct {
	db *gorm.DB
}

func NewCollectCountModel(db *gorm.DB) *CollectCountModel {
	return &CollectCountModel{
		db: db,
	}
}

func (m *CollectCountModel) Insert(ctx context.Context, data *CollectCount) error {
	return m.db.Create(data).Error
}

func (m *CollectCountModel) FindOne(ctx context.Context, id int64) (*CollectCount, error) {
	var result CollectCount
	err := m.db.Where("id = ?", id).First(&result).Error
	return &result, err
}

func (m *CollectCountModel) Update(ctx context.Context, data *CollectCount) error {
	return m.db.Save(data).Error
}

func (m *CollectCountModel) IncrCollectCount(ctx context.Context, bizid string, objId int64) error {
	return m.db.WithContext(ctx).
		Exec("INSERT INTO collect_count (biz_id,obj_id, like_num) VALUES (?, 1) ON DUPLICATE KEY UPDATE collect_num = collect_num + 1", bizid, objId).
		Error
}

func (m *CollectCountModel) DecrCollectCount(ctx context.Context, bizid string, objId int64) error {
	return m.db.WithContext(ctx).
		Exec("UPDATE collect_count SET collect_num = collect_num - 1 WHERE biz_id = ? AND obj_id = ? AND collect_num > 0", bizid, objId).
		Error
}
