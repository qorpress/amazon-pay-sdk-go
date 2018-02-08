package amazonpay

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/url"
	"sort"
	"strings"
	"time"
)

// Post post API info
func (amazonPay *AmazonPay) Post(path string, params Params) error {
	if _, ok := params.Get("AWSAccessKeyId"); !ok {
		params.Set("AWSAccessKeyId", amazonPay.Config.AccessKey)
	}

	if _, ok := params.Get("SellerId"); !ok {
		params.Set("SellerId", amazonPay.Config.MerchantID)
	}

	if _, ok := params.Get("SignatureMethod"); !ok {
		params.Set("SignatureMethod", "HmacSHA256")
	}

	if _, ok := params.Get("SignatureVersion"); !ok {
		params.Set("SignatureVersion", "2")
	}

	if _, ok := params.Get("Timestamp"); !ok {
		params.Set("Timestamp", time.Now().UTC().Format("2006-01-02T03:04:05Z"))
	}

	params.Set("Signature", params.Sign())

	if _, ok := params.Get("Version"); !ok {
		params.Set("Version", "2013-01-01")
	}

	// retry
	return nil
}

// buildPostURL build post URL
func (amazonPay *AmazonPay) buildPostURL(params Params) string {
	apiParams := []string{}

	for key, value := range params {
		if str := fmt.Sprint(value); str != "" {
			apiParams = append(apiParams, key+"="+url.QueryEscape(str))
		}
	}

	sort.Strings(apiParams)
	postURL := strings.Join(apiParams, "&")
	postURL += "&Signature=" + amazonPay.Sign(strings.Join([]string{"POST", amazonPay.Config.Endpoint, fmt.Sprintf("/%v/%v", amazonPay.Config.ModePath, amazonPay.Config.APIVersion), postURL}, "\n"))
	return postURL
}

// Sign sign messages
func (amazonPay *AmazonPay) Sign(message string) string {
	key := []byte(amazonPay.Config.SecretKey)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(message))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
