// computerextensionattributes_resource.go
package computerextensionattributes

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

// ResourceJamfProComputerExtensionAttributes defines the schema and CRUD operations (Create, Read, Update, Delete)
// for managing Jamf Pro Computer Extension Attributes in Terraform.
type ResourceJamfProComputerExtensionAttributes struct{}

func (r *ResourceJamfProComputerExtensionAttributes) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_jamf_pro_computer_extension_attributes"
}

// ResourceJamfProComputerExtensionAttributes defines the schema and CRUD operations (Create, Read, Update, Delete)
// for managing Jamf Pro Computer Extension Attributes in Terraform.
func (r *ResourceJamfProComputerExtensionAttributes) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Resource for managing Jamf Pro Computer Extension Attributes in Terraform.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The unique identifier of the computer extension attribute.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The unique name of the Jamf Pro computer extension attribute.",
			},
			"enabled": schema.BoolAttribute{
				Required:            true,
				MarkdownDescription: "Indicates if the computer extension attribute is enabled.",
			},
			"description": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Description of the computer extension attribute.",
			},
			"data_type": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Data type of the computer extension attribute. Can be String / Integer / Date (YYYY-MM-DD hh:mm:ss)",
				Validators: []tfsdk.AttributeValidator{
					dataTypeValidator{},
				},
			},
			"input_type": schema.ListNestedAttributes(map[string]schema.Attribute{
				"type": schema.StringAttribute{
					Required:            true,
					MarkdownDescription: "Type of input for the computer extension attribute.",
					Validators: []tfsdk.AttributeValidator{
						stringInSliceValidator([]string{"script", "Text Field", "LDAP Mapping", "Pop-up Menu"}),
					},
				},
				"platform": schema.StringAttribute{
					Optional:            true,
					Computed:            true,
					MarkdownDescription: "Platform type for the computer extension attribute.",
					Validators: []tfsdk.AttributeValidator{
						stringInSliceValidator([]string{"Mac", "Windows"}),
					},
				},
				"script": schema.StringAttribute{
					Optional:            true,
					Computed:            true,
					MarkdownDescription: "Script associated with the computer extension attribute.",
				},
				"choices": schema.ListNestedAttributes(map[string]schema.Attribute{
					"value": schema.StringAttribute{
						Optional: true,
					},
				}, schema.ListNestedAttributesOptions{}),
			}, schema.ListNestedAttributesOptions{}),
			"inventory_display": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Display details for inventory for the computer extension attribute.",
				Validators: []tfsdk.AttributeValidator{
					stringInSliceValidator([]string{"General", "Hardware", "Operating System", "User and Location", "Purchasing", "Extension Attributes"}),
				},
			},
			"recon_display": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Display details for recon for the computer extension attribute.",
			},
		},
	}
}

func (r *ResourceJamfProComputerExtensionAttributes) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Implement the creation logic here
}

func (r *ResourceJamfProComputerExtensionAttributes) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Implement the read logic here
}

func (r *ResourceJamfProComputerExtensionAttributes) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Implement the update logic here
}

func (r *ResourceJamfProComputerExtensionAttributes) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Implement the deletion logic here
}
