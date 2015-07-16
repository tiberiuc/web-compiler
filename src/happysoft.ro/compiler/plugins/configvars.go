package plugins

import (
    "fmt"
    "io/ioutil"
    "os"
    "path/filepath"
    "strings"
)

type Configvars struct {
    m       *Manager
    files   string
    vars    map[string]string
    include []string
    exclude []string
}

func NewConfigvars(m *Manager) *Configvars {
    p := &Configvars{m: m}

    p.files = m.Config.Configvars.Files

    files := filepath.SplitList(p.files)

    p.vars = make(map[string]string, 0)

    for _, file := range files {
        src, err := ioutil.ReadFile(file)
        if err != nil {
            fmt.Fprintf(os.Stderr, "Error loading Configvars file %s\n\n%s\n", file, err.Error())
            os.Exit(2)
        }

        lines := strings.Split(string(src), "\n")
        for _, line := range lines {
            p.processLineVar(line)
        }
    }

    // process added vars
    for _, line := range m.Config.Configvars.addedvars{
        p.processLineVar(line)
    }

    p.processExtensions()

    return p
}

func (p *Configvars) processExtensions() {
    if p.m.Config.Configvars.Include != "" {
        matches := strings.Split(p.m.Config.Configvars.Include, ",")
        for _, match := range matches {
            match = strings.Trim(match, " \t")
            p.include = append(p.include, match)
        }
    }
    if p.m.Config.Configvars.Exclude != "" {
        matches := strings.Split(p.m.Config.Configvars.Exclude, ",")
        for _, match := range matches {
            match = strings.Trim(match, " \t")
            p.exclude = append(p.exclude, match)
        }
    }
}

func (p *Configvars) processLineVar(line string) {
    line = strings.Trim(line, " \t")
    if line != "" {
        idx := strings.Index(line, "=")
        varName := ""
        varValue := ""
        if idx != -1 {
            varName = line[:idx]
            if idx+1 < len(line) {
                varValue = line[idx+1:]
            }
        } else {
            varName = line
        }
        varName = strings.Trim(varName, " \t")
        varValue = strings.Trim(varValue, " \t")

        if varName != "" {
            p.vars[varName] = varValue
        }
    }

}

func (p *Configvars) ProcessData(src []byte, relPath string, fileName string) ([]byte, error) {
    shouldProcess := true

    if len(p.include) > 0 {
        shouldProcess = false
        for _, match := range p.include {
            if matched, _ := filepath.Match(match, fileName); matched {
                shouldProcess = true
                break
            }
        }
    } else {
        for _, match := range p.exclude {
            if matched, _ := filepath.Match(match, fileName); matched {
                shouldProcess = false
                break
            }
        }
    }

    if shouldProcess {
        data := string(src)
        for name, value := range p.vars {
            data = strings.Replace(data, "{{"+name+"}}", value, -1)
        }
        return []byte(data), nil
    }

    return src, nil
}
