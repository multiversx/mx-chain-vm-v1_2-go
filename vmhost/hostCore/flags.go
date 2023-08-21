package hostCore

import "github.com/multiversx/mx-chain-core-go/core"

const (
	SCDeployFlag            core.EnableEpochFlag = "SCDeployFlag"
	BuiltInFunctionsFlag    core.EnableEpochFlag = "BuiltInFunctionsFlag"
	RepairCallbackFlag      core.EnableEpochFlag = "RepairCallbackFlag"
	AheadOfTimeGasUsageFlag core.EnableEpochFlag = "AheadOfTimeGasUsageFlag"
)

// allFlags must have all flags used by mx-chain-vm-v1_2-go in the current version
var allFlags = []core.EnableEpochFlag{
	SCDeployFlag,
	BuiltInFunctionsFlag,
	RepairCallbackFlag,
	AheadOfTimeGasUsageFlag,
}
