package model

type PostInfo struct {
	Long  float64 `description:"经度"`
	Lat   float64 `description:"维度"`
	Place string  `description:"定位名"`
}
