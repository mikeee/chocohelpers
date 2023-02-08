package chocohelpers

import (
	"context"
	"regexp"

	"github.com/chromedp/cdproto/dom"
	"github.com/chromedp/chromedp"
)

func RegexMatch(url string, matchString string, matchIndex int) (rawRegex string, err error) {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	var res string

	err = chromedp.Run(
		ctx,
		chromedp.Navigate(url),
		chromedp.ActionFunc(
			func(ctx context.Context) error {
				node, err := dom.GetDocument().Do(ctx)
				if err == nil {
					res, err = dom.GetOuterHTML().WithNodeID(node.NodeID).Do(ctx)
				}
				return err
			}),
	)

	if err == nil {
		re := regexp.MustCompile(matchString)

		rawRegex = re.FindStringSubmatch(res)[matchIndex]

		return
	}

	return "", err
}
