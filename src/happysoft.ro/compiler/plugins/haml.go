package plugins

import (
    "path"

    "os/exec"
    "strings"
    "path/filepath"
)

type Haml struct {
    m *Manager
    exePath string
}

func NewHaml(m *Manager) *Haml {
    return &Haml{m: m}
}

func (p *Haml) IsValid() bool {
    exePath, err := exec.LookPath("haml")
    p.exePath = exePath

    if err != nil {
        return false
    }
    return true
}

func (p *Haml) GetDependencies(path string) []string {
    return nil
}

func (p *Haml) getCmd(file string, realPath string) *exec.Cmd {
    dirToRun,_ := filepath.Abs(realPath)

    app := p.exePath
    if file != "" {
        file = filepath.Base(file)
        cmd := exec.Command(app, file)
        cmd.Dir = dirToRun
        return cmd
    }

    cmd := exec.Command(app, file)
    cmd.Dir = dirToRun
    return cmd
}

func (p *Haml) ProcessedFileName(fileName string) (string, bool) {
    ext := path.Ext(fileName)
    if ext == ".haml" {
        baseName := fileName[:len(fileName)-len(ext)]
        return baseName + ".html", true
    }
    return fileName, false
}

func (p *Haml) ReverseFileName(fileName string) (string, bool) {
    ext := path.Ext(fileName)
    if ext == ".html" {
        baseName := fileName[:len(fileName)-len(ext)]
        return baseName + ".haml", true
    }
    return fileName, false
}

func (p *Haml) ProcessFile(file string, relPath string, fileName string, realPath string) ([]byte, error) {
    cmd := p.getCmd(file, realPath)
    out, err := cmd.CombinedOutput()

    return []byte(out), err
}

func (p *Haml) ProcessData(src []byte, relPath string, realPath string) ([]byte, error) {
    cmd := p.getCmd("", realPath)
    cmd.Stdin = strings.NewReader(string(src))
    out, err := cmd.CombinedOutput()

    return []byte(out), err
}
