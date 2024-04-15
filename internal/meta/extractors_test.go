package meta_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/meta"
	"github.com/stretchr/testify/assert"
)

func TestAcc_GetRawConfigForKey(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					resource scaleway_mock mock {
						root_string = "a random string"
						root_bool = false
						root_int = 209
						list = [ "terraform-test", "core", "get-raw-config" ]
						map = {
							foo = "bar"
							elem = "420"
						}
						object {
							nested_string = "another random string"
							nested_bool = true
							nested_int = 2147364647
						}
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					// Primitive types
					assertGetRawConfigResults(t, "root_string", "a random string", true, cty.String),
					assertGetRawConfigResults(t, "root_bool", false, true, cty.Bool),
					assertGetRawConfigResults(t, "root_int", 209, true, cty.Number),
					assertGetRawConfigResults(t, "key_not_set", nil, false, cty.Number),
					// List
					assertGetRawConfigResults(t, "list.0", "terraform-test", true, cty.String),
					assertGetRawConfigResults(t, "list.1", "core", true, cty.String),
					assertGetRawConfigResults(t, "list.2", "get-raw-config", true, cty.String),
					assertGetRawConfigResults(t, "list.3", nil, false, cty.String),
					// Map
					assertGetRawConfigResults(t, "map.foo", "bar", true, cty.String),
					assertGetRawConfigResults(t, "map.elem", "420", true, cty.String),
					assertGetRawConfigResults(t, "map.not_in_map", nil, false, cty.String),
					// Object
					assertGetRawConfigResults(t, "object.0.nested_string", "another random string", true, cty.String),
					assertGetRawConfigResults(t, "object.0.nested_bool", true, true, cty.Bool),
					assertGetRawConfigResults(t, "object.0.nested_int", 2147483647, true, cty.Number),
					assertGetRawConfigResults(t, "object.0.not_in_object", nil, false, cty.String),
					assertGetRawConfigResults(t, "object.#.nested_bool", true, true, cty.Bool),
					// Default values
					assertGetRawConfigResults(t, "not_in_config_string", "default string value", false, cty.String),
					assertGetRawConfigResults(t, "not_in_config_int", -1, false, cty.Number),
					assertGetRawConfigResults(t, "not_in_config_bool", true, false, cty.Bool),
				),
			},
			{
				Config: `
					resource scaleway_mock mock {
						root_string = "string has changed"
						root_bool = true
						root_int = 2024
						list = [ "terraform-test", "core", "get-raw-config-for-key" ]
						map = {
							food = "bar"
							elem = "240"
						}
						object {
							nested_string = "this string has also changed"
							nested_bool = false
							nested_int = -2147364648
						}
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					// Primitive types
					assertGetRawConfigResults(t, "root_string", "string has changed", true, cty.String),
					assertGetRawConfigResults(t, "root_bool", true, true, cty.Bool),
					assertGetRawConfigResults(t, "root_int", 2024, true, cty.Number),
					assertGetRawConfigResults(t, "key_not_set", nil, false, cty.Number),
					// List
					assertGetRawConfigResults(t, "list.0", "terraform-test", true, cty.String),
					assertGetRawConfigResults(t, "list.1", "core", true, cty.String),
					assertGetRawConfigResults(t, "list.2", "get-raw-config-for-key", true, cty.String),
					assertGetRawConfigResults(t, "list.3", nil, false, cty.String),
					// Map
					assertGetRawConfigResults(t, "map.food", "bar", true, cty.String),
					assertGetRawConfigResults(t, "map.elem", "240", true, cty.String),
					assertGetRawConfigResults(t, "map.not_in_map", nil, false, cty.String),
					// Object
					assertGetRawConfigResults(t, "object.0.nested_string", "this string has also changed", true, cty.String),
					assertGetRawConfigResults(t, "object.0.nested_bool", false, true, cty.Bool),
					assertGetRawConfigResults(t, "object.0.nested_int", -2147483648, true, cty.Number),
					assertGetRawConfigResults(t, "object.0.not_in_object", nil, false, cty.String),
					assertGetRawConfigResults(t, "object.#.nested_bool", false, true, cty.Bool),
					// Default values
					assertGetRawConfigResults(t, "not_in_config_string", "default string value", false, cty.String),
					assertGetRawConfigResults(t, "not_in_config_int", -1, false, cty.Number),
					assertGetRawConfigResults(t, "not_in_config_bool", true, false, cty.Bool),
				),
			},
			{
				Config: `
					resource scaleway_mock mock {
						root_string = ""
						root_int = 1234567890
						list = [ "terraform-test", "get-raw-config-for-key" ]
						map = {
							foo = "bard"
							elem = "240"
						}
						object {
							nested_string = "this string has also changed"
							nested_int = 0
						}
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					// Primitive types
					assertGetRawConfigResults(t, "root_string", "", true, cty.String),
					assertGetRawConfigResults(t, "root_bool", true, false, cty.Bool),
					assertGetRawConfigResults(t, "root_int", 1234567890, true, cty.Number),
					assertGetRawConfigResults(t, "key_not_set", nil, false, cty.Number),
					// List
					assertGetRawConfigResults(t, "list.0", "terraform-test", true, cty.String),
					assertGetRawConfigResults(t, "list.1", "get-raw-config-for-key", true, cty.String),
					assertGetRawConfigResults(t, "list.2", nil, false, cty.String),
					// Map
					assertGetRawConfigResults(t, "map.foo", "bard", true, cty.String),
					assertGetRawConfigResults(t, "map.elem", "240", true, cty.String),
					assertGetRawConfigResults(t, "map.not_in_map", nil, false, cty.String),
					// Object
					assertGetRawConfigResults(t, "object.0.nested_string", "this string has also changed", true, cty.String),
					assertGetRawConfigResults(t, "object.0.nested_bool", false, true, cty.Bool),
					assertGetRawConfigResults(t, "object.0.nested_int", -2147483648, true, cty.Number),
					assertGetRawConfigResults(t, "object.0.not_in_object", nil, false, cty.String),
					assertGetRawConfigResults(t, "object.#.nested_bool", false, true, cty.Bool),
					// Default values
					assertGetRawConfigResults(t, "not_in_config_string", "default string value", false, cty.String),
					assertGetRawConfigResults(t, "not_in_config_int", -1, false, cty.Number),
					assertGetRawConfigResults(t, "not_in_config_bool", true, false, cty.Bool),
				),
			},
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
