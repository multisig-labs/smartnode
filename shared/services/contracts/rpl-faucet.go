// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contracts

import (
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
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// GGPFaucetABI is the input ABI used to generate the binding from.
const GGPFaucetABI = "[{\"constant\":true,\"inputs\":[],\"name\":\"maxWithdrawalPerPeriod\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"withdrawalFee\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"withdrawalPeriod\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"_ggpTokenAddress\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"created\",\"type\":\"uint256\"}],\"name\":\"Withdrawal\",\"type\":\"event\"},{\"constant\":false,\"inputs\":[{\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"withdraw\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"withdrawTo\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getBalance\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getAllowance\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_address\",\"type\":\"address\"}],\"name\":\"getAllowanceFor\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getWithdrawalPeriodStart\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_withdrawalPeriod\",\"type\":\"uint256\"}],\"name\":\"setWithdrawalPeriod\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_maxWithdrawalPerPeriod\",\"type\":\"uint256\"}],\"name\":\"setMaxWithdrawalPerPeriod\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_withdrawalFee\",\"type\":\"uint256\"}],\"name\":\"setWithdrawalFee\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// GGPFaucet is an auto generated Go binding around an Ethereum contract.
type GGPFaucet struct {
	GGPFaucetCaller     // Read-only binding to the contract
	GGPFaucetTransactor // Write-only binding to the contract
	GGPFaucetFilterer   // Log filterer for contract events
}

// GGPFaucetCaller is an auto generated read-only Go binding around an Ethereum contract.
type GGPFaucetCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GGPFaucetTransactor is an auto generated write-only Go binding around an Ethereum contract.
type GGPFaucetTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GGPFaucetFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type GGPFaucetFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GGPFaucetSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type GGPFaucetSession struct {
	Contract     *GGPFaucet        // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// GGPFaucetCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type GGPFaucetCallerSession struct {
	Contract *GGPFaucetCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts    // Call options to use throughout this session
}

// GGPFaucetTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type GGPFaucetTransactorSession struct {
	Contract     *GGPFaucetTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// GGPFaucetRaw is an auto generated low-level Go binding around an Ethereum contract.
type GGPFaucetRaw struct {
	Contract *GGPFaucet // Generic contract binding to access the raw methods on
}

// GGPFaucetCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type GGPFaucetCallerRaw struct {
	Contract *GGPFaucetCaller // Generic read-only contract binding to access the raw methods on
}

// GGPFaucetTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type GGPFaucetTransactorRaw struct {
	Contract *GGPFaucetTransactor // Generic write-only contract binding to access the raw methods on
}

// NewGGPFaucet creates a new instance of GGPFaucet, bound to a specific deployed contract.
func NewGGPFaucet(address common.Address, backend bind.ContractBackend) (*GGPFaucet, error) {
	contract, err := bindGGPFaucet(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &GGPFaucet{GGPFaucetCaller: GGPFaucetCaller{contract: contract}, GGPFaucetTransactor: GGPFaucetTransactor{contract: contract}, GGPFaucetFilterer: GGPFaucetFilterer{contract: contract}}, nil
}

// NewGGPFaucetCaller creates a new read-only instance of GGPFaucet, bound to a specific deployed contract.
func NewGGPFaucetCaller(address common.Address, caller bind.ContractCaller) (*GGPFaucetCaller, error) {
	contract, err := bindGGPFaucet(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &GGPFaucetCaller{contract: contract}, nil
}

// NewGGPFaucetTransactor creates a new write-only instance of GGPFaucet, bound to a specific deployed contract.
func NewGGPFaucetTransactor(address common.Address, transactor bind.ContractTransactor) (*GGPFaucetTransactor, error) {
	contract, err := bindGGPFaucet(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &GGPFaucetTransactor{contract: contract}, nil
}

// NewGGPFaucetFilterer creates a new log filterer instance of GGPFaucet, bound to a specific deployed contract.
func NewGGPFaucetFilterer(address common.Address, filterer bind.ContractFilterer) (*GGPFaucetFilterer, error) {
	contract, err := bindGGPFaucet(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &GGPFaucetFilterer{contract: contract}, nil
}

// bindGGPFaucet binds a generic wrapper to an already deployed contract.
func bindGGPFaucet(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(GGPFaucetABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_GGPFaucet *GGPFaucetRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _GGPFaucet.Contract.GGPFaucetCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_GGPFaucet *GGPFaucetRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _GGPFaucet.Contract.GGPFaucetTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_GGPFaucet *GGPFaucetRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _GGPFaucet.Contract.GGPFaucetTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_GGPFaucet *GGPFaucetCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _GGPFaucet.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_GGPFaucet *GGPFaucetTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _GGPFaucet.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_GGPFaucet *GGPFaucetTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _GGPFaucet.Contract.contract.Transact(opts, method, params...)
}

// GetAllowance is a free data retrieval call binding the contract method 0x973e9b8b.
//
// Solidity: function getAllowance() view returns(uint256)
func (_GGPFaucet *GGPFaucetCaller) GetAllowance(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _GGPFaucet.contract.Call(opts, &out, "getAllowance")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetAllowance is a free data retrieval call binding the contract method 0x973e9b8b.
//
// Solidity: function getAllowance() view returns(uint256)
func (_GGPFaucet *GGPFaucetSession) GetAllowance() (*big.Int, error) {
	return _GGPFaucet.Contract.GetAllowance(&_GGPFaucet.CallOpts)
}

// GetAllowance is a free data retrieval call binding the contract method 0x973e9b8b.
//
// Solidity: function getAllowance() view returns(uint256)
func (_GGPFaucet *GGPFaucetCallerSession) GetAllowance() (*big.Int, error) {
	return _GGPFaucet.Contract.GetAllowance(&_GGPFaucet.CallOpts)
}

// GetAllowanceFor is a free data retrieval call binding the contract method 0x7639a24b.
//
// Solidity: function getAllowanceFor(address _address) view returns(uint256)
func (_GGPFaucet *GGPFaucetCaller) GetAllowanceFor(opts *bind.CallOpts, _address common.Address) (*big.Int, error) {
	var out []interface{}
	err := _GGPFaucet.contract.Call(opts, &out, "getAllowanceFor", _address)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetAllowanceFor is a free data retrieval call binding the contract method 0x7639a24b.
//
// Solidity: function getAllowanceFor(address _address) view returns(uint256)
func (_GGPFaucet *GGPFaucetSession) GetAllowanceFor(_address common.Address) (*big.Int, error) {
	return _GGPFaucet.Contract.GetAllowanceFor(&_GGPFaucet.CallOpts, _address)
}

// GetAllowanceFor is a free data retrieval call binding the contract method 0x7639a24b.
//
// Solidity: function getAllowanceFor(address _address) view returns(uint256)
func (_GGPFaucet *GGPFaucetCallerSession) GetAllowanceFor(_address common.Address) (*big.Int, error) {
	return _GGPFaucet.Contract.GetAllowanceFor(&_GGPFaucet.CallOpts, _address)
}

// GetBalance is a free data retrieval call binding the contract method 0x12065fe0.
//
// Solidity: function getBalance() view returns(uint256)
func (_GGPFaucet *GGPFaucetCaller) GetBalance(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _GGPFaucet.contract.Call(opts, &out, "getBalance")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetBalance is a free data retrieval call binding the contract method 0x12065fe0.
//
// Solidity: function getBalance() view returns(uint256)
func (_GGPFaucet *GGPFaucetSession) GetBalance() (*big.Int, error) {
	return _GGPFaucet.Contract.GetBalance(&_GGPFaucet.CallOpts)
}

// GetBalance is a free data retrieval call binding the contract method 0x12065fe0.
//
// Solidity: function getBalance() view returns(uint256)
func (_GGPFaucet *GGPFaucetCallerSession) GetBalance() (*big.Int, error) {
	return _GGPFaucet.Contract.GetBalance(&_GGPFaucet.CallOpts)
}

// GetWithdrawalPeriodStart is a free data retrieval call binding the contract method 0xfc65bc4f.
//
// Solidity: function getWithdrawalPeriodStart() view returns(uint256)
func (_GGPFaucet *GGPFaucetCaller) GetWithdrawalPeriodStart(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _GGPFaucet.contract.Call(opts, &out, "getWithdrawalPeriodStart")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetWithdrawalPeriodStart is a free data retrieval call binding the contract method 0xfc65bc4f.
//
// Solidity: function getWithdrawalPeriodStart() view returns(uint256)
func (_GGPFaucet *GGPFaucetSession) GetWithdrawalPeriodStart() (*big.Int, error) {
	return _GGPFaucet.Contract.GetWithdrawalPeriodStart(&_GGPFaucet.CallOpts)
}

// GetWithdrawalPeriodStart is a free data retrieval call binding the contract method 0xfc65bc4f.
//
// Solidity: function getWithdrawalPeriodStart() view returns(uint256)
func (_GGPFaucet *GGPFaucetCallerSession) GetWithdrawalPeriodStart() (*big.Int, error) {
	return _GGPFaucet.Contract.GetWithdrawalPeriodStart(&_GGPFaucet.CallOpts)
}

// MaxWithdrawalPerPeriod is a free data retrieval call binding the contract method 0x203bf056.
//
// Solidity: function maxWithdrawalPerPeriod() view returns(uint256)
func (_GGPFaucet *GGPFaucetCaller) MaxWithdrawalPerPeriod(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _GGPFaucet.contract.Call(opts, &out, "maxWithdrawalPerPeriod")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MaxWithdrawalPerPeriod is a free data retrieval call binding the contract method 0x203bf056.
//
// Solidity: function maxWithdrawalPerPeriod() view returns(uint256)
func (_GGPFaucet *GGPFaucetSession) MaxWithdrawalPerPeriod() (*big.Int, error) {
	return _GGPFaucet.Contract.MaxWithdrawalPerPeriod(&_GGPFaucet.CallOpts)
}

// MaxWithdrawalPerPeriod is a free data retrieval call binding the contract method 0x203bf056.
//
// Solidity: function maxWithdrawalPerPeriod() view returns(uint256)
func (_GGPFaucet *GGPFaucetCallerSession) MaxWithdrawalPerPeriod() (*big.Int, error) {
	return _GGPFaucet.Contract.MaxWithdrawalPerPeriod(&_GGPFaucet.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_GGPFaucet *GGPFaucetCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _GGPFaucet.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_GGPFaucet *GGPFaucetSession) Owner() (common.Address, error) {
	return _GGPFaucet.Contract.Owner(&_GGPFaucet.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_GGPFaucet *GGPFaucetCallerSession) Owner() (common.Address, error) {
	return _GGPFaucet.Contract.Owner(&_GGPFaucet.CallOpts)
}

// WithdrawalFee is a free data retrieval call binding the contract method 0x8bc7e8c4.
//
// Solidity: function withdrawalFee() view returns(uint256)
func (_GGPFaucet *GGPFaucetCaller) WithdrawalFee(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _GGPFaucet.contract.Call(opts, &out, "withdrawalFee")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// WithdrawalFee is a free data retrieval call binding the contract method 0x8bc7e8c4.
//
// Solidity: function withdrawalFee() view returns(uint256)
func (_GGPFaucet *GGPFaucetSession) WithdrawalFee() (*big.Int, error) {
	return _GGPFaucet.Contract.WithdrawalFee(&_GGPFaucet.CallOpts)
}

// WithdrawalFee is a free data retrieval call binding the contract method 0x8bc7e8c4.
//
// Solidity: function withdrawalFee() view returns(uint256)
func (_GGPFaucet *GGPFaucetCallerSession) WithdrawalFee() (*big.Int, error) {
	return _GGPFaucet.Contract.WithdrawalFee(&_GGPFaucet.CallOpts)
}

// WithdrawalPeriod is a free data retrieval call binding the contract method 0xbca7093d.
//
// Solidity: function withdrawalPeriod() view returns(uint256)
func (_GGPFaucet *GGPFaucetCaller) WithdrawalPeriod(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _GGPFaucet.contract.Call(opts, &out, "withdrawalPeriod")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// WithdrawalPeriod is a free data retrieval call binding the contract method 0xbca7093d.
//
// Solidity: function withdrawalPeriod() view returns(uint256)
func (_GGPFaucet *GGPFaucetSession) WithdrawalPeriod() (*big.Int, error) {
	return _GGPFaucet.Contract.WithdrawalPeriod(&_GGPFaucet.CallOpts)
}

// WithdrawalPeriod is a free data retrieval call binding the contract method 0xbca7093d.
//
// Solidity: function withdrawalPeriod() view returns(uint256)
func (_GGPFaucet *GGPFaucetCallerSession) WithdrawalPeriod() (*big.Int, error) {
	return _GGPFaucet.Contract.WithdrawalPeriod(&_GGPFaucet.CallOpts)
}

// SetMaxWithdrawalPerPeriod is a paid mutator transaction binding the contract method 0xc0ac9128.
//
// Solidity: function setMaxWithdrawalPerPeriod(uint256 _maxWithdrawalPerPeriod) returns()
func (_GGPFaucet *GGPFaucetTransactor) SetMaxWithdrawalPerPeriod(opts *bind.TransactOpts, _maxWithdrawalPerPeriod *big.Int) (*types.Transaction, error) {
	return _GGPFaucet.contract.Transact(opts, "setMaxWithdrawalPerPeriod", _maxWithdrawalPerPeriod)
}

// SetMaxWithdrawalPerPeriod is a paid mutator transaction binding the contract method 0xc0ac9128.
//
// Solidity: function setMaxWithdrawalPerPeriod(uint256 _maxWithdrawalPerPeriod) returns()
func (_GGPFaucet *GGPFaucetSession) SetMaxWithdrawalPerPeriod(_maxWithdrawalPerPeriod *big.Int) (*types.Transaction, error) {
	return _GGPFaucet.Contract.SetMaxWithdrawalPerPeriod(&_GGPFaucet.TransactOpts, _maxWithdrawalPerPeriod)
}

// SetMaxWithdrawalPerPeriod is a paid mutator transaction binding the contract method 0xc0ac9128.
//
// Solidity: function setMaxWithdrawalPerPeriod(uint256 _maxWithdrawalPerPeriod) returns()
func (_GGPFaucet *GGPFaucetTransactorSession) SetMaxWithdrawalPerPeriod(_maxWithdrawalPerPeriod *big.Int) (*types.Transaction, error) {
	return _GGPFaucet.Contract.SetMaxWithdrawalPerPeriod(&_GGPFaucet.TransactOpts, _maxWithdrawalPerPeriod)
}

// SetWithdrawalFee is a paid mutator transaction binding the contract method 0xac1e5025.
//
// Solidity: function setWithdrawalFee(uint256 _withdrawalFee) returns()
func (_GGPFaucet *GGPFaucetTransactor) SetWithdrawalFee(opts *bind.TransactOpts, _withdrawalFee *big.Int) (*types.Transaction, error) {
	return _GGPFaucet.contract.Transact(opts, "setWithdrawalFee", _withdrawalFee)
}

// SetWithdrawalFee is a paid mutator transaction binding the contract method 0xac1e5025.
//
// Solidity: function setWithdrawalFee(uint256 _withdrawalFee) returns()
func (_GGPFaucet *GGPFaucetSession) SetWithdrawalFee(_withdrawalFee *big.Int) (*types.Transaction, error) {
	return _GGPFaucet.Contract.SetWithdrawalFee(&_GGPFaucet.TransactOpts, _withdrawalFee)
}

// SetWithdrawalFee is a paid mutator transaction binding the contract method 0xac1e5025.
//
// Solidity: function setWithdrawalFee(uint256 _withdrawalFee) returns()
func (_GGPFaucet *GGPFaucetTransactorSession) SetWithdrawalFee(_withdrawalFee *big.Int) (*types.Transaction, error) {
	return _GGPFaucet.Contract.SetWithdrawalFee(&_GGPFaucet.TransactOpts, _withdrawalFee)
}

// SetWithdrawalPeriod is a paid mutator transaction binding the contract method 0x973b294f.
//
// Solidity: function setWithdrawalPeriod(uint256 _withdrawalPeriod) returns()
func (_GGPFaucet *GGPFaucetTransactor) SetWithdrawalPeriod(opts *bind.TransactOpts, _withdrawalPeriod *big.Int) (*types.Transaction, error) {
	return _GGPFaucet.contract.Transact(opts, "setWithdrawalPeriod", _withdrawalPeriod)
}

// SetWithdrawalPeriod is a paid mutator transaction binding the contract method 0x973b294f.
//
// Solidity: function setWithdrawalPeriod(uint256 _withdrawalPeriod) returns()
func (_GGPFaucet *GGPFaucetSession) SetWithdrawalPeriod(_withdrawalPeriod *big.Int) (*types.Transaction, error) {
	return _GGPFaucet.Contract.SetWithdrawalPeriod(&_GGPFaucet.TransactOpts, _withdrawalPeriod)
}

// SetWithdrawalPeriod is a paid mutator transaction binding the contract method 0x973b294f.
//
// Solidity: function setWithdrawalPeriod(uint256 _withdrawalPeriod) returns()
func (_GGPFaucet *GGPFaucetTransactorSession) SetWithdrawalPeriod(_withdrawalPeriod *big.Int) (*types.Transaction, error) {
	return _GGPFaucet.Contract.SetWithdrawalPeriod(&_GGPFaucet.TransactOpts, _withdrawalPeriod)
}

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
//
// Solidity: function withdraw(uint256 _amount) payable returns(bool)
func (_GGPFaucet *GGPFaucetTransactor) Withdraw(opts *bind.TransactOpts, _amount *big.Int) (*types.Transaction, error) {
	return _GGPFaucet.contract.Transact(opts, "withdraw", _amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
//
// Solidity: function withdraw(uint256 _amount) payable returns(bool)
func (_GGPFaucet *GGPFaucetSession) Withdraw(_amount *big.Int) (*types.Transaction, error) {
	return _GGPFaucet.Contract.Withdraw(&_GGPFaucet.TransactOpts, _amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
//
// Solidity: function withdraw(uint256 _amount) payable returns(bool)
func (_GGPFaucet *GGPFaucetTransactorSession) Withdraw(_amount *big.Int) (*types.Transaction, error) {
	return _GGPFaucet.Contract.Withdraw(&_GGPFaucet.TransactOpts, _amount)
}

// WithdrawTo is a paid mutator transaction binding the contract method 0x205c2878.
//
// Solidity: function withdrawTo(address _to, uint256 _amount) payable returns(bool)
func (_GGPFaucet *GGPFaucetTransactor) WithdrawTo(opts *bind.TransactOpts, _to common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _GGPFaucet.contract.Transact(opts, "withdrawTo", _to, _amount)
}

// WithdrawTo is a paid mutator transaction binding the contract method 0x205c2878.
//
// Solidity: function withdrawTo(address _to, uint256 _amount) payable returns(bool)
func (_GGPFaucet *GGPFaucetSession) WithdrawTo(_to common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _GGPFaucet.Contract.WithdrawTo(&_GGPFaucet.TransactOpts, _to, _amount)
}

// WithdrawTo is a paid mutator transaction binding the contract method 0x205c2878.
//
// Solidity: function withdrawTo(address _to, uint256 _amount) payable returns(bool)
func (_GGPFaucet *GGPFaucetTransactorSession) WithdrawTo(_to common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _GGPFaucet.Contract.WithdrawTo(&_GGPFaucet.TransactOpts, _to, _amount)
}

// GGPFaucetWithdrawalIterator is returned from FilterWithdrawal and is used to iterate over the raw logs and unpacked data for Withdrawal events raised by the GGPFaucet contract.
type GGPFaucetWithdrawalIterator struct {
	Event *GGPFaucetWithdrawal // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *GGPFaucetWithdrawalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GGPFaucetWithdrawal)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(GGPFaucetWithdrawal)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *GGPFaucetWithdrawalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GGPFaucetWithdrawalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GGPFaucetWithdrawal represents a Withdrawal event raised by the GGPFaucet contract.
type GGPFaucetWithdrawal struct {
	To      common.Address
	Value   *big.Int
	Created *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterWithdrawal is a free log retrieval operation binding the contract event 0xdf273cb619d95419a9cd0ec88123a0538c85064229baa6363788f743fff90deb.
//
// Solidity: event Withdrawal(address indexed to, uint256 value, uint256 created)
func (_GGPFaucet *GGPFaucetFilterer) FilterWithdrawal(opts *bind.FilterOpts, to []common.Address) (*GGPFaucetWithdrawalIterator, error) {

	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _GGPFaucet.contract.FilterLogs(opts, "Withdrawal", toRule)
	if err != nil {
		return nil, err
	}
	return &GGPFaucetWithdrawalIterator{contract: _GGPFaucet.contract, event: "Withdrawal", logs: logs, sub: sub}, nil
}

// WatchWithdrawal is a free log subscription operation binding the contract event 0xdf273cb619d95419a9cd0ec88123a0538c85064229baa6363788f743fff90deb.
//
// Solidity: event Withdrawal(address indexed to, uint256 value, uint256 created)
func (_GGPFaucet *GGPFaucetFilterer) WatchWithdrawal(opts *bind.WatchOpts, sink chan<- *GGPFaucetWithdrawal, to []common.Address) (event.Subscription, error) {

	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _GGPFaucet.contract.WatchLogs(opts, "Withdrawal", toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GGPFaucetWithdrawal)
				if err := _GGPFaucet.contract.UnpackLog(event, "Withdrawal", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseWithdrawal is a log parse operation binding the contract event 0xdf273cb619d95419a9cd0ec88123a0538c85064229baa6363788f743fff90deb.
//
// Solidity: event Withdrawal(address indexed to, uint256 value, uint256 created)
func (_GGPFaucet *GGPFaucetFilterer) ParseWithdrawal(log types.Log) (*GGPFaucetWithdrawal, error) {
	event := new(GGPFaucetWithdrawal)
	if err := _GGPFaucet.contract.UnpackLog(event, "Withdrawal", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
