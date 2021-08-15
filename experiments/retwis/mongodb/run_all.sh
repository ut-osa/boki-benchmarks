#!/bin/bash
BASE_DIR=`realpath $(dirname $0)`
ROOT_DIR=`realpath $BASE_DIR/../../..`

HELPER_SCRIPT=$ROOT_DIR/scripts/exp_helper

$HELPER_SCRIPT start-machines --base-dir=$BASE_DIR

$BASE_DIR/run_once.sh con64  64
$BASE_DIR/run_once.sh con96  96
$BASE_DIR/run_once.sh con128 128
$BASE_DIR/run_once.sh con192 192

$HELPER_SCRIPT stop-machines --base-dir=$BASE_DIR
