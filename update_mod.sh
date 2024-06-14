dir="/Users/joseph/Github/terraform-provider-jamfpro"

if [ -d "$dir" ]; then
  echo "Directory exists. Proceeding with updates..."

  echo "Updating integration..."
  GOPROXY=direct go get -u github.com/deploymenttheory/go-api-http-client-integration-jamfpro/jamfprointegration

  echo "Updating SDK..."
  GOPROXY=direct go get -u github.com/deploymenttheory/go-api-sdk-jamfpro@dev-jl-httpclientv2

  # echo "Updating HTTP client..."
  # GOPROXY=direct go get -u github.com/deploymenttheory/go-api-http-client@dev-jl-version2
else
  echo "Directory does not exist"
fi
