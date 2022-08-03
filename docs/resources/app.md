---
page_title: "appetize_app Resource - terraform-provider-appetize"
subcategory: ""
description: |-
The App resource allows you to upload and create an App
---

# appetize_app (Resource)

## Example Usage

```terraform
resource "appetize_app" "app" {
  url       = "https://ftp.com/my/android/app"
  platform  = "android"
  
  note = "this is a note"
}
```

## Schema

### Required

- `platform` (String)
- `url` (String)

### Optional

- `button_text` (String)
- `disable_home` (Boolean)
- `disabled` (Boolean)
- `file_type` (String)
- `launch_url` (String)
- `note` (String)
- `post_session_button_text` (String)
- `use_last_frame` (Boolean)

### Read-Only

- `id` (String) The ID of this resource.
- `name` (String)
- `private_key` (String)
- `public_key` (String)
