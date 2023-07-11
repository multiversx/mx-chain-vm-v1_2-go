package mock

// EnableEpochsHandlerStub -
type EnableEpochsHandlerStub struct {
	GetCurrentEpochField                  uint32
	IsBuiltInFunctionsFlagEnabledField    bool
	IsSCDeployFlagEnabledField            bool
	IsRepairCallbackFlagEnabledField      bool
	IsAheadOfTimeGasUsageFlagEnabledField bool
}

// GetCurrentEpoch -
func (stub *EnableEpochsHandlerStub) GetCurrentEpoch() uint32 {
	return stub.GetCurrentEpochField
}

// IsBuiltInFunctionsFlagEnabledInEpoch -
func (stub *EnableEpochsHandlerStub) IsBuiltInFunctionsFlagEnabledInEpoch(_ uint32) bool {
	return stub.IsBuiltInFunctionsFlagEnabledField
}

// IsSCDeployFlagEnabledInEpoch -
func (stub *EnableEpochsHandlerStub) IsSCDeployFlagEnabledInEpoch(_ uint32) bool {
	return stub.IsSCDeployFlagEnabledField
}

// IsRepairCallbackFlagEnabledInEpoch -
func (stub *EnableEpochsHandlerStub) IsRepairCallbackFlagEnabledInEpoch(_ uint32) bool {
	return stub.IsRepairCallbackFlagEnabledField
}

// IsAheadOfTimeGasUsageFlagEnabledInEpoch -
func (stub *EnableEpochsHandlerStub) IsAheadOfTimeGasUsageFlagEnabledInEpoch(_ uint32) bool {
	return stub.IsAheadOfTimeGasUsageFlagEnabledField
}

// IsInterfaceNil -
func (stub *EnableEpochsHandlerStub) IsInterfaceNil() bool {
	return stub == nil
}
