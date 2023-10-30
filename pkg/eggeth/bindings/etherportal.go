// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package bindings

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

// EtherPortalMetaData contains all meta data concerning the EtherPortal contract.
var EtherPortalMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"contractIInputBox\",\"name\":\"_inputBox\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"EtherTransferFailed\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_dapp\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"_execLayerData\",\"type\":\"bytes\"}],\"name\":\"depositEther\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getInputBox\",\"outputs\":[{\"internalType\":\"contractIInputBox\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// EtherPortalABI is the input ABI used to generate the binding from.
// Deprecated: Use EtherPortalMetaData.ABI instead.
var EtherPortalABI = EtherPortalMetaData.ABI

// EtherPortal is an auto generated Go binding around an Ethereum contract.
type EtherPortal struct {
	EtherPortalCaller     // Read-only binding to the contract
	EtherPortalTransactor // Write-only binding to the contract
	EtherPortalFilterer   // Log filterer for contract events
}

// EtherPortalCaller is an auto generated read-only Go binding around an Ethereum contract.
type EtherPortalCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EtherPortalTransactor is an auto generated write-only Go binding around an Ethereum contract.
type EtherPortalTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EtherPortalFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type EtherPortalFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EtherPortalSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type EtherPortalSession struct {
	Contract     *EtherPortal      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// EtherPortalCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type EtherPortalCallerSession struct {
	Contract *EtherPortalCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// EtherPortalTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type EtherPortalTransactorSession struct {
	Contract     *EtherPortalTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// EtherPortalRaw is an auto generated low-level Go binding around an Ethereum contract.
type EtherPortalRaw struct {
	Contract *EtherPortal // Generic contract binding to access the raw methods on
}

// EtherPortalCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type EtherPortalCallerRaw struct {
	Contract *EtherPortalCaller // Generic read-only contract binding to access the raw methods on
}

// EtherPortalTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type EtherPortalTransactorRaw struct {
	Contract *EtherPortalTransactor // Generic write-only contract binding to access the raw methods on
}

// NewEtherPortal creates a new instance of EtherPortal, bound to a specific deployed contract.
func NewEtherPortal(address common.Address, backend bind.ContractBackend) (*EtherPortal, error) {
	contract, err := bindEtherPortal(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &EtherPortal{EtherPortalCaller: EtherPortalCaller{contract: contract}, EtherPortalTransactor: EtherPortalTransactor{contract: contract}, EtherPortalFilterer: EtherPortalFilterer{contract: contract}}, nil
}

// NewEtherPortalCaller creates a new read-only instance of EtherPortal, bound to a specific deployed contract.
func NewEtherPortalCaller(address common.Address, caller bind.ContractCaller) (*EtherPortalCaller, error) {
	contract, err := bindEtherPortal(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &EtherPortalCaller{contract: contract}, nil
}

// NewEtherPortalTransactor creates a new write-only instance of EtherPortal, bound to a specific deployed contract.
func NewEtherPortalTransactor(address common.Address, transactor bind.ContractTransactor) (*EtherPortalTransactor, error) {
	contract, err := bindEtherPortal(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &EtherPortalTransactor{contract: contract}, nil
}

// NewEtherPortalFilterer creates a new log filterer instance of EtherPortal, bound to a specific deployed contract.
func NewEtherPortalFilterer(address common.Address, filterer bind.ContractFilterer) (*EtherPortalFilterer, error) {
	contract, err := bindEtherPortal(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &EtherPortalFilterer{contract: contract}, nil
}

// bindEtherPortal binds a generic wrapper to an already deployed contract.
func bindEtherPortal(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := EtherPortalMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_EtherPortal *EtherPortalRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _EtherPortal.Contract.EtherPortalCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_EtherPortal *EtherPortalRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _EtherPortal.Contract.EtherPortalTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_EtherPortal *EtherPortalRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _EtherPortal.Contract.EtherPortalTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_EtherPortal *EtherPortalCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _EtherPortal.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_EtherPortal *EtherPortalTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _EtherPortal.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_EtherPortal *EtherPortalTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _EtherPortal.Contract.contract.Transact(opts, method, params...)
}

// GetInputBox is a free data retrieval call binding the contract method 0x00aace9a.
//
// Solidity: function getInputBox() view returns(address)
func (_EtherPortal *EtherPortalCaller) GetInputBox(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _EtherPortal.contract.Call(opts, &out, "getInputBox")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetInputBox is a free data retrieval call binding the contract method 0x00aace9a.
//
// Solidity: function getInputBox() view returns(address)
func (_EtherPortal *EtherPortalSession) GetInputBox() (common.Address, error) {
	return _EtherPortal.Contract.GetInputBox(&_EtherPortal.CallOpts)
}

// GetInputBox is a free data retrieval call binding the contract method 0x00aace9a.
//
// Solidity: function getInputBox() view returns(address)
func (_EtherPortal *EtherPortalCallerSession) GetInputBox() (common.Address, error) {
	return _EtherPortal.Contract.GetInputBox(&_EtherPortal.CallOpts)
}

// DepositEther is a paid mutator transaction binding the contract method 0x938c054f.
//
// Solidity: function depositEther(address _dapp, bytes _execLayerData) payable returns()
func (_EtherPortal *EtherPortalTransactor) DepositEther(opts *bind.TransactOpts, _dapp common.Address, _execLayerData []byte) (*types.Transaction, error) {
	return _EtherPortal.contract.Transact(opts, "depositEther", _dapp, _execLayerData)
}

// DepositEther is a paid mutator transaction binding the contract method 0x938c054f.
//
// Solidity: function depositEther(address _dapp, bytes _execLayerData) payable returns()
func (_EtherPortal *EtherPortalSession) DepositEther(_dapp common.Address, _execLayerData []byte) (*types.Transaction, error) {
	return _EtherPortal.Contract.DepositEther(&_EtherPortal.TransactOpts, _dapp, _execLayerData)
}

// DepositEther is a paid mutator transaction binding the contract method 0x938c054f.
//
// Solidity: function depositEther(address _dapp, bytes _execLayerData) payable returns()
func (_EtherPortal *EtherPortalTransactorSession) DepositEther(_dapp common.Address, _execLayerData []byte) (*types.Transaction, error) {
	return _EtherPortal.Contract.DepositEther(&_EtherPortal.TransactOpts, _dapp, _execLayerData)
}
