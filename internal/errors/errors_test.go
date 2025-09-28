package errors

import (
	"fmt"
	"strings"
	"testing"
)

func TestNewError(t *testing.T) {
	err := NewError(ErrorTypeValidation, "TEST_ERROR", "测试错误")
	if err.Code != "TEST_ERROR" {
		t.Errorf("期望错误代码为 'TEST_ERROR'，实际为 '%s'", err.Code)
	}
	if err.Message != "测试错误" {
		t.Errorf("期望错误消息为 '测试错误'，实际为 '%s'", err.Message)
	}
}

func TestNewErrorWithDetails(t *testing.T) {
	err := NewErrorWithDetails(ErrorTypeValidation, "CODE", "message", "详细信息")

	if err.Code != "CODE" {
		t.Errorf("期望错误代码为 'CODE'，实际为 '%s'", err.Code)
	}
	if err.Details != "详细信息" {
		t.Errorf("期望详情为 '详细信息'，实际为 '%s'", err.Details)
	}
}

func TestNewHTTPError(t *testing.T) {
	err := NewHTTPError(404, "Not Found")
	if err.StatusCode != 404 {
		t.Errorf("期望状态码为 404，实际为 %d", err.StatusCode)
	}
	if err.Code != "HTTP404" {
		t.Errorf("期望错误代码为 'HTTP404'，实际为 '%s'", err.Code)
	}
}

func TestIsRetryable(t *testing.T) {
	retryableErr := &HuaweiCloudError{Type: ErrorTypeNetwork, Code: "NETWORK_ERROR", Retryable: true}
	if !retryableErr.IsRetryable() {
		t.Error("NETWORK_ERROR应该是可重试的")
	}

	nonRetryableErr := &HuaweiCloudError{Type: ErrorTypeAuth, Code: "AUTH_ERROR", Retryable: false}
	if nonRetryableErr.IsRetryable() {
		t.Error("AUTH_ERROR不应该是可重试的")
	}
}

func TestErrorString(t *testing.T) {
	err := NewValidationError("测试错误")
	expected := "[InvalidParameter] 测试错误"
	if err.Error() != expected {
		t.Errorf("期望错误字符串为 '%s'，实际为 '%s'", expected, err.Error())
	}
}

func TestErrorWithRequestID(t *testing.T) {
	// 测试带RequestID的错误格式 - 对调试API调用很重要
	err := NewValidationError("参数错误")
	err.RequestID = "req-12345"

	expected := "[InvalidParameter] 参数错误 (RequestID: req-12345)"
	if err.Error() != expected {
		t.Errorf("带RequestID的错误格式不正确，期望: %s, 实际: %s", expected, err.Error())
	}
}

func TestErrorDetailedInformation(t *testing.T) {
	// 测试错误的详细信息管理 - 用于记录完整的错误上下文
	rootCause := NewValidationError("原始验证错误")
	wrappedError := NewErrorWithDetails(ErrorTypeServer, "ServerError", "服务器处理失败",
		fmt.Sprintf("原因: %s, 阶段: 数据验证", rootCause.Error()))

	// 验证错误包装
	if wrappedError.Error() == "" {
		t.Error("包装后的错误不应该为空")
	}

	// 验证详细信息是否正确设置
	if wrappedError.Details == "" {
		t.Error("错误详细信息不应该为空")
	}

	// 验证详细信息包含原始错误
	if !strings.Contains(wrappedError.Details, rootCause.Error()) {
		t.Errorf("错误详细信息应该包含原始错误: %s", rootCause.Error())
	}

	// 验证错误类型
	if wrappedError.Type != ErrorTypeServer {
		t.Errorf("错误类型不匹配，期望: %s, 实际: %s", ErrorTypeServer, wrappedError.Type)
	}
}

func TestErrorClassificationBusinessLogic(t *testing.T) {
	// 测试错误分类的业务逻辑 - 决定如何处理不同类型的错误
	testCases := []struct {
		name         string
		errorFunc    func(string) *HuaweiCloudError
		expectedCode string
		shouldRetry  bool
		description  string
	}{
		{
			name:         "认证错误",
			errorFunc:    NewAuthError,
			expectedCode: "AuthenticationFailed",
			shouldRetry:  false, // 认证错误通常不应该重试
			description:  "认证失败通常需要用户干预",
		},
		{
			name:         "权限错误",
			errorFunc:    NewPermissionError,
			expectedCode: "PermissionDenied",
			shouldRetry:  false, // 权限错误不应该重试
			description:  "权限不足需要配置更改",
		},
		{
			name:         "验证错误",
			errorFunc:    NewValidationError,
			expectedCode: "InvalidParameter",
			shouldRetry:  false, // 参数错误不应该重试
			description:  "参数错误需要修正输入",
		},
		{
			name:         "网络错误",
			errorFunc:    NewNetworkError,
			expectedCode: "NetworkError",
			shouldRetry:  true, // 网络错误可以重试
			description:  "网络问题可能是临时的",
		},
		{
			name:         "服务器错误",
			errorFunc:    NewServerError,
			expectedCode: "ServerError",
			shouldRetry:  true, // 服务器错误可以重试
			description:  "服务器错误可能是临时的",
		},
		{
			name:         "资源未找到",
			errorFunc:    NewNotFoundError,
			expectedCode: "NotFound",
			shouldRetry:  false, // 资源不存在不应该重试
			description:  "资源不存在需要检查输入",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.errorFunc("测试消息")

			// 验证错误代码
			if err.Code != tc.expectedCode {
				t.Errorf("错误代码不匹配，期望: %s, 实际: %s", tc.expectedCode, err.Code)
			}

			// 验证重试策略
			shouldRetry := err.IsRetryable()
			if shouldRetry != tc.shouldRetry {
				t.Errorf("重试策略不匹配 (%s)，期望: %v, 实际: %v",
					tc.description, tc.shouldRetry, shouldRetry)
			}

			// 验证错误格式
			errorStr := err.Error()
			if !strings.Contains(errorStr, tc.expectedCode) {
				t.Errorf("错误字符串应该包含错误代码 %s", tc.expectedCode)
			}
		})
	}
}

func TestParseHuaweiCloudErrorPatterns(t *testing.T) {
	// 测试华为云错误解析的不同模式 - 核心业务逻辑
	testCases := []struct {
		name         string
		statusCode   int
		body         string
		expectedType ErrorType
		expectedCode string
		description  string
	}{
		{
			name:         "无效参数响应",
			statusCode:   400,
			body:         `{"error": {"message": "Invalid parameter value"}}`,
			expectedType: ErrorTypeValidation,
			expectedCode: "InvalidParameter",
			description:  "应该识别参数验证错误",
		},
		{
			name:         "认证失败响应",
			statusCode:   401,
			body:         `{"error": {"message": "Unauthorized access"}}`,
			expectedType: ErrorTypeAuth,
			expectedCode: "Unauthorized",
			description:  "应该识别认证错误",
		},
		{
			name:         "权限不足响应",
			statusCode:   403,
			body:         `{"error": {"message": "Forbidden operation"}}`,
			expectedType: ErrorTypePermission,
			expectedCode: "Forbidden",
			description:  "应该识别权限错误",
		},
		{
			name:         "服务器错误响应",
			statusCode:   500,
			body:         `{"error": {"message": "Internal server error"}}`,
			expectedType: ErrorTypeServer,
			expectedCode: "HTTP500",
			description:  "应该识别服务器错误",
		},
		{
			name:         "通用HTTP错误",
			statusCode:   429,
			body:         "Too many requests",
			expectedType: ErrorTypeThrottle,
			expectedCode: "HTTP429",
			description:  "应该识别限流错误",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := ParseHuaweiCloudError(tc.statusCode, tc.body)

			if result.Type != tc.expectedType {
				t.Errorf("%s: 错误类型不匹配，期望 %s, 实际 %s",
					tc.description, tc.expectedType, result.Type)
			}

			if result.Code != tc.expectedCode {
				t.Errorf("%s: 错误代码不匹配，期望 %s, 实际 %s",
					tc.description, tc.expectedCode, result.Code)
			}

			if result.StatusCode != tc.statusCode {
				t.Errorf("%s: 状态码不匹配，期望 %d, 实际 %d",
					tc.description, tc.statusCode, result.StatusCode)
			}
		})
	}
}

func TestHTTPErrorMapping(t *testing.T) {
	// 测试HTTP状态码到错误类型的映射 - 这是业务逻辑的核心部分
	testCases := []struct {
		statusCode   int
		expectedType ErrorType
		shouldRetry  bool
		description  string
	}{
		{400, ErrorTypeValidation, false, "400应该映射到验证错误且不可重试"},
		{401, ErrorTypeAuth, false, "401应该映射到认证错误且不可重试"},
		{403, ErrorTypePermission, false, "403应该映射到权限错误且不可重试"},
		{404, ErrorTypeNotFound, false, "404应该映射到未找到错误且不可重试"},
		{429, ErrorTypeThrottle, true, "429应该映射到限流错误且可重试"},
		{500, ErrorTypeServer, true, "500应该映射到服务器错误且可重试"},
		{502, ErrorTypeServer, true, "502应该映射到服务器错误且可重试"},
		{503, ErrorTypeServer, true, "503应该映射到服务器错误且可重试"},
		{504, ErrorTypeServer, true, "504应该映射到服务器错误且可重试"},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("HTTP%d", tc.statusCode), func(t *testing.T) {
			err := NewHTTPError(tc.statusCode, "test body")

			if err.Type != tc.expectedType {
				t.Errorf("%s: 类型不匹配，期望 %s, 实际 %s",
					tc.description, tc.expectedType, err.Type)
			}

			if err.IsRetryable() != tc.shouldRetry {
				t.Errorf("%s: 重试策略不匹配，期望 %v, 实际 %v",
					tc.description, tc.shouldRetry, err.IsRetryable())
			}
		})
	}
}

func TestErrorTypeConsistency(t *testing.T) {
	// 测试错误类型的一致性 - 确保所有构造函数都返回正确的类型
	typeTests := []struct {
		constructor  func(string) *HuaweiCloudError
		expectedType ErrorType
		name         string
	}{
		{NewAuthError, ErrorTypeAuth, "认证错误"},
		{NewPermissionError, ErrorTypePermission, "权限错误"},
		{NewValidationError, ErrorTypeValidation, "验证错误"},
		{NewNetworkError, ErrorTypeNetwork, "网络错误"},
		{NewServerError, ErrorTypeServer, "服务器错误"},
		{NewNotFoundError, ErrorTypeNotFound, "未找到错误"},
	}

	for _, test := range typeTests {
		t.Run(test.name, func(t *testing.T) {
			err := test.constructor("test message")
			if err.Type != test.expectedType {
				t.Errorf("%s类型不匹配，期望: %s, 实际: %s",
					test.name, test.expectedType, err.Type)
			}
		})
	}
}

func TestNotFoundErrorFormatting(t *testing.T) {
	// 测试资源未找到错误的特殊格式 - 包含资源名称
	resourceName := "用户123"
	err := NewNotFoundError(resourceName)

	expectedMessage := fmt.Sprintf("资源 '%s' 不存在", resourceName)
	if err.Message != expectedMessage {
		t.Errorf("未找到错误格式不正确，期望: %s, 实际: %s", expectedMessage, err.Message)
	}

	if err.Code != "NotFound" {
		t.Errorf("未找到错误代码不正确，期望: NotFound, 实际: %s", err.Code)
	}
}
