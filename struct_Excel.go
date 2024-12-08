package main

type Descriptions struct {
	Style    string
	Size     string // Dimensions
	Color    string
	Material string
}

type ExcelImages struct {
	URL   string
	Saved bool
	Error string
}

type Excel struct {
	Sku          string
	Price        int
	Inventory    string
	Images       []ExcelImages
	Tags         []string
	Descriptions Descriptions
	Completed    bool
	Duplicated   bool
}
