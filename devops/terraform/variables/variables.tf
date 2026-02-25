variable "example_string" {
  type        = string
  description = "A simple string variable"
  default     = "hello"
}

variable "example_number" {
  type        = number
  description = "A numeric variable (int or float)"
  default     = 42
}

variable "example_bool" {
  type        = bool
  description = "A boolean variable"
  default     = true
}

variable "example_list_string" {
  type        = list(string)
  description = "A list of strings"
  default     = ["a", "b", "c"]
}


variable "example_set_string" {
  type        = set(string)
  description = "A set of strings (unique, unordered)"
  default     = ["x", "y", "z"]
}

variable "example_map_number" {
  type        = map(number)
  description = "A map with number values"
  default = {
    port_http  = 80
    port_https = 443
  }
}

variable "example_object" {
  type = object({
    name    = string
    age     = number
    enabled = bool
  })
  description = "An object variable with typed attributes"
  default = {
    name    = "Alice"
    age     = 30
    enabled = true
  }
}

variable "example_tuple" {
  type        = tuple([string, number, bool])
  description = "A tuple with fixed-length, mixed types"
  default     = ["item", 1, false]
}

variable "example_list_of_objects" {
  type = list(object({
    name = string
    port = number
  }))
  description = "A list of objects"
  default = [
    { name = "app", port = 8080 },
    { name = "api", port = 9090 }
  ]
}

variable "example_map_of_objects" {
  type = map(object({
    instance_type = string
    count         = number
  }))
  description = "A map of objects"
  default = {
    web = { instance_type = "t3.micro", count = 2 }
    db  = { instance_type = "t3.medium", count = 1 }
  }
}

variable "example_any" {
  type        = any
  description = "A variable with no type constraint"
  default     = "can be anything"
}

variable "example_sensitive" {
  type        = string
  description = "A sensitive variable (e.g. password, token)"
  sensitive   = true
  default     = "s3cr3t"
}

variable "example_nullable" {
  type        = string
  description = "A nullable variable that can be set to null"
  nullable    = true
  default     = null
}

variable "example_no_default" {
  type        = string
  description = "A required variable with no default (must be provided)"
}

variable "example_with_validation" {
  type        = number
  description = "A variable with a validation rule"
  default     = 8080

  validation {
    condition     = var.example_with_validation >= 1 && var.example_with_validation <= 65535
    error_message = "Port must be between 1 and 65535."
  }
}


