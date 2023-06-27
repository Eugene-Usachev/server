package fasthttp_utils

import (
	"github.com/Eugene-Usachev/fastbytes"
	"github.com/valyala/fasthttp"
)

var (
	authHeader = "Authorization"
)

func GetAuthorizationHeader(ctx *fasthttp.RequestCtx) []byte {
	var header []byte
	//faster than Peek
	ctx.Request.Header.VisitAll(func(k, v []byte) {
		if fastbytes.B2S(k) == authHeader {
			header = v
		}
	})
	return header
}
