package transfer

import "github.com/carlosrodriguesf/bank-api/pkg/model"

type transferWrapper struct {
	AccountOrigin *model.Account
	AccountTarget *model.Account
	Transfer      model.Transfer
}
