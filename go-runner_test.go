package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

func TestRunCmd(t *testing.T) {
	assert := assert.New(t)
	e := echo.New()
	f := make(url.Values)
	f.Set("args", "-l")

	// Positive test
	req := httptest.NewRequest(echo.POST, "/v2/", strings.NewReader(f.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/run/:command")
	c.SetParamNames("command")
	c.SetParamValues("ls")

	// Assertions
	if assert.NoError(runCmd(c)) {
		assert.Equal(http.StatusOK, rec.Code)
	}

	// Positive JSON test
	q := make(url.Values)
	q.Set("json", "true")
	req = httptest.NewRequest(echo.POST, "/v2/?"+q.Encode(), strings.NewReader(f.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetPath("/run/:command")
	c.SetParamNames("command")
	c.SetParamValues("ls")

	// Assertions
	if assert.NoError(runCmd(c)) {
		assert.Equal(http.StatusOK, rec.Code)
		assert.Regexp(regexp.MustCompile("^\\{.*\\}$"), rec.Body.String())
	}

	// Negative test
	req = httptest.NewRequest(echo.POST, "/v2/", strings.NewReader(f.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetPath("/run/:command")
	c.SetParamNames("command")
	c.SetParamValues("lss") // does not exist

	// Assertions
	if assert.NoError(runCmd(c)) {
		assert.Equal(http.StatusUnprocessableEntity, rec.Code)
		assert.Equal("Can not run comand: lss -l\n", rec.Body.String())
	}

	// Negative JSON test
	req = httptest.NewRequest(echo.POST, "/v2/?"+q.Encode(), strings.NewReader(f.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetPath("/run/:command")
	c.SetParamNames("command")
	c.SetParamValues("lss") // does not exist

	// Assertions
	if assert.NoError(runCmd(c)) {
		assert.Equal(http.StatusUnprocessableEntity, rec.Code)
		assert.Equal("{\"output\":\"Can not run comand: lss -l\\n\"}", rec.Body.String())
	}
}

func TestRunHistory(t *testing.T) {
	assert := assert.New(t)
	e := echo.New()

	// Positive test
	req := httptest.NewRequest(echo.POST, "/v2/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/history/:id")
	c.SetParamNames("id")
	c.SetParamValues("0")

	// Assertions
	if assert.NoError(runHistory(c)) {
		assert.Equal(http.StatusOK, rec.Code)
	}

	// Positive JSON test
	q := make(url.Values)
	q.Set("json", "true")
	req = httptest.NewRequest(echo.POST, "/v2/?"+q.Encode(), nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetPath("/history/:id")
	c.SetParamNames("id")
	c.SetParamValues("0")

	// Assertions
	if assert.NoError(runHistory(c)) {
		assert.Equal(http.StatusOK, rec.Code)
		assert.Regexp(regexp.MustCompile("^\\{.*\\}$"), rec.Body.String())
	}

	// Negative test
	req = httptest.NewRequest(echo.POST, "/v2/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetPath("/run/:id")
	c.SetParamNames("id")
	c.SetParamValues("1000") // does not exist

	// Assertions
	if assert.NoError(runHistory(c)) {
		assert.Equal(http.StatusUnprocessableEntity, rec.Code)
		assert.Equal("ID does not exist\n", rec.Body.String())
	}

	// Negative JSON test
	req = httptest.NewRequest(echo.POST, "/v2/?"+q.Encode(), nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetPath("/run/:id")
	c.SetParamNames("id")
	c.SetParamValues("1000") // does not exist

	// Assertions
	if assert.NoError(runHistory(c)) {
		assert.Equal(http.StatusUnprocessableEntity, rec.Code)
		assert.Equal("{\"output\":\"ID does not exist\\n\"}", rec.Body.String())
	}

	// Run command from previous history no longer exists
	req = httptest.NewRequest(echo.POST, "/v2/?"+q.Encode(), nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetPath("/run/:command")
	c.SetParamNames("command")
	c.SetParamValues("lss")
	runCmd(c)

	invalidCmdID := strconv.Itoa(len(historyEntries) - 1)

	req = httptest.NewRequest(echo.POST, "/v2/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetPath("/run/:id")
	c.SetParamNames("id")
	c.SetParamValues(invalidCmdID) // command does not exist

	// Assertions
	if assert.NoError(runHistory(c)) {
		assert.Equal(http.StatusUnprocessableEntity, rec.Code)
		assert.Equal("Command (lss) with history id 6 does not exist\n", rec.Body.String())
	}

	// JSON test
	req = httptest.NewRequest(echo.POST, "/v2/?"+q.Encode(), nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetPath("/run/:id")
	c.SetParamNames("id")
	c.SetParamValues(invalidCmdID) // command does not exist

	// Assertions
	if assert.NoError(runHistory(c)) {
		assert.Equal(http.StatusUnprocessableEntity, rec.Code)
		assert.Equal("{\"output\":\"Command (lss) with history id 6 does not exist\\n\"}", rec.Body.String())
	}
}

func TestHistory(t *testing.T) {
	e := echo.New()
	assert := assert.New(t)

	// Normal test
	req := httptest.NewRequest(echo.GET, "/v2/history", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Assertions
	if assert.NoError(history(c)) {
		assert.Equal(http.StatusOK, rec.Code)
	}

	// JSON test
	q := make(url.Values)
	q.Set("json", "true")
	req = httptest.NewRequest(echo.GET, "/v2/history?"+q.Encode(), nil)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)

	// Assertions
	if assert.NoError(history(c)) {
		assert.Equal(http.StatusOK, rec.Code)
		assert.Regexp(regexp.MustCompile("^\\[.*\\]$"), rec.Body.String())
	}

	// Empty History
	historyEntries = historyEntries[:0]

	req = httptest.NewRequest(echo.GET, "/v2/history", nil)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)

	// Assertions
	if assert.NoError(history(c)) {
		assert.Equal(http.StatusOK, rec.Code)
		assert.Empty(rec.Body)
	}

	// JSON test
	req = httptest.NewRequest(echo.GET, "/v2/history?"+q.Encode(), nil)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)

	// Assertions
	if assert.NoError(history(c)) {
		assert.Equal(http.StatusOK, rec.Code)
		assert.Regexp(regexp.MustCompile("^\\{.*\\}$"), rec.Body.String())
	}
}
