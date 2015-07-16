package plugins

import (
    "io/ioutil"
    "os"
    "path"
    "path/filepath"
    "sort"
    "strings"
    "time"

)

var pluginManager *Manager

func StartPluginManager() (*Manager) {
    pluginManager = new(Manager)
    pluginManager.Init()

    return pluginManager
}

type Plugin interface {
    GetDependencies(path string) []string
    IsValid() bool
    ProcessedFileName(fileName string) (string, bool)
    ReverseFileName(fileName string) (string, bool)
    ProcessFile(file string, relPath string, fileName string, realPath string) ([]byte, error)
    ProcessData(src []byte, relPath string, realPath string) ([]byte, error)
}

type PostProcessPlugin interface {
    ProcessData(src []byte, relPath string, fileName string) ([]byte, error)
}

type Manager struct {
    PathList           []string
    Plugins            []Plugin
    PostProcessPlugins []PostProcessPlugin
    Config             *Config
    Cache              *Cache
}

func RemoveDuplicates(xs *[]string) {
    found := make(map[string]bool)
    j := 0
    for i, x := range *xs {
        if !found[x] {
            found[x] = true
            (*xs)[j] = (*xs)[i]
            j++
        }
    }
    *xs = (*xs)[:j]
}

func (m *Manager) Init() {
//    m.Config = new(Config)
    m.Config = PluginsConfig
    m.Config.getConfig()

    m.PathList = filepath.SplitList(m.Config.General.Path)
    m.initPlugins()

    m.Cache = NewCache(m)
}

func (m *Manager) initPlugins() {

    if !m.Config.Coffeescript.Disabled {
        p := NewCoffeescript(m)
        if p.IsValid() {
            m.Plugins = append(m.Plugins, p)
        }
    }
    if !m.Config.Haml.Disabled {
        p := NewHaml(m)
        if p.IsValid() {
            m.Plugins = append(m.Plugins, p)
        }
    }
    if !m.Config.Jade.Disabled {
        p := NewJade(m)
        if p.IsValid() {
            m.Plugins = append(m.Plugins, p)
        }
    }
    if !m.Config.Stylus.Disabled {
        p := NewStylus(m)
        if p.IsValid() {
            m.Plugins = append(m.Plugins, p)
        }
    }
    if !m.Config.Sass.Disabled {
        p := NewSass(m)
        if p.IsValid() {
            m.Plugins = append(m.Plugins, p)
        }
    }
    if !m.Config.Scss.Disabled {
        p := NewScss(m)
        if p.IsValid() {
            m.Plugins = append(m.Plugins, p)
        }
    }
    if !m.Config.Jsmin.Disabled {
        p := NewJSMin(m)
        if p.IsValid() {
            m.Plugins = append(m.Plugins, p)
        }
    }
    if !m.Config.Concatenate.Disabled {
        p := NewConcatenate(m)
        if p.IsValid() {
            m.Plugins = append(m.Plugins, p)
        }
    }


    if !m.Config.Configvars.Disabled {
        m.PostProcessPlugins = append(m.PostProcessPlugins, NewConfigvars(m))
    }
}

func (m *Manager) ListFolders(relPath string) ([]string, error) {
    list := make([]string, 0)

    // Verifica daca in calea relativa sunt foldere ascunse ( incep cu _ )
    folders := strings.Split(relPath, "/")
    for _, f := range folders {
        if len(f) != 0 && f[0] == '_' {
            return list, nil
        }
    }

    relPath = m.FixHiddenDirs(relPath)

    for _, rootPath := range m.PathList {
        listPath := path.Join(rootPath, relPath)

        d, err := os.Open(listPath)
        if err == nil {

            fis, err := d.Readdir(-1)
            if err == nil {

                for _, fi := range fis {
                    if fi.IsDir() {
                        if fi.Name()[0] != '_' {
                            list = append(list, fi.Name())
                        } else if m.Config.General.Development {
                            list = append(list, "." + fi.Name())
                        }
                    }
                }
            }
        }
    }

    RemoveDuplicates(&list)

    return list, nil
}

func (m *Manager) ProcessFileName(fileName string, list *[]string, firstRun bool) (string, bool) {
    var ok bool

    processed := false
    newFile := fileName
    for _, p := range m.Plugins {
        newFile, ok = p.ProcessedFileName(newFile)
        if ok {
            processed = true
            newFile, _ = m.ProcessFileName(newFile, list, false)
            break
        }
    }

    if m.Config.General.Development && processed && !firstRun {
        *list = append(*list, "."+fileName)
    }
    return newFile, processed
}

func (m *Manager) FixHiddenDirs( relPath string) string {
    newPath := ""

    // Verifica daca in calea relativa sunt foldere ascunse ( incep cu _ )
    folders := strings.Split(relPath, "/")
    for _, f := range folders {
        p := f
        if len(f) >= 2 && f[0] == '.' && f[1] == '_' {
            p = string(f[1:])
        }
        newPath = path.Join(newPath, p)
    }

    return newPath
}


func (m *Manager) ListFiles(relPath string) []string {
    list := make([]string, 0)

    // Verifica daca in calea relativa sunt foldere ascunse ( incep cu _ )
    folders := strings.Split(relPath, "/")
    for _, f := range folders {
        if len(f) != 0 && f[0] == '_' {
            return list
        }
    }

    relPath = m.FixHiddenDirs(relPath)

    // citesc toate fisierele din director
    for _, rootPath := range m.PathList {
        listPath := path.Join(rootPath, relPath)

        d, err := os.Open(listPath)
        if err == nil {

            fis, err := d.Readdir(-1)
            if err == nil {

                for _, fi := range fis {
                    if fi.Mode().IsRegular() {
                        if fi.Name()[0] != '_' {
                            list = append(list, fi.Name())
                        }
                    }
                }
            }
        }
    }

    RemoveDuplicates(&list)

    // Procesez fiecare fisier
    processedList := make([]string, 0)

    for _, f := range list {
        newFile, _ := m.ProcessFileName(f, &processedList, true)

        processedList = append(processedList, newFile)
    }

    RemoveDuplicates(&processedList)

    sort.Strings(processedList)

    return processedList
}

func (m *Manager) getRealFilePathFromRelPath(relPath string, fileName string) (string, bool) {
    for _, rootPath := range m.PathList {
        if _, err := os.Stat(path.Join(rootPath, relPath, fileName)); err == nil {
            return path.Join(rootPath, relPath, fileName), true
        }
    }

    return "", false
}

func (m *Manager) GetFileContents(relPath string, fileName string) ([]byte, bool) {
    fName, ok := m.getRealFilePathFromRelPath(relPath, fileName)
    if ok {
        data, _ := ioutil.ReadFile(fName)
        return data, true
    }
    return nil, false
}

func (m *Manager) recursiveProcessData(relPath string, fileName string, prevFileName string, data []byte, origFileName string, realPath string) ([]byte, bool) {

    if fileName == origFileName {
        return data, true
    }

    var err error
    var ok bool
    var revName string
    var newData []byte

    for _, p := range m.Plugins {
        revName, ok = p.ProcessedFileName(fileName)
        if ok {
            newData, err = p.ProcessData(data, relPath, realPath)
            if err != nil {
                return newData, false
            }

            deps := []string{filepath.Join(relPath, revName)}
            fileDeps := p.GetDependencies(filepath.Join(relPath, revName))
            for _,fileDep := range fileDeps {
                deps = append(deps, fileDep)
            }
            m.Cache.Set(filepath.Join(relPath, prevFileName), time.Now(), true, &deps, data)
            return m.recursiveProcessData(relPath, revName, fileName, newData, origFileName, realPath)
        }
    }

    return data, true

}

func (m *Manager) recursiveProcessFile(relPath string, fileName string, prevFileName string, origFileName string) ([]byte, bool) {

    var ok bool
    var revName string
    var data []byte
    var prevData []byte
    var err error

    for _, p := range m.Plugins {
        revName, ok = p.ReverseFileName(fileName)
        if ok {
            fName, fileExist := m.getRealFilePathFromRelPath(relPath, revName)
            if fileExist {
                prevData, err = p.ProcessFile(fName, relPath, fileName, path.Dir(fName))
                if err != nil {
                    return prevData, false
                }
                data, ok = m.recursiveProcessData(relPath, fileName, fileName, prevData, origFileName, path.Dir(fName))
                deps := &[]string{fName}

                m.Cache.Set(filepath.Join(relPath, fileName), time.Now(), true, deps, prevData)
                fInfo, _ := os.Stat(fName)
                fTime := fInfo.ModTime()

                fileDeps := p.GetDependencies(fName)
                m.Cache.Set(fName, fTime, false, &fileDeps, nil)
                return data, ok
            } else {
                data, ok = m.recursiveProcessFile(relPath, revName, fileName, origFileName)
                if ok {
                    deps := []string{filepath.Join(relPath, revName)}
                    fileDeps := p.GetDependencies(filepath.Join(relPath, revName))
                    for _,fileDep := range fileDeps {
                        deps = append(deps, fileDep)
                    }
                    m.Cache.Set(filepath.Join(relPath, prevFileName), time.Now(), true, &deps, data)
                    return data, true
                } else if data != nil {
                    return data, false
                }
            }
        }
    }

    return nil, false
}

func (m *Manager) ProcessFile(relPath string, fileName string) ([]byte, bool) {

    var err error
    isIntermediary := false

    relPath = m.FixHiddenDirs(relPath)

    // pentru un fisier intermediar sau un hidden ( incepe cu _ )
    // in cache apare fara . in fata
    cachedFileName := fileName
    if cachedFileName[0] == '.' {
        cachedFileName = string(fileName[1:])
        isIntermediary = true
    }

    data, ok := m.Cache.Get(filepath.Join(relPath, cachedFileName))
    if ok {
        println("------ DIN CACHE -----", filepath.Join(relPath, fileName))
        return data, true
    }


    // trebuie sa verific pentru un fisier ce incepe cu .
    // sa nu fie un fisier pe disc efectiv
    data, ok = m.Cache.Get(filepath.Join(relPath, fileName))
    if ok {
        println("------  REAL FILE DIN CACHE -----", filepath.Join(relPath, fileName))
        data, ok = m.GetFileContents(relPath, fileName)
        return data, true
    }


    data, ok = m.GetFileContents(relPath, fileName)
    if ok {
        realName, _ := m.getRealFilePathFromRelPath( relPath, fileName)
        fInfo, _ := os.Stat(realName)
        fTime := fInfo.ModTime()
        m.Cache.Set(realName, fTime, false, nil, nil,)
        m.Cache.Set(filepath.Join(relPath, fileName), fTime, true, &[]string{realName}, nil)
        return data, true
    }

    if fileName[0] == '.' {
        fileName = string(fileName[1:])
    }

    data, ok = m.recursiveProcessFile(relPath, fileName, fileName, fileName)


    if ok && !isIntermediary {
        for _, p := range m.PostProcessPlugins {
            data, err = p.ProcessData( data, relPath, fileName)
            if err != nil {
                return data, false
            }

        }
        m.Cache.UpdateCacheData(filepath.Join(relPath, fileName), data)
    }

    return data, ok
}
