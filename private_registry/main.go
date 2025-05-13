// // HashiCorp's public PGP key - https://www.hashicorp.com/.well-known/pgp-key.txt?ajs_aid=a4e1422e-4cd3-4b83-86c8-96aaab547d3a&product_intent=terraform&utm_source=WEBSITE&utm_medium=WEB_IO&utm_offer=ARTICLE_PAGE&utm_content=DOCS

// package main

// import (
// 	"bufio"
// 	"bytes"
// 	"crypto/sha256"
// 	"encoding/hex"
// 	"encoding/json"
// 	"fmt"
// 	"io"
// 	"net/http"
// 	"os"
// 	"path/filepath"
// 	"strings"
// 	"sync/atomic"
// )

// type Config struct {
// 	GithubAPIToken         string `json:"github_api_token"`
// 	TerraformCloudAPIToken string `json:"terraform_cloud_api_token"`
// 	OrganizationName       string `json:"organization_name"`
// 	ProviderName           string `json:"provider_name"`
// 	RepoDetails            string `json:"repo_details"`
// 	Version                string `json:"version"`
// 	GpgPublicKey           string `json:"gpg_public_key"`
// }

// // Define a struct to unmarshal the GitHub API response
// type GitHubRelease struct {
// 	Assets []struct {
// 		Name               string `json:"name"`
// 		BrowserDownloadURL string `json:"browser_download_url"`
// 	} `json:"assets"`
// }

// // PlatformDetails holds information about the target operating system and architecture for a binary file.
// type PlatformDetails struct {
// 	OS   string // Operating System (e.g., linux, darwin, windows)
// 	Arch string // Architecture (e.g., amd64, 386, arm64)
// }

// // Define necessary structs for dynamic JSON creation
// type VersionPayload struct {
// 	Data struct {
// 		Type       string `json:"type"`
// 		Attributes struct {
// 			Version   string   `json:"version"`
// 			Protocols []string `json:"protocols"`
// 			KeyID     string   `json:"key-id,omitempty"` // Include this only if you're associating a GPG key with each version
// 		} `json:"attributes"`
// 	} `json:"data"`
// }

// type PlatformPayload struct {
// 	Data struct {
// 		Type       string `json:"type"`
// 		Attributes struct {
// 			OS       string `json:"os"`
// 			Arch     string `json:"arch"`
// 			Filename string `json:"filename"`
// 			Shasum   string `json:"shasum"`
// 		} `json:"attributes"`
// 	} `json:"data"`
// }

// // GitHubReleaseAsset represents an asset of a GitHub release.
// type GitHubReleaseAsset struct {
// 	Name               string `json:"name"`
// 	BrowserDownloadURL string `json:"browser_download_url"`
// }

// func main() {
// 	const configFilePath = "config.json"

// 	config, err := loadConfigFromFile(configFilePath)
// 	if err != nil || config == nil {
// 		config = &Config{}
// 	}

// 	reader := bufio.NewReader(os.Stdin)
// 	prompts := map[string]*string{
// 		"GitHub API Token: ": &config.GithubAPIToken,
// 		"Terraform Cloud API Token (team not org token): ":    &config.TerraformCloudAPIToken,
// 		"Terraform Organization Name: ":                       &config.OrganizationName,
// 		"Terraform Provider Name: ":                           &config.ProviderName,
// 		"Source GitHub Repository (format: owner/repoName): ": &config.RepoDetails,
// 		"Terraform Provider Version: ":                        &config.Version,
// 		"GPG Public Key: ":                                    &config.GpgPublicKey,
// 	}

// 	for prompt, value := range prompts {
// 		if *value == "" {
// 			fmt.Print(prompt)
// 			input, err := reader.ReadString('\n')
// 			if err != nil {
// 				fmt.Printf("Error reading input: %v\n", err)
// 				os.Exit(1)
// 			}
// 			*value = strings.TrimSpace(input)
// 		}
// 	}

// 	// Save the updated configuration for future use
// 	err = saveConfigToFile(config, configFilePath)
// 	if err != nil {
// 		fmt.Printf("Error saving configuration: %v\n", err)
// 		os.Exit(1)
// 	}

// 	// Explicitly assign each configuration value to a specific variable for clarity
// 	githubAPIToken := config.GithubAPIToken
// 	terraformCloudAPIToken := config.TerraformCloudAPIToken
// 	organizationName := config.OrganizationName
// 	providerName := config.ProviderName
// 	repoDetails := strings.Split(config.RepoDetails, "/")
// 	version := config.Version
// 	gpgPublicKey := config.GpgPublicKey

// 	// Validate repository details format
// 	if len(repoDetails) != 2 {
// 		fmt.Println("Invalid repository details. Please specify as <owner>/<repoName>.")
// 		os.Exit(1)
// 	}
// 	repoOwner, repoName := repoDetails[0], repoDetails[1]

// 	// Step 1: Fetch tf provider release information from GitHub
// 	release := fetchGithubReleaseInfo(repoOwner, repoName, version, githubAPIToken)
// 	if release == nil {
// 		fmt.Println("Failed to fetch release info from GitHub.")
// 		return
// 	}

// 	// Step 2: Download all assets found in the GitHub release
// 	for _, asset := range release.Assets {
// 		fmt.Printf("Downloading asset: %s\n", asset.Name)
// 		// Use downloadAsset with version parameter to manage downloads
// 		if err := downloadAsset(asset.BrowserDownloadURL, asset.Name, version); err != nil {
// 			fmt.Printf("Failed to download asset %s: %v\n", asset.Name, err)

// 		}
// 	}

// 	// Step 3: Create the provider on Terraform registry if not already present
// 	// https://developer.hashicorp.com/terraform/cloud-docs/registry/publish-providers#create-the-provider
// 	createProvider(terraformCloudAPIToken, organizationName, providerName)

// 	// Step 4: Upload the GPG public key to Terraform registry from the file
// 	// https://developer.hashicorp.com/terraform/cloud-docs/registry/publish-providers#add-your-public-key
// 	keyID, err := uploadGPGSigningKey(terraformCloudAPIToken, organizationName, gpgPublicKey)
// 	if err != nil {
// 		fmt.Printf("Error uploading GPG signing key: %v\n", err)
// 	}

// 	// Step 5: create a provider version and platform
// 	// Call createProviderVersion to create a version for your provider
// 	shasumsUploadURL, shasumsSigUploadURL, err := createProviderVersion(terraformCloudAPIToken, organizationName, providerName, version, keyID)
// 	if err != nil {
// 		fmt.Printf("Error creating provider version: %v\n", err)
// 	}

// 	// Step 6: Upload the SHA256SUMS and SHA256SUMS.sig files to the Terraform Registry
// 	assetsDir := "release-assets" // The directory where your assets are stored
// 	if err := uploadChecksumFiles(providerName, shasumsUploadURL, shasumsSigUploadURL, version, assetsDir); err != nil {
// 		fmt.Printf("Error uploading checksum files: %v\n", err)
// 	}

// 	// Step 7: Create platforms for each binary file and upload the binaries
// 	err = createProviderPlatforms(terraformCloudAPIToken, organizationName, providerName, version, assetsDir)
// 	if err != nil {
// 		fmt.Printf("Error creating provider platforms: %v\n", err)
// 	}

// 	// Use providerBinaryUploadURL as needed
// 	fmt.Println("Provider binary uploaded successfully.")
// }

// func loadConfigFromFile(filePath string) (*Config, error) {
// 	file, err := os.ReadFile(filePath)
// 	if err != nil {
// 		return nil, err
// 	}
// 	var config Config
// 	err = json.Unmarshal(file, &config)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &config, nil
// }

// func saveConfigToFile(config *Config, filePath string) error {
// 	file, err := json.MarshalIndent(config, "", " ")
// 	if err != nil {
// 		return err
// 	}
// 	return os.WriteFile(filePath, file, 0644)
// }

// // fetchReleaseInfo makes an HTTP GET request to the GitHub API to fetch details of a specific release by its tag name.
// // It takes the repository owner's name, the repository name, the release version, and a GitHub API token as parameters.
// // The function constructs the GitHub API URL using the repository owner's name, repository name, and release version.
// // It then sets up an HTTP request, adding the Authorization header with the provided GitHub API token to authenticate the request.
// // The function sends the request using an HTTP client and checks the response status code. If the response indicates an error (status code >= 400),
// // it logs the HTTP error and returns nil, indicating that the release information could not be fetched.
// // If the request is successful, the function reads the response body and attempts to unmarshal the JSON content into a GitHubRelease struct.
// // The GitHubRelease struct is designed to match the expected JSON structure of the GitHub API response, focusing on the 'assets' part of a release,
// // which includes the name and download URL of each asset.
// // If the JSON unmarshalling is successful, the function returns a pointer to the populated GitHubRelease struct, providing access to the release assets' details.
// // In case of any errors during the process (such as an HTTP request error, a response body read error, or a JSON unmarshalling error),
// // the function logs the error and returns nil. This allows the caller to check for a nil return value to determine if the function succeeded.
// func fetchGithubReleaseInfo(repoOwner, repoName, version, apiToken string) *GitHubRelease {
// 	// Construct the GitHub API URL for fetching release information by tag name
// 	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/tags/%s", repoOwner, repoName, version)

// 	// Create a new HTTP GET request
// 	req, err := http.NewRequest("GET", url, nil)
// 	if err != nil {
// 		fmt.Printf("Error creating request: %v\n", err)
// 		return nil
// 	}

// 	// Add the Authorization header with the GitHub API token
// 	req.Header.Set("Authorization", "token "+apiToken)
// 	// Specify the media type for the GitHub API version
// 	req.Header.Set("Accept", "application/vnd.github.v3+json")

// 	// Initialize an HTTP client and send the request
// 	client := &http.Client{}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		fmt.Printf("Error making request: %v\n", err)
// 		return nil
// 	}
// 	defer resp.Body.Close()

// 	// Check for HTTP error response
// 	if resp.StatusCode >= 400 {
// 		fmt.Printf("HTTP error: %s\n", resp.Status)
// 		return nil
// 	}

// 	// Read the response body
// 	body, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		fmt.Printf("Error reading response body: %v\n", err)
// 		return nil
// 	}

// 	// Unmarshal the JSON response into the GitHubRelease struct
// 	var release GitHubRelease
// 	if err := json.Unmarshal(body, &release); err != nil {
// 		fmt.Printf("Error unmarshalling response: %v\n", err)
// 		return nil
// 	}

// 	// Print the names of the assets to the terminal
// 	fmt.Print("Assets found in release: ")
// 	for _, asset := range release.Assets {
// 		fmt.Printf("%s, ", asset.Name)
// 	}
// 	fmt.Println() // Print a newline after listing all assets

// 	// Return the populated GitHubRelease struct
// 	return &release
// }

// // downloadAsset downloads a single asset from a given URL and saves it with the specified filename.
// // It checks for file existence and creates necessary directories based on the version number.
// // It includes a progress meter to display the download progress in kilobytes.
// func downloadAsset(url, filename, version string) error {
// 	// Create the directory structure based on the version number
// 	dirPath := filepath.Join("release-assets", version)
// 	if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
// 		return fmt.Errorf("failed to create directory %s: %v", dirPath, err)
// 	}

// 	// Construct the full file path
// 	fullPath := filepath.Join(dirPath, filename)

// 	// Check if the file already exists
// 	if _, err := os.Stat(fullPath); err == nil {
// 		fmt.Printf("File already exists, skipping download: %s\n", fullPath)
// 		return nil // File exists, skip download
// 	} else if !os.IsNotExist(err) {
// 		return fmt.Errorf("failed to check if file exists %s: %v", fullPath, err)
// 	}

// 	resp, err := http.Get(url)
// 	if err != nil {
// 		return fmt.Errorf("failed to download asset from %s: %v", url, err)
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusOK {
// 		return fmt.Errorf("received non-200 response status: %d %s", resp.StatusCode, resp.Status)
// 	}

// 	out, err := os.Create(fullPath)
// 	if err != nil {
// 		return fmt.Errorf("failed to create file %s: %v", fullPath, err)
// 	}
// 	defer out.Close()

// 	// Set up progress tracking
// 	counter := &writeCounter{}
// 	proxyReader := &proxyReader{reader: resp.Body, counter: counter}

// 	_, err = io.Copy(out, proxyReader)
// 	if err != nil {
// 		return fmt.Errorf("failed to write data to file %s: %v", fullPath, err)
// 	}

// 	fmt.Printf("\nAsset downloaded successfully: %s\n", fullPath)
// 	return nil
// }

// // writeCounter tracks the number of bytes transferred and prints progress.
// type writeCounter struct {
// 	total uint64
// }

// // Write increments the counter and prints the current progress.
// func (wc *writeCounter) Write(p []byte) (int, error) {
// 	n := len(p)
// 	atomic.AddUint64(&wc.total, uint64(n))
// 	wc.PrintProgress()
// 	return n, nil
// }

// // PrintProgress prints the current progress to stdout.
// func (wc *writeCounter) PrintProgress() {
// 	// Convert bytes to kilobytes for display
// 	fmt.Printf("\rDownloaded: %d KB", atomic.LoadUint64(&wc.total)/1024)
// }

// // proxyReader wraps an io.Reader and reports progress through a writeCounter.
// type proxyReader struct {
// 	reader  io.Reader
// 	counter *writeCounter
// }

// // Read reads data from the wrapped reader, reports progress, and passes the data through.
// func (pr *proxyReader) Read(p []byte) (int, error) {
// 	n, err := pr.reader.Read(p)
// 	if n > 0 {
// 		_, _ = pr.counter.Write(p[:n]) // Ignoring errors from the counter as it's only for progress reporting
// 	}
// 	return n, err
// }

// // createProvider registers a new provider in the Terraform Cloud's private registry for a given organization.
// // It requires an 'apiToken' for authentication, 'organizationName' to specify which organization the provider belongs to,
// // and 'jsonFilePath' which is the path to a JSON file containing the provider details.
// // The function constructs an HTTP POST request to the Terraform Cloud API, sending the contents of the JSON file as the request body.
// // On success, it logs a message indicating the provider was created successfully along with any additional response details.
// // If an error occurs during any step (file reading, request creation, or API communication), it logs the error and exits.
// func createProvider(apiToken, organizationName, providerName string) {
// 	// Ensure the URL matches the documentation format
// 	url := fmt.Sprintf("https://app.terraform.io/api/v2/organizations/%s/registry-providers", organizationName)

// 	// Prepare the JSON payload according to the documentation for a private provider
// 	payload := map[string]interface{}{
// 		"data": map[string]interface{}{
// 			"type": "registry-providers",
// 			"attributes": map[string]interface{}{
// 				"name":          providerName,
// 				"namespace":     organizationName, // For private providers, this is the same as the organization name
// 				"registry-name": "private",        // Specify 'private' for a private provider
// 			},
// 		},
// 	}

// 	// Convert the payload to JSON format
// 	payloadBytes, err := json.Marshal(payload)
// 	if err != nil {
// 		fmt.Printf("Failed to marshal payload: %v\n", err)
// 		return
// 	}

// 	// Create the HTTP POST request with the JSON payload
// 	request, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
// 	if err != nil {
// 		fmt.Printf("Failed to create request: %v\n", err)
// 		return
// 	}

// 	// Set the necessary headers
// 	request.Header.Set("Authorization", "Bearer "+apiToken)
// 	request.Header.Set("Content-Type", "application/vnd.api+json")

// 	// Initialize an HTTP client and send the request
// 	client := &http.Client{}
// 	response, err := client.Do(request)
// 	if err != nil {
// 		fmt.Printf("Failed to send request: %v\n", err)
// 		return
// 	}
// 	defer response.Body.Close()

// 	// Check the response status code
// 	if response.StatusCode != 201 { // Expecting a 201 status code for successful creation
// 		responseBody, _ := io.ReadAll(response.Body) // Read the response body for more details on the error
// 		fmt.Printf("Failed to create provider, HTTP error: %s, Detail: %s\n", response.Status, string(responseBody))
// 		return
// 	}

// 	// Read and log the response body for confirmation
// 	responseBody, err := io.ReadAll(response.Body)
// 	if err != nil {
// 		fmt.Printf("Failed to read response body: %v\n", err)
// 		return
// 	}
// 	fmt.Printf("Provider created successfully, response: %s\n", string(responseBody))
// }

// // uploadGPGSigningKey uploads a GPG signing key to the Terraform Cloud's private registry
// // for a given organization.
// // This function is essential for ensuring that releases are signed with a GPG key, enhancing
// // the security and integrity of the provider packages.
// //
// // Parameters:
// // - apiToken: A string representing the API token used for authenticating with the Terraform
// // Cloud API. This token should have the necessary permissions to manage GPG keys within the
// // specified organization.
// // - organizationName: A string representing the name of the organization in Terraform Cloud
// // where the GPG key will be uploaded. This organization should already exist in Terraform Cloud,
// //
// //	and the API token should be associated with a user who has permissions to manage GPG keys
// //
// // in this organization.
// // - publicKey: A string containing the ASCII armored representation of the GPG public key.
// // This key is used to sign provider releases, and its upload enables Terraform Cloud to verify
// // the authenticity of the signed files. The ASCII armored format includes the header and footer
// // lines (-----BEGIN PGP PUBLIC KEY BLOCK----- and -----END PGP PUBLIC KEY BLOCK-----, respectively),
// // along with the base64-encoded public key content.
// //
// // The function constructs a JSON payload that includes the organization name and the ASCII
// // armored public key. It then sends a POST request to the Terraform Cloud API endpoint
// // responsible for adding GPG keys. The 'Authorization' header of the request includes
// // the API token for authentication, and the 'Content-Type' header is set to
// // 'application/vnd.api+json' to indicate the format of the request body.
// //
// // Upon a successful request, the function prints a success message indicating that the GPG key
// //
// //	was uploaded successfully. If the request fails due to client-side issues (e.g., invalid
// //
// // parameters or payload) or server-side issues (e.g., authentication failure, insufficient
// // permissions, or other API errors), the function prints an error message with details about
// // the failure.
// //
// // Note: This function does not return any value or error. It handles all success and error
// // logging internally. Consider modifying the function if you need to handle errors or responses
// // differently in your application context.
// func uploadGPGSigningKey(apiToken, organizationName, publicKeyFilePath string) (string, error) {
// 	file, err := os.Open(publicKeyFilePath)
// 	if err != nil {
// 		return "", fmt.Errorf("failed to open public key file: %v", err)
// 	}
// 	defer file.Close()

// 	publicKeyBytes, err := io.ReadAll(file)
// 	if err != nil {
// 		return "", fmt.Errorf("failed to read public key file: %v", err)
// 	}
// 	publicKey := string(publicKeyBytes)

// 	url := "https://app.terraform.io/api/registry/private/v2/gpg-keys"

// 	payload := map[string]interface{}{
// 		"data": map[string]interface{}{
// 			"type": "gpg-keys",
// 			"attributes": map[string]interface{}{
// 				"namespace":   organizationName,
// 				"ascii-armor": publicKey,
// 			},
// 		},
// 	}

// 	payloadBytes, err := json.Marshal(payload)
// 	if err != nil {
// 		return "", fmt.Errorf("failed to marshal payload: %v", err)
// 	}

// 	request, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
// 	if err != nil {
// 		return "", fmt.Errorf("failed to create request: %v", err)
// 	}

// 	request.Header.Set("Authorization", "Bearer "+apiToken)
// 	request.Header.Set("Content-Type", "application/vnd.api+json")

// 	client := &http.Client{}
// 	response, err := client.Do(request)
// 	if err != nil {
// 		return "", fmt.Errorf("failed to send request: %v", err)
// 	}
// 	defer response.Body.Close()

// 	responseBodyBytes, err := io.ReadAll(response.Body)
// 	if err != nil {
// 		return "", fmt.Errorf("failed to read response body: %v", err)
// 	}

// 	if response.StatusCode >= 400 {
// 		return "", fmt.Errorf("failed to add GPG key, HTTP error: %s, Detail: %s", response.Status, string(responseBodyBytes))
// 	}

// 	var responseJSON struct {
// 		Data struct {
// 			ID string `json:"id"`
// 		} `json:"data"`
// 	}
// 	if err := json.Unmarshal(responseBodyBytes, &responseJSON); err != nil {
// 		return "", fmt.Errorf("failed to parse response JSON: %v", err)
// 	}

// 	return responseJSON.Data.ID, nil
// }

// // createProviderVersion creates a new version for a provider in the Terraform Cloud's private registry.
// // It requires an API token for authentication, the name of the Terraform Cloud organization,
// // the provider name, the version of the provider to create, and the ID of the GPG key associated with the provider.
// // The function constructs a JSON payload with these details and sends a POST request to the Terraform Cloud API.
// // On a successful request, it parses the response to extract URLs for uploading the SHA256SUMS and SHA256SUMS.sig files.
// // These URLs are then printed to the console.
// // The function returns the upload URLs and any error encountered during the process.
// func createProviderVersion(apiToken, organizationName, providerName, version, keyID string) (shasumsUploadURL, shasumsSigUploadURL string, err error) {
// 	// Ensure the URL is correctly formatted according to the documentation
// 	url := fmt.Sprintf("https://app.terraform.io/api/v2/organizations/%s/registry-providers/private/%s/%s/versions", organizationName, organizationName, providerName)

// 	// Update the keyID to a valid GPG key string if necessary
// 	keyID = "DB95CA76A94A208C" // Example, replace with a valid GPG key ID if necessary

// 	// Remove the prefixed "v" from the version parameter if it exists
// 	providerVersion := strings.TrimPrefix(version, "v")

// 	// Prepare the JSON payload according to the documentation
// 	payload := map[string]interface{}{
// 		"data": map[string]interface{}{
// 			"type": "registry-provider-versions",
// 			"attributes": map[string]interface{}{
// 				"version":   providerVersion,
// 				"key-id":    keyID,
// 				"protocols": []string{"5.0"}, // Ensure this matches the required or supported protocol versions
// 			},
// 		},
// 	}

// 	// Marshal the payload to JSON
// 	payloadBytes, err := json.Marshal(payload)
// 	if err != nil {
// 		return "", "", fmt.Errorf("failed to marshal payload: %v", err)
// 	}

// 	// Create the HTTP POST request
// 	request, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
// 	if err != nil {
// 		return "", "", fmt.Errorf("failed to create request: %v", err)
// 	}

// 	// Set the necessary headers
// 	request.Header.Set("Authorization", "Bearer "+apiToken)
// 	request.Header.Set("Content-Type", "application/vnd.api+json")

// 	// Initialize an HTTP client and send the request
// 	client := &http.Client{}
// 	response, err := client.Do(request)
// 	if err != nil {
// 		return "", "", fmt.Errorf("failed to send request: %v", err)
// 	}
// 	defer response.Body.Close()

// 	// Read the response body
// 	responseBodyBytes, err := io.ReadAll(response.Body)
// 	if err != nil {
// 		return "", "", fmt.Errorf("failed to read response body: %v", err)
// 	}

// 	// Check the response status code and handle errors
// 	if response.StatusCode != http.StatusCreated { // Expecting 201 status code for success
// 		return "", "", fmt.Errorf("failed to create provider version, HTTP error: %s, Detail: %s", response.Status, string(responseBodyBytes))
// 	}

// 	// Parse the response body to extract the upload URLs
// 	var responseData struct {
// 		Data struct {
// 			Links struct {
// 				ShasumsUpload    string `json:"shasums-upload"`
// 				ShasumsSigUpload string `json:"shasums-sig-upload"`
// 			} `json:"links"`
// 		} `json:"data"`
// 	}

// 	if err := json.Unmarshal(responseBodyBytes, &responseData); err != nil {
// 		return "", "", fmt.Errorf("failed to parse response: %v, Response Body: %s", err, string(responseBodyBytes))
// 	}

// 	// Print the URLs for uploading SHA256SUMS and SHA256.sig files to the console
// 	fmt.Printf("Provider version created successfully.\n")
// 	fmt.Printf("SHA256SUMS upload URL: %s\n", responseData.Data.Links.ShasumsUpload)
// 	fmt.Printf("SHA256.sig upload URL: %s\n", responseData.Data.Links.ShasumsSigUpload)

// 	// Return the URLs for uploading SHA256SUMS and SHA256.sig files
// 	return responseData.Data.Links.ShasumsUpload, responseData.Data.Links.ShasumsSigUpload, nil
// }

// func createProviderPlatforms(apiToken, organizationName, providerName, version, assetsDir string) error {
// 	providerVersion := strings.TrimPrefix(version, "v")
// 	platformURL := fmt.Sprintf("https://app.terraform.io/api/v2/organizations/%s/registry-providers/private/%s/%s/versions/%s/platforms", organizationName, organizationName, providerName, providerVersion)

// 	// List all files in the assets directory, including zip files and manifest json files
// 	assetFiles, err := listAssetFiles(assetsDir, version)
// 	if err != nil {
// 		return fmt.Errorf("error listing asset files: %v", err)
// 	}

// 	// Calculate SHA256 checksums for all asset files
// 	checksums, err := calculateSHA256ForFiles(assetFiles)
// 	if err != nil {
// 		return fmt.Errorf("error calculating SHA256 checksums: %v", err)
// 	}

// 	// Iterate over asset files to create and upload platforms for each one
// 	for _, filePath := range assetFiles {
// 		filename := filepath.Base(filePath)
// 		if strings.HasSuffix(filename, "_manifest.json") {
// 			// Handle manifest json file upload
// 			fmt.Printf("Uploading manifest file: %s\n", filename)
// 			uploadURL := fmt.Sprintf("https://app.terraform.io/api/v2/organizations/%s/registry-providers/private/%s/%s/versions/%s/manifest", organizationName, organizationName, providerName, providerVersion)
// 			err = uploadBinary(uploadURL, filePath)
// 			if err != nil {
// 				return fmt.Errorf("error uploading manifest file: %v", err)
// 			}
// 			fmt.Printf("Manifest file uploaded successfully: %s\n", filename)
// 		} else {
// 			// Handle binary file upload
// 			os, arch, _, err := extractPlatformDetails(filename)
// 			if err != nil {
// 				fmt.Printf("Skipping file due to error extracting details: %v\n", err)
// 				continue
// 			}

// 			// Prepare the payload with the extracted details and SHA256 checksum
// 			payload := PlatformPayload{
// 				Data: struct {
// 					Type       string `json:"type"`
// 					Attributes struct {
// 						OS       string `json:"os"`
// 						Arch     string `json:"arch"`
// 						Filename string `json:"filename"`
// 						Shasum   string `json:"shasum"`
// 					} `json:"attributes"`
// 				}{
// 					Type: "registry-provider-version-platforms",
// 					Attributes: struct {
// 						OS       string `json:"os"`
// 						Arch     string `json:"arch"`
// 						Filename string `json:"filename"`
// 						Shasum   string `json:"shasum"`
// 					}{
// 						OS:       os,
// 						Arch:     arch,
// 						Filename: filename,
// 						Shasum:   checksums[filePath],
// 					},
// 				},
// 			}

// 			payloadBytes, err := json.Marshal(payload)
// 			if err != nil {
// 				return fmt.Errorf("error marshaling payload: %v", err)
// 			}

// 			req, err := http.NewRequest("POST", platformURL, bytes.NewBuffer(payloadBytes))
// 			if err != nil {
// 				return fmt.Errorf("error creating request: %v", err)
// 			}

// 			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiToken))
// 			req.Header.Set("Content-Type", "application/vnd.api+json")

// 			client := &http.Client{}
// 			resp, err := client.Do(req)
// 			if err != nil {
// 				return fmt.Errorf("error sending request: %v", err)
// 			}
// 			defer resp.Body.Close()

// 			responseBodyBytes, err := io.ReadAll(resp.Body)
// 			if err != nil {
// 				return fmt.Errorf("failed to read response body: %v", err)
// 			}

// 			if resp.StatusCode != http.StatusCreated {
// 				return fmt.Errorf("error creating provider platform, status: %d, response: %s", resp.StatusCode, string(responseBodyBytes))
// 			}

// 			var responsePayload struct {
// 				Data struct {
// 					Links struct {
// 						ProviderBinaryUpload string `json:"provider-binary-upload"`
// 					} `json:"links"`
// 				} `json:"data"`
// 			}

// 			if err := json.Unmarshal(responseBodyBytes, &responsePayload); err != nil {
// 				return fmt.Errorf("error parsing response JSON: %v", err)
// 			}

// 			fmt.Printf("Platform created successfully for %s/%s. Binary upload URL: %s, Platform Binary file name: %s\n", os, arch, responsePayload.Data.Links.ProviderBinaryUpload, filename)

// 			// Now, upload the binary file associated with this platform
// 			uploadURL := responsePayload.Data.Links.ProviderBinaryUpload
// 			err = uploadBinary(uploadURL, filePath)
// 			if err != nil {
// 				return fmt.Errorf("error uploading binary file: %v", err)
// 			}
// 		}
// 	}

// 	return nil
// }

// // listAssetFiles lists both zip files and manifest json files in the assets directory for a given version
// func listAssetFiles(dir, version string) ([]string, error) {
// 	var assetFiles []string
// 	versionedDir := filepath.Join(dir, version) // Construct the path to the versioned directory

// 	err := filepath.Walk(versionedDir, func(path string, info os.FileInfo, err error) error {
// 		if err != nil {
// 			return err
// 		}

// 		// Check if the file is not a directory and has the desired suffixes
// 		if !info.IsDir() && (strings.HasSuffix(info.Name(), ".zip") || strings.HasSuffix(info.Name(), "_manifest.json")) {
// 			assetFiles = append(assetFiles, path) // Add the file path to the list
// 		}

// 		return nil
// 	})

// 	return assetFiles, err
// }

// func uploadBinary(uploadURL, filePath string) error {
// 	file, err := os.Open(filePath)
// 	if err != nil {
// 		return err
// 	}
// 	defer file.Close()

// 	// Wrap the file in a proxyReader to track upload progress
// 	counter := &writeCounter{}
// 	proxyReader := &proxyReader{reader: file, counter: counter}

// 	req, err := http.NewRequest("PUT", uploadURL, proxyReader)
// 	if err != nil {
// 		return err
// 	}

// 	client := &http.Client{}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		return err
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusOK {
// 		return fmt.Errorf("error uploading binary file, status: %d", resp.StatusCode)
// 	}

// 	fmt.Println("\nBinary uploaded successfully")
// 	return nil
// }

// func extractPlatformDetails(filename string) (os, arch, filenameOnly string, err error) {
// 	// Split the filename by underscores
// 	parts := strings.Split(filename, "_")
// 	if len(parts) < 3 {
// 		return "", "", "", fmt.Errorf("filename does not follow the expected pattern: %s", filename)
// 	}

// 	// Extract OS and Arch from filename
// 	os = parts[len(parts)-2]
// 	archWithExt := parts[len(parts)-1]
// 	arch = strings.Split(archWithExt, ".")[0] // Remove the .zip part to get the architecture

// 	// Return the OS, Arch, and the original filename
// 	return os, arch, filename, nil
// }

// func listZipFiles(dir string) ([]string, error) {
// 	var zipFiles []string

// 	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
// 		if err != nil {
// 			return err
// 		}

// 		if !info.IsDir() && strings.HasSuffix(info.Name(), ".zip") {
// 			zipFiles = append(zipFiles, path)
// 		}

// 		return nil
// 	})

// 	return zipFiles, err
// }

// func calculateSHA256ForFiles(filePaths []string) (map[string]string, error) {
// 	checksums := make(map[string]string)

// 	for _, filePath := range filePaths {
// 		file, err := os.Open(filePath)
// 		if err != nil {
// 			return nil, err
// 		}
// 		defer file.Close()

// 		hash := sha256.New()
// 		if _, err := io.Copy(hash, file); err != nil {
// 			return nil, err
// 		}

// 		checksums[filePath] = hex.EncodeToString(hash.Sum(nil))
// 	}

// 	return checksums, nil
// }

// func uploadChecksumFiles(providerName, shasumsUploadURL, shasumsSigUploadURL, version, assetsDir string) error {
// 	// Remove the "v" prefix from the version for file names
// 	providerVersion := strings.TrimPrefix(version, "v")

// 	// Construct the file names without the "v" prefix
// 	shaFileName := fmt.Sprintf("%s_%s_SHA256SUMS", providerName, providerVersion)
// 	sigFileName := fmt.Sprintf("%s_%s_SHA256SUMS.sig", providerName, providerVersion)

// 	// Construct the full file paths using the original version with the "v" prefix for the directory
// 	shaFilePath := filepath.Join(assetsDir, version, shaFileName)
// 	sigFilePath := filepath.Join(assetsDir, version, sigFileName)

// 	// Upload the SHA256SUMS file
// 	if err := uploadFileToURL(shaFilePath, shasumsUploadURL); err != nil {
// 		return fmt.Errorf("failed to upload SHA256SUMS file: %v", err)
// 	}

// 	// Upload the SHA256SUMS.sig file
// 	if err := uploadFileToURL(sigFilePath, shasumsSigUploadURL); err != nil {
// 		return fmt.Errorf("failed to upload SHA256SUMS.sig file: %v", err)
// 	}

// 	fmt.Println("Successfully uploaded SHA256SUMS and SHA256SUMS.sig files.")
// 	return nil
// }

// func uploadFileToURL(filePath, uploadURL string) error {
// 	file, err := os.Open(filePath)
// 	if err != nil {
// 		return fmt.Errorf("failed to open file %s: %v", filePath, err)
// 	}
// 	defer file.Close()

// 	req, err := http.NewRequest("PUT", uploadURL, file)
// 	if err != nil {
// 		return fmt.Errorf("failed to create PUT request for %s: %v", filePath, err)
// 	}

// 	client := &http.Client{}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		return fmt.Errorf("failed to upload file %s: %v", filePath, err)
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode >= 400 {
// 		return fmt.Errorf("failed to upload file %s, received HTTP status %s", filePath, resp.Status)
// 	}

// 	return nil
// }
