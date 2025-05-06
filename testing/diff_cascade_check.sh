          
terraform apply -auto-approve 

terraform_output=$(terraform plan -detailed-exitcode -out=tfplan_post 2>&1)

diff_cascade=$(echo "$terraform_output" | grep 'Plan:' && echo true || echo false)

terraform destroy -auto-approve

if [[ $diff_cascade == "true" ]]; then
echo "::error::Unexpected changes detected after apply"
exit 1
else
echo "No changes detected after apply - infrastructure is consistent"
fi

# Generate diff