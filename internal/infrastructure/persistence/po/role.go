package po

import "time"

type Role struct {
	ID        int `gorm:"primaryKey"`
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
