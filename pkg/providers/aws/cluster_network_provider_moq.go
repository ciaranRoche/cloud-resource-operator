// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package aws

import (
	"context"
	"net"
	"sync"
)

var (
	lockNetworkManagerMockCreateNetwork            sync.RWMutex
	lockNetworkManagerMockCreateNetworkConnection  sync.RWMutex
	lockNetworkManagerMockCreateNetworkPeering     sync.RWMutex
	lockNetworkManagerMockDeleteNetwork            sync.RWMutex
	lockNetworkManagerMockDeleteNetworkConnection  sync.RWMutex
	lockNetworkManagerMockDeleteNetworkPeering     sync.RWMutex
	lockNetworkManagerMockGetClusterNetworkPeering sync.RWMutex
	lockNetworkManagerMockIsEnabled                sync.RWMutex
)

// Ensure, that NetworkManagerMock does implement NetworkManager.
// If this is not the case, regenerate this file with moq.
var _ NetworkManager = &NetworkManagerMock{}

// NetworkManagerMock is a mock implementation of NetworkManager.
//
//     func TestSomethingThatUsesNetworkManager(t *testing.T) {
//
//         // make and configure a mocked NetworkManager
//         mockedNetworkManager := &NetworkManagerMock{
//             CreateNetworkFunc: func(in1 context.Context, in2 *net.IPNet) (*Network, error) {
// 	               panic("mock out the CreateNetwork method")
//             },
//             CreateNetworkConnectionFunc: func(in1 context.Context, in2 *Network) (*NetworkConnection, error) {
// 	               panic("mock out the CreateNetworkConnection method")
//             },
//             CreateNetworkPeeringFunc: func(in1 context.Context, in2 *Network) (*NetworkPeering, error) {
// 	               panic("mock out the CreateNetworkPeering method")
//             },
//             DeleteNetworkFunc: func(in1 context.Context) error {
// 	               panic("mock out the DeleteNetwork method")
//             },
//             DeleteNetworkConnectionFunc: func(in1 context.Context, in2 *NetworkPeering) error {
// 	               panic("mock out the DeleteNetworkConnection method")
//             },
//             DeleteNetworkPeeringFunc: func(in1 *NetworkPeering) error {
// 	               panic("mock out the DeleteNetworkPeering method")
//             },
//             GetClusterNetworkPeeringFunc: func(in1 context.Context) (*NetworkPeering, error) {
// 	               panic("mock out the GetClusterNetworkPeering method")
//             },
//             IsEnabledFunc: func(in1 context.Context) (bool, error) {
// 	               panic("mock out the IsEnabled method")
//             },
//         }
//
//         // use mockedNetworkManager in code that requires NetworkManager
//         // and then make assertions.
//
//     }
type NetworkManagerMock struct {
	// CreateNetworkFunc mocks the CreateNetwork method.
	CreateNetworkFunc func(in1 context.Context, in2 *net.IPNet) (*Network, error)

	// CreateNetworkConnectionFunc mocks the CreateNetworkConnection method.
	CreateNetworkConnectionFunc func(in1 context.Context, in2 *Network) (*NetworkConnection, error)

	// CreateNetworkPeeringFunc mocks the CreateNetworkPeering method.
	CreateNetworkPeeringFunc func(in1 context.Context, in2 *Network) (*NetworkPeering, error)

	// DeleteNetworkFunc mocks the DeleteNetwork method.
	DeleteNetworkFunc func(in1 context.Context) error

	// DeleteNetworkConnectionFunc mocks the DeleteNetworkConnection method.
	DeleteNetworkConnectionFunc func(in1 context.Context, in2 *NetworkPeering) error

	// DeleteNetworkPeeringFunc mocks the DeleteNetworkPeering method.
	DeleteNetworkPeeringFunc func(in1 *NetworkPeering) error

	// GetClusterNetworkPeeringFunc mocks the GetClusterNetworkPeering method.
	GetClusterNetworkPeeringFunc func(in1 context.Context) (*NetworkPeering, error)

	// IsEnabledFunc mocks the IsEnabled method.
	IsEnabledFunc func(in1 context.Context) (bool, error)

	// calls tracks calls to the methods.
	calls struct {
		// CreateNetwork holds details about calls to the CreateNetwork method.
		CreateNetwork []struct {
			// In1 is the in1 argument value.
			In1 context.Context
			// In2 is the in2 argument value.
			In2 *net.IPNet
		}
		// CreateNetworkConnection holds details about calls to the CreateNetworkConnection method.
		CreateNetworkConnection []struct {
			// In1 is the in1 argument value.
			In1 context.Context
			// In2 is the in2 argument value.
			In2 *Network
		}
		// CreateNetworkPeering holds details about calls to the CreateNetworkPeering method.
		CreateNetworkPeering []struct {
			// In1 is the in1 argument value.
			In1 context.Context
			// In2 is the in2 argument value.
			In2 *Network
		}
		// DeleteNetwork holds details about calls to the DeleteNetwork method.
		DeleteNetwork []struct {
			// In1 is the in1 argument value.
			In1 context.Context
		}
		// DeleteNetworkConnection holds details about calls to the DeleteNetworkConnection method.
		DeleteNetworkConnection []struct {
			// In1 is the in1 argument value.
			In1 context.Context
			// In2 is the in2 argument value.
			In2 *NetworkPeering
		}
		// DeleteNetworkPeering holds details about calls to the DeleteNetworkPeering method.
		DeleteNetworkPeering []struct {
			// In1 is the in1 argument value.
			In1 *NetworkPeering
		}
		// GetClusterNetworkPeering holds details about calls to the GetClusterNetworkPeering method.
		GetClusterNetworkPeering []struct {
			// In1 is the in1 argument value.
			In1 context.Context
		}
		// IsEnabled holds details about calls to the IsEnabled method.
		IsEnabled []struct {
			// In1 is the in1 argument value.
			In1 context.Context
		}
	}
}

// CreateNetwork calls CreateNetworkFunc.
func (mock *NetworkManagerMock) CreateNetwork(in1 context.Context, in2 *net.IPNet) (*Network, error) {
	if mock.CreateNetworkFunc == nil {
		panic("NetworkManagerMock.CreateNetworkFunc: method is nil but NetworkManager.CreateNetwork was just called")
	}
	callInfo := struct {
		In1 context.Context
		In2 *net.IPNet
	}{
		In1: in1,
		In2: in2,
	}
	lockNetworkManagerMockCreateNetwork.Lock()
	mock.calls.CreateNetwork = append(mock.calls.CreateNetwork, callInfo)
	lockNetworkManagerMockCreateNetwork.Unlock()
	return mock.CreateNetworkFunc(in1, in2)
}

// CreateNetworkCalls gets all the calls that were made to CreateNetwork.
// Check the length with:
//     len(mockedNetworkManager.CreateNetworkCalls())
func (mock *NetworkManagerMock) CreateNetworkCalls() []struct {
	In1 context.Context
	In2 *net.IPNet
} {
	var calls []struct {
		In1 context.Context
		In2 *net.IPNet
	}
	lockNetworkManagerMockCreateNetwork.RLock()
	calls = mock.calls.CreateNetwork
	lockNetworkManagerMockCreateNetwork.RUnlock()
	return calls
}

// CreateNetworkConnection calls CreateNetworkConnectionFunc.
func (mock *NetworkManagerMock) CreateNetworkConnection(in1 context.Context, in2 *Network) (*NetworkConnection, error) {
	if mock.CreateNetworkConnectionFunc == nil {
		panic("NetworkManagerMock.CreateNetworkConnectionFunc: method is nil but NetworkManager.CreateNetworkConnection was just called")
	}
	callInfo := struct {
		In1 context.Context
		In2 *Network
	}{
		In1: in1,
		In2: in2,
	}
	lockNetworkManagerMockCreateNetworkConnection.Lock()
	mock.calls.CreateNetworkConnection = append(mock.calls.CreateNetworkConnection, callInfo)
	lockNetworkManagerMockCreateNetworkConnection.Unlock()
	return mock.CreateNetworkConnectionFunc(in1, in2)
}

// CreateNetworkConnectionCalls gets all the calls that were made to CreateNetworkConnection.
// Check the length with:
//     len(mockedNetworkManager.CreateNetworkConnectionCalls())
func (mock *NetworkManagerMock) CreateNetworkConnectionCalls() []struct {
	In1 context.Context
	In2 *Network
} {
	var calls []struct {
		In1 context.Context
		In2 *Network
	}
	lockNetworkManagerMockCreateNetworkConnection.RLock()
	calls = mock.calls.CreateNetworkConnection
	lockNetworkManagerMockCreateNetworkConnection.RUnlock()
	return calls
}

// CreateNetworkPeering calls CreateNetworkPeeringFunc.
func (mock *NetworkManagerMock) CreateNetworkPeering(in1 context.Context, in2 *Network) (*NetworkPeering, error) {
	if mock.CreateNetworkPeeringFunc == nil {
		panic("NetworkManagerMock.CreateNetworkPeeringFunc: method is nil but NetworkManager.CreateNetworkPeering was just called")
	}
	callInfo := struct {
		In1 context.Context
		In2 *Network
	}{
		In1: in1,
		In2: in2,
	}
	lockNetworkManagerMockCreateNetworkPeering.Lock()
	mock.calls.CreateNetworkPeering = append(mock.calls.CreateNetworkPeering, callInfo)
	lockNetworkManagerMockCreateNetworkPeering.Unlock()
	return mock.CreateNetworkPeeringFunc(in1, in2)
}

// CreateNetworkPeeringCalls gets all the calls that were made to CreateNetworkPeering.
// Check the length with:
//     len(mockedNetworkManager.CreateNetworkPeeringCalls())
func (mock *NetworkManagerMock) CreateNetworkPeeringCalls() []struct {
	In1 context.Context
	In2 *Network
} {
	var calls []struct {
		In1 context.Context
		In2 *Network
	}
	lockNetworkManagerMockCreateNetworkPeering.RLock()
	calls = mock.calls.CreateNetworkPeering
	lockNetworkManagerMockCreateNetworkPeering.RUnlock()
	return calls
}

// DeleteNetwork calls DeleteNetworkFunc.
func (mock *NetworkManagerMock) DeleteNetwork(in1 context.Context) error {
	if mock.DeleteNetworkFunc == nil {
		panic("NetworkManagerMock.DeleteNetworkFunc: method is nil but NetworkManager.DeleteNetwork was just called")
	}
	callInfo := struct {
		In1 context.Context
	}{
		In1: in1,
	}
	lockNetworkManagerMockDeleteNetwork.Lock()
	mock.calls.DeleteNetwork = append(mock.calls.DeleteNetwork, callInfo)
	lockNetworkManagerMockDeleteNetwork.Unlock()
	return mock.DeleteNetworkFunc(in1)
}

// DeleteNetworkCalls gets all the calls that were made to DeleteNetwork.
// Check the length with:
//     len(mockedNetworkManager.DeleteNetworkCalls())
func (mock *NetworkManagerMock) DeleteNetworkCalls() []struct {
	In1 context.Context
} {
	var calls []struct {
		In1 context.Context
	}
	lockNetworkManagerMockDeleteNetwork.RLock()
	calls = mock.calls.DeleteNetwork
	lockNetworkManagerMockDeleteNetwork.RUnlock()
	return calls
}

// DeleteNetworkConnection calls DeleteNetworkConnectionFunc.
func (mock *NetworkManagerMock) DeleteNetworkConnection(in1 context.Context, in2 *NetworkPeering) error {
	if mock.DeleteNetworkConnectionFunc == nil {
		panic("NetworkManagerMock.DeleteNetworkConnectionFunc: method is nil but NetworkManager.DeleteNetworkConnection was just called")
	}
	callInfo := struct {
		In1 context.Context
		In2 *NetworkPeering
	}{
		In1: in1,
		In2: in2,
	}
	lockNetworkManagerMockDeleteNetworkConnection.Lock()
	mock.calls.DeleteNetworkConnection = append(mock.calls.DeleteNetworkConnection, callInfo)
	lockNetworkManagerMockDeleteNetworkConnection.Unlock()
	return mock.DeleteNetworkConnectionFunc(in1, in2)
}

// DeleteNetworkConnectionCalls gets all the calls that were made to DeleteNetworkConnection.
// Check the length with:
//     len(mockedNetworkManager.DeleteNetworkConnectionCalls())
func (mock *NetworkManagerMock) DeleteNetworkConnectionCalls() []struct {
	In1 context.Context
	In2 *NetworkPeering
} {
	var calls []struct {
		In1 context.Context
		In2 *NetworkPeering
	}
	lockNetworkManagerMockDeleteNetworkConnection.RLock()
	calls = mock.calls.DeleteNetworkConnection
	lockNetworkManagerMockDeleteNetworkConnection.RUnlock()
	return calls
}

// DeleteNetworkPeering calls DeleteNetworkPeeringFunc.
func (mock *NetworkManagerMock) DeleteNetworkPeering(in1 *NetworkPeering) error {
	if mock.DeleteNetworkPeeringFunc == nil {
		panic("NetworkManagerMock.DeleteNetworkPeeringFunc: method is nil but NetworkManager.DeleteNetworkPeering was just called")
	}
	callInfo := struct {
		In1 *NetworkPeering
	}{
		In1: in1,
	}
	lockNetworkManagerMockDeleteNetworkPeering.Lock()
	mock.calls.DeleteNetworkPeering = append(mock.calls.DeleteNetworkPeering, callInfo)
	lockNetworkManagerMockDeleteNetworkPeering.Unlock()
	return mock.DeleteNetworkPeeringFunc(in1)
}

// DeleteNetworkPeeringCalls gets all the calls that were made to DeleteNetworkPeering.
// Check the length with:
//     len(mockedNetworkManager.DeleteNetworkPeeringCalls())
func (mock *NetworkManagerMock) DeleteNetworkPeeringCalls() []struct {
	In1 *NetworkPeering
} {
	var calls []struct {
		In1 *NetworkPeering
	}
	lockNetworkManagerMockDeleteNetworkPeering.RLock()
	calls = mock.calls.DeleteNetworkPeering
	lockNetworkManagerMockDeleteNetworkPeering.RUnlock()
	return calls
}

// GetClusterNetworkPeering calls GetClusterNetworkPeeringFunc.
func (mock *NetworkManagerMock) GetClusterNetworkPeering(in1 context.Context) (*NetworkPeering, error) {
	if mock.GetClusterNetworkPeeringFunc == nil {
		panic("NetworkManagerMock.GetClusterNetworkPeeringFunc: method is nil but NetworkManager.GetClusterNetworkPeering was just called")
	}
	callInfo := struct {
		In1 context.Context
	}{
		In1: in1,
	}
	lockNetworkManagerMockGetClusterNetworkPeering.Lock()
	mock.calls.GetClusterNetworkPeering = append(mock.calls.GetClusterNetworkPeering, callInfo)
	lockNetworkManagerMockGetClusterNetworkPeering.Unlock()
	return mock.GetClusterNetworkPeeringFunc(in1)
}

// GetClusterNetworkPeeringCalls gets all the calls that were made to GetClusterNetworkPeering.
// Check the length with:
//     len(mockedNetworkManager.GetClusterNetworkPeeringCalls())
func (mock *NetworkManagerMock) GetClusterNetworkPeeringCalls() []struct {
	In1 context.Context
} {
	var calls []struct {
		In1 context.Context
	}
	lockNetworkManagerMockGetClusterNetworkPeering.RLock()
	calls = mock.calls.GetClusterNetworkPeering
	lockNetworkManagerMockGetClusterNetworkPeering.RUnlock()
	return calls
}

// IsEnabled calls IsEnabledFunc.
func (mock *NetworkManagerMock) IsEnabled(in1 context.Context) (bool, error) {
	if mock.IsEnabledFunc == nil {
		panic("NetworkManagerMock.IsEnabledFunc: method is nil but NetworkManager.IsEnabled was just called")
	}
	callInfo := struct {
		In1 context.Context
	}{
		In1: in1,
	}
	lockNetworkManagerMockIsEnabled.Lock()
	mock.calls.IsEnabled = append(mock.calls.IsEnabled, callInfo)
	lockNetworkManagerMockIsEnabled.Unlock()
	return mock.IsEnabledFunc(in1)
}

// IsEnabledCalls gets all the calls that were made to IsEnabled.
// Check the length with:
//     len(mockedNetworkManager.IsEnabledCalls())
func (mock *NetworkManagerMock) IsEnabledCalls() []struct {
	In1 context.Context
} {
	var calls []struct {
		In1 context.Context
	}
	lockNetworkManagerMockIsEnabled.RLock()
	calls = mock.calls.IsEnabled
	lockNetworkManagerMockIsEnabled.RUnlock()
	return calls
}
