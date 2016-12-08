### Disclaimer

This is the worst code I've written in a long time. It doesn't use types, it's
badly commented, there are short variable names and zero tests.

TL;DR: Here be dragons

### The problem

Whilst working on some [Terraform](https://www.terraform.io/) files I lost quite
a lot of time to the fact that I'd copied and pasted a resource and forgotten to 
change it's name. This meant that when Terraform tried to create the second one
it errored out with a non-helpful error message that the resource was already in
use.

### The solution

This script reads HCL from `STDIN` and outputs the name of any resources that
have the same value for a specified field. You'll need to edit `main.go` to 
specify your searches (sorry about that)

Given the following input:

**main.go**
```
searches["azurerm_network_security_group"]
```

Running the following command:

```
cat /path/to/project/*.tf| go run main.go
```

Will output something that looks like the following if there are duplicates:

```
#########################
## azurerm_network_security_group
#########################
 - name :: [api_sg db_sg]
```

### Why Go?
I wanted to reuse the actual Terraform HCL parser and not write one myself.

