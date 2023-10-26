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

// ERC1155SinglePortalMetaData contains all meta data concerning the ERC1155SinglePortal contract.
var ERC1155SinglePortalMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"contractIInputBox\",\"name\":\"_inputBox\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"contractIERC1155\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_dapp\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_tokenId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_value\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"_baseLayerData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"_execLayerData\",\"type\":\"bytes\"}],\"name\":\"depositSingleERC1155Token\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getInputBox\",\"outputs\":[{\"internalType\":\"contractIInputBox\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// ERC1155SinglePortalABI is the input ABI used to generate the binding from.
// Deprecated: Use ERC1155SinglePortalMetaData.ABI instead.
var ERC1155SinglePortalABI = ERC1155SinglePortalMetaData.ABI

// ERC1155SinglePortal is an auto generated Go binding around an Ethereum contract.
type ERC1155SinglePortal struct {
	ERC1155SinglePortalCaller     // Read-only binding to the contract
	ERC1155SinglePortalTransactor // Write-only binding to the contract
	ERC1155SinglePortalFilterer   // Log filterer for contract events
}

// ERC1155SinglePortalCaller is an auto generated read-only Go binding around an Ethereum contract.
type ERC1155SinglePortalCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC1155SinglePortalTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ERC1155SinglePortalTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC1155SinglePortalFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ERC1155SinglePortalFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC1155SinglePortalSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ERC1155SinglePortalSession struct {
	Contract     *ERC1155SinglePortal // Generic contract binding to set the session for
	CallOpts     bind.CallOpts        // Call options to use throughout this session
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// ERC1155SinglePortalCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ERC1155SinglePortalCallerSession struct {
	Contract *ERC1155SinglePortalCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts              // Call options to use throughout this session
}

// ERC1155SinglePortalTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ERC1155SinglePortalTransactorSession struct {
	Contract     *ERC1155SinglePortalTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts              // Transaction auth options to use throughout this session
}

// ERC1155SinglePortalRaw is an auto generated low-level Go binding around an Ethereum contract.
type ERC1155SinglePortalRaw struct {
	Contract *ERC1155SinglePortal // Generic contract binding to access the raw methods on
}

// ERC1155SinglePortalCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ERC1155SinglePortalCallerRaw struct {
	Contract *ERC1155SinglePortalCaller // Generic read-only contract binding to access the raw methods on
}

// ERC1155SinglePortalTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ERC1155SinglePortalTransactorRaw struct {
	Contract *ERC1155SinglePortalTransactor // Generic write-only contract binding to access the raw methods on
}

// NewERC1155SinglePortal creates a new instance of ERC1155SinglePortal, bound to a specific deployed contract.
func NewERC1155SinglePortal(address common.Address, backend bind.ContractBackend) (*ERC1155SinglePortal, error) {
	contract, err := bindERC1155SinglePortal(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ERC1155SinglePortal{ERC1155SinglePortalCaller: ERC1155SinglePortalCaller{contract: contract}, ERC1155SinglePortalTransactor: ERC1155SinglePortalTransactor{contract: contract}, ERC1155SinglePortalFilterer: ERC1155SinglePortalFilterer{contract: contract}}, nil
}

// NewERC1155SinglePortalCaller creates a new read-only instance of ERC1155SinglePortal, bound to a specific deployed contract.
func NewERC1155SinglePortalCaller(address common.Address, caller bind.ContractCaller) (*ERC1155SinglePortalCaller, error) {
	contract, err := bindERC1155SinglePortal(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ERC1155SinglePortalCaller{contract: contract}, nil
}

// NewERC1155SinglePortalTransactor creates a new write-only instance of ERC1155SinglePortal, bound to a specific deployed contract.
func NewERC1155SinglePortalTransactor(address common.Address, transactor bind.ContractTransactor) (*ERC1155SinglePortalTransactor, error) {
	contract, err := bindERC1155SinglePortal(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ERC1155SinglePortalTransactor{contract: contract}, nil
}

// NewERC1155SinglePortalFilterer creates a new log filterer instance of ERC1155SinglePortal, bound to a specific deployed contract.
func NewERC1155SinglePortalFilterer(address common.Address, filterer bind.ContractFilterer) (*ERC1155SinglePortalFilterer, error) {
	contract, err := bindERC1155SinglePortal(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ERC1155SinglePortalFilterer{contract: contract}, nil
}

// bindERC1155SinglePortal binds a generic wrapper to an already deployed contract.
func bindERC1155SinglePortal(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ERC1155SinglePortalMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ERC1155SinglePortal *ERC1155SinglePortalRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ERC1155SinglePortal.Contract.ERC1155SinglePortalCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ERC1155SinglePortal *ERC1155SinglePortalRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ERC1155SinglePortal.Contract.ERC1155SinglePortalTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ERC1155SinglePortal *ERC1155SinglePortalRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ERC1155SinglePortal.Contract.ERC1155SinglePortalTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ERC1155SinglePortal *ERC1155SinglePortalCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ERC1155SinglePortal.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ERC1155SinglePortal *ERC1155SinglePortalTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ERC1155SinglePortal.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ERC1155SinglePortal *ERC1155SinglePortalTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ERC1155SinglePortal.Contract.contract.Transact(opts, method, params...)
}

// GetInputBox is a free data retrieval call binding the contract method 0x00aace9a.
//
// Solidity: function getInputBox() view returns(address)
func (_ERC1155SinglePortal *ERC1155SinglePortalCaller) GetInputBox(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ERC1155SinglePortal.contract.Call(opts, &out, "getInputBox")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetInputBox is a free data retrieval call binding the contract method 0x00aace9a.
//
// Solidity: function getInputBox() view returns(address)
func (_ERC1155SinglePortal *ERC1155SinglePortalSession) GetInputBox() (common.Address, error) {
	return _ERC1155SinglePortal.Contract.GetInputBox(&_ERC1155SinglePortal.CallOpts)
}

// GetInputBox is a free data retrieval call binding the contract method 0x00aace9a.
//
// Solidity: function getInputBox() view returns(address)
func (_ERC1155SinglePortal *ERC1155SinglePortalCallerSession) GetInputBox() (common.Address, error) {
	return _ERC1155SinglePortal.Contract.GetInputBox(&_ERC1155SinglePortal.CallOpts)
}

// DepositSingleERC1155Token is a paid mutator transaction binding the contract method 0xdec07dca.
//
// Solidity: function depositSingleERC1155Token(address _token, address _dapp, uint256 _tokenId, uint256 _value, bytes _baseLayerData, bytes _execLayerData) returns()
func (_ERC1155SinglePortal *ERC1155SinglePortalTransactor) DepositSingleERC1155Token(opts *bind.TransactOpts, _token common.Address, _dapp common.Address, _tokenId *big.Int, _value *big.Int, _baseLayerData []byte, _execLayerData []byte) (*types.Transaction, error) {
	return _ERC1155SinglePortal.contract.Transact(opts, "depositSingleERC1155Token", _token, _dapp, _tokenId, _value, _baseLayerData, _execLayerData)
}

// DepositSingleERC1155Token is a paid mutator transaction binding the contract method 0xdec07dca.
//
// Solidity: function depositSingleERC1155Token(address _token, address _dapp, uint256 _tokenId, uint256 _value, bytes _baseLayerData, bytes _execLayerData) returns()
func (_ERC1155SinglePortal *ERC1155SinglePortalSession) DepositSingleERC1155Token(_token common.Address, _dapp common.Address, _tokenId *big.Int, _value *big.Int, _baseLayerData []byte, _execLayerData []byte) (*types.Transaction, error) {
	return _ERC1155SinglePortal.Contract.DepositSingleERC1155Token(&_ERC1155SinglePortal.TransactOpts, _token, _dapp, _tokenId, _value, _baseLayerData, _execLayerData)
}

// DepositSingleERC1155Token is a paid mutator transaction binding the contract method 0xdec07dca.
//
// Solidity: function depositSingleERC1155Token(address _token, address _dapp, uint256 _tokenId, uint256 _value, bytes _baseLayerData, bytes _execLayerData) returns()
func (_ERC1155SinglePortal *ERC1155SinglePortalTransactorSession) DepositSingleERC1155Token(_token common.Address, _dapp common.Address, _tokenId *big.Int, _value *big.Int, _baseLayerData []byte, _execLayerData []byte) (*types.Transaction, error) {
	return _ERC1155SinglePortal.Contract.DepositSingleERC1155Token(&_ERC1155SinglePortal.TransactOpts, _token, _dapp, _tokenId, _value, _baseLayerData, _execLayerData)
}
