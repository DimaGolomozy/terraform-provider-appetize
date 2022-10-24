# terraform-appetize-provider
test
please see basic usage here [appetize-docs](https://appetize.io/docs)

### resources:
`appetize_app`  
`appetize_direct_file_app`  

### Argument Reference
```
url         - (Required) A publicly accessible link to your .zip, .tar.gz, or .apk file
              Only for appetize_app resource
file_path   - (Required) Specify a file location on the local filesystem
              Only for appetize_direct_file_app resource.
platform    - (Required) ios or android
disabled    - (Optional) Disables streaming for this app
note        - (Optional) A note for your own purposes
button_text - (Optional) Customize the message prompting the user to start the session
post_session_button_text - (Optional) Customize the message prompting the user to restart the session
```

### Attributes Reference

In addition to all arguments above, the following attributes are exported:

```
name        - Name of the app
public_key  - Public key of the app 
private_key - Private key of the app 
```

## Example Usage

Do not keep your authentication api token in HCL for production environments, use Terraform environment variables.

```terraform
provider "appetize" {
  api_token = "tokentokentoken" 
}

resource "appetize_app" "app" {
  url       = "https://ftp.com/my/android/app"
  platform  = "android"
  
  note = "this is a note"
}
```
