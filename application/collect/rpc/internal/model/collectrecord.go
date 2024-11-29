package model

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
)

type CollectRecord struct {
	ID            int64 `gorm:"primary_key"`
	BizID         string
	ObjID         int64
	UserID        int64
	CollectStatus int
	CreateTime    time.Time
	UpdateTime    time.Time
}

func (m *CollectRecord) TableName() string {
	return "collect_record"
}

type CollectRecordModel struct {
	db *gorm.DB
}

func NewCollectRecordModel(db *gorm.DB) *CollectRecordModel {
	return &CollectRecordModel{
		db: db,
	}
}

func (m *CollectRecordModel) Insert(ctx context.Context, data *CollectRecord) error {
	return m.db.WithContext(ctx).Create(data).Error
}

func (m *CollectRecordModel) FindOne(ctx context.Context, id int64) (*CollectRecord, error) {
	var result CollectRecord
	err := m.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
	return &result, err
}

func (m *CollectRecordModel) Update(ctx context.Context, data *CollectRecord) error {
	return m.db.WithContext(ctx).Save(data).Error
}

func (m *CollectRecordModel) UpdateFields(ctx context.Context, id int64, values map[string]interface{}) error {
	return m.db.WithContext(ctx).Model(&CollectRecord{}).Where("id = ?", id).Updates(values).Error
}

func (m *CollectRecordModel) FindByBizIDAndUserID(ctx context.Context, bizId string, userId int64, limit int) ([]*CollectRecord, error) {
	var result []*CollectRecord
	err := m.db.WithContext(ctx).
		Where("biz_id = ? AND user_id = ? AND collect_status = 1", bizId, userId).
		Order("update_time desc").
		Limit(limit).
		Find(&result).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return result, err
}

func (m *CollectRecordModel) FindByBizIDAndObjIDAndUserID(ctx context.Context, bizId string, objId, userId int64) (*CollectRecord, error) {
	var result CollectRecord
	err := m.db.WithContext(ctx).
		Where("biz_id = ? AND obj_Id = ? AND user_id = ?", bizId, objId, userId).
		First(&result).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}

	return &result, err
}

//func (m *CollectRecordModel) FindByUserId(ctx context.Context, userId int64, limit int) ([]*CollectRecord, error) {
//	var result []*CollectRecord
//	err := m.db.WithContext(ctx).
//		Where("user_id = ? AND Collect_status = ?", userId, 1).
//		Order("id desc").
//		Limit(limit).
//		Find(&result).Error
//
//	return result, err
//}

//func (m *CollectRecordModel) FindByCollectedUserIds(ctx context.Context, userId int64, collectedUserIds []int64) ([]*CollectRecord, error) {
//	var result []*CollectRecord
//	err := m.db.WithContext(ctx).
//		Where("user_id = ?", userId).
//		Where("collected_user_id in (?)", collectedUserIds).
//		Find(&result).Error
//
//	return result, err
//}
//
//func (m *CollectRecordModel) FindByCollectedUserId(ctx context.Context, userId int64, limit int) ([]*CollectRecord, error) {
//	var result []*CollectRecord
//	err := m.db.WithContext(ctx).
//		Where("collected_user_id = ? AND collect_status = ?", userId, 1).
//		Order("id desc").
//		Limit(limit).
//		Find(&result).Error
//	return result, err
//}
