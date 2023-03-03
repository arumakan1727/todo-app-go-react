package restapi

func toResp[RESP any, SRC any](src *SRC, fill func(*RESP, *SRC)) RESP {
	var r RESP
	fill(&r, src)
	return r
}

func toRespArray[ITEM any, SRC any](src []SRC, fill func(*ITEM, *SRC)) []ITEM {
	a := make([]ITEM, len(src))

	// 構造体のコピーを避けるためにiでイテレート
	for i := range src {
		fill(&a[i], &src[i])
	}
	return a
}
