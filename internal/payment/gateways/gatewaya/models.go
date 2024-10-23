package gatewaya

type Request struct {
	Amount      float64 `json:"amount"`
	CallbackURL string  `json:"callback_url"`
}

type Response struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}
