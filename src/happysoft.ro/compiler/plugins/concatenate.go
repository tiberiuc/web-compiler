package plugins

import (
    "io/ioutil"
    "path"
    "errors"
    "strings"
    "bufio"
    "os"
)

type Concatenate struct {
    m *Manager
}

func NewConcatenate(m *Manager) *Concatenate {
    return &Concatenate{m: m}
}

func (p *Concatenate) IsValid() bool {
    return true
}

func (p *Concatenate) GetDependencies(path string) []string {
    file, err := os.Open(path)
    if err != nil {
        return nil
    }

    defer file.Close()

    lines := []string{}

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        lines = append(lines, scanner.Text())
    }

    return lines
}

func (p *Concatenate) ProcessedFileName(fileName string) (string, bool) {
    ext := path.Ext(fileName)
    if ext == ".concat" {
        baseName := fileName[:len(fileName)-len(ext)]
        return baseName, true
    }
    return fileName, false
}

func (p *Concatenate) ReverseFileName(fileName string) (string, bool) {
    ext := path.Ext(fileName)
    if ext == ".concat" {
        return fileName, false
    }
    return fileName + ".concat", true
}

func (p *Concatenate) processLine(line string, relPath string, realPath string) (string, bool) {

    fileName := strings.Trim(line, " \n\t")
    if fileName == "" {
        return "", true
    }

    newRelPath := path.Clean(path.Join(relPath, path.Dir(fileName))) + "/"


    result, ok := p.m.ProcessFile(newRelPath, path.Base(fileName))

    return string(result)+"\n", ok
}

func (p *Concatenate) ProcessFile(file string, relPath string, fileName string, realPath string) ([]byte, error) {
    content, err := ioutil.ReadFile(file)
    if err != nil {
        return content, err
    }

    result, err := p.ProcessData([]byte(content), relPath, realPath)
    if err != nil {
        return result, errors.New(err.Error() + " in " + file)
    }

    return result, nil
}

func (p *Concatenate) ProcessData(src []byte, relPath string, realPath string) ([]byte, error) {
    result := ""
    lines := strings.Split(string(src), "\n")
    for _, line := range lines {
        out, ok := p.processLine(line, relPath, realPath)
        result += out
        if !ok {

            errText := "Error processing " + relPath + line
            result += "\n" + errText
            return []byte(result), errors.New(errText)
        }
    }

    return []byte(result), nil
}
