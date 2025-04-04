package model

import (
	"context"
	"gorm.io/gorm"
	"time"
)

type Meeting struct {
	ID         uint64         `json:"id"`         // 会议ID
	Name       string         `json:"name"`       // 会议名称
	Location   string         `json:"location"`   // 会议地点
	Time       time.Time      `json:"time"`       // 会议时间
	KeyWords   string         `json:"key_words"`  // 关键词
	Highlights string         `json:"highlights"` // 会议要点
	Content    string         `json:"content"`    // 会议内容
	Minutes    string         `json:"minutes"`    // 会议纪要
	CreatedAt  time.Time      `json:"created_at"` // 创建时间
	UpdatedAt  time.Time      `json:"updated_at"` // 更新时间
	DeletedAt  gorm.DeletedAt `json:"deleted_at"` // 删除时间
}

func (m Meeting) TableName() string {
	return "meeting"
}

type MeetingQuery struct {
	ctx context.Context
	db  *gorm.DB
}

func NewMeetingQuery(ctx context.Context, db *gorm.DB) MeetingQuery {
	return MeetingQuery{ctx: ctx, db: db}
}

func (m *MeetingQuery) Create(meeting Meeting) (Meeting, error) {
	err := m.db.WithContext(m.ctx).Create(&meeting).Error
	return meeting, err
}

func (m *MeetingQuery) Update(id uint64, meeting Meeting) (Meeting, error) {
	err := m.db.WithContext(m.ctx).Model(&Meeting{}).Where(&Meeting{ID: id}).Updates(&meeting).Error
	return meeting, err
}

func (m *MeetingQuery) Delete(meetingId string) error {
	return m.db.WithContext(m.ctx).Delete(&Meeting{}, meetingId).Error
}

func (m *MeetingQuery) GetById(meetingId uint64) (meeting Meeting, err error) {
	err = m.db.WithContext(m.ctx).Model(&Meeting{}).Where(&Meeting{ID: meetingId}).First(&meeting).Error
	return
}

func (m *MeetingQuery) GetAll(num int32, size int64) (meetings []Meeting, total int64, err error) {
	// 首先计算总数
	err = m.db.WithContext(m.ctx).Model(&Meeting{}).
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	// 然后查询数据
	query := m.db.WithContext(m.ctx).Model(&Meeting{}).Order("created_at desc")
	if num != 0 || size != 0 {
		query = query.Limit(int(size)).Offset(int(num))
	}
	err = query.Find(&meetings).Error
	return meetings, total, err
}
