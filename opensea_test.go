package opensea

import (
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"reflect"
	"testing"

	"github.com/cheekybits/is"
)

var (
	o          = &Opensea{}
	owner      = "0xd868711BD9a2C6F1548F5f4737f71DA67d821090"
	contract   = "0xD1E5b0FF1287aA9f9A268759062E4Ab08b9Dacbe"
	tokenID, _ = new(big.Int).SetString("68193319175094895046676294033579301732477745586436552446948983324346430262893", 0)
)

func TestGetSingleAsset(t *testing.T) {
	is := initializeTest(t)

	ret, err := o.GetSingleAsset(contract, tokenID)
	is.Nil(err)

	print(*ret)
}
func TestGetAssets(t *testing.T) {
	is := initializeTest(t)
	ret, err := o.GetAssetDetail("0xfacb4bc7bfa5007d91095564e05918c82079befb")
	is.Nil(err)

	print(ret.Collection.Slug)
}
func TestGetCollections(t *testing.T) {
	is := initializeTest(t)
	res, err := o.GetCollections(10000, 300)
	is.Nil(err)
	by, _ := json.Marshal(res)
	fmt.Println(len(res), string(by))

	cs, err := o.GetSingleCollection(res[0].Slug)
	is.Nil(err)
	fmt.Println(cs)
}

func initializeTest(t *testing.T) is.I {
	is := is.New(t)
	var err error

	o, err = NewOpensea(os.Getenv("API_KEY"))
	is.Nil(err)
	return is
}

func print(in interface{}) {
	if reflect.TypeOf(in).Kind() == reflect.Struct {
		in, _ = json.Marshal(in)
		in = string(in.([]byte))
	}
	fmt.Println(in)
}
