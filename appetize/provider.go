package appetize

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_token": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("APPETIZE_API_TOKEN", nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"appetize_app":             resourceAppetizeApp(),
			"appetize_direct_file_app": resourceAppetizeDirectFileApp(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	apiToken := d.Get("api_token").(string)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	appetizer := NewAppetizer(apiToken)
	return appetizer, diags
}
