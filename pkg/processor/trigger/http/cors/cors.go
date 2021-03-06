/*
Copyright 2018 The Nuclio Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cors

import (
	"strconv"
	"strings"

	"github.com/nuclio/nuclio/pkg/common"

	"github.com/valyala/fasthttp"
)

type CORS struct {
	Enabled bool

	// allow configuration
	AllowOrigin      string
	AllowMethods     []string
	AllowHeaders     []string
	AllowCredentials bool

	// preflight
	PreflightRequestMethod string
	PreflightMaxAgeSeconds int64

	// computed
	allowMethodsStr           string
	allowHeadersStr           string
	preflightMaxAgeSecondsStr string
	allowCredentialsStr       string
	simpleMethods             []string
}

func NewCORS() *CORS {
	return &CORS{
		Enabled:     true,
		AllowOrigin: "*",
		AllowMethods: []string{
			fasthttp.MethodHead,
			fasthttp.MethodGet,
			fasthttp.MethodPost,
			fasthttp.MethodPut,
			fasthttp.MethodDelete,
			fasthttp.MethodOptions,
		},
		AllowHeaders: []string{
			fasthttp.HeaderAccept,
			fasthttp.HeaderContentLength,
			fasthttp.HeaderContentType,

			// nuclio custom
			"X-nuclio-log-level",
		},
		AllowCredentials:       false,
		PreflightRequestMethod: fasthttp.MethodOptions,
		PreflightMaxAgeSeconds: -1, // disable cache by default
	}
}

func (c *CORS) OriginAllowed(origin string) bool {
	if origin == "" {
		return false
	}
	return c.AllowOrigin == "*" || origin == c.AllowOrigin
}

func (c *CORS) MethodAllowed(method string) bool {
	return method != "" &&
		(method == c.PreflightRequestMethod || common.StringSliceContainsString(c.AllowMethods, method))
}

func (c *CORS) HeadersAllowed(headers []string) bool {
	for _, header := range headers {
		if !common.StringSliceContainsStringCaseInsensitive(c.AllowHeaders, header) {
			return false
		}
	}
	return true
}

func (c *CORS) EncodedAllowMethods() string {
	if c.allowMethodsStr == "" {
		c.allowMethodsStr = strings.Join(c.AllowMethods, ", ")
	}
	return c.allowMethodsStr
}

func (c *CORS) EncodeAllowHeaders() string {
	if c.allowHeadersStr == "" {
		c.allowHeadersStr = strings.Join(c.AllowHeaders, ", ")
	}
	return c.allowHeadersStr
}

func (c *CORS) EncodeAllowCredentialsHeader() string {
	if c.allowCredentialsStr == "" {
		c.allowCredentialsStr = strconv.FormatBool(c.AllowCredentials)
	}
	return c.allowHeadersStr
}

func (c *CORS) EncodePreflightMaxAgeSeconds() string {
	if c.preflightMaxAgeSecondsStr == "" {
		c.preflightMaxAgeSecondsStr = strconv.FormatInt(c.PreflightMaxAgeSeconds, 10)
	}
	return c.preflightMaxAgeSecondsStr
}
