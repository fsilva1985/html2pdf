package main

import (
	"context"
	"io/ioutil"
	"log"

	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"github.com/thatisuday/commando"
)

func main() {
	commando.
		SetExecutableName("html2pdf").
		SetVersion("alpha").
		SetDescription("Google Chrome Command Line Tool to convert html to pdf | html2pdf")

	commando.
		Register(nil).
		AddArgument("source", "Url or local file.", "").
		AddArgument("output", "File output", "").
		SetAction(func(args map[string]commando.ArgValue, flags map[string]commando.FlagValue) {
			generate(args["source"].Value, args["output"].Value)
		})

	commando.Parse(nil)
}

func generate(source string, output string) {
	taskCtx, cancel := chromedp.NewContext(
		context.Background(),
		chromedp.WithLogf(log.Printf),
	)
	defer cancel()
	var pdfBuffer []byte
	if err := chromedp.Run(taskCtx, getContent(source, "html", &pdfBuffer)); err != nil {
		log.Fatal(err)
	}
	if err := ioutil.WriteFile(output, pdfBuffer, 0644); err != nil {
		log.Fatal(err)
	}
}

func getContent(source string, sel string, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		emulation.SetUserAgentOverride("WebScraper 1.0"),
		chromedp.Navigate(source),
		chromedp.WaitVisible(sel, chromedp.ByQuery),
		chromedp.ActionFunc(func(ctx context.Context) error {
			buf, _, err := page.PrintToPDF().WithPrintBackground(true).Do(ctx)
			if err != nil {
				return err
			}
			*res = buf
			return nil
		}),
	}
}
