package htmlgopdf

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

// Generator handles PDF generation from HTML content
type Generator struct {
	options *PDFOptions
}

// NewGenerator creates a new PDF generator with the given options
func NewGenerator(options *PDFOptions) *Generator {
	if options == nil {
		options = DefaultOptions()
	}
	return &Generator{
		options: options,
	}
}

// FromHTML generates a PDF from HTML content string
func (g *Generator) FromHTML(htmlContent string) ([]byte, error) {
	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), g.options.Timeout)
	defer cancel()

	// Create a new browser context
	ctx, cancel = chromedp.NewContext(ctx)
	defer cancel()

	var pdfData []byte
	var err error

	dataURL := "data:text/html;charset=utf-8," + url.PathEscape(htmlContent)

	// Execute the browser automation
	err = chromedp.Run(ctx,
		chromedp.Navigate(dataURL),
		chromedp.WaitReady("body"),
		g.waitForConditions(),
		chromedp.ActionFunc(func(ctx context.Context) error {
			pdfData, err = g.generatePDF(ctx)
			return err
		}),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to generate PDF: %w", err)
	}

	return pdfData, nil
}

// FromURL generates a PDF from a URL
func (g *Generator) FromURL(url string) ([]byte, error) {
	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), g.options.Timeout)
	defer cancel()

	// Create a new browser context
	ctx, cancel = chromedp.NewContext(ctx)
	defer cancel()

	var pdfData []byte
	var err error

	// Execute the browser automation
	err = chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.WaitReady("body"),
		g.waitForConditions(),
		chromedp.ActionFunc(func(ctx context.Context) error {
			pdfData, err = g.generatePDF(ctx)
			return err
		}),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to generate PDF from URL: %w", err)
	}

	return pdfData, nil
}

// waitForConditions handles waiting for specific conditions before PDF generation
func (g *Generator) waitForConditions() chromedp.Action {
	var actions []chromedp.Action

	// Wait for specific selector if provided
	if g.options.WaitForSelector != "" {
		actions = append(actions, chromedp.WaitVisible(g.options.WaitForSelector))
	}

	// Additional wait time
	if g.options.WaitTime > 0 {
		actions = append(actions, chromedp.Sleep(g.options.WaitTime))
	}

	if len(actions) == 0 {
		return chromedp.Sleep(500 * time.Millisecond) // Default small wait
	}

	return chromedp.Tasks(actions)
}

// generatePDF generates the actual PDF using Chrome DevTools Protocol
func (g *Generator) generatePDF(ctx context.Context) ([]byte, error) {
	// Build PDF parameters using the correct chromedp API
	params := page.PrintToPDFParams{
		PrintBackground:     g.options.PrintBackground,
		Landscape:           g.options.Landscape,
		DisplayHeaderFooter: g.options.DisplayHeaderFooter,
		Scale:               g.options.Scale,
	}

	// Set paper size based on format or custom dimensions
	if g.options.Format != "" {
		switch g.options.Format {
		case "A4":
			params.PaperWidth = 8.27  // A4 width in inches
			params.PaperHeight = 11.7 // A4 height in inches
		case "A3":
			params.PaperWidth = 11.7
			params.PaperHeight = 16.5
		case "Letter":
			params.PaperWidth = 8.5
			params.PaperHeight = 11.0
		case "Legal":
			params.PaperWidth = 8.5
			params.PaperHeight = 14.0
		case "Tabloid":
			params.PaperWidth = 11.0
			params.PaperHeight = 17.0
		}
	} else if g.options.Width > 0 && g.options.Height > 0 {
		params.PaperWidth = g.options.Width
		params.PaperHeight = g.options.Height
	}

	// Set margins
	params.MarginTop = g.options.MarginTop
	params.MarginBottom = g.options.MarginBottom
	params.MarginLeft = g.options.MarginLeft
	params.MarginRight = g.options.MarginRight

	// Set header and footer templates
	if g.options.HeaderTemplate != "" {
		params.HeaderTemplate = g.options.HeaderTemplate
	}
	if g.options.FooterTemplate != "" {
		params.FooterTemplate = g.options.FooterTemplate
	}

	// Generate PDF using the correct chromedp method
	var pdfData []byte
	err := chromedp.Run(ctx, chromedp.ActionFunc(func(ctx context.Context) error {
		var err error
		pdfData, _, err = params.Do(ctx)
		return err
	}))

	if err != nil {
		return nil, fmt.Errorf("failed to generate PDF: %w", err)
	}

	return pdfData, nil
}

// Convenience functions for common use cases

// FromHTML is a convenience function for basic HTML to PDF conversion
func FromHTML(htmlContent string) ([]byte, error) {
	generator := NewGenerator(DefaultOptions())
	return generator.FromHTML(htmlContent)
}

// FromURL is a convenience function for basic URL to PDF conversion
func FromURL(url string) ([]byte, error) {
	generator := NewGenerator(DefaultOptions())
	return generator.FromURL(url)
}
