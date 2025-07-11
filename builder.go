package htmlgopdf

import "time"

// WithOptions creates a generator with custom options - builder pattern
func WithOptions() *OptionsBuilder {
	return &OptionsBuilder{options: DefaultOptions()}
}

// OptionsBuilder provides a fluent interface for building PDF options
type OptionsBuilder struct {
	options *PDFOptions
}

// Format sets the paper format (A4, A3, Letter, Legal, Tabloid)
func (b *OptionsBuilder) Format(format string) *OptionsBuilder {
	b.options.Format = format
	return b
}

// Size sets custom paper size in inches
func (b *OptionsBuilder) Size(width, height float64) *OptionsBuilder {
	b.options.Width = width
	b.options.Height = height
	b.options.Format = "" // Clear format when using custom size
	return b
}

// Margins sets all margins in inches
func (b *OptionsBuilder) Margins(top, bottom, left, right float64) *OptionsBuilder {
	b.options.MarginTop = top
	b.options.MarginBottom = bottom
	b.options.MarginLeft = left
	b.options.MarginRight = right
	return b
}

// Landscape sets the orientation to landscape
func (b *OptionsBuilder) Landscape() *OptionsBuilder {
	b.options.Landscape = true
	return b
}

// Portrait sets the orientation to portrait
func (b *OptionsBuilder) Portrait() *OptionsBuilder {
	b.options.Landscape = false
	return b
}

// Scale sets the scale factor (0.1 to 2.0)
func (b *OptionsBuilder) Scale(scale float64) *OptionsBuilder {
	b.options.Scale = scale
	return b
}

// PrintBackground enables/disables background printing
func (b *OptionsBuilder) PrintBackground(enable bool) *OptionsBuilder {
	b.options.PrintBackground = enable
	return b
}

// HeaderFooter enables header and footer with templates
func (b *OptionsBuilder) HeaderFooter(header, footer string) *OptionsBuilder {
	b.options.DisplayHeaderFooter = true
	b.options.HeaderTemplate = header
	b.options.FooterTemplate = footer
	return b
}

// WaitFor sets a CSS selector to wait for before generating PDF
func (b *OptionsBuilder) WaitFor(selector string) *OptionsBuilder {
	b.options.WaitForSelector = selector
	return b
}

// WaitTime sets additional wait time before generating PDF
func (b *OptionsBuilder) WaitTime(duration time.Duration) *OptionsBuilder {
	b.options.WaitTime = duration
	return b
}

// Timeout sets the context timeout for PDF generation
func (b *OptionsBuilder) Timeout(duration time.Duration) *OptionsBuilder {
	b.options.Timeout = duration
	return b
}

// Build creates the PDF generator with the configured options
func (b *OptionsBuilder) Build() *Generator {
	return NewGenerator(b.options)
}

// Generate generates PDF from HTML using the configured options
func (b *OptionsBuilder) Generate(htmlContent string) ([]byte, error) {
	return b.Build().FromHTML(htmlContent)
}

// GenerateFromURL generates PDF from URL using the configured options
func (b *OptionsBuilder) GenerateFromURL(url string) ([]byte, error) {
	return b.Build().FromURL(url)
}
