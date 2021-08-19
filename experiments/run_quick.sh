#!/bin/bash
ROOT_DIR=`realpath $(dirname $0)/..`

HELPER_SCRIPT=$ROOT_DIR/scripts/exp_helper

# This IAM role has DynamoDB read/write access
BOKI_MACHINE_IAM=boki-ae-experiments


echo "====== Start running BokiQueue experiments ======"

BASE_DIR=$ROOT_DIR/experiments/queue/boki

$HELPER_SCRIPT start-machines --base-dir=$BASE_DIR --instance-iam-role $BOKI_MACHINE_IAM

$BASE_DIR/run_once.sh p128c128 128 6 1 128
$BASE_DIR/run_once.sh p128c32  32  8 1 128
$BASE_DIR/run_once.sh p32c128  128 3 1 32

$HELPER_SCRIPT stop-machines --base-dir=$BASE_DIR

echo "====== Finish running BokiQueue experiments ======"
echo ""


echo "====== Start running Pulsar experiments ======"

BASE_DIR=$ROOT_DIR/experiments/queue/pulsar

$HELPER_SCRIPT start-machines --base-dir=$BASE_DIR --instance-iam-role $BOKI_MACHINE_IAM

$BASE_DIR/run_once.sh p128c128 6 128 128
$BASE_DIR/run_once.sh p128c32  8 128 32
$BASE_DIR/run_once.sh p32c128  3 32  128

$HELPER_SCRIPT stop-machines --base-dir=$BASE_DIR

echo "====== Finish running Pulsar experiments ======"
echo ""


echo "====== Start running BokiFlow experiments ======"

BASE_DIR=$ROOT_DIR/experiments/workflow/boki-hotel

$HELPER_SCRIPT start-machines --base-dir=$BASE_DIR --instance-iam-role $BOKI_MACHINE_IAM

$BASE_DIR/run_once.sh qps100 100
$BASE_DIR/run_once.sh qps200 200
$BASE_DIR/run_once.sh qps300 300

$HELPER_SCRIPT stop-machines --base-dir=$BASE_DIR

BASE_DIR=$ROOT_DIR/experiments/workflow/boki-movie

$HELPER_SCRIPT start-machines --base-dir=$BASE_DIR --instance-iam-role $BOKI_MACHINE_IAM

$BASE_DIR/run_once.sh qps50  50
$BASE_DIR/run_once.sh qps100 100
$BASE_DIR/run_once.sh qps150 150

$HELPER_SCRIPT stop-machines --base-dir=$BASE_DIR

echo "====== Finish running BokiFlow experiments ======"
echo ""


echo "====== Start running Beldi experiments ======"

BASE_DIR=$ROOT_DIR/experiments/workflow/beldi-hotel

$HELPER_SCRIPT start-machines --base-dir=$BASE_DIR --instance-iam-role $BOKI_MACHINE_IAM

$BASE_DIR/run_once.sh qps100 100
$BASE_DIR/run_once.sh qps200 200
$BASE_DIR/run_once.sh qps300 300

$HELPER_SCRIPT stop-machines --base-dir=$BASE_DIR

BASE_DIR=$ROOT_DIR/experiments/workflow/beldi-movie

$HELPER_SCRIPT start-machines --base-dir=$BASE_DIR --instance-iam-role $BOKI_MACHINE_IAM

$BASE_DIR/run_once.sh qps50  50
$BASE_DIR/run_once.sh qps100 100
$BASE_DIR/run_once.sh qps150 150

$HELPER_SCRIPT stop-machines --base-dir=$BASE_DIR

echo "====== Finish running Beldi experiments ======"
echo ""
