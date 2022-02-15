package cache_handler

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strings"

	"github.com/StevenCyb/cache_handler/store"
)

// cacheManager to record the response body from the ResponseWriter
type cacheManager struct {
	Store             store.Store
	IncludeKeyOptions []Options
	BypassOptions     []Options
}

// useOptions let the manager use given options
func (cm *cacheManager) useOptions(opts ...Options) {
	if cm.IncludeKeyOptions == nil {
		cm.IncludeKeyOptions = []Options{}
	}
	if cm.BypassOptions == nil {
		cm.BypassOptions = []Options{}
	}

	for _, opt := range opts {
		switch optT := opt.(type) {
		case AllowBypassHeader:
			cm.BypassOptions = append(cm.BypassOptions, opt)
		case AllowBypassMethod:
			optT.Key = strings.ToLower(optT.Key)
			cm.BypassOptions = append(cm.BypassOptions, optT)
		case UsePathKey:
			cm.IncludeKeyOptions = append(cm.IncludeKeyOptions, opt)
		case UseMethodKey:
			cm.IncludeKeyOptions = append(cm.IncludeKeyOptions, opt)
		case UseQueryParamsKey:
			cm.IncludeKeyOptions = append(cm.IncludeKeyOptions, opt)
		case UseHeaderKey:
			cm.IncludeKeyOptions = append(cm.IncludeKeyOptions, opt)
		}
	}
}

// keyFromRequest generate key based on request and configured options
func (cm cacheManager) keyFromRequest(r *http.Request) string {
	keyParts := []string{}
	for _, inc := range cm.IncludeKeyOptions {
		keyParts = append(keyParts, inc.ExtractString(r))
	}

	h := sha256.New()
	h.Write([]byte(strings.Join(keyParts, "/")))
	return hex.EncodeToString(h.Sum(nil))
}

// canBypass return if bypass allowed
func (cm cacheManager) canBypass(r *http.Request) bool {
	for _, bo := range cm.BypassOptions {
		if bo.ExtractBool(r) {
			return true
		}
	}

	return false
}
