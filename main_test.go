package main

import (
	"reflect"
	"testing"
)

func TestExtractDescriptions_01(t *testing.T) {
	got := extractDescriptions("<span data-mce-fragment=\"1\">Dimensions: H:19\" W:28\" D:28\" </span><br data-mce-fragment=\"1\"><span data-mce-fragment=\"1\">Colour: White / Cream</span><br data-mce-fragment=\"1\"><span data-mce-fragment=\"1\">Material: Concrete</span>")

	wantColor := "WHITE / CREAM"
	if got.Color != wantColor {
		t.Errorf("got %q, wanted %q", got.Color, wantColor)
	}

	wantSize := "H:19 W:28 D:28"
	if got.Size != wantSize {
		t.Errorf("got %q, wanted %q", got.Size, wantSize)
	}

	wantMaterial := "CONCRETE"
	if got.Material != wantMaterial {
		t.Errorf("got %q, wanted %q", got.Material, wantMaterial)
	}
}

func TestExtractDescriptions_02(t *testing.T) {
	got := extractDescriptions("<p>Dimensions: 8.5\"w x 18\"d x 21\"h<br><span style=\"font-size: 0.875rem;\">Material: Brass / Iron<br></span><span data-mce-fragment=\"1\">Style / Era: Georgian / Early to Mid 20th Century<br><br></span></p>")

	wantSize := "8.5W X 18D X 21H"
	if got.Size != wantSize {
		t.Errorf("got %q, wanted %q", got.Size, wantSize)
	}

	wantStyle := "GEORGIAN / EARLY TO MID 20TH CENTURY"
	if got.Style != wantStyle {
		t.Errorf("got %q, wanted %q", got.Style, wantStyle)
	}

	wantMaterial := "BRASS / IRON"
	if got.Material != wantMaterial {
		t.Errorf("got %q, wanted %q", got.Material, wantMaterial)
	}
}

func TestExtractDescriptions_03(t *testing.T) {
	got := extractDescriptions("<p>Style: Persian\n<br>Size: 9'8 W x 16'2 L\n<br>Main Colour: Blue/ Purple</p>")

	wantSize := "9'8 W X 16'2 L"
	if got.Size != wantSize {
		t.Errorf("got %q, wanted %q", got.Size, wantSize)
	}

	wantColor := "BLUE/ PURPLE"
	if got.Color != wantColor {
		t.Errorf("got %q, wanted %q", got.Color, wantColor)
	}
}

func TestExtractDescriptions_04(t *testing.T) {
	got := extractDescriptions("<span data-mce-fragment=\"1\"> Dimensions: 29 h x 43 w x 21 d </span><br data-mce-fragment=\"1\"><span data-mce-fragment=\"1\"> Colour: White / Red</span><br data-mce-fragment=\"1\"><span data-mce-fragment=\"1\"> Material: Wood</span><br data-mce-fragment=\"1\"><span data-mce-fragment=\"1\"> Era/ Style: 19th / 20th Century<br>Matches: <a href=\"https://mountpf.com/collections/dresser/products/drsr005\" target=\"_blank\" title=\"DRSR005\" rel=\"noopener noreferrer\">Dresser DRSR005</a></span>")

	wantSize := "29 H X 43 W X 21 D"
	if got.Size != wantSize {
		t.Errorf("got %q, wanted %q", got.Size, wantSize)
	}

	wantStyle := "19TH / 20TH CENTURY"
	if got.Style != wantStyle {
		t.Errorf("got %q, wanted %q", got.Style, wantStyle)
	}

	wantMaterial := "WOOD"
	if got.Material != wantMaterial {
		t.Errorf("got %q, wanted %q", got.Material, wantMaterial)
	}

	wantColor := "WHITE / RED"
	if got.Color != wantColor {
		t.Errorf("got %q, wanted %q", got.Color, wantColor)
	}
}

func TestExtractDescriptions_05(t *testing.T) {
	got := extractDescriptions("<span data-mce-fragment=\"1\">Dimensions: 27.5 w x 40.5 h x 22 d </span><br data-mce-fragment=\"1\"><span data-mce-fragment=\"1\">Colour: Dark </span><br data-mce-fragment=\"1\"><span data-mce-fragment=\"1\">Material: Wood, Cane, </span><br data-mce-fragment=\"1\"><span data-mce-fragment=\"1\">Era/ Style: 20th Century</span>")

	wantSize := "27.5 W X 40.5 H X 22 D"
	if got.Size != wantSize {
		t.Errorf("got %q, wanted %q", got.Size, wantSize)
	}

	wantStyle := "20TH CENTURY"
	if got.Style != wantStyle {
		t.Errorf("got %q, wanted %q", got.Style, wantStyle)
	}

	wantMaterial := "WOOD, CANE,"
	if got.Material != wantMaterial {
		t.Errorf("got %q, wanted %q", got.Material, wantMaterial)
	}

	wantColor := "DARK"
	if got.Color != wantColor {
		t.Errorf("got %q, wanted %q", got.Color, wantColor)
	}
}

func TestExtractInventoryAmount(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Returning x2",
			input:    "Library Table T074 (x2)",
			expected: "X2",
		},
		{
			name:     "Returning x3",
			input:    "Library Table T074 (x3)",
			expected: "X3",
		},
		{
			name:     "Empty Material",
			input:    "Library Table T074",
			expected: "X1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractInventoryAmount(tt.input)
			if result != tt.expected {
				t.Errorf("extractInventoryAmount(%q) = %q; want %q", tt.input, result, tt.expected)
			}
		})
	}

}

func TestRemoveIntFromTags(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name:     "Cleaning out 3 numbers",
			input:    []string{"001", "01", "1", "carpet", "largecarpet"},
			expected: []string{"CARPET", "LARGECARPET"},
		},
		{
			name:     "Correct values if numbers are in the tag",
			input:    []string{"19th century", "EXECUTIVE DESK"},
			expected: []string{"19TH CENTURY", "EXECUTIVE DESK"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := removeIntFromTags(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("removeIntFromTags(%v) = %v; want %v", tt.input, result, tt.expected)
			}
		})
	}
}
