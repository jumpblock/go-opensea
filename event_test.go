package opensea

import (
	"fmt"
	"testing"
)

func TestRetrievingEvents(t *testing.T) {
	is := initializeTest(t)

	params := NewRetrievingEventsParams()
	//err := params.SetAssetContractAddress(contract)
	//is.Nil(err)

	//params.OccurredAfter = time.Now().Unix() - 86400
	//params.OccurredBefore = time.Now().Unix()
	params.OccurredAfter = 1650112349
	params.OccurredBefore = 1650119549
	params.EventType = EventTypeCreated
	//params.AssetContractAddress = "0x7Bd29408f11D2bFC23c34f18275bBf23bB716Bc7"
	ret, err := o.RetrievingEvents(params)
	is.Nil(err)
	print(len(ret))
	for _, v := range ret {
		fmt.Println("===>:", v.Asset.TokenID, v.EventType, v.AuctionType, v.StartingPrice, v.EndingPrice, v.CreatedDate.Time(), v.ListingTime)
	}
}
