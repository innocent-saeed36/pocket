package modules

import (
	cryptoPocket "github.com/pokt-network/pocket/shared/crypto"
)

type Module interface {
	InitializableModule
	IntegratableModule
	InterruptableModule
}

// IntegratableModule is a module that integrates with the bus.
// Essentially it's a module that is capable of communicating with the `bus` (see `shared/modules/bus_module.go`) for additional details.
type IntegratableModule interface {
	SetBus(Bus)
	GetBus() Bus
}

// InitializableModule is a module that has some basic lifecycle logic. Specifically, it can be started and stopped.
type InterruptableModule interface {
	Start() error
	Stop() error
}

// ModuleOption is a function that configures a module when it is created.
// It uses a widely used pattern in Go called functional options.
// See https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis
// for more information.
//
// It is used to provide optional parameters to the module constructor for all the cases
// where there is no configuration, which is often the case for sub-modules that are used
// and configured at runtime.
//
// It accepts an InitializableModule as a parameter, because in order to create a module with these options,
// at a minimum, the module must implement the InitializableModule interface.
//
// Example:
//
//	func WithFoo(foo string) ModuleOption {
//	  return func(m InitializableModule) {
//	    m.(*MyModule).foo = foo
//	  }
//	}
//
//	func NewMyModule(options ...ModuleOption) (Module, error) {
//	  m := &MyModule{}
//	  for _, option := range options {
//	    option(m)
//	  }
//	  return m, nil
//	}
type ModuleOption func(InitializableModule)

// InitializableModule is a module that can be created via the standardized `Create` method and that has a name
// that can be used to identify it (see `shared\modules\modules_registry_module.go`) for additional details.
type InitializableModule interface {
	GetModuleName() string
	Create(bus Bus, options ...ModuleOption) (Module, error)
}

// KeyholderModule is a module that can provide a private key.
type KeyholderModule interface {
	GetPrivateKey() (cryptoPocket.PrivateKey, error)
}

// P2PAddressableModule is a module that can provide a P2P address.
type P2PAddressableModule interface {
	GetP2PAddress() cryptoPocket.Address
}

// ObservableModule is a module that can provide observability via a Logger.
type ObservableModule interface {
	GetLogger() Logger
}
