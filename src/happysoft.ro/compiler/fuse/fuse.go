package fuse

import (
    "fmt"
    "log"
    "os"
    "path"

    "bazil.org/fuse"
    "bazil.org/fuse/fs"
    "golang.org/x/net/context"

    "happysoft.ro/compiler/plugins"
)

var pluginManager *plugins.Manager

func StartFuse(manager *plugins.Manager) (*fuse.Conn, bool) {
    pluginManager = manager

    conn, err := fuse.Mount(pluginManager.Config.General.Mountpoint)
    if err != nil {
        log.Fatal(err)
        return nil,false
    }

    return conn,true
}

func UnmountFS() {
    fuse.Unmount(pluginManager.Config.General.Mountpoint)
    //cmd := exec.Command("fusermount", "-u", pluginManager.Config.General.Mountpoint)
    //cmd.Run()
}

func StartServing(conn *fuse.Conn) {
    fs.Serve(conn, FS{})
}

// FS implements the hello world file system.
type FS struct{}

func (FS) Root() (fs.Node, error) {
    return Dir{Path: "/"}, nil
}

// Dir implements both Node and Handle for the root directory.
type Dir struct {
    Path string
    Name string
}

func (d Dir) Attr(a *fuse.Attr) {
    a.Inode = 1
    a.Mode = os.ModeDir | 0555
}

func stringInSlice(a string, list []string) bool {
    for _, b := range list {
        if b == a {
            return true
        }
    }
    return false
}

func (d Dir) Lookup(ctx context.Context, name string) (fs.Node, error) {
    println("-=-=- FAC LOOKUP ")
    dirs, _ := pluginManager.ListFolders(d.Path)
    if stringInSlice(name, dirs) {
        return Dir{Path: d.Path + name + "/", Name: name}, nil
    }

    files := pluginManager.ListFiles(d.Path)
    if stringInSlice(name, files) {
        return File{Path: d.Path, Name: name}, nil
    }

    return nil, fuse.ENOENT
}

func (d Dir) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {
    dirDirs1 := make([]fuse.Dirent, 0)

    dirs, _ := pluginManager.ListFolders(d.Path)
    for _, d := range dirs {
        dirDirs1 = append(dirDirs1, fuse.Dirent{Inode: 1, Name: d, Type: fuse.DT_Dir})
    }

    files := pluginManager.ListFiles(d.Path)
    for _, f := range files {
        dirDirs1 = append(dirDirs1, fuse.Dirent{Inode: 2, Name: f, Type: fuse.DT_File})
    }
    return dirDirs1, nil
}

// File implements both Node and Handle for the hello file.
type File struct {
    Path        string
    Name        string
}

func (f File) Attr(a *fuse.Attr) {
    println(" fac attr")
    data, _ := pluginManager.ProcessFile(f.Path, f.Name)

    length := uint64(len(data))
    a.Mode = 0444
    a.Size = length
}

func (f File) ReadAll(ctx context.Context) ([]byte, error) {
    println("Fac readall")
    data, ok := pluginManager.ProcessFile(f.Path, f.Name)
    if !ok {
        fmt.Fprintf(os.Stderr, "------- ERROR on %s ------\n", path.Join(f.Path, f.Name))
        fmt.Fprintf(os.Stderr, "%s", string(data))
    }
    pluginManager.Cache.Log()
    return data, nil
}
