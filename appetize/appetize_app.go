package appetize

import (
	"context"

	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAppetizeApp() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAppetizeAppCreate,
		ReadContext:   resourceAppetizeAppRead,
		UpdateContext: resourceAppetizeAppUpdate,
		DeleteContext: resourceAppetizeAppDelete,

		Schema: map[string]*schema.Schema{
			"url": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.IsURLWithHTTPorHTTPS,
			},
			"platform": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"ios", "android"}, false),
			},
			"file_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			//"timeout": &schema.Schema{
			//	Type:     schema.TypeInt,
			//	Optional: true,
			//},
			"disable_home": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"disabled": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"use_last_frame": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"launch_url": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"note": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"button_text": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"post_session_button_text": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"public_key": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"private_key": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAppetizeAppCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	appetizer := m.(*Appetize)
	app, err := appetizer.CreateApp(NewAppOptions(d))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(app.PublicKey)
	return append(diags, resourceAppetizeAppRead(ctx, d, m)...)
}

func resourceAppetizeAppRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	appetizer := m.(*Appetize)
	app, err := appetizer.GetApp(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	if app == nil {
		return diag.FromErr(fmt.Errorf("failed getting app with id [%v]", d.Id()))
	}

	d.Set("public_key", app.PublicKey)
	d.Set("private_key", app.PrivateKey)
	d.Set("platform", app.Platform)
	d.Set("disabled", app.Disabled)

	if app.Name != nil {
		d.Set("name", app.Name)
	}
	if app.ButtonText != nil {
		d.Set("button_text", app.ButtonText)
	}
	if app.PostSessionButtonText != nil {
		d.Set("post_session_button_text", app.PostSessionButtonText)
	}
	if app.Note != nil {
		d.Set("note", app.Note)
	}
	if app.FileType != nil {
		d.Set("file_type", app.FileType)
	}
	if app.UseLastFrame != nil {
		d.Set("use_last_frame", app.UseLastFrame)
	}
	if app.DisabledHome != nil {
		d.Set("disabled_home", app.DisabledHome)
	}
	if app.LaunchUrl != nil {
		d.Set("launch_url", app.LaunchUrl)
	}
	//if app.Timeout != nil {
	//	d.Set("timeout", app.Timeout)
	//}

	return diags
}

func resourceAppetizeAppUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	isChanged := d.HasChanges("url", "platform", "disabled", "note", "button_text", "post_session_button_text")
	if isChanged {
		appetizer := m.(*Appetize)
		_, err := appetizer.UpdateApp(d.Id(), NewAppOptions(d))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceAppetizeAppRead(ctx, d, m)
}

func resourceAppetizeAppDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	appetizer := m.(*Appetize)
	err := appetizer.DeleteApp(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}
