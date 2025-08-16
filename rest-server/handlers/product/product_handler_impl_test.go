package product

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/IainMosima/gomart/domains/product/schema"
	"github.com/IainMosima/gomart/domains/product/service"
	"github.com/IainMosima/gomart/rest-server/dtos"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestNewProductHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := service.NewMockProductService(ctrl)
	handler := NewProductHandler(mockService)

	assert.NotNil(t, handler)
	assert.IsType(t, &ProductHandlerImpl{}, handler)
}

func TestProductHandlerImpl_CreateProduct(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := service.NewMockProductService(ctrl)
	handler := NewProductHandler(mockService)

	gin.SetMode(gin.TestMode)

	productID := uuid.New()
	categoryID := uuid.New()
	now := time.Now()

	tests := []struct {
		name           string
		requestBody    interface{}
		setup          func()
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name: "successful creation",
			requestBody: dtos.CreateProductRequestDTO{
				ProductName:   "iPhone 15",
				Description:   stringPtr("Latest iPhone"),
				Price:         999.99,
				SKU:           "IPHONE15-001",
				StockQuantity: 100,
				CategoryID:    categoryID,
				IsActive:      boolPtr(true),
			},
			setup: func() {
				expectedReq := &schema.CreateProductRequest{
					ProductName:   "iPhone 15",
					Description:   stringPtr("Latest iPhone"),
					Price:         999.99,
					SKU:           "IPHONE15-001",
					StockQuantity: 100,
					CategoryID:    categoryID,
					IsActive:      boolPtr(true),
				}

				mockService.EXPECT().
					CreateProduct(gomock.Any(), expectedReq).
					Return(&schema.ProductResponse{
						ProductID:     productID,
						ProductName:   "iPhone 15",
						Description:   stringPtr("Latest iPhone"),
						Price:         999.99,
						SKU:           "IPHONE15-001",
						StockQuantity: 100,
						CategoryID:    categoryID,
						IsActive:      true,
						CreatedAt:     now,
					}, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedBody: &dtos.ProductResponseDTO{
				ProductID:     productID,
				ProductName:   "iPhone 15",
				Description:   stringPtr("Latest iPhone"),
				Price:         999.99,
				SKU:           "IPHONE15-001",
				StockQuantity: 100,
				CategoryID:    categoryID,
				IsActive:      true,
				CreatedAt:     now,
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
			requestBody: dtos.CreateProductRequestDTO{
				ProductName:   "iPhone 15",
				Price:         999.99,
				SKU:           "IPHONE15-001",
				StockQuantity: 100,
				CategoryID:    categoryID,
			},
			setup: func() {
				mockService.EXPECT().
					CreateProduct(gomock.Any(), gomock.Any()).
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

			c.Request = httptest.NewRequest(http.MethodPost, "/products", bytes.NewBuffer(bodyBytes))
			c.Request.Header.Set("Content-Type", "application/json")

			handler.(*ProductHandlerImpl).CreateProduct(c)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var response interface{}
			json.Unmarshal(w.Body.Bytes(), &response)

			if tt.expectedStatus == http.StatusCreated {
				expectedDTO := tt.expectedBody.(*dtos.ProductResponseDTO)
				actualResponse := response.(map[string]interface{})
				assert.Equal(t, expectedDTO.ProductName, actualResponse["product_name"])
				assert.Equal(t, expectedDTO.SKU, actualResponse["sku"])
				assert.NotEmpty(t, actualResponse["product_id"])
			} else {
				expectedError := tt.expectedBody.(gin.H)
				actualResponse := response.(map[string]interface{})
				assert.Contains(t, actualResponse["error"], expectedError["error"])
			}
		})
	}
}

func TestProductHandlerImpl_GetProduct(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := service.NewMockProductService(ctrl)
	handler := NewProductHandler(mockService)

	gin.SetMode(gin.TestMode)

	productID := uuid.New()
	categoryID := uuid.New()
	now := time.Now()

	tests := []struct {
		name           string
		productID      string
		setup          func()
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name:      "successful get with valid UUID",
			productID: productID.String(),
			setup: func() {
				mockService.EXPECT().
					GetProduct(gomock.Any(), productID).
					Return(&schema.ProductResponse{
						ProductID:     productID,
						ProductName:   "iPhone 15",
						Price:         999.99,
						SKU:           "IPHONE15-001",
						StockQuantity: 100,
						CategoryID:    categoryID,
						IsActive:      true,
						CreatedAt:     now,
					}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: &dtos.ProductResponseDTO{
				ProductID:     productID,
				ProductName:   "iPhone 15",
				Price:         999.99,
				SKU:           "IPHONE15-001",
				StockQuantity: 100,
				CategoryID:    categoryID,
				IsActive:      true,
				CreatedAt:     now,
			},
		},
		{
			name:           "invalid product ID format",
			productID:      "invalid-uuid",
			setup:          func() {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   gin.H{"error": "Invalid product ID"},
		},
		{
			name:      "product not found",
			productID: productID.String(),
			setup: func() {
				mockService.EXPECT().
					GetProduct(gomock.Any(), productID).
					Return(nil, errors.New("product not found"))
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   gin.H{"error": "product not found"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Request = httptest.NewRequest(http.MethodGet, "/products/"+tt.productID, nil)
			c.Params = gin.Params{{Key: "id", Value: tt.productID}}

			handler.(*ProductHandlerImpl).GetProduct(c)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var response interface{}
			json.Unmarshal(w.Body.Bytes(), &response)

			if tt.expectedStatus == http.StatusOK {
				expectedDTO := tt.expectedBody.(*dtos.ProductResponseDTO)
				actualResponse := response.(map[string]interface{})
				assert.Equal(t, expectedDTO.ProductName, actualResponse["product_name"])
				assert.Equal(t, expectedDTO.SKU, actualResponse["sku"])
				assert.NotEmpty(t, actualResponse["product_id"])
			} else {
				expectedError := tt.expectedBody.(gin.H)
				actualResponse := response.(map[string]interface{})
				assert.Contains(t, actualResponse["error"], expectedError["error"])
			}
		})
	}
}

func TestProductHandlerImpl_UpdateProduct(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := service.NewMockProductService(ctrl)
	handler := NewProductHandler(mockService)

	gin.SetMode(gin.TestMode)

	productID := uuid.New()
	categoryID := uuid.New()
	now := time.Now()

	tests := []struct {
		name           string
		productID      string
		requestBody    interface{}
		setup          func()
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name:      "successful update",
			productID: productID.String(),
			requestBody: dtos.UpdateProductRequestDTO{
				ProductName: stringPtr("iPhone 15 Pro"),
				Price:       float64Ptr(1199.99),
			},
			setup: func() {
				mockService.EXPECT().
					UpdateProduct(gomock.Any(), productID, gomock.Any()).
					Return(&schema.ProductResponse{
						ProductID:     productID,
						ProductName:   "iPhone 15 Pro",
						Price:         1199.99,
						SKU:           "IPHONE15-001",
						StockQuantity: 100,
						CategoryID:    categoryID,
						IsActive:      true,
						CreatedAt:     now,
					}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: &dtos.ProductResponseDTO{
				ProductID:     productID,
				ProductName:   "iPhone 15 Pro",
				Price:         1199.99,
				SKU:           "IPHONE15-001",
				StockQuantity: 100,
				CategoryID:    categoryID,
				IsActive:      true,
				CreatedAt:     now,
			},
		},
		{
			name:           "invalid product ID format",
			productID:      "invalid-uuid",
			requestBody:    dtos.UpdateProductRequestDTO{},
			setup:          func() {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   gin.H{"error": "Invalid product ID"},
		},
		{
			name:      "service error",
			productID: productID.String(),
			requestBody: dtos.UpdateProductRequestDTO{
				ProductName: stringPtr("iPhone 15 Pro"),
			},
			setup: func() {
				mockService.EXPECT().
					UpdateProduct(gomock.Any(), productID, gomock.Any()).
					Return(nil, errors.New("update failed"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   gin.H{"error": "update failed"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			bodyBytes, _ := json.Marshal(tt.requestBody)
			c.Request = httptest.NewRequest(http.MethodPut, "/products/"+tt.productID, bytes.NewBuffer(bodyBytes))
			c.Request.Header.Set("Content-Type", "application/json")
			c.Params = gin.Params{{Key: "id", Value: tt.productID}}

			handler.(*ProductHandlerImpl).UpdateProduct(c)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var response interface{}
			json.Unmarshal(w.Body.Bytes(), &response)

			if tt.expectedStatus == http.StatusOK {
				expectedDTO := tt.expectedBody.(*dtos.ProductResponseDTO)
				actualResponse := response.(map[string]interface{})
				assert.Equal(t, expectedDTO.ProductName, actualResponse["product_name"])
				assert.Equal(t, expectedDTO.SKU, actualResponse["sku"])
				assert.NotEmpty(t, actualResponse["product_id"])
			} else {
				expectedError := tt.expectedBody.(gin.H)
				actualResponse := response.(map[string]interface{})
				assert.Contains(t, actualResponse["error"], expectedError["error"])
			}
		})
	}
}

func TestProductHandlerImpl_DeleteProduct(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := service.NewMockProductService(ctrl)
	handler := NewProductHandler(mockService)

	gin.SetMode(gin.TestMode)

	productID := uuid.New()

	tests := []struct {
		name           string
		productID      string
		setup          func()
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name:      "successful delete",
			productID: productID.String(),
			setup: func() {
				mockService.EXPECT().
					DeleteProduct(gomock.Any(), productID).
					Return(nil)
			},
			expectedStatus: http.StatusNoContent,
			expectedBody:   nil,
		},
		{
			name:           "invalid product ID format",
			productID:      "invalid-uuid",
			setup:          func() {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   gin.H{"error": "Invalid product ID"},
		},
		{
			name:      "service error",
			productID: productID.String(),
			setup: func() {
				mockService.EXPECT().
					DeleteProduct(gomock.Any(), productID).
					Return(errors.New("delete failed"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   gin.H{"error": "delete failed"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Request = httptest.NewRequest(http.MethodDelete, "/products/"+tt.productID, nil)
			c.Params = gin.Params{{Key: "id", Value: tt.productID}}

			handler.(*ProductHandlerImpl).DeleteProduct(c)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedBody != nil {
				var response interface{}
				json.Unmarshal(w.Body.Bytes(), &response)

				expectedError := tt.expectedBody.(gin.H)
				actualResponse := response.(map[string]interface{})
				assert.Contains(t, actualResponse["error"], expectedError["error"])
			}
		})
	}
}

// Helper functions
func stringPtr(s string) *string {
	return &s
}

func boolPtr(b bool) *bool {
	return &b
}

func float64Ptr(f float64) *float64 {
	return &f
}

func TestProductHandlerImpl_ListProducts(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := service.NewMockProductService(ctrl)
	handler := NewProductHandler(mockService)

	gin.SetMode(gin.TestMode)

	productID := uuid.New()
	categoryID := uuid.New()
	now := time.Now()

	tests := []struct {
		name           string
		queryParams    string
		setup          func()
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name:        "successful list without filters",
			queryParams: "",
			setup: func() {
				mockService.EXPECT().
					ListProducts(gomock.Any(), gomock.Any()).
					Return(&schema.ProductListResponse{
						Products: []*schema.ProductResponse{{
							ProductID:     productID,
							ProductName:   "iPhone 15",
							Price:         35000.00,
							SKU:           "PHONE-IPHONE15-001",
							StockQuantity: 50,
							CategoryID:    categoryID,
							IsActive:      true,
							CreatedAt:     now,
						}},
						Total:   1,
						Page:    1,
						Limit:   10,
						HasNext: false,
					}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:        "service error",
			queryParams: "",
			setup: func() {
				mockService.EXPECT().
					ListProducts(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("database error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Request = httptest.NewRequest(http.MethodGet, "/products?"+tt.queryParams, nil)

			handler.(*ProductHandlerImpl).ListProducts(c)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestProductHandlerImpl_GetProductsByCategory(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := service.NewMockProductService(ctrl)
	handler := NewProductHandler(mockService)

	gin.SetMode(gin.TestMode)

	productID := uuid.New()
	categoryID := uuid.New()
	now := time.Now()

	tests := []struct {
		name           string
		categoryID     string
		setup          func()
		expectedStatus int
	}{
		{
			name:       "successful get products by category",
			categoryID: categoryID.String(),
			setup: func() {
				mockService.EXPECT().
					GetProductsByCategory(gomock.Any(), categoryID).
					Return(&schema.ProductListResponse{
						Products: []*schema.ProductResponse{{
							ProductID:     productID,
							ProductName:   "iPhone 15",
							Price:         35000.00,
							SKU:           "PHONE-IPHONE15-001",
							StockQuantity: 50,
							CategoryID:    categoryID,
							IsActive:      true,
							CreatedAt:     now,
						}},
						Total:   1,
						Page:    1,
						Limit:   10,
						HasNext: false,
					}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "invalid category ID format",
			categoryID:     "invalid-uuid",
			setup:          func() {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:       "service error",
			categoryID: categoryID.String(),
			setup: func() {
				mockService.EXPECT().
					GetProductsByCategory(gomock.Any(), categoryID).
					Return(nil, errors.New("database error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Request = httptest.NewRequest(http.MethodGet, "/products/category/"+tt.categoryID, nil)
			c.Params = gin.Params{{Key: "categoryId", Value: tt.categoryID}}

			handler.(*ProductHandlerImpl).GetProductsByCategory(c)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}
