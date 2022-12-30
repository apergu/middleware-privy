package queries

import (
	"net/http"

	"gitlab.com/rteja-library3/rhelper"
)

func QueryStringKeyExist(r *http.Request, key string) bool {
	return r.URL.Query().Has(key)
}

func QueryStringInt64s(r *http.Request, key string, def int64) []int64 {
	if !QueryStringKeyExist(r, key) {
		return nil
	}

	keys := rhelper.QueryStrings(r, key)
	resp := make([]int64, len(keys))
	for i, key := range keys {
		resp[i] = rhelper.ToInt64(key, def)
	}

	return resp
}

func QueryStringToPointerInt8(r *http.Request, key string, def int8) *int8 {
	if !QueryStringKeyExist(r, key) {
		return nil
	}

	s := rhelper.QueryString(r, key)
	v := int8(rhelper.ToInt64(s, int64(def)))
	return &v
}

func QueryStringToPointerInt64(r *http.Request, key string, def int64) *int64 {
	if !QueryStringKeyExist(r, key) {
		return nil
	}

	s := rhelper.QueryString(r, key)
	v := rhelper.ToInt64(s, def)
	return &v
}

func QueryStringToPointerString(r *http.Request, key string) *string {
	if !QueryStringKeyExist(r, key) {
		return nil
	}

	v := rhelper.QueryString(r, key)
	return &v
}
