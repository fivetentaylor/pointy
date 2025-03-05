package main

import (
	"context"
	"io/ioutil"
	"log"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

func main() {
	// Start a new Chromedp context
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// The HTML content is converted to a string and URL encoded to be used in a data URL
	// htmlDataUrl := "data:text/html," + "<html><body><h1>Hello, World!</h1></body></html>"

	// Variable to store the screenshot
	var screenshot []byte

	// Set the desired viewport size here
	width, height := 800, 1600 // Example: 800x600 pixels

	headers := map[string]interface{}{
		"X-Reviso": "1",
	}

	// Run tasks
	err := chromedp.Run(ctx,
		chromedp.EmulateViewport(int64(width), int64(height)),
		network.Enable(),
		network.SetExtraHTTPHeaders(network.Headers(headers)),
		chromedp.Navigate("https://www.reviso.dev/screenshot/af4e3548-e6de-4f31-8283-03074f71d0a1"),
		chromedp.WaitVisible("#screenshot-page", chromedp.ByQuery), // Ensure the page is loaded
		chromedp.Screenshot("#screenshot-page", &screenshot, chromedp.NodeVisible, chromedp.ByQuery),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Write the screenshot to a file
	if err := ioutil.WriteFile("screenshot.png", screenshot, 0644); err != nil {
		log.Fatal(err)
	}

	log.Println("Screenshot taken and saved as screenshot.png")
}
