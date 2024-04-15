package meta

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

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
	return resourceMockRead(ctx, d, i)
}
func resourceMockRead(_ context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
	_ = d.Set("root_string", d.Get("root_string"))
	_ = d.Set("root_int", d.Get("root_int"))
	_ = d.Set("root_bool", d.Get("root_bool"))
	_ = d.Set("list", d.Get("list"))
	_ = d.Set("map", d.Get("map"))
	_ = d.Set("not_in_config_string", d.Get("not_in_config_string"))
	_ = d.Set("not_in_config_int", d.Get("not_in_config_int"))
	_ = d.Set("not_in_config_bool", d.Get("not_in_config_bool"))
	return nil
}
func resourceMockUpdate(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	return resourceMockRead(ctx, d, i)
}
func resourceMockDelete(_ context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
	d.SetId("")
	return nil
}
