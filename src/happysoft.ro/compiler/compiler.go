package main

import (
    "log"
    "os"
    "os/signal"

    "happysoft.ro/compiler/fuse"
    "happysoft.ro/compiler/plugins"
    "happysoft.ro/compiler/web"
)


func main() {

    pluginManager := plugins.StartPluginManager()
    conn,ok := fuse.StartFuse(pluginManager)
    if !ok {
        return
    }

    sc := make(chan os.Signal, 1)
    signal.Notify(sc, os.Interrupt)
    go func() {
        for sig := range sc {
            log.Printf("captured %v, unmounting file system and exiting..", sig)
            fuse.UnmountFS()
            os.Exit(1)
            return
        }
    }()

    web.StartWebServer(pluginManager)
    println("Ready for action !")
    fuse.StartServing(conn)
}

