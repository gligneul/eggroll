// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package eggeth

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// ERC1155BatchPortalMetaData contains all meta data concerning the ERC1155BatchPortal contract.
var ERC1155BatchPortalMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"contractIInputBox\",\"name\":\"_inputBox\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"contractIERC1155\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_dapp\",\"type\":\"address\"},{\"internalType\":\"uint256[]\",\"name\":\"_tokenIds\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"_values\",\"type\":\"uint256[]\"},{\"internalType\":\"bytes\",\"name\":\"_baseLayerData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"_execLayerData\",\"type\":\"bytes\"}],\"name\":\"depositBatchERC1155Token\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getInputBox\",\"outputs\":[{\"internalType\":\"contractIInputBox\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// ERC1155BatchPortalABI is the input ABI used to generate the binding from.
// Deprecated: Use ERC1155BatchPortalMetaData.ABI instead.
var ERC1155BatchPortalABI = ERC1155BatchPortalMetaData.ABI

// ERC1155BatchPortal is an auto generated Go binding around an Ethereum contract.
type ERC1155BatchPortal struct {
	ERC1155BatchPortalCaller     // Read-only binding to the contract
	ERC1155BatchPortalTransactor // Write-only binding to the contract
	ERC1155BatchPortalFilterer   // Log filterer for contract events
}

// ERC1155BatchPortalCaller is an auto generated read-only Go binding around an Ethereum contract.
type ERC1155BatchPortalCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC1155BatchPortalTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ERC1155BatchPortalTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC1155BatchPortalFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ERC1155BatchPortalFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC1155BatchPortalSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ERC1155BatchPortalSession struct {
	Contract     *ERC1155BatchPortal // Generic contract binding to set the session for
	CallOpts     bind.CallOpts       // Call options to use throughout this session
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// ERC1155BatchPortalCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ERC1155BatchPortalCallerSession struct {
	Contract *ERC1155BatchPortalCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts             // Call options to use throughout this session
}

// ERC1155BatchPortalTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ERC1155BatchPortalTransactorSession struct {
	Contract     *ERC1155BatchPortalTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts             // Transaction auth options to use throughout this session
}

// ERC1155BatchPortalRaw is an auto generated low-level Go binding around an Ethereum contract.
type ERC1155BatchPortalRaw struct {
	Contract *ERC1155BatchPortal // Generic contract binding to access the raw methods on
}

// ERC1155BatchPortalCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ERC1155BatchPortalCallerRaw struct {
	Contract *ERC1155BatchPortalCaller // Generic read-only contract binding to access the raw methods on
}

// ERC1155BatchPortalTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ERC1155BatchPortalTransactorRaw struct {
	Contract *ERC1155BatchPortalTransactor // Generic write-only contract binding to access the raw methods on
}

// NewERC1155BatchPortal creates a new instance of ERC1155BatchPortal, bound to a specific deployed contract.
func NewERC1155BatchPortal(address common.Address, backend bind.ContractBackend) (*ERC1155BatchPortal, error) {
	contract, err := bindERC1155BatchPortal(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ERC1155BatchPortal{ERC1155BatchPortalCaller: ERC1155BatchPortalCaller{contract: contract}, ERC1155BatchPortalTransactor: ERC1155BatchPortalTransactor{contract: contract}, ERC1155BatchPortalFilterer: ERC1155BatchPortalFilterer{contract: contract}}, nil
}

// NewERC1155BatchPortalCaller creates a new read-only instance of ERC1155BatchPortal, bound to a specific deployed contract.
func NewERC1155BatchPortalCaller(address common.Address, caller bind.ContractCaller) (*ERC1155BatchPortalCaller, error) {
	contract, err := bindERC1155BatchPortal(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ERC1155BatchPortalCaller{contract: contract}, nil
}

// NewERC1155BatchPortalTransactor creates a new write-only instance of ERC1155BatchPortal, bound to a specific deployed contract.
func NewERC1155BatchPortalTransactor(address common.Address, transactor bind.ContractTransactor) (*ERC1155BatchPortalTransactor, error) {
	contract, err := bindERC1155BatchPortal(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ERC1155BatchPortalTransactor{contract: contract}, nil
}

// NewERC1155BatchPortalFilterer creates a new log filterer instance of ERC1155BatchPortal, bound to a specific deployed contract.
func NewERC1155BatchPortalFilterer(address common.Address, filterer bind.ContractFilterer) (*ERC1155BatchPortalFilterer, error) {
	contract, err := bindERC1155BatchPortal(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ERC1155BatchPortalFilterer{contract: contract}, nil
}

// bindERC1155BatchPortal binds a generic wrapper to an already deployed contract.
func bindERC1155BatchPortal(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ERC1155BatchPortalMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ERC1155BatchPortal *ERC1155BatchPortalRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ERC1155BatchPortal.Contract.ERC1155BatchPortalCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ERC1155BatchPortal *ERC1155BatchPortalRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ERC1155BatchPortal.Contract.ERC1155BatchPortalTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ERC1155BatchPortal *ERC1155BatchPortalRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ERC1155BatchPortal.Contract.ERC1155BatchPortalTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ERC1155BatchPortal *ERC1155BatchPortalCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ERC1155BatchPortal.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ERC1155BatchPortal *ERC1155BatchPortalTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ERC1155BatchPortal.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ERC1155BatchPortal *ERC1155BatchPortalTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ERC1155BatchPortal.Contract.contract.Transact(opts, method, params...)
}

// GetInputBox is a free data retrieval call binding the contract method 0x00aace9a.
//
// Solidity: function getInputBox() view returns(address)
func (_ERC1155BatchPortal *ERC1155BatchPortalCaller) GetInputBox(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ERC1155BatchPortal.contract.Call(opts, &out, "getInputBox")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetInputBox is a free data retrieval call binding the contract method 0x00aace9a.
//
// Solidity: function getInputBox() view returns(address)
func (_ERC1155BatchPortal *ERC1155BatchPortalSession) GetInputBox() (common.Address, error) {
	return _ERC1155BatchPortal.Contract.GetInputBox(&_ERC1155BatchPortal.CallOpts)
}

// GetInputBox is a free data retrieval call binding the contract method 0x00aace9a.
//
// Solidity: function getInputBox() view returns(address)
func (_ERC1155BatchPortal *ERC1155BatchPortalCallerSession) GetInputBox() (common.Address, error) {
	return _ERC1155BatchPortal.Contract.GetInputBox(&_ERC1155BatchPortal.CallOpts)
}

// DepositBatchERC1155Token is a paid mutator transaction binding the contract method 0x24d15c67.
//
// Solidity: function depositBatchERC1155Token(address _token, address _dapp, uint256[] _tokenIds, uint256[] _values, bytes _baseLayerData, bytes _execLayerData) returns()
func (_ERC1155BatchPortal *ERC1155BatchPortalTransactor) DepositBatchERC1155Token(opts *bind.TransactOpts, _token common.Address, _dapp common.Address, _tokenIds []*big.Int, _values []*big.Int, _baseLayerData []byte, _execLayerData []byte) (*types.Transaction, error) {
	return _ERC1155BatchPortal.contract.Transact(opts, "depositBatchERC1155Token", _token, _dapp, _tokenIds, _values, _baseLayerData, _execLayerData)
}

// DepositBatchERC1155Token is a paid mutator transaction binding the contract method 0x24d15c67.
//
// Solidity: function depositBatchERC1155Token(address _token, address _dapp, uint256[] _tokenIds, uint256[] _values, bytes _baseLayerData, bytes _execLayerData) returns()
func (_ERC1155BatchPortal *ERC1155BatchPortalSession) DepositBatchERC1155Token(_token common.Address, _dapp common.Address, _tokenIds []*big.Int, _values []*big.Int, _baseLayerData []byte, _execLayerData []byte) (*types.Transaction, error) {
	return _ERC1155BatchPortal.Contract.DepositBatchERC1155Token(&_ERC1155BatchPortal.TransactOpts, _token, _dapp, _tokenIds, _values, _baseLayerData, _execLayerData)
}

// DepositBatchERC1155Token is a paid mutator transaction binding the contract method 0x24d15c67.
//
// Solidity: function depositBatchERC1155Token(address _token, address _dapp, uint256[] _tokenIds, uint256[] _values, bytes _baseLayerData, bytes _execLayerData) returns()
func (_ERC1155BatchPortal *ERC1155BatchPortalTransactorSession) DepositBatchERC1155Token(_token common.Address, _dapp common.Address, _tokenIds []*big.Int, _values []*big.Int, _baseLayerData []byte, _execLayerData []byte) (*types.Transaction, error) {
	return _ERC1155BatchPortal.Contract.DepositBatchERC1155Token(&_ERC1155BatchPortal.TransactOpts, _token, _dapp, _tokenIds, _values, _baseLayerData, _execLayerData)
}
