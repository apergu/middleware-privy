package model

import "math"

type Meta struct {
	TotalData int64 `json:"totalData"`
	TotalPage int64 `json:"totalPage"`
	LastID    int64 `json:"lastId"`
}

func NewMeta(total, limit, lastId int64) Meta {
	if limit <= 0 {
		return Meta{
			TotalData: total,
			LastID:    lastId,
		}
	}

	page := math.Ceil(float64(total) / float64(limit))

	return Meta{
		TotalData: total,
		TotalPage: int64(page),
		LastID:    lastId,
	}
}
