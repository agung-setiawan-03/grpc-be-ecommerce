package utils

import "github.com/AgungSetiawan/grpc-be-ecommerce/pb/common"

func SuccessResponse(message string) *common.BaseResponse {
	return &common.BaseResponse{
		StatusCode: 200,
		Message:    message,
	}
}

func ValidationErrorResponse(validationErrors []*common.ValidationError) *common.BaseResponse {
	return &common.BaseResponse{
		StatusCode:       400,
		Message:          "Validation error",
		IsError:          true,
		ValidationErrors: validationErrors,
	}
}
