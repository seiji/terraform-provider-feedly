package feedly

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/seiji/feedly"
)

func dataSourceFeedlyProfile() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceFeedlyProfileRead,
		Schema: map[string]*schema.Schema{
			"client":             {Type: schema.TypeString, Computed: true},
			"dropbox_connected":  {Type: schema.TypeBool, Computed: true},
			"email":              {Type: schema.TypeString, Computed: true},
			"evernote_connected": {Type: schema.TypeBool, Computed: true},
			"facebook_connected": {Type: schema.TypeBool, Computed: true},
			"family_name":        {Type: schema.TypeString, Computed: true},
			"full_name":          {Type: schema.TypeString, Computed: true},
			"gender":             {Type: schema.TypeString, Computed: true},
			"given_name":         {Type: schema.TypeString, Computed: true},
			"google":             {Type: schema.TypeString, Computed: true},
			"id":                 {Type: schema.TypeString, Computed: true},
			"locale":             {Type: schema.TypeString, Computed: true},
			"logins": {
				Type: schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"full_name":   {Type: schema.TypeString, Computed: true},
						"id":          {Type: schema.TypeString, Computed: true},
						"picture":     {Type: schema.TypeString, Computed: true},
						"provider":    {Type: schema.TypeString, Computed: true},
						"provider_id": {Type: schema.TypeString, Computed: true},
						"verified":    {Type: schema.TypeBool, Computed: true},
					},
				},
				Computed: true,
			},
			"payment_provider_id": {
				Type:     schema.TypeMap,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"payment_subscription_id": {
				Type:     schema.TypeMap,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"picture":                       {Type: schema.TypeString, Computed: true},
			"pocket_connected":              {Type: schema.TypeBool, Computed: true},
			"product":                       {Type: schema.TypeString, Computed: true},
			"product_expiration":            {Type: schema.TypeInt, Computed: true},
			"subscription_payment_provider": {Type: schema.TypeString, Computed: true},
			"subscription_status":           {Type: schema.TypeString, Computed: true},
			"twitter_connected":             {Type: schema.TypeBool, Computed: true},
			"upgrade_date":                  {Type: schema.TypeInt, Computed: true},
			"wave":                          {Type: schema.TypeString, Computed: true},
			"windows_live_connected":        {Type: schema.TypeBool, Computed: true},
			"wordpress_connected":           {Type: schema.TypeBool, Computed: true},
		},
	}
}

func dataSourceFeedlyProfileRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	api := m.(*feedlyAPI)

	var profile *feedly.Profile
	var err error
	if profile, err = api.ProfileGet(ctx); err != nil {
		return diag.FromErr(err)
	}

	logins := make([]map[string]interface{}, len(profile.Logins))
	for i, v := range profile.Logins {
		logins[i] = map[string]interface{}{
			"full_name":   v.FullName,
			"id":          v.ID,
			"picture":     v.Picture,
			"provider":    v.Provider,
			"provider_id": v.ProviderID,
			"verified":    v.Verified,
		}
	}
	var paymentProviderID map[string]string
	if profile.PaymentProviderID != nil {
		paymentProviderID = map[string]string{
			"paypal": profile.PaymentProviderID.Paypal,
		}
	}
	var paymentSubscriptionID map[string]string
	if profile.PaymentSubscriptionID != nil {
		paymentSubscriptionID = map[string]string{
			"paypal": profile.PaymentSubscriptionID.Paypal,
		}
	}
	for k, v := range map[string]interface{}{
		"client":                        profile.Client,
		"dropbox_connected":             profile.DropboxConnected,
		"email":                         profile.Email,
		"evernote_connected":            profile.EvernoteConnected,
		"facebook_connected":            profile.FacebookConnected,
		"family_name":                   profile.FamilyName,
		"full_name":                     profile.FullName,
		"gender":                        profile.Gender,
		"given_name":                    profile.GivenName,
		"google":                        profile.Google,
		"id":                            profile.ID,
		"locale":                        profile.Locale,
		"logins":                        logins,
		"payment_provider_id":           paymentProviderID,
		"payment_subscription_id":       paymentSubscriptionID,
		"picture":                       profile.Picture,
		"pocket_connected":              profile.PocketConnected,
		"product":                       profile.Product,
		"product_expiration":            profile.ProductExpiration,
		"subscription_payment_provider": profile.SubscriptionPaymentProvider,
		"subscription_status":           profile.SubscriptionStatus,
		"twitter_connected":             profile.TwitterConnected,
		"upgrade_date":                  profile.UpgradeDate,
		"wave":                          profile.Wave,
		"windows_live_connected":        profile.WindowsLiveConnected,
		"wordpress_connected":           profile.WordPressConnected,
	} {
		if err = d.Set(k, v); err != nil {
			return diag.FromErr(err)
		}
	}
	d.SetId(profile.ID)

	return nil
}
