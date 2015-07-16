package plugins

import (
    "path"

    "os/exec"
    "io/ioutil"
    "strings"

    "path/filepath"

    "errors"

    "regexp"
)

type Sass struct {
    m *Manager
    exePath string
}

func NewSass(m *Manager) *Sass {
    return &Sass{m: m}
}

func (p *Sass) IsValid() bool {
    exePath, err := exec.LookPath("sass")
    p.exePath = exePath

    if err != nil {
        return false
    }
    return true
}

func (p *Sass) GetDependencies(path string) []string {

    if data, err := ioutil.ReadFile(path); err == nil {
        deps := p.getImports(string(data))
        return deps
    }

    return nil
}

func (p* Sass) getImports(data string) []string {
    re := regexp.MustCompile("^[ \t]*@import[ \t]+['\"]{1}(.+)['\"]{1}")

    deps := []string{}

    for _,line := range strings.Split(string(data), "\n") {
        match := re.FindStringSubmatch(line)
        if match != nil {
            deps = append(deps, match[1])
        }
    }

    return deps
}

func (p *Sass) getCmd(file string, realPath string) *exec.Cmd {
    dirToRun,_ := filepath.Abs(realPath)

    app := p.exePath
    arg0 := "--no-cache"
    arg1 := "--style"
    arg2 := "compressed"
    if file != "" {
        file = filepath.Base(file)
        cmd := exec.Command(app, arg0, arg1, arg2, file)
        cmd.Dir = dirToRun
        return cmd
    }

    cmd := exec.Command(app, arg0, arg1, arg2, file)
    cmd.Dir = dirToRun
    return cmd
}

func (p *Sass) ProcessedFileName(fileName string) (string, bool) {
    ext := path.Ext(fileName)
    if ext == ".sass" {
        baseName := fileName[:len(fileName)-len(ext)]
        return baseName + ".css", true
    }
    return fileName, false
}

func (p *Sass) ReverseFileName(fileName string) (string, bool) {
    ext := path.Ext(fileName)
    if ext == ".css" {
        baseName := fileName[:len(fileName)-len(ext)]
        return baseName + ".sass", true
    }
    return fileName, false
}

func (p *Sass) processImports(relPath string, data []byte) bool {
    deps := p.getImports(string(data))

    if deps != nil {
        for _, dep := range deps {
            dep := strings.Trim(dep, " \n\t")
            if dep != "" {
                newRelPath := path.Clean(path.Join(relPath, path.Dir(dep))) + "/"

                _, ok := p.m.ProcessFile(newRelPath, path.Base(dep))
                if !ok {
                    // Intorc true chiar daca este o eroare doar pentru
                    // crearea dependintelor
                    // eroarea efectiva se trimite de catre SCSS
                    return true
                }
            }

        }
    }

    return true

}

func (p *Sass) ProcessFile(file string, relPath string, fileName string, realPath string) ([]byte, error) {

    if data, err := ioutil.ReadFile(file); err == nil {
        ok := p.processImports(relPath, data)
        if !ok {
            return nil, errors.New("Error procesing imports")
        }
    }

    cmd := p.getCmd(file, realPath)
    out, err := cmd.CombinedOutput()

    return []byte(out), err
}

func (p *Sass) ProcessData(src []byte, relPath string, realPath string) ([]byte, error) {

    ok := p.processImports(relPath, src)
    if !ok {
        return nil, errors.New("Error procesing imports")
    }

    cmd := p.getCmd("", realPath)
    cmd.Stdin = strings.NewReader(string(src))
    out, err := cmd.CombinedOutput()

    return []byte(out), err
}
