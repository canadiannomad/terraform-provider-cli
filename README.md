terraform-provider-cli
=========================

Terraform provider plugin to replace null\_resource local\_exec.

Usage
-------------------------

```
provider "cli" {
  // Shell to run (defaults to "sh"
  shell = "bash"
}
resource "cli" "example" {
  // Any changes to this will trigger an update
  triggers = {
    key = "value"
  }

  // If specified, this command will be run on creation
  create_cmd = "echo \"Create$(pwd)\"; echo CreateBye 1>&2; exit 1"

  // If false then assume success, even if exit code is non-zero (default: true)
  create_break_on_error   = false

  // If specified, the read command is compared with the previous run and will
  // trigger an update if it differs.
  read_cmd   = "echo \"Update\"; echo UpdateBye 1>&2; exit 3"
  // If false then assume success, even if exit code is non-zero (default: true)
  read_break_on_error     = false
  // If there is an error then assume the resource has been destroyed
  // out-of-band.
  read_destroyed_on_error = false

  // If specified, this command will be run on update
  update_cmd = "echo \"Update\"; echo UpdateBye 1>&2; exit 3"
  // If false then assume success, even if exit code is non-zero (default: true)
  update_break_on_error   = false

  // If specified, this command will be run on destroy
  delete_cmd = "exit 4"
  // If false then assume success, even if exit code is non-zero (default: true)
  delete_break_on_error   = false

  // Trim the whitespace from the start and end of the output (default: false)
  trim_output = true
}

output "stdout" {
  value = "${cli.example.stdout}"
}

output "stderr" {
  value = "${cli.example.stderr}"
}

output "retval" {
  value = "${cli.example.retval}"
}
```

Installation
-------------------------

```
go get -u github.com/canadiannomad/terraform-provider-cli
```


License
-------------------------

[MIT License](./LICENSE)
