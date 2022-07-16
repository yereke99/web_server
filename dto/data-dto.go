package dto

type DataCreateDTO struct {
	Covid       string `json:"covid" form:"covid" binding:"required"`
	Description string `json:"description" form:"description" binding:"required"`
	UserID      uint64 `json:"user_id,omitempty"  form:"user_id,omitempty"`
}

type DataUpdateDTO struct {
	ID          uint64 `json:"id" form:"id" binding:"required"`
	Covid       string `json:"covid" form:"covid" binding:"required"`
	Description string `json:"description" form:"description" binding:"required"`
	UserID      uint64 `json:"user_id,omitempty" form:"user_id,omitempty"`
}
