# Generate and target diego's manifest:

~/workspace/cf-release/generate_deployment_manifest aws /home/ubuntu/workspace/deployments/microbosh/director.yml ~/workspace/diego-release/templates/enable_diego_in_cc.yml /home/ubuntu/workspace/deployments/microbosh/terraform_settings.yml  ~/workspace/deployments/microbosh/director.yml > ~/workspace/deployments/microbosh/cf.yml

# Dance some more:
bosh create release --force
bosh -n upload release
bosh -n deploy

