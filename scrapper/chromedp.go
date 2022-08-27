package scrapper

import (
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
	"github.com/unluckythoughts/go-microservice/tools/web"
)

func _startBrowser(ctx web.Context) (context.Context, context.CancelFunc) {
	dir, err := ioutil.TempDir("", "chromedp-example")
	if err != nil {
		log.Fatal(err)
	}

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.DisableGPU,
		chromedp.UserDataDir(dir),
	)

	newCtx, _ := ctx.(context.Context)
	dpCtx, cancel1 := chromedp.NewExecAllocator(newCtx, opts...)

	dpCtx, cancel2 := chromedp.NewContext(dpCtx, chromedp.WithErrorf(ctx.Logger().Debugf))

	return dpCtx, context.CancelFunc(func() { cancel2(); cancel1(); os.RemoveAll(dir) })
}

func GetImagesAsDataUrls(ctx web.Context, urls []string) []string {
	dpCtx, cancel := _startBrowser(ctx)
	defer cancel()

	js := `
	const toDataURL = url => fetch(url)
		.then(response => response.blob())
		.then(blob => new Promise((resolve, reject) => {
			const reader = new FileReader()
			reader.onloadend = () => resolve(reader.result)
			reader.onerror = reject
			reader.readAsDataURL(blob)
		}))

		Promise.all("` + strings.Join(urls, ",") + `".split(",").map(toDataURL))
	`

	waitForCloudFlareChallenge := func(aCtx context.Context) error {
		for {
			p := network.GetAllCookies()
			cookies, _ := p.Do(aCtx)
			for _, c := range cookies {
				if c.Name == "cf_clearance" {
					return nil
				}
			}
			time.Sleep(500 * time.Millisecond)
		}
	}

	var data []string
	if err := chromedp.Run(dpCtx,
		chromedp.Navigate(urls[0]),
		chromedp.ActionFunc(waitForCloudFlareChallenge),
		chromedp.Evaluate(js, &data, func(p *runtime.EvaluateParams) *runtime.EvaluateParams {
			return p.WithAwaitPromise(true)
		}),
	); err != nil {
		ctx.Logger().With(zap.Error(err)).Error("error getting image data")
		return data
	}

	return data
}

func GetImagesAsDataUrlsFromChapter(ctx web.Context, chapterURL string) []string {
	dpCtx, cancel := _startBrowser(ctx)
	defer cancel()

	reqIDs := []network.RequestID{}
	chromedp.ListenTarget(dpCtx, func(v interface{}) {
		switch ev := v.(type) {
		case *network.EventRequestWillBeSent:
			ctx.Logger().Debugf("Event of type %+v", ev.Request.URL)
			if strings.HasPrefix(ev.Request.URL, "https://img.mghubcdn.com/file/imghub/kingdom/727") {
				reqIDs = append(reqIDs, ev.RequestID)
			}
		}
	})

	var data []string
	if err := chromedp.Run(dpCtx, chromedp.Navigate(chapterURL)); err != nil {
		ctx.Logger().With(zap.Error(err)).Error("error getting image data")
		return data
	}

	if err := chromedp.Run(dpCtx, chromedp.ActionFunc(func(dpCtx context.Context) error {
		for _, reqID := range reqIDs {
			buf, err := network.GetResponseBody(reqID).Do(dpCtx)
			if err != nil {
				ctx.Logger().With(zap.Error(err)).Error("error getting image data")
				return err
			}
			data := base64.StdEncoding.EncodeToString(buf)
			fmt.Println("$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$")
			fmt.Println(data)
			fmt.Println("$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$")
		}

		return nil
	})); err != nil {
		ctx.Logger().With(zap.Error(err)).Error("error getting image data")
		return data
	}

	return data
}
