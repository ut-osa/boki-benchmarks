Benchmark workloads of [Boki](https://github.com/ut-osa/boki)
==================================

This repository includes source code of evaluation workloads of Boki,
and scripts for running experiments.
It includes all materials for the artifact evaluation of our SOSP '21 paper.

### Structure of this repository ###

* `dockerfiles`: Dockerfiles for building relevant Docker containers.
* `workloads`: source code of workloads for evaluating BokiFlow, BokiStore, and BokiQueue.
* `experiments`: setup scripts for running experiments of individual workloads.
* `scripts`: helper scripts for building Docker containers, and provisioning EC2 instances for experiments.

### Hardware and software dependencies ###

Our evaluation workloads run on AWS EC2 instances in `us-east-2` region.

EC2 VMs for running experiments use a public AMI (`ami-0c6de836734de3280`) built by us,
which is based on Ubuntu 20.04 with necessary dependencies installed.

### Environment setup ###

#### Setting up the controller machine ####

A controller machine in AWS `us-east-2` region is required for running scripts executing experiment workflows.
The controller machine can use very small EC2 instance type, as it only provisions and controls experiment VMs,
but does not affect experimental results.
In our own setup, we use a `t3.micro` EC2 instance installed with Ubuntu 20.04 as the controller machine.

The controller machine needs `python3`, `rsync`, and AWS CLI version 1 installed.
`python3` and `rsync` can be installed with `apt`,
and this [documentation](https://docs.aws.amazon.com/cli/latest/userguide/install-linux.html)
details the recommanded way for installing AWS CLI version 1.
Once installed, AWS CLI has to be configured with region `us-east-2` and access key
(see this [documentation](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-quickstart.html)).

Then on the controller machine, clone this repository with all git submodules
```
git clone --recursive https://github.com/ut-osa/boki-benchmarks.git
```
Finally, execute `scripts/setup_sshkey.sh` to setup SSH keys that will be used to access experiment VMs.
Please read the notice in `scripts/setup_sshkey.sh` before executing it to see if this script works for your setup.

#### Setting up EC2 security group and placement group ####

Our VM provisioning script creates EC2 instances with security group `boki` and placement group `boki-experiments`.
The security group includes firewall rules for experiment VMs (including allowing the controller machine to SSH into them),
while the placement group instructs AWS to place experiment VMs close together.

Executing `scripts/aws_provision.sh` on the controller machine creates these groups with correct configurations.

#### Building Docker images ####
We also provide the script (`scripts/docker_images.sh`) for building Docker images relevant to experiments in this artifact.
As we already pushed all compiled images to DockerHub, there is no need to run this script
as long as you do not modify source code of Boki (in `boki` directory) and evaluation workloads (in `workloads` directory).

### Experiment workflow ###

Each sub-directory within `experiments` corresponds to one experiment.
Within each experiment directory, a `config.json` file describes machine configuration and placement assignment of
individual Docker containers for this experiment.

`run_once.sh` script is the entry point of one experiment run, which performs workload-specific setups,
runs the benchmark with configured options, and stores results in  `results` directory.

Before executing `run_once.sh` script, VM provisioning is done by `scripts/exp_helper`
with sub-command `start-machines`.
After EC2 instances are up, the script then sets up Docker engines on newly created
VMs to form a Docker cluster in [swarm](https://docs.docker.com/engine/swarm/) mode.

`experiments/run_quick.sh` is a push-button script for running all experiments with a small set
of options. This script can be used to quickly test all setups, and learn the experiment workflow.

### Evaluation and expected result ###

Within individual result directory, a `results.log` or `latency.txt` file describes the metrics
of this run. `scripts/summarize_results.py` is a simple script to print a result summary for full inspection.

Within the directory of each workload experiment, we provide `exepcted_results` directory, including some examples of experiment results.

### License ###

* [Boki](https://github.com/ut-osa/boki) is licensed under Apache License 2.0.
* BokiFlow (`workloads/workflow`) derives from [Beldi codebase](https://github.com/eniac/Beldi). BokiFlow is licensed under MIT License, in accordance with Beldi.
* All other source code in this repository is licensed under Apache License 2.0.
