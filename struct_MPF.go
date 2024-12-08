package main

import "time"

type MPF_EXPORT struct {
	Data struct {
		Query            string  `json:"query,omitempty"`
		TotalCollections int     `json:"totalCollections,omitempty"`
		Collections      any     `json:"collections,omitempty"`
		TotalPages       int     `json:"totalPages,omitempty"`
		Pages            any     `json:"pages,omitempty"`
		Suggestions      any     `json:"suggestions,omitempty"`
		Total            int     `json:"total,omitempty"`
		Items            []Items `json:"items,omitempty"`
		Facets           Facets  `json:"facets,omitempty"`
		Extra            struct {
			Collections []Collections `json:"collections,omitempty"`
		} `json:"extra,omitempty"`
		Currency      Curency `json:"currency,omitempty"`
		PopularSearch any     `json:"popularSearch,omitempty"`
	} `json:"data,omitempty"`
}

type Facets struct {
	ID                          int    `json:"id,omitempty"`
	Title                       string `json:"title,omitempty"`
	FacetName                   string `json:"facetName,omitempty"`
	Labels                      []any  `json:"labels,omitempty"`
	RangeFormat                 string `json:"rangeFormat,omitempty"`
	Multiple                    int    `json:"multiple,omitempty"`
	Display                     string `json:"display,omitempty"`
	MaxHeight                   string `json:"maxHeight,omitempty"`
	Range                       []int  `json:"range,omitempty"`
	SliderColor                 string `json:"sliderColor,omitempty"`
	SliderValueSymbols          string `json:"sliderValueSymbols,omitempty"`
	SliderPrefix                any    `json:"sliderPrefix,omitempty"`
	SliderSuffix                string `json:"sliderSuffix,omitempty"`
	ShowSliderInputPrefixSuffix bool   `json:"showSliderInputPrefixSuffix,omitempty"`
	Min                         int    `json:"min,omitempty"`
	Max                         int    `json:"max,omitempty"`
	NumericRange                int    `json:"numericRange,omitempty"`
	Sort                        int    `json:"sort,omitempty"`
}

type Items struct {
	ID                int64        `json:"id,omitempty"`
	ProductType       string       `json:"productType,omitempty"`
	Title             string       `json:"title,omitempty"`
	Description       string       `json:"description,omitempty"`
	Collections       []int64      `json:"collections,omitempty"`
	Tags              []string     `json:"tags,omitempty"`
	URLName           string       `json:"urlName,omitempty"`
	Vendor            string       `json:"vendor,omitempty"`
	Date              time.Time    `json:"date,omitempty"`
	Variants          []Variants   `json:"variants,omitempty"`
	SelectedVariantID any          `json:"selectedVariantId,omitempty"`
	Images            []Images     `json:"images,omitempty"`
	Metafields        []Metafields `json:"metafields,omitempty"`
	Options           []any        `json:"options,omitempty"`
	Review            int          `json:"review,omitempty"`
	ReviewCount       int          `json:"reviewCount,omitempty"`
	Extra             any          `json:"extra,omitempty"`
}

type Images struct {
	URL    string `json:"url,omitempty"`
	Alt    any    `json:"alt,omitempty"`
	Width  int    `json:"width,omitempty"`
	Height int    `json:"height,omitempty"`
}

type Variants struct {
	ID             int64  `json:"id,omitempty"`
	Sku            string `json:"sku,omitempty"`
	Barcode        any    `json:"barcode,omitempty"`
	Available      int    `json:"available,omitempty"`
	Price          int    `json:"price,omitempty"`
	Weight         int    `json:"weight,omitempty"`
	CompareAtPrice int    `json:"compareAtPrice,omitempty"`
	ImageIndex     int    `json:"imageIndex,omitempty"`
	Options        []any  `json:"options,omitempty"`
	Metafields     any    `json:"metafields,omitempty"`
	Flags          int    `json:"flags,omitempty"`
}

type Metafields struct {
	Key       string `json:"key,omitempty"`
	Value     string `json:"value,omitempty"`
	Multiple  bool   `json:"multiple,omitempty"`
	ValueType string `json:"valueType,omitempty"`
	Namespace string `json:"namespace,omitempty"`
}

type Curency struct {
	Format            string `json:"format,omitempty"`
	LongFormat        string `json:"longFormat,omitempty"`
	DecimalSeparator  string `json:"decimalSeparator,omitempty"`
	HasDecimals       bool   `json:"hasDecimals,omitempty"`
	ThousandSeparator string `json:"thousandSeparator,omitempty"`
}

type Collections struct {
	ID          int64  `json:"id,omitempty"`
	Title       string `json:"title,omitempty"`
	URLName     string `json:"urlName,omitempty"`
	Description any    `json:"description,omitempty"`
	ImageURL    any    `json:"imageUrl,omitempty"`
	SortOrder   string `json:"sortOrder,omitempty"`
}
