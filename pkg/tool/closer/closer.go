package closer

import (
	"github.com/carlosrodriguesf/bank-api/pkg/tool/logger"
	"io"
)

func MustClose(log logger.Logger, closer io.Closer) {
	err := closer.Close()
	if err != nil {
		log.Fatal(err)
	}
}
