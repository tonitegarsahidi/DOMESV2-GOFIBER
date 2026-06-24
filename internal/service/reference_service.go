package service

import (
	"domesv2/internal/model"
	"domesv2/internal/repository"
)

type ReferenceService interface {
	GetAgencies() ([]model.Agency, error)
	GetSdgs() ([]model.Sdg, error)
	GetSectors() ([]model.Sector, error)
	GetLanguages() ([]model.Language, error)
	GetJointProgrammes() ([]model.JointProgramme, error)
	GetLnobs() ([]model.Lnob, error)
	GetNonUnPartners() ([]model.NonUnPartner, error)
	GetOrganizations() ([]model.Organization, error)
}

type referenceService struct {
	refRepo repository.ReferenceRepository
}

func NewReferenceService(refRepo repository.ReferenceRepository) ReferenceService {
	return &referenceService{
		refRepo: refRepo,
	}
}

func (s *referenceService) GetAgencies() ([]model.Agency, error) {
	return s.refRepo.GetAgencies()
}

func (s *referenceService) GetSdgs() ([]model.Sdg, error) {
	return s.refRepo.GetSdgs()
}

func (s *referenceService) GetSectors() ([]model.Sector, error) {
	return s.refRepo.GetSectors()
}

func (s *referenceService) GetLanguages() ([]model.Language, error) {
	return s.refRepo.GetLanguages()
}

func (s *referenceService) GetJointProgrammes() ([]model.JointProgramme, error) {
	return s.refRepo.GetJointProgrammes()
}

func (s *referenceService) GetLnobs() ([]model.Lnob, error) {
	return s.refRepo.GetLnobs()
}

func (s *referenceService) GetNonUnPartners() ([]model.NonUnPartner, error) {
	return s.refRepo.GetNonUnPartners()
}

func (s *referenceService) GetOrganizations() ([]model.Organization, error) {
	return s.refRepo.GetOrganizations()
}
