# Favicon Generator

A command-line tool to generate all necessary favicon files for a website from a single input image.

## Installation

First, make sure you have Go installed on your system. Then install the required dependencies:

```bash
go mod download
```

Then build the executable:

```bash
go build -o favicon-go
```

## Usage

```bash
./favicon-go -input your-image.png [-output favicons]
```

## Alternative install
```bash
go install https://github.com/naruebaet/favicon-go@latest
```

## Usage by global install
```bash
favicon-go -input your-image.png [-output favicons]
```

### Parameters

- `-input`: Required. Path to the input image (PNG, JPEG, etc.).
- `-output`: Optional. Directory where favicon files will be saved. Defaults to "favicons".

## Generated Files

The tool generates:

- favicon.ico (16x16, 32x32, 48x48)
- PNG files in various sizes (16x16 to 512x512)
- apple-touch-icon.png
- android-chrome-192x192.png and android-chrome-512x512.png
- site.webmanifest
- browserconfig.xml

## HTML Code

After generation, the tool will output HTML code that you should add to your website's `<head>` section to properly reference the favicon files.

## Requirements

- Go 1.16 or later
- The input image should ideally be square and at least 512x512 pixels in size for best results.
