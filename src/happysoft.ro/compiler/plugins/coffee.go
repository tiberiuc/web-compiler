package plugins

import (
    "os/exec"
    "path"
    "path/filepath"
    "strings"
)

type Coffeescript struct {
    m *Manager
    exePath string
}

func NewCoffeescript(m *Manager) *Coffeescript {
    return &Coffeescript{m: m}
}

func (p *Coffeescript) IsValid() bool {

    exePath, err := exec.LookPath("coffee")
    p.exePath = exePath

    if err != nil {
        return false
    }

    return true
}

func (p *Coffeescript) GetDependencies(path string) []string {
    return nil
}

func (p *Coffeescript) getCmd(file string, realPath string) *exec.Cmd {
    dirToRun,_ := filepath.Abs(realPath)

    app := p.exePath
    arg0 := "-c"
    arg1 := "-p"
    if file == "" {
        file = "-s"
    } else {
        file = filepath.Base(file)
    }

    cmd := exec.Command(app, arg0, arg1, file)
    cmd.Dir = dirToRun

    return cmd
}

func (p *Coffeescript) ProcessedFileName(fileName string) (string, bool) {
    ext := path.Ext(fileName)
    if ext == ".coffee" {
        baseName := fileName[:len(fileName)-len(ext)]
        return baseName + ".js", true
    }
    return fileName, false
}

func (p *Coffeescript) ReverseFileName(fileName string) (string, bool) {
    ext := path.Ext(fileName)
    if ext == ".js" {
        baseName := fileName[:len(fileName)-len(ext)]
        return baseName + ".coffee", true
    }
    return fileName, false
}

func (p *Coffeescript) ProcessFile(file string, relPath string, fileName string, realPath string) ([]byte, error) {
    cmd := p.getCmd(file, realPath)
    out, err := cmd.CombinedOutput()

    return []byte(out), err
}

func (p *Coffeescript) ProcessData(src []byte, relPath string, realPath string) ([]byte, error) {
    cmd := p.getCmd("", realPath)
    cmd.Stdin = strings.NewReader(string(src))
    out, err := cmd.CombinedOutput()

    return []byte(out), err
}
