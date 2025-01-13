package hyperliquid

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

func ToTypedSig(r [32]byte, s [32]byte, v byte) RsvSignature {
	return RsvSignature{
		R: hexutil.Encode(r[:]),
		S: hexutil.Encode(s[:]),
		V: v,
	}
}

func ArrayAppend(data []byte, toAppend []byte) []byte {
	return append(data, toAppend...)
}

func HexToBytes(addr string) []byte {
	if strings.HasPrefix(addr, "0x") {
		fAddr := strings.Replace(addr, "0x", "", 1)
		b, _ := hex.DecodeString(fAddr)
		return b
	} else {
		b, _ := hex.DecodeString(addr)
		return b
	}
}
func HexToInt(hexString string) (*big.Int, error) {
	value := new(big.Int)
	if len(hexString) > 1 && hexString[:2] == "0x" {
		hexString = hexString[2:]
	}
	_, success := value.SetString(hexString, 16)
	if !success {
		return nil, fmt.Errorf("invalid hexadecimal string: %s", hexString)
	}
	return value, nil
}

func IntToHex(value *big.Int) string {
	return "0x" + value.Text(16) 
}

func OrderWiresToOrderAction(orders []OrderWire, grouping Grouping) PlaceOrderAction {
	return PlaceOrderAction{
		Type:     "order",
		Grouping: grouping,
		Orders:   orders,
	}
}

func OrderRequestToWire(req OrderRequest, meta map[string]AssetInfo) OrderWire {
	info := meta[req.Coin]
	return OrderWire{
		Asset:      info.AssetId,
		IsBuy:      req.IsBuy,
		LimitPx:    FloatToWire(req.LimitPx, nil),
		SizePx:     FloatToWire(req.Sz, &info.SzDecimals),
		ReduceOnly: req.ReduceOnly,
		OrderType:  OrderTypeToWire(req.OrderType),
		Cloid:      req.Cloid,
	}
}
func ModifyOrderRequestToWire(req ModifyOrderRequest, meta map[string]AssetInfo) ModifyOrderWire {
	info := meta[req.Coin]
	return ModifyOrderWire{
		OrderId: req.OrderId,
		Order: OrderWire{
			Asset:      info.AssetId,
			IsBuy:      req.IsBuy,
			LimitPx:    FloatToWire(req.LimitPx, nil),
			SizePx:     FloatToWire(req.Sz, &info.SzDecimals),
			ReduceOnly: req.ReduceOnly,
			OrderType:  OrderTypeToWire(req.OrderType),
			Cloid:      req.Cloid,
		},
	}
}

func OrderTypeToWire(orderType OrderType) OrderTypeWire {
	if orderType.Limit != nil {
		return OrderTypeWire{
			Limit: &LimitOrderType{
				Tif: orderType.Limit.Tif,
			},
			Trigger: nil,
		}
	} else if orderType.Trigger != nil {
		return OrderTypeWire{
			Trigger: &TriggerOrderType{
				TpSl:      orderType.Trigger.TpSl,
				TriggerPx: orderType.Trigger.TriggerPx,
				IsMarket:  orderType.Trigger.IsMarket,
			},
			Limit: nil,
		}
	}
	return OrderTypeWire{}
}

// Format the float with custom decimal places, default is 6.
// Hyperliquid only allows at most 6 digits.
func FloatToWire(x float64, szDecimals *int) string {
	bigf := big.NewFloat(x)
	var maxDecSz uint
	if szDecimals != nil {
		maxDecSz = uint(*szDecimals)
	} else {
		intPart, _ := bigf.Int64()
		intSize := len(strconv.FormatInt(intPart, 10))
		if intSize >= 6 {
			maxDecSz = 0
		} else {
			maxDecSz = uint(6 - intSize)
		}
	}
	x, _ = bigf.Float64()
	rounded := fmt.Sprintf("%.*f", maxDecSz, x)
	for strings.HasSuffix(rounded, "0") {
		rounded = strings.TrimSuffix(rounded, "0")
	}
	rounded = strings.TrimSuffix(rounded, ".")
	return rounded
}

// To sign raw messages via EIP-712
func StructToMap(strct any) (res map[string]interface{}, err error) {
	a, err := json.Marshal(strct)
	if err != nil {
		return map[string]interface{}{}, err
	}
	json.Unmarshal(a, &res)
	return res, nil
}
