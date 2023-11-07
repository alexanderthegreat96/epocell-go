package core

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type ColumnInfo struct {
	Letter string
	Name   string
}

type Config struct {
	EpocellFile                         string `json:"epocell_file"`
	EpocellStartsAtRow                  int    `json:"epocell_starts_at_row"`
	StoreNameCellLetter                 string `json:"store_name_cell_letter"`
	StoreSalesUnitsCellLetter           string `json:"store_sales_units_cell_letter"`
	StoreSalesUnitesDeviationCellLetter string `json:"store_sales_units_deviation_cell_letter"`
	RawDataFile                         string `json:"raw_data_file"`
	RawDataCategoryKeywords             string `json:"raw_data_category_keywords"`
	RawCategoryCellLetter               string `json:"raw_data_category_cell_letter"`
	RawDataSumOfSumSales                string `json:"raw_data_sum_of_sum_sales_letter"`
}

type StoreCache struct {
	Index                int      `json:"index"`
	Status               bool     `json:"status"`
	StoreName            string   `json:"store_name"`
	SalesUnits           string   `json:"sales_units"`
	SalesUnitsDeviation  string   `json:"sales_units_deviation"`
	CsvFile              string   `json:"csv_file"`
	Keywords             []string `json:"keywords"`
	SkipKeywords         []string `json:"skip_keywords"`
	SheetId              int      `json:"sheet_id"`
	DataStartsAtRow      int      `json:"data_starts_at_row"`
	CategoryCellLetter   string   `json:"category_cell_letter"`
	ProductCellLetter    string   `json:"product_cell_letter"`
	SumOfSalesCellLetter string   `json:"sum_of_sales_cell_letter"`
}

type TableModel struct {
	Columns []string
	Rows    [][]string
}

type clickableLabel struct {
	*widget.Label
	onClicked func()
}

type StoreKeywordCache struct {
	Keyword string `json:"keyword"`
	Count   int    `json:"count"`
}

type LazyTabContent struct {
	ContentFunc func() fyne.CanvasObject
	Loaded      bool
}
