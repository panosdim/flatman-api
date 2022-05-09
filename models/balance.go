package models

type Balance struct {
	ID      uint    `gorm:"primary_key;autoIncrement" json:"id"`
	FlatID  uint    `gorm:"not null" json:"flat_id"`
	Date    string  `gorm:"not null" json:"date"`
	Amount  float64 `gorm:"not null" json:"amount"`
	Comment string  `gorm:"not null" json:"comment"`
}
