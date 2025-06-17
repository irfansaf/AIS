package response

type response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type SuccessResponse struct {
	response
}

type SuccessDataResponse struct {
	response
	Data interface{} `json:"data"`
}

type SuccessMultiDataResponse struct {
	response
	Data      interface{} `json:"data,omitempty"`
	TotalPage int64       `json:"total_page,omitempty"`
	Page      int64       `json:"page,omitempty"`
	Results   int64       `json:"results,omitempty"`
}

type SuccessMultiDataResponseScylla struct {
	response
	Data      interface{} `json:"data"`
	NextState string      `json:"next_state"`
}

type ErrorResponse struct {
	response
}

// NewSuccessDataResponse Utility function to create a success response
func NewSuccessDataResponse(message string, data interface{}) SuccessDataResponse {
	return SuccessDataResponse{
		response: response{
			Success: true,
			Message: message,
		},
		Data: data,
	}
}

func NewSuccessMultiDataResponse(message string, data interface{}, totalPage, totalResults, page int64) SuccessMultiDataResponse {
	return SuccessMultiDataResponse{
		response: response{
			Success: true,
			Message: message,
		},
		Data:      data,
		TotalPage: totalPage,
		Page:      page,
		Results:   totalResults,
	}
}

func NewSuccessMultiDataResponseScylla(message string, data interface{}, nextState string) SuccessMultiDataResponseScylla {
	return SuccessMultiDataResponseScylla{
		response: response{
			Success: true,
			Message: message,
		},
		Data:      data,
		NextState: nextState,
	}
}

// NewSuccessResponse Utility function to create a success response
func NewSuccessResponse(message string) SuccessResponse {
	return SuccessResponse{
		response: response{
			Success: true,
			Message: message,
		},
	}
}

// NewErrorResponse Utility function to create an error response
func NewErrorResponse(message string) ErrorResponse {
	return ErrorResponse{
		response: response{
			Success: false,
			Message: message,
		},
	}
}
