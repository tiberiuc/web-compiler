package plugins

import (
    "path"

    "os/exec"
    "strings"

    "path/filepath"

    "io/ioutil"

    "regexp"
)

type Jade struct {
    m *Manager
    exePath string
}

func NewJade(m *Manager) *Jade {
    return &Jade{m: m}
}

func (p *Jade) IsValid() bool {
    exePath, err := exec.LookPath("jade")
    p.exePath = exePath

    if err != nil {
        return false
    }
    return true
}

func (p *Jade) GetDependencies(path string) []string {

    if data, err := ioutil.ReadFile(path); err == nil {
        deps := p.getImports(string(data))
        return deps
    }

    return nil
}

func (p* Jade) getImports(data string) []string {
    re := regexp.MustCompile("^[ \t]*include[ \t]+(\\S+)")

    deps := []string{}

    for _,line := range strings.Split(string(data), "\n") {
        match := re.FindStringSubmatch(line)
        if match != nil {
            deps = append(deps, match[1])
        }
    }

    return deps
}

func (p *Jade) getCmd(file string, realPath string) *exec.Cmd {
    dirToRun,_ := filepath.Abs(realPath)

    app := p.exePath
    arg0 := "-P"
    arg1 := "-p"
    arg2 := filepath.Join(dirToRun, file)

    if file != "" {
        cmd := exec.Command(app, arg0, arg1, arg2, file)
        cmd.Dir = dirToRun
        return cmd
    }

    arg2 = filepath.Join(arg2, ".__fakefile__.jade")

    cmd := exec.Command(app, arg0, arg1, arg2)
    cmd.Dir = dirToRun
    return cmd
}

func (p *Jade) ProcessedFileName(fileName string) (string, bool) {
    ext := path.Ext(fileName)
    if ext == ".jade" {
        baseName := fileName[:len(fileName)-len(ext)]
        return baseName + ".html", true
    }
    return fileName, false
}

func (p *Jade) ReverseFileName(fileName string) (string, bool) {
    ext := path.Ext(fileName)
    if ext == ".html" {
        baseName := fileName[:len(fileName)-len(ext)]
        return baseName + ".jade", true
    }
    return fileName, false
}

func (p *Jade) ProcessFile(file string, relPath string, fileName string, realPath string) ([]byte, error) {
    cmd := p.getCmd("", realPath)
    data,_ := ioutil.ReadFile(file)
    cmd.Stdin = strings.NewReader(string(data))
    out, _ := cmd.CombinedOutput()

    return []byte(out), nil
}

func (p *Jade) ProcessData(src []byte, relPath string, realPath string) ([]byte, error) {
    cmd := p.getCmd("", realPath)
    cmd.Stdin = strings.NewReader(string(src))
    out, err := cmd.CombinedOutput()

    return []byte(out), err
}
