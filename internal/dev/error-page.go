package dev

import (
	"fmt"
	"net/http"
	
	"github.com/creamsensation/gox"
)

func CreateErrorPage(statusCode int, err error) gox.Node {
	return gox.Html(
		gox.Head(
			gox.Title(gox.Text(http.StatusText(statusCode))),
			gox.Style(
				gox.Element(),
				gox.Raw("*{padding:0;margin:0;box-sizing:border-box;}html{font-family: 'Helvetica', sans-serif;}"),
			),
		),
		gox.Body(
			gox.Style(gox.Text("background:#1e293b;color:white;width:100vw;height:100vh;overflow-y:auto;")),
			gox.Div(
				gox.Style(gox.Text("background:#b91c1c;width:100%;height:100px;display:flex;align-items:center;font-weight:bold;font-size:24px;")),
				gox.Div(
					gox.Style(gox.Text("margin: 0 auto;width:900px;")),
					gox.Text(fmt.Sprintf("Error: %s", http.StatusText(statusCode))),
				),
			),
			gox.Div(
				gox.Style(gox.Text("margin: 0 auto;width:900px;padding: 32px 0;font-size:14px;")),
				gox.Div(
					gox.Style(gox.Text("font-weight:bold;font-size:18px;margin-bottom:8px;")),
					gox.Text("Description:"),
				),
				gox.Div(gox.Text(err)),
				gox.Div(
					gox.Style(gox.Text("font-weight:bold;font-size:18px;margin-bottom:8px;margin-top:32px;")),
					gox.Text("Trace:"),
				),
				gox.Range(
					GetErrorTrace(), func(item ErrorTrace, _ int) gox.Node {
						return gox.Div(
							gox.Style(gox.Text("margin-bottom:8px;")),
							gox.Div(
								gox.Style(gox.Text("margin-bottom:4px;color:#ff0000;")),
								gox.Text(item.Path),
								gox.Text(":"),
								gox.Text(item.Line),
							),
							gox.Div(
								gox.Style(
									gox.Text("border:1px solid white;padding:8px;"),
								),
								gox.Range(
									item.Rows, func(row string, i int) gox.Node {
										return gox.Pre(
											gox.Style(
												gox.If(i == 1, gox.Text("background:#b91c1c;color:white;padding:4px;")),
												gox.Text("margin-bottom:4px;"),
											),
											gox.Raw(row),
										)
									},
								),
							),
						)
					},
				),
			),
		),
	)
}
