package worldmock

import "bytes"

// GenerateMockAddress simulates creation of a new address by the protocol.
func GenerateMockAddress(creatorAddress []byte, creatorNonce uint64) []byte {
	result := make([]byte, 32)
	result[10] = 0x11
	result[11] = 0x11
	result[12] = 0x11
	result[13] = 0x11
	copy(result[14:29], creatorAddress)

	result[29] = byte(creatorNonce)

	copy(result[30:], creatorAddress[30:])
	return result
}

// CreateMockWorldNewAddress creates a new address, simulating the protocol's address generation.
func CreateMockWorldNewAddress(world *MockWorld) func(address []byte, nonce uint64, vmType []byte) ([]byte, error) {
	return func(address []byte, nonce uint64, vmType []byte) ([]byte, error) {
		// custom error
		if world.Err != nil {
			return nil, world.Err
		}

		// explicit new address mocks
		// matched by creator address and nonce
		for _, newAddressMock := range world.NewAddressMocks {
			if bytes.Equal(address, newAddressMock.CreatorAddress) && nonce == newAddressMock.CreatorNonce {
				world.LastCreatedContractAddress = newAddressMock.NewAddress
				return newAddressMock.NewAddress, nil
			}
		}

		// If a mock address wasn't registered for the specified creatorAddress, generate one automatically.
		// This is not the real algorithm but it's simple and close enough.
		result := GenerateMockAddress(address, nonce)
		world.LastCreatedContractAddress = result
		return result, nil
	}
}
