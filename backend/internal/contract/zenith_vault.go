// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contract

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

// ZenithVaultMetaData contains all meta data concerning the ZenithVault contract.
var ZenithVaultMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_initialSigner\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"backendSigner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"balances\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"deposit\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"nonces\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setSigner\",\"inputs\":[{\"name\":\"_newSigner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"withdraw\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"nonce\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"signature\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"Deposit\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"token\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Withdraw\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"token\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"nonce\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"ECDSAInvalidSignature\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ECDSAInvalidSignatureLength\",\"inputs\":[{\"name\":\"length\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ECDSAInvalidSignatureS\",\"inputs\":[{\"name\":\"s\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"OwnableInvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnableUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]}]",
}

// ZenithVaultABI is the input ABI used to generate the binding from.
// Deprecated: Use ZenithVaultMetaData.ABI instead.
var ZenithVaultABI = ZenithVaultMetaData.ABI

// ZenithVault is an auto generated Go binding around an Ethereum contract.
type ZenithVault struct {
	ZenithVaultCaller     // Read-only binding to the contract
	ZenithVaultTransactor // Write-only binding to the contract
	ZenithVaultFilterer   // Log filterer for contract events
}

// ZenithVaultCaller is an auto generated read-only Go binding around an Ethereum contract.
type ZenithVaultCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ZenithVaultTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ZenithVaultTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ZenithVaultFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ZenithVaultFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ZenithVaultSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ZenithVaultSession struct {
	Contract     *ZenithVault      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ZenithVaultCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ZenithVaultCallerSession struct {
	Contract *ZenithVaultCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// ZenithVaultTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ZenithVaultTransactorSession struct {
	Contract     *ZenithVaultTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// ZenithVaultRaw is an auto generated low-level Go binding around an Ethereum contract.
type ZenithVaultRaw struct {
	Contract *ZenithVault // Generic contract binding to access the raw methods on
}

// ZenithVaultCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ZenithVaultCallerRaw struct {
	Contract *ZenithVaultCaller // Generic read-only contract binding to access the raw methods on
}

// ZenithVaultTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ZenithVaultTransactorRaw struct {
	Contract *ZenithVaultTransactor // Generic write-only contract binding to access the raw methods on
}

// NewZenithVault creates a new instance of ZenithVault, bound to a specific deployed contract.
func NewZenithVault(address common.Address, backend bind.ContractBackend) (*ZenithVault, error) {
	contract, err := bindZenithVault(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ZenithVault{ZenithVaultCaller: ZenithVaultCaller{contract: contract}, ZenithVaultTransactor: ZenithVaultTransactor{contract: contract}, ZenithVaultFilterer: ZenithVaultFilterer{contract: contract}}, nil
}

// NewZenithVaultCaller creates a new read-only instance of ZenithVault, bound to a specific deployed contract.
func NewZenithVaultCaller(address common.Address, caller bind.ContractCaller) (*ZenithVaultCaller, error) {
	contract, err := bindZenithVault(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ZenithVaultCaller{contract: contract}, nil
}

// NewZenithVaultTransactor creates a new write-only instance of ZenithVault, bound to a specific deployed contract.
func NewZenithVaultTransactor(address common.Address, transactor bind.ContractTransactor) (*ZenithVaultTransactor, error) {
	contract, err := bindZenithVault(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ZenithVaultTransactor{contract: contract}, nil
}

// NewZenithVaultFilterer creates a new log filterer instance of ZenithVault, bound to a specific deployed contract.
func NewZenithVaultFilterer(address common.Address, filterer bind.ContractFilterer) (*ZenithVaultFilterer, error) {
	contract, err := bindZenithVault(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ZenithVaultFilterer{contract: contract}, nil
}

// bindZenithVault binds a generic wrapper to an already deployed contract.
func bindZenithVault(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ZenithVaultMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ZenithVault *ZenithVaultRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ZenithVault.Contract.ZenithVaultCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ZenithVault *ZenithVaultRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ZenithVault.Contract.ZenithVaultTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ZenithVault *ZenithVaultRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ZenithVault.Contract.ZenithVaultTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ZenithVault *ZenithVaultCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ZenithVault.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ZenithVault *ZenithVaultTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ZenithVault.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ZenithVault *ZenithVaultTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ZenithVault.Contract.contract.Transact(opts, method, params...)
}

// BackendSigner is a free data retrieval call binding the contract method 0x65d65e86.
//
// Solidity: function backendSigner() view returns(address)
func (_ZenithVault *ZenithVaultCaller) BackendSigner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ZenithVault.contract.Call(opts, &out, "backendSigner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// BackendSigner is a free data retrieval call binding the contract method 0x65d65e86.
//
// Solidity: function backendSigner() view returns(address)
func (_ZenithVault *ZenithVaultSession) BackendSigner() (common.Address, error) {
	return _ZenithVault.Contract.BackendSigner(&_ZenithVault.CallOpts)
}

// BackendSigner is a free data retrieval call binding the contract method 0x65d65e86.
//
// Solidity: function backendSigner() view returns(address)
func (_ZenithVault *ZenithVaultCallerSession) BackendSigner() (common.Address, error) {
	return _ZenithVault.Contract.BackendSigner(&_ZenithVault.CallOpts)
}

// Balances is a free data retrieval call binding the contract method 0xc23f001f.
//
// Solidity: function balances(address , address ) view returns(uint256)
func (_ZenithVault *ZenithVaultCaller) Balances(opts *bind.CallOpts, arg0 common.Address, arg1 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _ZenithVault.contract.Call(opts, &out, "balances", arg0, arg1)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Balances is a free data retrieval call binding the contract method 0xc23f001f.
//
// Solidity: function balances(address , address ) view returns(uint256)
func (_ZenithVault *ZenithVaultSession) Balances(arg0 common.Address, arg1 common.Address) (*big.Int, error) {
	return _ZenithVault.Contract.Balances(&_ZenithVault.CallOpts, arg0, arg1)
}

// Balances is a free data retrieval call binding the contract method 0xc23f001f.
//
// Solidity: function balances(address , address ) view returns(uint256)
func (_ZenithVault *ZenithVaultCallerSession) Balances(arg0 common.Address, arg1 common.Address) (*big.Int, error) {
	return _ZenithVault.Contract.Balances(&_ZenithVault.CallOpts, arg0, arg1)
}

// Nonces is a free data retrieval call binding the contract method 0x7ecebe00.
//
// Solidity: function nonces(address ) view returns(uint256)
func (_ZenithVault *ZenithVaultCaller) Nonces(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _ZenithVault.contract.Call(opts, &out, "nonces", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Nonces is a free data retrieval call binding the contract method 0x7ecebe00.
//
// Solidity: function nonces(address ) view returns(uint256)
func (_ZenithVault *ZenithVaultSession) Nonces(arg0 common.Address) (*big.Int, error) {
	return _ZenithVault.Contract.Nonces(&_ZenithVault.CallOpts, arg0)
}

// Nonces is a free data retrieval call binding the contract method 0x7ecebe00.
//
// Solidity: function nonces(address ) view returns(uint256)
func (_ZenithVault *ZenithVaultCallerSession) Nonces(arg0 common.Address) (*big.Int, error) {
	return _ZenithVault.Contract.Nonces(&_ZenithVault.CallOpts, arg0)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ZenithVault *ZenithVaultCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ZenithVault.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ZenithVault *ZenithVaultSession) Owner() (common.Address, error) {
	return _ZenithVault.Contract.Owner(&_ZenithVault.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ZenithVault *ZenithVaultCallerSession) Owner() (common.Address, error) {
	return _ZenithVault.Contract.Owner(&_ZenithVault.CallOpts)
}

// Deposit is a paid mutator transaction binding the contract method 0x47e7ef24.
//
// Solidity: function deposit(address token, uint256 amount) returns()
func (_ZenithVault *ZenithVaultTransactor) Deposit(opts *bind.TransactOpts, token common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ZenithVault.contract.Transact(opts, "deposit", token, amount)
}

// Deposit is a paid mutator transaction binding the contract method 0x47e7ef24.
//
// Solidity: function deposit(address token, uint256 amount) returns()
func (_ZenithVault *ZenithVaultSession) Deposit(token common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ZenithVault.Contract.Deposit(&_ZenithVault.TransactOpts, token, amount)
}

// Deposit is a paid mutator transaction binding the contract method 0x47e7ef24.
//
// Solidity: function deposit(address token, uint256 amount) returns()
func (_ZenithVault *ZenithVaultTransactorSession) Deposit(token common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ZenithVault.Contract.Deposit(&_ZenithVault.TransactOpts, token, amount)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ZenithVault *ZenithVaultTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ZenithVault.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ZenithVault *ZenithVaultSession) RenounceOwnership() (*types.Transaction, error) {
	return _ZenithVault.Contract.RenounceOwnership(&_ZenithVault.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ZenithVault *ZenithVaultTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _ZenithVault.Contract.RenounceOwnership(&_ZenithVault.TransactOpts)
}

// SetSigner is a paid mutator transaction binding the contract method 0x6c19e783.
//
// Solidity: function setSigner(address _newSigner) returns()
func (_ZenithVault *ZenithVaultTransactor) SetSigner(opts *bind.TransactOpts, _newSigner common.Address) (*types.Transaction, error) {
	return _ZenithVault.contract.Transact(opts, "setSigner", _newSigner)
}

// SetSigner is a paid mutator transaction binding the contract method 0x6c19e783.
//
// Solidity: function setSigner(address _newSigner) returns()
func (_ZenithVault *ZenithVaultSession) SetSigner(_newSigner common.Address) (*types.Transaction, error) {
	return _ZenithVault.Contract.SetSigner(&_ZenithVault.TransactOpts, _newSigner)
}

// SetSigner is a paid mutator transaction binding the contract method 0x6c19e783.
//
// Solidity: function setSigner(address _newSigner) returns()
func (_ZenithVault *ZenithVaultTransactorSession) SetSigner(_newSigner common.Address) (*types.Transaction, error) {
	return _ZenithVault.Contract.SetSigner(&_ZenithVault.TransactOpts, _newSigner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ZenithVault *ZenithVaultTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _ZenithVault.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ZenithVault *ZenithVaultSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _ZenithVault.Contract.TransferOwnership(&_ZenithVault.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ZenithVault *ZenithVaultTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _ZenithVault.Contract.TransferOwnership(&_ZenithVault.TransactOpts, newOwner)
}

// Withdraw is a paid mutator transaction binding the contract method 0x0bd0ca9a.
//
// Solidity: function withdraw(address token, uint256 amount, uint256 nonce, bytes signature) returns()
func (_ZenithVault *ZenithVaultTransactor) Withdraw(opts *bind.TransactOpts, token common.Address, amount *big.Int, nonce *big.Int, signature []byte) (*types.Transaction, error) {
	return _ZenithVault.contract.Transact(opts, "withdraw", token, amount, nonce, signature)
}

// Withdraw is a paid mutator transaction binding the contract method 0x0bd0ca9a.
//
// Solidity: function withdraw(address token, uint256 amount, uint256 nonce, bytes signature) returns()
func (_ZenithVault *ZenithVaultSession) Withdraw(token common.Address, amount *big.Int, nonce *big.Int, signature []byte) (*types.Transaction, error) {
	return _ZenithVault.Contract.Withdraw(&_ZenithVault.TransactOpts, token, amount, nonce, signature)
}

// Withdraw is a paid mutator transaction binding the contract method 0x0bd0ca9a.
//
// Solidity: function withdraw(address token, uint256 amount, uint256 nonce, bytes signature) returns()
func (_ZenithVault *ZenithVaultTransactorSession) Withdraw(token common.Address, amount *big.Int, nonce *big.Int, signature []byte) (*types.Transaction, error) {
	return _ZenithVault.Contract.Withdraw(&_ZenithVault.TransactOpts, token, amount, nonce, signature)
}

// ZenithVaultDepositIterator is returned from FilterDeposit and is used to iterate over the raw logs and unpacked data for Deposit events raised by the ZenithVault contract.
type ZenithVaultDepositIterator struct {
	Event *ZenithVaultDeposit // Event containing the contract specifics and raw log

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
func (it *ZenithVaultDepositIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ZenithVaultDeposit)
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
		it.Event = new(ZenithVaultDeposit)
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
func (it *ZenithVaultDepositIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ZenithVaultDepositIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ZenithVaultDeposit represents a Deposit event raised by the ZenithVault contract.
type ZenithVaultDeposit struct {
	User   common.Address
	Token  common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterDeposit is a free log retrieval operation binding the contract event 0x5548c837ab068cf56a2c2479df0882a4922fd203edb7517321831d95078c5f62.
//
// Solidity: event Deposit(address indexed user, address indexed token, uint256 amount)
func (_ZenithVault *ZenithVaultFilterer) FilterDeposit(opts *bind.FilterOpts, user []common.Address, token []common.Address) (*ZenithVaultDepositIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _ZenithVault.contract.FilterLogs(opts, "Deposit", userRule, tokenRule)
	if err != nil {
		return nil, err
	}
	return &ZenithVaultDepositIterator{contract: _ZenithVault.contract, event: "Deposit", logs: logs, sub: sub}, nil
}

// WatchDeposit is a free log subscription operation binding the contract event 0x5548c837ab068cf56a2c2479df0882a4922fd203edb7517321831d95078c5f62.
//
// Solidity: event Deposit(address indexed user, address indexed token, uint256 amount)
func (_ZenithVault *ZenithVaultFilterer) WatchDeposit(opts *bind.WatchOpts, sink chan<- *ZenithVaultDeposit, user []common.Address, token []common.Address) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _ZenithVault.contract.WatchLogs(opts, "Deposit", userRule, tokenRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ZenithVaultDeposit)
				if err := _ZenithVault.contract.UnpackLog(event, "Deposit", log); err != nil {
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

// ParseDeposit is a log parse operation binding the contract event 0x5548c837ab068cf56a2c2479df0882a4922fd203edb7517321831d95078c5f62.
//
// Solidity: event Deposit(address indexed user, address indexed token, uint256 amount)
func (_ZenithVault *ZenithVaultFilterer) ParseDeposit(log types.Log) (*ZenithVaultDeposit, error) {
	event := new(ZenithVaultDeposit)
	if err := _ZenithVault.contract.UnpackLog(event, "Deposit", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ZenithVaultOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the ZenithVault contract.
type ZenithVaultOwnershipTransferredIterator struct {
	Event *ZenithVaultOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *ZenithVaultOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ZenithVaultOwnershipTransferred)
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
		it.Event = new(ZenithVaultOwnershipTransferred)
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
func (it *ZenithVaultOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ZenithVaultOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ZenithVaultOwnershipTransferred represents a OwnershipTransferred event raised by the ZenithVault contract.
type ZenithVaultOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_ZenithVault *ZenithVaultFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*ZenithVaultOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _ZenithVault.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &ZenithVaultOwnershipTransferredIterator{contract: _ZenithVault.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_ZenithVault *ZenithVaultFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *ZenithVaultOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _ZenithVault.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ZenithVaultOwnershipTransferred)
				if err := _ZenithVault.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_ZenithVault *ZenithVaultFilterer) ParseOwnershipTransferred(log types.Log) (*ZenithVaultOwnershipTransferred, error) {
	event := new(ZenithVaultOwnershipTransferred)
	if err := _ZenithVault.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ZenithVaultWithdrawIterator is returned from FilterWithdraw and is used to iterate over the raw logs and unpacked data for Withdraw events raised by the ZenithVault contract.
type ZenithVaultWithdrawIterator struct {
	Event *ZenithVaultWithdraw // Event containing the contract specifics and raw log

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
func (it *ZenithVaultWithdrawIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ZenithVaultWithdraw)
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
		it.Event = new(ZenithVaultWithdraw)
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
func (it *ZenithVaultWithdrawIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ZenithVaultWithdrawIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ZenithVaultWithdraw represents a Withdraw event raised by the ZenithVault contract.
type ZenithVaultWithdraw struct {
	User   common.Address
	Token  common.Address
	Amount *big.Int
	Nonce  *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterWithdraw is a free log retrieval operation binding the contract event 0xf341246adaac6f497bc2a656f546ab9e182111d630394f0c57c710a59a2cb567.
//
// Solidity: event Withdraw(address indexed user, address indexed token, uint256 amount, uint256 nonce)
func (_ZenithVault *ZenithVaultFilterer) FilterWithdraw(opts *bind.FilterOpts, user []common.Address, token []common.Address) (*ZenithVaultWithdrawIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _ZenithVault.contract.FilterLogs(opts, "Withdraw", userRule, tokenRule)
	if err != nil {
		return nil, err
	}
	return &ZenithVaultWithdrawIterator{contract: _ZenithVault.contract, event: "Withdraw", logs: logs, sub: sub}, nil
}

// WatchWithdraw is a free log subscription operation binding the contract event 0xf341246adaac6f497bc2a656f546ab9e182111d630394f0c57c710a59a2cb567.
//
// Solidity: event Withdraw(address indexed user, address indexed token, uint256 amount, uint256 nonce)
func (_ZenithVault *ZenithVaultFilterer) WatchWithdraw(opts *bind.WatchOpts, sink chan<- *ZenithVaultWithdraw, user []common.Address, token []common.Address) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _ZenithVault.contract.WatchLogs(opts, "Withdraw", userRule, tokenRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ZenithVaultWithdraw)
				if err := _ZenithVault.contract.UnpackLog(event, "Withdraw", log); err != nil {
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

// ParseWithdraw is a log parse operation binding the contract event 0xf341246adaac6f497bc2a656f546ab9e182111d630394f0c57c710a59a2cb567.
//
// Solidity: event Withdraw(address indexed user, address indexed token, uint256 amount, uint256 nonce)
func (_ZenithVault *ZenithVaultFilterer) ParseWithdraw(log types.Log) (*ZenithVaultWithdraw, error) {
	event := new(ZenithVaultWithdraw)
	if err := _ZenithVault.contract.UnpackLog(event, "Withdraw", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
