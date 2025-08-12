package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/IainMosima/gomart/domains/category/schema"
	"github.com/IainMosima/gomart/domains/category/service"
	"github.com/IainMosima/gomart/rest-server/dtos"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestNewCategoryHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := service.NewMockCategoryService(ctrl)
	handler := NewCategoryHandler(mockService)

	assert.NotNil(t, handler)
	assert.IsType(t, &CategoryHandlerImpl{}, handler)
}

func TestCategoryHandlerImpl_CreateCategory(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := service.NewMockCategoryService(ctrl)
	handler := NewCategoryHandler(mockService)

	gin.SetMode(gin.TestMode)

	categoryID := uuid.New()
	parentID := uuid.New()
	now := time.Now()

	tests := []struct {
		name           string
		requestBody    interface{}
		setup          func()
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name: "successful creation with parent",
			requestBody: dtos.CreateCategoryRequestDTO{
				CategoryName: "Electronics",
				ParentID:     &parentID,
			},
			setup: func() {
				expectedReq := &schema.CreateCategoryRequest{
					CategoryName: "Electronics",
					ParentID:     &parentID,
				}

				mockService.EXPECT().
					CreateCategory(gomock.Any(), expectedReq).
					Return(&schema.CategoryResponse{
						CategoryID:   categoryID,
						CategoryName: "Electronics",
						ParentID:     &parentID,
						CreatedAt:    now,
					}, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedBody: &dtos.CategoryResponseDTO{
				CategoryID:   categoryID,
				CategoryName: "Electronics",
				ParentID:     &parentID,
				CreatedAt:    now,
			},
		},
		{
			name: "successful creation without parent (root category)",
			requestBody: dtos.CreateCategoryRequestDTO{
				CategoryName: "Electronics",
				ParentID:     nil,
			},
			setup: func() {
				expectedReq := &schema.CreateCategoryRequest{
					CategoryName: "Electronics",
					ParentID:     nil,
				}

				mockService.EXPECT().
					CreateCategory(gomock.Any(), expectedReq).
					Return(&schema.CategoryResponse{
						CategoryID:   categoryID,
						CategoryName: "Electronics",
						ParentID:     nil,
						CreatedAt:    now,
					}, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedBody: &dtos.CategoryResponseDTO{
				CategoryID:   categoryID,
				CategoryName: "Electronics",
				ParentID:     nil,
				CreatedAt:    now,
			},
		},
		{
			name:           "invalid JSON request body",
			requestBody:    "invalid-json",
			setup:          func() {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   gin.H{"error": "invalid character 'i' looking for beginning of value"},
		},
		{
			name: "service error",
			requestBody: dtos.CreateCategoryRequestDTO{
				CategoryName: "Electronics",
				ParentID:     nil,
			},
			setup: func() {
				mockService.EXPECT().
					CreateCategory(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("database connection failed"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   gin.H{"error": "database connection failed"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			var bodyBytes []byte
			if str, ok := tt.requestBody.(string); ok {
				bodyBytes = []byte(str)
			} else {
				bodyBytes, _ = json.Marshal(tt.requestBody)
			}

			c.Request = httptest.NewRequest(http.MethodPost, "/categories", bytes.NewBuffer(bodyBytes))
			c.Request.Header.Set("Content-Type", "application/json")

			handler.(*CategoryHandlerImpl).CreateCategory(c)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var response interface{}
			json.Unmarshal(w.Body.Bytes(), &response)

			if tt.expectedStatus == http.StatusCreated {
				expectedDTO := tt.expectedBody.(*dtos.CategoryResponseDTO)
				actualResponse := response.(map[string]interface{})
				assert.Equal(t, expectedDTO.CategoryName, actualResponse["category_name"])
				assert.NotEmpty(t, actualResponse["category_id"])
			} else {
				expectedError := tt.expectedBody.(gin.H)
				actualResponse := response.(map[string]interface{})
				assert.Contains(t, actualResponse["error"], expectedError["error"])
			}
		})
	}
}

func TestCategoryHandlerImpl_GetCategory(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := service.NewMockCategoryService(ctrl)
	handler := NewCategoryHandler(mockService)

	gin.SetMode(gin.TestMode)

	categoryID := uuid.New()
	parentID := uuid.New()
	now := time.Now()

	tests := []struct {
		name           string
		categoryID     string
		setup          func()
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name:       "successful get with valid UUID",
			categoryID: categoryID.String(),
			setup: func() {
				mockService.EXPECT().
					GetCategory(gomock.Any(), categoryID).
					Return(&schema.CategoryResponse{
						CategoryID:   categoryID,
						CategoryName: "Electronics",
						ParentID:     &parentID,
						CreatedAt:    now,
					}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: &dtos.CategoryResponseDTO{
				CategoryID:   categoryID,
				CategoryName: "Electronics",
				ParentID:     &parentID,
				CreatedAt:    now,
			},
		},
		{
			name:           "invalid category ID format",
			categoryID:     "invalid-uuid",
			setup:          func() {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   gin.H{"error": "Invalid category ID"},
		},
		{
			name:       "category not found",
			categoryID: categoryID.String(),
			setup: func() {
				mockService.EXPECT().
					GetCategory(gomock.Any(), categoryID).
					Return(nil, errors.New("category not found"))
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   gin.H{"error": "category not found"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/categories/%s", tt.categoryID), nil)
			c.Params = []gin.Param{{Key: "id", Value: tt.categoryID}}

			handler.(*CategoryHandlerImpl).GetCategory(c)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var response interface{}
			json.Unmarshal(w.Body.Bytes(), &response)

			if tt.expectedStatus == http.StatusOK {
				expectedDTO := tt.expectedBody.(*dtos.CategoryResponseDTO)
				actualResponse := response.(map[string]interface{})
				assert.Equal(t, expectedDTO.CategoryName, actualResponse["category_name"])
				assert.Equal(t, expectedDTO.CategoryID.String(), actualResponse["category_id"])
			} else {
				expectedError := tt.expectedBody.(gin.H)
				actualResponse := response.(map[string]interface{})
				assert.Contains(t, actualResponse["error"], expectedError["error"])
			}
		})
	}
}

func TestCategoryHandlerImpl_UpdateCategory(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := service.NewMockCategoryService(ctrl)
	handler := NewCategoryHandler(mockService)

	gin.SetMode(gin.TestMode)

	categoryID := uuid.New()
	parentID := uuid.New()
	now := time.Now()
	updated := now.Add(time.Hour)

	tests := []struct {
		name           string
		categoryID     string
		requestBody    interface{}
		setup          func()
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name:       "successful update",
			categoryID: categoryID.String(),
			requestBody: dtos.UpdateCategoryRequestDTO{
				CategoryName: "Updated Electronics",
			},
			setup: func() {
				expectedReq := &schema.UpdateCategoryRequest{
					CategoryName: "Updated Electronics",
				}

				mockService.EXPECT().
					UpdateCategory(gomock.Any(), categoryID, expectedReq).
					Return(&schema.CategoryResponse{
						CategoryID:   categoryID,
						CategoryName: "Updated Electronics",
						ParentID:     &parentID,
						CreatedAt:    now,
						UpdatedAt:    &updated,
					}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: &dtos.CategoryResponseDTO{
				CategoryID:   categoryID,
				CategoryName: "Updated Electronics",
				ParentID:     &parentID,
				CreatedAt:    now,
				UpdatedAt:    &updated,
			},
		},
		{
			name:           "invalid category ID format",
			categoryID:     "invalid-uuid",
			requestBody:    dtos.UpdateCategoryRequestDTO{CategoryName: "Updated"},
			setup:          func() {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   gin.H{"error": "Invalid category ID"},
		},
		{
			name:           "invalid JSON request body",
			categoryID:     categoryID.String(),
			requestBody:    "invalid-json",
			setup:          func() {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   gin.H{"error": "invalid character 'i' looking for beginning of value"},
		},
		{
			name:       "service error",
			categoryID: categoryID.String(),
			requestBody: dtos.UpdateCategoryRequestDTO{
				CategoryName: "Updated Electronics",
			},
			setup: func() {
				mockService.EXPECT().
					UpdateCategory(gomock.Any(), categoryID, gomock.Any()).
					Return(nil, errors.New("category not found"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   gin.H{"error": "category not found"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			var bodyBytes []byte
			if str, ok := tt.requestBody.(string); ok {
				bodyBytes = []byte(str)
			} else {
				bodyBytes, _ = json.Marshal(tt.requestBody)
			}

			c.Request = httptest.NewRequest(http.MethodPut, fmt.Sprintf("/categories/%s", tt.categoryID), bytes.NewBuffer(bodyBytes))
			c.Request.Header.Set("Content-Type", "application/json")
			c.Params = []gin.Param{{Key: "id", Value: tt.categoryID}}

			handler.(*CategoryHandlerImpl).UpdateCategory(c)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var response interface{}
			json.Unmarshal(w.Body.Bytes(), &response)

			if tt.expectedStatus == http.StatusOK {
				expectedDTO := tt.expectedBody.(*dtos.CategoryResponseDTO)
				actualResponse := response.(map[string]interface{})
				assert.Equal(t, expectedDTO.CategoryName, actualResponse["category_name"])
				assert.Equal(t, expectedDTO.CategoryID.String(), actualResponse["category_id"])
			} else {
				expectedError := tt.expectedBody.(gin.H)
				actualResponse := response.(map[string]interface{})
				assert.Contains(t, actualResponse["error"], expectedError["error"])
			}
		})
	}
}

func TestCategoryHandlerImpl_ListCategories(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := service.NewMockCategoryService(ctrl)
	handler := NewCategoryHandler(mockService)

	gin.SetMode(gin.TestMode)

	cat1ID := uuid.New()
	cat2ID := uuid.New()
	parentID := uuid.New()
	now := time.Now()

	tests := []struct {
		name           string
		method         string
		url            string
		requestBody    interface{}
		setup          func()
		expectedStatus int
		expectedCount  int
		hasError       bool
	}{
		{
			name:   "successful list all categories (GET no params)",
			method: "GET",
			url:    "/categories",
			setup: func() {
				mockService.EXPECT().
					ListCategories(gomock.Any(), &schema.ListCategoriesRequest{}).
					Return(&schema.CategoryListResponse{
						Categories: []*schema.CategoryResponse{
							{
								CategoryID:   cat1ID,
								CategoryName: "Electronics",
								ParentID:     nil,
								CreatedAt:    now,
							},
							{
								CategoryID:   cat2ID,
								CategoryName: "Smartphones",
								ParentID:     &parentID,
								CreatedAt:    now,
							},
						},
						Total: 2,
					}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedCount:  2,
			hasError:       false,
		},
		{
			name:   "successful list root categories only (GET query params)",
			method: "GET",
			url:    "/categories?root_only=true",
			setup: func() {
				mockService.EXPECT().
					ListCategories(gomock.Any(), &schema.ListCategoriesRequest{RootOnly: true}).
					Return(&schema.CategoryListResponse{
						Categories: []*schema.CategoryResponse{
							{
								CategoryID:   cat1ID,
								CategoryName: "Electronics",
								ParentID:     nil,
								CreatedAt:    now,
							},
						},
						Total: 1,
					}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedCount:  1,
			hasError:       false,
		},
		{
			name:   "successful list by parent ID (GET query params)",
			method: "GET",
			url:    fmt.Sprintf("/categories?parent_id=%s", parentID.String()),
			setup: func() {
				mockService.EXPECT().
					ListCategories(gomock.Any(), &schema.ListCategoriesRequest{ParentID: &parentID}).
					Return(&schema.CategoryListResponse{
						Categories: []*schema.CategoryResponse{
							{
								CategoryID:   cat2ID,
								CategoryName: "Smartphones",
								ParentID:     &parentID,
								CreatedAt:    now,
							},
						},
						Total: 1,
					}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedCount:  1,
			hasError:       false,
		},
		{
			name:   "successful list via POST request",
			method: "POST",
			url:    "/categories",
			requestBody: dtos.ListCategoriesRequestDTO{
				RootOnly: true,
			},
			setup: func() {
				mockService.EXPECT().
					ListCategories(gomock.Any(), &schema.ListCategoriesRequest{RootOnly: true}).
					Return(&schema.CategoryListResponse{
						Categories: []*schema.CategoryResponse{
							{
								CategoryID:   cat1ID,
								CategoryName: "Electronics",
								ParentID:     nil,
								CreatedAt:    now,
							},
						},
						Total: 1,
					}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedCount:  1,
			hasError:       false,
		},
		{
			name:   "successful list with empty result",
			method: "GET",
			url:    "/categories",
			setup: func() {
				mockService.EXPECT().
					ListCategories(gomock.Any(), &schema.ListCategoriesRequest{}).
					Return(&schema.CategoryListResponse{
						Categories: []*schema.CategoryResponse{},
						Total:      0,
					}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedCount:  0,
			hasError:       false,
		},
		{
			name:   "invalid query parameter",
			method: "GET",
			url:    "/categories?parent_id=invalid-uuid",
			setup: func() {
			},
			expectedStatus: http.StatusBadRequest,
			expectedCount:  0,
			hasError:       true,
		},
		{
			name:   "service error",
			method: "GET",
			url:    "/categories",
			setup: func() {
				mockService.EXPECT().
					ListCategories(gomock.Any(), &schema.ListCategoriesRequest{}).
					Return(nil, errors.New("database connection failed"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedCount:  0,
			hasError:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			var req *http.Request
			if tt.requestBody != nil {
				jsonBody, _ := json.Marshal(tt.requestBody)
				req = httptest.NewRequest(tt.method, tt.url, bytes.NewBuffer(jsonBody))
				req.Header.Set("Content-Type", "application/json")
			} else {
				req = httptest.NewRequest(tt.method, tt.url, nil)
			}
			c.Request = req

			handler.(*CategoryHandlerImpl).ListCategories(c)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var response map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &response)

			if tt.hasError {
				assert.Contains(t, response, "error")
			} else {
				assert.Contains(t, response, "categories")
				assert.Contains(t, response, "total")

				categories := response["categories"].([]interface{})
				assert.Len(t, categories, tt.expectedCount)

				total := response["total"].(float64)
				assert.Equal(t, float64(tt.expectedCount), total)
			}
		})
	}
}

func TestCategoryHandlerImpl_GetCategoryChildren(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := service.NewMockCategoryService(ctrl)
	handler := NewCategoryHandler(mockService)

	gin.SetMode(gin.TestMode)

	parentID := uuid.New()
	child1ID := uuid.New()
	child2ID := uuid.New()
	now := time.Now()

	tests := []struct {
		name           string
		categoryID     string
		setup          func()
		expectedStatus int
		expectedCount  int
		hasError       bool
	}{
		{
			name:       "successful get children with valid parent ID",
			categoryID: parentID.String(),
			setup: func() {
				mockService.EXPECT().
					GetCategoryChildren(gomock.Any(), &parentID).
					Return(&schema.CategoryListResponse{
						Categories: []*schema.CategoryResponse{
							{
								CategoryID:   child1ID,
								CategoryName: "Smartphones",
								ParentID:     &parentID,
								CreatedAt:    now,
							},
							{
								CategoryID:   child2ID,
								CategoryName: "Laptops",
								ParentID:     &parentID,
								CreatedAt:    now,
							},
						},
						Total: 2,
					}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedCount:  2,
			hasError:       false,
		},
		{
			name:       "successful get root categories (empty parent ID)",
			categoryID: "",
			setup: func() {
				mockService.EXPECT().
					GetCategoryChildren(gomock.Any(), (*uuid.UUID)(nil)).
					Return(&schema.CategoryListResponse{
						Categories: []*schema.CategoryResponse{
							{
								CategoryID:   child1ID,
								CategoryName: "Electronics",
								ParentID:     nil,
								CreatedAt:    now,
							},
						},
						Total: 1,
					}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedCount:  1,
			hasError:       false,
		},
		{
			name:       "no children found",
			categoryID: parentID.String(),
			setup: func() {
				mockService.EXPECT().
					GetCategoryChildren(gomock.Any(), &parentID).
					Return(&schema.CategoryListResponse{
						Categories: []*schema.CategoryResponse{},
						Total:      0,
					}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedCount:  0,
			hasError:       false,
		},
		{
			name:           "invalid parent ID format",
			categoryID:     "invalid-uuid",
			setup:          func() {},
			expectedStatus: http.StatusBadRequest,
			expectedCount:  0,
			hasError:       true,
		},
		{
			name:       "service error",
			categoryID: parentID.String(),
			setup: func() {
				mockService.EXPECT().
					GetCategoryChildren(gomock.Any(), &parentID).
					Return(nil, errors.New("database connection failed"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedCount:  0,
			hasError:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/categories/%s/children", tt.categoryID), nil)
			c.Params = []gin.Param{{Key: "id", Value: tt.categoryID}}

			handler.(*CategoryHandlerImpl).GetCategoryChildren(c)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var response map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &response)

			if tt.hasError && tt.expectedStatus != http.StatusOK {
				assert.Contains(t, response, "error")
			} else {
				assert.Contains(t, response, "categories")
				assert.Contains(t, response, "total")

				categories := response["categories"].([]interface{})
				assert.Len(t, categories, tt.expectedCount)

				total := response["total"].(float64)
				assert.Equal(t, float64(tt.expectedCount), total)
			}
		})
	}
}

func TestCategoryHandlerImpl_GetCategoryAverageProductPrice(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := service.NewMockCategoryService(ctrl)
	handler := NewCategoryHandler(mockService)

	gin.SetMode(gin.TestMode)

	categoryID := uuid.New()

	tests := []struct {
		name           string
		categoryID     string
		setup          func()
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name:       "successful get average price",
			categoryID: categoryID.String(),
			setup: func() {
				mockService.EXPECT().
					GetCategoryAverageProductPrice(gomock.Any(), categoryID).
					Return(&schema.CategoryAverageProductPriceResponse{
						CategoryID:   categoryID,
						CategoryName: "Electronics",
						AveragePrice: 299.99,
					}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: &dtos.CategoryAverageProductPriceResponseDTO{
				CategoryID:   categoryID,
				CategoryName: "Electronics",
				AveragePrice: 299.99,
			},
		},
		{
			name:       "successful get average price with zero (no products)",
			categoryID: categoryID.String(),
			setup: func() {
				mockService.EXPECT().
					GetCategoryAverageProductPrice(gomock.Any(), categoryID).
					Return(&schema.CategoryAverageProductPriceResponse{
						CategoryID:   categoryID,
						CategoryName: "Empty Category",
						AveragePrice: 0.0,
					}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: &dtos.CategoryAverageProductPriceResponseDTO{
				CategoryID:   categoryID,
				CategoryName: "Empty Category",
				AveragePrice: 0.0,
			},
		},
		{
			name:           "invalid category ID format",
			categoryID:     "invalid-uuid",
			setup:          func() {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   gin.H{"error": "Invalid category ID"},
		},
		{
			name:       "category not found",
			categoryID: categoryID.String(),
			setup: func() {
				mockService.EXPECT().
					GetCategoryAverageProductPrice(gomock.Any(), categoryID).
					Return(nil, errors.New("category not found"))
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   gin.H{"error": "category not found"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/categories/%s/average-price", tt.categoryID), nil)
			c.Params = []gin.Param{{Key: "id", Value: tt.categoryID}}

			handler.(*CategoryHandlerImpl).GetCategoryAverageProductPrice(c)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var response interface{}
			json.Unmarshal(w.Body.Bytes(), &response)

			if tt.expectedStatus == http.StatusOK {
				expectedDTO := tt.expectedBody.(*dtos.CategoryAverageProductPriceResponseDTO)
				actualResponse := response.(map[string]interface{})
				assert.Equal(t, expectedDTO.CategoryName, actualResponse["category_name"])
				assert.Equal(t, expectedDTO.CategoryID.String(), actualResponse["category_id"])
				assert.Equal(t, expectedDTO.AveragePrice, actualResponse["average_price"])
			} else {
				expectedError := tt.expectedBody.(gin.H)
				actualResponse := response.(map[string]interface{})
				assert.Contains(t, actualResponse["error"], expectedError["error"])
			}
		})
	}
}
