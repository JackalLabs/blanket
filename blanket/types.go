package blanket

type SpaceResponse struct {
	TotalSpace int64 `json:"total_space"`
	UsedSpace  int64 `json:"used_space"`
	FreeSpace  int64 `json:"free_space"`
}

type BalanceResponse struct {
	Balance Balance `json:"balance"`
}

type PriceResponse struct {
	Price float64 `json:"price"`
}

type Balance struct {
	Denom  string `json:"denom"`
	Amount string `json:"amount"`
}

type APIResponse struct {
	Pagination Pagination `json:"pagination"`
}

type Pagination struct {
	Count string `json:"total"`
}
