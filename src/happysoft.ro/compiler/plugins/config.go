package plugins

import (
    "flag"
    "fmt"
    "os"
    "path/filepath"

    "code.google.com/p/gcfg"
)


var PluginsConfig *Config

func init() {


    PluginsConfig = new(Config)

    c := PluginsConfig

    flag.Usage = Usage

    flag.BoolVar(&c.General.Development, "development", true, "Set development mode ( default True )")
    flag.StringVar(&c.General.Path, "path", "", "Set root paths ( for multiple path separate with :)")
    flag.StringVar(&c.General.Mountpoint, "mountpoint", "", "Set mount point")
    flag.StringVar(&c.General.Configs, "configs", "", "Load config from files ( for multiple config files, separate with : )")

    flag.BoolVar(&c.Sass.Disabled, "sass.disabled", false, "Disable Sass plugin (default False)")

    flag.BoolVar(&c.Scss.Disabled, "scss.disabled", false, "Disable Scss plugin (default False)")

    flag.BoolVar(&c.Stylus.Disabled, "stylus.disabled", false, "Disable Stylus plugin  (default False)")
    flag.BoolVar(&c.Coffeescript.Disabled, "coffeescript.disabled", false, "Disable Coffeescript plugin (default False)")

    flag.BoolVar(&c.Jsmin.Disabled, "jsmin.disabled", false, "Disable JS Minifier plugin (default False)")
    flag.BoolVar(&c.Jsmin.Removelogs, "jsmin.removelogs", true, "In production mode delete console.log from JS before minify (default true)")

    flag.BoolVar(&c.Haml.Disabled, "haml.disabled", false, "Disable Haml plugin  (default False)")
    flag.BoolVar(&c.Jade.Disabled, "jade.disabled", false, "Disable Jade plugin  (default False)")

    flag.BoolVar(&c.Concatenate.Disabled, "concatenate.disabled", false, "Disable file Concatenate plugin (default False)")

    flag.BoolVar(&c.Configvars.Disabled, "configvars.disabled", false, "Disable file Configvars plugin (default False)")
    flag.StringVar(&c.Configvars.Files, "configvars.files", "", "Load vars from files ( for multiple files, separate with : )")
    flag.StringVar(&c.Configvars.Include, "configvars.include", "", "Comma separated list with extensions to be included")
    flag.StringVar(&c.Configvars.Exclude, "configvars.exclude", "", "Comma separated list with extensions to be excluded ( Include have priority in case is defined )")
    flag.Var(&c.Configvars.addedvars, "configvars.addvar", "Add new var to the list")

}

type addedVars []string

type Config struct {
    General struct {
        Development bool
        Path        string
        Mountpoint  string
        Configs     string
    }

    Sass struct {
        Disabled bool
    }

    Scss struct {
        Disabled bool
    }

    Coffeescript struct {
        Disabled bool
    }

    Jsmin struct {
        Disabled   bool
        Removelogs bool
    }

    Haml struct {
        Disabled bool
    }

    Jade struct {
        Disabled bool
    }

    Stylus struct {
        Disabled bool
    }

    Concatenate struct {
        Disabled bool
    }

    Configvars struct {
        Disabled  bool
        Files     string
        Include   string
        Exclude   string
        addedvars addedVars
    }
}

var Usage = func() {
    fmt.Fprintf(os.Stderr, "Usage of %s:\n", "compiler")
    fmt.Fprintf(os.Stderr, "  %s -mountpoint=MOUNTPOINT | -configs=CONFIG_FILES\n\n", "compiler")
    flag.PrintDefaults()
}

func (v *addedVars) String() string {
    vars := ""
    for _, a := range *v {
        vars += a + "\n"
    }

    return vars
}

func (v *addedVars) Set(value string) error {
    *v = append(*v, value)
    return nil
}

func (c *Config) getConfig() error {
    flag.Parse()

    if c.General.Configs != "" {
        configs := filepath.SplitList(c.General.Configs)
        for _, config := range configs {
            err := gcfg.ReadFileInto(c, config)
            if err != nil {
                fmt.Fprintf(os.Stderr, "Error loading config: %s\n\n%s\n", config, err.Error())
                os.Exit(2)
            }
        }
    }

    if c.General.Path == "" {
        fmt.Fprintf(os.Stderr, "Missing path\n")
        flag.Usage()
        os.Exit(2)
    }
    if c.General.Mountpoint == "" {
        fmt.Fprintf(os.Stderr, "Missing mount point\n")
        flag.Usage()
        os.Exit(2)
    }

    return nil
}
