package result

import (
	"github.com/creamsensation/cp/env"
	"github.com/creamsensation/cp/internal/constant/contentType"
	"github.com/creamsensation/cp/internal/dev"
	"github.com/creamsensation/gox"
)

type Result struct {
	ResultType  int
	StatusCode  int
	ContentType string
	Content     string
	Data        []byte
}

func CreateHtml(content string, statusCode int) Result {
	return Result{
		ResultType:  Render,
		StatusCode:  statusCode,
		ContentType: contentType.Html,
		Content:     content,
	}
}

func CreateError(content string, statusCode int, err error) Result {
	contentExist := len(content) > 0
	if env.Development() && !contentExist {
		content = gox.Render(dev.CreateErrorPage(statusCode, err))
	}
	if !env.Development() && !contentExist {
		content = err.Error()
	}
	return Result{
		ResultType:  Error,
		StatusCode:  statusCode,
		ContentType: contentType.Html,
		Content:     content,
	}
}

func CreateRedirect(path string, statusCode int) Result {
	return Result{
		ResultType: Redirect,
		StatusCode: statusCode,
		Content:    path,
	}
}

func CreateJson(content string, statusCode int) Result {
	return Result{
		ResultType:  Json,
		StatusCode:  statusCode,
		ContentType: contentType.Json,
		Content:     content,
	}
}

func CreateText(content string, statusCode int) Result {
	return Result{
		ResultType:  Text,
		StatusCode:  statusCode,
		ContentType: contentType.Text,
		Content:     content,
	}
}

func CreateStream(name string, data []byte, statusCode int) Result {
	return Result{
		ResultType:  Stream,
		StatusCode:  statusCode,
		ContentType: contentType.OctetStream,
		Content:     name,
		Data:        data,
	}
}
