#!/bin/bash
BASE_DIR=`realpath $(dirname $0)`
ROOT_DIR=`realpath $BASE_DIR/../../..`

HELPER_SCRIPT=$ROOT_DIR/scripts/exp_helper

$HELPER_SCRIPT start-machines --base-dir=$BASE_DIR

$BASE_DIR/run_once.sh p64c64   8  64  64
$BASE_DIR/run_once.sh p128c128 10 128 128
$BASE_DIR/run_once.sh p256c256 8  256 256

$BASE_DIR/run_once.sh p64c16  22 64  16
$BASE_DIR/run_once.sh p128c32 24 128 32
$BASE_DIR/run_once.sh p256c64 24 256 64

$BASE_DIR/run_once.sh p16c64  6 16 64
$BASE_DIR/run_once.sh p32c128 7 32 128
$BASE_DIR/run_once.sh p64c256 7 64 256

$HELPER_SCRIPT stop-machines --base-dir=$BASE_DIR
