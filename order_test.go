package opensea

import (
	"fmt"
	"testing"
)

func TestGetOrders(t *testing.T) {
	is := initializeTest(t)

	//since := time.Now().Unix() - 86400
	params := OrderParams{
		AssetContractAddress: "0x4af69be25f4eb13dab39af246d607d643fe71968",
		TokenIds:             []string{"637", "901"},
		Side:                 fmt.Sprintf("%d", Sell),
	}
	ret, err := o.GetOrders(params, false)
	is.Nil(err)

	print(len(ret))
	for _, v := range ret {
		fmt.Println(v.Asset.TokenID, v)
	}
}
