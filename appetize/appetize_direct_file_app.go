package appetize

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAppetizeDirectFileApp() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAppetizeDirectFileAppCreate,
		ReadContext:   resourceAppetizeAppRead,
		DeleteContext: resourceAppetizeAppDelete,

		Schema: func() map[string]*schema.Schema {
			s := resourceAppetizeApp().Schema

			delete(s, "url")
			s["file_path"] = &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			}

			// Everything is ForceNew
			for _, v := range s {
				v.ForceNew = true
			}

			return s
		}(),
	}
}

func resourceAppetizeDirectFileAppCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	appetizer := m.(*Appetize)
	err := appetizer.DeleteApp(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}
