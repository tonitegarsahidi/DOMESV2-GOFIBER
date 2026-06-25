package service

import (
	"domesv2/internal/model"
	"domesv2/internal/repository"
)

type MasterService interface {
	GetAgencies() ([]model.Agency, error)
	GetSdgs() ([]model.Sdg, error)
	GetSectors() ([]model.Sector, error)
	GetLanguages() ([]model.Language, error)
	GetJointProgrammes() ([]model.JointProgramme, error)
	GetLnobs() ([]model.Lnob, error)
	GetNonUnPartners() ([]model.NonUnPartner, error)
	GetOrganizations() ([]model.Organization, error)
	GetThematicAreas() ([]model.ThematicArea, error)
}

type masterService struct {
	refRepo repository.MasterRepository
}

func NewMasterService(refRepo repository.MasterRepository) MasterService {
	return &masterService{
		refRepo: refRepo,
	}
}

func (s *masterService) GetAgencies() ([]model.Agency, error) {
	return s.refRepo.GetAgencies()
}

func (s *masterService) GetSdgs() ([]model.Sdg, error) {
	return s.refRepo.GetSdgs()
}

func (s *masterService) GetSectors() ([]model.Sector, error) {
	return s.refRepo.GetSectors()
}

func (s *masterService) GetLanguages() ([]model.Language, error) {
	return s.refRepo.GetLanguages()
}

func (s *masterService) GetJointProgrammes() ([]model.JointProgramme, error) {
	return s.refRepo.GetJointProgrammes()
}

func (s *masterService) GetLnobs() ([]model.Lnob, error) {
	return s.refRepo.GetLnobs()
}

func (s *masterService) GetNonUnPartners() ([]model.NonUnPartner, error) {
	return s.refRepo.GetNonUnPartners()
}

func (s *masterService) GetOrganizations() ([]model.Organization, error) {
	return s.refRepo.GetOrganizations()
}

func (s *masterService) GetThematicAreas() ([]model.ThematicArea, error) {
	return s.refRepo.GetThematicAreas()
}
