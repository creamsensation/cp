package cp

import (
	"github.com/creamsensation/config"
	"github.com/creamsensation/exporter"
)

type Export interface {
	Csv() exporter.Csv
	Excel() exporter.Excel
	Pdf() exporter.Pdf
}

type export struct {
	exporter exporter.Exporter
	config   config.Export
}

func createExport(config config.Export) export {
	return export{
		exporter: exporter.New(),
		config:   config,
	}
}

func (e export) Csv() exporter.Csv {
	return e.exporter.Csv()
}

func (e export) Excel() exporter.Excel {
	return e.exporter.Excel()
}

func (e export) Pdf() exporter.Pdf {
	return e.exporter.Pdf(e.config.Gotenberg.Endpoint)
}
