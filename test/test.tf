provider "cli" {
  shell = "bash"
}
resource "cli" "init" {
  triggers = {
    key = "value3"
  }
  //working_dir = "/usr"
  create_cmd = "echo \"Create$(pwd)\"; echo CreateBye 1>&2; exit 1"
  read_cmd   = "cat file.txt; exit 0"
  update_cmd = "echo \"Update\"; echo UpdateBye 1>&2; exit 3"
  delete_cmd = "exit 4"
  create_break_on_error   = false
  read_break_on_error     = false
  read_destroyed_on_error = false
  read_on_create_update   = true
  update_break_on_error   = false
  delete_break_on_error   = false
  trim_output = true
}

output "create_stdout" {
  value = "${cli.init.stdout}"
}

output "create_stderr" {
  value = "${cli.init.stderr}"
}

output "create_retval" {
  value = "${cli.init.retval}"
}
