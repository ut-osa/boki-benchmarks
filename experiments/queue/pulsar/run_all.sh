#!/bin/bash
BASE_DIR=`realpath $(dirname $0)`
ROOT_DIR=`realpath $BASE_DIR/../../..`

HELPER_SCRIPT=$ROOT_DIR/scripts/exp_helper

$HELPER_SCRIPT start-machines --base-dir=$BASE_DIR

$BASE_DIR/run_once.sh p64c64   6 64  64
$BASE_DIR/run_once.sh p128c128 6 128 128
$BASE_DIR/run_once.sh p256c256 8 256 256

$BASE_DIR/run_once.sh p64c16  7  64  16
$BASE_DIR/run_once.sh p128c32 8  128 32
$BASE_DIR/run_once.sh p256c64 12 256 64

$BASE_DIR/run_once.sh p16c64  3 16 64
$BASE_DIR/run_once.sh p32c128 3 32 128
$BASE_DIR/run_once.sh p64c256 4 64 256

$HELPER_SCRIPT stop-machines --base-dir=$BASE_DIR
