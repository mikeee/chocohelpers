package chocohelpers

import (
	"context"
	"regexp"

	"github.com/chromedp/cdproto/dom"
	"github.com/chromedp/chromedp"
)

type Options struct {
	UserAgent string
}

type optionFunction func(*Options)

func WithUserAgent(ua string) optionFunction {
	return func(o *Options) {
		o.UserAgent = ua
	}
}

func RegexMatch(url string, matchString string, matchIndex int) (rawRegex string, err error) {
	return RegexMatchWithOpts(url, matchString, matchIndex)
}

func RegexMatchWithOpts(url string, matchString string, matchIndex int, opts ...optionFunction) (rawRegex string,
	err error,
) {
	options := Options{}
	for _, opt := range opts {
		opt(&options)
	}

	chromeDPOptions := chromedp.DefaultExecAllocatorOptions[:]

	if options.UserAgent != "" {
		chromeDPOptions = append(chromeDPOptions, chromedp.UserAgent(options.UserAgent))
	}

	// specify a new context set up for NewContext
	ctx, cancel := chromedp.NewExecAllocator(context.Background(), chromeDPOptions...)
	defer cancel()

	ctx, cancel = chromedp.NewContext(context.Background())
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
