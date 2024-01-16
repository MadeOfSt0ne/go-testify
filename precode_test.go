package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	correctCity   = "moscow"
	incorrectCity = "test"
)

func TestMainHandlerWhenCountMoreThanTotal(t *testing.T) {
	totalCount := 10
	expected := 4

	res := getResponse(fmt.Sprintf("/cafe?city=%s&count=%d", correctCity, totalCount))
	defer res.Body.Close()

	require.Equal(t, http.StatusOK, res.StatusCode)
	data, err := io.ReadAll(res.Body)
	require.Nil(t, err)
	require.NotEmpty(t, res.Body)
	assert.Equal(t, expected, len(strings.Split(string(data), ",")))
}

func TestMainHandlerWhenRequestIsCorrect(t *testing.T) {
	totalCount := 3
	expected := 3

	res := getResponse(fmt.Sprintf("/cafe?city=%s&count=%d", correctCity, totalCount))
	defer res.Body.Close()

	require.Equal(t, http.StatusOK, res.StatusCode)
	data, err := io.ReadAll(res.Body)
	require.Nil(t, err)
	assert.Equal(t, expected, len(strings.Split(string(data), ",")))
}

func TestMainHandlerWhenCityIsUnknown(t *testing.T) {
	totalCount := 3
	expected := "wrong city value"

	res := getResponse(fmt.Sprintf("/cafe?city=%s&count=%d", incorrectCity, totalCount))
	defer res.Body.Close()

	require.Equal(t, http.StatusBadRequest, res.StatusCode)
	data, err := io.ReadAll(res.Body)
	require.Nil(t, err)
	assert.Equal(t, expected, string(data))
}

func TestMainHandlerWhenCountIsNegative(t *testing.T) {
	totalCount := -1
	expected := "wrong count value"

	res := getResponse(fmt.Sprintf("/cafe?city=%s&count=%d", correctCity, totalCount))
	defer res.Body.Close()

	require.Equal(t, http.StatusBadRequest, res.StatusCode)
	data, err := io.ReadAll(res.Body)
	require.Nil(t, err)
	assert.Equal(t, expected, string(data))
}

func TestMainHandlerWhenCountIsMissing(t *testing.T) {
	totalCount := ""
	expected := "count missing"

	res := getResponse(fmt.Sprintf("/cafe?city=%s&count=%s", correctCity, totalCount))
	defer res.Body.Close()

	require.Equal(t, http.StatusBadRequest, res.StatusCode)
	data, err := io.ReadAll(res.Body)
	require.Nil(t, err)
	assert.Equal(t, expected, string(data))
}

func getResponse(url string) *http.Response {
	req := httptest.NewRequest(http.MethodGet, url, nil)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	return responseRecorder.Result()
}
