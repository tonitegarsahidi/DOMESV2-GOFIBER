package model

type Agency struct {
	V2Base
	Code    string `json:"code" gorm:"uniqueIndex;size:100;column:code"`
	Name    string `json:"name" gorm:"size:255;not null;column:name"`
	LogoURL string `json:"logo_url" gorm:"size:255;column:logo_url"`
}

func (Agency) TableName() string {
	return "V2Agencies"
}

type Sdg struct {
	V2Base
	Code  string `json:"code" gorm:"uniqueIndex;size:100;column:code"`
	Name  string `json:"name" gorm:"size:255;not null;column:name"`
	Icon  string `json:"icon" gorm:"size:255;column:icon"`
	Color string `json:"color" gorm:"size:50;column:color"`
}

func (Sdg) TableName() string {
	return "V2Sdgs"
}

type Sector struct {
	V2Base
	Code string `json:"code" gorm:"uniqueIndex;size:100;column:code"`
	Name string `json:"name" gorm:"size:255;not null;column:name"`
}

func (Sector) TableName() string {
	return "V2Sectors"
}

type Language struct {
	V2Base
	Code string `json:"code" gorm:"uniqueIndex;size:100;column:code"`
	Name string `json:"name" gorm:"size:255;not null;column:name"`
}

func (Language) TableName() string {
	return "V2Languages"
}

type JointProgramme struct {
	V2Base
	Code string `json:"code" gorm:"uniqueIndex;size:100;column:code"`
	Name string `json:"name" gorm:"size:255;not null;column:name"`
}

func (JointProgramme) TableName() string {
	return "V2JointProgrammes"
}

type Lnob struct {
	V2Base
	Code string `json:"code" gorm:"uniqueIndex;size:100;column:code"`
	Name string `json:"name" gorm:"size:255;not null;column:name"`
}

func (Lnob) TableName() string {
	return "V2Lnobs"
}

type NonUnPartner struct {
	V2Base
	Code string `json:"code" gorm:"uniqueIndex;size:100;column:code"`
	Name string `json:"name" gorm:"size:255;not null;column:name"`
}

func (NonUnPartner) TableName() string {
	return "V2NonUnPartners"
}

type Organization struct {
	V2Base
	Code string `json:"code" gorm:"uniqueIndex;size:100;column:code"`
	Name string `json:"name" gorm:"size:255;not null;column:name"`
}

func (Organization) TableName() string {
	return "V2Organizations"
}

type ReferenceRequest struct {
	Code    string `json:"code"`
	Name    string `json:"name"`
	LogoURL string `json:"logo_url,omitempty"`
	Icon    string `json:"icon,omitempty"`
	Color   string `json:"color,omitempty"`
}
