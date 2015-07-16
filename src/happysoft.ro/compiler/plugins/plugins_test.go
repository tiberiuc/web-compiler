
package plugins

import (
    "testing"
    "flag"

    "io/ioutil"

    "os"
    "path"
    "path/filepath"

)

var pluginManager *Manager
var outputPath = "../../../../testdata/output"

func init() {
    flag.Set("configs", "../../../../testdata/config.cfg")
}

func compare(X, Y []string) bool {

    m := make(map[string]int, 0)

    for _,y := range Y {
        m[y]++
    }

    for _,x := range X {
        if m[x] > 0 {
            m[x]--
            continue
        }

        return false
    }

    for _,val := range m {
        if val != 0 {
            return false
        }
    }

    return true
}

func testFolderContents(relPath string, t *testing.T) bool {

    var outputFiles   []string
    var outputFolders []string

    d, err := os.Open(path.Join(outputPath, relPath))
    if err != nil {
        t.Error("Could not open", path.Join(outputPath, relPath))
        return false
    }

    fis, err := d.Readdir(-1)
    if err != nil {
        t.Error("Could not read dir", path.Join(outputPath, relPath))
        return false
    }

    for _, fi := range fis {
        if fi.Mode().IsRegular() {
            outputFiles = append(outputFiles, fi.Name())
        }

        if fi.IsDir() {
            outputFolders = append(outputFolders, fi.Name())
        }
    }

    files := pluginManager.ListFiles(relPath)

    equal := compare(outputFiles, files)

    if !equal {
        t.Error("Files lists are not identical", relPath)
        return false
    }

    for _, file := range files {
        if ok := testFile( path.Join(relPath, file), t); !ok {
            return false
        }
    }

    folders,err := pluginManager.ListFolders(relPath)

    if err != nil {
        t.Error("Error on listing output folders", relPath, err)
        return false
    }

    equal = compare(outputFolders, folders)

    if !equal {
        t.Error("Folders lists are not identical",relPath)
        return false
    }

    for _, folder := range folders {
        if ok := testFolderContents( path.Join(relPath, folder), t); !ok {
            return false
        }
    }

    return true
}

func testFile(path string, t *testing.T) bool {

    relPath := filepath.Dir(path)
    fileName := filepath.Base(path)

    data, _ := pluginManager.ProcessFile(relPath, fileName);
    //if !ok {
        //t.Error("Could not process file ", path)
        //return false
    //}


    dataOutput,err := ioutil.ReadFile(filepath.Join(outputPath, relPath, fileName))
    if err != nil {
        t.Error("Could not read the output corespondent file ", path)
        return false
    }

    if string(data) != string(dataOutput) {
        t.Error("Output not identical ", path,"\n", string(data), "\n---------\n", string(dataOutput))
        return false
    }

    return true
}



func TestStartPluginManager(t *testing.T) {
    pluginManager = new(Manager)
    pluginManager.Init()
}


func TestCompile(t *testing.T) {
    if ok := testFolderContents("", t); !ok {
        t.Fail()
    }
}
