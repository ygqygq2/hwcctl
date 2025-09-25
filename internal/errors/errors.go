package errors

import (
	"fmt"
	"net/http"
	"strings"
)

// ErrorType 错误类型
type ErrorType string

const (
	// 认证错误
	ErrorTypeAuth ErrorType = "AuthenticationError"
	// 授权错误
	ErrorTypePermission ErrorType = "PermissionError"
	// 网络错误
	ErrorTypeNetwork ErrorType = "NetworkError"
	// 参数错误
	ErrorTypeValidation ErrorType = "ValidationError"
	// 服务器错误
	ErrorTypeServer ErrorType = "ServerError"
	// 资源未找到
	ErrorTypeNotFound ErrorType = "NotFoundError"
	// 限流错误
	ErrorTypeThrottle ErrorType = "ThrottleError"
	// 未知错误
	ErrorTypeUnknown ErrorType = "UnknownError"
)

// HuaweiCloudError 华为云错误结构
type HuaweiCloudError struct {
	Type       ErrorType `json:"type"`
	Code       string    `json:"code"`
	Message    string    `json:"message"`
	Details    string    `json:"details,omitempty"`
	RequestID  string    `json:"request_id,omitempty"`
	StatusCode int       `json:"status_code,omitempty"`
	Retryable  bool      `json:"retryable"`
}

// Error 实现 error 接口
func (e *HuaweiCloudError) Error() string {
	if e.RequestID != "" {
		return fmt.Sprintf("[%s] %s (RequestID: %s)", e.Code, e.Message, e.RequestID)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// IsRetryable 判断是否可重试
func (e *HuaweiCloudError) IsRetryable() bool {
	return e.Retryable
}

// NewError 创建新的华为云错误
func NewError(errorType ErrorType, code, message string) *HuaweiCloudError {
	return &HuaweiCloudError{
		Type:      errorType,
		Code:      code,
		Message:   message,
		Retryable: isRetryableError(errorType, code),
	}
}

// NewErrorWithDetails 创建带详细信息的错误
func NewErrorWithDetails(errorType ErrorType, code, message, details string) *HuaweiCloudError {
	return &HuaweiCloudError{
		Type:      errorType,
		Code:      code,
		Message:   message,
		Details:   details,
		Retryable: isRetryableError(errorType, code),
	}
}

// NewHTTPError 从 HTTP 响应创建错误
func NewHTTPError(statusCode int, body string) *HuaweiCloudError {
	errorType := getErrorTypeFromStatusCode(statusCode)
	code := fmt.Sprintf("HTTP%d", statusCode)
	message := getMessageFromStatusCode(statusCode)
	
	return &HuaweiCloudError{
		Type:       errorType,
		Code:       code,
		Message:    message,
		Details:    body,
		StatusCode: statusCode,
		Retryable:  isRetryableStatusCode(statusCode),
	}
}

// ParseHuaweiCloudError 解析华为云 API 错误响应
func ParseHuaweiCloudError(statusCode int, body string) *HuaweiCloudError {
	// 这里可以解析华为云特定的错误格式
	// 例如：{"error": {"code": "CDN.0001", "message": "Invalid parameter"}}
	
	// 尝试从响应体中提取错误信息
	if strings.Contains(body, "Invalid") || strings.Contains(body, "invalid") {
		return &HuaweiCloudError{
			Type:       ErrorTypeValidation,
			Code:       "InvalidParameter",
			Message:    "请求参数无效",
			Details:    body,
			StatusCode: statusCode,
			Retryable:  false,
		}
	}
	
	if strings.Contains(body, "Unauthorized") || strings.Contains(body, "unauthorized") {
		return &HuaweiCloudError{
			Type:       ErrorTypeAuth,
			Code:       "Unauthorized",
			Message:    "认证失败，请检查访问密钥",
			Details:    body,
			StatusCode: statusCode,
			Retryable:  false,
		}
	}
	
	if strings.Contains(body, "Forbidden") || strings.Contains(body, "forbidden") {
		return &HuaweiCloudError{
			Type:       ErrorTypePermission,
			Code:       "Forbidden",
			Message:    "权限不足，请检查账户权限",
			Details:    body,
			StatusCode: statusCode,
			Retryable:  false,
		}
	}
	
	// 默认错误
	return NewHTTPError(statusCode, body)
}

// getErrorTypeFromStatusCode 根据状态码确定错误类型
func getErrorTypeFromStatusCode(statusCode int) ErrorType {
	switch {
	case statusCode == http.StatusUnauthorized:
		return ErrorTypeAuth
	case statusCode == http.StatusForbidden:
		return ErrorTypePermission
	case statusCode == http.StatusNotFound:
		return ErrorTypeNotFound
	case statusCode == http.StatusTooManyRequests:
		return ErrorTypeThrottle
	case statusCode >= 400 && statusCode < 500:
		return ErrorTypeValidation
	case statusCode >= 500:
		return ErrorTypeServer
	default:
		return ErrorTypeUnknown
	}
}

// getMessageFromStatusCode 根据状态码获取错误消息
func getMessageFromStatusCode(statusCode int) string {
	switch statusCode {
	case http.StatusBadRequest:
		return "请求参数错误"
	case http.StatusUnauthorized:
		return "认证失败"
	case http.StatusForbidden:
		return "权限不足"
	case http.StatusNotFound:
		return "资源不存在"
	case http.StatusTooManyRequests:
		return "请求频率过高"
	case http.StatusInternalServerError:
		return "服务器内部错误"
	case http.StatusBadGateway:
		return "网关错误"
	case http.StatusServiceUnavailable:
		return "服务不可用"
	case http.StatusGatewayTimeout:
		return "网关超时"
	default:
		return fmt.Sprintf("HTTP 错误 %d", statusCode)
	}
}

// isRetryableError 判断错误是否可重试
func isRetryableError(errorType ErrorType, code string) bool {
	switch errorType {
	case ErrorTypeNetwork, ErrorTypeServer, ErrorTypeThrottle:
		return true
	case ErrorTypeAuth, ErrorTypePermission, ErrorTypeValidation, ErrorTypeNotFound:
		return false
	default:
		// 特定错误码的重试策略
		retryableCodes := []string{
			"RequestTimeout",
			"ServiceUnavailable",
			"InternalError",
			"ThrottleException",
		}
		
		for _, retryableCode := range retryableCodes {
			if strings.Contains(code, retryableCode) {
				return true
			}
		}
		
		return false
	}
}

// isRetryableStatusCode 判断 HTTP 状态码是否可重试
func isRetryableStatusCode(statusCode int) bool {
	retryableStatusCodes := []int{
		http.StatusInternalServerError,     // 500
		http.StatusBadGateway,              // 502
		http.StatusServiceUnavailable,      // 503
		http.StatusGatewayTimeout,          // 504
		http.StatusTooManyRequests,         // 429
	}
	
	for _, code := range retryableStatusCodes {
		if statusCode == code {
			return true
		}
	}
	
	return false
}

// Common error constructors

// NewAuthError 创建认证错误
func NewAuthError(message string) *HuaweiCloudError {
	return NewError(ErrorTypeAuth, "AuthenticationFailed", message)
}

// NewPermissionError 创建权限错误
func NewPermissionError(message string) *HuaweiCloudError {
	return NewError(ErrorTypePermission, "PermissionDenied", message)
}

// NewValidationError 创建参数验证错误
func NewValidationError(message string) *HuaweiCloudError {
	return NewError(ErrorTypeValidation, "InvalidParameter", message)
}

// NewNetworkError 创建网络错误
func NewNetworkError(message string) *HuaweiCloudError {
	return NewError(ErrorTypeNetwork, "NetworkError", message)
}

// NewServerError 创建服务器错误
func NewServerError(message string) *HuaweiCloudError {
	return NewError(ErrorTypeServer, "ServerError", message)
}

// NewNotFoundError 创建资源未找到错误
func NewNotFoundError(resource string) *HuaweiCloudError {
	return NewError(ErrorTypeNotFound, "NotFound", fmt.Sprintf("资源 '%s' 不存在", resource))
}
