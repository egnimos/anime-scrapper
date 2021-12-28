package servers

import (
	"context"

	"github.com/chromedp/chromedp"
	"github.com/egnimos/anime-scrapper/src/utility"
)

type task func() chromedp.Tasks

func InitializeChromeDp(submit task) utility.RestError {
	opts := append(chromedp.DefaultExecAllocatorOptions[:], chromedp.Flag("headless", false))
	ctx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	//new context
	nctx, cancel := chromedp.NewContext(ctx)
	defer cancel()

	//run the task
	if err := chromedp.Run(nctx, submit()); err != nil {
		return utility.NewInternalServerError(err.Error())
	}

	return nil
}
