#!/usr/bin/python3
import os
import glob
import sys
import re

if __name__ == '__main__':
    exp_dir = os.path.realpath(os.path.join(__file__, '../../experiments'))

    result_files = glob.glob('%s/**/results/*/results.log' % (exp_dir,), recursive=True)
    result_files += glob.glob('%s/**/results/*/latency.txt' % (exp_dir,), recursive=True)
    result_files.sort()

    for file_path in result_files:
        parts = file_path.split('/')
        workload = parts[-5]
        system = parts[-4]
        subexp = parts[-2]
        with open(file_path) as fin:
            contents = fin.read()
        sys.stdout.write('===== %s (%s), sub-experiment %s =====\n' % (workload, system, subexp))
        sys.stdout.write(contents + '\n')
