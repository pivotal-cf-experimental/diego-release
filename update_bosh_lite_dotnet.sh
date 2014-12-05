# # Generate a deployment stub with the BOSH director UUID
# 
# mkdir -p ~/deployments/bosh-lite
# cd ~/workspace/diego-release
# ./scripts/print-director-stub > ~/deployments/bosh-lite/director.yml
# 
# # Generate and target cf-release manifest:
# 
# cd ~/workspace/cf-release
# ./generate_deployment_manifest warden \
#     ~/deployments/bosh-lite/director.yml \
#     ~/workspace/diego-release/templates/enable_diego_in_cc.yml > \
#     ~/deployments/bosh-lite/cf.yml
# bosh deployment ~/deployments/bosh-lite/cf.yml
# 
# # Do the BOSH dance:
# 
# cd ~/workspace/cf-release
# bosh create release --force
# bosh -n upload release
# bosh -n deploy

# Generate and target diego's manifest:

cd ~/workspace/diego-release
./scripts/generate-deployment-manifest bosh-lite ../cf-release \
    ~/deployments/bosh-lite/director.yml > \
    ~/deployments/bosh-lite/diego.yml
bosh deployment ~/deployments/bosh-lite/diego.yml

# Dance some more:

bosh create release --force
bosh -n upload release
bosh -n deploy

