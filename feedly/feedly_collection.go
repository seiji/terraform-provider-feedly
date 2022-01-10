package feedly

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/seiji/feedly"
)

func feedlyCollectionSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"customizable":  {Type: schema.TypeBool, Computed: true},
		"collection_id": {Type: schema.TypeString, ForceNew: true, Required: true},
		"description":   {Type: schema.TypeString, Computed: true},
		"enterprise":    {Type: schema.TypeBool, Computed: true},
		"feed": {
			Type: schema.TypeList,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"content_type":         {Type: schema.TypeString, Computed: true},
					"description":          {Type: schema.TypeString, Computed: true},
					"estimated_engagement": {Type: schema.TypeFloat, Computed: true},
					"feed_id":              {Type: schema.TypeString, Computed: true},
					"icon_url":             {Type: schema.TypeString, Computed: true},
					"id":                   {Type: schema.TypeString, Required: true},
					"language":             {Type: schema.TypeString, Computed: true},
					"partial":              {Type: schema.TypeBool, Computed: true},
					"subscribers":          {Type: schema.TypeFloat, Computed: true},
					"title":                {Type: schema.TypeString, Computed: true, Optional: true},
					"topics":               {Type: schema.TypeString, Computed: true},
					"updated":              {Type: schema.TypeFloat, Computed: true},
					"velocity":             {Type: schema.TypeFloat, Computed: true},
					"visual_url":           {Type: schema.TypeString, Computed: true},
					"website":              {Type: schema.TypeString, Computed: true},
				},
			},
			Computed: true,
			Optional: true,
		},
		"label":     {Type: schema.TypeString, ForceNew: true, Required: true},
		"num_feeds": {Type: schema.TypeFloat, Computed: true},
	}
}

func (api *feedlyAPI) setCollection(ctx context.Context, d *schema.ResourceData, id string) diag.Diagnostics {
	var err error
	var collections feedly.Collections
	if collections, err = api.CollectionsGet(ctx, id); err != nil {
		return diag.FromErr(err)
	}
	var collection *feedly.Collection
	for i, c := range collections {
		if strings.EqualFold(c.ID, id) {
			collection = &collections[i]
			break
		}
	}
	if collection == nil {
		return diag.Errorf("collection '%s' is not found", id)
	}

	feeds := make([]map[string]interface{}, len(collection.Feeds))
	for i, v := range collection.Feeds {
		feeds[i] = map[string]interface{}{
			"content_type":         v.ContentType,
			"description":          v.Description,
			"estimated_engagement": v.EstimatedEngagement,
			"feed_id":              v.FeedID,
			"icon_url":             v.IconURL,
			"id":                   v.ID,
			"language":             v.Language,
			"partial":              v.Partial,
			"subscribers":          v.Subscribers,
			"title":                v.Title,
			"topics":               strings.Join(v.Topics, ","),
			"updated":              v.Updated,
			"velocity":             v.Velocity,
			"visual_url":           v.VisualURL,
			"website":              v.Website,
		}
	}
	for k, v := range map[string]interface{}{
		"collection_id": collection.ID,
		"customizable":  collection.Customizable,
		"description":   collection.Description,
		"enterprise":    collection.Enterprise,
		"feed":          feeds,
		"label":         collection.Label,
		"num_feeds":     collection.NumFeeds,
	} {
		if err = d.Set(k, v); err != nil {
			return diag.FromErr(err)
		}
	}
	d.SetId(collection.ID)
	return nil
}
