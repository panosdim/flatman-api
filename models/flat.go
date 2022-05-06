package models

type Flat struct {
	ID      uint   `gorm:"primary_key;autoIncrement" json:"id"`
	UserID  uint   `json:"-"`
	Name    string `gorm:"size:255;not null;unique" json:"name"`
	Address string `gorm:"not null" json:"address"`
	Floor   uint   `json:"floor"`
}
