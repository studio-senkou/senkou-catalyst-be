package dtos

type CreateMerchantRequestDTO struct {
	Name string `json:"name" validate:"required"`
}

func (dto *CreateMerchantRequestDTO) ErrorMessages() map[string]string {
	return map[string]string{
		"name.required": "Merchant name is required",
	}
}

type UpdateMerchantRequestDTO struct {
	Name string `json:"name" validate:"required"`
}

func (dto *UpdateMerchantRequestDTO) ErrorMessages() map[string]string {
	return map[string]string{
		"name.required": "Merchant name is required",
	}
}
