package websockets

import (
	"encoding/json"
	"fmt"
	"sync/atomic"

	"github.com/ffddw/ripple/data"
)

var counter uint64

type Syncer interface {
	Done()
	Fail(message string)
}

type CommandError struct {
	Name      string `json:"error"`
	Code      int    `json:"error_code"`
	Message   string `json:"error_message"`
	Exception string `json:"error_exception"`
}

type Command struct {
	*CommandError
	Id     uint64        `json:"id"`
	Name   string        `json:"command"`
	Type   string        `json:"type,omitempty"`
	Status string        `json:"status,omitempty"`
	Ready  chan struct{} `json:"-"`
}

func (c *Command) Done() {
	c.Ready <- struct{}{}
}

func (c *Command) Fail(message string) {
	c.CommandError = &CommandError{
		Name:    "Client Error",
		Code:    -1,
		Message: message,
	}
	c.Ready <- struct{}{}
}

func (c *Command) IncrementId() {
	c.Id = atomic.AddUint64(&counter, 1)
}

func (e *CommandError) Error() string {
	return fmt.Sprintf("%s %d %s %s", e.Name, e.Code, e.Message, e.Exception)
}

func newCommand(command string) *Command {
	return &Command{
		Id:    atomic.AddUint64(&counter, 1),
		Name:  command,
		Ready: make(chan struct{}),
	}
}

type AccountTxCommand struct {
	*Command
	Account   data.Account           `json:"account"`
	MinLedger int64                  `json:"ledger_index_min"`
	MaxLedger int64                  `json:"ledger_index_max"`
	Binary    bool                   `json:"binary,omitempty"`
	Forward   bool                   `json:"forward,omitempty"`
	Limit     int                    `json:"limit,omitempty"`
	Marker    map[string]interface{} `json:"marker,omitempty"`
	Result    *AccountTxResult       `json:"result,omitempty"`
}

type AccountTxResult struct {
	Marker       map[string]interface{} `json:"marker,omitempty"`
	Transactions data.TransactionSlice  `json:"transactions,omitempty"`
}

func newAccountTxCommand(account data.Account, pageSize int, marker map[string]interface{}, minLedger, maxLedger int64) *AccountTxCommand {
	return &AccountTxCommand{
		Command:   newCommand("account_tx"),
		Account:   account,
		MinLedger: minLedger,
		MaxLedger: maxLedger,
		Limit:     pageSize,
		Marker:    marker,
	}
}

func newBinaryLedgerDataCommand(ledger interface{}, marker *data.Hash256) *BinaryLedgerDataCommand {
	return &BinaryLedgerDataCommand{
		Command: newCommand("ledger_data"),
		Ledger:  ledger,
		Binary:  true,
		Marker:  marker,
	}
}

type TxCommand struct {
	*Command
	Transaction data.Hash256 `json:"transaction"`
	Result      *TxResult    `json:"result,omitempty"`
}

type TxResult struct {
	data.TransactionWithMetaData
	Validated bool `json:"validated"`
}

// A shim to populate the Validated field before passing
// control on to TransactionWithMetaData.UnmarshalJSON
func (txr *TxResult) UnmarshalJSON(b []byte) error {
	var extract map[string]interface{}
	if err := json.Unmarshal(b, &extract); err != nil {
		return err
	}
	// "validated" can be absent, when tx result is provisional.
	validated, ok := extract["validated"]
	if !ok {
		txr.Validated = false
	} else {
		txr.Validated = validated.(bool)
	}
	return json.Unmarshal(b, &txr.TransactionWithMetaData)
}

type SubmitCommand struct {
	*Command
	TxBlob string        `json:"tx_blob"`
	Result *SubmitResult `json:"result,omitempty"`
}

type SubmitResultTxJSON struct {
	Account            string `json:"Account"`
	Domain             string `json:"Domain"`
	Fee                string `json:"Fee"`
	Flags              uint32 `json:"Flags"`
	LastLedgerSequence uint32 `json:"LastLedgerSequence"`
	Sequence           uint32 `json:"Sequence"`
	SigningPubKey      string `json:"SigningPubKey"`
	TransactionType    string `json:"TransactionType"`
	TxnSignature       string `json:"TxnSignature"`
	Hash               string `json:"hash"`
}

type SubmitResult struct {
	EngineResult        data.TransactionResult `json:"engine_result"`
	EngineResultCode    int                    `json:"engine_result_code"`
	EngineResultMessage string                 `json:"engine_result_message"`
	TxBlob              string                 `json:"tx_blob"`
	Tx                  *SubmitResultTxJSON    `json:"tx_json"`
}

type LedgerCommand struct {
	*Command
	LedgerIndex  interface{}   `json:"ledger_index"`
	Accounts     bool          `json:"accounts"`
	Transactions bool          `json:"transactions"`
	Expand       bool          `json:"expand"`
	Result       *LedgerResult `json:"result,omitempty"`
}

type LedgerResult struct {
	Ledger data.Ledger
}

type LedgerHeaderCommand struct {
	*Command
	Ledger interface{} `json:"ledger"`
	Result *LedgerHeaderResult
}

type LedgerHeaderResult struct {
	Ledger         data.Ledger
	LedgerSequence uint32              `json:"ledger_index"`
	Hash           *data.Hash256       `json:"ledger_hash,omitempty"`
	LedgerData     data.VariableLength `json:"ledger_data"`
}

type LedgerDataCommand struct {
	*Command
	Ledger interface{}       `json:"ledger"`
	Marker *data.Hash256     `json:"marker,omitempty"`
	Result *LedgerDataResult `json:"result,omitempty"`
}

type BinaryLedgerDataCommand struct {
	*Command
	Ledger interface{}             `json:"ledger"`
	Binary bool                    `json:"binary"`
	Marker *data.Hash256           `json:"marker,omitempty"`
	Result *BinaryLedgerDataResult `json:"result,omitempty"`
}

type LedgerDataResult struct {
	LedgerSequence uint32                `json:"ledger_index"`
	Hash           data.Hash256          `json:"ledger_hash"`
	Marker         *data.Hash256         `json:"marker"`
	State          data.LedgerEntrySlice `json:"state"`
}

type BinaryLedgerData struct {
	Data  string `json:"data"`
	Index string `json:"index"`
}

type BinaryLedgerDataResult struct {
	LedgerSequence uint32             `json:"ledger_index"`
	Hash           data.Hash256       `json:"ledger_hash"`
	Marker         *data.Hash256      `json:"marker"`
	State          []BinaryLedgerData `json:"state"`
}

type RipplePathFindCommand struct {
	*Command
	SrcAccount    data.Account          `json:"source_account"`
	SrcCurrencies *[]data.Currency      `json:"source_currencies,omitempty"`
	DestAccount   data.Account          `json:"destination_account"`
	DestAmount    data.Amount           `json:"destination_amount"`
	Result        *RipplePathFindResult `json:"result,omitempty"`
}

type RipplePathFindResult struct {
	Alternatives []struct {
		SrcAmount      data.Amount  `json:"source_amount"`
		PathsComputed  data.PathSet `json:"paths_computed,omitempty"`
		PathsCanonical data.PathSet `json:"paths_canonical,omitempty"`
	}
	DestAccount    data.Account    `json:"destination_account"`
	DestCurrencies []data.Currency `json:"destination_currencies"`
}

type AccountInfoCommand struct {
	*Command
	Account     data.Account       `json:"account"`
	LedgerIndex interface{}        `json:"ledger_index,omitempty"`
	Result      *AccountInfoResult `json:"result,omitempty"`
}

type AccountInfoResult struct {
	LedgerSequence uint32           `json:"ledger_current_index"`
	AccountData    data.AccountRoot `json:"account_data"`
}

type AccountLinesCommand struct {
	*Command
	Account     data.Account        `json:"account"`
	Limit       uint32              `json:"limit"`
	LedgerIndex interface{}         `json:"ledger_index,omitempty"`
	Marker      *data.Hash256       `json:"marker,omitempty"`
	Result      *AccountLinesResult `json:"result,omitempty"`
}

type AccountLinesResult struct {
	LedgerSequence *uint32               `json:"ledger_index"`
	Account        data.Account          `json:"account"`
	Marker         *data.Hash256         `json:"marker"`
	Lines          data.AccountLineSlice `json:"lines"`
}

type AccountOffersCommand struct {
	*Command
	Account     data.Account         `json:"account"`
	Limit       uint32               `json:"limit"`
	LedgerIndex interface{}          `json:"ledger_index,omitempty"`
	Marker      *data.Hash256        `json:"marker,omitempty"`
	Result      *AccountOffersResult `json:"result,omitempty"`
}

type AccountOffersResult struct {
	LedgerSequence *uint32                `json:"ledger_index"`
	Account        data.Account           `json:"account"`
	Marker         *data.Hash256          `json:"marker"`
	Offers         data.AccountOfferSlice `json:"offers"`
}

type BookOffersCommand struct {
	*Command
	LedgerIndex interface{}  `json:"ledger_index,omitempty"`
	Taker       data.Account `json:"taker"`
	TakerPays   data.Asset   `json:"taker_pays"`
	TakerGets   data.Asset   `json:"taker_gets"`
	Limit       uint32       `json:"limit"`
	Result      *BookOffersResult
}

type BookOffersResult struct {
	LedgerSequence uint32                `json:"ledger_index"`
	Offers         []data.OrderBookOffer `json:"offers"`
}

type FeeCommand struct {
	*Command
	Result *FeeResult
}

type FeeResult struct {
	CurrentLedgerSize uint32 `json:"current_ledger_size,string"`
	CurrentQueueSize  uint32 `json:"current_queue_size,string"`
	Drops             struct {
		BaseFee       data.Value `json:"base_fee"`
		MedianFee     data.Value `json:"median_fee"`
		MinimumFee    data.Value `json:"minimum_fee"`
		OpenLedgerFee data.Value `json:"open_ledger_fee"`
	} `json:"drops"`
	ExpectedLedgerSize uint32 `json:"expected_ledger_size,string"`
	Levels             struct {
		MedianLevel     data.Value `json:"median_level"`
		MinimumLevel    data.Value `json:"minimum_level"`
		OpenLedgerLevel data.Value `json:"open_ledger_level"`
		ReferenceLevel  data.Value `json:"reference_level"`
	} `json:"levels"`
	MaxQueueSize uint32 `json:"max_queue_size,string"`
	Status       string `json:"status"`
}

type ServerStateCommand struct {
	*Command
	Result *ServerStateResult `json:"result,omitempty"`
}

type ServerStateResult struct {
	State struct {
		BuildVersion          string `json:"build_version"`
		CompleteLedgers       string `json:"complete_ledgers"`
		InitialSyncDurationUs string `json:"initial_sync_duration_us"`
		IoLatencyMs           int    `json:"io_latency_ms"`
		JqTransOverflow       string `json:"jq_trans_overflow"`
		LastClose             struct {
			ConvergeTime int `json:"converge_time"`
			Proposers    int `json:"proposers"`
		} `json:"last_close"`
		LoadBase                 int    `json:"load_base"`
		LoadFactor               int    `json:"load_factor"`
		LoadFactorFeeEscalation  int    `json:"load_factor_fee_escalation"`
		LoadFactorFeeQueue       int    `json:"load_factor_fee_queue"`
		LoadFactorFeeReference   int    `json:"load_factor_fee_reference"`
		LoadFactorServer         int    `json:"load_factor_server"`
		NetworkId                int    `json:"network_id"`
		PeerDisconnects          string `json:"peer_disconnects"`
		PeerDisconnectsResources string `json:"peer_disconnects_resources"`
		Peers                    int    `json:"peers"`
		Ports                    []struct {
			Port     string   `json:"port"`
			Protocol []string `json:"protocol"`
		} `json:"ports"`
		PubkeyNode            string `json:"pubkey_node"`
		ServerState           string `json:"server_state"`
		ServerStateDurationUs string `json:"server_state_duration_us"`
		StateAccounting       struct {
			Connected struct {
				DurationUs  string `json:"duration_us"`
				Transitions string `json:"transitions"`
			} `json:"connected"`
			Disconnected struct {
				DurationUs  string `json:"duration_us"`
				Transitions string `json:"transitions"`
			} `json:"disconnected"`
			Full struct {
				DurationUs  string `json:"duration_us"`
				Transitions string `json:"transitions"`
			} `json:"full"`
			Syncing struct {
				DurationUs  string `json:"duration_us"`
				Transitions string `json:"transitions"`
			} `json:"syncing"`
			Tracking struct {
				DurationUs  string `json:"duration_us"`
				Transitions string `json:"transitions"`
			} `json:"tracking"`
		} `json:"state_accounting"`
		Time            string `json:"time"`
		Uptime          int    `json:"uptime"`
		ValidatedLedger struct {
			BaseFee     int    `json:"base_fee"`
			CloseTime   int    `json:"close_time"`
			Hash        string `json:"hash"`
			ReserveBase int    `json:"reserve_base"`
			ReserveInc  int    `json:"reserve_inc"`
			Seq         int    `json:"seq"`
		} `json:"validated_ledger"`
		ValidationQuorum int `json:"validation_quorum"`
	} `json:"state"`
}
