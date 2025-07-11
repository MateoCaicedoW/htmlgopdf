package main

import (
	"fmt"
	"log"
	"os"

	"github.com/MateoCaicedoW/htmlgopdf"
)

func main() {
	simpleExample()
}

func simpleExample() {
	fmt.Println("=== Simple HTML to PDF ===")

	html := `
	<!DOCTYPE html>
	<html>
	<head>
		<meta charset="UTF-8">
		<style>
			body { 
				font-family: Arial, sans-serif; 
				margin: 20px; 
				background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
				color: white;
				min-height: 50vh;
			}
			.container {
				background: rgba(255,255,255,0.1);
				padding: 30px;
				border-radius: 10px;
				backdrop-filter: blur(10px);
			}
			h1 { color: #fff; text-align: center; }
		</style>
	</head>
	<body>
		<div class="container">
			<h1>Hello ChromePDF!</h1>
			<p>This is a simple PDF generated from HTML with background colors.</p>
		</div>
	</body>
	</html>
	`

	pdfData, err := htmlgopdf.FromHTML(html)
	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile("simple.pdf", pdfData, 0644)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("✅ Simple PDF generated: %d bytes\n", len(pdfData))

	pdfData, err = htmlgopdf.WithOptions().
		Landscape().
		Size(8.5, 11).
		Margins(0.5, 0.5, 0.5, 0.5).
		PrintBackground(true).
		Scale(1.0).
		Format(htmlgopdf.FormatTabloid).Build().FromHTML(html)
	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile("simple2.pdf", pdfData, 0644)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("✅ Simple PDF generated: %d bytes\n", len(pdfData))
}
