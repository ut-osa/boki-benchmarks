package cayonlib

import (
	"cs.utexas.edu/zjia/faas/types"
	"context"
)

type LogEntry struct {
	SeqNum uint64
	Data   map[string]interface{}
}

type Env struct {
	LambdaId    string
	InstanceId  string
	StepNumber  int32
	Input       interface{}
	TxnId       string
	Instruction string
	Baseline    bool
	FaasCtx     context.Context
	FaasEnv     types.Environment
	Fsm         *IntentFsm
}
