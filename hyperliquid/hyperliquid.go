package hyperliquid

type IHyperliquid interface {
	IExchangeAPI
	IInfoAPI
}

type Hyperliquid struct {
	ExchangeAPI
	InfoAPI
}

// HyperliquidClientConfig is a configuration struct for Hyperliquid API.
// PrivateKey can be empty if you only need to use the public endpoints.
// AccountAddress is the default account address for the API that can be changed with SetAccountAddress().
// AccountAddress may be different from the address build from the private key due to Hyperliquid's account system.
type HyperliquidClientConfig struct {
	IsMainnet      bool
	PrivateKey     string
	AccountAddress string
}

func NewHyperliquid(config *HyperliquidClientConfig) *Hyperliquid {
	var defaultConfig *HyperliquidClientConfig
	if config == nil {
		defaultConfig = &HyperliquidClientConfig{
			IsMainnet:      true,
			PrivateKey:     "",
			AccountAddress: "",
		}
	} else {
		defaultConfig = config
	}
	exchangeAPI := NewExchangeAPI(defaultConfig.IsMainnet)
	exchangeAPI.SetPrivateKey(defaultConfig.PrivateKey)
	exchangeAPI.SetAccountAddress(defaultConfig.AccountAddress)
	infoAPI := NewInfoAPI(defaultConfig.IsMainnet)
	infoAPI.SetAccountAddress(defaultConfig.AccountAddress)

	roleResp, err := infoAPI.GetUserRole(config.AccountAddress)
	if err != nil {
		exchangeAPI.debug("Error querying user role: %s", err)
	} else {
		exchangeAPI.role = roleResp.Role
	}

	return &Hyperliquid{
		ExchangeAPI: *exchangeAPI,
		InfoAPI:     *infoAPI,
	}
}

func (h *Hyperliquid) SetDebugActive() {
	h.ExchangeAPI.SetDebugActive()
	h.InfoAPI.SetDebugActive()
}

func (h *Hyperliquid) SetPrivateKey(privateKey string) error {
	err := h.ExchangeAPI.SetPrivateKey(privateKey)
	if err != nil {
		return err
	}
	return nil
}

func (h *Hyperliquid) SetAccountAddress(accountAddress string) {
	h.ExchangeAPI.SetAccountAddress(accountAddress)
	h.InfoAPI.SetAccountAddress(accountAddress)
}

func (h *Hyperliquid) AccountAddress() string {
	return h.ExchangeAPI.AccountAddress()
}

func (h *Hyperliquid) IsMainnet() bool {
	return h.ExchangeAPI.IsMainnet()
}
