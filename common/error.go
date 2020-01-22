package common

import "errors"

//common errors
var (
	ERR_INVALID_ADDRESS       	= errors.New("Incorrect Address Format")
	ERR_ADDRESS_LEN 		  	= errors.New("Address Len Invalid")
	ERR_INVALID_PRIVATEKEY   	= errors.New("privateKey invalid")
	ERR_INVALID_PRIVATEKEY_LEN  = errors.New("privateKey Len invalid")
)

// WaykiChain errors
var (
	ERR_INVALID_MNEMONIC      = errors.New("Invalid Mnemonic")
	ERR_INVALID_NETWORK       = errors.New("Invalid Network type")
	ERR_NEGATIVE_VALID_HEIGHT = errors.New("ValidHeight can not be negative")
	ERR_INVALID_SRC_REG_ID    = errors.New("SrcRegId must be a valid RegID")
	ERR_INVALID_DEST_ADDR     = errors.New("DestAddr must be a valid RegID or Address")
	ERR_RANGE_VALUES          = errors.New("Values out of range")
	ERR_RANGE_FEE             = errors.New("Fees out of range")
	ERR_FEE_SMALLER_MIN       = errors.New("Fees smaller than MinTxFee")
	ERR_EMPTY_VOTES           = errors.New("Votes can be not empty")
	ERR_INVALID_VOTE_PUBKEY   = errors.New("Vote PubKey invalid, PubKey len must equal 33")
	ERR_RANGE_VOTE_VALUE      = errors.New("VoteValue out of range")
	ERR_INVALID_APP_ID        = errors.New("AppId must be a valid RegID")
	ERR_INVALID_CONTRACT_HEX  = errors.New("ContractHex must be valid hex format")
	ERR_INVALID_SCRIPT        = errors.New("Script can not be empty or is too large")
	ERR_INVALID_SCRIPT_DESC   = errors.New("Description of script is too large")

	ERR_CDP_TX_HASH      = errors.New("CDP tx hash error")
	ERR_CDP_STAKE_NUMBER = errors.New("CDP stake number error")
	ERR_COIN_TYPE        = errors.New("Coin type error")
	ERR_USER_PUBLICKEY   = errors.New("PublicKey invalid")

	ERR_ASK_PRICE   = errors.New("Ask Price invalid")
	ERR_SIGNATURE_ERROR       = errors.New("Signature error")
	ERR_SYMBOL_ERROR       = errors.New("Symbol Capital letter A-Z 6-7 digits [A_Z] error")
	ERR_ASSET_NAME_ERROR       = errors.New("Asset Name error")
	ERR_TOTAl_SUPPLY_ERROR       = errors.New("Asset Total Supply error")
	ERR_ASSET_UPDATE_TYPE_ERROR       = errors.New("Asset Update Type error")
	ERR_ASSET_UPDATE_OWNER_ERROR       = errors.New("Asset Update Owner error")

	ERR_PUBLICKEY_SIGNATURE_ERROR       = errors.New("The Length of PublicKey or Signature error")
)

// ETH errors
var (

)