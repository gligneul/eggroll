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

// ERC721PortalMetaData contains all meta data concerning the ERC721Portal contract.
var ERC721PortalMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"contractIInputBox\",\"name\":\"_inputBox\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"contractIERC721\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_dapp\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_tokenId\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"_baseLayerData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"_execLayerData\",\"type\":\"bytes\"}],\"name\":\"depositERC721Token\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getInputBox\",\"outputs\":[{\"internalType\":\"contractIInputBox\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// ERC721PortalABI is the input ABI used to generate the binding from.
// Deprecated: Use ERC721PortalMetaData.ABI instead.
var ERC721PortalABI = ERC721PortalMetaData.ABI

// ERC721Portal is an auto generated Go binding around an Ethereum contract.
type ERC721Portal struct {
	ERC721PortalCaller     // Read-only binding to the contract
	ERC721PortalTransactor // Write-only binding to the contract
	ERC721PortalFilterer   // Log filterer for contract events
}

// ERC721PortalCaller is an auto generated read-only Go binding around an Ethereum contract.
type ERC721PortalCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC721PortalTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ERC721PortalTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC721PortalFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ERC721PortalFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC721PortalSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ERC721PortalSession struct {
	Contract     *ERC721Portal     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ERC721PortalCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ERC721PortalCallerSession struct {
	Contract *ERC721PortalCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// ERC721PortalTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ERC721PortalTransactorSession struct {
	Contract     *ERC721PortalTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// ERC721PortalRaw is an auto generated low-level Go binding around an Ethereum contract.
type ERC721PortalRaw struct {
	Contract *ERC721Portal // Generic contract binding to access the raw methods on
}

// ERC721PortalCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ERC721PortalCallerRaw struct {
	Contract *ERC721PortalCaller // Generic read-only contract binding to access the raw methods on
}

// ERC721PortalTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ERC721PortalTransactorRaw struct {
	Contract *ERC721PortalTransactor // Generic write-only contract binding to access the raw methods on
}

// NewERC721Portal creates a new instance of ERC721Portal, bound to a specific deployed contract.
func NewERC721Portal(address common.Address, backend bind.ContractBackend) (*ERC721Portal, error) {
	contract, err := bindERC721Portal(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ERC721Portal{ERC721PortalCaller: ERC721PortalCaller{contract: contract}, ERC721PortalTransactor: ERC721PortalTransactor{contract: contract}, ERC721PortalFilterer: ERC721PortalFilterer{contract: contract}}, nil
}

// NewERC721PortalCaller creates a new read-only instance of ERC721Portal, bound to a specific deployed contract.
func NewERC721PortalCaller(address common.Address, caller bind.ContractCaller) (*ERC721PortalCaller, error) {
	contract, err := bindERC721Portal(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ERC721PortalCaller{contract: contract}, nil
}

// NewERC721PortalTransactor creates a new write-only instance of ERC721Portal, bound to a specific deployed contract.
func NewERC721PortalTransactor(address common.Address, transactor bind.ContractTransactor) (*ERC721PortalTransactor, error) {
	contract, err := bindERC721Portal(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ERC721PortalTransactor{contract: contract}, nil
}

// NewERC721PortalFilterer creates a new log filterer instance of ERC721Portal, bound to a specific deployed contract.
func NewERC721PortalFilterer(address common.Address, filterer bind.ContractFilterer) (*ERC721PortalFilterer, error) {
	contract, err := bindERC721Portal(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ERC721PortalFilterer{contract: contract}, nil
}

// bindERC721Portal binds a generic wrapper to an already deployed contract.
func bindERC721Portal(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ERC721PortalMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ERC721Portal *ERC721PortalRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ERC721Portal.Contract.ERC721PortalCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ERC721Portal *ERC721PortalRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ERC721Portal.Contract.ERC721PortalTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ERC721Portal *ERC721PortalRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ERC721Portal.Contract.ERC721PortalTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ERC721Portal *ERC721PortalCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ERC721Portal.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ERC721Portal *ERC721PortalTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ERC721Portal.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ERC721Portal *ERC721PortalTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ERC721Portal.Contract.contract.Transact(opts, method, params...)
}

// GetInputBox is a free data retrieval call binding the contract method 0x00aace9a.
//
// Solidity: function getInputBox() view returns(address)
func (_ERC721Portal *ERC721PortalCaller) GetInputBox(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ERC721Portal.contract.Call(opts, &out, "getInputBox")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetInputBox is a free data retrieval call binding the contract method 0x00aace9a.
//
// Solidity: function getInputBox() view returns(address)
func (_ERC721Portal *ERC721PortalSession) GetInputBox() (common.Address, error) {
	return _ERC721Portal.Contract.GetInputBox(&_ERC721Portal.CallOpts)
}

// GetInputBox is a free data retrieval call binding the contract method 0x00aace9a.
//
// Solidity: function getInputBox() view returns(address)
func (_ERC721Portal *ERC721PortalCallerSession) GetInputBox() (common.Address, error) {
	return _ERC721Portal.Contract.GetInputBox(&_ERC721Portal.CallOpts)
}

// DepositERC721Token is a paid mutator transaction binding the contract method 0x28911e83.
//
// Solidity: function depositERC721Token(address _token, address _dapp, uint256 _tokenId, bytes _baseLayerData, bytes _execLayerData) returns()
func (_ERC721Portal *ERC721PortalTransactor) DepositERC721Token(opts *bind.TransactOpts, _token common.Address, _dapp common.Address, _tokenId *big.Int, _baseLayerData []byte, _execLayerData []byte) (*types.Transaction, error) {
	return _ERC721Portal.contract.Transact(opts, "depositERC721Token", _token, _dapp, _tokenId, _baseLayerData, _execLayerData)
}

// DepositERC721Token is a paid mutator transaction binding the contract method 0x28911e83.
//
// Solidity: function depositERC721Token(address _token, address _dapp, uint256 _tokenId, bytes _baseLayerData, bytes _execLayerData) returns()
func (_ERC721Portal *ERC721PortalSession) DepositERC721Token(_token common.Address, _dapp common.Address, _tokenId *big.Int, _baseLayerData []byte, _execLayerData []byte) (*types.Transaction, error) {
	return _ERC721Portal.Contract.DepositERC721Token(&_ERC721Portal.TransactOpts, _token, _dapp, _tokenId, _baseLayerData, _execLayerData)
}

// DepositERC721Token is a paid mutator transaction binding the contract method 0x28911e83.
//
// Solidity: function depositERC721Token(address _token, address _dapp, uint256 _tokenId, bytes _baseLayerData, bytes _execLayerData) returns()
func (_ERC721Portal *ERC721PortalTransactorSession) DepositERC721Token(_token common.Address, _dapp common.Address, _tokenId *big.Int, _baseLayerData []byte, _execLayerData []byte) (*types.Transaction, error) {
	return _ERC721Portal.Contract.DepositERC721Token(&_ERC721Portal.TransactOpts, _token, _dapp, _tokenId, _baseLayerData, _execLayerData)
}
