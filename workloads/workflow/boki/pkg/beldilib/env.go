package beldilib

import (
	"cs.utexas.edu/zjia/faas/types"
	"context"
)

type Env struct {
	LambdaId    string
	InstanceId  string
	LogTable    string
	IntentTable string
	LocalTable  string
	StepNumber  int32
	Input       interface{}
	TxnId       string
	Instruction string
	Baseline    bool
	FaasCtx     context.Context
	FaasEnv     types.Environment
}
