package repository

import (
	"domesv2/config/database"
	"domesv2/internal/model"
	"domesv2/pkg/errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ReferenceRepository interface {
	GetAgencies() ([]model.Agency, error)
	GetSdgs() ([]model.Sdg, error)
	GetSectors() ([]model.Sector, error)
	GetLanguages() ([]model.Language, error)
	GetJointProgrammes() ([]model.JointProgramme, error)
	GetLnobs() ([]model.Lnob, error)
	GetNonUnPartners() ([]model.NonUnPartner, error)
	GetOrganizations() ([]model.Organization, error)
}

type referenceRepository struct {
	db *gorm.DB
}

func NewReferenceRepository() ReferenceRepository {
	return &referenceRepository{
		db: database.GetDB(),
	}
}

func (r *referenceRepository) GetAgencies() ([]model.Agency, error) {
	var agencies []model.Agency
	if err := r.db.Order("code asc").Find(&agencies).Error; err != nil {
		zap.L().Error("Failed to fetch agencies", zap.Error(err))
		return nil, errors.NewInternalServerError("Failed to fetch agencies", "DATABASE_ERROR")
	}
	return agencies, nil
}

func (r *referenceRepository) GetSdgs() ([]model.Sdg, error) {
	var sdgs []model.Sdg
	// Sort by number from GOAL XX.
	// Since code is "GOAL 1", "GOAL 10", we can sort them properly.
	// In MySQL: ORDER BY CAST(SUBSTRING(code, 6) AS UNSIGNED)
	if err := r.db.Order("CAST(SUBSTRING(code, 6) AS UNSIGNED) asc").Find(&sdgs).Error; err != nil {
		zap.L().Error("Failed to fetch SDGs", zap.Error(err))
		return nil, errors.NewInternalServerError("Failed to fetch SDGs", "DATABASE_ERROR")
	}
	return sdgs, nil
}

func (r *referenceRepository) GetSectors() ([]model.Sector, error) {
	var sectors []model.Sector
	if err := r.db.Order("name asc").Find(&sectors).Error; err != nil {
		zap.L().Error("Failed to fetch sectors", zap.Error(err))
		return nil, errors.NewInternalServerError("Failed to fetch sectors", "DATABASE_ERROR")
	}
	return sectors, nil
}

func (r *referenceRepository) GetLanguages() ([]model.Language, error) {
	var languages []model.Language
	if err := r.db.Order("name asc").Find(&languages).Error; err != nil {
		zap.L().Error("Failed to fetch languages", zap.Error(err))
		return nil, errors.NewInternalServerError("Failed to fetch languages", "DATABASE_ERROR")
	}
	return languages, nil
}

func (r *referenceRepository) GetJointProgrammes() ([]model.JointProgramme, error) {
	var jps []model.JointProgramme
	if err := r.db.Order("name asc").Find(&jps).Error; err != nil {
		zap.L().Error("Failed to fetch joint programmes", zap.Error(err))
		return nil, errors.NewInternalServerError("Failed to fetch joint programmes", "DATABASE_ERROR")
	}
	return jps, nil
}

func (r *referenceRepository) GetLnobs() ([]model.Lnob, error) {
	var lnobs []model.Lnob
	if err := r.db.Order("name asc").Find(&lnobs).Error; err != nil {
		zap.L().Error("Failed to fetch LNOBs", zap.Error(err))
		return nil, errors.NewInternalServerError("Failed to fetch LNOB groups", "DATABASE_ERROR")
	}
	return lnobs, nil
}

func (r *referenceRepository) GetNonUnPartners() ([]model.NonUnPartner, error) {
	var partners []model.NonUnPartner
	if err := r.db.Order("name asc").Find(&partners).Error; err != nil {
		zap.L().Error("Failed to fetch non-UN partners", zap.Error(err))
		return nil, errors.NewInternalServerError("Failed to fetch non-UN partner types", "DATABASE_ERROR")
	}
	return partners, nil
}

func (r *referenceRepository) GetOrganizations() ([]model.Organization, error) {
	var orgs []model.Organization
	if err := r.db.Order("name asc").Find(&orgs).Error; err != nil {
		zap.L().Error("Failed to fetch organizations", zap.Error(err))
		return nil, errors.NewInternalServerError("Failed to fetch organizations", "DATABASE_ERROR")
	}
	return orgs, nil
}
