/*
 * Payment API
 *
 * Demo code
 *
 * API version: 1.0.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package swagger

type PaymentAttributesFx struct {

	ContractReference string `json:"contract_reference,omitempty"`

	ExchangeRate string `json:"exchange_rate,omitempty"`

	OriginalAmount string `json:"original_amount,omitempty"`

	OriginalCurrency string `json:"original_currency,omitempty"`
}