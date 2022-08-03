---
page_title: "Appetize Provider"
subcategory: ""
description: "Terraform provider for interacting with [Appetize](https://appetize.io/) API"
---

# Appetize Provider

Terraform provider for interacting with [Appetize](https://appetize.io/) API.

## Example Usage

Do not keep your authentication api token in HCL for production environments, use Terraform environment variables.

```terraform
provider "appetize" {
  api_token = "tokentokentoken" 
}
```
