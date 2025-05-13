package main

import (
	"flag"
	"fmt"
	"image"
	"image/png"
	"log"
	"os"
	"path/filepath"

	ico "github.com/Kodeworks/golang-image-ico"
	"github.com/disintegration/imaging"
)

var sizes = []int{16, 32, 48, 57, 60, 72, 76, 96, 114, 120, 144, 152, 180, 192, 512}

func main() {
	inputFile := flag.String("input", "", "Input image file (required)")
	outputDir := flag.String("output", "favicons", "Output directory for favicon files")
	flag.Parse()

	if *inputFile == "" {
		fmt.Println("Error: Input file is required")
		fmt.Println("Usage: favicon-go -input <image-file> [-output <output-directory>]")
		os.Exit(1)
	}

	// Create output directory if it doesn't exist
	if err := os.MkdirAll(*outputDir, 0755); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	// Open source image
	src, err := imaging.Open(*inputFile)
	if err != nil {
		log.Fatalf("Failed to open input image: %v", err)
	}

	// Generate favicon.ico (combining 16, 32, 48 px)
	generateFaviconIco(src, *outputDir)

	// Generate various size PNG files
	generatePngIcons(src, *outputDir)

	// Generate manifest.json and browserconfig.xml
	generateManifest(*outputDir)
	generateBrowserConfig(*outputDir)

	fmt.Println("Favicon generation complete!")
	fmt.Println("Add the following to your HTML <head> section:")
	printHtmlCode(*outputDir)
}

func generateFaviconIco(src image.Image, outputDir string) {
	// Create resized image for ICO (using 32x32 as the standard size)
	resized := imaging.Resize(src, 32, 32, imaging.Lanczos)

	// Create favicon.ico
	outputPath := filepath.Join(outputDir, "favicon.ico")
	f, err := os.Create(outputPath)
	if err != nil {
		log.Fatalf("Failed to create favicon.ico: %v", err)
	}
	defer f.Close()

	// The ico.Encode function expects a single image, not a slice of images
	if err := ico.Encode(f, resized); err != nil {
		log.Fatalf("Failed to encode favicon.ico: %v", err)
	}

	fmt.Printf("Created: %s\n", outputPath)

	// Additionally create favicon-16x16.png and favicon-32x32.png for HTML references
	for _, size := range []int{16, 32} {
		resized := imaging.Resize(src, size, size, imaging.Lanczos)
		outputPath := filepath.Join(outputDir, fmt.Sprintf("favicon-%dx%d.png", size, size))
		f, err := os.Create(outputPath)
		if err != nil {
			log.Fatalf("Failed to create %s: %v", outputPath, err)
		}

		if err := png.Encode(f, resized); err != nil {
			f.Close()
			log.Fatalf("Failed to encode %s: %v", outputPath, err)
		}

		f.Close()
		fmt.Printf("Created: %s\n", outputPath)
	}
}

func generatePngIcons(src image.Image, outputDir string) {
	for _, size := range sizes {
		resized := imaging.Resize(src, size, size, imaging.Lanczos)

		outputPath := filepath.Join(outputDir, fmt.Sprintf("favicon-%dx%d.png", size, size))
		f, err := os.Create(outputPath)
		if err != nil {
			log.Fatalf("Failed to create %s: %v", outputPath, err)
		}

		if err := png.Encode(f, resized); err != nil {
			f.Close()
			log.Fatalf("Failed to encode %s: %v", outputPath, err)
		}

		f.Close()
		fmt.Printf("Created: %s\n", outputPath)

		// Create special names for certain sizes
		if size == 180 {
			copyFile(outputPath, filepath.Join(outputDir, "apple-touch-icon.png"))
		} else if size == 192 {
			copyFile(outputPath, filepath.Join(outputDir, "android-chrome-192x192.png"))
		} else if size == 512 {
			copyFile(outputPath, filepath.Join(outputDir, "android-chrome-512x512.png"))
		}
	}
}

func copyFile(srcPath, dstPath string) {
	input, err := os.ReadFile(srcPath)
	if err != nil {
		log.Fatalf("Failed to read file %s: %v", srcPath, err)
	}

	if err := os.WriteFile(dstPath, input, 0644); err != nil {
		log.Fatalf("Failed to write file %s: %v", dstPath, err)
	}

	fmt.Printf("Created: %s\n", dstPath)
}

func generateManifest(outputDir string) {
	manifest := `{
  "name": "Your Website",
  "short_name": "Website",
  "icons": [
    {
      "src": "/android-chrome-192x192.png",
      "sizes": "192x192",
      "type": "image/png"
    },
    {
      "src": "/android-chrome-512x512.png",
      "sizes": "512x512",
      "type": "image/png"
    }
  ],
  "theme_color": "#ffffff",
  "background_color": "#ffffff",
  "display": "standalone"
}`

	writeConfigFile(outputDir, "site.webmanifest", manifest)
}

func generateBrowserConfig(outputDir string) {
	config := `<?xml version="1.0" encoding="utf-8"?>
<browserconfig>
    <msapplication>
        <tile>
            <square150x150logo src="/favicon-144x144.png"/>
            <TileColor>#ffffff</TileColor>
        </tile>
    </msapplication>
</browserconfig>`

	writeConfigFile(outputDir, "browserconfig.xml", config)
}

func writeConfigFile(outputDir, filename, content string) {
	outputPath := filepath.Join(outputDir, filename)
	if err := os.WriteFile(outputPath, []byte(content), 0644); err != nil {
		log.Fatalf("Failed to write %s file: %v", filename, err)
	}

	fmt.Printf("Created: %s\n", outputPath)
}

func printHtmlCode(outputDir string) {
	relativePath := filepath.Base(outputDir)

	elements := []struct {
		tag       string
		attrs     map[string]string
		selfClose bool
	}{
		{
			tag:       "link",
			attrs:     map[string]string{"rel": "apple-touch-icon", "sizes": "180x180", "href": fmt.Sprintf("/%s/apple-touch-icon.png", relativePath)},
			selfClose: true,
		},
		{
			tag:       "link",
			attrs:     map[string]string{"rel": "icon", "type": "image/png", "sizes": "32x32", "href": fmt.Sprintf("/%s/favicon-32x32.png", relativePath)},
			selfClose: true,
		},
		{
			tag:       "link",
			attrs:     map[string]string{"rel": "icon", "type": "image/png", "sizes": "16x16", "href": fmt.Sprintf("/%s/favicon-16x16.png", relativePath)},
			selfClose: true,
		},
		{
			tag:       "link",
			attrs:     map[string]string{"rel": "manifest", "href": fmt.Sprintf("/%s/site.webmanifest", relativePath)},
			selfClose: true,
		},
		{
			tag:       "meta",
			attrs:     map[string]string{"name": "msapplication-config", "content": fmt.Sprintf("/%s/browserconfig.xml", relativePath)},
			selfClose: true,
		},
		{
			tag:       "meta",
			attrs:     map[string]string{"name": "theme-color", "content": "#ffffff"},
			selfClose: true,
		},
	}

	for _, elem := range elements {
		html := "<" + elem.tag
		for key, value := range elem.attrs {
			html += fmt.Sprintf(` %s="%s"`, key, value)
		}
		if elem.selfClose {
			html += ">"
		}
		fmt.Println(html)
	}
}
