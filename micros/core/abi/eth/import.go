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

// EthImportCMetaData contains all meta data concerning the EthImportC contract.
var EthImportCMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_gov\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Imported\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"networkId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Withdraw\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"AUTHENTICATOR\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"BLACK_HOLE\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MANAGER_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"VERSION\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"archive\",\"outputs\":[{\"internalType\":\"contractIImArchive\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_requestId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"_sig\",\"type\":\"bytes\"}],\"name\":\"claim\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"gov\",\"outputs\":[{\"internalType\":\"contractIGovernance\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_archive\",\"type\":\"address\"}],\"name\":\"setArchive\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_gov\",\"type\":\"address\"}],\"name\":\"setGov\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_networkId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"withdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// EthImportCABI is the input ABI used to generate the binding from.
// Deprecated: Use EthImportCMetaData.ABI instead.
var EthImportCABI = EthImportCMetaData.ABI

// EthImportC is an auto generated Go binding around an Ethereum contract.
type EthImportC struct {
	EthImportCCaller     // Read-only binding to the contract
	EthImportCTransactor // Write-only binding to the contract
	EthImportCFilterer   // Log filterer for contract events
}

// EthImportCCaller is an auto generated read-only Go binding around an Ethereum contract.
type EthImportCCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EthImportCTransactor is an auto generated write-only Go binding around an Ethereum contract.
type EthImportCTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EthImportCFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type EthImportCFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EthImportCSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type EthImportCSession struct {
	Contract     *EthImportC       // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// EthImportCCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type EthImportCCallerSession struct {
	Contract *EthImportCCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts     // Call options to use throughout this session
}

// EthImportCTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type EthImportCTransactorSession struct {
	Contract     *EthImportCTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// EthImportCRaw is an auto generated low-level Go binding around an Ethereum contract.
type EthImportCRaw struct {
	Contract *EthImportC // Generic contract binding to access the raw methods on
}

// EthImportCCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type EthImportCCallerRaw struct {
	Contract *EthImportCCaller // Generic read-only contract binding to access the raw methods on
}

// EthImportCTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type EthImportCTransactorRaw struct {
	Contract *EthImportCTransactor // Generic write-only contract binding to access the raw methods on
}

// NewEthImportC creates a new instance of EthImportC, bound to a specific deployed contract.
func NewEthImportC(address common.Address, backend bind.ContractBackend) (*EthImportC, error) {
	contract, err := bindEthImportC(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &EthImportC{EthImportCCaller: EthImportCCaller{contract: contract}, EthImportCTransactor: EthImportCTransactor{contract: contract}, EthImportCFilterer: EthImportCFilterer{contract: contract}}, nil
}

// NewEthImportCCaller creates a new read-only instance of EthImportC, bound to a specific deployed contract.
func NewEthImportCCaller(address common.Address, caller bind.ContractCaller) (*EthImportCCaller, error) {
	contract, err := bindEthImportC(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &EthImportCCaller{contract: contract}, nil
}

// NewEthImportCTransactor creates a new write-only instance of EthImportC, bound to a specific deployed contract.
func NewEthImportCTransactor(address common.Address, transactor bind.ContractTransactor) (*EthImportCTransactor, error) {
	contract, err := bindEthImportC(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &EthImportCTransactor{contract: contract}, nil
}

// NewEthImportCFilterer creates a new log filterer instance of EthImportC, bound to a specific deployed contract.
func NewEthImportCFilterer(address common.Address, filterer bind.ContractFilterer) (*EthImportCFilterer, error) {
	contract, err := bindEthImportC(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &EthImportCFilterer{contract: contract}, nil
}

// bindEthImportC binds a generic wrapper to an already deployed contract.
func bindEthImportC(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(EthImportCABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_EthImportC *EthImportCRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _EthImportC.Contract.EthImportCCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_EthImportC *EthImportCRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _EthImportC.Contract.EthImportCTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_EthImportC *EthImportCRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _EthImportC.Contract.EthImportCTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_EthImportC *EthImportCCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _EthImportC.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_EthImportC *EthImportCTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _EthImportC.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_EthImportC *EthImportCTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _EthImportC.Contract.contract.Transact(opts, method, params...)
}

// AUTHENTICATOR is a free data retrieval call binding the contract method 0xc6186181.
//
// Solidity: function AUTHENTICATOR() view returns(bytes32)
func (_EthImportC *EthImportCCaller) AUTHENTICATOR(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _EthImportC.contract.Call(opts, &out, "AUTHENTICATOR")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// AUTHENTICATOR is a free data retrieval call binding the contract method 0xc6186181.
//
// Solidity: function AUTHENTICATOR() view returns(bytes32)
func (_EthImportC *EthImportCSession) AUTHENTICATOR() ([32]byte, error) {
	return _EthImportC.Contract.AUTHENTICATOR(&_EthImportC.CallOpts)
}

// AUTHENTICATOR is a free data retrieval call binding the contract method 0xc6186181.
//
// Solidity: function AUTHENTICATOR() view returns(bytes32)
func (_EthImportC *EthImportCCallerSession) AUTHENTICATOR() ([32]byte, error) {
	return _EthImportC.Contract.AUTHENTICATOR(&_EthImportC.CallOpts)
}

// BLACKHOLE is a free data retrieval call binding the contract method 0x55eda4e8.
//
// Solidity: function BLACK_HOLE() view returns(address)
func (_EthImportC *EthImportCCaller) BLACKHOLE(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _EthImportC.contract.Call(opts, &out, "BLACK_HOLE")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// BLACKHOLE is a free data retrieval call binding the contract method 0x55eda4e8.
//
// Solidity: function BLACK_HOLE() view returns(address)
func (_EthImportC *EthImportCSession) BLACKHOLE() (common.Address, error) {
	return _EthImportC.Contract.BLACKHOLE(&_EthImportC.CallOpts)
}

// BLACKHOLE is a free data retrieval call binding the contract method 0x55eda4e8.
//
// Solidity: function BLACK_HOLE() view returns(address)
func (_EthImportC *EthImportCCallerSession) BLACKHOLE() (common.Address, error) {
	return _EthImportC.Contract.BLACKHOLE(&_EthImportC.CallOpts)
}

// MANAGERROLE is a free data retrieval call binding the contract method 0xec87621c.
//
// Solidity: function MANAGER_ROLE() view returns(bytes32)
func (_EthImportC *EthImportCCaller) MANAGERROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _EthImportC.contract.Call(opts, &out, "MANAGER_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// MANAGERROLE is a free data retrieval call binding the contract method 0xec87621c.
//
// Solidity: function MANAGER_ROLE() view returns(bytes32)
func (_EthImportC *EthImportCSession) MANAGERROLE() ([32]byte, error) {
	return _EthImportC.Contract.MANAGERROLE(&_EthImportC.CallOpts)
}

// MANAGERROLE is a free data retrieval call binding the contract method 0xec87621c.
//
// Solidity: function MANAGER_ROLE() view returns(bytes32)
func (_EthImportC *EthImportCCallerSession) MANAGERROLE() ([32]byte, error) {
	return _EthImportC.Contract.MANAGERROLE(&_EthImportC.CallOpts)
}

// VERSION is a free data retrieval call binding the contract method 0xffa1ad74.
//
// Solidity: function VERSION() view returns(string)
func (_EthImportC *EthImportCCaller) VERSION(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _EthImportC.contract.Call(opts, &out, "VERSION")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// VERSION is a free data retrieval call binding the contract method 0xffa1ad74.
//
// Solidity: function VERSION() view returns(string)
func (_EthImportC *EthImportCSession) VERSION() (string, error) {
	return _EthImportC.Contract.VERSION(&_EthImportC.CallOpts)
}

// VERSION is a free data retrieval call binding the contract method 0xffa1ad74.
//
// Solidity: function VERSION() view returns(string)
func (_EthImportC *EthImportCCallerSession) VERSION() (string, error) {
	return _EthImportC.Contract.VERSION(&_EthImportC.CallOpts)
}

// Archive is a free data retrieval call binding the contract method 0x02a21460.
//
// Solidity: function archive() view returns(address)
func (_EthImportC *EthImportCCaller) Archive(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _EthImportC.contract.Call(opts, &out, "archive")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Archive is a free data retrieval call binding the contract method 0x02a21460.
//
// Solidity: function archive() view returns(address)
func (_EthImportC *EthImportCSession) Archive() (common.Address, error) {
	return _EthImportC.Contract.Archive(&_EthImportC.CallOpts)
}

// Archive is a free data retrieval call binding the contract method 0x02a21460.
//
// Solidity: function archive() view returns(address)
func (_EthImportC *EthImportCCallerSession) Archive() (common.Address, error) {
	return _EthImportC.Contract.Archive(&_EthImportC.CallOpts)
}

// Gov is a free data retrieval call binding the contract method 0x12d43a51.
//
// Solidity: function gov() view returns(address)
func (_EthImportC *EthImportCCaller) Gov(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _EthImportC.contract.Call(opts, &out, "gov")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Gov is a free data retrieval call binding the contract method 0x12d43a51.
//
// Solidity: function gov() view returns(address)
func (_EthImportC *EthImportCSession) Gov() (common.Address, error) {
	return _EthImportC.Contract.Gov(&_EthImportC.CallOpts)
}

// Gov is a free data retrieval call binding the contract method 0x12d43a51.
//
// Solidity: function gov() view returns(address)
func (_EthImportC *EthImportCCallerSession) Gov() (common.Address, error) {
	return _EthImportC.Contract.Gov(&_EthImportC.CallOpts)
}

// Claim is a paid mutator transaction binding the contract method 0x2ada8a32.
//
// Solidity: function claim(address _token, uint256 _requestId, uint256 _amount, bytes _sig) returns()
func (_EthImportC *EthImportCTransactor) Claim(opts *bind.TransactOpts, _token common.Address, _requestId *big.Int, _amount *big.Int, _sig []byte) (*types.Transaction, error) {
	return _EthImportC.contract.Transact(opts, "claim", _token, _requestId, _amount, _sig)
}

// Claim is a paid mutator transaction binding the contract method 0x2ada8a32.
//
// Solidity: function claim(address _token, uint256 _requestId, uint256 _amount, bytes _sig) returns()
func (_EthImportC *EthImportCSession) Claim(_token common.Address, _requestId *big.Int, _amount *big.Int, _sig []byte) (*types.Transaction, error) {
	return _EthImportC.Contract.Claim(&_EthImportC.TransactOpts, _token, _requestId, _amount, _sig)
}

// Claim is a paid mutator transaction binding the contract method 0x2ada8a32.
//
// Solidity: function claim(address _token, uint256 _requestId, uint256 _amount, bytes _sig) returns()
func (_EthImportC *EthImportCTransactorSession) Claim(_token common.Address, _requestId *big.Int, _amount *big.Int, _sig []byte) (*types.Transaction, error) {
	return _EthImportC.Contract.Claim(&_EthImportC.TransactOpts, _token, _requestId, _amount, _sig)
}

// SetArchive is a paid mutator transaction binding the contract method 0x499dfd71.
//
// Solidity: function setArchive(address _archive) returns()
func (_EthImportC *EthImportCTransactor) SetArchive(opts *bind.TransactOpts, _archive common.Address) (*types.Transaction, error) {
	return _EthImportC.contract.Transact(opts, "setArchive", _archive)
}

// SetArchive is a paid mutator transaction binding the contract method 0x499dfd71.
//
// Solidity: function setArchive(address _archive) returns()
func (_EthImportC *EthImportCSession) SetArchive(_archive common.Address) (*types.Transaction, error) {
	return _EthImportC.Contract.SetArchive(&_EthImportC.TransactOpts, _archive)
}

// SetArchive is a paid mutator transaction binding the contract method 0x499dfd71.
//
// Solidity: function setArchive(address _archive) returns()
func (_EthImportC *EthImportCTransactorSession) SetArchive(_archive common.Address) (*types.Transaction, error) {
	return _EthImportC.Contract.SetArchive(&_EthImportC.TransactOpts, _archive)
}

// SetGov is a paid mutator transaction binding the contract method 0xcfad57a2.
//
// Solidity: function setGov(address _gov) returns()
func (_EthImportC *EthImportCTransactor) SetGov(opts *bind.TransactOpts, _gov common.Address) (*types.Transaction, error) {
	return _EthImportC.contract.Transact(opts, "setGov", _gov)
}

// SetGov is a paid mutator transaction binding the contract method 0xcfad57a2.
//
// Solidity: function setGov(address _gov) returns()
func (_EthImportC *EthImportCSession) SetGov(_gov common.Address) (*types.Transaction, error) {
	return _EthImportC.Contract.SetGov(&_EthImportC.TransactOpts, _gov)
}

// SetGov is a paid mutator transaction binding the contract method 0xcfad57a2.
//
// Solidity: function setGov(address _gov) returns()
func (_EthImportC *EthImportCTransactorSession) SetGov(_gov common.Address) (*types.Transaction, error) {
	return _EthImportC.Contract.SetGov(&_EthImportC.TransactOpts, _gov)
}

// Withdraw is a paid mutator transaction binding the contract method 0x7bfe950c.
//
// Solidity: function withdraw(address _token, address _to, uint256 _networkId, uint256 _value) returns()
func (_EthImportC *EthImportCTransactor) Withdraw(opts *bind.TransactOpts, _token common.Address, _to common.Address, _networkId *big.Int, _value *big.Int) (*types.Transaction, error) {
	return _EthImportC.contract.Transact(opts, "withdraw", _token, _to, _networkId, _value)
}

// Withdraw is a paid mutator transaction binding the contract method 0x7bfe950c.
//
// Solidity: function withdraw(address _token, address _to, uint256 _networkId, uint256 _value) returns()
func (_EthImportC *EthImportCSession) Withdraw(_token common.Address, _to common.Address, _networkId *big.Int, _value *big.Int) (*types.Transaction, error) {
	return _EthImportC.Contract.Withdraw(&_EthImportC.TransactOpts, _token, _to, _networkId, _value)
}

// Withdraw is a paid mutator transaction binding the contract method 0x7bfe950c.
//
// Solidity: function withdraw(address _token, address _to, uint256 _networkId, uint256 _value) returns()
func (_EthImportC *EthImportCTransactorSession) Withdraw(_token common.Address, _to common.Address, _networkId *big.Int, _value *big.Int) (*types.Transaction, error) {
	return _EthImportC.Contract.Withdraw(&_EthImportC.TransactOpts, _token, _to, _networkId, _value)
}

// EthImportCImportedIterator is returned from FilterImported and is used to iterate over the raw logs and unpacked data for Imported events raised by the EthImportC contract.
type EthImportCImportedIterator struct {
	Event *EthImportCImported // Event containing the contract specifics and raw log

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
func (it *EthImportCImportedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EthImportCImported)
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
		it.Event = new(EthImportCImported)
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
func (it *EthImportCImportedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EthImportCImportedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EthImportCImported represents a Imported event raised by the EthImportC contract.
type EthImportCImported struct {
	RequestId *big.Int
	Token     common.Address
	User      common.Address
	Amount    *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterImported is a free log retrieval operation binding the contract event 0x2fe9a6fc5aaedd15e012b68eca99fd27304c05a0bfb0080eef5fbc60c511a02c.
//
// Solidity: event Imported(uint256 indexed requestId, address indexed token, address indexed user, uint256 amount)
func (_EthImportC *EthImportCFilterer) FilterImported(opts *bind.FilterOpts, requestId []*big.Int, token []common.Address, user []common.Address) (*EthImportCImportedIterator, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}
	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}
	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _EthImportC.contract.FilterLogs(opts, "Imported", requestIdRule, tokenRule, userRule)
	if err != nil {
		return nil, err
	}
	return &EthImportCImportedIterator{contract: _EthImportC.contract, event: "Imported", logs: logs, sub: sub}, nil
}

// WatchImported is a free log subscription operation binding the contract event 0x2fe9a6fc5aaedd15e012b68eca99fd27304c05a0bfb0080eef5fbc60c511a02c.
//
// Solidity: event Imported(uint256 indexed requestId, address indexed token, address indexed user, uint256 amount)
func (_EthImportC *EthImportCFilterer) WatchImported(opts *bind.WatchOpts, sink chan<- *EthImportCImported, requestId []*big.Int, token []common.Address, user []common.Address) (event.Subscription, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}
	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}
	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _EthImportC.contract.WatchLogs(opts, "Imported", requestIdRule, tokenRule, userRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EthImportCImported)
				if err := _EthImportC.contract.UnpackLog(event, "Imported", log); err != nil {
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

// ParseImported is a log parse operation binding the contract event 0x2fe9a6fc5aaedd15e012b68eca99fd27304c05a0bfb0080eef5fbc60c511a02c.
//
// Solidity: event Imported(uint256 indexed requestId, address indexed token, address indexed user, uint256 amount)
func (_EthImportC *EthImportCFilterer) ParseImported(log types.Log) (*EthImportCImported, error) {
	event := new(EthImportCImported)
	if err := _EthImportC.contract.UnpackLog(event, "Imported", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EthImportCWithdrawIterator is returned from FilterWithdraw and is used to iterate over the raw logs and unpacked data for Withdraw events raised by the EthImportC contract.
type EthImportCWithdrawIterator struct {
	Event *EthImportCWithdraw // Event containing the contract specifics and raw log

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
func (it *EthImportCWithdrawIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EthImportCWithdraw)
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
		it.Event = new(EthImportCWithdraw)
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
func (it *EthImportCWithdrawIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EthImportCWithdrawIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EthImportCWithdraw represents a Withdraw event raised by the EthImportC contract.
type EthImportCWithdraw struct {
	Token     common.Address
	From      common.Address
	To        common.Address
	NetworkId *big.Int
	Amount    *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterWithdraw is a free log retrieval operation binding the contract event 0xfbde797d201c681b91056529119e0b02407c7bb96a4a2c75c01fc9667232c8db.
//
// Solidity: event Withdraw(address indexed token, address indexed from, address indexed to, uint256 networkId, uint256 amount)
func (_EthImportC *EthImportCFilterer) FilterWithdraw(opts *bind.FilterOpts, token []common.Address, from []common.Address, to []common.Address) (*EthImportCWithdrawIterator, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _EthImportC.contract.FilterLogs(opts, "Withdraw", tokenRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &EthImportCWithdrawIterator{contract: _EthImportC.contract, event: "Withdraw", logs: logs, sub: sub}, nil
}

// WatchWithdraw is a free log subscription operation binding the contract event 0xfbde797d201c681b91056529119e0b02407c7bb96a4a2c75c01fc9667232c8db.
//
// Solidity: event Withdraw(address indexed token, address indexed from, address indexed to, uint256 networkId, uint256 amount)
func (_EthImportC *EthImportCFilterer) WatchWithdraw(opts *bind.WatchOpts, sink chan<- *EthImportCWithdraw, token []common.Address, from []common.Address, to []common.Address) (event.Subscription, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _EthImportC.contract.WatchLogs(opts, "Withdraw", tokenRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EthImportCWithdraw)
				if err := _EthImportC.contract.UnpackLog(event, "Withdraw", log); err != nil {
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

// ParseWithdraw is a log parse operation binding the contract event 0xfbde797d201c681b91056529119e0b02407c7bb96a4a2c75c01fc9667232c8db.
//
// Solidity: event Withdraw(address indexed token, address indexed from, address indexed to, uint256 networkId, uint256 amount)
func (_EthImportC *EthImportCFilterer) ParseWithdraw(log types.Log) (*EthImportCWithdraw, error) {
	event := new(EthImportCWithdraw)
	if err := _EthImportC.contract.UnpackLog(event, "Withdraw", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
