package appetize

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"appetize_app":             resourceAppetizeApp(),
			"appetize_direct_file_app": resourceAppetizeDirectFileApp(),
		},
	}
}
