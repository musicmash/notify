package db

import "time"

type Notification struct {
	ID        int `gorm:"primary_key" sql:"AUTO_INCREMENT"`
	Date      time.Time
	UserName  string `gorm:"unique_index:idx_user_name_release_id"`
	ReleaseID uint64 `gorm:"unique_index:idx_user_name_release_id"`
}

type NotificationMgr interface {
	CreateNotification(notification *Notification) error
	GetNotificationsForUser(userName string) ([]*Notification, error)
}

func (mgr *AppDatabaseMgr) GetNotificationsForUser(userName string) ([]*Notification, error) {
	notifications := []*Notification{}
	if err := mgr.db.Where("user_name = ?", userName).Find(&notifications).Error; err != nil {
		return nil, err
	}
	return notifications, nil
}

func (mgr *AppDatabaseMgr) CreateNotification(notification *Notification) error {
	return mgr.db.Create(&notification).Error
}
