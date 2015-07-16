package plugins

import (
    "path"

    "os/exec"
    "strings"

    "path/filepath"

    "io/ioutil"

    "regexp"
)

type Stylus struct {
    m *Manager
    exePath string
}

func NewStylus(m *Manager) *Stylus {
    return &Stylus{m: m}
}

func (p *Stylus) IsValid() bool {
    exePath, err := exec.LookPath("stylus")
    p.exePath = exePath

    if err != nil {
        return false
    }
    return true
}

func (p *Stylus) GetDependencies(path string) []string {

    if data, err := ioutil.ReadFile(path); err == nil {
        deps := p.getImports(string(data))
        return deps
    }

    return nil
}

func (p* Stylus) getImports(data string) []string {
    re := regexp.MustCompile("@(import|require)[ \t]+['\"](\\S+)['\"]")

    deps := []string{}

    for _,line := range strings.Split(string(data), "\n") {
        match := re.FindStringSubmatch(line)
        if match != nil {
            deps = append(deps, match[1])
        }
    }

    return deps
}

func (p *Stylus) getCmd(file string, realPath string) *exec.Cmd {
    dirToRun,_ := filepath.Abs(realPath)

    app := p.exePath
    //arg0 := "--compress"
    arg0 := "--disable-cache"
    arg1 := "-u"
    arg2 := "nib"
    arg3 := "-u"
    arg4 := "jeet"
    arg5 := "-u"
    arg6 := "axis"
    arg7 := "-u"
    arg8 := "rupture"
    arg9 := "<"
    arg10 := filepath.Join(dirToRun, file)

    if file != "" {
        cmd := exec.Command(app, arg0, arg1, arg2, arg3, arg4, arg5, arg6, arg7, arg8, arg9, arg10)
        cmd.Dir = dirToRun
        return cmd
    }

    cmd := exec.Command(app, arg0, arg1, arg2, arg3, arg4, arg5, arg6, arg7, arg8)
    cmd.Dir = dirToRun
    return cmd
}

func (p *Stylus) ProcessedFileName(fileName string) (string, bool) {
    ext := path.Ext(fileName)
    if ext == ".styl" {
        baseName := fileName[:len(fileName)-len(ext)]
        return baseName + ".css", true
    }
    return fileName, false
}

func (p *Stylus) ReverseFileName(fileName string) (string, bool) {
    ext := path.Ext(fileName)
    if ext == ".css" {
        baseName := fileName[:len(fileName)-len(ext)]
        return baseName + ".styl", true
    }
    return fileName, false
}

func (p *Stylus) ProcessFile(file string, relPath string, fileName string, realPath string) ([]byte, error) {
    cmd := p.getCmd("", realPath)
    data,_ := ioutil.ReadFile(file)
    cmd.Stdin = strings.NewReader(string(data))
    out, _ := cmd.CombinedOutput()

    return []byte(out), nil
}

func (p *Stylus) ProcessData(src []byte, relPath string, realPath string) ([]byte, error) {
    cmd := p.getCmd("", realPath)
    cmd.Stdin = strings.NewReader(string(src))
    out, err := cmd.CombinedOutput()

    return []byte(out), err
}
