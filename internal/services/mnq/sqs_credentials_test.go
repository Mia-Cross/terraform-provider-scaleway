package mnq_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	mnqSDK "github.com/scaleway/scaleway-sdk-go/api/mnq/v1beta1"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/httperrors"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/services/mnq"
)

func TestAccSQSCredentials_Basic(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      isSQSCredentialsDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource scaleway_account_project main {
						name = "tf_tests_mnq_sqs_credentials_basic"
					}

					resource scaleway_mnq_sqs main {
						project_id = scaleway_account_project.main.id
					}

					resource scaleway_mnq_sqs_credentials main {
						project_id = scaleway_mnq_sqs.main.project_id
						name = "test-mnq-sqs-credentials-basic"
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					isSQSCredentialsPresent(tt, "scaleway_mnq_sqs_credentials.main"),
					acctest.CheckResourceAttrUUID("scaleway_mnq_sqs_credentials.main", "id"),
					resource.TestCheckResourceAttr("scaleway_mnq_sqs_credentials.main", "name", "test-mnq-sqs-credentials-basic"),
					resource.TestCheckResourceAttrSet("scaleway_mnq_sqs_credentials.main", "access_key"),
					resource.TestCheckResourceAttrSet("scaleway_mnq_sqs_credentials.main", "secret_key"),
				),
			},
			{
				Config: `
					resource scaleway_account_project main {
						name = "tf_tests_mnq_sqs_credentials_basic"
					}

					resource scaleway_mnq_sqs main {
						project_id = scaleway_account_project.main.id
					}

					resource scaleway_mnq_sqs_credentials main {
						project_id = scaleway_mnq_sqs.main.project_id
						name = "test-mnq-sqs-credentials-basic"
						permissions {
							can_manage = true
							can_receive = false
							can_publish = true
						}
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					isSQSCredentialsPresent(tt, "scaleway_mnq_sqs_credentials.main"),
					resource.TestCheckResourceAttr("scaleway_mnq_sqs_credentials.main", "permissions.0.can_manage", "true"),
					resource.TestCheckResourceAttr("scaleway_mnq_sqs_credentials.main", "permissions.0.can_receive", "false"),
					resource.TestCheckResourceAttr("scaleway_mnq_sqs_credentials.main", "permissions.0.can_publish", "true"),
				),
			},
			{
				Config: `
					resource scaleway_account_project main {
						name = "tf_tests_mnq_sqs_credentials_basic"
					}

					resource scaleway_mnq_sqs main {
						project_id = scaleway_account_project.main.id
					}

					resource scaleway_mnq_sqs_credentials main {
						project_id = scaleway_mnq_sqs.main.project_id
						name = "test-mnq-sqs-credentials-basic"
						permissions {
							can_manage = false
							can_receive = true
							can_publish = false
						}
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					isSQSCredentialsPresent(tt, "scaleway_mnq_sqs_credentials.main"),
					resource.TestCheckResourceAttr("scaleway_mnq_sqs_credentials.main", "permissions.0.can_manage", "false"),
					resource.TestCheckResourceAttr("scaleway_mnq_sqs_credentials.main", "permissions.0.can_receive", "true"),
					resource.TestCheckResourceAttr("scaleway_mnq_sqs_credentials.main", "permissions.0.can_publish", "false"),
				),
			},
		},
	})
}

func isSQSCredentialsPresent(tt *acctest.TestTools, n string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource not found: %s", n)
		}

		api, region, id, err := mnq.NewSQSAPIWithRegionAndID(tt.Meta, rs.Primary.ID)
		if err != nil {
			return err
		}

		_, err = api.GetSqsCredentials(&mnqSDK.SqsAPIGetSqsCredentialsRequest{
			SqsCredentialsID: id,
			Region:           region,
		})
		if err != nil {
			return err
		}

		return nil
	}
}

func isSQSCredentialsDestroyed(tt *acctest.TestTools) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		for _, rs := range state.RootModule().Resources {
			if rs.Type != "scaleway_mnq_sqs_credentials" {
				continue
			}

			api, region, id, err := mnq.NewSQSAPIWithRegionAndID(tt.Meta, rs.Primary.ID)
			if err != nil {
				return err
			}

			err = api.DeleteSqsCredentials(&mnqSDK.SqsAPIDeleteSqsCredentialsRequest{
				SqsCredentialsID: id,
				Region:           region,
			})

			if err == nil {
				return fmt.Errorf("mnq sqs credentials (%s) still exists", rs.Primary.ID)
			}

			if !httperrors.Is404(err) {
				return err
			}
		}

		return nil
	}
}
