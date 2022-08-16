variable "project_id" {
  default = "test-project-123"
}

variable "service_account" {
  type    = string
  default = ""
}

variable "instances" {
  type = map(object({
    machine_type = string
    zone         = string
    description  = string
  }))
}

resource "google_compute_instance" "vm" {
  for_each     = var.instances
  project      = var.project_id
  name         = each.key
  description  = each.value.description
  machine_type = each.value.machine_type
  zone         = each.value.zone
  hostname     = "${each.key}-${each.value.zone}"

  dynamic "service_account" {
    // This is will use the service_account block of code if the var.service_account has been define.
    // You could use this as an argument in a terraform module to only use this block of code if the var.service_account is defined. 
    
    for_each = var.service_account != "" ? [1] : []
    content {
      email  = var.service_account
      scopes = var.scopes
    }
  }
}

output "vms" {
  value = {
    for k, v in google_compute_instance.vm : k => tomap({ v.name = { hostname = v.hostname } })
  }
}