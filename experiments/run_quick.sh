#!/bin/bash
ROOT_DIR=`realpath $(dirname $0)/..`

HELPER_SCRIPT=$ROOT_DIR/scripts/exp_helper

echo "====== Start running BokiQueue experiments ======"

BASE_DIR=$ROOT_DIR/experiments/queue/boki

$HELPER_SCRIPT start-machines --base-dir=$BASE_DIR --use-spot-instances

$BASE_DIR/run_once.sh p128c128 128 6 1 128
$BASE_DIR/run_once.sh p128c32  32  8 1 128
$BASE_DIR/run_once.sh p32c128  128 3 1 32

$HELPER_SCRIPT stop-machines --base-dir=$BASE_DIR

echo "====== Finish running BokiQueue experiments ======"
