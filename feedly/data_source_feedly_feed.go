package feedly

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/seiji/feedly"
)

func dataSourceFeedlyFeed() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceFeedlyFeedRead,
		Schema: map[string]*schema.Schema{
			"content_type":         {Type: schema.TypeString, Computed: true},
			"description":          {Type: schema.TypeString, Computed: true},
			"estimated_engagement": {Type: schema.TypeFloat, Computed: true},
			"feed_id":              {Type: schema.TypeString, Required: true},
			"icon_url":             {Type: schema.TypeString, Computed: true},
			"language":             {Type: schema.TypeString, Computed: true},
			"partial":              {Type: schema.TypeBool, Computed: true},
			"subscribers":          {Type: schema.TypeFloat, Computed: true},
			"title":                {Type: schema.TypeString, Computed: true},
			"topics":               {Type: schema.TypeString, Computed: true},
			"updated":              {Type: schema.TypeFloat, Computed: true},
			"velocity":             {Type: schema.TypeFloat, Computed: true},
			"visual_url":           {Type: schema.TypeString, Computed: true},
			"website":              {Type: schema.TypeString, Computed: true},
		},
	}
}

func dataSourceFeedlyFeedRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	api := m.(*feedlyAPI)
	feedID := d.Get("feed_id").(string)

	var feed *feedly.Feed
	var err error
	if feed, err = api.FeedsGet(ctx, feedID); err != nil {
		return diag.FromErr(err)
	}

	for k, v := range map[string]interface{}{
		"content_type":         feed.ContentType,
		"description":          feed.Description,
		"estimated_engagement": feed.EstimatedEngagement,
		"feed_id":              feed.FeedID,
		"icon_url":             feed.IconURL,
		"language":             feed.Language,
		"partial":              feed.Partial,
		"subscribers":          feed.Subscribers,
		"title":                feed.Title,
		"topics":               strings.Join(feed.Topics, ","),
		"updated":              feed.Updated,
		"velocity":             feed.Velocity,
		"visual_url":           feed.VisualURL,
		"website":              feed.Website,
	} {
		if err = d.Set(k, v); err != nil {
			return diag.FromErr(err)
		}
	}
	d.SetId(feed.ID)

	return nil
}
