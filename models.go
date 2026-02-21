package intelligentdata

// ── Address Validation ────────────────────────────────────────────────────

type AddressRequest struct {
	AddressLine1 string `json:"addressLine1"`
	AddressLine2 string `json:"addressLine2,omitempty"`
	City         string `json:"city"`
	State        string `json:"state,omitempty"`
	PostalCode   string `json:"postalCode,omitempty"`
	Country      string `json:"country"`
}

type AddressResponse struct {
	IsValid              bool              `json:"isValid"`
	ConfidenceScore      float64           `json:"confidenceScore"`
	StandardizedAddress  map[string]string `json:"standardizedAddress"`
	Raw                  map[string]interface{}
}

// ── Tax ID Validation ─────────────────────────────────────────────────────

type TaxIdRequest struct {
	TaxID     string `json:"taxId"`
	Country   string `json:"country"`
	TaxIDType string `json:"taxIdType,omitempty"`
}

type TaxIdResponse struct {
	IsValid        bool   `json:"isValid"`
	TaxIDType      string `json:"taxIdType"`
	Country        string `json:"country"`
	RegisteredName string `json:"registeredName"`
	Raw            map[string]interface{}
}

// ── Bank Account Validation ───────────────────────────────────────────────

type BankAccountRequest struct {
	AccountNumber string `json:"accountNumber"`
	Country       string `json:"country"`
	RoutingNumber string `json:"routingNumber,omitempty"`
	IBAN          string `json:"iban,omitempty"`
	BankCode      string `json:"bankCode,omitempty"`
}

type BankAccountResponse struct {
	IsValid     bool   `json:"isValid"`
	BankName    string `json:"bankName"`
	AccountType string `json:"accountType"`
	Raw         map[string]interface{}
}

// ── Business Lookup ───────────────────────────────────────────────────────

type BusinessLookupRequest struct {
	CompanyName        string `json:"companyName"`
	Country            string `json:"country"`
	RegistrationNumber string `json:"registrationNumber,omitempty"`
	State              string `json:"state,omitempty"`
}

type BusinessLookupResponse struct {
	Found              bool              `json:"found"`
	CompanyName        string            `json:"companyName"`
	RegistrationNumber string            `json:"registrationNumber"`
	Status             string            `json:"status"`
	Address            map[string]string `json:"address"`
	Raw                map[string]interface{}
}

// ── Sanctions Screening ──────────────────────────────────────────────────

type SanctionsRequest struct {
	EntityName string `json:"entityName"`
	EntityType string `json:"entityType,omitempty"`
	Country    string `json:"country,omitempty"`
}

type SanctionsResponse struct {
	HasMatches   bool                     `json:"hasMatches"`
	Matches      []map[string]interface{} `json:"matches"`
	ScreenedLists []string                `json:"screenedLists"`
	Raw          map[string]interface{}
}

// ── Directors Check ───────────────────────────────────────────────────────

type DirectorsRequest struct {
	CompanyName        string `json:"companyName"`
	Country            string `json:"country"`
	RegistrationNumber string `json:"registrationNumber,omitempty"`
}

type DirectorsResponse struct {
	HasDisqualified bool                     `json:"hasDisqualified"`
	Directors       []map[string]interface{} `json:"directors"`
	Raw             map[string]interface{}
}
