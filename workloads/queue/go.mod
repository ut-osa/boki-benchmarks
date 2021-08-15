module cs.utexas.edu/zjia/faas-queue

go 1.14

require (
	cs.utexas.edu/zjia/faas v0.0.0
	cs.utexas.edu/zjia/faas/slib v0.0.0
	github.com/apache/pulsar-client-go v0.4.0
	github.com/aws/aws-sdk-go v1.37.20
	github.com/golang/snappy v0.0.2 // indirect
	github.com/google/uuid v1.2.0
	github.com/montanaflynn/stats v0.6.3
)

replace cs.utexas.edu/zjia/faas => /src/boki/worker/golang

replace cs.utexas.edu/zjia/faas/slib => /src/boki/slib
