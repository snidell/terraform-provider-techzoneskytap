---
page_title: "skytap_project Resource - terraform-provider-skytap"
subcategory: ""
description: |-
  Provides a Skytap Project resource.
---

# skytap_project (Resource)

Provides a Skytap Project resource. Projects are an access permissions model used to share environments, 
templates, and assets with other users.

## Example Usage

```hcl
# Create a new project
resource "skytap_project" "project" {
  name = "Terraform Example"
  summary = "Skytap terraform provider example project."
  show_project_members = false
  auto_add_role_name = "participant"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- **name** (String) User-defined project name

### Optional

- **auto_add_role_name** (String) If this field is set to `viewer`, `participant`, `editor`, or `manager`, new users added to your Skytap account are automatically added to this project with the specified project role. Existing users aren’t affected by this setting. For additional details, see [Automatically adding new users to a project](https://help.skytap.com/csh-project-automatic-role.html)
- **environment_ids** (Set of String) A list of environments to add to the project
- **id** (String) The ID of this resource.
- **show_project_members** (Boolean) Whether project members can view a list of other project members
- **summary** (String) User-defined description of the project
- **timeouts** (Block, Optional) (see [below for nested schema](#nestedblock--timeouts))

<a id="nestedblock--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- **create** (String)
- **delete** (String)
- **update** (String)
