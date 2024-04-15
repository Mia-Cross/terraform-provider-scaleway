package meta_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/meta"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/provider"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/transport"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func NewMockTestTools(t *testing.T) *acctest.TestTools {
	t.Helper()
	ctx := context.Background()

	folder, err := os.Getwd()
	if err != nil {
		t.Fatalf("cannot detect working directory for testing")
	}

	// Create a http client with recording capabilities
	httpClient, cleanup, err := acctest.GetHTTPRecoder(t, folder, *acctest.UpdateCassettes)
	require.NoError(t, err)

	// Create meta that will be passed in the provider config
	m, err := meta.NewMeta(ctx, &meta.Config{
		ProviderSchema:   nil,
		TerraformVersion: "terraform-tests",
		HTTPClient:       httpClient,
	})
	require.NoError(t, err)

	if !*acctest.UpdateCassettes {
		tmp := 0 * time.Second
		transport.DefaultWaitRetryInterval = &tmp
	}

	return &acctest.TestTools{
		T:    t,
		Meta: m,
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"scaleway": func() (*schema.Provider, error) {
				return MockProvider(&provider.Config{Meta: m})(), nil
			},
		},
		Cleanup: cleanup,
	}
}

// MockProvider returns a terraform.ResourceProvider.
func MockProvider(config *provider.Config) plugin.ProviderFunc {
	return func() *schema.Provider {
		p := &schema.Provider{
			Schema: map[string]*schema.Schema{
				//"access_key": {
				//	Type:        schema.TypeString,
				//	Optional:    true,
				//	Description: "The Scaleway access key.",
				//},
				//"secret_key": {
				//	Type:         schema.TypeString,
				//	Optional:     true, // To allow user to use deprecated `token`.
				//	Description:  "The Scaleway secret Key.",
				//	ValidateFunc: verify.IsUUID(),
				//},
				//"profile": {
				//	Type:        schema.TypeString,
				//	Optional:    true, // To allow user to use `access_key`, `secret_key`, `project_id`...
				//	Description: "The Scaleway profile to use.",
				//},
				//"project_id": {
				//	Type:         schema.TypeString,
				//	Optional:     true, // To allow user to use organization instead of project
				//	Description:  "The Scaleway project ID.",
				//	ValidateFunc: verify.IsUUID(),
				//},
				//"organization_id": {
				//	Type:         schema.TypeString,
				//	Optional:     true,
				//	Description:  "The Scaleway organization ID.",
				//	ValidateFunc: verify.IsUUID(),
				//},
				//"region": regional.Schema(),
				//"zone":   zonal.Schema(),
				//"api_url": {
				//	Type:        schema.TypeString,
				//	Optional:    true,
				//	Description: "The Scaleway API URL to use.",
				//},
			},
			ResourcesMap: map[string]*schema.Resource{
				"scaleway_mock": ResourceMock(),
			},
		}
		return p
	}
}

const mockTimeout = 30 * time.Second

func ResourceMock() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMockCreate,
		ReadContext:   resourceMockRead,
		UpdateContext: resourceMockUpdate,
		DeleteContext: resourceMockDelete,
		Timeouts: &schema.ResourceTimeout{
			Create:  schema.DefaultTimeout(mockTimeout),
			Read:    schema.DefaultTimeout(mockTimeout),
			Update:  schema.DefaultTimeout(mockTimeout),
			Delete:  schema.DefaultTimeout(mockTimeout),
			Default: schema.DefaultTimeout(mockTimeout),
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		SchemaVersion: 0,
		Schema: map[string]*schema.Schema{
			"root_string": {
				Type:     schema.TypeString,
				Required: true,
			},
			"root_int": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"root_bool": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"list": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"map": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"object": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"nested_string": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"nested_int": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"nested_bool": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
					},
				},
			},
			"not_in_config_string": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "default string value",
			},
			"not_in_config_int": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  -1,
			},
			"not_in_config_bool": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
		},
	}
}

func resourceMockCreate(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	d.SetId("scw-mock")
	Mock.RootString = d.Get("root_string").(string)
	return resourceMockRead(ctx, d, i)
}
func resourceMockRead(_ context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
	_ = d.Set("root_string", Mock.RootString)

	//_ = d.Set("root_int", d.Get("root_int"))
	//_ = d.Set("root_bool", d.Get("root_bool"))
	//_ = d.Set("list", d.Get("list"))
	//_ = d.Set("map", d.Get("map"))
	//_ = d.Set("not_in_config_string", d.Get("not_in_config_string"))
	//_ = d.Set("not_in_config_int", d.Get("not_in_config_int"))
	//_ = d.Set("not_in_config_bool", d.Get("not_in_config_bool"))
	return nil
}
func resourceMockUpdate(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	if d.HasChange("root_string") {
		Mock.RootString = d.Get("root_string").(string)
	}
	return resourceMockRead(ctx, d, i)
}
func resourceMockDelete(_ context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
	d.SetId("")
	return nil
}

type MockResource struct {
	RootString string
}

var Mock MockResource

func TestAcc_GetRawConfigForKey(t *testing.T) {
	tt := NewMockTestTools(t)
	defer tt.Cleanup()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		Steps: []resource.TestStep{
			{
				//Config: `
				//	resource scaleway_mock mock {
				//		root_string = "a random string"
				//		root_bool = false
				//		root_int = 209
				//		list = [ "terraform-test", "core", "get-raw-config" ]
				//		map = {
				//			foo = "bar"
				//			elem = "420"
				//		}
				//		object {
				//			nested_string = "another random string"
				//			nested_bool = true
				//			nested_int = 2147364647
				//		}
				//	}
				//`,
				Config: `
					resource scaleway_mock mock {
						root_string = "a random string"
					}
				`,
				Check: resource.ComposeTestCheckFunc(
				// Primitive types
				//assertGetRawConfigResults(t, "root_string", "a random string", true, cty.String),
				//assertGetRawConfigResults(t, "root_bool", false, true, cty.Bool),
				//assertGetRawConfigResults(t, "root_int", 209, true, cty.Number),
				//assertGetRawConfigResults(t, "key_not_set", nil, false, cty.Number),
				//// List
				//assertGetRawConfigResults(t, "list.0", "terraform-test", true, cty.String),
				//assertGetRawConfigResults(t, "list.1", "core", true, cty.String),
				//assertGetRawConfigResults(t, "list.2", "get-raw-config", true, cty.String),
				//assertGetRawConfigResults(t, "list.3", nil, false, cty.String),
				//// Map
				//assertGetRawConfigResults(t, "map.foo", "bar", true, cty.String),
				//assertGetRawConfigResults(t, "map.elem", "420", true, cty.String),
				//assertGetRawConfigResults(t, "map.not_in_map", nil, false, cty.String),
				//// Object
				//assertGetRawConfigResults(t, "object.0.nested_string", "another random string", true, cty.String),
				//assertGetRawConfigResults(t, "object.0.nested_bool", true, true, cty.Bool),
				//assertGetRawConfigResults(t, "object.0.nested_int", 2147483647, true, cty.Number),
				//assertGetRawConfigResults(t, "object.0.not_in_object", nil, false, cty.String),
				//assertGetRawConfigResults(t, "object.#.nested_bool", true, true, cty.Bool),
				//// Default values
				//assertGetRawConfigResults(t, "not_in_config_string", "default string value", false, cty.String),
				//assertGetRawConfigResults(t, "not_in_config_int", -1, false, cty.Number),
				//assertGetRawConfigResults(t, "not_in_config_bool", true, false, cty.Bool),
				),
			},
			{
				//Config: `
				//	resource scaleway_mock mock {
				//		root_string = "string has changed"
				//		root_bool = true
				//		root_int = 2024
				//		list = [ "terraform-test", "core", "get-raw-config-for-key" ]
				//		map = {
				//			food = "bar"
				//			elem = "240"
				//		}
				//		object {
				//			nested_string = "this string has also changed"
				//			nested_bool = false
				//			nested_int = -2147364648
				//		}
				//	}
				//`,
				Config: `
					resource scaleway_mock mock {
						root_string = "string has changed"
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					// Primitive types
					assertGetRawConfigResults(t, "root_string", "string has changed", true, cty.String),
					//assertGetRawConfigResults(t, "root_bool", true, true, cty.Bool),
					//assertGetRawConfigResults(t, "root_int", 2024, true, cty.Number),
					//assertGetRawConfigResults(t, "key_not_set", nil, false, cty.Number),
					//// List
					//assertGetRawConfigResults(t, "list.0", "terraform-test", true, cty.String),
					//assertGetRawConfigResults(t, "list.1", "core", true, cty.String),
					//assertGetRawConfigResults(t, "list.2", "get-raw-config-for-key", true, cty.String),
					//assertGetRawConfigResults(t, "list.3", nil, false, cty.String),
					//// Map
					//assertGetRawConfigResults(t, "map.food", "bar", true, cty.String),
					//assertGetRawConfigResults(t, "map.elem", "240", true, cty.String),
					//assertGetRawConfigResults(t, "map.not_in_map", nil, false, cty.String),
					//// Object
					//assertGetRawConfigResults(t, "object.0.nested_string", "this string has also changed", true, cty.String),
					//assertGetRawConfigResults(t, "object.0.nested_bool", false, true, cty.Bool),
					//assertGetRawConfigResults(t, "object.0.nested_int", -2147483648, true, cty.Number),
					//assertGetRawConfigResults(t, "object.0.not_in_object", nil, false, cty.String),
					//assertGetRawConfigResults(t, "object.#.nested_bool", false, true, cty.Bool),
					//// Default values
					//assertGetRawConfigResults(t, "not_in_config_string", "default string value", false, cty.String),
					//assertGetRawConfigResults(t, "not_in_config_int", -1, false, cty.Number),
					//assertGetRawConfigResults(t, "not_in_config_bool", true, false, cty.Bool),
				),
			},
			//{
			//	Config: `
			//		resource scaleway_mock mock {
			//			root_string = ""
			//			root_int = 1234567890
			//			list = [ "terraform-test", "get-raw-config-for-key" ]
			//			map = {
			//				foo = "bard"
			//				elem = "240"
			//			}
			//			object {
			//				nested_string = "this string has also changed"
			//				nested_int = 0
			//			}
			//		}
			//	`,
			//	Check: resource.ComposeTestCheckFunc(
			//		// Primitive types
			//		assertGetRawConfigResults(t, "root_string", "", true, cty.String),
			//		//assertGetRawConfigResults(t, "root_bool", true, false, cty.Bool),
			//		//assertGetRawConfigResults(t, "root_int", 1234567890, true, cty.Number),
			//		//assertGetRawConfigResults(t, "key_not_set", nil, false, cty.Number),
			//		//// List
			//		//assertGetRawConfigResults(t, "list.0", "terraform-test", true, cty.String),
			//		//assertGetRawConfigResults(t, "list.1", "get-raw-config-for-key", true, cty.String),
			//		//assertGetRawConfigResults(t, "list.2", nil, false, cty.String),
			//		//// Map
			//		//assertGetRawConfigResults(t, "map.foo", "bard", true, cty.String),
			//		//assertGetRawConfigResults(t, "map.elem", "240", true, cty.String),
			//		//assertGetRawConfigResults(t, "map.not_in_map", nil, false, cty.String),
			//		//// Object
			//		//assertGetRawConfigResults(t, "object.0.nested_string", "this string has also changed", true, cty.String),
			//		//assertGetRawConfigResults(t, "object.0.nested_bool", false, true, cty.Bool),
			//		//assertGetRawConfigResults(t, "object.0.nested_int", -2147483648, true, cty.Number),
			//		//assertGetRawConfigResults(t, "object.0.not_in_object", nil, false, cty.String),
			//		//assertGetRawConfigResults(t, "object.#.nested_bool", false, true, cty.Bool),
			//		//// Default values
			//		//assertGetRawConfigResults(t, "not_in_config_string", "default string value", false, cty.String),
			//		//assertGetRawConfigResults(t, "not_in_config_int", -1, false, cty.Number),
			//		//assertGetRawConfigResults(t, "not_in_config_bool", true, false, cty.Bool),
			//	),
			//},
		},
	})
}

func assertGetRawConfigResults(t *testing.T, key string, expectedValue any, expectedSet bool, ty cty.Type) resource.TestCheckFunc {
	t.Helper()
	return func(s *terraform.State) error {
		resourceTFName := "scaleway_mock.mock"
		rs, ok := s.RootModule().Resources[resourceTFName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceTFName)
		}
		if rs.Primary.RawConfig.IsNull() {
			return nil
		}

		actualValue, actualSet := meta.GetKeyInRawConfigMap(rs.Primary.RawConfig.AsValueMap(), key, ty)
		assert.Equal(t, expectedSet, actualSet)
		assert.Equal(t, expectedValue, actualValue)

		return nil
	}
}
