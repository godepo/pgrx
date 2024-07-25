package example

import (
	"time"

	"github.com/google/uuid"
)

type KindOfPayment string

const (
	KOPIncoming KindOfPayment = "INCOMING"
	KOPOutgoing KindOfPayment = "OUTGOING"
	KOPInternal KindOfPayment = "INTERNAL"
	KOPUnknown  KindOfPayment = "UNKNOWN"
)

type KindOfState string

const (
	KOSCreated    KindOfState = "CREATED"
	KOSDeclined   KindOfState = "DECLINED"
	KOSSucceeded  KindOfState = "SUCCEEDED"
	KOSProcessing KindOfState = "PROCESSING"
	KOSDelayed    KindOfState = "DELAYED"
	KOSAborted    KindOfState = "ABORTED"
	KOSUnknown    KindOfState = "UNKNOWN"
)

type KindOfCurrency string

const (
	USD KindOfCurrency = "USD"
	EUR KindOfCurrency = "EUR"
	GBP KindOfCurrency = "GBP"
)

type Payment struct {
	ID          uuid.UUID      `json:"id" db:"id"`
	UserID      string         `json:"user_id" db:"user_id"`
	Kind        KindOfPayment  `json:"kind" db:"kind"`
	State       KindOfState    `json:"state" db:"state"`
	Currency    KindOfCurrency `json:"currency" db:"currency"`
	Amount      int64          `json:"amount" db:"amount"`
	CreatedAt   time.Time      `json:"created" db:"created_at"`
	ProcessedAt *time.Time     `json:"processed" db:"processed_at"`
}
