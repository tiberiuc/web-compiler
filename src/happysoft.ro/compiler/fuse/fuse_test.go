
package fuse

import (
    "testing"

    "flag"
)

func init(){
    flag.Set("configs", "../../../../testdata/config.cfg")
}

func TestConnect(t *testing.T) {
    t.Log("testing connect fuse")
    StartPluginManager()
}


func TestDisconnect(t *testing.T) {
    t.Log("testing disconnect fuse")
    UnmountFS()
}
