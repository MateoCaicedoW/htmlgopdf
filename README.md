# htmlgopdf

A Go package for generating PDF documents from HTML content using Chrome/Chromium via the Chrome DevTools Protocol.

## Features

- Generate PDFs from HTML strings or URLs
- Customizable page formats (A4, A3, Letter, Legal, Tabloid)
- Custom paper sizes and margins
- Portrait and landscape orientations
- Background graphics support
- Header and footer templates
- Wait conditions for dynamic content
- Fluent builder pattern for easy configuration
- Context-based timeout handling

## Installation

```bash
go get github.com/MateoCaicedoW/htmlgopdf
```

## Requirements

This package requires Chrome or Chromium to be installed on your system as it uses the Chrome DevTools Protocol via [chromedp](https://github.com/chromedp/chromedp).

### Docker Setup

When using this package in a Docker container, you'll need to install Chromium. Here are examples for different base images:

#### Alpine Linux

```dockerfile
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git

# Copy and build your application
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o app .

# Runtime stage
FROM alpine:latest

# Install Chromium and required dependencies
RUN apk add --no-cache \
    chromium \
    nss \
    freetype \
    harfbuzz \
    ca-certificates \
    ttf-freefont \
    && rm -rf /var/cache/apk/*

# Tell Puppeteer to skip installing Chrome. We'll be using the installed package.
ENV PUPPETEER_SKIP_CHROMIUM_DOWNLOAD=true \
    PUPPETEER_EXECUTABLE_PATH=/usr/bin/chromium-browser

# Copy the binary from builder stage
COPY --from=builder /app/app /usr/local/bin/app

# Create a non-root user
RUN adduser -D -s /bin/sh appuser
USER appuser

EXPOSE 8080
CMD ["app"]
```

#### Ubuntu/Debian

```dockerfile
FROM golang:1.21 AS builder

WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o app .

# Runtime stage
FROM ubuntu:22.04

# Install Chromium and dependencies
RUN apt-get update && apt-get install -y \
    chromium-browser \
    ca-certificates \
    fonts-liberation \
    libappindicator3-1 \
    libasound2 \
    libatk-bridge2.0-0 \
    libdrm2 \
    libgtk-3-0 \
    libnspr4 \
    libnss3 \
    libxcomposite1 \
    libxdamage1 \
    libxrandr2 \
    xdg-utils \
    && rm -rf /var/lib/apt/lists/*

# Copy the binary from builder stage
COPY --from=builder /app/app /usr/local/bin/app

# Create a non-root user
RUN useradd -m -s /bin/bash appuser
USER appuser

EXPOSE 8080
CMD ["app"]
```

#### Distroless (Recommended for production)

```dockerfile
FROM golang:1.21 AS builder

WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o app .

# Get Chromium
FROM ubuntu:22.04 AS chromium
RUN apt-get update && apt-get install -y \
    chromium-browser \
    && rm -rf /var/lib/apt/lists/*

# Runtime stage
FROM gcr.io/distroless/base-debian11

# Copy Chromium from the chromium stage
COPY --from=chromium /usr/bin/chromium-browser /usr/bin/chromium-browser
COPY --from=chromium /usr/lib/chromium-browser /usr/lib/chromium-browser
COPY --from=chromium /etc/chromium-browser /etc/chromium-browser

# Copy the binary from builder stage
COPY --from=builder /app/app /app

EXPOSE 8080
ENTRYPOINT ["/app"]
```

#### Docker Compose Example

```yaml
version: '3.8'
services:
  pdf-service:
    build: .
    ports:
      - "8080:8080"
    environment:
      - CHROMIUM_PATH=/usr/bin/chromium-browser
    # Add shared memory size for Chrome
    shm_size: '2gb'
    # Or use tmpfs
    tmpfs:
      - /tmp:size=2G
```

#### Important Docker Notes

1. **Shared Memory**: Chrome needs adequate shared memory. Use `shm_size: '2gb'` or mount tmpfs
2. **Security**: Run as non-root user when possible
3. **Fonts**: Install font packages for better text rendering
4. **Environment Variables**: You may need to set Chrome flags for containerized environments

#### Chrome Flags for Docker

You might need to configure Chrome with specific flags for containerized environments:

```go
// In your Go code, you can set Chrome flags via chromedp
opts := append(chromedp.DefaultExecAllocatorOptions[:],
    chromedp.Flag("no-sandbox", true),
    chromedp.Flag("disable-gpu", true),
    chromedp.Flag("disable-dev-shm-usage", true),
    chromedp.Flag("disable-extensions", true),
    chromedp.Flag("no-first-run", true),
)

ctx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
defer cancel()

ctx, cancel = chromedp.NewContext(ctx)
defer cancel()

// Use this context for PDF generation
```

## Quick Start

### Basic Usage

```go
package main

import (
    "fmt"
    "os"
    "github.com/MateoCaicedoW/htmlgopdf"
)

func main() {
    // Simple HTML to PDF conversion
    html := `
    <html>
        <body>
            <h1>Hello World</h1>
            <p>This is a PDF generated from HTML!</p>
        </body>
    </html>`
    
    pdfData, err := htmlgopdf.FromHTML(html)
    if err != nil {
        panic(err)
    }
    
    // Save to file
    err = os.WriteFile("output.pdf", pdfData, 0644)
    if err != nil {
        panic(err)
    }
    
    fmt.Println("PDF generated successfully!")
}
```

### Generate PDF from URL

```go
pdfData, err := htmlgopdf.FromURL("https://example.com")
if err != nil {
    panic(err)
}
```

## Advanced Usage

### Using the Builder Pattern

```go
package main

import (
    "time"
    "github.com/MateoCaicedoW/htmlgopdf"
)

func main() {
    html := `<html><body><h1>Custom PDF</h1></body></html>`
    
    pdfData, err := htmlgopdf.WithOptions().
        Format("A3").
        Landscape().
        Margins(0.5, 0.5, 0.5, 0.5).
        Scale(0.8).
        PrintBackground(true).
        WaitTime(time.Second * 3).
        Generate(html)
    
    if err != nil {
        panic(err)
    }
    
    // Save or use pdfData...
}
```

### Custom Generator with Options

```go
package main

import (
    "time"
    "github.com/MateoCaicedoW/htmlgopdf"
)

func main() {
    // Create custom options
    options := &htmlgopdf.PDFOptions{
        Format:          htmlgopdf.FormatLetter,
        MarginTop:       0.75,
        MarginBottom:    0.75,
        MarginLeft:      0.5,
        MarginRight:     0.5,
        Landscape:       false,
        PrintBackground: true,
        Scale:           1.2,
        WaitTime:        time.Second * 5,
        Timeout:         time.Minute * 2,
    }
    
    generator := htmlgopdf.NewGenerator(options)
    
    html := `<html><body><h1>Custom Configuration</h1></body></html>`
    pdfData, err := generator.FromHTML(html)
    if err != nil {
        panic(err)
    }
    
    // Use pdfData...
}
```

### Headers and Footers

```go
headerHTML := `
<div style="font-size: 10px; text-align: center; width: 100%;">
    <span>Document Title</span>
</div>`

footerHTML := `
<div style="font-size: 10px; text-align: center; width: 100%;">
    <span>Page <span class="pageNumber"></span> of <span class="totalPages"></span></span>
</div>`

pdfData, err := htmlgopdf.WithOptions().
    Format("A4").
    HeaderFooter(headerHTML, footerHTML).
    Generate(html)
```

### Wait for Dynamic Content

```go
// Wait for a specific element to be visible
pdfData, err := htmlgopdf.WithOptions().
    WaitFor("#content-loaded").
    WaitTime(time.Second * 2).
    Generate(html)

// Or wait for a specific amount of time
pdfData, err := htmlgopdf.WithOptions().
    WaitTime(time.Second * 5).
    Generate(html)
```

## Configuration Options

### PDFOptions

| Field | Type | Description | Default |
|-------|------|-------------|---------|
| `Format` | `string` | Paper format (A4, A3, Letter, Legal, Tabloid) | `"A4"` |
| `Width` | `float64` | Custom paper width in inches | `0` |
| `Height` | `float64` | Custom paper height in inches | `0` |
| `MarginTop` | `float64` | Top margin in inches | `0.4` |
| `MarginBottom` | `float64` | Bottom margin in inches | `0.4` |
| `MarginLeft` | `float64` | Left margin in inches | `0.4` |
| `MarginRight` | `float64` | Right margin in inches | `0.4` |
| `Landscape` | `bool` | Landscape orientation | `false` |
| `PrintBackground` | `bool` | Include background graphics | `true` |
| `Scale` | `float64` | Scale factor (0.1 to 2.0) | `1.0` |
| `DisplayHeaderFooter` | `bool` | Display header and footer | `false` |
| `HeaderTemplate` | `string` | HTML template for header | `""` |
| `FooterTemplate` | `string` | HTML template for footer | `""` |
| `WaitForSelector` | `string` | CSS selector to wait for | `""` |
| `WaitTime` | `time.Duration` | Additional wait time | `2s` |
| `Timeout` | `time.Duration` | Context timeout | `30s` |

### Builder Methods

| Method | Description |
|--------|-------------|
| `Format(string)` | Set paper format |
| `Size(width, height float64)` | Set custom paper size |
| `Margins(top, bottom, left, right float64)` | Set all margins |
| `Landscape()` | Set landscape orientation |
| `Portrait()` | Set portrait orientation |
| `Scale(float64)` | Set scale factor |
| `PrintBackground(bool)` | Enable/disable background printing |
| `HeaderFooter(header, footer string)` | Set header and footer templates |
| `WaitFor(selector string)` | Wait for CSS selector |
| `WaitTime(duration)` | Set additional wait time |
| `Timeout(duration)` | Set context timeout |

## Paper Formats

The following predefined formats are available:

- `htmlgopdf.FormatA4` - A4 (8.27" × 11.7")
- `htmlgopdf.FormatA3` - A3 (11.7" × 16.5")
- `htmlgopdf.FormatLetter` - Letter (8.5" × 11")
- `htmlgopdf.FormatLegal` - Legal (8.5" × 14")
- `htmlgopdf.FormatTabloid` - Tabloid (11" × 17")

## Error Handling

The package returns detailed error messages for common issues:

```go
pdfData, err := htmlgopdf.FromHTML(html)
if err != nil {
    // Handle specific error cases
    switch {
    case strings.Contains(err.Error(), "timeout"):
        // Handle timeout
    case strings.Contains(err.Error(), "navigate"):
        // Handle navigation errors
    default:
        // Handle other errors
    }
}
```

## Best Practices

1. **Set appropriate timeouts** - Complex pages may need longer timeouts
2. **Use wait conditions** - For dynamic content, wait for specific elements or use additional wait time
3. **Optimize HTML** - Use print-friendly CSS and avoid complex animations
4. **Handle errors gracefully** - Always check for errors and handle them appropriately
5. **Resource management** - The package automatically handles Chrome instance lifecycle

## CSS Print Styles

For better PDF output, consider using CSS print styles:

```html
<style>
@media print {
    body {
        font-family: Arial, sans-serif;
        font-size: 12pt;
        line-height: 1.4;
    }
    
    .no-print {
        display: none;
    }
    
    .page-break {
        page-break-before: always;
    }
}
</style>
```

## Dependencies

- [chromedp](https://github.com/chromedp/chromedp) - Chrome DevTools Protocol client

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Support

If you encounter any issues or have questions, please open an issue on GitHub.