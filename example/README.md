### Example configuration
The templest YAML config can create an arbitrarily deep directory structure by simply nesting YAML keys. Variables for each directory level are defined by the special `_vars` key, like so:
```
  vpc:
    _vars:
      sourcePath: "terraform-aws-modules/vpc/aws"
```
This config would look for a directory named `vpc` in the templates directory at the configured level, then create a variable `sourcePath` to be written to any file with a `.tmpl` suffix in that directory, e.g. `main.tf.tmpl`. In a template file, that variable substitution would be defined like this:
```
source = "{{ .Vars.sourcePath }}"
```
Subdirectories are created by just adding the directory name as a sub-key, then defining whatever variables you want for that directory under its own `_vars`:
```
  vpc:
    _vars:
      sourcePath: "terraform-aws-modules/vpc/aws"
    vpc_endpoints:
      _vars:
        createDns: true
```
templest will only create and write to the output directories as specified in the YAML input file. If there are directories in the template dir that aren't included in the YAML config, they will be skipped.

See the `out/` directory for what this looks like once you run the tool. This was generated from the parent directory with:
`./templest -config example/config.yml -templates example/templates -out example/out`

### Note:
* `_vars` is not a valid subdirectory because it is reserved for defining variables
* `_meta` is also a reserved keyword.
* This terraform code in the example likely won't actually deploy, it's just a demonstration of templest functionality.