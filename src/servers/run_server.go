package servers

import (
	"context"

	"github.com/chromedp/chromedp"
)

type task func() chromedp.Tasks

func InitializeChromeDp(submit task) {
	opts := append(chromedp.DefaultExecAllocatorOptions[:], chromedp.Flag("headless", false))
	ctx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	//new context
	nctx, cancel := chromedp.NewContext(ctx)
	defer cancel()

	//run the task
	if err := chromedp.Run(nctx, submit()); err != nil {
		panic(err)
	}
}
