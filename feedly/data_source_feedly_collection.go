package feedly

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceFeedlyCollection() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceFeedlyCollectionRead,
		Schema:      feedlyCollectionSchema(),
	}
}

func dataSourceFeedlyCollectionRead(ctx context.Context, d *schema.ResourceData, m interface{}) (dd diag.Diagnostics) {
	api := m.(*feedlyAPI)
	id := d.Get("collection_id").(string)
	return api.setCollection(ctx, d, id)
}
