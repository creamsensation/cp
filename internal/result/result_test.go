package result

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"testing"
	
	"github.com/stretchr/testify/assert"
	
	"github.com/creamsensation/cp/internal/constant/contentType"
)

func TestResultCreator(t *testing.T) {
	t.Run(
		"html", func(t *testing.T) {
			htmlResult := "<h1>Hello world</h1>"
			result := CreateHtml(htmlResult, http.StatusOK)
			
			assert.Equal(t, Render, result.ResultType)
			assert.Equal(t, htmlResult, result.Content)
			assert.Equal(t, http.StatusOK, result.StatusCode)
			assert.Equal(t, contentType.Html, result.ContentType)
		},
	)
	t.Run(
		"redirect", func(t *testing.T) {
			redirectResult := "/test/2"
			result := CreateRedirect(redirectResult, http.StatusFound)
			
			assert.Equal(t, Redirect, result.ResultType)
			assert.Equal(t, redirectResult, result.Content)
			assert.Equal(t, http.StatusFound, result.StatusCode)
		},
	)
	t.Run(
		"json", func(t *testing.T) {
			dataMap := map[string]int{
				"a": 1,
				"b": 2,
				"c": 3,
			}
			dataBts, err := json.Marshal(dataMap)
			assert.Nil(t, err)
			result := CreateJson(string(dataBts), http.StatusOK)
			
			assert.Equal(t, Json, result.ResultType)
			assert.Equal(t, string(dataBts), result.Content)
			assert.Equal(t, http.StatusOK, result.StatusCode)
			assert.Equal(t, contentType.Json, result.ContentType)
		},
	)
	t.Run(
		"error", func(t *testing.T) {
			err := errors.New("test not found")
			result := CreateError("", http.StatusNotFound, err)
			
			assert.Equal(t, Error, result.ResultType)
			assert.Equal(t, err.Error(), result.Content)
			assert.Equal(t, http.StatusNotFound, result.StatusCode)
			assert.Equal(t, contentType.Html, result.ContentType)
		},
	)
	t.Run(
		"text", func(t *testing.T) {
			result := CreateText("test", http.StatusOK)
			
			assert.Equal(t, Text, result.ResultType)
			assert.Equal(t, "test", result.Content)
			assert.Equal(t, http.StatusOK, result.StatusCode)
			assert.Equal(t, contentType.Text, result.ContentType)
		},
	)
	t.Run(
		"stream", func(t *testing.T) {
			data := bytes.Repeat([]byte("test"), 1<<8)
			result := CreateStream("test.txt", data, http.StatusOK)
			
			assert.Equal(t, Stream, result.ResultType)
			assert.Equal(t, "test.txt", result.Content)
			assert.Equal(t, data, result.Data)
			assert.Equal(t, http.StatusOK, result.StatusCode)
			assert.Equal(t, contentType.OctetStream, result.ContentType)
		},
	)
}
