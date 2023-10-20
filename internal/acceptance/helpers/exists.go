package helpers

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/clients"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func DoesNotExistInJamfPro(client *clients.Client, testResource TestResource, resourceName string) schema.TestCheckFunc {
	return existsFunc(false)(client, testResource, resourceName)
}

func ExistsInJamfPro(client *clients.Client, testResource TestResource, resourceName string) schema.TestCheckFunc {
	return existsFunc(true)(client, testResource, resourceName)
}

func existsFunc(shouldExist bool) func(*clients.Client, TestResource, string) schema.TestCheckFunc {
	return func(client *clients.Client, testResource TestResource, resourceName string) schema.TestCheckFunc {
		return func(s *terraform.State) error {
			// even with rate limiting - an exists function should never take more than 5m, so should be safe
			ctx, cancel := context.WithDeadline(client.StopContext, time.Now().Add(5*time.Minute))
			defer cancel()

			rs, ok := s.RootModule().Resources[resourceName]
			if !ok {
				return fmt.Errorf("%q was not found in the state", resourceName)
			}

			result, err := testResource.Exists(ctx, client, rs.Primary)
			if err != nil {
				return fmt.Errorf("running exists func for %q: %+v", resourceName, err)
			}
			if result == nil {
				return fmt.Errorf("received nil for exists for %q", resourceName)
			}

			if *result != shouldExist {
				if !shouldExist {
					return fmt.Errorf("%q still exists", resourceName)
				}

				return fmt.Errorf("%q did not exist", resourceName)
			}

			return nil
		}
	}
}

// Placeholder for TestResource interface - You should define the behavior of this interface based on your requirements.
type TestResource interface {
	Exists(ctx context.Context, client *clients.Client, state *terraform.InstanceState) (*bool, error)
}
