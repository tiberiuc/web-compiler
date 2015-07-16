package plugins

import (
    "time"
)

type FileCache struct {
    Filename    string
    Timestamp   time.Time
    IsVirtual   bool
    Depedencies *[]string
    Data        []byte
    LastCheck   time.Time
}


