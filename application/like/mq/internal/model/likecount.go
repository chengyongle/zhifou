package model

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type LikeCount struct {
	ID         int64 `gorm:"primary_key"`
	BizID      string
	ObjID      int64
	LikeNum    int
	CreateTime time.Time
	UpdateTime time.Time
}

func (m *LikeCount) TableName() string {
	return "like_count"
}

type LikeCountModel struct {
	db *gorm.DB
}

func NewLikeCountModel(db *gorm.DB) *LikeCountModel {
	return &LikeCountModel{
		db: db,
	}
}

func (m *LikeCountModel) Insert(ctx context.Context, data *LikeCount) error {
	return m.db.Create(data).Error
}

func (m *LikeCountModel) FindOne(ctx context.Context, id int64) (*LikeCount, error) {
	var result LikeCount
	err := m.db.Where("id = ?", id).First(&result).Error
	return &result, err
}

func (m *LikeCountModel) Update(ctx context.Context, data *LikeCount) error {
	return m.db.Save(data).Error
}

func (m *LikeCountModel) IncrLikeCount(ctx context.Context, bizid string, objId int64) error {
	return m.db.WithContext(ctx).
		Exec("INSERT INTO like_count (biz_id,obj_id, like_num) VALUES (?,?,1) ON DUPLICATE KEY UPDATE like_num = like_num + 1", bizid, objId).
		Error
}

func (m *LikeCountModel) DecrLikeCount(ctx context.Context, bizid string, objId int64) error {
	return m.db.WithContext(ctx).
		Exec("UPDATE like_count SET like_num = like_num - 1 WHERE biz_id = ? AND obj_id = ? AND like_num > 0", bizid, objId).
		Error
}
