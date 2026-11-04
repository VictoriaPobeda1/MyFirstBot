package apidata

type DataCoin struct {
	Data []Coins
}

type Coins struct {
	Name  string
	Quote CoinQuote
}

type CoinQuote struct {
	USD USDDetails
}

type USDDetails struct {
	Price, Percent_change_1h float64
}
