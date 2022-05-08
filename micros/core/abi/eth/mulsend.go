package eth

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
)

// EthMultiSenderCMetaData contains all meta data concerning the EthMultiSenderC contract.
var EthMultiSenderCMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Decline\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"receivers\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"uint256[]\",\"name\":\"amounts\",\"type\":\"uint256[]\"}],\"name\":\"Disperse\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"NativeTx\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"VERSION\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"deprecated\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"disable\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"_receivers\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"_values\",\"type\":\"uint256[]\"}],\"name\":\"disperse\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// EthMultiSenderCABI is the input ABI used to generate the binding from.
// Deprecated: Use EthMultiSenderCMetaData.ABI instead.
var EthMultiSenderCABI = EthMultiSenderCMetaData.ABI

// EthMultiSenderC is an auto generated Go binding around an Ethereum contract.
type EthMultiSenderC struct {
	EthMultiSenderCCaller     // Read-only binding to the contract
	EthMultiSenderCTransactor // Write-only binding to the contract
	EthMultiSenderCFilterer   // Log filterer for contract events
}

// EthMultiSenderCCaller is an auto generated read-only Go binding around an Ethereum contract.
type EthMultiSenderCCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EthMultiSenderCTransactor is an auto generated write-only Go binding around an Ethereum contract.
type EthMultiSenderCTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EthMultiSenderCFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type EthMultiSenderCFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EthMultiSenderCSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type EthMultiSenderCSession struct {
	Contract     *EthMultiSenderC  // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// EthMultiSenderCCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type EthMultiSenderCCallerSession struct {
	Contract *EthMultiSenderCCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts          // Call options to use throughout this session
}

// EthMultiSenderCTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type EthMultiSenderCTransactorSession struct {
	Contract     *EthMultiSenderCTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// EthMultiSenderCRaw is an auto generated low-level Go binding around an Ethereum contract.
type EthMultiSenderCRaw struct {
	Contract *EthMultiSenderC // Generic contract binding to access the raw methods on
}

// EthMultiSenderCCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type EthMultiSenderCCallerRaw struct {
	Contract *EthMultiSenderCCaller // Generic read-only contract binding to access the raw methods on
}

// EthMultiSenderCTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type EthMultiSenderCTransactorRaw struct {
	Contract *EthMultiSenderCTransactor // Generic write-only contract binding to access the raw methods on
}

// NewEthMultiSenderC creates a new instance of EthMultiSenderC, bound to a specific deployed contract.
func NewEthMultiSenderC(address common.Address, backend bind.ContractBackend) (*EthMultiSenderC, error) {
	contract, err := bindEthMultiSenderC(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &EthMultiSenderC{EthMultiSenderCCaller: EthMultiSenderCCaller{contract: contract}, EthMultiSenderCTransactor: EthMultiSenderCTransactor{contract: contract}, EthMultiSenderCFilterer: EthMultiSenderCFilterer{contract: contract}}, nil
}

// NewEthMultiSenderCCaller creates a new read-only instance of EthMultiSenderC, bound to a specific deployed contract.
func NewEthMultiSenderCCaller(address common.Address, caller bind.ContractCaller) (*EthMultiSenderCCaller, error) {
	contract, err := bindEthMultiSenderC(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &EthMultiSenderCCaller{contract: contract}, nil
}

// NewEthMultiSenderCTransactor creates a new write-only instance of EthMultiSenderC, bound to a specific deployed contract.
func NewEthMultiSenderCTransactor(address common.Address, transactor bind.ContractTransactor) (*EthMultiSenderCTransactor, error) {
	contract, err := bindEthMultiSenderC(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &EthMultiSenderCTransactor{contract: contract}, nil
}

// NewEthMultiSenderCFilterer creates a new log filterer instance of EthMultiSenderC, bound to a specific deployed contract.
func NewEthMultiSenderCFilterer(address common.Address, filterer bind.ContractFilterer) (*EthMultiSenderCFilterer, error) {
	contract, err := bindEthMultiSenderC(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &EthMultiSenderCFilterer{contract: contract}, nil
}

// bindEthMultiSenderC binds a generic wrapper to an already deployed contract.
func bindEthMultiSenderC(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(EthMultiSenderCABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_EthMultiSenderC *EthMultiSenderCRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _EthMultiSenderC.Contract.EthMultiSenderCCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_EthMultiSenderC *EthMultiSenderCRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _EthMultiSenderC.Contract.EthMultiSenderCTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_EthMultiSenderC *EthMultiSenderCRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _EthMultiSenderC.Contract.EthMultiSenderCTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_EthMultiSenderC *EthMultiSenderCCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _EthMultiSenderC.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_EthMultiSenderC *EthMultiSenderCTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _EthMultiSenderC.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_EthMultiSenderC *EthMultiSenderCTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _EthMultiSenderC.Contract.contract.Transact(opts, method, params...)
}

// VERSION is a free data retrieval call binding the contract method 0xffa1ad74.
//
// Solidity: function VERSION() view returns(string)
func (_EthMultiSenderC *EthMultiSenderCCaller) VERSION(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _EthMultiSenderC.contract.Call(opts, &out, "VERSION")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// VERSION is a free data retrieval call binding the contract method 0xffa1ad74.
//
// Solidity: function VERSION() view returns(string)
func (_EthMultiSenderC *EthMultiSenderCSession) VERSION() (string, error) {
	return _EthMultiSenderC.Contract.VERSION(&_EthMultiSenderC.CallOpts)
}

// VERSION is a free data retrieval call binding the contract method 0xffa1ad74.
//
// Solidity: function VERSION() view returns(string)
func (_EthMultiSenderC *EthMultiSenderCCallerSession) VERSION() (string, error) {
	return _EthMultiSenderC.Contract.VERSION(&_EthMultiSenderC.CallOpts)
}

// Deprecated is a free data retrieval call binding the contract method 0x0e136b19.
//
// Solidity: function deprecated() view returns(bool)
func (_EthMultiSenderC *EthMultiSenderCCaller) Deprecated(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _EthMultiSenderC.contract.Call(opts, &out, "deprecated")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Deprecated is a free data retrieval call binding the contract method 0x0e136b19.
//
// Solidity: function deprecated() view returns(bool)
func (_EthMultiSenderC *EthMultiSenderCSession) Deprecated() (bool, error) {
	return _EthMultiSenderC.Contract.Deprecated(&_EthMultiSenderC.CallOpts)
}

// Deprecated is a free data retrieval call binding the contract method 0x0e136b19.
//
// Solidity: function deprecated() view returns(bool)
func (_EthMultiSenderC *EthMultiSenderCCallerSession) Deprecated() (bool, error) {
	return _EthMultiSenderC.Contract.Deprecated(&_EthMultiSenderC.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_EthMultiSenderC *EthMultiSenderCCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _EthMultiSenderC.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_EthMultiSenderC *EthMultiSenderCSession) Owner() (common.Address, error) {
	return _EthMultiSenderC.Contract.Owner(&_EthMultiSenderC.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_EthMultiSenderC *EthMultiSenderCCallerSession) Owner() (common.Address, error) {
	return _EthMultiSenderC.Contract.Owner(&_EthMultiSenderC.CallOpts)
}

// Disable is a paid mutator transaction binding the contract method 0x2f2770db.
//
// Solidity: function disable() returns()
func (_EthMultiSenderC *EthMultiSenderCTransactor) Disable(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _EthMultiSenderC.contract.Transact(opts, "disable")
}

// Disable is a paid mutator transaction binding the contract method 0x2f2770db.
//
// Solidity: function disable() returns()
func (_EthMultiSenderC *EthMultiSenderCSession) Disable() (*types.Transaction, error) {
	return _EthMultiSenderC.Contract.Disable(&_EthMultiSenderC.TransactOpts)
}

// Disable is a paid mutator transaction binding the contract method 0x2f2770db.
//
// Solidity: function disable() returns()
func (_EthMultiSenderC *EthMultiSenderCTransactorSession) Disable() (*types.Transaction, error) {
	return _EthMultiSenderC.Contract.Disable(&_EthMultiSenderC.TransactOpts)
}

// Disperse is a paid mutator transaction binding the contract method 0xc87b1ae3.
//
// Solidity: function disperse(address _token, address[] _receivers, uint256[] _values) payable returns()
func (_EthMultiSenderC *EthMultiSenderCTransactor) Disperse(opts *bind.TransactOpts, _token common.Address, _receivers []common.Address, _values []*big.Int) (*types.Transaction, error) {
	return _EthMultiSenderC.contract.Transact(opts, "disperse", _token, _receivers, _values)
}

// Disperse is a paid mutator transaction binding the contract method 0xc87b1ae3.
//
// Solidity: function disperse(address _token, address[] _receivers, uint256[] _values) payable returns()
func (_EthMultiSenderC *EthMultiSenderCSession) Disperse(_token common.Address, _receivers []common.Address, _values []*big.Int) (*types.Transaction, error) {
	return _EthMultiSenderC.Contract.Disperse(&_EthMultiSenderC.TransactOpts, _token, _receivers, _values)
}

// Disperse is a paid mutator transaction binding the contract method 0xc87b1ae3.
//
// Solidity: function disperse(address _token, address[] _receivers, uint256[] _values) payable returns()
func (_EthMultiSenderC *EthMultiSenderCTransactorSession) Disperse(_token common.Address, _receivers []common.Address, _values []*big.Int) (*types.Transaction, error) {
	return _EthMultiSenderC.Contract.Disperse(&_EthMultiSenderC.TransactOpts, _token, _receivers, _values)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_EthMultiSenderC *EthMultiSenderCTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _EthMultiSenderC.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_EthMultiSenderC *EthMultiSenderCSession) RenounceOwnership() (*types.Transaction, error) {
	return _EthMultiSenderC.Contract.RenounceOwnership(&_EthMultiSenderC.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_EthMultiSenderC *EthMultiSenderCTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _EthMultiSenderC.Contract.RenounceOwnership(&_EthMultiSenderC.TransactOpts)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address _to, uint256 _value) returns()
func (_EthMultiSenderC *EthMultiSenderCTransactor) Transfer(opts *bind.TransactOpts, _to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _EthMultiSenderC.contract.Transact(opts, "transfer", _to, _value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address _to, uint256 _value) returns()
func (_EthMultiSenderC *EthMultiSenderCSession) Transfer(_to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _EthMultiSenderC.Contract.Transfer(&_EthMultiSenderC.TransactOpts, _to, _value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address _to, uint256 _value) returns()
func (_EthMultiSenderC *EthMultiSenderCTransactorSession) Transfer(_to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _EthMultiSenderC.Contract.Transfer(&_EthMultiSenderC.TransactOpts, _to, _value)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_EthMultiSenderC *EthMultiSenderCTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _EthMultiSenderC.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_EthMultiSenderC *EthMultiSenderCSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _EthMultiSenderC.Contract.TransferOwnership(&_EthMultiSenderC.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_EthMultiSenderC *EthMultiSenderCTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _EthMultiSenderC.Contract.TransferOwnership(&_EthMultiSenderC.TransactOpts, newOwner)
}

// EthMultiSenderCDeclineIterator is returned from FilterDecline and is used to iterate over the raw logs and unpacked data for Decline events raised by the EthMultiSenderC contract.
type EthMultiSenderCDeclineIterator struct {
	Event *EthMultiSenderCDecline // Event containing the contract specifics and raw log

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
func (it *EthMultiSenderCDeclineIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EthMultiSenderCDecline)
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
		it.Event = new(EthMultiSenderCDecline)
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
func (it *EthMultiSenderCDeclineIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EthMultiSenderCDeclineIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EthMultiSenderCDecline represents a Decline event raised by the EthMultiSenderC contract.
type EthMultiSenderCDecline struct {
	To     common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterDecline is a free log retrieval operation binding the contract event 0x678197bf918152597106edcd055b1e253a0e65adcd9c7917ab52151628ffafc3.
//
// Solidity: event Decline(address indexed to, uint256 amount)
func (_EthMultiSenderC *EthMultiSenderCFilterer) FilterDecline(opts *bind.FilterOpts, to []common.Address) (*EthMultiSenderCDeclineIterator, error) {

	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _EthMultiSenderC.contract.FilterLogs(opts, "Decline", toRule)
	if err != nil {
		return nil, err
	}
	return &EthMultiSenderCDeclineIterator{contract: _EthMultiSenderC.contract, event: "Decline", logs: logs, sub: sub}, nil
}

// WatchDecline is a free log subscription operation binding the contract event 0x678197bf918152597106edcd055b1e253a0e65adcd9c7917ab52151628ffafc3.
//
// Solidity: event Decline(address indexed to, uint256 amount)
func (_EthMultiSenderC *EthMultiSenderCFilterer) WatchDecline(opts *bind.WatchOpts, sink chan<- *EthMultiSenderCDecline, to []common.Address) (event.Subscription, error) {

	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _EthMultiSenderC.contract.WatchLogs(opts, "Decline", toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EthMultiSenderCDecline)
				if err := _EthMultiSenderC.contract.UnpackLog(event, "Decline", log); err != nil {
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

// ParseDecline is a log parse operation binding the contract event 0x678197bf918152597106edcd055b1e253a0e65adcd9c7917ab52151628ffafc3.
//
// Solidity: event Decline(address indexed to, uint256 amount)
func (_EthMultiSenderC *EthMultiSenderCFilterer) ParseDecline(log types.Log) (*EthMultiSenderCDecline, error) {
	event := new(EthMultiSenderCDecline)
	if err := _EthMultiSenderC.contract.UnpackLog(event, "Decline", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EthMultiSenderCDisperseIterator is returned from FilterDisperse and is used to iterate over the raw logs and unpacked data for Disperse events raised by the EthMultiSenderC contract.
type EthMultiSenderCDisperseIterator struct {
	Event *EthMultiSenderCDisperse // Event containing the contract specifics and raw log

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
func (it *EthMultiSenderCDisperseIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EthMultiSenderCDisperse)
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
		it.Event = new(EthMultiSenderCDisperse)
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
func (it *EthMultiSenderCDisperseIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EthMultiSenderCDisperseIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EthMultiSenderCDisperse represents a Disperse event raised by the EthMultiSenderC contract.
type EthMultiSenderCDisperse struct {
	Token     common.Address
	Receivers []common.Address
	Amounts   []*big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterDisperse is a free log retrieval operation binding the contract event 0x814fa8436d1f10ae4171a791465b8129893bcdc91cf1dbd72abca9475b3ad1be.
//
// Solidity: event Disperse(address indexed token, address[] receivers, uint256[] amounts)
func (_EthMultiSenderC *EthMultiSenderCFilterer) FilterDisperse(opts *bind.FilterOpts, token []common.Address) (*EthMultiSenderCDisperseIterator, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _EthMultiSenderC.contract.FilterLogs(opts, "Disperse", tokenRule)
	if err != nil {
		return nil, err
	}
	return &EthMultiSenderCDisperseIterator{contract: _EthMultiSenderC.contract, event: "Disperse", logs: logs, sub: sub}, nil
}

// WatchDisperse is a free log subscription operation binding the contract event 0x814fa8436d1f10ae4171a791465b8129893bcdc91cf1dbd72abca9475b3ad1be.
//
// Solidity: event Disperse(address indexed token, address[] receivers, uint256[] amounts)
func (_EthMultiSenderC *EthMultiSenderCFilterer) WatchDisperse(opts *bind.WatchOpts, sink chan<- *EthMultiSenderCDisperse, token []common.Address) (event.Subscription, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _EthMultiSenderC.contract.WatchLogs(opts, "Disperse", tokenRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EthMultiSenderCDisperse)
				if err := _EthMultiSenderC.contract.UnpackLog(event, "Disperse", log); err != nil {
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

// ParseDisperse is a log parse operation binding the contract event 0x814fa8436d1f10ae4171a791465b8129893bcdc91cf1dbd72abca9475b3ad1be.
//
// Solidity: event Disperse(address indexed token, address[] receivers, uint256[] amounts)
func (_EthMultiSenderC *EthMultiSenderCFilterer) ParseDisperse(log types.Log) (*EthMultiSenderCDisperse, error) {
	event := new(EthMultiSenderCDisperse)
	if err := _EthMultiSenderC.contract.UnpackLog(event, "Disperse", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EthMultiSenderCNativeTxIterator is returned from FilterNativeTx and is used to iterate over the raw logs and unpacked data for NativeTx events raised by the EthMultiSenderC contract.
type EthMultiSenderCNativeTxIterator struct {
	Event *EthMultiSenderCNativeTx // Event containing the contract specifics and raw log

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
func (it *EthMultiSenderCNativeTxIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EthMultiSenderCNativeTx)
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
		it.Event = new(EthMultiSenderCNativeTx)
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
func (it *EthMultiSenderCNativeTxIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EthMultiSenderCNativeTxIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EthMultiSenderCNativeTx represents a NativeTx event raised by the EthMultiSenderC contract.
type EthMultiSenderCNativeTx struct {
	To     common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterNativeTx is a free log retrieval operation binding the contract event 0x311279b6436c79139d7d8a1dc58c6cd2499083a8cc78bfdd8444c5b6fd89ba0d.
//
// Solidity: event NativeTx(address indexed to, uint256 amount)
func (_EthMultiSenderC *EthMultiSenderCFilterer) FilterNativeTx(opts *bind.FilterOpts, to []common.Address) (*EthMultiSenderCNativeTxIterator, error) {

	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _EthMultiSenderC.contract.FilterLogs(opts, "NativeTx", toRule)
	if err != nil {
		return nil, err
	}
	return &EthMultiSenderCNativeTxIterator{contract: _EthMultiSenderC.contract, event: "NativeTx", logs: logs, sub: sub}, nil
}

// WatchNativeTx is a free log subscription operation binding the contract event 0x311279b6436c79139d7d8a1dc58c6cd2499083a8cc78bfdd8444c5b6fd89ba0d.
//
// Solidity: event NativeTx(address indexed to, uint256 amount)
func (_EthMultiSenderC *EthMultiSenderCFilterer) WatchNativeTx(opts *bind.WatchOpts, sink chan<- *EthMultiSenderCNativeTx, to []common.Address) (event.Subscription, error) {

	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _EthMultiSenderC.contract.WatchLogs(opts, "NativeTx", toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EthMultiSenderCNativeTx)
				if err := _EthMultiSenderC.contract.UnpackLog(event, "NativeTx", log); err != nil {
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

// ParseNativeTx is a log parse operation binding the contract event 0x311279b6436c79139d7d8a1dc58c6cd2499083a8cc78bfdd8444c5b6fd89ba0d.
//
// Solidity: event NativeTx(address indexed to, uint256 amount)
func (_EthMultiSenderC *EthMultiSenderCFilterer) ParseNativeTx(log types.Log) (*EthMultiSenderCNativeTx, error) {
	event := new(EthMultiSenderCNativeTx)
	if err := _EthMultiSenderC.contract.UnpackLog(event, "NativeTx", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EthMultiSenderCOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the EthMultiSenderC contract.
type EthMultiSenderCOwnershipTransferredIterator struct {
	Event *EthMultiSenderCOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *EthMultiSenderCOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EthMultiSenderCOwnershipTransferred)
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
		it.Event = new(EthMultiSenderCOwnershipTransferred)
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
func (it *EthMultiSenderCOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EthMultiSenderCOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EthMultiSenderCOwnershipTransferred represents a OwnershipTransferred event raised by the EthMultiSenderC contract.
type EthMultiSenderCOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_EthMultiSenderC *EthMultiSenderCFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*EthMultiSenderCOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _EthMultiSenderC.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &EthMultiSenderCOwnershipTransferredIterator{contract: _EthMultiSenderC.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_EthMultiSenderC *EthMultiSenderCFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *EthMultiSenderCOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _EthMultiSenderC.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EthMultiSenderCOwnershipTransferred)
				if err := _EthMultiSenderC.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_EthMultiSenderC *EthMultiSenderCFilterer) ParseOwnershipTransferred(log types.Log) (*EthMultiSenderCOwnershipTransferred, error) {
	event := new(EthMultiSenderCOwnershipTransferred)
	if err := _EthMultiSenderC.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
