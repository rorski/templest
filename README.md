## Templizer
Templizer is a tool to build a directory structure from go templates, using a layout and variables defined in a YAML configuration. It is similar in intent to [cookiecutter](https://github.com/cookiecutter/cookiecutter), with some inspiration drawn from [helm](https://github.com/helm/helm).

In short, Templizer:
1. Takes a YAML configuration and template directory as inputs
2. Parses the variables defined in that YAML configuration for each YAML key (i.e., each subdirectory)
3. Creates the directory defined by that key at the appropriate place in the output directory
4. Renders any [go template files](https://pkg.go.dev/text/template) (those with a `.tmpl` suffix) in the template directory to the output directory with the `.tmpl` suffix stripped off and the variables filled in
5. Copies any non-template files (those without a `.tmpl` suffix) to the output directory, including symlinks.

The basic usage is:
```
Usage of templizer:

  -config string
        path to yaml formatted stack configuration
  -out string
        path to render directories and files to
  -templates string
        path to templates dir
```
For example, `./templizer -config example/config.yml -templates example/templates -out example/out`

See [the examples directory](https://github.com/rorski/templizer/examples) for a sample of a YAML configuration and the output it produces.

### Why though?
This was initially created to help quickly create [terraform module structures](https://www.terraform.io/language/modules/develop/structure) for various clusters based on a simple input file that could generate that project structure.

Cookiecutter is a feasible tool for this but had some limitations (at least as of v1.7) that inspired this work:
* Cookiecutter [only supports JSON configuration](https://github.com/cookiecutter/cookiecutter/issues/970#issuecomment-336695070), which in my opinion is not a good configuration language. Templizer uses YAML.
* I prefered a tool that runs as a single binary (hence Go)
* Cookiecutter doesn't support copying symlinks natively - you need to use a hook for that.
* I prefer Go templates to Jinja2 mostly due to familiarity, but also because the latter [doesn't allow for dashes in variable names](https://stackoverflow.com/questions/27779024/setting-data-attributes-on-a-wtforms-field) which is frustrating if your template variable happens to have a dash in it.

Cookiecutter also provides a number of features that this project doesn't, such as pointing at a cookiecutter formatted git repository and pre/post-hooks. Some of those features may be implemeted in future versions of templizer.