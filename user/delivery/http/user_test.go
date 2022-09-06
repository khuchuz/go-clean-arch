package http_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/khuchuz/faker"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	userHttp "github.com/khuchuz/go-clean-arch/user/delivery/http"
	"github.com/khuchuz/go-clean-arch/domain"
	"github.com/khuchuz/go-clean-arch/domain/mocks"
)

func TestFetch(t *testing.T) {
	var mockUser domain.User
	err := faker.FakeData(&mockUser)
	assert.NoError(t, err)
	mockUCase := new(mocks.UserUsecase)
	mockListUser := make([]domain.User, 0)
	mockListUser = append(mockListUser, mockUser)
	num := 1
	cursor := "2"
	mockUCase.On("Fetch", mock.Anything, cursor, int64(num)).Return(mockListUser, "10", nil)

	e := echo.New()
	req, err := http.NewRequest(echo.GET, "/user?num=1&cursor="+cursor, strings.NewReader(""))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	handler := userHttp.UserHandler{
		AUsecase: mockUCase,
	}
	err = handler.FetchUser(c)
	require.NoError(t, err)

	responseCursor := rec.Header().Get("X-Cursor")
	assert.Equal(t, "10", responseCursor)
	assert.Equal(t, http.StatusOK, rec.Code)
	mockUCase.AssertExpectations(t)
}

func TestFetchError(t *testing.T) {
	mockUCase := new(mocks.UserUsecase)
	num := 1
	cursor := "2"
	mockUCase.On("Fetch", mock.Anything, cursor, int64(num)).Return(nil, "", domain.ErrInternalServerError)

	e := echo.New()
	req, err := http.NewRequest(echo.GET, "/user?num=1&cursor="+cursor, strings.NewReader(""))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	handler := userHttp.UserHandler{
		AUsecase: mockUCase,
	}
	err = handler.FetchUser(c)
	require.NoError(t, err)

	responseCursor := rec.Header().Get("X-Cursor")
	assert.Equal(t, "", responseCursor)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	mockUCase.AssertExpectations(t)
}

func TestGetByID(t *testing.T) {
	var mockUser domain.User
	err := faker.FakeData(&mockUser)
	assert.NoError(t, err)

	mockUCase := new(mocks.UserUsecase)

	num := int(mockUser.ID)

	mockUCase.On("GetByID", mock.Anything, int64(num)).Return(mockUser, nil)

	e := echo.New()
	req, err := http.NewRequest(echo.GET, "/user/"+strconv.Itoa(num), strings.NewReader(""))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("user/:id")
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(num))
	handler := userHttp.UserHandler{
		AUsecase: mockUCase,
	}
	err = handler.GetByID(c)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, rec.Code)
	mockUCase.AssertExpectations(t)
}

func TestStore(t *testing.T) {
	mockUser := domain.User{
		Name:      "Name",
		Password:  "Password",
		Email:     "Email",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	tempMockUser := mockUser
	tempMockUser.ID = 0
	mockUCase := new(mocks.UserUsecase)

	j, err := json.Marshal(tempMockUser)
	assert.NoError(t, err)

	mockUCase.On("Store", mock.Anything, mock.AnythingOfType("*domain.User")).Return(nil)

	e := echo.New()
	req, err := http.NewRequest(echo.POST, "/user", strings.NewReader(string(j)))
	assert.NoError(t, err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/user")

	handler := userHttp.UserHandler{
		AUsecase: mockUCase,
	}
	err = handler.Store(c)
	require.NoError(t, err)

	assert.Equal(t, http.StatusCreated, rec.Code)
	mockUCase.AssertExpectations(t)
}

func TestDelete(t *testing.T) {
	var mockUser domain.User
	err := faker.FakeData(&mockUser)
	assert.NoError(t, err)

	mockUCase := new(mocks.UserUsecase)

	num := int(mockUser.ID)

	mockUCase.On("Delete", mock.Anything, int64(num)).Return(nil)

	e := echo.New()
	req, err := http.NewRequest(echo.DELETE, "/user/"+strconv.Itoa(num), strings.NewReader(""))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("user/:id")
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(num))
	handler := userHttp.UserHandler{
		AUsecase: mockUCase,
	}
	err = handler.Delete(c)
	require.NoError(t, err)

	assert.Equal(t, http.StatusNoContent, rec.Code)
	mockUCase.AssertExpectations(t)

}
