package transfer

type postTransferBody struct {
	TargetAccountID string `json:"account_destination_id"`
	Amount          int64  `json:"amount"`
}
