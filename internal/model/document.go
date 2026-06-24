package model

import "time"

type Document struct {
	ID                    uint             `json:"id" gorm:"primaryKey;column:id"`
	Code                  string           `json:"code" gorm:"size:100;column:code"`
	Slug                  string           `json:"slug" gorm:"uniqueIndex;size:255;column:slug"`
	Title                 string           `json:"title" gorm:"size:255;not null;column:title"`
	Description           string           `json:"description" gorm:"type:text;column:description"`
	Abstract              string           `json:"abstract" gorm:"type:text;column:abstract"`
	Summary               string           `json:"summary" gorm:"type:text;column:summary"`
	Year                  int              `json:"year" gorm:"column:year"`
	DateOfPublication     string           `json:"date_of_publication" gorm:"size:50;column:date_of_publication"`
	TotalPages            int              `json:"total_pages" gorm:"column:total_pages"`
	Language              string           `json:"language" gorm:"size:255;column:language"`
	FileURL               string           `json:"file_url" gorm:"size:255;column:file_url"`
	FileSize              string           `json:"file_size" gorm:"size:50;column:file_size"`
	CoverImage            string           `json:"cover_image" gorm:"size:255;column:cover_image"`
	ExternalURL           string           `json:"external_url" gorm:"size:255;column:external_url"`
	Status                string           `json:"status" gorm:"size:50;default:'pending_review';column:status"`
	Views                 int              `json:"views" gorm:"default:0;column:views"`
	Downloads             int              `json:"downloads" gorm:"default:0;column:downloads"`
	FocalPointName        string           `json:"focal_point_name" gorm:"size:255;column:focal_point_name"`
	FocalPointEmail       string           `json:"focal_point_email" gorm:"size:255;column:focal_point_email"`
	FocalPointPhone       string           `json:"focal_point_phone" gorm:"size:100;column:focal_point_phone"`
	FocalPointDepartment  string           `json:"focal_point_department" gorm:"size:255;column:focal_point_department"`
	LeadAgencyCode        string           `json:"lead_agency_code" gorm:"size:100;column:lead_agency_code"`
	LeadAgency            *Agency          `json:"lead_agency,omitempty" gorm:"foreignKey:LeadAgencyCode;constraint:false"`
	JointProgrammeCode    string           `json:"joint_programme_code" gorm:"size:100;column:joint_programme_code"`
	JointProgramme        *JointProgramme  `json:"joint_programme,omitempty" gorm:"foreignKey:JointProgrammeCode;constraint:false"`
	GeographicScope       string           `json:"geographic_scope" gorm:"size:255;column:geographic_scope"`
	ThematicAreas         string           `json:"thematic_areas" gorm:"type:text;column:thematic_areas"` // JSON list
	Tags                  string           `json:"tags" gorm:"type:text;column:tags"`                     // JSON list
	OtherAgencies         string           `json:"other_agencies" gorm:"type:text;column:other_agencies"`   // JSON list
	NonUnPartners         string           `json:"non_un_partners" gorm:"type:text;column:non_un_partners"` // JSON list
	SupportingFiles       string           `json:"supporting_files" gorm:"type:text;column:supporting_files"` // JSON list
	AuthorID              uint             `json:"author_id" gorm:"column:author_id"`
	Author                *User            `json:"author,omitempty" gorm:"foreignKey:AuthorID;constraint:false"`
	Sdgs                  []Sdg            `json:"sdgs" gorm:"many2many:DocumentSdgs;constraint:false;"`
	Sectors               []Sector         `json:"sectors" gorm:"many2many:DocumentSectors;constraint:false;"`
	Lnobs                 []Lnob           `json:"lnob_groups" gorm:"many2many:DocumentLnobs;constraint:false;"`
	CreatedAt             time.Time        `json:"created_at" gorm:"column:createdAt"`
	UpdatedAt             time.Time        `json:"updated_at" gorm:"column:updatedAt"`
}

func (Document) TableName() string {
	return "Documents"
}

// Request payload structures
type SubmissionRequest struct {
	Title             string           `json:"title"`
	ShortDescription  string           `json:"short_description"`
	Abstract          string           `json:"abstract"`
	DetailedSummary   string           `json:"detailed_summary"`
	DateOfPublication string           `json:"date_of_publication"`
	TotalPages        int              `json:"total_pages"`
	Language          string           `json:"language"`
	PublicationStatus string           `json:"publication_status"`
	Tags              []string         `json:"tags"`
	FileURL           string           `json:"file_url"`
	FileSize          string           `json:"file_size"`
	CoverImageURL     string           `json:"cover_image_url"`
	ExternalURL       string           `json:"external_url"`
	SupportingFiles   []SupportingFile `json:"supporting_files"`
	Agency            string           `json:"agency"`
	FocalPoint        FocalPointDTO    `json:"focal_point"`
	Sdgs              []string         `json:"sdgs"`
	Sectors           []string         `json:"sectors"`
	LnobGroups        []string         `json:"lnob_groups"`
	JointProgramme    string           `json:"joint_programme"`
	OtherAgencies     []string         `json:"other_agencies"`
	NonUnPartners     []PartnerDTO     `json:"non_un_partners"`
	ThematicAreas     []string         `json:"thematic_areas"`
	GeographicScope   string           `json:"geographic_scope"`
}

type DraftRequest struct {
	Step int         `json:"step"`
	Data interface{} `json:"data"`
}

type SupportingFile struct {
	URL         string `json:"url"`
	Type        string `json:"type"`
	Description string `json:"description"`
}

type FocalPointDTO struct {
	Name       string `json:"name"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
	Department string `json:"department"`
}

type PartnerDTO struct {
	Type string `json:"type"`
	Name string `json:"name"`
}

// Detailed response structs
type DocumentResponse struct {
	ID              uint             `json:"id"`
	Code            string           `json:"code"`
	Slug            string           `json:"slug"`
	Title           string           `json:"title"`
	Agency          string           `json:"agency"`
	Year            int              `json:"year"`
	Language        string           `json:"language"`
	FileURL         string           `json:"file_url"`
	FileSize        string           `json:"file_size"`
	DateAdded       string           `json:"date_added"`
	Type            string           `json:"type"`
	TotalPages      int              `json:"total_pages"`
	PubStatus       string           `json:"pub_status"`
	CoverImage      string           `json:"cover_image"`
	Abstract        string           `json:"abstract"`
	Summary         string           `json:"summary"`
	Sdgs            []SdgDTO         `json:"sdgs"`
	Tags            []string         `json:"tags"`
	ThematicAreas   []string         `json:"thematic_areas"`
	Sectors         []string         `json:"sectors"`
	LnobGroups      []string         `json:"lnob_groups"`
	Classification  Classification   `json:"classification"`
	FocalPoint      FocalPointDTO    `json:"focal_point"`
	Views           int              `json:"views"`
	Downloads       int              `json:"downloads"`
	CreatedAt       time.Time        `json:"created_at"`
	UpdatedAt       time.Time        `json:"updated_at"`
	SupportingFiles []SupportingFile `json:"supporting_files,omitempty"`
}

type SdgDTO struct {
	Code string `json:"code"`
	Name string `json:"name"`
	Icon string `json:"icon"`
}

type Classification struct {
	LeadAgency      string       `json:"lead_agency"`
	OtherAgencies   []string     `json:"other_agencies"`
	JointProgramme  string       `json:"joint_programme"`
	GeographicScope string       `json:"geographic_scope"`
	NonUnPartners   []PartnerDTO `json:"non_un_partners"`
}

type DocumentListItem struct {
	ID         uint     `json:"id"`
	Title      string   `json:"title"`
	Slug       string   `json:"slug"`
	Description string  `json:"description"`
	Agency     string   `json:"agency"`
	Year       int      `json:"year"`
	Language   string   `json:"language"`
	FileSize   string   `json:"file_size"`
	TotalPages int      `json:"total_pages"`
	Type       string   `json:"type"`
	PubStatus  string   `json:"pub_status"`
	CoverImage string   `json:"cover_image"`
	Sdgs       []string `json:"sdgs"`
	Tags       []string `json:"tags"`
	Views      int      `json:"views"`
	Downloads  int      `json:"downloads"`
	CreatedAt  time.Time `json:"created_at"`
}

type DocumentListResponse struct {
	Items      []DocumentListItem `json:"items"`
	Pagination Pagination         `json:"pagination"`
}

type Pagination struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	TotalItems int `json:"totalItems"`
	TotalPages int `json:"totalPages"`
}
