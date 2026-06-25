package controller

import (
	"github.com/gofiber/fiber/v2"
	"domesv2/internal/service"
	"domesv2/pkg/response"
)

type MasterController struct {
	refService service.MasterService
}

func NewMasterController(refService service.MasterService) *MasterController {
	return &MasterController{
		refService: refService,
	}
}

func (ctrl *MasterController) GetAgencies(c *fiber.Ctx) error {
	result, err := ctrl.refService.GetAgencies()
	if err != nil {
		return response.Error(c, err)
	}
	return response.Success(c, result, "Agencies retrieved successfully")
}

func (ctrl *MasterController) GetSdgs(c *fiber.Ctx) error {
	result, err := ctrl.refService.GetSdgs()
	if err != nil {
		return response.Error(c, err)
	}
	return response.Success(c, result, "SDGs retrieved successfully")
}

func (ctrl *MasterController) GetSectors(c *fiber.Ctx) error {
	result, err := ctrl.refService.GetSectors()
	if err != nil {
		return response.Error(c, err)
	}
	return response.Success(c, result, "Sectors retrieved successfully")
}

func (ctrl *MasterController) GetLanguages(c *fiber.Ctx) error {
	result, err := ctrl.refService.GetLanguages()
	if err != nil {
		return response.Error(c, err)
	}
	return response.Success(c, result, "Languages retrieved successfully")
}

func (ctrl *MasterController) GetJointProgrammes(c *fiber.Ctx) error {
	result, err := ctrl.refService.GetJointProgrammes()
	if err != nil {
		return response.Error(c, err)
	}
	return response.Success(c, result, "Joint programmes retrieved successfully")
}

func (ctrl *MasterController) GetLnobs(c *fiber.Ctx) error {
	result, err := ctrl.refService.GetLnobs()
	if err != nil {
		return response.Error(c, err)
	}
	return response.Success(c, result, "LNOB groups retrieved successfully")
}

func (ctrl *MasterController) GetNonUnPartners(c *fiber.Ctx) error {
	result, err := ctrl.refService.GetNonUnPartners()
	if err != nil {
		return response.Error(c, err)
	}
	return response.Success(c, result, "Non-UN partner types retrieved successfully")
}

func (ctrl *MasterController) GetOrganizations(c *fiber.Ctx) error {
	result, err := ctrl.refService.GetOrganizations()
	if err != nil {
		return response.Error(c, err)
	}
	return response.Success(c, result, "Organizations retrieved successfully")
}

func (ctrl *MasterController) GetThematicAreas(c *fiber.Ctx) error {
	result, err := ctrl.refService.GetThematicAreas()
	if err != nil {
		return response.Error(c, err)
	}
	return response.Success(c, result, "Thematic areas retrieved successfully")
}
