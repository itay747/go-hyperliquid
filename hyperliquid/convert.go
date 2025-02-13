package hyperliquid

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
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

func (req *OrderRequest) isSpot() bool {
	return strings.ContainsAny(req.Coin, "@-")
}

func (req *ModifyOrderRequest) isSpot() bool {
	return strings.ContainsAny(req.Coin, "@-")
}

// ToWire (OrderRequest) converts an OrderRequest to an OrderWire using the provided metadata.
// https://hyperliquid.gitbook.io/hyperliquid-docs/for-developers/api/asset-ids
func (req *OrderRequest) ToWire(meta map[string]AssetInfo) OrderWire {
	info := meta[req.Coin]
	return OrderWire{
		Asset:      info.AssetID,
		IsBuy:      req.IsBuy,
		LimitPx:    FloatToWire(req.LimitPx, nil),
		SizePx:     FloatToWire(req.Sz, &info.SzDecimals),
		ReduceOnly: req.ReduceOnly,
		OrderType:  OrderTypeToWire(req.OrderType),
		Cloid:      req.Cloid,
	}
}

// ToWire (ModifyOrderRequest) converts a ModifyOrderRequest to a ModifyOrderWire using the provided metadata.
func (req *ModifyOrderRequest) ToWire(meta map[string]AssetInfo) ModifyOrderWire {
	info := meta[req.Coin]
	return ModifyOrderWire{
		OrderID: req.OrderID,
		Order: OrderWire{
			Asset:      info.AssetID,
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

// RoundOrderSize rounds the order size to the nearest tick size
func RoundOrderSize(x float64, szDecimals int) string {
	newX := math.Round(x*math.Pow10(szDecimals)) / math.Pow10(szDecimals)
	// TODO: add rounding
	return big.NewFloat(newX).Text('f', szDecimals)
}

// RoundOrderPrice rounds the order price to the nearest tick size
func RoundOrderPrice(x float64, szDecimals int, maxDecimals int) string {
	maxSignFigures := 5
	allowedDecimals := maxDecimals - szDecimals
	numberOfDigitsInIntegerPart := len(strconv.Itoa(int(x)))
	if numberOfDigitsInIntegerPart >= maxSignFigures {
		return RoundOrderSize(x, 0)
	}
	allowedSignFigures := maxSignFigures - numberOfDigitsInIntegerPart
	if x >= 1 {
		return RoundOrderSize(x, min(allowedSignFigures, allowedDecimals))
	}

	text := RoundOrderSize(x, allowedDecimals)
	startSignFigures := false
	for i := 2; i < len(text); i++ {
		if text[i] == '0' && !startSignFigures {
			continue
		}
		startSignFigures = true
		allowedSignFigures--
		if allowedSignFigures == 0 {
			return text[:i+1]
		}
	}
	return text
}
