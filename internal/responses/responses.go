package responses

type (
	ErrorResponse struct {
		Code        int    `json:"-"`
		Error       string `json:"error"`
		Description string `json:"description"`
		Internal    error  `json:"-"` // Stores the error returned by an external dependency
	}
)
