package htmlgopdf

import (
	"time"
)

// PDFOptions represents configuration options for PDF generation
type PDFOptions struct {
	// Page settings
	Format string  `json:"format,omitempty"` // A4, A3, Letter, etc.
	Width  float64 `json:"width,omitempty"`  // Paper width in inches
	Height float64 `json:"height,omitempty"` // Paper height in inches

	// Margins in inches
	MarginTop    float64 `json:"marginTop,omitempty"`    // Top margin
	MarginBottom float64 `json:"marginBottom,omitempty"` // Bottom margin
	MarginLeft   float64 `json:"marginLeft,omitempty"`   // Left margin
	MarginRight  float64 `json:"marginRight,omitempty"`  // Right margin

	// Layout options
	Landscape       bool `json:"landscape,omitempty"`       // Landscape orientation
	PrintBackground bool `json:"printBackground,omitempty"` // Include background graphics

	// Scale and quality
	Scale float64 `json:"scale,omitempty"` // Scale of the webpage rendering (0.1 to 2)

	// Header and footer
	DisplayHeaderFooter bool   `json:"displayHeaderFooter,omitempty"` // Display header and footer
	HeaderTemplate      string `json:"headerTemplate,omitempty"`      // HTML template for header
	FooterTemplate      string `json:"footerTemplate,omitempty"`      // HTML template for footer

	// Wait conditions
	WaitForSelector string        `json:"-"` // CSS selector to wait for before generating PDF
	WaitTime        time.Duration `json:"-"` // Additional wait time

	// Timeout
	Timeout time.Duration `json:"-"` // Context timeout
}

// DefaultOptions returns sensible defaults for PDF generation
func DefaultOptions() *PDFOptions {
	return &PDFOptions{
		Format:          "A4",
		MarginTop:       0.4,
		MarginBottom:    0.4,
		MarginLeft:      0.4,
		MarginRight:     0.4,
		Landscape:       false,
		PrintBackground: true,
		Scale:           1.0,
		WaitTime:        time.Second * 2,
		Timeout:         time.Second * 30,
	}
}
