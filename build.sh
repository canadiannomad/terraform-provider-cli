#!/bin/bash
set -e
go build
cp terraform-provider-cli ~/.terraform.d/plugins/terraform-provider-cli_v0.0.1
cd test
#rm -rf terraform.tfstate* .terraform
#export TF_LOG=TRACE
export TF_LOG=DEBUG
echo "STEP: init" > log.txt &&
terraform init >> log.txt 2>&1 &&
echo "STEP: apply 1" >> log.txt 2>&1 &&
terraform apply --auto-approve >> log.txt 2>&1 &&
echo "STEP: state show 1" >> log.txt 2>&1 &&
terraform state show >> log.txt 2>&1 &&
echo "STEP: apply 2" >> log.txt 2>&1 &&
terraform apply --auto-approve >> log.txt 2>&1 &&
echo "STEP: state show 2" >> log.txt 2>&1 &&
terraform state show >> log.txt 2>&1 &&
#echo "STEP: destroy" >> log.txt 2>&1 &&
#terraform destroy -force >> log.txt 2>&1 &&
true
less -R log.txt
