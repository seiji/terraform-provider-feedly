package feedly

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/seiji/feedly"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"access_token": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("FEEDLY_ACCESS_TOKEN", nil),
				Description: "The Feedly Access Token",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"feedly_collection": resourceFeedlyCollection(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"feedly_collection": dataSourceFeedlyCollection(),
			"feedly_profile":    dataSourceFeedlyProfile(),
		},
		ConfigureFunc: func(d *schema.ResourceData) (interface{}, error) {
			// accessToken := d.Get("access_token").(string)
			return &feedlyAPI{feedly.NewAPI(nil)}, nil
		},
		TerraformVersion: "",
	}
}

type feedlyAPI struct {
	feedly.API
}
