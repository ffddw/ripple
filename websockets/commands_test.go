package websockets

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/ffddw/ripple/data"
	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type MessagesSuite struct{}

var _ = Suite(&MessagesSuite{})

func readResponseFile(c *C, msg interface{}, path string) {
	b, err := os.ReadFile(path)
	if err != nil {
		c.Error(err)
	}

	if err = json.Unmarshal(b, msg); err != nil {
		c.Error(err)
	}
}

func (s *MessagesSuite) TestLedgerResponse(c *C) {
	msg := &LedgerCommand{}
	readResponseFile(c, msg, "testdata/ledger.json")

	// Response fields
	c.Assert(msg.Status, Equals, "success")
	c.Assert(msg.Type, Equals, "response")

	// Result fields
	c.Assert(msg.Result.Ledger.LedgerSequence, Equals, uint32(6917762))
	c.Assert(msg.Result.Ledger.Accepted, Equals, true)
	c.Assert(msg.Result.Ledger.CloseTime.String(), Equals, "2014-May-30 13:11:50 UTC")
	c.Assert(msg.Result.Ledger.Closed, Equals, true)
	c.Assert(msg.Result.Ledger.Hash.String(), Equals, "0C5C5B39EA40D40ACA6EB47E50B2B85FD516D1A2BA67BA3E050349D3EF3632A4")
	c.Assert(msg.Result.Ledger.PreviousLedger.String(), Equals, "F8F0363803C30E659AA24D6A62A6512BA24BEA5AC52A29731ABA1E2D80796E8B")
	c.Assert(msg.Result.Ledger.TotalXRP, Equals, uint64(99999990098968782))
	c.Assert(msg.Result.Ledger.StateHash.String(), Equals, "46D3E36FE845B9A18293F4C0F134D7DAFB06D4D9A1C7E4CB03F8B293CCA45FA0")
	c.Assert(msg.Result.Ledger.TransactionHash.String(), Equals, "757CCB586D44F3C58E366EC7618988C0596277D3D5D0B412E49563B5EEDF04FF")

	c.Assert(msg.Result.Ledger.Transactions, HasLen, 7)
	tx0 := msg.Result.Ledger.Transactions[0]
	c.Assert(tx0.GetHash().String(), Equals, "2D0CE11154B655A2BFE7F3F857AAC344622EC7DAB11B1EBD920DCDB00E8646FF")
	c.Assert(tx0.MetaData.AffectedNodes, HasLen, 4)
}

func (s *MessagesSuite) TestLedgerHeaderResponse(c *C) {
	msg := &LedgerHeaderCommand{}
	readResponseFile(c, msg, "testdata/ledger_header.json")

	// Response fields
	c.Assert(msg.Status, Equals, "success")
	c.Assert(msg.Type, Equals, "response")

	// Result fields
	c.Assert(len(msg.Result.LedgerData), Equals, 118)
	c.Assert(msg.Result.LedgerSequence, Equals, uint32(32570))
	c.Assert(msg.Result.Ledger.LedgerSequence, Equals, uint32(32570))
	c.Assert(msg.Result.Ledger.Accepted, Equals, true)
	c.Assert(msg.Result.Ledger.CloseTime.String(), Equals, "2013-Jan-01 03:21:10 UTC")
	c.Assert(msg.Result.Ledger.Closed, Equals, true)
	c.Assert(msg.Result.Ledger.Hash.String(), Equals, "4109C6F2045FC7EFF4CDE8F9905D19C28820D86304080FF886B299F0206E42B5")
	c.Assert(msg.Result.Ledger.PreviousLedger.String(), Equals, "60A01EBF11537D8394EA1235253293508BDA7131D5F8710EFE9413AA129653A2")
	c.Assert(msg.Result.Ledger.TotalXRP, Equals, uint64(99999999999996320))
	c.Assert(msg.Result.Ledger.StateHash.String(), Equals, "3806AF8F22037DE598D30D38C8861FADF391171D26F7DE34ACFA038996EA6BEB")
	c.Assert(msg.Result.Ledger.TransactionHash.String(), Equals, "0000000000000000000000000000000000000000000000000000000000000000")
	c.Assert(msg.Result.Ledger.Transactions, HasLen, 0)
}

func (s *MessagesSuite) TestTxResponse(c *C) {
	msg := &TxCommand{}
	readResponseFile(c, msg, "testdata/tx.json")

	// Response fields
	c.Assert(msg.Status, Equals, "success")
	c.Assert(msg.Type, Equals, "response")

	// Result fields
	c.Assert(msg.Result.Date.String(), Equals, "2014-May-30 13:11:50 UTC")
	c.Assert(msg.Result.Validated, Equals, true)
	c.Assert(msg.Result.MetaData.AffectedNodes, HasLen, 4)
	c.Assert(msg.Result.MetaData.TransactionResult.String(), Equals, "tesSUCCESS")

	offer := msg.Result.Transaction.(*data.OfferCreate)
	c.Assert(msg.Result.GetHash().String(), Equals, "2D0CE11154B655A2BFE7F3F857AAC344622EC7DAB11B1EBD920DCDB00E8646FF")
	c.Assert(offer.GetType(), Equals, "OfferCreate")
	c.Assert(offer.Account.String(), Equals, "rwpxNWdpKu2QVgrh5LQXEygYLshhgnRL1Y")
	c.Assert(offer.Fee.String(), Equals, "0.00001")
	c.Assert(offer.SigningPubKey.String(), Equals, "02BD6F0CFD0182F2F408512286A0D935C58FF41169DAC7E721D159D711695DFF85")
	c.Assert(offer.TxnSignature.String(), Equals, "30440220216D42DF672C1CC7EF0CA9C7840838A2AF5FEDD4DEFCBA770C763D7509703C8702203C8D831BFF8A8BC2CC993BECB4E6C7BE1EA9D394AB7CE7C6F7542B6CDA781467")
	c.Assert(offer.Sequence, Equals, uint32(1681497))
}

func (s *MessagesSuite) TestTxResponseOracleSet(c *C) {
	msg := &TxCommand{}
	readResponseFile(c, msg, "testdata/tx_oracle_set.json")

	// Response fields
	c.Assert(msg.Status, Equals, "success")
	c.Assert(msg.Type, Equals, "response")

	// Result fields
	c.Assert(msg.Result.Date.String(), Equals, "2025-Feb-11 05:35:32 UTC")
	c.Assert(msg.Result.Validated, Equals, true)
	c.Assert(msg.Result.MetaData.AffectedNodes, HasLen, 2)
	c.Assert(msg.Result.MetaData.TransactionResult.String(), Equals, "tesSUCCESS")

	oracle := msg.Result.Transaction.(*data.OracleSet)
	c.Assert(msg.Result.GetHash().String(), Equals, "940E8BE41222353F2971C90B572BBC4549F5B823E1F056A13C43055407B4F77D")
	c.Assert(oracle.GetType(), Equals, "OracleSet")
	c.Assert(oracle.Account.String(), Equals, "rrJPYwVRyWFcwfaNMm83QEaCexEpKnkEg")
	c.Assert(oracle.Fee.String(), Equals, "0.000012")
	c.Assert(oracle.SigningPubKey.String(), Equals, "039971FA8F08E76B59DE4D4525451C2200E585D0CC5AC52C8E9EBB753B2FC6A6C9")
	c.Assert(oracle.TxnSignature.String(), Equals, "3044022034CBED2017EC8938EE60CA44C039BF1103C1EE46A57BDFE6FF7A8E3199DD427D02200C07CA8CE664649BA7FE735A505046E7FED19AC61AEFD91D226FEFC97EC39A3D")
	c.Assert(oracle.Sequence, Equals, uint32(69611178))
}

func (s *MessagesSuite) TestTxResponseOracleDelete(c *C) {
	msg := &TxCommand{}
	readResponseFile(c, msg, "testdata/tx_oracle_delete.json")

	// Response fields
	c.Assert(msg.Status, Equals, "success")
	c.Assert(msg.Type, Equals, "response")

	// Result fields
	c.Assert(msg.Result.Date.String(), Equals, "2025-Feb-11 05:53:11 UTC")
	c.Assert(msg.Result.Validated, Equals, true)
	c.Assert(msg.Result.MetaData.AffectedNodes, HasLen, 3)
	c.Assert(msg.Result.MetaData.TransactionResult.String(), Equals, "tesSUCCESS")

	oracle := msg.Result.Transaction.(*data.OracleDelete)
	c.Assert(msg.Result.GetHash().String(), Equals, "C5E277B47209744103E13E5C6CAC57E9F7FC2E8119EC8309E50C846DB852CC7D")
	c.Assert(oracle.GetType(), Equals, "OracleDelete")
	c.Assert(oracle.Account.String(), Equals, "roosteri9aGNFRXZrJNYQKVBfxHiE5abg")
	c.Assert(oracle.Fee.String(), Equals, "0.000012")
	c.Assert(oracle.SigningPubKey.String(), Equals, "03E5B78E49930BCC391F94C64875879E8BED4A924C9110C7A200F8827544957A92")
	c.Assert(oracle.TxnSignature.String(), Equals, "304502210089B24B7A836EE3CB87F9065851112FAF7DCEBA7B8487893FBE7B9CFAAAACB6B302206EF3912B1F2F2282DC6210BBDE35AF04C935BA5697AAC7EE5ADE32B0ADFF4C44")
	c.Assert(oracle.Sequence, Equals, uint32(92545197))
}

func (s *MessagesSuite) TestAccountTxResponse(c *C) {
	msg := &AccountTxCommand{}
	readResponseFile(c, msg, "testdata/account_tx.json")

	// Response fields
	c.Assert(msg.Status, Equals, "success")
	c.Assert(msg.Type, Equals, "response")

	c.Assert(len(msg.Result.Transactions), Equals, 2)
	c.Assert(msg.Result.Transactions[1].Date.String(), Equals, "2014-Jun-19 14:14:40 UTC")
	offer := msg.Result.Transactions[1].Transaction.(*data.OfferCreate)
	c.Assert(offer.TakerPays.String(), Equals, "0.034800328/BTC/rvYAfWj5gh67oV6fW32ZzP3Aw4Eubs59B")
}

func (s *MessagesSuite) TestLedgerDataResponse(c *C) {
	msg := &LedgerDataCommand{}
	readResponseFile(c, msg, "testdata/ledger_data.json")

	// Response fields
	c.Assert(msg.Status, Equals, "success")
	c.Assert(msg.Type, Equals, "response")

	c.Assert(msg.Result.LedgerSequence, Equals, uint32(6281820))
	c.Assert(msg.Result.Hash.String(), Equals, "83CC350B1CDD9792D47F60D3DBB7673518FD6E71821070673E6EAE65DE69086B")
	c.Assert(msg.Result.Marker.String(), Equals, "02DE1A2AD4332A1AF01C59F16E45218FA70E5792BD963B6D7ACF188D6D150607")
	c.Assert(len(msg.Result.State), Equals, 2048)
	c.Assert(msg.Result.State[0].GetType(), Equals, "AccountRoot")
	c.Assert(msg.Result.State[0].GetLedgerIndex().String(), Equals, "00001A2969BE1FC85F1D7A55282FA2E6D95C71D2E4B9C0FDD3D9994F3C00FF8F")
}

func (s *MessagesSuite) TestRipplePathFindResponse(c *C) {
	msg := &RipplePathFindCommand{}
	readResponseFile(c, msg, "testdata/ripple_path_find.json")

	// Response fields
	c.Assert(msg.Status, Equals, "success")
	c.Assert(msg.Type, Equals, "response")

	c.Assert(msg.Result.DestAccount.String(), Equals, "r9Dr5xwkeLegBeXq6ujinjSBLQzQ1zQGjH")
	c.Assert(msg.Result.DestCurrencies, HasLen, 6)
	c.Assert(msg.Result.Alternatives, HasLen, 1)
	c.Assert(msg.Result.Alternatives[0].SrcAmount.String(), Equals, "0.9940475268/USD/rvYAfWj5gh67oV6fW32ZzP3Aw4Eubs59B")
	c.Assert(msg.Result.Alternatives[0].PathsComputed, HasLen, 4)
	c.Assert(msg.Result.Alternatives[0].PathsCanonical, HasLen, 0)
	c.Assert(msg.Result.Alternatives[0].PathsComputed[0], HasLen, 2)
	c.Assert(msg.Result.Alternatives[0].PathsComputed[0].String(), Equals, "XRP => SGD/r9Dr5xwkeLegBeXq6ujinjSBLQzQ1zQGjH")
}

func (s *MessagesSuite) TestAccountInfoResponse(c *C) {
	msg := &AccountInfoCommand{}
	readResponseFile(c, msg, "testdata/account_info.json")

	// Response fields
	c.Assert(msg.Status, Equals, "success")
	c.Assert(msg.Type, Equals, "response")

	c.Assert(msg.Result.LedgerSequence, Equals, uint32(7636529))
	c.Assert(*msg.Result.AccountData.TransferRate, Equals, uint32(1002000000))
	c.Assert(msg.Result.AccountData.LedgerEntryType, Equals, data.ACCOUNT_ROOT)
	c.Assert(*msg.Result.AccountData.Sequence, Equals, uint32(546))
	c.Assert(msg.Result.AccountData.Balance.String(), Equals, "10321199.422233")
}

func (s *MessagesSuite) TestServerStateResponse(c *C) {
	msg := &ServerStateCommand{}
	readResponseFile(c, msg, "testdata/server_state.json")

	// Response fields
	c.Assert(msg.Status, Equals, "success")
	c.Assert(msg.Type, Equals, "response")

	c.Assert(msg.Result.State.CompleteLedgers, Equals, "93610338-93958930")
	c.Assert(msg.Result.State.NetworkId, Equals, 0)
	c.Assert(msg.Result.State.PubkeyNode, Equals, "n9LRU27FWoEwEhxKLeJKanjZq7ekfdx3dTQhnjnmuVnuubYvLSwh")
	c.Assert(msg.Result.State.ServerState, Equals, "full")
	c.Assert(msg.Result.State.Ports, HasLen, 2)
	c.Assert(msg.Result.State.Ports[1].Protocol[0], Equals, "peer")
	c.Assert(msg.Result.State.StateAccounting.Connected.DurationUs, Equals, "1086992087")
	c.Assert(msg.Result.State.StateAccounting.Disconnected.DurationUs, Equals, "1056155")
	c.Assert(msg.Result.State.Time, Equals, "2025-Feb-06 09:10:45.542639 UTC")
	c.Assert(msg.Result.State.Uptime, Equals, 270655)
	c.Assert(msg.Result.State.ValidatedLedger.CloseTime, Equals, 792148241)
	c.Assert(msg.Result.State.ValidatedLedger.Hash, Equals, "0A6DFE7685BB62B5D434DDF8939DA3C617C6988BC9BF65409DB4BE035A7C1001")
}

func (s *MessagesSuite) TestSubmitResultResponse(c *C) {
	msg := &SubmitCommand{}
	readResponseFile(c, msg, "testdata/submit_result.json")

	// Response fields
	c.Assert(msg.Status, Equals, "success")
	c.Assert(msg.Type, Equals, "response")

	c.Assert(msg.Result.Tx.Account, Equals, "r9cZA1mLK5R5Am25ArfXFmqgNwjZgnfk59")
	c.Assert(msg.Result.Tx.Domain, Equals, "726970706C652E636F6D")
	c.Assert(msg.Result.Tx.Fee, Equals, "12")
	c.Assert(msg.Result.Tx.Flags, Equals, uint32(2147483648))
	c.Assert(msg.Result.Tx.LastLedgerSequence, Equals, uint32(8820051))
	c.Assert(msg.Result.Tx.Sequence, Equals, uint32(23))
	c.Assert(msg.Result.Tx.SigningPubKey, Equals, "02F89EAEC7667B30F33D0687BBA86C3FE2A08CCA40A9186C5BDE2DAA6FA97A37D8")
	c.Assert(msg.Result.Tx.TransactionType, Equals, "AccountSet")
	c.Assert(msg.Result.Tx.TxnSignature, Equals, "3045022100BDE09A1F6670403F341C21A77CF35BA47E45CDE974096E1AA5FC39811D8269E702203D60291B9A27F1DCABA9CF5DED307B4F23223E0B6F156991DB601DFB9C41CE1C")
	c.Assert(msg.Result.Tx.Hash, Equals, "02ACE87F1996E3A23690A5BB7F1774BF71CCBA68F79805831B42ABAD5913D6F4")
}
