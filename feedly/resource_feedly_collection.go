package feedly

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/seiji/feedly"
)

func resourceFeedlyCollection() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFeedlyCollectionCreate,
		ReadContext:   resourceFeedlyCollectionRead,
		UpdateContext: resourceFeedlyCollectionUpdate,
		DeleteContext: resourceFeedlyCollectionDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: feedlyCollectionSchema(),
	}
}

func resourceFeedlyCollectionCreate(ctx context.Context, d *schema.ResourceData, m interface{}) (dd diag.Diagnostics) {
	api := m.(*feedlyAPI)
	id := d.Get("collection_id").(string)
	label := d.Get("label").(string)
	description := d.Get("description").(string)
	feed := d.Get("feed").([]interface{})

	feeds := make([]feedly.CollectionFeedCreate, len(feed))
	for i, v := range feed {
		m := v.(map[string]interface{})
		feeds[i] = feedly.CollectionFeedCreate{
			ID:    m["id"].(string),
			Title: m["title"].(string),
		}
	}

	var collections feedly.Collections
	var err error
	if collections, err = api.CollectionsCreate(ctx, &feedly.CollectionCreate{
		Description: description,
		Feeds:       feeds,
		ID:          id,
		Label:       label,
	}); err != nil {
		return diag.FromErr(err)
	}
	var collection *feedly.Collection
	for i, c := range collections {
		if strings.EqualFold(c.ID, id) {
			collection = &collections[i]
			break
		}
	}
	d.SetId(collection.ID)
	return resourceFeedlyCollectionRead(ctx, d, m)
}

func resourceFeedlyCollectionRead(ctx context.Context, d *schema.ResourceData, m interface{}) (dd diag.Diagnostics) {
	api := m.(*feedlyAPI)
	id := d.Id()
	return api.setCollection(ctx, d, id)
}

func resourceFeedlyCollectionUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	api := m.(*feedlyAPI)
	id := d.Id()

	if d.HasChange("feed") {
		var err error
		o, n := d.GetChange("feed")
		oldFeeds, newFeeds := o.([]interface{}), n.([]interface{})
		if len(oldFeeds) > 0 {
			oldParams := make(feedly.CollectionFeedDeletes, len(oldFeeds))
			for i, v := range oldFeeds {
				feed := v.(map[string]interface{})
				oldParams[i] = feedly.CollectionFeedDelete{ID: feed["id"].(string)}
			}
			if err = api.CollectionsFeedsMDelete(ctx, id, oldParams); err != nil {
				return diag.FromErr(err)
			}
		}
		if len(newFeeds) > 0 {
			newParams := make(feedly.CollectionFeedCreates, len(newFeeds))
			for i, v := range newFeeds {
				feed := v.(map[string]interface{})
				newParams[i] = feedly.CollectionFeedCreate{
					ID:    feed["id"].(string),
					Title: feed["title"].(string),
				}
			}
			if _, err = api.CollectionsFeedsMPut(ctx, id, newParams); err != nil {
				return diag.FromErr(err)
			}
		}
	}
	return resourceFeedlyCollectionRead(ctx, d, m)
}

func resourceFeedlyCollectionDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	label := d.Get("label").(string)
	id := d.Get("collection_id").(string)
	api := m.(*feedlyAPI)

	var err error
	if _, err = api.CollectionsCreate(ctx, &feedly.CollectionCreate{
		Description: "",
		Feeds:       []feedly.CollectionFeedCreate{},
		ID:          id,
		Label:       label,
	}); err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return nil
}
