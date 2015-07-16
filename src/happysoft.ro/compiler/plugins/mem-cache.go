package plugins

import (
    "fmt"
    "os"
    "time"
    "path"
)

type Cache struct {
    Data map[string]*FileCache
    manager *Manager
}

func NewCache(m *Manager) *Cache {
    cache := &Cache{}
    cache.Data = make(map[string]*FileCache)
    cache.manager = m

    return cache
}

func (c *Cache) Log() {
    println("========== LOG CACHE =======")
    for id, cache := range c.Data {
        haveData := "NOT empty"
        if cache.Data == nil {
            haveData = "empty"
        }
        fmt.Println(id, "( Virtual:", cache.IsVirtual, "Data:", haveData, "Timestamp:", cache.Timestamp.String(), ")")
        if cache.Depedencies != nil {
            for _, deps := range *cache.Depedencies {
                fmt.Println("\t", deps)
            }
        }
    }
    println("========== END LOG CACHE =======")
}

func (c *Cache) Set(pathStr string, timestamp time.Time, isVirtual bool, deps *[]string, data []byte) {
    c.Data[pathStr] = &FileCache{Filename: pathStr, Timestamp: timestamp, IsVirtual: isVirtual, Depedencies: deps, Data: data, LastCheck: time.Now()}
}

func checkTimeAfter(t1 time.Time, t2 time.Time) bool {
    t1 = t1.Add(time.Millisecond * 100)
    after := t1.After(t2) || t1.Equal(t2)
    return after
}

func (c *Cache) ShouldRebuild(pathStr string, timestamp time.Time) bool {
    fileCache, ok := c.Data[pathStr]
    if !ok {
        return true
    }

    // cand se apeleaza foarte repede ( un read dupa un attr ) sa nu mai refaca verificarile
    lastCheckDiff := time.Since(fileCache.LastCheck)/time.Millisecond
    if lastCheckDiff < 100 {
        return false
    }


    depTimestamp := fileCache.Timestamp

    if !fileCache.IsVirtual {
        fInfo, _ := os.Stat(fileCache.Filename)
        depTimestamp = fInfo.ModTime()

        // Pentru fisierele de pe disk reactualizez data in cache
        c.Data[pathStr].Timestamp = depTimestamp
    }


    if !(fileCache.IsVirtual && fileCache.Data == nil) {
        if ok = checkTimeAfter(timestamp, depTimestamp); !ok {
            return true
        }
    }

    if fileCache.Depedencies != nil {
        for _, dep := range *fileCache.Depedencies {
            // verific daca vreo dependinta trebuie reactualizata
            if ok = c.ShouldRebuild(dep, timestamp); ok {
                return true
            }

            // verific daca nu cumva o dependinta recursiva e deja actualizata dar e mai noua decat
            // fisierul dependinta
            if ok = checkTimeAfter(timestamp, c.Data[dep].Timestamp); !ok {
                return true
            }
        }
    }

    println(" should be ok from cache")

    return false

}

func (c *Cache) UpdateCacheData(pathStr string, data []byte) bool {
    fmt.Println(" set cache data ", pathStr, " from cache")
    _, ok := c.Data[pathStr]
    if ok {
        c.Data[pathStr].Data = data
        c.Data[pathStr].Timestamp = time.Now()
    }

    return ok
}

func (c *Cache) Get(pathStr string) ([]byte, bool) {
    fmt.Println(" get ", pathStr, " from cache")
    data, ok := c.Data[pathStr]
    if ok && !c.ShouldRebuild(pathStr, data.Timestamp) {
        data.LastCheck = time.Now()
        fileData := data.Data
        if fileData == nil {
            dir := path.Dir(pathStr)
            fileName := path.Base(pathStr)
            fileData,_ = c.manager.GetFileContents(dir, fileName)
        }
        return fileData, true
    }

    println(" ------- SHOULD REBUILD : ",pathStr, " ---------")

    return nil, false

}
