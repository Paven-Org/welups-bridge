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

// EthGovMetaData contains all meta data concerning the EthGov contract.
var EthGovMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_superAdmin\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"previousAdminRole\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"newAdminRole\",\"type\":\"bytes32\"}],\"name\":\"RoleAdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleGranted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleRevoked\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"AUTHENTICATOR\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"DEFAULT_ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MANAGER_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_networkId\",\"type\":\"uint256\"}],\"name\":\"addNetwork\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_opt\",\"type\":\"uint256\"}],\"name\":\"addToList\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"eService\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"exports\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"}],\"name\":\"getRoleAdmin\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"}],\"name\":\"getRoleMember\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"}],\"name\":\"getRoleMemberCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"grantRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"hasRole\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"iService\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"imports\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"locked\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"networks\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_opt\",\"type\":\"uint256\"}],\"name\":\"removeFromList\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_networkId\",\"type\":\"uint256\"}],\"name\":\"removeNetwork\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"renounceRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"revokeRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bool\",\"name\":\"_state\",\"type\":\"bool\"}],\"name\":\"setLock\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_export\",\"type\":\"address\"}],\"name\":\"updateEService\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_import\",\"type\":\"address\"}],\"name\":\"updateIService\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// EthGovABI is the input ABI used to generate the binding from.
// Deprecated: Use EthGovMetaData.ABI instead.
var EthGovABI = EthGovMetaData.ABI

// EthGov is an auto generated Go binding around an Ethereum contract.
type EthGov struct {
	EthGovCaller     // Read-only binding to the contract
	EthGovTransactor // Write-only binding to the contract
	EthGovFilterer   // Log filterer for contract events
}

// EthGovCaller is an auto generated read-only Go binding around an Ethereum contract.
type EthGovCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EthGovTransactor is an auto generated write-only Go binding around an Ethereum contract.
type EthGovTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EthGovFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type EthGovFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EthGovSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type EthGovSession struct {
	Contract     *EthGov           // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// EthGovCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type EthGovCallerSession struct {
	Contract *EthGovCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// EthGovTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type EthGovTransactorSession struct {
	Contract     *EthGovTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// EthGovRaw is an auto generated low-level Go binding around an Ethereum contract.
type EthGovRaw struct {
	Contract *EthGov // Generic contract binding to access the raw methods on
}

// EthGovCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type EthGovCallerRaw struct {
	Contract *EthGovCaller // Generic read-only contract binding to access the raw methods on
}

// EthGovTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type EthGovTransactorRaw struct {
	Contract *EthGovTransactor // Generic write-only contract binding to access the raw methods on
}

// NewEthGov creates a new instance of EthGov, bound to a specific deployed contract.
func NewEthGov(address common.Address, backend bind.ContractBackend) (*EthGov, error) {
	contract, err := bindEthGov(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &EthGov{EthGovCaller: EthGovCaller{contract: contract}, EthGovTransactor: EthGovTransactor{contract: contract}, EthGovFilterer: EthGovFilterer{contract: contract}}, nil
}

// NewEthGovCaller creates a new read-only instance of EthGov, bound to a specific deployed contract.
func NewEthGovCaller(address common.Address, caller bind.ContractCaller) (*EthGovCaller, error) {
	contract, err := bindEthGov(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &EthGovCaller{contract: contract}, nil
}

// NewEthGovTransactor creates a new write-only instance of EthGov, bound to a specific deployed contract.
func NewEthGovTransactor(address common.Address, transactor bind.ContractTransactor) (*EthGovTransactor, error) {
	contract, err := bindEthGov(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &EthGovTransactor{contract: contract}, nil
}

// NewEthGovFilterer creates a new log filterer instance of EthGov, bound to a specific deployed contract.
func NewEthGovFilterer(address common.Address, filterer bind.ContractFilterer) (*EthGovFilterer, error) {
	contract, err := bindEthGov(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &EthGovFilterer{contract: contract}, nil
}

// bindEthGov binds a generic wrapper to an already deployed contract.
func bindEthGov(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(EthGovABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_EthGov *EthGovRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _EthGov.Contract.EthGovCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_EthGov *EthGovRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _EthGov.Contract.EthGovTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_EthGov *EthGovRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _EthGov.Contract.EthGovTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_EthGov *EthGovCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _EthGov.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_EthGov *EthGovTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _EthGov.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_EthGov *EthGovTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _EthGov.Contract.contract.Transact(opts, method, params...)
}

// AUTHENTICATOR is a free data retrieval call binding the contract method 0xc6186181.
//
// Solidity: function AUTHENTICATOR() view returns(bytes32)
func (_EthGov *EthGovCaller) AUTHENTICATOR(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _EthGov.contract.Call(opts, &out, "AUTHENTICATOR")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// AUTHENTICATOR is a free data retrieval call binding the contract method 0xc6186181.
//
// Solidity: function AUTHENTICATOR() view returns(bytes32)
func (_EthGov *EthGovSession) AUTHENTICATOR() ([32]byte, error) {
	return _EthGov.Contract.AUTHENTICATOR(&_EthGov.CallOpts)
}

// AUTHENTICATOR is a free data retrieval call binding the contract method 0xc6186181.
//
// Solidity: function AUTHENTICATOR() view returns(bytes32)
func (_EthGov *EthGovCallerSession) AUTHENTICATOR() ([32]byte, error) {
	return _EthGov.Contract.AUTHENTICATOR(&_EthGov.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_EthGov *EthGovCaller) DEFAULTADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _EthGov.contract.Call(opts, &out, "DEFAULT_ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_EthGov *EthGovSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _EthGov.Contract.DEFAULTADMINROLE(&_EthGov.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_EthGov *EthGovCallerSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _EthGov.Contract.DEFAULTADMINROLE(&_EthGov.CallOpts)
}

// MANAGERROLE is a free data retrieval call binding the contract method 0xec87621c.
//
// Solidity: function MANAGER_ROLE() view returns(bytes32)
func (_EthGov *EthGovCaller) MANAGERROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _EthGov.contract.Call(opts, &out, "MANAGER_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// MANAGERROLE is a free data retrieval call binding the contract method 0xec87621c.
//
// Solidity: function MANAGER_ROLE() view returns(bytes32)
func (_EthGov *EthGovSession) MANAGERROLE() ([32]byte, error) {
	return _EthGov.Contract.MANAGERROLE(&_EthGov.CallOpts)
}

// MANAGERROLE is a free data retrieval call binding the contract method 0xec87621c.
//
// Solidity: function MANAGER_ROLE() view returns(bytes32)
func (_EthGov *EthGovCallerSession) MANAGERROLE() ([32]byte, error) {
	return _EthGov.Contract.MANAGERROLE(&_EthGov.CallOpts)
}

// EService is a free data retrieval call binding the contract method 0x1f1c4a13.
//
// Solidity: function eService() view returns(address)
func (_EthGov *EthGovCaller) EService(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _EthGov.contract.Call(opts, &out, "eService")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// EService is a free data retrieval call binding the contract method 0x1f1c4a13.
//
// Solidity: function eService() view returns(address)
func (_EthGov *EthGovSession) EService() (common.Address, error) {
	return _EthGov.Contract.EService(&_EthGov.CallOpts)
}

// EService is a free data retrieval call binding the contract method 0x1f1c4a13.
//
// Solidity: function eService() view returns(address)
func (_EthGov *EthGovCallerSession) EService() (common.Address, error) {
	return _EthGov.Contract.EService(&_EthGov.CallOpts)
}

// Exports is a free data retrieval call binding the contract method 0x935aa419.
//
// Solidity: function exports(address ) view returns(bool)
func (_EthGov *EthGovCaller) Exports(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var out []interface{}
	err := _EthGov.contract.Call(opts, &out, "exports", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Exports is a free data retrieval call binding the contract method 0x935aa419.
//
// Solidity: function exports(address ) view returns(bool)
func (_EthGov *EthGovSession) Exports(arg0 common.Address) (bool, error) {
	return _EthGov.Contract.Exports(&_EthGov.CallOpts, arg0)
}

// Exports is a free data retrieval call binding the contract method 0x935aa419.
//
// Solidity: function exports(address ) view returns(bool)
func (_EthGov *EthGovCallerSession) Exports(arg0 common.Address) (bool, error) {
	return _EthGov.Contract.Exports(&_EthGov.CallOpts, arg0)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_EthGov *EthGovCaller) GetRoleAdmin(opts *bind.CallOpts, role [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _EthGov.contract.Call(opts, &out, "getRoleAdmin", role)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_EthGov *EthGovSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _EthGov.Contract.GetRoleAdmin(&_EthGov.CallOpts, role)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_EthGov *EthGovCallerSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _EthGov.Contract.GetRoleAdmin(&_EthGov.CallOpts, role)
}

// GetRoleMember is a free data retrieval call binding the contract method 0x9010d07c.
//
// Solidity: function getRoleMember(bytes32 role, uint256 index) view returns(address)
func (_EthGov *EthGovCaller) GetRoleMember(opts *bind.CallOpts, role [32]byte, index *big.Int) (common.Address, error) {
	var out []interface{}
	err := _EthGov.contract.Call(opts, &out, "getRoleMember", role, index)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetRoleMember is a free data retrieval call binding the contract method 0x9010d07c.
//
// Solidity: function getRoleMember(bytes32 role, uint256 index) view returns(address)
func (_EthGov *EthGovSession) GetRoleMember(role [32]byte, index *big.Int) (common.Address, error) {
	return _EthGov.Contract.GetRoleMember(&_EthGov.CallOpts, role, index)
}

// GetRoleMember is a free data retrieval call binding the contract method 0x9010d07c.
//
// Solidity: function getRoleMember(bytes32 role, uint256 index) view returns(address)
func (_EthGov *EthGovCallerSession) GetRoleMember(role [32]byte, index *big.Int) (common.Address, error) {
	return _EthGov.Contract.GetRoleMember(&_EthGov.CallOpts, role, index)
}

// GetRoleMemberCount is a free data retrieval call binding the contract method 0xca15c873.
//
// Solidity: function getRoleMemberCount(bytes32 role) view returns(uint256)
func (_EthGov *EthGovCaller) GetRoleMemberCount(opts *bind.CallOpts, role [32]byte) (*big.Int, error) {
	var out []interface{}
	err := _EthGov.contract.Call(opts, &out, "getRoleMemberCount", role)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetRoleMemberCount is a free data retrieval call binding the contract method 0xca15c873.
//
// Solidity: function getRoleMemberCount(bytes32 role) view returns(uint256)
func (_EthGov *EthGovSession) GetRoleMemberCount(role [32]byte) (*big.Int, error) {
	return _EthGov.Contract.GetRoleMemberCount(&_EthGov.CallOpts, role)
}

// GetRoleMemberCount is a free data retrieval call binding the contract method 0xca15c873.
//
// Solidity: function getRoleMemberCount(bytes32 role) view returns(uint256)
func (_EthGov *EthGovCallerSession) GetRoleMemberCount(role [32]byte) (*big.Int, error) {
	return _EthGov.Contract.GetRoleMemberCount(&_EthGov.CallOpts, role)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_EthGov *EthGovCaller) HasRole(opts *bind.CallOpts, role [32]byte, account common.Address) (bool, error) {
	var out []interface{}
	err := _EthGov.contract.Call(opts, &out, "hasRole", role, account)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_EthGov *EthGovSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _EthGov.Contract.HasRole(&_EthGov.CallOpts, role, account)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_EthGov *EthGovCallerSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _EthGov.Contract.HasRole(&_EthGov.CallOpts, role, account)
}

// IService is a free data retrieval call binding the contract method 0xf603bc8a.
//
// Solidity: function iService() view returns(address)
func (_EthGov *EthGovCaller) IService(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _EthGov.contract.Call(opts, &out, "iService")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// IService is a free data retrieval call binding the contract method 0xf603bc8a.
//
// Solidity: function iService() view returns(address)
func (_EthGov *EthGovSession) IService() (common.Address, error) {
	return _EthGov.Contract.IService(&_EthGov.CallOpts)
}

// IService is a free data retrieval call binding the contract method 0xf603bc8a.
//
// Solidity: function iService() view returns(address)
func (_EthGov *EthGovCallerSession) IService() (common.Address, error) {
	return _EthGov.Contract.IService(&_EthGov.CallOpts)
}

// Imports is a free data retrieval call binding the contract method 0xf7e0d46f.
//
// Solidity: function imports(address ) view returns(bool)
func (_EthGov *EthGovCaller) Imports(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var out []interface{}
	err := _EthGov.contract.Call(opts, &out, "imports", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Imports is a free data retrieval call binding the contract method 0xf7e0d46f.
//
// Solidity: function imports(address ) view returns(bool)
func (_EthGov *EthGovSession) Imports(arg0 common.Address) (bool, error) {
	return _EthGov.Contract.Imports(&_EthGov.CallOpts, arg0)
}

// Imports is a free data retrieval call binding the contract method 0xf7e0d46f.
//
// Solidity: function imports(address ) view returns(bool)
func (_EthGov *EthGovCallerSession) Imports(arg0 common.Address) (bool, error) {
	return _EthGov.Contract.Imports(&_EthGov.CallOpts, arg0)
}

// Locked is a free data retrieval call binding the contract method 0xcf309012.
//
// Solidity: function locked() view returns(bool)
func (_EthGov *EthGovCaller) Locked(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _EthGov.contract.Call(opts, &out, "locked")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Locked is a free data retrieval call binding the contract method 0xcf309012.
//
// Solidity: function locked() view returns(bool)
func (_EthGov *EthGovSession) Locked() (bool, error) {
	return _EthGov.Contract.Locked(&_EthGov.CallOpts)
}

// Locked is a free data retrieval call binding the contract method 0xcf309012.
//
// Solidity: function locked() view returns(bool)
func (_EthGov *EthGovCallerSession) Locked() (bool, error) {
	return _EthGov.Contract.Locked(&_EthGov.CallOpts)
}

// Networks is a free data retrieval call binding the contract method 0x8bb0a17c.
//
// Solidity: function networks(uint256 ) view returns(bool)
func (_EthGov *EthGovCaller) Networks(opts *bind.CallOpts, arg0 *big.Int) (bool, error) {
	var out []interface{}
	err := _EthGov.contract.Call(opts, &out, "networks", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Networks is a free data retrieval call binding the contract method 0x8bb0a17c.
//
// Solidity: function networks(uint256 ) view returns(bool)
func (_EthGov *EthGovSession) Networks(arg0 *big.Int) (bool, error) {
	return _EthGov.Contract.Networks(&_EthGov.CallOpts, arg0)
}

// Networks is a free data retrieval call binding the contract method 0x8bb0a17c.
//
// Solidity: function networks(uint256 ) view returns(bool)
func (_EthGov *EthGovCallerSession) Networks(arg0 *big.Int) (bool, error) {
	return _EthGov.Contract.Networks(&_EthGov.CallOpts, arg0)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_EthGov *EthGovCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _EthGov.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_EthGov *EthGovSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _EthGov.Contract.SupportsInterface(&_EthGov.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_EthGov *EthGovCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _EthGov.Contract.SupportsInterface(&_EthGov.CallOpts, interfaceId)
}

// AddNetwork is a paid mutator transaction binding the contract method 0xe5d36ede.
//
// Solidity: function addNetwork(uint256 _networkId) returns()
func (_EthGov *EthGovTransactor) AddNetwork(opts *bind.TransactOpts, _networkId *big.Int) (*types.Transaction, error) {
	return _EthGov.contract.Transact(opts, "addNetwork", _networkId)
}

// AddNetwork is a paid mutator transaction binding the contract method 0xe5d36ede.
//
// Solidity: function addNetwork(uint256 _networkId) returns()
func (_EthGov *EthGovSession) AddNetwork(_networkId *big.Int) (*types.Transaction, error) {
	return _EthGov.Contract.AddNetwork(&_EthGov.TransactOpts, _networkId)
}

// AddNetwork is a paid mutator transaction binding the contract method 0xe5d36ede.
//
// Solidity: function addNetwork(uint256 _networkId) returns()
func (_EthGov *EthGovTransactorSession) AddNetwork(_networkId *big.Int) (*types.Transaction, error) {
	return _EthGov.Contract.AddNetwork(&_EthGov.TransactOpts, _networkId)
}

// AddToList is a paid mutator transaction binding the contract method 0xda2dca10.
//
// Solidity: function addToList(address _token, uint256 _opt) returns()
func (_EthGov *EthGovTransactor) AddToList(opts *bind.TransactOpts, _token common.Address, _opt *big.Int) (*types.Transaction, error) {
	return _EthGov.contract.Transact(opts, "addToList", _token, _opt)
}

// AddToList is a paid mutator transaction binding the contract method 0xda2dca10.
//
// Solidity: function addToList(address _token, uint256 _opt) returns()
func (_EthGov *EthGovSession) AddToList(_token common.Address, _opt *big.Int) (*types.Transaction, error) {
	return _EthGov.Contract.AddToList(&_EthGov.TransactOpts, _token, _opt)
}

// AddToList is a paid mutator transaction binding the contract method 0xda2dca10.
//
// Solidity: function addToList(address _token, uint256 _opt) returns()
func (_EthGov *EthGovTransactorSession) AddToList(_token common.Address, _opt *big.Int) (*types.Transaction, error) {
	return _EthGov.Contract.AddToList(&_EthGov.TransactOpts, _token, _opt)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_EthGov *EthGovTransactor) GrantRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _EthGov.contract.Transact(opts, "grantRole", role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_EthGov *EthGovSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _EthGov.Contract.GrantRole(&_EthGov.TransactOpts, role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_EthGov *EthGovTransactorSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _EthGov.Contract.GrantRole(&_EthGov.TransactOpts, role, account)
}

// RemoveFromList is a paid mutator transaction binding the contract method 0xdfaa21d0.
//
// Solidity: function removeFromList(address _token, uint256 _opt) returns()
func (_EthGov *EthGovTransactor) RemoveFromList(opts *bind.TransactOpts, _token common.Address, _opt *big.Int) (*types.Transaction, error) {
	return _EthGov.contract.Transact(opts, "removeFromList", _token, _opt)
}

// RemoveFromList is a paid mutator transaction binding the contract method 0xdfaa21d0.
//
// Solidity: function removeFromList(address _token, uint256 _opt) returns()
func (_EthGov *EthGovSession) RemoveFromList(_token common.Address, _opt *big.Int) (*types.Transaction, error) {
	return _EthGov.Contract.RemoveFromList(&_EthGov.TransactOpts, _token, _opt)
}

// RemoveFromList is a paid mutator transaction binding the contract method 0xdfaa21d0.
//
// Solidity: function removeFromList(address _token, uint256 _opt) returns()
func (_EthGov *EthGovTransactorSession) RemoveFromList(_token common.Address, _opt *big.Int) (*types.Transaction, error) {
	return _EthGov.Contract.RemoveFromList(&_EthGov.TransactOpts, _token, _opt)
}

// RemoveNetwork is a paid mutator transaction binding the contract method 0xb74d4c05.
//
// Solidity: function removeNetwork(uint256 _networkId) returns()
func (_EthGov *EthGovTransactor) RemoveNetwork(opts *bind.TransactOpts, _networkId *big.Int) (*types.Transaction, error) {
	return _EthGov.contract.Transact(opts, "removeNetwork", _networkId)
}

// RemoveNetwork is a paid mutator transaction binding the contract method 0xb74d4c05.
//
// Solidity: function removeNetwork(uint256 _networkId) returns()
func (_EthGov *EthGovSession) RemoveNetwork(_networkId *big.Int) (*types.Transaction, error) {
	return _EthGov.Contract.RemoveNetwork(&_EthGov.TransactOpts, _networkId)
}

// RemoveNetwork is a paid mutator transaction binding the contract method 0xb74d4c05.
//
// Solidity: function removeNetwork(uint256 _networkId) returns()
func (_EthGov *EthGovTransactorSession) RemoveNetwork(_networkId *big.Int) (*types.Transaction, error) {
	return _EthGov.Contract.RemoveNetwork(&_EthGov.TransactOpts, _networkId)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_EthGov *EthGovTransactor) RenounceRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _EthGov.contract.Transact(opts, "renounceRole", role, account)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_EthGov *EthGovSession) RenounceRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _EthGov.Contract.RenounceRole(&_EthGov.TransactOpts, role, account)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_EthGov *EthGovTransactorSession) RenounceRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _EthGov.Contract.RenounceRole(&_EthGov.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_EthGov *EthGovTransactor) RevokeRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _EthGov.contract.Transact(opts, "revokeRole", role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_EthGov *EthGovSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _EthGov.Contract.RevokeRole(&_EthGov.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_EthGov *EthGovTransactorSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _EthGov.Contract.RevokeRole(&_EthGov.TransactOpts, role, account)
}

// SetLock is a paid mutator transaction binding the contract method 0x619d5194.
//
// Solidity: function setLock(bool _state) returns()
func (_EthGov *EthGovTransactor) SetLock(opts *bind.TransactOpts, _state bool) (*types.Transaction, error) {
	return _EthGov.contract.Transact(opts, "setLock", _state)
}

// SetLock is a paid mutator transaction binding the contract method 0x619d5194.
//
// Solidity: function setLock(bool _state) returns()
func (_EthGov *EthGovSession) SetLock(_state bool) (*types.Transaction, error) {
	return _EthGov.Contract.SetLock(&_EthGov.TransactOpts, _state)
}

// SetLock is a paid mutator transaction binding the contract method 0x619d5194.
//
// Solidity: function setLock(bool _state) returns()
func (_EthGov *EthGovTransactorSession) SetLock(_state bool) (*types.Transaction, error) {
	return _EthGov.Contract.SetLock(&_EthGov.TransactOpts, _state)
}

// UpdateEService is a paid mutator transaction binding the contract method 0xb77a014f.
//
// Solidity: function updateEService(address _export) returns()
func (_EthGov *EthGovTransactor) UpdateEService(opts *bind.TransactOpts, _export common.Address) (*types.Transaction, error) {
	return _EthGov.contract.Transact(opts, "updateEService", _export)
}

// UpdateEService is a paid mutator transaction binding the contract method 0xb77a014f.
//
// Solidity: function updateEService(address _export) returns()
func (_EthGov *EthGovSession) UpdateEService(_export common.Address) (*types.Transaction, error) {
	return _EthGov.Contract.UpdateEService(&_EthGov.TransactOpts, _export)
}

// UpdateEService is a paid mutator transaction binding the contract method 0xb77a014f.
//
// Solidity: function updateEService(address _export) returns()
func (_EthGov *EthGovTransactorSession) UpdateEService(_export common.Address) (*types.Transaction, error) {
	return _EthGov.Contract.UpdateEService(&_EthGov.TransactOpts, _export)
}

// UpdateIService is a paid mutator transaction binding the contract method 0x1a0c8a7a.
//
// Solidity: function updateIService(address _import) returns()
func (_EthGov *EthGovTransactor) UpdateIService(opts *bind.TransactOpts, _import common.Address) (*types.Transaction, error) {
	return _EthGov.contract.Transact(opts, "updateIService", _import)
}

// UpdateIService is a paid mutator transaction binding the contract method 0x1a0c8a7a.
//
// Solidity: function updateIService(address _import) returns()
func (_EthGov *EthGovSession) UpdateIService(_import common.Address) (*types.Transaction, error) {
	return _EthGov.Contract.UpdateIService(&_EthGov.TransactOpts, _import)
}

// UpdateIService is a paid mutator transaction binding the contract method 0x1a0c8a7a.
//
// Solidity: function updateIService(address _import) returns()
func (_EthGov *EthGovTransactorSession) UpdateIService(_import common.Address) (*types.Transaction, error) {
	return _EthGov.Contract.UpdateIService(&_EthGov.TransactOpts, _import)
}

// EthGovRoleAdminChangedIterator is returned from FilterRoleAdminChanged and is used to iterate over the raw logs and unpacked data for RoleAdminChanged events raised by the EthGov contract.
type EthGovRoleAdminChangedIterator struct {
	Event *EthGovRoleAdminChanged // Event containing the contract specifics and raw log

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
func (it *EthGovRoleAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EthGovRoleAdminChanged)
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
		it.Event = new(EthGovRoleAdminChanged)
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
func (it *EthGovRoleAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EthGovRoleAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EthGovRoleAdminChanged represents a RoleAdminChanged event raised by the EthGov contract.
type EthGovRoleAdminChanged struct {
	Role              [32]byte
	PreviousAdminRole [32]byte
	NewAdminRole      [32]byte
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterRoleAdminChanged is a free log retrieval operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_EthGov *EthGovFilterer) FilterRoleAdminChanged(opts *bind.FilterOpts, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (*EthGovRoleAdminChangedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var previousAdminRoleRule []interface{}
	for _, previousAdminRoleItem := range previousAdminRole {
		previousAdminRoleRule = append(previousAdminRoleRule, previousAdminRoleItem)
	}
	var newAdminRoleRule []interface{}
	for _, newAdminRoleItem := range newAdminRole {
		newAdminRoleRule = append(newAdminRoleRule, newAdminRoleItem)
	}

	logs, sub, err := _EthGov.contract.FilterLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return &EthGovRoleAdminChangedIterator{contract: _EthGov.contract, event: "RoleAdminChanged", logs: logs, sub: sub}, nil
}

// WatchRoleAdminChanged is a free log subscription operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_EthGov *EthGovFilterer) WatchRoleAdminChanged(opts *bind.WatchOpts, sink chan<- *EthGovRoleAdminChanged, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var previousAdminRoleRule []interface{}
	for _, previousAdminRoleItem := range previousAdminRole {
		previousAdminRoleRule = append(previousAdminRoleRule, previousAdminRoleItem)
	}
	var newAdminRoleRule []interface{}
	for _, newAdminRoleItem := range newAdminRole {
		newAdminRoleRule = append(newAdminRoleRule, newAdminRoleItem)
	}

	logs, sub, err := _EthGov.contract.WatchLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EthGovRoleAdminChanged)
				if err := _EthGov.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
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

// ParseRoleAdminChanged is a log parse operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_EthGov *EthGovFilterer) ParseRoleAdminChanged(log types.Log) (*EthGovRoleAdminChanged, error) {
	event := new(EthGovRoleAdminChanged)
	if err := _EthGov.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EthGovRoleGrantedIterator is returned from FilterRoleGranted and is used to iterate over the raw logs and unpacked data for RoleGranted events raised by the EthGov contract.
type EthGovRoleGrantedIterator struct {
	Event *EthGovRoleGranted // Event containing the contract specifics and raw log

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
func (it *EthGovRoleGrantedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EthGovRoleGranted)
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
		it.Event = new(EthGovRoleGranted)
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
func (it *EthGovRoleGrantedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EthGovRoleGrantedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EthGovRoleGranted represents a RoleGranted event raised by the EthGov contract.
type EthGovRoleGranted struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleGranted is a free log retrieval operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_EthGov *EthGovFilterer) FilterRoleGranted(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*EthGovRoleGrantedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _EthGov.contract.FilterLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &EthGovRoleGrantedIterator{contract: _EthGov.contract, event: "RoleGranted", logs: logs, sub: sub}, nil
}

// WatchRoleGranted is a free log subscription operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_EthGov *EthGovFilterer) WatchRoleGranted(opts *bind.WatchOpts, sink chan<- *EthGovRoleGranted, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _EthGov.contract.WatchLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EthGovRoleGranted)
				if err := _EthGov.contract.UnpackLog(event, "RoleGranted", log); err != nil {
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

// ParseRoleGranted is a log parse operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_EthGov *EthGovFilterer) ParseRoleGranted(log types.Log) (*EthGovRoleGranted, error) {
	event := new(EthGovRoleGranted)
	if err := _EthGov.contract.UnpackLog(event, "RoleGranted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EthGovRoleRevokedIterator is returned from FilterRoleRevoked and is used to iterate over the raw logs and unpacked data for RoleRevoked events raised by the EthGov contract.
type EthGovRoleRevokedIterator struct {
	Event *EthGovRoleRevoked // Event containing the contract specifics and raw log

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
func (it *EthGovRoleRevokedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EthGovRoleRevoked)
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
		it.Event = new(EthGovRoleRevoked)
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
func (it *EthGovRoleRevokedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EthGovRoleRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EthGovRoleRevoked represents a RoleRevoked event raised by the EthGov contract.
type EthGovRoleRevoked struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleRevoked is a free log retrieval operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_EthGov *EthGovFilterer) FilterRoleRevoked(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*EthGovRoleRevokedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _EthGov.contract.FilterLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &EthGovRoleRevokedIterator{contract: _EthGov.contract, event: "RoleRevoked", logs: logs, sub: sub}, nil
}

// WatchRoleRevoked is a free log subscription operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_EthGov *EthGovFilterer) WatchRoleRevoked(opts *bind.WatchOpts, sink chan<- *EthGovRoleRevoked, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _EthGov.contract.WatchLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EthGovRoleRevoked)
				if err := _EthGov.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
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

// ParseRoleRevoked is a log parse operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_EthGov *EthGovFilterer) ParseRoleRevoked(log types.Log) (*EthGovRoleRevoked, error) {
	event := new(EthGovRoleRevoked)
	if err := _EthGov.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
