#!/bin/bash
BASE_DIR=`realpath $(dirname $0)`
ROOT_DIR=`realpath $BASE_DIR/../../..`

HELPER_SCRIPT=$ROOT_DIR/scripts/exp_helper

$HELPER_SCRIPT start-machines --base-dir=$BASE_DIR

$BASE_DIR/run_once.sh p64c64   64  6 1 64
$BASE_DIR/run_once.sh p128c128 128 6 1 128
$BASE_DIR/run_once.sh p256c256 256 8 1 256

$BASE_DIR/run_once.sh p64c16  16 7  1 64
$BASE_DIR/run_once.sh p128c32 32 8  1 128
$BASE_DIR/run_once.sh p256c64 64 12 1 256

$BASE_DIR/run_once.sh p16c64  64  3 1 16
$BASE_DIR/run_once.sh p32c128 128 3 1 32
$BASE_DIR/run_once.sh p64c256 256 4 1 64

$HELPER_SCRIPT stop-machines --base-dir=$BASE_DIR
