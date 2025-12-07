package jamf_cloud_ip_address_list

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Read fetches the Jamf Cloud IP address list and applies any configured filters to the results.
func (d *JamfCloudIPAddressListDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state JamfCloudIPAddressListDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ipList, err := getJamfCloudIPList(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to get Jamf Cloud IP address list",
			fmt.Sprintf("Error fetching IP list: %s", err),
		)
		return
	}

	serviceFilter := state.ServiceFilter.ValueString()
	providerFilter := state.ProviderFilter.ValueString()
	trafficFilter := state.TrafficFilter.ValueString()
	regionFilter := state.RegionFilter.ValueString()

	filteredIPs := filterPublicIPs(ipList.PublicIPs, serviceFilter, providerFilter, trafficFilter, regionFilter)

	state.PublishDate = types.StringValue(ipList.PublishDate)
	state.PublicIPs = make([]PublicIPEntryModel, 0, len(filteredIPs))

	for _, entry := range filteredIPs {
		ipPrefixes := make([]types.String, len(entry.IPPrefixes))
		for i, ip := range entry.IPPrefixes {
			ipPrefixes[i] = types.StringValue(ip)
		}

		fqdns := make([]types.String, len(entry.FQDNs))
		for i, fqdn := range entry.FQDNs {
			fqdns[i] = types.StringValue(fqdn)
		}

		state.PublicIPs = append(state.PublicIPs, PublicIPEntryModel{
			Service:    types.StringValue(entry.Service),
			Provider:   types.StringValue(entry.Provider),
			Traffic:    types.StringValue(entry.Traffic),
			Region:     types.StringValue(entry.Region),
			IPPrefixes: ipPrefixes,
			FQDNs:      fqdns,
		})
	}

	idStr := fmt.Sprintf("%s-%s-%s-%s-%s", serviceFilter, providerFilter, trafficFilter, regionFilter, time.Now().UTC().String())
	state.ID = types.StringValue(fmt.Sprintf("%x", sha256.Sum256([]byte(idStr))))

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// getJamfCloudIPList retrieves the IP address list from the Jamf public URL.
func getJamfCloudIPList(ctx context.Context) (*ResponseJamfCloudIPAddressList, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, jamfCloudIPListURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get IP list: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var ipList ResponseJamfCloudIPAddressList
	if err := json.Unmarshal(body, &ipList); err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	return &ipList, nil
}

// filterPublicIPs filters the public IP entries based on the provided filter criteria.
func filterPublicIPs(entries []ResponsePublicIPEntry, service, provider, traffic, region string) []ResponsePublicIPEntry {
	filtered := make([]ResponsePublicIPEntry, 0)
	for _, entry := range entries {
		if (service == "" || entry.Service == service) &&
			(provider == "" || entry.Provider == provider) &&
			(traffic == "" || entry.Traffic == traffic) &&
			(region == "" || entry.Region == region) {
			filtered = append(filtered, entry)
		}
	}
	return filtered
}
