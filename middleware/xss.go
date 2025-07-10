package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"store/utils"

	"github.com/gin-gonic/gin"
	"github.com/microcosm-cc/bluemonday"
)

type CustomReadCloser struct {
	Body []byte
}

func (c *CustomReadCloser) Read(p []byte) (n int, err error) {
	copy(p, c.Body)
	return len(c.Body), io.EOF
}

func (c *CustomReadCloser) Close() error {
	return nil
}

func XSSProtectionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		p := bluemonday.UGCPolicy()

		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read request body"})
			return
		}

		//for body param
		if len(body) != 0 {
			var inputMap map[string]any
			if err := json.Unmarshal(body, &inputMap); err != nil {
				utils.Response(c, http.StatusBadRequest, nil, nil, nil, err.Error(), nil, true)
				return
			}
			detectMaliousContent(inputMap, p, c)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
		}

		//for form param
		for _, values := range c.Request.Form {
			for _, value := range values {
				clean := p.Sanitize(value)
				if clean != value {
					utils.Response(c, http.StatusBadRequest, nil, nil, nil, fmt.Sprintf("XSS-like content detected:%s", value), nil, true)
					return
				}
			}
		}

		//for query param
		for _, values := range c.Request.URL.Query() {
			for _, value := range values {
				clean := p.Sanitize(value)
				if clean != value {
					utils.Response(c, http.StatusBadRequest, nil, nil, nil, fmt.Sprintf("XSS-like content detected:%s", value), nil, true)
					return
				}
			}
		}

		for _, param := range c.Params {
			clean := p.Sanitize(param.Value)
			if clean != param.Value {
				utils.Response(c, http.StatusBadRequest, nil, nil, nil, fmt.Sprintf("XSS-like content detected:%s", param.Value), nil,true)
				return
			}
		}

		c.Next()
	}
}

func detectMaliousContent(data map[string]any, policy *bluemonday.Policy, c *gin.Context) bool {
	for key, value := range data {
		if key == "description" {
			continue
		}
		switch v := value.(type) {
		case string:
			// Sanitize the string field
			sanitizedValue := policy.Sanitize(v)

			// If the sanitized value is different from the original, it indicates malicious content
			if sanitizedValue != v {
				handleMaliciousContent(key, v, c)
				return false // Stop processing further if malicious content is found
			}

			// Update the map with sanitized value
			data[key] = sanitizedValue

		case map[string]any:
			// Recursively sanitize nested maps
			if !detectMaliousContent(v, policy, c) {
				return false
			}

		case []any:
			// Sanitize elements of arrays if necessary
			for i, elem := range v {
				switch e := elem.(type) {
				case string:
					sanitizedElem := policy.Sanitize(e)
					// If the sanitized element is different from the original, indicate malicious content
					if sanitizedElem != e {
						handleMaliciousContent(fmt.Sprintf("%d", i), e, c)
						return false
					}
					v[i] = sanitizedElem

				case map[string]any:
					// Recursively sanitize nested maps in arrays
					if !detectMaliousContent(e, policy, c) {
						return false
					}
				}
			}
		}
	}

	return true
}

func handleMaliciousContent(key string, value string, c *gin.Context) {
	utils.Response(c, http.StatusBadRequest, nil, nil, nil, fmt.Sprintf("Malicious content detected in field: %s, value: %s", key, value), nil, true)
}
