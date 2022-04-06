package main

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// Render takes an input template configuration of type Layout and creates a directory
// structure with rendered Go templates from the configuration variables specified.
func (c Config) Render(layout map[string]interface{}) error {
	parsed, err := parseLayout(layout, nil)
	if err != nil {
		return fmt.Errorf("error parsing data: %v", err)
	}

	err = c.walkLayout(parsed, "")
	if err != nil {
		return fmt.Errorf("error creating dirs: %v", err)
	}

	return nil
}

// parseLayout creates a linked list of variables
func parseLayout(yamlConfig map[string]interface{}, v *Vars) (*Vars, error) {
	if v == nil {
		v = new(Vars)
	}
	for k := range yamlConfig {
		switch k {
		case "tmpl_vars":
			v.TmplVars = yamlConfig[k].(map[string]interface{})
		default:
			// recursively parse subdirectories
			if yamlConfig[k] != nil {
				sub, err := parseLayout(yamlConfig[k].(map[string]interface{}), &Vars{Name: k})
				if err != nil {
					return nil, fmt.Errorf("error parsing layout: %v", err)
				}
				v.SubDirs = append(v.SubDirs, sub)
			} else {
				// if there are no template variables defined, just create the subdir name,
				// so any files in that template directory will still be copied to the output
				v.SubDirs = append(v.SubDirs, &Vars{Name: k})
			}
		}
	}

	return v, nil
}

// walkLayout walks the parsed YAML config, rendering and copying files at each directory level
func (c Config) walkLayout(v *Vars, path string) error {
	err := os.MkdirAll(filepath.Join(c.OutPath, path, v.Name), 0755)
	if err != nil {
		return fmt.Errorf("error creating directory: %v", err)
	}
	err = c.handleFiles(v, path)
	if err != nil {
		return fmt.Errorf("error handling file at %s: %v", path, err)
	}
	if len(v.SubDirs) != 0 {
		// recursively walk the parsed yaml for each subdirectory defined
		for _, s := range v.SubDirs {
			err := c.walkLayout(s, filepath.Join(path, v.Name))
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (c Config) handleFiles(v *Vars, path string) error {
	tmplDir := filepath.Join(c.TemplatePath, path, v.Name)
	outDir := filepath.Join(c.OutPath, path, v.Name)

	files, err := os.ReadDir(tmplDir)
	if err != nil {
		return fmt.Errorf("error reading files in directory: %v", err)
	}

	for _, f := range files {
		srcPath := filepath.Join(tmplDir, f.Name())

		fileInfo, err := f.Info()
		if err != nil {
			return fmt.Errorf("error getting file info for %s: %v", f.Name(), err)
		}

		switch fileInfo.Mode() & os.ModeType {
		case fs.ModeDir:
			continue
		case fs.ModeSymlink:
			link, err := os.Readlink(srcPath)
			if err != nil {
				return fmt.Errorf("could not read destination file for link %s: %v", srcPath, err)
			}
			err = createLink(filepath.Join(outDir, fileInfo.Name()), link)
			if err != nil {
				return err
			}
		default:
			if strings.HasSuffix(f.Name(), ".tmpl") {
				err = v.renderTemplate(f, tmplDir, outDir)
				if err != nil {
					return fmt.Errorf("error rendering template %s: %v", f.Name(), err)
				}
			} else {
				err = copyFiles(f.Name(), srcPath, outDir)
				if err != nil {
					return fmt.Errorf("error copying file %s: %v", f.Name(), err)
				}
			}
		}
	}

	return nil
}

func (v *Vars) renderTemplate(file fs.DirEntry, tmplDir string, outDir string) error {
	t, err := template.New("").Funcs(template.FuncMap{
		"HCLJoin": HCLJoin,
	}).ParseFiles(filepath.Join(tmplDir, file.Name()))
	if err != nil {
		return fmt.Errorf("error parsing template file %s: %v", file.Name(), err)
	}
	// create the rendered output file without the ".tmpl" suffix
	outFile, err := os.Create(filepath.Join(outDir, strings.TrimSuffix(file.Name(), ".tmpl")))
	if err != nil {
		return fmt.Errorf("could not create output file: %v", err)
	}
	// execute the template with the base name as parsed above
	// see https://stackoverflow.com/questions/44979276/the-go-template-parsefiles-func-parsing-multiple-files
	err = t.ExecuteTemplate(outFile, file.Name(), v)
	if err != nil {
		return fmt.Errorf("could not execute template for %s: %v", outFile.Name(), err)
	}
	defer outFile.Close()

	return nil
}

// copy over any regular files to the output from the template dir
func copyFiles(fileName, srcPath, outDir string) error {
	src, err := os.Open(srcPath)
	if err != nil {
		return fmt.Errorf("error opening source file %s: %v", fileName, err)
	}
	defer src.Close()

	dst, err := os.Create(filepath.Join(outDir, fileName))
	if err != nil {
		return fmt.Errorf("error creating file %s: %v", fileName, err)
	}
	defer dst.Close()
	_, err = io.Copy(dst, src)
	if err != nil {
		return err
	}
	return nil
}

func createLink(src, dst string) error {
	err := os.Symlink(dst, src)
	if os.IsExist(err) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("could not create link %s: %v", src, err)
	}
	return nil
}

// HCLJoin converts a list of interfaces to a list of strings in HCL/Terraform friendly format
// For example, []string{"this", "that", "the other"} wouid produce: ["this", "that", "the other"]
func HCLJoin(values []interface{}) string {
	var s []string
	for _, v := range values {
		v := fmt.Sprintf("\"%s\"", v.(string))
		s = append(s, v)
	}

	return fmt.Sprintf("[%s]", strings.Join(s, ", "))
}
