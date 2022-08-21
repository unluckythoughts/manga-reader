package scrapper

import (
	"log"
	"net/url"
	"os/exec"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/raff/godet"
	"github.com/unluckythoughts/go-microservice/tools/web"
	"go.uber.org/zap"
)

func startChromium(ctx web.Context) {
	chromeapp := "/usr/bin/chromium"
	chromeappArg := []string{"--headless", "--hide-scrollbars", "--remote-debugging-port=9222", "--disable-gpu", "--allow-insecure-localhost"}
	cmd := exec.Command(chromeapp, chromeappArg...)
	err := cmd.Start()
	if err != nil {
		ctx.Logger().Error("cannot start browser", err)
	}
	ctx.Logger().Debug("Started chromium headless")
}

func killChromium(ctx web.Context) {
	cmd := exec.Command("pkill", "-9", "chromium")
	err := cmd.Start()
	if err != nil {
		ctx.Logger().Error("cannot kill processes", err)
	}
	ctx.Logger().Debug("Killed chromium headless")
}

func connectToDebugger(ctx web.Context) (*godet.RemoteDebugger, error) {
	for {
		time.Sleep(100 * time.Millisecond)
		remote, err := godet.Connect("localhost:9222", false)
		if err != nil {
			_, ok := err.(*url.Error)
			if ok {
				continue
			}
			ctx.Logger().With(zap.Error(err)).Error("cannot connect to Chrome instance")
			return nil, err
		}

		// settings
		remote.LogEvents(false)
		remote.PageEvents(true)
		remote.DOMEvents(true)

		return remote, nil
	}
}

func injectCallback(ctx web.Context, remote *godet.RemoteDebugger, script string, data *respData) godet.EventCallback {
	return func(params godet.Params) {
		scriptResp, err := remote.EvaluateWrap(script)
		if err != nil {
			log.Println("Error executing js block:", err)
			return
		}

		vals := []string{}
		if val, ok := scriptResp.([]interface{}); ok {
			for _, v := range val {
				if str, ok := v.(string); ok {
					vals = append(vals, str)
				}
			}
		}

		data.data = vals
	}
}

func validOutput(data []string) bool {
	if len(data) > 0 {
		return true
	}

	return false
}

func GetInjectionScript(selector string) string {
	innerHTMLSnippet := `data.push(e.InnerHtml())`

	attrSnippet := `
		if (e.getAttribute('###Attribute###') !== '') {
			data.push(e.getAttribute('###Attribute###'))
		}
	`

	script := `
		var data = []
		document.querySelectorAll('###Selector####').forEach(e => {
			###Snippet###
		})
		return data
	`

	selector, attrs, ok := hasDataInAttr(selector)
	if !ok {
		script = strings.Replace(script, "###Selector####", selector, -1)
		return strings.Replace(script, "###Snippet###", innerHTMLSnippet, -1)
	}

	script = strings.Replace(script, "###Selector####", selector, -1)
	if len(attrs) == 1 {
		attrSnippet = strings.Replace(attrSnippet, "###Attribute###", attrs[0], -1)
		return strings.Replace(script, "###Snippet###", attrSnippet, -1)
	}

	attrSnippets := []string{}
	for _, attr := range attrs {
		attrSnippets = append(attrSnippets, strings.Replace(attrSnippet, "###Attribute###", attr, -1))
	}
	attrSnippet = strings.Join(attrSnippets, " else ")
	return strings.Replace(script, "###Snippet###", attrSnippet, 1)
}

type respData struct {
	data []string
}

func SimulateBrowser(ctx web.Context, url, injectScript string) ([]string, error) {
	startChromium(ctx)
	defer killChromium(ctx)

	// connect to Chromium instance
	remote, err := connectToDebugger(ctx)
	if err != nil {
		return []string{}, err
	}

	// disconnect when done
	defer remote.Close()

	_, err = remote.Navigate(url)
	if err != nil {
		return []string{}, err
	}

	resp := respData{}
	remote.CallbackEvent("Page.frameStoppedLoading", injectCallback(
		ctx, remote, injectScript, &resp,
	))

	for i := 10; i > 0; i-- {
		output := resp.data
		if validOutput(output) {
			return output, nil
		}
		time.Sleep(1 * time.Second)
	}

	return []string{}, errors.Errorf("could not get data for url %s", url)
}
