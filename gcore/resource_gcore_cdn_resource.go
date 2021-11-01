package gcore

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	gcdn "github.com/skyeng/gcorelabscdn-go/gcore"
	"github.com/skyeng/gcorelabscdn-go/resources"
)

func resourceCDNResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"cname": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "A CNAME that will be used to deliver content though a CDN",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "An optional comment that describes this resource",
			},
			"origin_group": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ExactlyOneOf: []string{
					"origin_group",
					"origin",
				},
				Description: "ID of the Origins Group. Use one of your Origins Group or create a new one. You can use either 'origin' parameter or 'originGroup' in the resource definition.",
			},
			"origin": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ExactlyOneOf: []string{
					"origin_group",
					"origin",
				},
				Description: "A domain name or IP of your origin source. Specify a port if custom. You can use either 'origin' parameter or 'originGroup' in the resource definition.",
			},
			"origin_protocol": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "This option defines the protocol that will be used by CDN servers to request content from an origin source. If not specified, we will use HTTP to connect to an origin server. Possible values are: HTTPS, HTTP, MATCH.",
			},
			"secondary_hostnames": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of additional CNAMEs.",
			},
			"ssl_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Use HTTPS protocol for content delivery.",
			},
			"ssl_data": {
				Type:         schema.TypeInt,
				Optional:     true,
				RequiredWith: []string{"ssl_enabled"},
				Description:  "Specify the SSL Certificate ID which should be used for the CDN Resource.",
			},
			"options": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Computed:    true,
				Description: "Each option in CDN resource settings. Each option added to CDN resource settings should have the following mandatory request fields: enabled, value.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"edge_cache_settings": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Computed:    true,
							Description: "The cache expiration time for CDN servers.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:     schema.TypeBool,
										Required: true,
									},
									"value": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Caching time for a response with codes 200, 206, 301, 302. Responses with codes 4xx, 5xx will not be cached. Use '0s' disable to caching. Use custom_values field to specify a custom caching time for a response with specific codes.",
									},
									"default": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Content will be cached according to origin cache settings. The value applies for a response with codes 200, 201, 204, 206, 301, 302, 303, 304, 307, 308 if an origin server does not have caching HTTP headers. Responses with other codes will not be cached.",
									},
									"custom_values": {
										Type:        schema.TypeMap,
										Optional:    true,
										Elem:        schema.TypeString,
										Description: "Caching time for a response with specific codes. These settings have a higher priority than the value field. Response code ('304', '404' for example). Use 'any' to specify caching time for all response codes. Caching time in seconds ('0s', '600s' for example). Use '0s' to disable caching for a specific response code.",
									},
								},
							},
						},
						"browser_cache_settings": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Computed:    true,
							Description: "The cache expiration time for customers' browsers in seconds.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:     schema.TypeBool,
										Required: true,
									},
									"value": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The value applies for a response with codes 200, 201, 204, 206, 301, 302, 303, 304, 307, 308. Responses with other codes will not be cached. Use '0s' to disable caching.",
									},
								},
							},
						},
						"host_header": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Specify the Host header that CDN servers use when request content from an origin server. Your server must be able to process requests with the chosen header. If the option is in NULL state Host Header value is taken from the CNAME field.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:     schema.TypeBool,
										Required: true,
									},
									"value": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
						"webp": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "The option allows automatically convert JPG and PNG images into WebP format",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:     schema.TypeBool,
										Required: true,
									},
									"jpg_quality": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"png_quality": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"png_lossless": {
										Type:     schema.TypeBool,
										Optional: true,
									},
								},
							},
						},
						"rewrite": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Defines nginx rewrite directive.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:     schema.TypeBool,
										Required: true,
									},
									"body": {
										Type:     schema.TypeString,
										Required: true,
									},
									"flag": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
						//"redirect_http_to_https": {
						//	Type:        schema.TypeList,
						//	MaxItems:    1,
						//	Optional:    true,
						//	Description: "When enabled redirects HTTP requests to HTTPS",
						//	Elem: &schema.Resource{
						//		Schema: map[string]*schema.Schema{
						//			"enabled": {
						//				Type:     schema.TypeBool,
						//				Required: true,
						//			},
						//			"value": {
						//				Type:     schema.TypeBool,
						//				Required: true,
						//			},
						//		},
						//	},
						//},
						//"cors": {
						//	Type:        schema.TypeList,
						//	MaxItems:    1,
						//	Optional:    true,
						//	Description: "The option adds the Access-Control-Allow-Origin header to responses from CDN servers",
						//	Elem: &schema.Resource{
						//		Schema: map[string]*schema.Schema{
						//			"enabled": {
						//				Type:     schema.TypeBool,
						//				Required: true,
						//			},
						//			//"value": {
						//			//	Type:     schema.TypeList,
						//			//	Required: true,
						//			//	Elem: &schema.Schema{
						//			//		Type: schema.TypeString,
						//			//	},
						//			//},
						//		},
						//	},
						//},
					},
				},
			},
			"active": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "The setting allows to enable or disable a CDN Resource",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of a CDN resource content availability. Possible values are: Active, Suspended, Processed.",
			},
		},
		CreateContext: resourceCDNResourceCreate,
		ReadContext:   resourceCDNResourceRead,
		UpdateContext: resourceCDNResourceUpdate,
		DeleteContext: resourceCDNRresourceDelete,
		Description:   "Represent cdn resource",
	}
}

func resourceCDNResourceCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start CDN Resource creating")
	config := m.(*Config)
	client := config.CDNClient

	var req resources.CreateRequest
	req.Cname = d.Get("cname").(string)
	req.Description = d.Get("description").(string)
	req.Origin = d.Get("origin").(string)
	req.OriginGroup = d.Get("origin_group").(int)

	for _, hostname := range d.Get("secondary_hostnames").(*schema.Set).List() {
		req.SecondaryHostnames = append(req.SecondaryHostnames, hostname.(string))
	}

	result, err := client.Resources().Create(ctx, &req)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d", result.ID))
	resourceCDNResourceRead(ctx, d, m)

	log.Printf("[DEBUG] Finish CDN Resource creating (id=%d)\n", result.ID)
	return nil
}

func resourceCDNResourceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	resourceID := d.Id()
	log.Printf("[DEBUG] Start CDN Resource reading (id=%s)\n", resourceID)
	config := m.(*Config)
	client := config.CDNClient

	id, err := strconv.ParseInt(resourceID, 10, 64)
	if err != nil {
		return diag.FromErr(err)
	}

	result, err := client.Resources().Get(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("cname", result.Cname)
	d.Set("description", result.Description)
	d.Set("origin_group", result.OriginGroup)
	d.Set("origin_protocol", result.OriginProtocol)
	d.Set("secondary_hostnames", result.SecondaryHostnames)
	d.Set("ssl_enabled", result.SSlEnabled)
	d.Set("ssl_data", result.SSLData)
	d.Set("status", result.Status)
	d.Set("active", result.Active)
	if err := d.Set("options", optionsToList(result.Options)); err != nil {
		return diag.FromErr(err)
	}

	log.Println("[DEBUG] Finish CDN Resource reading")
	return nil
}

func resourceCDNResourceUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	resourceID := d.Id()
	log.Printf("[DEBUG] Start CDN Resource updating (id=%s)\n", resourceID)
	config := m.(*Config)
	client := config.CDNClient

	id, err := strconv.ParseInt(resourceID, 10, 64)
	if err != nil {
		return diag.FromErr(err)
	}

	var req resources.UpdateRequest
	req.Active = d.Get("active").(bool)
	req.Description = d.Get("description").(string)
	req.SSlEnabled = d.Get("ssl_enabled").(bool)
	req.SSLData = d.Get("ssl_data").(int)
	req.OriginGroup = d.Get("origin_group").(int)
	req.OriginProtocol = resources.Protocol(d.Get("origin_protocol").(string))
	req.Options = listToOptions(d.Get("options").([]interface{}))
	for _, hostname := range d.Get("secondary_hostnames").(*schema.Set).List() {
		req.SecondaryHostnames = append(req.SecondaryHostnames, hostname.(string))
	}

	if _, err := client.Resources().Update(ctx, id, &req); err != nil {
		return diag.FromErr(err)
	}

	log.Println("[DEBUG] Finish CDN Resource updating")
	return resourceCDNResourceRead(ctx, d, m)
}

func resourceCDNRresourceDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	resourceID := d.Id()
	log.Printf("[DEBUG] Start CDN Resource deleting (id=%s)\n", resourceID)
	config := m.(*Config)
	client := config.CDNClient

	id, err := strconv.ParseInt(resourceID, 10, 64)
	if err != nil {
		return diag.FromErr(err)
	}

	var req resources.UpdateRequest
	req.Active = false
	req.OriginGroup = d.Get("origin_group").(int)

	if _, err := client.Resources().Update(ctx, id, &req); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	log.Println("[DEBUG] Finish CDN Resource deleting")
	return nil
}

func listToOptions(l []interface{}) *gcdn.Options {
	if len(l) == 0 {
		return nil
	}

	var opts gcdn.Options
	fields := l[0].(map[string]interface{})
	if opt, ok := getOptByName(fields, "edge_cache_settings"); ok {
		rawCustomVals := opt["custom_values"].(map[string]interface{})
		customVals := make(map[string]string, len(rawCustomVals))
		for key, value := range rawCustomVals {
			customVals[key] = value.(string)
		}

		opts.EdgeCacheSettings = &gcdn.EdgeCacheSettings{
			Enabled:      opt["enabled"].(bool),
			Value:        opt["value"].(string),
			CustomValues: customVals,
			Default:      opt["default"].(string),
		}
	}
	if opt, ok := getOptByName(fields, "browser_cache_settings"); ok {
		opts.BrowserCacheSettings = &gcdn.BrowserCacheSettings{
			Enabled: opt["enabled"].(bool),
			Value:   opt["value"].(string),
		}
	}
	if opt, ok := getOptByName(fields, "host_header"); ok {
		opts.HostHeader = &gcdn.HostHeader{
			Enabled: opt["enabled"].(bool),
			Value:   opt["value"].(string),
		}
	}
	if opt, ok := getOptByName(fields, "webp"); ok {
		opts.Webp = &gcdn.Webp{
			Enabled:     opt["enabled"].(bool),
			JPGQuality:  opt["jpg_quality"].(int),
			PNGQuality:  opt["png_quality"].(int),
			PNGLossless: opt["png_lossless"].(bool),
		}
	}
	if opt, ok := getOptByName(fields, "rewrite"); ok {
		opts.Rewrite = &gcdn.Rewrite{
			Enabled: opt["enabled"].(bool),
			Body:    opt["body"].(string),
			Flag:    opt["flag"].(string),
		}
	}
	//if opt, ok := getOptByName(fields, "redirect_http_to_https"); ok {
	//	opts.RedirectHttpToHttps = &gcdn.RedirectHttpToHttps{
	//		Enabled: opt["enabled"].(bool),
	//		Value:   opt["value"].(bool),
	//	}
	//}
	//if opt, ok := getOptByName(fields, "cors"); ok {
	//	opts.Cors = &gcdn.Cors{
	//		Enabled: opt["enabled"].(bool),
	//		Value:   opt["value"].([]string),
	//	}
	//}

	return &opts
}

func getOptByName(fields map[string]interface{}, name string) (map[string]interface{}, bool) {
	if _, ok := fields[name]; !ok {
		return nil, false
	}

	container, ok := fields[name].([]interface{})
	if !ok {
		return nil, false
	}

	if len(container) == 0 {
		return nil, false
	}

	opt, ok := container[0].(map[string]interface{})
	if !ok {
		return nil, false
	}

	return opt, true
}

func optionsToList(options *gcdn.Options) []interface{} {
	result := make(map[string][]interface{})
	if options.EdgeCacheSettings != nil {
		m := structToMap(options.EdgeCacheSettings)
		result["edge_cache_settings"] = []interface{}{m}
	}
	if options.BrowserCacheSettings != nil {
		m := structToMap(options.BrowserCacheSettings)
		result["browser_cache_settings"] = []interface{}{m}
	}
	if options.HostHeader != nil {
		m := structToMap(options.HostHeader)
		result["host_header"] = []interface{}{m}
	}
	if options.Webp != nil {
		m := structToMap(options.Webp)
		result["webp"] = []interface{}{m}
	}
	if options.Rewrite != nil {
		m := structToMap(options.Rewrite)
		result["rewrite"] = []interface{}{m}
	}
	//if options.RedirectHttpToHttps != nil {
	//	m := structToMap(options.RedirectHttpToHttps)
	//	result["redirect_http_to_https"] = []interface{}{m}
	//}
	//if options.Cors != nil {
	//	m := structToMap(options.Cors)
	//	result["cors"] = []interface{}{m}
	//}
	return []interface{}{result}
}

func structToMap(item interface{}) map[string]interface{} {

	res := map[string]interface{}{}
	if item == nil {
		return res
	}
	v := reflect.TypeOf(item)
	reflectValue := reflect.ValueOf(item)
	reflectValue = reflect.Indirect(reflectValue)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	for i := 0; i < v.NumField(); i++ {
		tag := v.Field(i).Tag.Get("json")
		field := reflectValue.Field(i).Interface()
		if tag != "" && tag != "-" {
			if v.Field(i).Type.Kind() == reflect.Struct {
				res[tag] = structToMap(field)
			} else {
				res[tag] = field
			}
		}
	}
	return res
}
