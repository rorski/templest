## Templest
Templest is a tool to build a directory structure from go templates, using a layout and variables defined in a YAML configuration. It is similar in intent to [cookiecutter](https://github.com/cookiecutter/cookiecutter), with some inspiration drawn from [helm](https://github.com/helm/helm).

In short, templest:
1. Takes a YAML configuration and template directory as inputs
2. Parses the variables defined in that YAML configuration for each YAML key (i.e., each subdirectory)
3. Creates the directory defined by that key at the appropriate place in the output directory
4. Renders any [go template files](https://pkg.go.dev/text/template) (those with a `.tmpl` suffix) in the template directory to the output directory with the `.tmpl` suffix stripped off and the variables filled in
5. Copies any non-template files (those without a `.tmpl` suffix) to the output directory, including symlinks.

The basic usage is:
```
Usage of templest:

  -config string
        path to yaml formatted stack configuration
  -out string
        path to render directories and files to
  -templates string
        path to templates dir
```
For example, `./templest -config example/config.yml -templates example/templates -out example/out`

See [the example directory](https://github.com/rorski/templest/example) for a sample of a YAML configuration and the output it produces.
### Building
`make build` or just `go build`
### Why though?
This was initially created to help quickly create [terraform module structures](https://www.terraform.io/language/modules/develop/structure) for various clusters based on a simple input file that could generate that project structure.

Cookiecutter is a feasible tool for this but had some limitations (at least as of v1.7) that inspired this work:
* Cookiecutter [only supports JSON configuration](https://github.com/cookiecutter/cookiecutter/issues/970#issuecomment-336695070), which in my opinion is not a good configuration language. Templest uses YAML.
* I prefered a tool that runs as a single binary (hence Go)
* Cookiecutter doesn't support copying symlinks natively - you need to use a hook for that.

Cookiecutter also provides a number of features that this project doesn't, such as pointing at a cookiecutter formatted git repository and pre/post-hooks. Some of those features may be implemeted in future versions of templest.