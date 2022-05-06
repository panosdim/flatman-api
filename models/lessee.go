package models

type Lessee struct {
	ID         uint    `gorm:"primary_key;autoIncrement" json:"id"`
	FlatID     uint    `json:"flat_id"`
	Name       string  `gorm:"size:255;not null;unique" json:"name"`
	Address    string  `gorm:"not null" json:"address"`
	PostalCode string  `gorm:"size:5;not null" json:"postal_code"`
	From       string  `gorm:"not null" json:"from"`
	Until      string  `gorm:"not null" json:"until"`
	Tin        string  `gorm:"size:9;not null" json:"tin"`
	Rent       float64 `gorm:"not null" json:"rent"`
}
