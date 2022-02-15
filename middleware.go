package cache_handler

import (
	"net/http"

	"github.com/StevenCyb/cache_handler/store"
)

func NewMiddleware(next http.HandlerFunc, store store.Store, opts ...Options) http.HandlerFunc {
	cm := cacheManager{
		Store:             store,
		IncludeKeyOptions: []Options{UsePathKey{}},
	}
	cm.useOptions(opts...)

	return func(w http.ResponseWriter, r *http.Request) {
		key := cm.keyFromRequest(r)
		cachedData, err := cm.Store.Get(key)

		if cm.canBypass(r) || err != nil {
			rec := NewHttpRecorder(w)
			next.ServeHTTP(rec, r)
			cm.Store.Set(key, rec.Body.Bytes())
		} else {
			w.Write(cachedData)
		}
	}
}
