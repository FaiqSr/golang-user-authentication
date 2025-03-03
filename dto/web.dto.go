package dto

type (
	WebResopnse struct {
		Message string      `json:"message"`
		Data    interface{} `json:"data"`
	}
)
