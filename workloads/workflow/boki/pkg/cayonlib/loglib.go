package cayonlib

import (
	"fmt"
	// "log"
	"encoding/json"
	// "cs.utexas.edu/zjia/faas/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/golang/snappy"
	// "context"
)

type IntentLogEntry struct {
	SeqNum     uint64        `json:"-"`
	InstanceId string        `json:"instanceId"`
	StepNumber int32         `json:"step"`
	PostStep   bool          `json:"postStep"`
	Data       aws.JSONValue `json:"data"`
}

type IntentFsm struct {
	instanceId   string
	stepNumber   int32
	tail         *IntentLogEntry
	stepLogs     map[int32]*IntentLogEntry
	postStepLogs map[int32]*IntentLogEntry
}

func NewIntentFsm(instanceId string) *IntentFsm {
	return &IntentFsm{
		instanceId:   instanceId,
		stepNumber:   0,
		tail:         nil,
		stepLogs:     make(map[int32]*IntentLogEntry),
		postStepLogs: make(map[int32]*IntentLogEntry),
	}
}

func (fsm *IntentFsm) applyLog(intentLog *IntentLogEntry) {
	fsm.tail = intentLog
	step := intentLog.StepNumber
	if intentLog.PostStep {
		if _, exists := fsm.postStepLogs[step]; !exists {
			fsm.postStepLogs[step] = intentLog
		}
	} else {
		if _, exists := fsm.stepLogs[step]; !exists {
			if step != fsm.stepNumber {
				panic(fmt.Sprintf("StepNumber is not monotonic: expected=%d, seen=%d", fsm.stepNumber, step))
			}
			fsm.stepNumber += 1
			fsm.stepLogs[step] = intentLog
		}
	}
}

func (fsm *IntentFsm) Catch(env *Env) {
	tag := IntentStepStreamTag(fsm.instanceId)
	seqNum := uint64(0)
	if fsm.tail != nil {
		seqNum = fsm.tail.SeqNum + 1
	}
	for {
		logEntry, err := env.FaasEnv.SharedLogReadNext(env.FaasCtx, tag, seqNum)
		CHECK(err)
		if logEntry == nil {
			break
		}
		decoded, err := snappy.Decode(nil, logEntry.Data)
		CHECK(err)
		var intentLog IntentLogEntry
		err = json.Unmarshal(decoded, &intentLog)
		CHECK(err)
		if intentLog.InstanceId == fsm.instanceId {
			// log.Printf("[INFO] Found my log: seqnum=%d, step=%d", logEntry.SeqNum, intentLog.StepNumber)
			intentLog.SeqNum = logEntry.SeqNum
			fsm.applyLog(&intentLog)
		}
		seqNum = logEntry.SeqNum + 1
	}
}

func (fsm *IntentFsm) GetStepLog(stepNumber int32) *IntentLogEntry {
	if log, exists := fsm.stepLogs[stepNumber]; exists {
		return log
	} else {
		return nil
	}
}

func (fsm *IntentFsm) GetPostStepLog(stepNumber int32) *IntentLogEntry {
	if log, exists := fsm.postStepLogs[stepNumber]; exists {
		return log
	} else {
		return nil
	}
}

func ProposeNextStep(env *Env, data aws.JSONValue) (bool, *IntentLogEntry) {
	step := env.StepNumber
	env.StepNumber += 1
	intentLog := env.Fsm.GetStepLog(step)
	if intentLog != nil {
		return false, intentLog
	}
	intentLog = &IntentLogEntry{
		InstanceId: env.InstanceId,
		StepNumber: step,
		PostStep:   false,
		Data:       data,
	}
	seqNum := LibAppendLog(env, IntentStepStreamTag(env.InstanceId), &intentLog)
	env.Fsm.Catch(env)
	intentLog = env.Fsm.GetStepLog(step)
	if intentLog == nil {
		panic(fmt.Sprintf("Cannot find intent log for step %d", step))
	}
	return seqNum == intentLog.SeqNum, intentLog
}

func LogStepResult(env *Env, instanceId string, stepNumber int32, data aws.JSONValue) {
	LibAppendLog(env, IntentStepStreamTag(instanceId), &IntentLogEntry{
		InstanceId: instanceId,
		StepNumber: stepNumber,
		PostStep:   true,
		Data:       data,
	})
}

func FetchStepResultLog(env *Env, stepNumber int32, catch bool) *IntentLogEntry {
	intentLog := env.Fsm.GetPostStepLog(stepNumber)
	if intentLog != nil {
		return intentLog
	}
	if catch {
		env.Fsm.Catch(env)
	} else {
		return nil
	}
	return env.Fsm.GetPostStepLog(stepNumber)
}

func LibAppendLog(env *Env, tag uint64, data interface{}) uint64 {
	serializedData, err := json.Marshal(data)
	CHECK(err)
	encoded := snappy.Encode(nil, serializedData)
	seqNum, err := env.FaasEnv.SharedLogAppend(env.FaasCtx, []uint64{tag}, encoded)
	CHECK(err)
	return seqNum
}

func CheckLogDataField(intentLog *IntentLogEntry, field string, expected string) {
	if tmp := intentLog.Data[field].(string); tmp != expected {
		panic(fmt.Sprintf("Field %s mismatch: expected=%s, have=%s", field, expected, tmp))
	}
}
