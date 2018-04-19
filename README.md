terraform-provider-cli
=========================

Terraform provider plugin to replace null\_resource local\_exec.

Usage
-------------------------

```
provider "cli" {
  shell = "bash"
}
resource "cli" "example" {
  triggers = {
    key = "value"
  }
  create_cmd = "echo \"Create$(pwd)\"; echo CreateBye 1>&2; exit 1"
  create_break_on_error   = false

  read_cmd   = "echo \"Update\"; echo UpdateBye 1>&2; exit 3"
  read_break_on_error     = false
  read_destroyed_on_error = false

  update_cmd = "echo \"Update\"; echo UpdateBye 1>&2; exit 3"
  update_break_on_error   = false

  delete_cmd = "exit 4"
  delete_break_on_error   = false

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
go get -u github.com/takebayashi/terraform-provider-dozens
```


License
-------------------------

[MIT License](./LICENSE)
