package opensea

import (
	"testing"
)

func TestGetOrders(t *testing.T) {
	is := initializeTest(t)

	//since := time.Now().Unix() - 86400
	params := OrderParams{
		AssetContractAddress: "0x91673149FFae3274b32997288395D07A8213e41F",
		TokenIds:             []string{"2012", "6574"},
	}
	ret, err := o.GetOrders(params, true)
	is.Nil(err)

	print(len(ret))

	print(*ret[0])
	print(ret[0].IsPrivate())
	print(ret[0].BasePrice.Big())
}
