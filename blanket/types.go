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

type IndexResponse struct {
	Address string `json:"address"`
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

type Providers struct {
	Providers Provider `json:"providers"`
}

type Provider struct {
	BurnedContracts string `json:"burned_contracts"`
}

type BlockResponse struct {
	Block Block `json:"block"`
}

type Block struct {
	Header BlockHeader `json:"header"`
}

type BlockHeader struct {
	Height string `json:"height"`
}

type ParamResponse struct {
	Params Params `json:"params"`
}

type Params struct {
	ProofWindow string `json:"proof_window"`
}
