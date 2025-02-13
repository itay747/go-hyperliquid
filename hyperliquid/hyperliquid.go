package hyperliquid

type IHyperliquid interface {
	IExchangeAPI
	IInfoAPI
}

type Hyperliquid struct {
	ExchangeAPI
	InfoAPI
}

// ClientConfig is a configuration struct for Hyperliquid API.
// PrivateKey can be empty if you only need to use the public endpoints.
// AccountAddress is the default account address for the API that can be changed with SetAccountAddress().
// AccountAddress may be different from the address build from the private key due to Hyperliquid's account system.
type ClientConfig struct {
	IsMainnet      bool
	PrivateKey     string
	AccountAddress string
}

// NewHyperliquid creates a new Hyperliquid API client.
func NewHyperliquid(config *ClientConfig) *Hyperliquid {
	var defaultConfig *ClientConfig
	if config == nil {
		defaultConfig = &ClientConfig{
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
	hl := &Hyperliquid{
		ExchangeAPI: *exchangeAPI,
		InfoAPI:     *infoAPI,
	}
	hl.UpdateAccountVaultAddress()
	return hl
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

func (h *Hyperliquid) UpdateAccountVaultAddress() {
	h.UpdateVaultAddress(h.AccountAddress())
}

func (h *Hyperliquid) UpdateVaultAddress(address string) error {
	role, err := h.InfoAPI.GetUserRole(address)
	if err != nil {
		return err
	}
	h.ExchangeAPI.SetUserRole(role.UserRole.Role)
	return nil
}
func (h *Hyperliquid) AccountAddress() string {
	return h.ExchangeAPI.AccountAddress()
}

func (h *Hyperliquid) IsMainnet() bool {
	return h.ExchangeAPI.IsMainnet()
}
