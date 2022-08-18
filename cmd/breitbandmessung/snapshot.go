package breitbandmessung

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/chromedp/cdproto/browser"
	"github.com/chromedp/chromedp"
	"github.com/spf13/cobra"
)

const (
	defaultTimeout  = 2 * time.Minute
	defaultFilePerm = 0755
)

func newSnapshotCommand() *cobra.Command {
	var (
		debugMode  bool
		baseURL    string
		resultsDir string
	)

	cmd := &cobra.Command{
		Use:   "snapshot",
		Short: "snapshot current measurements",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
			defer cancel()

			if err := os.MkdirAll(resultsDir, defaultFilePerm); err != nil {
				return fmt.Errorf("make results directory: %w", err)
			}

			ctx, cancel = chromedp.NewExecAllocator(
				ctx,
				chromedp.Headless,
				chromedp.NoSandbox,
				chromedp.Flag("disable-software-rasterizer", true),
				chromedp.Flag("disable-dev-shm-usage", true),
				chromedp.Flag("use-gl", "swiftshader"),
			)
			defer cancel()

			debugF := log.Printf
			if !debugMode {
				debugF = nil
			}
			ctx, cancel = chromedp.NewContext(ctx, chromedp.WithDebugf(debugF))
			defer cancel()

			if err := chromedp.Run(
				ctx,
				chromedp.EmulateViewport(2160, 2160),
				browser.GrantPermissions([]browser.PermissionType{}),
				browser.SetDownloadBehavior(browser.SetDownloadBehaviorBehaviorAllow).
					WithDownloadPath(resultsDir).
					WithEventsEnabled(true),
			); err != nil {
				return fmt.Errorf("setup browser: %w", err)
			}

			if err := chromedp.Run(
				ctx,
				chromedp.Navigate(baseURL+"/test"),
			); err != nil {
				return fmt.Errorf("navigate: %w", err)
			}

			const startButtonSelector = "#root > div > div > div > div > div > button"
			if err := chromedp.Run(
				ctx,
				chromedp.Click(startButtonSelector, chromedp.ByQuery, chromedp.NodeVisible),
			); err != nil {
				return fmt.Errorf("click start button: %w", err)
			}

			const acceptButtonSelector = "#root > div > div > div > div > div.justify-content-between.modal-footer > button:nth-child(2)"
			if err := chromedp.Run(
				ctx,
				chromedp.Click(acceptButtonSelector, chromedp.ByQuery, chromedp.NodeVisible),
			); err != nil {
				return fmt.Errorf("click accept button: %w", err)
			}

			var suggestedFilename string
			chromedp.ListenTarget(ctx, func(v interface{}) {
				if ev, ok := v.(*browser.EventDownloadWillBegin); ok {
					suggestedFilename = ev.SuggestedFilename
				}
			})

			done := make(chan struct{})

			chromedp.ListenTarget(ctx, func(v interface{}) {
				if ev, ok := v.(*browser.EventDownloadProgress); ok {
					if ev.State == browser.DownloadProgressStateCompleted {
						done <- struct{}{}
					}
				}
			})

			const csvResultsButtonSelector = "#root > div > div > div > div > div.messung-options.col.col-12.text-md-right > button.px-0.px-sm-4.btn.btn-link"
			if err := chromedp.Run(
				ctx,
				chromedp.Click(csvResultsButtonSelector, chromedp.ByQuery, chromedp.NodeVisible),
			); err != nil {
				return fmt.Errorf("click accept policy button: %w", err)
			}

			select {
			case <-ctx.Done():
				return nil
			case <-done:
			}

			var screenshotBuf []byte

			if err := chromedp.Run(
				ctx,
				chromedp.CaptureScreenshot(&screenshotBuf),
			); err != nil {
				return fmt.Errorf("click start button: %w", err)
			}

			screenshotPath := filepath.Join(
				resultsDir,
				fmt.Sprintf("%s.jpg", suggestedFilename),
			)
			if err := ioutil.WriteFile(screenshotPath, screenshotBuf, defaultFilePerm); err != nil {
				return fmt.Errorf("save screenshot: %w", err)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&baseURL, "base-url", "u", "https://breitbandmessung.de", "Base URL.")
	cmd.Flags().StringVarP(&resultsDir, "results-dir", "d", "./results", "Results directory.")
	cmd.Flags().BoolVar(&debugMode, "debug", false, "Debug mode.")

	return cmd
}
