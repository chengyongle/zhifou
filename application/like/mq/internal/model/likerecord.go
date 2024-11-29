package model

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type LikeRecord struct {
	ID         int64 `gorm:"primary_key"`
	BizID      string
	ObjID      int64
	UserID     int64
	LikeStatus int
	CreateTime time.Time
	UpdateTime time.Time
}

func (m *LikeRecord) TableName() string {
	return "like_record"
}

type LikeRecordModel struct {
	db *gorm.DB
}

func NewLikeRecordModel(db *gorm.DB) *LikeRecordModel {
	return &LikeRecordModel{
		db: db,
	}
}

func (m *LikeRecordModel) Insert(ctx context.Context, data *LikeRecord) error {
	return m.db.WithContext(ctx).Create(data).Error
}

func (m *LikeRecordModel) FindOne(ctx context.Context, id int64) (*LikeRecord, error) {
	var result LikeRecord
	err := m.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
	return &result, err
}

func (m *LikeRecordModel) Update(ctx context.Context, data *LikeRecord) error {
	return m.db.WithContext(ctx).Save(data).Error
}

func (m *LikeRecordModel) UpdateFields(ctx context.Context, id int64, values map[string]interface{}) error {
	return m.db.WithContext(ctx).Model(&LikeRecord{}).Where("id = ?", id).Updates(values).Error
}

func (m *LikeRecordModel) FindByBizIDAndObjIDAndUserID(ctx context.Context, bizId string, objId, userId int64) (*LikeRecord, error) {
	var result LikeRecord
	err := m.db.WithContext(ctx).
		Where("biz_id = ? AND obj_Id = ? AND user_id = ?", bizId, objId, userId).
		First(&result).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}

	return &result, err
}

//func (m *LikeRecordModel) FindByUserId(ctx context.Context, userId int64, limit int) ([]*LikeRecord, error) {
//	var result []*LikeRecord
//	err := m.db.WithContext(ctx).
//		Where("user_id = ? AND like_status = ?", userId, 1).
//		Order("id desc").
//		Limit(limit).
//		Find(&result).Error
//
//	return result, err
//}

//func (m *LikeRecordModel) FindByFollowedUserIds(ctx context.Context, userId int64, followedUserIds []int64) ([]*LikeRecord, error) {
//	var result []*LikeRecord
//	err := m.db.WithContext(ctx).
//		Where("user_id = ?", userId).
//		Where("followed_user_id in (?)", followedUserIds).
//		Find(&result).Error
//
//	return result, err
//}
//
//func (m *LikeRecordModel) FindByFollowedUserId(ctx context.Context, userId int64, limit int) ([]*LikeRecord, error) {
//	var result []*LikeRecord
//	err := m.db.WithContext(ctx).
//		Where("followed_user_id = ? AND follow_status = ?", userId, 1).
//		Order("id desc").
//		Limit(limit).
//		Find(&result).Error
//	return result, err
//}
