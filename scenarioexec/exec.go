package scenarioexec

import (
	"fmt"
	"path/filepath"

	logger "github.com/multiversx/mx-chain-logger-go"
	vmi "github.com/multiversx/mx-chain-vm-common-go"
	"github.com/multiversx/mx-chain-vm-common-go/mock"
	"github.com/multiversx/mx-chain-vm-v1_2-go/config"
	worldhook "github.com/multiversx/mx-chain-vm-v1_2-go/mock/world"
	mc "github.com/multiversx/mx-chain-vm-v1_2-go/scenarios/controller"
	er "github.com/multiversx/mx-chain-vm-v1_2-go/scenarios/expression/reconstructor"
	fr "github.com/multiversx/mx-chain-vm-v1_2-go/scenarios/fileresolver"
	mj "github.com/multiversx/mx-chain-vm-v1_2-go/scenarios/json/model"
	"github.com/multiversx/mx-chain-vm-v1_2-go/vmhost"
	"github.com/multiversx/mx-chain-vm-v1_2-go/vmhost/hostCore"
)

var log = logger.GetOrCreate("vm/scenarios")

// TestVMType is the VM type argument we use in tests.
var TestVMType = []byte{0, 0}

// VMTestExecutor parses, interprets and executes both .test.json tests and .scen.json scenarios with VM.
type VMTestExecutor struct {
	World                 *worldhook.MockWorld
	vm                    vmi.VMExecutionHandler
	checkGas              bool
	scenarioexecPath      string
	scenGasScheduleLoaded bool
	fileResolver          fr.FileResolver
	exprReconstructor     er.ExprReconstructor
}

var _ mc.TestExecutor = (*VMTestExecutor)(nil)
var _ mc.ScenarioExecutor = (*VMTestExecutor)(nil)

// NewVMTestExecutor prepares a new VMTestExecutor instance.
func NewVMTestExecutor(scenarioexecPath string) (*VMTestExecutor, error) {
	world := worldhook.NewMockWorld()

	gasScheduleMap := config.MakeGasMapForTests()
	err := world.InitBuiltinFunctions(gasScheduleMap)
	if err != nil {
		return nil, err
	}

	blockGasLimit := uint64(10000000)
	vm, err := hostCore.NewVMHost(world, &vmhost.VMHostParameters{
		VMType:                   TestVMType,
		BlockGasLimit:            blockGasLimit,
		GasSchedule:              gasScheduleMap,
		ProtocolBuiltinFunctions: world.GetBuiltinFunctionNames(),
		ProtectedKeyPrefix:       []byte(ProtectedKeyPrefix),
		EnableEpochsHandler: &mock.EnableEpochsHandlerStub{
			IsSCDeployFlagEnabledField:            true,
			IsAheadOfTimeGasUsageFlagEnabledField: true,
			IsRepairCallbackFlagEnabledField:      true,
			IsBuiltInFunctionsFlagEnabledField:    true,
		},
	})
	if err != nil {
		return nil, err
	}

	return &VMTestExecutor{
		World:                 world,
		vm:                    vm,
		checkGas:              true,
		scenarioexecPath:      scenarioexecPath,
		scenGasScheduleLoaded: false,
		fileResolver:          nil,
		exprReconstructor:     er.ExprReconstructor{},
	}, nil
}

// GetVM yields a reference to the VMExecutionHandler used.
func (ae *VMTestExecutor) GetVM() vmi.VMExecutionHandler {
	return ae.vm
}

func (ae *VMTestExecutor) gasScheduleMapFromScenarios(scenGasSchedule mj.GasSchedule) (config.GasScheduleMap, error) {
	switch scenGasSchedule {
	case mj.GasScheduleDefault:
		return hostCore.LoadGasScheduleConfig(filepath.Join(ae.scenarioexecPath, "gasSchedules/gasScheduleV3.toml"))
	case mj.GasScheduleDummy:
		return config.MakeGasMapForTests(), nil
	case mj.GasScheduleV1:
		return hostCore.LoadGasScheduleConfig(filepath.Join(ae.scenarioexecPath, "gasSchedules/gasScheduleV1.toml"))
	case mj.GasScheduleV2:
		return hostCore.LoadGasScheduleConfig(filepath.Join(ae.scenarioexecPath, "gasSchedules/gasScheduleV2.toml"))
	case mj.GasScheduleV3:
		return hostCore.LoadGasScheduleConfig(filepath.Join(ae.scenarioexecPath, "gasSchedules/gasScheduleV3.toml"))
	default:
		return nil, fmt.Errorf("unknown scenario GasSchedule: %d", scenGasSchedule)
	}
}

// SetScenariosGasSchedule updates the gas costs based on the scenario config
// only changes the gas schedule once,
// this prevents subsequent gasSchedule declarations in externalSteps to overwrite
func (ae *VMTestExecutor) SetScenariosGasSchedule(newGasSchedule mj.GasSchedule) error {
	if ae.scenGasScheduleLoaded {
		return nil
	}
	ae.scenGasScheduleLoaded = true
	gasSchedule, err := ae.gasScheduleMapFromScenarios(newGasSchedule)
	if err != nil {
		return err
	}
	ae.vm.GasScheduleChange(gasSchedule)
	return nil
}
