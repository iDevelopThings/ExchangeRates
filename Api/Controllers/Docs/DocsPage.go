package Docs

import "Currencies/Api/Controllers"

func DocsPage(ctx Controllers.RequestContextImpl) {
	ctx.Html("info")
}
