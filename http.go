package main

import (
	"net/http"

	fhttp "github.com/bogdanfinn/fhttp"
	"github.com/elazarl/goproxy"
)

func convertHttpReqToFhttpReq(req *http.Request) (*fhttp.Request, error) {
	fhttpReq, err := fhttp.NewRequest(req.Method, req.URL.String(), req.Body)
	if err != nil {
		return fhttpReq, err
	}

	// fhttpReq.Method = req.Method
	// fhttpReq.URL = req.URL
	fhttpReq.Proto = req.Proto
	fhttpReq.ProtoMajor = req.ProtoMajor
	fhttpReq.ProtoMinor = req.ProtoMinor
	fhttpReq.Header = fhttp.Header(req.Header)
	// fhttpReq.Body = req.Body
	fhttpReq.ContentLength = req.ContentLength
	fhttpReq.TransferEncoding = req.TransferEncoding
	fhttpReq.Close = req.Close
	// fhttpReq.Host = req.Host
	fhttpReq.Form = req.Form
	fhttpReq.PostForm = req.PostForm
	fhttpReq.MultipartForm = req.MultipartForm
	fhttpReq.Trailer = fhttp.Header(req.Trailer)
	fhttpReq.RemoteAddr = req.RemoteAddr
	fhttpReq.RequestURI = req.RequestURI

	return fhttpReq, nil
}

func convertFhttpRespToHttpResp(resp *fhttp.Response, req *http.Request) *http.Response {
	return &http.Response{
		Status:           resp.Status,
		StatusCode:       resp.StatusCode,
		Proto:            resp.Proto,
		ProtoMajor:       resp.ProtoMajor,
		ProtoMinor:       resp.ProtoMinor,
		Header:           http.Header(resp.Header),
		Body:             resp.Body,
		ContentLength:    resp.ContentLength,
		TransferEncoding: resp.TransferEncoding,
		Close:            resp.Close,
		Uncompressed:     resp.Uncompressed,
		Trailer:          http.Header(resp.Trailer),
		Request:          req,
		TLS:              req.TLS,
	}
}

func invalidClientProfileResponse(req *http.Request, ctx *goproxy.ProxyCtx, clientProfile string) *http.Response {
	ctx.Logf("Client specified invalid client profile: %s", clientProfile)
	return goproxy.NewResponse(req, goproxy.ContentTypeText, http.StatusBadRequest, "Invalid client profile: "+clientProfile)
}

func invalidUpstreamProxyResponse(req *http.Request, ctx *goproxy.ProxyCtx, upstreamProxy string) *http.Response {
	ctx.Logf("Client specified invalid upstream proxy: %s", upstreamProxy)
	return goproxy.NewResponse(req, goproxy.ContentTypeText, http.StatusBadRequest, "Invalid upstream proxy: "+upstreamProxy)
}
