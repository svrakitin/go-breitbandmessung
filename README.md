# go-breitbandmessung

This project automates [breitbandmessung.de](https://breitbandmessung.de) speedtest (via [chromedp/chromedp](https://github.com/chromedp/chromedp)).

It started as a frustration with [PYUR (TeleColumbus)](https://www.pyur.com) having issues all the time.

```
Usage:
  breitbandmessung snapshot [flags]

Flags:
  -u, --base-url string      Base URL. (default "https://breitbandmessung.de")
      --debug                Debug mode.
  -h, --help                 help for snapshot
  -d, --results-dir string   Results directory. (default "./results")
```