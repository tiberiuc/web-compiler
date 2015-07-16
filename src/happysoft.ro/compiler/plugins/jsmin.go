package plugins

import (
    "path"
    "path/filepath"

    "io/ioutil"
    "os/exec"
    "strings"
    "regexp"
)

type JSMin struct {
    m *Manager
    exePath string
}

func NewJSMin(m *Manager) *JSMin {
    return &JSMin{m: m}
}

func (p *JSMin) IsValid() bool {
    exePath, err := exec.LookPath("uglifyjs")
    p.exePath = exePath

    if err != nil {
        return false
    }
    return true
}

func (p *JSMin) getCmd(file string, realPath string) *exec.Cmd {
    dirToRun,_ := filepath.Abs(realPath)

    app := p.exePath
    arg0 := "-b"
    arg1 := "beautify=true"
    if file != "" {
        file = filepath.Base(file)
        cmd := exec.Command(app, arg0, arg1, file)
        cmd.Dir = dirToRun
        return cmd
    }

    cmd := exec.Command(app, arg0, arg1)
    cmd.Dir = dirToRun
    return cmd
}

func (p *JSMin) GetDependencies(path string) []string {
    return nil
}

func (p *JSMin) ProcessedFileName(fileName string) (string, bool) {
    ext := path.Ext(fileName)
    baseName := fileName[:len(fileName)-len(ext)]
    ext2 := path.Ext(baseName)
    if ext == ".js" && ext2 != ".min" {
        return baseName + ".min.js", true
    }
    return fileName, false
}

func (p *JSMin) ReverseFileName(fileName string) (string, bool) {
    ext := path.Ext(fileName)
    baseName := fileName[:len(fileName)-len(ext)]
    ext2 := path.Ext(baseName)
    baseName = baseName[:len(baseName)-len(ext2)]
    if ext == ".js" && ext2 == ".min" {
        return baseName + ".js", true
    }
    return fileName, false
}

func (p *JSMin) ProcessFile(file string, relPath string, fileName string, realPath string) ([]byte, error) {


    if p.m.Config.General.Development {
        cmd := p.getCmd(file, realPath)
        out, err := cmd.CombinedOutput()
        return []byte(out), err
    }

    out, err := ioutil.ReadFile(file)

    if err != nil {
        return out, err
    }

    return p.ProcessData(out, relPath, realPath)
}

func (p *JSMin) ProcessData(src []byte, relPath string, realPath string) ([]byte, error) {
    var err error
    var out []byte

    srcModified := string(src)

    if !p.m.Config.General.Development {
        //Scot console.log


        if p.m.Config.Jsmin.Removelogs {
            re1 := regexp.MustCompile("[ \\t]*return ((i?\\$|(parent\\.)?window|console)\\.t?(log|warn)\\(.*\\);)")
            re2 := regexp.MustCompile("[ \\t]*((i?\\$?|(parent\\.)?window|console|window\\.console)\\.t?(log|warn)\\(.*\\);)")
            srcModified = re1.ReplaceAllString(srcModified, "return undefined;")
            srcModified = re2.ReplaceAllString(srcModified, "")
        }

    }

    cmd := p.getCmd("", realPath)
    cmd.Stdin = strings.NewReader(srcModified)
    out, err = cmd.CombinedOutput()

    return out, err
}
