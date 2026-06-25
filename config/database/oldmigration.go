package database

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"domesv2/internal/model"
	"gorm.io/gorm"
)

// Legacy schema structures for decoding
type LegacyProp struct {
	Type string          `gorm:"column:type"`
	Data json.RawMessage `gorm:"column:data"`
}

type LegacyData struct {
	ID        int             `gorm:"column:id"`
	Data      string          `gorm:"column:data"`
	CreatedBy string          `gorm:"column:created_by"`
	CreatedAt time.Time       `gorm:"column:createdAt"`
	UpdatedAt time.Time       `gorm:"column:updatedAt"`
}

type PropOption struct {
	Key     string `json:"key"`
	Value   string `json:"value"`
	Label   string `json:"label"`
}

type FormConfig struct {
	Tag     string       `json:"tag"`
	Options []PropOption `json:"options"`
}

type LegacyJSON struct {
	ID               string   `json:"id"`
	Title            string   `json:"title"`
	DocumentLink     string   `json:"document_link"`
	Thumbnail        string   `json:"thumbnail"`
	Published        string   `json:"published"`
	Language         string   `json:"language"`
	GeographicScope  string   `json:"data"` // data_no or data_yes
	Unct             []string `json:"unct"`
	JointProgramme   string   `json:"joint_programme"`
	Sdg              []string `json:"sdg"`
	Sector           []string `json:"sector"`
	Lnob             []string `json:"lnob"`
	Thematic         []string `json:"thematic"`
	NonUnPartners    string   `json:"non_un_partners"`
	Email            string   `json:"email"`
	Phone            string   `json:"phone"`
	Position         string   `json:"position"`
	Nam              string   `json:"nam"`
}

var sectorCodes = []string{
	"agriculture-food",                         // 1
	"business-investment",                      // 2
	"conflict-violence-radicalism",             // 3
	"covid-19",                                 // 4
	"disability-vulnerability-social-welfare",  // 5
	"disaster-emergency",                       // 6
	"economic-development",                     // 7
	"education-culture",                        // 8
	"energy-natural-resources",                 // 9
	"environment-climate-change",               // 10
	"fishery-maritime",                         // 11
	"governance-corruption",                    // 12
	"gender-child-protection",                  // 13
	"health-nutrition",                         // 14
	"infrastructure-development",               // 15
	"innovation-technology",                    // 16
	"livelihood-employment",                    // 17
	"population-migration",                     // 18
	"poverty-social-exclusion",                 // 19
	"public-finance-tax-fiscal-policy",         // 20
	"rural-regional-development",               // 21
	"social-security-protection",               // 22
	"urban-development",                        // 23
	"water-sanitation",                         // 24
}

var thematicCodes = []string{
	"inclusive-economic-transformation",
	"environmental-development-climate-resilience",
	"human-development",
	"democratic-governance-security",
}

func generateSlug(title, id string) string {
	reg := regexp.MustCompile("[^a-zA-Z0-9]+")
	slug := reg.ReplaceAllString(strings.ToLower(title), "-")
	slug = strings.Trim(slug, "-")
	if len(slug) > 200 {
		slug = slug[:200]
	}
	suffix := ""
	if len(id) >= 8 {
		suffix = "-" + id[:8]
	}
	return slug + suffix
}

func parseYear(pubDate string) int {
	parts := strings.Split(pubDate, "/")
	if len(parts) == 3 {
		yearPart := parts[2]
		if len(yearPart) == 2 {
			y, err := strconv.Atoi(yearPart)
			if err == nil {
				if y > 50 {
					return 1900 + y
				}
				return 2000 + y
			}
		} else if len(yearPart) == 4 {
			y, err := strconv.Atoi(yearPart)
			if err == nil {
				return y
			}
		}
	}
	return time.Now().Year()
}

func MigrateLegacyData(db *gorm.DB) {
	// Clean existing V2 Documents to make migration repeatable
	log.Println("Cleaning existing V2 Document tables...")
	db.Exec("DELETE FROM v2_document_sdgs")
	db.Exec("DELETE FROM v2_document_sectors")
	db.Exec("DELETE FROM v2_document_lnobs")
	db.Exec("DELETE FROM V2Documents")

	// Fetch JP and LNOB lookups from Tableprops
	var prop LegacyProp
	if err := db.Table("Tableprops").Where("type = ?", "docs").First(&prop).Error; err != nil {
		log.Fatalf("Failed to fetch Tableprops lookups: %v", err)
	}

	var formConfigs []FormConfig
	if err := json.Unmarshal(prop.Data, &formConfigs); err != nil {
		log.Fatalf("Failed to unmarshal Tableprops data: %v", err)
	}

	jpLookup := make(map[string]string)
	lnobLookup := make(map[string]string)

	for _, fc := range formConfigs {
		if fc.Tag == "joint_programme" {
			for _, opt := range fc.Options {
				jpLookup[opt.Key] = opt.Value
			}
		} else if fc.Tag == "lnob" {
			for _, opt := range fc.Options {
				val := opt.Value
				if val == "" {
					val = opt.Label
				}
				lnobLookup[opt.Key] = val
			}
		}
	}

	// Load V2Master entities for referencing
	var agencies []model.Agency
	db.Find(&agencies)
	agencyMap := make(map[string]model.Agency)
	for _, a := range agencies {
		agencyMap[strings.ToLower(a.Code)] = a
	}

	var sdgs []model.Sdg
	db.Find(&sdgs)
	sdgMap := make(map[string]model.Sdg)
	for _, s := range sdgs {
		sdgMap[s.Code] = s
	}

	var sectors []model.Sector
	db.Find(&sectors)
	sectorMap := make(map[string]model.Sector)
	for _, sec := range sectors {
		sectorMap[sec.Code] = sec
	}

	var lnobs []model.Lnob
	db.Find(&lnobs)
	lnobMap := make(map[string]model.Lnob)
	for _, l := range lnobs {
		lnobMap[l.Code] = l
	}

	var jps []model.JointProgramme
	db.Find(&jps)
	jpMap := make(map[string]model.JointProgramme)
	for _, j := range jps {
		jpMap[strings.ToLower(j.Name)] = j
	}

	// Query legacy data
	var legacyRows []LegacyData
	if err := db.Table("Tabledatas").Find(&legacyRows).Error; err != nil {
		log.Fatalf("Failed to fetch Tabledatas: %v", err)
	}

	log.Printf("Fetched %d legacy records to migrate.", len(legacyRows))

	digitRegex := regexp.MustCompile(`\d+`)
	countSuccess := 0

	for _, row := range legacyRows {
		var lj LegacyJSON
		if err := json.Unmarshal([]byte(row.Data), &lj); err != nil {
			log.Printf("Warning: Failed to parse JSON for legacy row ID %d: %v", row.ID, err)
			continue
		}

		if lj.ID == "" || lj.Title == "" {
			continue
		}

		// Map basic info
		doc := model.Document{
			Code:              lj.ID,
			Slug:              generateSlug(lj.Title, lj.ID),
			Title:             lj.Title,
			Description:       lj.Title,
			Abstract:          "",
			Summary:           "",
			Year:              parseYear(lj.Published),
			DateOfPublication: lj.Published,
			TotalPages:        0,
			FileURL:           lj.DocumentLink,
			FileSize:          "0MB",
			CoverImage:        lj.Thumbnail,
			ExternalURL:       "",
			Views:             0,
			Downloads:         0,
			FocalPointName:    lj.Nam,
			FocalPointEmail:   lj.Email,
			FocalPointPhone:   lj.Phone,
			FocalPointDepartment: lj.Position,
			GeographicScope:   lj.GeographicScope,
			AuthorID:          133,
		}

		doc.ID = lj.ID
		doc.CreatedAt = &row.CreatedAt
		doc.UpdatedAt = &row.UpdatedAt
		sys := "System"
		doc.CreatedBy = &sys
		doc.UpdatedBy = &sys
		active := true
		doc.IsActive = &active

		// Language
		if lj.Language == "language_2" {
			doc.Language = "bahasa"
		} else {
			doc.Language = "english"
		}

		// Lead Agency / Other Agencies
		var otherAgencies []string
		for i, u := range lj.Unct {
			agCode := strings.ToLower(u)
			if agCode == "un women" {
				agCode = "un women"
			}
			if resolved, ok := agencyMap[agCode]; ok {
				if i == 0 {
					doc.LeadAgencyCode = resolved.Code
				} else {
					otherAgencies = append(otherAgencies, resolved.Code)
				}
			}
		}
		oaBytes, _ := json.Marshal(otherAgencies)
		doc.OtherAgencies = string(oaBytes)

		// Joint Programme
		if lj.JointProgramme != "" {
			jpName := lj.JointProgramme
			if strings.HasPrefix(lj.JointProgramme, "jp") {
				if resolvedName, ok := jpLookup[lj.JointProgramme]; ok {
					jpName = resolvedName
				}
			}
			if resolvedJp, ok := jpMap[strings.ToLower(jpName)]; ok {
				doc.JointProgrammeCode = resolvedJp.Code
			}
		}

		// Thematic Areas
		var thematicAreas []string
		for _, them := range lj.Thematic {
			themClean := strings.ReplaceAll(them, "!", "1")
			matches := digitRegex.FindAllString(themClean, -1)
			for _, m := range matches {
				idx, _ := strconv.Atoi(m)
				if idx >= 1 && idx <= 4 {
					thematicAreas = append(thematicAreas, thematicCodes[idx-1])
				}
			}
		}
		taBytes, _ := json.Marshal(thematicAreas)
		doc.ThematicAreas = string(taBytes)

		doc.Tags = "[]"
		doc.SupportingFiles = "[]"
		doc.NonUnPartners = "[]"

		// SDGs (Many-to-many)
		for _, s := range lj.Sdg {
			sClean := strings.ReplaceAll(s, "!", "1")
			matches := digitRegex.FindAllString(sClean, -1)
			for _, m := range matches {
				idx, _ := strconv.Atoi(m)
				if idx >= 1 && idx <= 17 {
					code := fmt.Sprintf("GOAL %d", idx)
					if sdgEntity, ok := sdgMap[code]; ok {
						doc.Sdgs = append(doc.Sdgs, sdgEntity)
					}
				}
			}
		}

		// Sectors (Many-to-many)
		for _, sec := range lj.Sector {
			matches := digitRegex.FindAllString(sec, -1)
			for _, m := range matches {
				idx, _ := strconv.Atoi(m)
				if idx >= 1 && idx <= 24 {
					code := sectorCodes[idx-1]
					if sectorEntity, ok := sectorMap[code]; ok {
						doc.Sectors = append(doc.Sectors, sectorEntity)
					}
				}
			}
		}

		// LNOBs (Many-to-many)
		for _, l := range lj.Lnob {
			lnobName := l
			if strings.HasPrefix(l, "lnob_") {
				if resolvedName, ok := lnobLookup[l]; ok {
					lnobName = resolvedName
				}
			}
			var code string
			switch strings.ToLower(lnobName) {
			case "women and girls":
				code = "women-girls"
			case "children", "youth":
				code = "youth-children"
			case "persons with disabilities":
				code = "disabilities"
			default:
				code = "others"
			}
			if lnobEntity, ok := lnobMap[code]; ok {
				doc.Lnobs = append(doc.Lnobs, lnobEntity)
			}
		}

		if err := db.Create(&doc).Error; err != nil {
			log.Printf("Error inserting document %s: %v", lj.ID, err)
		} else {
			countSuccess++
		}
	}

	log.Printf("Migration completed! Successfully migrated %d / %d records.", countSuccess, len(legacyRows))
}
