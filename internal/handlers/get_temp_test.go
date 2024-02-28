package handlers

import (
	"errors"
	"fmt"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/danielzinhors/cloudrun_go/internal/usecases"
	"github.com/danielzinhors/cloudrun_go/internal/usecases/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type GetTempHandlerTestSuite struct {
	suite.Suite
	controller                     *gomock.Controller
	getTemperatureByCepUseCaseMock *mocks.MockGetTemperatureByCepUseCase
}

func TestGetTempHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(GetTempHandlerTestSuite))
}

func (s *GetTempHandlerTestSuite) SetupTest() {
	s.controller = gomock.NewController(s.T())
	s.getTemperatureByCepUseCaseMock = mocks.NewMockGetTemperatureByCepUseCase(s.controller)
}

func (s *GetTempHandlerTestSuite) TestNewGetTempHandler() {
	handler := NewGetTemperatureByCepHandler(s.getTemperatureByCepUseCaseMock)
	assert.NotNil(s.T(), handler)
}

func (s *GetTempHandlerTestSuite) TestHandleInvalidCeps() {

	testCases := []string{
		"",
		"?cep=1",
		"?cep=ABC",
		"?cep=70160-90",
		"?cep=7016090",
	}

	for _, testCase := range testCases {
		url := fmt.Sprintf("http://testing%s", testCase)
		request := httptest.NewRequest("GET", url, nil)
		recorder := httptest.NewRecorder()

		handler := NewGetTemperatureByCepHandler(s.getTemperatureByCepUseCaseMock)
		handler.Handle(recorder, request)

		response := recorder.Result()
		responseStatus := response.StatusCode
		responseBody, err := io.ReadAll(response.Body)

		assert.Nil(s.T(), err)
		assert.Equal(s.T(), 422, responseStatus)
		assert.Equal(s.T(), "invalid zipcode", string(responseBody))
	}
}

func (s *GetTempHandlerTestSuite) TestHandleUseCaseErrZipcodeNotFound() {

	s.getTemperatureByCepUseCaseMock.EXPECT().
		Execute(gomock.Any(), gomock.Any()).Return(nil, errors.New("can not found zipcode"))

	request := httptest.NewRequest("GET", "http://testing?cep=01451-000", nil)
	recorder := httptest.NewRecorder()

	handler := NewGetTemperatureByCepHandler(s.getTemperatureByCepUseCaseMock)
	handler.Handle(recorder, request)

	response := recorder.Result()
	responseStatus := response.StatusCode
	responseBody, err := io.ReadAll(response.Body)

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 404, responseStatus)
	assert.Equal(s.T(), "can not found zipcode", string(responseBody))
}

func (s *GetTempHandlerTestSuite) TestHandleUseCaseErrGeneric() {

	s.getTemperatureByCepUseCaseMock.EXPECT().
		Execute(gomock.Any(), gomock.Any()).Return(nil, errors.New("generic error"))

	request := httptest.NewRequest("GET", "http://testing?cep=01451-000", nil)
	recorder := httptest.NewRecorder()

	handler := NewGetTemperatureByCepHandler(s.getTemperatureByCepUseCaseMock)
	handler.Handle(recorder, request)

	response := recorder.Result()
	responseStatus := response.StatusCode
	responseBody, err := io.ReadAll(response.Body)

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 500, responseStatus)
	assert.Equal(s.T(), "internal server error", string(responseBody))
}

func (s *GetTempHandlerTestSuite) TestHandleSuccess() {

	useCaseOutput := &usecases.TempOutput{
		TemperatureCelsius:    28.511111,
		TemperatureFahrenheit: 29.511111,
		TemperatureKelvin:     30.511111,
	}

	request := httptest.NewRequest("GET", "http://testing?cep=01451-000", nil)
	recorder := httptest.NewRecorder()

	s.getTemperatureByCepUseCaseMock.EXPECT().
		Execute(request.Context(), &usecases.TempInput{Cep: "01451-000"}).Return(useCaseOutput, nil)

	handler := NewGetTemperatureByCepHandler(s.getTemperatureByCepUseCaseMock)
	handler.Handle(recorder, request)

	response := recorder.Result()
	responseStatus := response.StatusCode
	responseBody, err := io.ReadAll(response.Body)

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 200, responseStatus)
	assert.Equal(s.T(), "{\"temp_C\":28.5,\"temp_F\":29.5,\"temp_K\":30.5}\n", string(responseBody))
}
