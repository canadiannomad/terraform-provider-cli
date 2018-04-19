package main

import (
	"github.com/hashicorp/terraform/helper/schema"
    "math/rand"
    "log"
    "fmt"
    "bytes"
    "strings"
    "syscall"
    "os/exec"
)
// https://godoc.org/github.com/hashicorp/terraform/helper/schema#Resource


func resourceServer() *schema.Resource {
	return &schema.Resource{
		Create: createAction,
		Read:   readAction,
		Update: updateAction,
		Delete: deleteAction,
		Schema: map[string]*schema.Schema{
            "triggers": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
			},

            // Directory to run the commands in
			"working_dir": &schema.Schema{
				Type:     schema.TypeString,
                Optional: true,
				Required: false,
			},
            // This will be run on creation
			"create_cmd": &schema.Schema{
				Type:     schema.TypeString,
                Optional: true,
				Required: false,
				ForceNew: false,
			},
			"create_break_on_error": &schema.Schema{
				Type:     schema.TypeBool,
                Optional: true,
				Required: false,
                Default:  true,
			},

			// The Read callback is used to sync the local state with the actual state (upstream).
			// This is called at various points by Terraform and should be a read-only operation.

			"read_cmd": &schema.Schema{
				Type:     schema.TypeString,
                Optional: true,
				Required: false,
				ForceNew: false,
			},

			// If exit code is non-zero, we will assume the resource no longer exists
			//   (maybe it was destroyed out of band).
            // If the results of this differ from the previous run then update will be triggered
			"read_destroyed_on_error": &schema.Schema{
				Type:     schema.TypeBool,
                Optional: true,
				Required: false,
                Default:  true,
			},

			"read_break_on_error": &schema.Schema{
				Type:     schema.TypeBool,
                Optional: true,
				Required: false,
                Default:  true,
			},

            // If this is set to true then the output of create and update will be the result of
            // the execution of the read command
			"read_on_create_update": &schema.Schema{
				Type:     schema.TypeBool,
                Optional: true,
				Required: false,
                Default:  false,
			},

			// This will be run if we are updating the state.
			"update_cmd": &schema.Schema{
				Type:     schema.TypeString,
                Optional: true,
				Required: false,
				ForceNew: false,
			},
			"update_break_on_error": &schema.Schema{
				Type:     schema.TypeBool,
                Optional: true,
				Required: false,
                Default:  true,
			},

			// This is called to destroy the resource. No output is collected.
			"delete_cmd": &schema.Schema{
				Type:     schema.TypeString,
                Optional: true,
				Required: false,
				ForceNew: false,
			},

			// If false assumes deleted even if error
			"delete_break_on_error": &schema.Schema{
				Type:     schema.TypeBool,
                Optional: true,
				Required: false,
                Default:  true,
			},

			// Trim whitespace from outputs
			"trim_output": &schema.Schema{
				Type:     schema.TypeBool,
                Optional: true,
				Required: false,
                Default:  false,
			},

            // Will set to the latest output of the commands
			"last_read_stdout": &schema.Schema{
				Type:     schema.TypeString,
                Optional: true,
                //DiffSuppressFunc: lastReadDiffFunc,
			},
			"last_read_stderr": &schema.Schema{
				Type:     schema.TypeString,
                Optional: true,
                //DiffSuppressFunc: lastReadDiffFunc,
			},
			"last_read_retval": &schema.Schema{
				Type:     schema.TypeInt,
                Optional: true,
                //DiffSuppressFunc: lastReadDiffFunc,
			},

			// Our return values
            "stdout": &schema.Schema{
                Type:     schema.TypeString,
                Computed: true,
            },
            "stderr": &schema.Schema{
                Type:     schema.TypeString,
                Computed: true,
            },
            "retval": &schema.Schema{
                Type:     schema.TypeInt,
                Computed: true,
            },
		},
	}
}

func lastReadDiffFunc (k, oldVal, newVal string, d *schema.ResourceData) bool {
    if (len(oldVal) > 0 && len(newVal) == 0) {
        return true
    } else {
        return false
    }
}

func createAction(d *schema.ResourceData, m interface{}) error {
    config := m.(*Config)

    log.Printf("cli Creating")

    d.SetId(fmt.Sprintf("%d", rand.Int()))

    if cmdv, ok := d.GetOk("create_cmd"); ok {
        cmd := cmdv.(string)
        working_dir := d.Get("working_dir").(string)
        trim_output := d.Get("trim_output").(bool)
        create_break_on_error := d.Get("create_break_on_error").(bool)
        stdout, stderr, retval := run(config.Shell, working_dir, cmd)
        if (create_break_on_error && retval != 0) {
            return fmt.Errorf("'%s' returned a non-zero exit code.", cmd)
        }
        if (trim_output) {
            stdout = strings.TrimSpace(stdout)
            stderr = strings.TrimSpace(stderr)
        }

        d.Set("stdout", stdout)
        d.Set("stderr", stderr)
        d.Set("retval", retval)

        if (d.Get("read_on_create_update").(bool)) {
            return readAction(d, m)
        }

        return nil

    }

    if (d.Get("read_on_create_update").(bool)) {
        return readAction(d, m)
    }


	return nil
}

func readAction(d *schema.ResourceData, m interface{}) error {
    config := m.(*Config)

    log.Printf("cli Reading")

    if cmdv, ok := d.GetOk("read_cmd"); ok {
        cmd := cmdv.(string)
        working_dir := d.Get("working_dir").(string)
        trim_output := d.Get("trim_output").(bool)
        read_destroyed_on_error := d.Get("read_destroyed_on_error").(bool)
        read_break_on_error := d.Get("read_break_on_error").(bool)
        stdout, stderr, retval := run(config.Shell, working_dir, cmd)
        if (read_destroyed_on_error && retval != 0) {
            d.SetId("")
            return nil
        }
        if (read_break_on_error && retval != 0) {
            return fmt.Errorf("'%s' returned a non-zero exit code.", cmd)
        }

        if (trim_output) {
            stdout = strings.TrimSpace(stdout)
            stderr = strings.TrimSpace(stderr)
        }
        d.SetId(d.Id())
        d.Set("last_read_stdout", stdout)
        d.Set("last_read_stderr", stderr)
        d.Set("last_read_retval", retval)

    }

	return nil
}

func updateAction(d *schema.ResourceData, m interface{}) error {
    config := m.(*Config)

    log.Printf("cli Updating")

    if cmdv, ok := d.GetOk("update_cmd"); ok {
        cmd := cmdv.(string)
        working_dir := d.Get("working_dir").(string)
        trim_output := d.Get("trim_output").(bool)
        update_break_on_error := d.Get("update_break_on_error").(bool)
        stdout, stderr, retval := run(config.Shell, working_dir, cmd)
        if (update_break_on_error && retval != 0) {
            return fmt.Errorf("'%s' returned a non-zero exit code.", cmd)
        }
        if (trim_output) {
            stdout = strings.TrimSpace(stdout)
            stderr = strings.TrimSpace(stderr)
        }

        d.Set("stdout", stdout)
        d.Set("stderr", stderr)
        d.Set("retval", retval)

        if (d.Get("read_on_create_update").(bool)) {
            return readAction(d, m)
        }

        return nil

    }

    if (d.Get("read_on_create_update").(bool)) {
        return readAction(d, m)
    }

	return nil
}

func deleteAction(d *schema.ResourceData, m interface{}) error {
    config := m.(*Config)

    log.Printf("cli Deleting")

    if cmdv, ok := d.GetOk("delete_cmd"); ok {
        cmd := cmdv.(string)
        working_dir := d.Get("working_dir").(string)
        delete_break_on_error := d.Get("delete_break_on_error").(bool)
        _, _, retval := run(config.Shell, working_dir, cmd)
        if (delete_break_on_error && retval != 0) {
            return fmt.Errorf("'%s' returned a non-zero exit code.", cmd)
        }
    }

    d.SetId("")
	return nil
}

func run(shell string, working_dir string, str string) (string, string, int) {
    var stdout_buf bytes.Buffer
    var stderr_buf bytes.Buffer

    cmd := exec.Command(shell, "-c", str)
    cmd.Stdout = &stdout_buf
    cmd.Stderr = &stderr_buf

    if (working_dir != "") {
        cmd.Dir = working_dir
    }
    retval := 0


    err := cmd.Start()
    if err != nil {
        log.Fatal(err)
    }

    err = cmd.Wait()
    if err != nil {
        if exiterr, ok := err.(*exec.ExitError); ok {
            // The program has exited with an exit code != 0
            if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
                retval = status.ExitStatus()
            }
        } else {
            log.Fatalf("cmd.Wait: %v", err)
        }
        //log.Fatal(err)
    }
    log.Printf("cmd.stdout: %s\ncmd.stderr: %s\ncmd.retval: %d\n", strings.TrimSpace(stdout_buf.String()), strings.TrimSpace(stderr_buf.String()), retval)
    return stdout_buf.String(), stderr_buf.String(), retval
}


