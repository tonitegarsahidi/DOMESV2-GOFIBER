package model

import "time"

type Agency struct {
	Code      string    `json:"code" gorm:"primaryKey;size:100;column:code"`
	Name      string    `json:"name" gorm:"size:255;not null;column:name"`
	LogoURL   string    `json:"logo_url" gorm:"size:255;column:logo_url"`
	CreatedAt time.Time `json:"-" gorm:"column:createdAt"`
	UpdatedAt time.Time `json:"-" gorm:"column:updatedAt"`
}

func (Agency) TableName() string {
	return "Agencies"
}

type Sdg struct {
	Code      string    `json:"code" gorm:"primaryKey;size:100;column:code"`
	Name      string    `json:"name" gorm:"size:255;not null;column:name"`
	Icon      string    `json:"icon" gorm:"size:255;column:icon"`
	Color     string    `json:"color" gorm:"size:50;column:color"`
	CreatedAt time.Time `json:"-" gorm:"column:createdAt"`
	UpdatedAt time.Time `json:"-" gorm:"column:updatedAt"`
}

func (Sdg) TableName() string {
	return "Sdgs"
}

type Sector struct {
	Code      string    `json:"code" gorm:"primaryKey;size:100;column:code"`
	Name      string    `json:"name" gorm:"size:255;not null;column:name"`
	CreatedAt time.Time `json:"-" gorm:"column:createdAt"`
	UpdatedAt time.Time `json:"-" gorm:"column:updatedAt"`
}

func (Sector) TableName() string {
	return "Sectors"
}

type Language struct {
	Code      string    `json:"code" gorm:"primaryKey;size:100;column:code"`
	Name      string    `json:"name" gorm:"size:255;not null;column:name"`
	CreatedAt time.Time `json:"-" gorm:"column:createdAt"`
	UpdatedAt time.Time `json:"-" gorm:"column:updatedAt"`
}

func (Language) TableName() string {
	return "Languages"
}

type JointProgramme struct {
	Code      string    `json:"code" gorm:"primaryKey;size:100;column:code"`
	Name      string    `json:"name" gorm:"size:255;not null;column:name"`
	CreatedAt time.Time `json:"-" gorm:"column:createdAt"`
	UpdatedAt time.Time `json:"-" gorm:"column:updatedAt"`
}

func (JointProgramme) TableName() string {
	return "JointProgrammes"
}

type Lnob struct {
	Code      string    `json:"code" gorm:"primaryKey;size:100;column:code"`
	Name      string    `json:"name" gorm:"size:255;not null;column:name"`
	CreatedAt time.Time `json:"-" gorm:"column:createdAt"`
	UpdatedAt time.Time `json:"-" gorm:"column:updatedAt"`
}

func (Lnob) TableName() string {
	return "Lnobs"
}

type NonUnPartner struct {
	Code      string    `json:"code" gorm:"primaryKey;size:100;column:code"`
	Name      string    `json:"name" gorm:"size:255;not null;column:name"`
	CreatedAt time.Time `json:"-" gorm:"column:createdAt"`
	UpdatedAt time.Time `json:"-" gorm:"column:updatedAt"`
}

func (NonUnPartner) TableName() string {
	return "NonUnPartners"
}

type Organization struct {
	Code      string    `json:"code" gorm:"primaryKey;size:100;column:code"`
	Name      string    `json:"name" gorm:"size:255;not null;column:name"`
	CreatedAt time.Time `json:"-" gorm:"column:createdAt"`
	UpdatedAt time.Time `json:"-" gorm:"column:updatedAt"`
}

func (Organization) TableName() string {
	return "Organizations"
}
