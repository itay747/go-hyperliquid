package hyperliquid

import (
	"os"
	"testing"
)

func GetHyperliquidAPI() *Hyperliquid {
	hl := NewHyperliquid(&ClientConfig{
		IsMainnet:      false,
		AccountAddress: os.Getenv("TEST_ADDRESS"),
		PrivateKey:     os.Getenv("TEST_PRIVATE_KEY"),
	})
	if GLOBAL_DEBUG {
		hl.SetDebugActive()
	}
	return hl
}

func TestHyperliquid_CheckFieldsConsistency(t *testing.T) {
	hl := GetHyperliquidAPI()
	if hl.ExchangeAPI.baseEndpoint != "/exchange" {
		t.Errorf("baseEndpoint = %v, want %v", hl.ExchangeAPI.baseEndpoint, "/exchange")
	}
	if hl.InfoAPI.baseEndpoint != "/info" {
		t.Errorf("baseEndpoint = %v, want %v", hl.InfoAPI.baseEndpoint, "/info")
	}
	if hl.InfoAPI.baseURL != "https://api.hyperliquid.xyz" {
		t.Errorf("baseUrl = %v, want %v", hl.InfoAPI.baseURL, "https://api.hyperliquid.com")
	}
	hl.SetDebugActive()
	if hl.InfoAPI.Debug != hl.ExchangeAPI.Debug {
		t.Errorf("debug = %v, want %v", hl.InfoAPI.Debug, hl.ExchangeAPI.Debug)
	}
	newAddress := "0x1234567890"
	hl.SetAccountAddress(newAddress)
	if hl.InfoAPI.AccountAddress() != newAddress {
		t.Errorf("AccountAddress = %v, want %v", hl.InfoAPI.AccountAddress(), newAddress)
	}
	if hl.ExchangeAPI.AccountAddress() != newAddress {
		t.Errorf("AccountAddress = %v, want %v", hl.ExchangeAPI.AccountAddress(), newAddress)
	}
	if hl.AccountAddress() != newAddress {
		t.Errorf("AccountAddress = %v, want %v", hl.AccountAddress(), newAddress)
	}
	if hl.infoAPI.AccountAddress() != newAddress {
		t.Errorf("AccountAddress = %v, want %v", hl.infoAPI.AccountAddress(), newAddress)
	}
}

func TestHyperliquid_MakeSomeTradingLogic(t *testing.T) {
	client := GetHyperliquidAPI()

	// Make limit order
	res1, err := client.LimitOrder(TifGtc, "ETH", 0.01, 1234.1, false)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	t.Logf("LimitOrder(TifIoc, ETH, 0.01, 1234.1, false): %v", res1)

	res2, err := client.LimitOrder(TifGtc, "ETH", 0.01, 1200.1, true)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	t.Logf("LimitOrder(TifGtc, ETH, 0.01, 1200.1, true): %v", res2)

	res3, err := client.LimitOrder(TifGtc, "ETH", -0.01, 5000.1, true)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	t.Logf("LimitOrder(TifGtc, ETH, -0.01, 5000.1, true): %v", res3)

	// Get all ordres
	res4, err := client.GetAccountOpenOrders()
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	t.Logf("GetAccountOpenOrders(): %v", res4)

	// Close all orders
	res5, err := client.CancelAllOrders()
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	t.Logf("CancelAllOrders(): %v", res5)

	// Make market order
	res6, err := client.MarketOrder("ETH", 0.01, nil)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	t.Logf("MarketOrder(ETH, 0.01, nil): %v", res6)

	// Close position
	res7, err := client.ClosePosition("ETH")
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	t.Logf("ClosePosition(ETH): %v", res7)

	// Get account balance
	res8, err := client.GetAccountState()
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	t.Logf("GetAccountState(): %v", res8)
}
