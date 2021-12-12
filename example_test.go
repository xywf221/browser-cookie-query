package browser_cookie_query

func ExampleBrowserCookieQuery() {
	// use Chrome browser
	queryChrome := NewBrowserCookieQuery("Chrome")
	queryChrome.Init()              // Read SQLite and get key
	queryChrome.Query(".baidu.com") //Query

	// use Microsoft Edge browser
	queryEdge := NewBrowserCookieQuery("Microsoft Edge")
	queryEdge.Init()
	queryEdge.Query(".baidu.com")
}
