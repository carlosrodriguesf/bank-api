package account

type postAccountBody struct {
	Name     string `json:"name"`
	Document string `json:"document"`
	Secret   string `json:"secret"`
	Balance  int64  `json:"balance"`
}
