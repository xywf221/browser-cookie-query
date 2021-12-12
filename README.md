# browser cookie query

Only supported temporarily macOS,theoretically,it supports all `chromium` browsers

valid browser : `Edge` `Chrome`

## Installation
exec `go get -u github.com/xywf221/browser-cookie-query`

## Usage
Query the cookie of Microsoft Edge browser `github.com`

```go
queryEdge := NewBrowserCookieQuery("Microsoft Edge")
queryEdge.Init()
queryEdge.Query(".github.com")
```


## Todo
- [ ] Windows supported


## Reference projects

[pycookiecheat](https://github.com/n8henrie/pycookiecheat)

[gookies](https://github.com/CCob/gookies)