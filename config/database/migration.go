package database

import (
	"os"
	"time"

	"domesv2/internal/model"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func MigrateAndSeed(db *gorm.DB) {
	if db == nil {
		zap.L().Error("Database client is nil, skipping migration")
		return
	}

	zap.L().Info("Running database migration check...")

	// 1. Update Users Table columns manually if they don't exist (Only if explicitly requested)
	if os.Getenv("RUN_USER_MIGRATION") == "true" {
		zap.L().Info("Running explicit Users table migration...")
		if !db.Migrator().HasColumn(&model.User{}, "role") {
			zap.L().Info("Adding 'role' column to Users table")
			db.Migrator().AddColumn(&model.User{}, "role")
		}
		if !db.Migrator().HasColumn(&model.User{}, "status") {
			zap.L().Info("Adding 'status' column to Users table")
			db.Migrator().AddColumn(&model.User{}, "status")
		}
		if !db.Migrator().HasColumn(&model.User{}, "avatar_url") {
			zap.L().Info("Adding 'avatar_url' column to Users table")
			db.Migrator().AddColumn(&model.User{}, "avatar_url")
		}
	} else {
		zap.L().Info("Skipping Users table migration (RUN_USER_MIGRATION is not true)")
	}

	// 2. Create NotificationPreferences Table
	if err := db.AutoMigrate(&model.NotificationPreference{}); err != nil {
		zap.L().Error("Failed to migrate NotificationPreferences", zap.Error(err))
	}

	// 3. Create AdminEmails Table
	if err := db.AutoMigrate(&model.AdminEmail{}); err != nil {
		zap.L().Error("Failed to migrate AdminEmails", zap.Error(err))
	}
	seedAdminEmails(db)

	// 4. Create Master Tables & seed them
	migrateAndSeedMasters(db)

	// 5. Create Documents Table
	if err := db.AutoMigrate(&model.Document{}); err != nil {
		zap.L().Error("Failed to migrate Documents", zap.Error(err))
	}

	// 6. Create Reports Table
	if err := db.AutoMigrate(&model.Report{}); err != nil {
		zap.L().Error("Failed to migrate Reports", zap.Error(err))
	}

	// 7. Seed sample documents if empty
	seedDocuments(db)

	zap.L().Info("Database migration check completed.")
}

func seedDocuments(db *gorm.DB) {
	var count int64
	db.Model(&model.Document{}).Where("status = ?", "published").Count(&count)
	if count > 0 {
		return
	}

	zap.L().Info("Seeding default documents...")

	var user model.User
	if err := db.First(&user).Error; err != nil {
		zap.L().Error("Failed to find user for document seeding", zap.Error(err))
		return
	}

	docs := []model.Document{
		{
			Code:              "UNDP-2023-SDG",
			Slug:              "annual-progress-report-sustainable-development-goals-indonesia-2023",
			Title:             "Annual Progress Report: Sustainable Development Goals in Indonesia 2023",
			Description:       "This comprehensive report outlines the milestones achieved in the past fiscal year, focusing on the impact of global initiatives on sustainable development goals across multiple governance levels...",
			Abstract:          "This comprehensive report outlines the milestones achieved in the past fiscal year, focusing on the impact of global initiatives on sustainable development goals across multiple governance levels.",
			Summary:           "Executive Overview: This report reviews the implementation and progress of the Sustainable Development Goals (SDGs) in Indonesia in 2023.",
			Year:              2023,
			DateOfPublication: "2023-10-15",
			TotalPages:        134,
			Language:          "English, Bahasa Indonesia",
			FileURL:           "/uploads/documents/sample-sdg.pdf",
			FileSize:          "5.2 MB",
			CoverImage:        "/images/doc-cover-sdg.png",
			Status:            "published",
			AuthorID:          user.ID,
			LeadAgencyCode:    "UNDP",
			Tags:              `["SDG", "Indonesia", "Annual Report"]`,
			Sdgs:              []model.Sdg{{Code: "GOAL 1"}, {Code: "GOAL 5"}, {Code: "GOAL 8"}, {Code: "GOAL 10"}},
			Sectors:           []model.Sector{{Code: "economic-development"}, {Code: "poverty-social-exclusion"}},
		},
		{
			Code:              "UNEP-2022-CLIMATE",
			Slug:              "climate-action-framework-maritime-biodiversity-archipelago",
			Title:             "Climate Action Framework: Maritime Biodiversity in the Archipelago",
			Description:       "A strategic roadmap for protecting marine and coastal habitats in the Indonesian archipelago through sustainable fishing practices and community-led conservation efforts...",
			Abstract:          "A strategic roadmap for protecting marine and coastal habitats in the Indonesian archipelago through sustainable fishing practices and community-led conservation efforts.",
			Summary:           "Maritime biodiversity is crucial for food security and climate resilience. This framework outlines policies for marine protection.",
			Year:              2022,
			DateOfPublication: "2022-08-20",
			TotalPages:        78,
			Language:          "English",
			FileURL:           "/uploads/documents/sample-climate.pdf",
			FileSize:          "3.8 MB",
			CoverImage:        "/images/doc-cover-ocean.png",
			Status:            "published",
			AuthorID:          user.ID,
			LeadAgencyCode:    "UNEP",
			Tags:              `["Climate Action", "Maritime", "Biodiversity"]`,
			Sdgs:              []model.Sdg{{Code: "GOAL 13"}, {Code: "GOAL 14"}},
			Sectors:           []model.Sector{{Code: "environment-climate-change"}, {Code: "fishery-maritime"}},
		},
		{
			Code:              "UNICEF-2023-CHILDREN",
			Slug:              "children-in-focus-socio-economic-protection-systems",
			Title:             "Children in Focus: Socio-Economic Protection Systems",
			Description:       "Analyzing the efficacy of social safety nets for vulnerable families across the outer islands of Indonesia, with policy recommendations for enhanced coverage and efficiency...",
			Abstract:          "Analyzing the efficacy of social safety nets for vulnerable families across the outer islands of Indonesia, with policy recommendations for enhanced coverage and efficiency.",
			Summary:           "This policy brief analyzes child poverty and social protection systems in remote regions of Indonesia.",
			Year:              2023,
			DateOfPublication: "2023-07-10",
			TotalPages:        92,
			Language:          "English, Bahasa Indonesia",
			FileURL:           "/uploads/documents/sample-children.pdf",
			FileSize:          "4.5 MB",
			CoverImage:        "/images/doc-cover-children.png",
			Status:            "published",
			AuthorID:          user.ID,
			LeadAgencyCode:    "UNICEF",
			Tags:              `["Children", "Social Protection", "Policy"]`,
			Sdgs:              []model.Sdg{{Code: "GOAL 4"}, {Code: "GOAL 10"}},
			Sectors:           []model.Sector{{Code: "education-culture"}, {Code: "social-security-protection"}},
		},
		{
			Code:              "UNWOMEN-2023-GENDER",
			Slug:              "gender-equality-workplace-progress-challenges",
			Title:             "Gender Equality in the Workplace: Progress and Challenges",
			Description:       "A national survey on female labor force participation, wage gaps, and policy interventions needed to foster inclusive economic growth and empower women in rural and urban areas...",
			Abstract:          "A national survey on female labor force participation, wage gaps, and policy interventions needed to foster inclusive economic growth and empower women in rural and urban areas.",
			Summary:           "This report presents empirical findings on gender gaps in employment, wages, and leadership roles in Indonesia.",
			Year:              2023,
			DateOfPublication: "2023-05-12",
			TotalPages:        45,
			Language:          "English",
			FileURL:           "/uploads/documents/sample-gender.pdf",
			FileSize:          "2.1 MB",
			CoverImage:        "/images/doc-cover-sdg.png",
			Status:            "published",
			AuthorID:          user.ID,
			LeadAgencyCode:    "UN Women",
			Tags:              `["Gender Equality", "Workplace", "Women Empowerment"]`,
			Sdgs:              []model.Sdg{{Code: "GOAL 5"}, {Code: "GOAL 8"}},
			Sectors:           []model.Sector{{Code: "gender-child-protection"}, {Code: "livelihood-employment"}},
		},
		{
			Code:              "FAO-2023-FOOD",
			Slug:              "sustainable-agriculture-food-security-resilience",
			Title:             "Sustainable Agriculture and Food Security Resilience",
			Description:       "Examining the impact of climate smart agriculture practices on crop yields and food security for smallholder farmers amidst changing weather patterns in Southeast Asia...",
			Abstract:          "Examining the impact of climate smart agriculture practices on crop yields and food security for smallholder farmers amidst changing weather patterns in Southeast Asia.",
			Summary:           "Food security is a critical pillar of sustainable development. This study presents best practices for climate-smart farming.",
			Year:              2023,
			DateOfPublication: "2023-06-18",
			TotalPages:        64,
			Language:          "English, Bahasa Indonesia",
			FileURL:           "/uploads/documents/sample-food.pdf",
			FileSize:          "3.1 MB",
			CoverImage:        "/images/doc-cover-ocean.png",
			Status:            "published",
			AuthorID:          user.ID,
			LeadAgencyCode:    "FAO",
			Tags:              `["Agriculture", "Food Security", "Resilience"]`,
			Sdgs:              []model.Sdg{{Code: "GOAL 2"}},
			Sectors:           []model.Sector{{Code: "agriculture-food"}, {Code: "environment-climate-change"}},
		},
	}

	for _, doc := range docs {
		var matchedSdgs []model.Sdg
		for _, s := range doc.Sdgs {
			var dbSdg model.Sdg
			if err := db.Where("code = ?", s.Code).First(&dbSdg).Error; err == nil {
				matchedSdgs = append(matchedSdgs, dbSdg)
			}
		}
		doc.Sdgs = matchedSdgs

		var matchedSectors []model.Sector
		for _, sec := range doc.Sectors {
			var dbSec model.Sector
			if err := db.Where("code = ?", sec.Code).First(&dbSec).Error; err == nil {
				matchedSectors = append(matchedSectors, dbSec)
			}
		}
		doc.Sectors = matchedSectors

		if err := db.Create(&doc).Error; err != nil {
			zap.L().Error("Failed to seed document", zap.String("code", doc.Code), zap.Error(err))
		}
	}
}


func seedAdminEmails(db *gorm.DB) {
	emails := []model.AdminEmail{
		{Email: "admin@un.org", AddedAt: time.Now()},
		{Email: "superadmin@un.org", AddedAt: time.Now()},
		{Email: "newadmin@un.org", AddedAt: time.Now()},
		{Email: "tonitegarsahidi@gmail.com", AddedAt: time.Now()},
	}
	for _, email := range emails {
		var count int64
		db.Model(&model.AdminEmail{}).Where("email = ?", email.Email).Count(&count)
		if count == 0 {
			db.Create(&email)
		}
	}
}

func migrateAndSeedMasters(db *gorm.DB) {
	// Agencies
	if err := db.AutoMigrate(&model.Agency{}); err == nil {
		var count int64
		db.Model(&model.Agency{}).Count(&count)
		if count == 0 {
			zap.L().Info("Seeding Agencies...")
			agencies := []model.Agency{
				{Code: "FAO", Name: "Food and Agriculture Organization", LogoURL: "/images/agency-logos/fao.png"},
				{Code: "IFAD", Name: "International Fund for Agricultural Development", LogoURL: "/images/agency-logos/ifad.png"},
				{Code: "ILO", Name: "International Labour Organization", LogoURL: "/images/agency-logos/ilo.png"},
				{Code: "IMF", Name: "International Monetary Fund", LogoURL: "/images/agency-logos/imf.png"},
				{Code: "IOM", Name: "International Organization for Migration", LogoURL: "/images/agency-logos/iom.png"},
				{Code: "ITU", Name: "International Telecommunication Union", LogoURL: "/images/agency-logos/itu.png"},
				{Code: "RCO", Name: "Resident Coordinator Office", LogoURL: "/images/agency-logos/rco.png"},
				{Code: "UNAIDS", Name: "Joint United Nations Programme on HIV/AIDS", LogoURL: "/images/agency-logos/unaids.png"},
				{Code: "UN Women", Name: "United Nations Entity for Gender Equality", LogoURL: "/images/agency-logos/un-women.png"},
				{Code: "UNDP", Name: "United Nations Development Programme", LogoURL: "/images/agency-logos/undp.png"},
				{Code: "UNEP", Name: "United Nations Environment Programme", LogoURL: "/images/agency-logos/unep.png"},
				{Code: "UNESCO", Name: "United Nations Educational, Scientific and Cultural Organization", LogoURL: "/images/agency-logos/unesco.png"},
				{Code: "UNFPA", Name: "United Nations Population Fund", LogoURL: "/images/agency-logos/unfpa.png"},
				{Code: "UN-HABITAT", Name: "United Nations Human Settlements Programme", LogoURL: "/images/agency-logos/un-habitat.png"},
				{Code: "UNHCR", Name: "United Nations High Commissioner for Refugees", LogoURL: "/images/agency-logos/unhcr.png"},
				{Code: "UNICEF", Name: "United Nations Children's Fund", LogoURL: "/images/agency-logos/unicef.png"},
				{Code: "UNIDO", Name: "United Nations Industrial Development Organization", LogoURL: "/images/agency-logos/unido.png"},
				{Code: "UNOCHA", Name: "Office for the Coordination of Humanitarian Affairs", LogoURL: "/images/agency-logos/unoocha.png"},
				{Code: "UNODC", Name: "United Nations Office on Drugs and Crime", LogoURL: "/images/agency-logos/unodc.png"},
				{Code: "UNOPS", Name: "United Nations Office for Project Services", LogoURL: "/images/agency-logos/unops.png"},
				{Code: "WFP", Name: "World Food Programme", LogoURL: "/images/agency-logos/wfp.png"},
				{Code: "WHO", Name: "World Health Organization", LogoURL: "/images/agency-logos/who.png"},
				{Code: "World Bank", Name: "World Bank Group", LogoURL: "/images/agency-logos/world-bank.png"},
				{Code: "Global Pulse/PLJ", Name: "UN Global Pulse / Pulse Lab Jakarta", LogoURL: "/images/agency-logos/global-pulse.png"},
			}
			db.Create(&agencies)
		}
	}

	// SDGs
	if err := db.AutoMigrate(&model.Sdg{}); err == nil {
		var count int64
		db.Model(&model.Sdg{}).Count(&count)
		if count == 0 {
			zap.L().Info("Seeding SDGs...")
			sdgs := []model.Sdg{
				{Code: "GOAL 1", Name: "No Poverty", Icon: "/images/SDG-logos/SDG-1.png", Color: "#E5243B"},
				{Code: "GOAL 2", Name: "Zero Hunger", Icon: "/images/SDG-logos/SDG-2.png", Color: "#DDA63A"},
				{Code: "GOAL 3", Name: "Good Health and Well-being", Icon: "/images/SDG-logos/SDG-3.png", Color: "#4C9F38"},
				{Code: "GOAL 4", Name: "Quality Education", Icon: "/images/SDG-logos/SDG-4.png", Color: "#C5192D"},
				{Code: "GOAL 5", Name: "Gender Equality", Icon: "/images/SDG-logos/SDG-5.png", Color: "#FF3A21"},
				{Code: "GOAL 6", Name: "Clean Water and Sanitation", Icon: "/images/SDG-logos/SDG-6.png", Color: "#26BDE2"},
				{Code: "GOAL 7", Name: "Affordable and Clean Energy", Icon: "/images/SDG-logos/SDG-7.png", Color: "#FCC30B"},
				{Code: "GOAL 8", Name: "Decent Work and Economic Growth", Icon: "/images/SDG-logos/SDG-8.png", Color: "#A21942"},
				{Code: "GOAL 9", Name: "Industry, Innovation and Infrastructure", Icon: "/images/SDG-logos/SDG-9.png", Color: "#FD6925"},
				{Code: "GOAL 10", Name: "Reduced Inequalities", Icon: "/images/SDG-logos/SDG-10.png", Color: "#DD1367"},
				{Code: "GOAL 11", Name: "Sustainable Cities and Communities", Icon: "/images/SDG-logos/SDG-11.png", Color: "#FD9D24"},
				{Code: "GOAL 12", Name: "Responsible Consumption and Production", Icon: "/images/SDG-logos/SDG-12.png", Color: "#BF8B2E"},
				{Code: "GOAL 13", Name: "Climate Action", Icon: "/images/SDG-logos/SDG-13.png", Color: "#3F7E44"},
				{Code: "GOAL 14", Name: "Life Below Water", Icon: "/images/SDG-logos/SDG-14.png", Color: "#0A97D9"},
				{Code: "GOAL 15", Name: "Life on Land", Icon: "/images/SDG-logos/SDG-15.png", Color: "#56C02B"},
				{Code: "GOAL 16", Name: "Peace, Justice and Strong Institutions", Icon: "/images/SDG-logos/SDG-16.png", Color: "#00689D"},
				{Code: "GOAL 17", Name: "Partnerships for the Goals", Icon: "/images/SDG-logos/SDG-17.png", Color: "#19486A"},
			}
			db.Create(&sdgs)
		}
	}

	// Sectors
	if err := db.AutoMigrate(&model.Sector{}); err == nil {
		var count int64
		db.Model(&model.Sector{}).Count(&count)
		if count == 0 {
			zap.L().Info("Seeding Sectors...")
			sectors := []model.Sector{
				{Code: "agriculture-food", Name: "Agriculture and Food"},
				{Code: "business-investment", Name: "Business and Investment"},
				{Code: "conflict-violence-radicalism", Name: "Conflict, Violence, and Radicalism"},
				{Code: "covid-19", Name: "COVID-19"},
				{Code: "disability-vulnerability-social-welfare", Name: "Disability and Vulnerability and Social Welfare"},
				{Code: "disaster-emergency", Name: "Disaster and Emergency"},
				{Code: "economic-development", Name: "Economic Development"},
				{Code: "education-culture", Name: "Education and Culture"},
				{Code: "energy-natural-resources", Name: "Energy and Natural Resources"},
				{Code: "environment-climate-change", Name: "Environment and Climate Change"},
				{Code: "fishery-maritime", Name: "Fishery and Maritime"},
				{Code: "gender-child-protection", Name: "Gender and Child Protection"},
				{Code: "governance-corruption", Name: "Governance and Corruption"},
				{Code: "health-nutrition", Name: "Health and Nutrition"},
				{Code: "infrastructure-development", Name: "Infrastructure Development"},
				{Code: "innovation-technology", Name: "Innovation and Technology"},
				{Code: "livelihood-employment", Name: "Livelihood and Employment"},
				{Code: "population-migration", Name: "Population and Migration"},
				{Code: "poverty-social-exclusion", Name: "Poverty and Social Exclusion"},
				{Code: "public-finance-tax-fiscal-policy", Name: "Public Finance, Tax, and Fiscal Policy"},
				{Code: "rural-regional-development", Name: "Rural and Regional Development"},
				{Code: "social-security-protection", Name: "Social Security and Protection"},
				{Code: "urban-development", Name: "Urban Development"},
				{Code: "water-sanitation", Name: "Water and Sanitation"},
			}
			db.Create(&sectors)
		}
	}

	// Languages
	if err := db.AutoMigrate(&model.Language{}); err == nil {
		var count int64
		db.Model(&model.Language{}).Count(&count)
		if count == 0 {
			zap.L().Info("Seeding Languages...")
			languages := []model.Language{
				{Code: "english", Name: "English"},
				{Code: "bahasa", Name: "Bahasa Indonesia"},
				{Code: "french", Name: "French"},
				{Code: "arabic", Name: "Arabic"},
				{Code: "spanish", Name: "Spanish"},
			}
			db.Create(&languages)
		}
	}

	// Joint Programmes
	if err := db.AutoMigrate(&model.JointProgramme{}); err == nil {
		var count int64
		db.Model(&model.JointProgramme{}).Count(&count)
		if count == 0 {
			zap.L().Info("Seeding JointProgrammes...")
			jp := []model.JointProgramme{
				{Code: "adlight", Name: "Advancing Indonesia's Lighting Market to High Efficient Technologies (ADLIGHT)"},
				{Code: "berani", Name: "Better Reproductive Health and Rights for All in Indonesia (BERANI)"},
				{Code: "berani-ii", Name: "Better Sexual and Reproductive Rights for All in Indonesia (BERANI II)"},
				{Code: "chemical-weapons-terrorism", Name: "Building a safer South-East Asia by preventing and responding to the use of chemical weapons by terrorists and other non-state actors in Indonesia (Chemical Weapons Terrorism Project)"},
				{Code: "proklim", Name: "Climate Village Project (PROKLIM)"},
				{Code: "assisst", Name: "Driving Public and Private Capital Towards Green and Social Investments in Indonesia / Accelerating SDGs Investments in Indonesia (ASSIST)"},
				{Code: "empower", Name: "EmPower: Women for Climate-Resilient Societies"},
				{Code: "eljp-covid19", Name: "Employment and Livelihood: An Inclusive Approach to Economic Empowerment of Women and Vulnerable Populations in Indonesia (ELJP, COVID-19)"},
				{Code: "folur", Name: "Food Systems, Land Use and Restoration (FOLUR) Impact Program"},
				{Code: "iom-undp-seed-I", Name: "Global IOM-UNDP Seed Funding Round I"},
				{Code: "iom-undp-seed-II", Name: "Global IOM-UNDP Seed Funding Round II"},
				{Code: "gpi", Name: "Global Peatlands Initiative (GPI)"},
				{Code: "hiv-aids", Name: "HIV/AIDS Joint Programme"},
				{Code: "asp-indonesia", Name: "Leaving No One Behind: Adaptive Social Protection (ASP) for All in Indonesia"},
				{Code: "migration-governance", Name: "Migration Governance for Sustainable Development in Indonesia"},
				{Code: "net-zero-nature-positive", Name: "Net Zero Nature Positive Accelerator"},
				{Code: "page", Name: "Partnership for Action on Green Economy (PAGE)"},
				{Code: "protect", Name: "Preventing Violent Extremism through Promoting Tolerance and Respect for Diversity (PROTECT) Project"},
				{Code: "unwaste", Name: "Project Unwaste: tackling waste trafficking to support a circular economy"},
				{Code: "respect", Name: "RESPECT - Preventing Violence against Women"},
				{Code: "spotlight", Name: "Safe and Fair Migration: Realizing women migrant workers' rights and opportunities in the ASEAN region (SPOTLIGHT)"},
				{Code: "ship-to-shore", Name: "Ship to Shore Rights Project"},
				{Code: "strive-asia", Name: "Strengthening Resilience Against Violent Extremism in Asia (STRIVE Asia)"},
				{Code: "social-protection-covid19", Name: "Supporting the Government of Indonesia and Key Stakeholders to Scale-Up Inclusive Social Protection Programmes in Response to COVID-19"},
				{Code: "shift-indonesia", Name: "Sustainable, Healthy and Inclusive Food Systems Transformation (SHIFT) Indonesia"},
				{Code: "guyub", Name: "Tackling the threat of violent extremism and its impact on human securities in East Java (The Guyub Project)"},
				{Code: "veps-parole", Name: "UN Joint Violent Extremist Prisoners (VEPs) Parole and Probation Project"},
				{Code: "un-redd", Name: "UN-REDD ASEAN Social Forestry initiative (UN-REDD)"},
			}
			db.Create(&jp)
		}
	}

	// LNOB Groups
	if err := db.AutoMigrate(&model.Lnob{}); err == nil {
		var count int64
		db.Model(&model.Lnob{}).Count(&count)
		if count == 0 {
			zap.L().Info("Seeding Lnobs...")
			lnobs := []model.Lnob{
				{Code: "women-girls", Name: "Women and Girls"},
				{Code: "youth-children", Name: "Youth and Children"},
				{Code: "disabilities", Name: "Persons with Disabilities"},
				{Code: "others", Name: "Others"},
			}
			db.Create(&lnobs)
		}
	}

	// Non-UN Partners
	if err := db.AutoMigrate(&model.NonUnPartner{}); err == nil {
		var count int64
		db.Model(&model.NonUnPartner{}).Count(&count)
		if count == 0 {
			zap.L().Info("Seeding NonUnPartners...")
			partners := []model.NonUnPartner{
				{Code: "government", Name: "Government"},
				{Code: "universities", Name: "Universities"},
				{Code: "bilateral-agency", Name: "Bilateral Agency"},
				{Code: "consulting-firm", Name: "Consulting Firm"},
				{Code: "think-tank", Name: "Think Tank / Research Institute"},
				{Code: "international-ngo", Name: "International NGO"},
				{Code: "local-ngo", Name: "Local NGO"},
				{Code: "others", Name: "Others"},
			}
			db.Create(&partners)
		}
	}

	// Organizations
	if err := db.AutoMigrate(&model.Organization{}); err == nil {
		var count int64
		db.Model(&model.Organization{}).Count(&count)
		if count == 0 {
			zap.L().Info("Seeding Organizations...")
			orgs := []model.Organization{
				{Code: "united-nations", Name: "UNITED NATIONS"},
				{Code: "fao", Name: "FAO"},
				{Code: "ifad", Name: "IFAD"},
				{Code: "ilo", Name: "ILO"},
				{Code: "iom", Name: "IOM"},
				{Code: "itu", Name: "ITU"},
				{Code: "unaids", Name: "UNAIDS"},
				{Code: "undp", Name: "UNDP"},
				{Code: "unep", Name: "UNEP"},
				{Code: "unesco", Name: "UNESCO"},
				{Code: "unfpa", Name: "UNFPA"},
				{Code: "unhcr", Name: "UNHCR"},
				{Code: "unicef", Name: "UNICEF"},
				{Code: "unido", Name: "UNIDO"},
				{Code: "unops", Name: "UNOPS"},
				{Code: "unv", Name: "UNV"},
				{Code: "un women", Name: "UN Women"},
				{Code: "wfp", Name: "WFP"},
				{Code: "who", Name: "WHO"},
				{Code: "world bank", Name: "World Bank"},
				{Code: "other", Name: "Other"},
			}
			db.Create(&orgs)
		}
	}

	// Thematic Areas (UNSDCF)
	if err := db.AutoMigrate(&model.ThematicArea{}); err == nil {
		var count int64
		db.Model(&model.ThematicArea{}).Count(&count)
		if count == 0 {
			zap.L().Info("Seeding Thematic Areas...")
			areas := []model.ThematicArea{
				{Code: "inclusive-economic-transformation", Name: "Inclusive Economic Transformation"},
				{Code: "environmental-development-climate-resilience", Name: "Environmental Development and Climate Change Resilience"},
				{Code: "human-development", Name: "Human Development"},
				{Code: "democratic-governance-security", Name: "Democratic Governance and Security"},
			}
			db.Create(&areas)
		}
	}
}
