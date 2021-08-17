package main

import (
	// "fmt"
	// "github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/eniac/Beldi/pkg/beldilib"
	"github.com/lithammer/shortuuid"
	"time"

	"cs.utexas.edu/zjia/faas"
)

var TXN = "DISABLE"

func Handler(env *beldilib.Env) interface{} {
	results := map[string]int64{}

	if TXN == "ENABLE" {
		panic("Not implemented")
	}
	if beldilib.TYPE == "BELDI" {
		a := shortuuid.New()
		start := time.Now()
		beldilib.Write(env, "singleop", "K", map[expression.NameBuilder]expression.OperandBuilder{
			expression.Name("V"): expression.Value(a),
		})
		results["latencyDWrite"] = time.Since(start).Microseconds()
		// fmt.Printf("DURATION DWrite %s\n", time.Since(start))

		start = time.Now()
		beldilib.CondWrite(env, "singleop", "K", map[expression.NameBuilder]expression.OperandBuilder{
			expression.Name("V2"): expression.Value(1),
		}, expression.Name("V").Equal(expression.Value(a)))
		results["latencyCWriteT"] = time.Since(start).Microseconds()
		// fmt.Printf("DURATION CWriteT %s\n", time.Since(start))

		start = time.Now()
		beldilib.CondWrite(env, "singleop", "K", map[expression.NameBuilder]expression.OperandBuilder{
			expression.Name("V2"): expression.Value(a),
		}, expression.Name("V").Equal(expression.Value(2)))
		results["latencyCWriteF"] = time.Since(start).Microseconds()
		// fmt.Printf("DURATION CWriteF %s\n", time.Since(start))

		start = time.Now()
		beldilib.Read(env, "singleop", "K")
		results["latencyRead"] = time.Since(start).Microseconds()
		// fmt.Printf("DURATION Read %s\n", time.Since(start))

		start = time.Now()
		beldilib.SyncInvoke(env, "nop", "")
		results["latencyCall"] = time.Since(start).Microseconds()
		// fmt.Printf("DURATION Call %s\n", time.Since(start))
	}
	if beldilib.TYPE == "BASELINE" {
		a := shortuuid.New()
		start := time.Now()
		beldilib.Write(env, "bsingleop", "K", map[expression.NameBuilder]expression.OperandBuilder{
			expression.Name("V"): expression.Value(a),
		})
		results["latencyDWrite"] = time.Since(start).Microseconds()
		// fmt.Printf("DURATION DWrite %s\n", time.Since(start))

		start = time.Now()
		beldilib.CondWrite(env, "bsingleop", "K", map[expression.NameBuilder]expression.OperandBuilder{
			expression.Name("V2"): expression.Value(1),
		}, expression.Name("V").Equal(expression.Value(a)))
		results["latencyCWriteT"] = time.Since(start).Microseconds()
		// fmt.Printf("DURATION CWriteT %s\n", time.Since(start))

		start = time.Now()
		beldilib.CondWrite(env, "bsingleop", "K", map[expression.NameBuilder]expression.OperandBuilder{
			expression.Name("V2"): expression.Value(a),
		}, expression.Name("V").Equal(expression.Value(2)))
		results["latencyCWriteF"] = time.Since(start).Microseconds()
		// fmt.Printf("DURATION CWriteF %s\n", time.Since(start))

		start = time.Now()
		beldilib.Read(env, "bsingleop", "K")
		results["latencyRead"] = time.Since(start).Microseconds()
		// fmt.Printf("DURATION Read %s\n", time.Since(start))

		start = time.Now()
		beldilib.SyncInvoke(env, "bnop", "")
		results["latencyCall"] = time.Since(start).Microseconds()
		// fmt.Printf("DURATION Call %s\n", time.Since(start))
	}
	return results
}

func main() {
	// lambda.Start(beldilib.Wrapper(Handler))
	faas.Serve(beldilib.CreateFuncHandlerFactory(Handler))
}
