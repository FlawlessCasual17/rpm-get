package cmds

import "fmt"

// VERSION is the current version of _rpm-get_.
const VERSION string = "0.0.1"

// Usage prints the help text when rpm-get is called without arguments,
// Or called with the `help` argument.
func Usage() {
	fmt.Printf(`
rpm-get version %s


Usage

rpm-get {update [--repos-only] [--quiet] | upgrade [--dg-only] | show <pkg list> | install <pkg list>
        | reinstall <pkg list> | remove [--remove-repo] <pkg list>
        | search [--include-unsupported] <regex> | cache | clean
        | list [--include-unsupported] [--raw|--installed|--not-installed]
        | prettylist [<repo>] | csvlist [<repo>] | fix-installed [--old-apps]
        | help | version}

rpm-get provides a high-level commandline interface for the package management
system to easily install and update packages published in 3rd party rpm
repositories or via direct download.

update
    update is used to resynchronize the package index files from their sources.
    When --repos-only is provided, only initialize and update rpm-get's
    external repositories, without updating rpm or looking for updates of
    installed packages.
    When --quiet is provided the fetching of rpm-get repository updates is done without progress feedback.

upgrade
    upgrade is used to install the newest versions of all packages currently
    installed on the system.
    When --dg-only is provided, only the packages which have been installed by rpm-get will be upgraded.

install
    install is followed by-repo is provided, also remove the rpm repository
    of rpm/ppa packages.

clean
    clean clears out the local repository (/var/cache/rpm-get) of retrieved
    package files.

search
    search for the given regex(7) term(s) from the list of available packages
    supported by rpm-get and display matches. When --include-unsupported is
    provided, include packages with unsupported architecture or upstream
    codename and include PPAs for Debian-derived distributions.

show
    show information about the given package (or a space-separated list of
    packages) including their install source and update mechanism.

list
    list the packages available via rpm-get. When no option is provided, list
    all supported packages and tell which ones are installed (slower). When
    --include-unsupported is provided, include packages with unsupported
    architecture or upstream codename and include PPAs for Debian-derived
    distributions (faster). When --raw is provided, list all packages and do
    not tell which ones are installed (faster). When --installed is provided,
    only list the packages installed (faster). When --not-installed is provided,
    only list the packages not installed (faster).

prettylist
    markdown formatted list the packages available in repo. repo defaults to
    01-main. If repo is 00-builtin or 01-main the packages from 00-builtin are
    included. Use this to update README.md.

csvlist
    csv formatted list the packages available in repo. repo defaults to
    01-main. If repo is 00-builtin or 01-main the packages from 00-builtin are
    included. Use this with 3rd party wrappers.

cache
    list the contents of the rpm-get cache (/var/cache/rpm-get).

fix-installed
    fix installed packages whose definitions were changed. When --old-apps is
provided, transition packages to new format. This command is only intended
    for internal use.

help
    show this help.

version
    show rpm-get version.
`, VERSION)
}
