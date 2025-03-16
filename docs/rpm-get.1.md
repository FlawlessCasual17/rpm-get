---
author: Martin Wimpress
date: February 14, 2022
footer: rpm-get
header: rpm-get User Manual
section: 1
title: rpm-get
---

# NAME

rpm-get - An installation manager for 3rd-Party rpm packages

# SYNOPSIS

**rpm-get** \[*COMMAND*\]...

```
rpm-get {update [--repos-only] [--quiet] | upgrade | show <pkg list> | install <pkg list>
        | reinstall <pkg list> | remove [--remove-repo] <pkg list>
        | purge [--remove-repo] <pkg list>
        | search [--include-unsupported] <regex> | cache | clean
        | list [--include-unsupported] [--raw|--installed|--not-installed]
        | prettylist [<repo>] | csvlist [<repo>] | fix-installed [--old-apps]
        | help | version}
```

# DESCRIPTION

rpm-get provides a high-level commandline interface for the package management
system to easily install and update packages published in 3rd party rpm
repositories or via direct download.

# COMMANDS

**update**
:    update is used to resynchronize the package index files from their sources. When --repos-only is provided, only initialize and update rpm-get's external repositories, without updating rpm or looking for updates of installed packages.

**upgrade**
:    upgrade is used to install the newest versions of all packages currently installed on the system.

**install**
:    install is followed by one package (or a space-separated list of packages) desired for installation or upgrading.

**reinstall**
:    reinstall is followed by one package (or a space-separated list of packages) desired for reinstallation.

**remove**
:    remove is identical to install except that packages are removed instead of installed. When --remove-repo is provided, also remove the rpm repository of rpm/ppa packages.

**purge**
:    purge is identical to remove except that packages are removed and purged (any configuration files are deleted too). When --remove-repo is provided, also remove the rpm repository of rpm/ppa packages.

**clean**
:    clean clears out the local repository (/var/cache/rpm-get) of retrieved package files.

**search**
:    search for the given regex(7) term(s) from the list of available packages supported by rpm-get and display matches. When --include-unsupported is provided, include packages with unsupported architecture or upstream codename and include PPAs for Debian-derived distributions.

**show**
:    show information about the given package (or a space-separated list of packages) including their install source and update mechanism.

**list**
:    list the packages available via rpm-get. When no option is provided, list all supported packages and tell which ones are installed (slower). When --include-unsupported is provided, include packages with unsupported architecture or upstream codename and include PPAs for Debian-derived distributions (faster). When --raw is provided, list all packages and do not tell which ones are installed (faster). When --installed is provided, only list the packages installed (faster). When --not-installed is provided, only list the packages not installed (faster).

**prettylist**
:    markdown formatted list the packages available in repo. repo defaults to 01-main. If repo is 00-builtin or 01-main the packages from 00-builtin are included. Use this to update README.md.

**csvlist**
:    csv formatted list the packages available in repo. repo defaults to 01-main. If repo is 00-builtin or 01-main the packages from 00-builtin are included. Use this with 3rd party wrappers.

**cache**
:    list the contents of the rpm-get cache (/var/cache/rpm-get).

**fix-installed**
:    fix installed packages whose definitions were changed. When --old-apps is provided, transition packages to new format. This command is only intended for internal use.

**help**
:    show this help.

**version**
:    show rpm-get version.

# EXAMPLES

## Introduction

`rpm-get` makes it easy to install and update `.rpms` published in 3rd party
rpm repositories or made available via direct download on websites or GitHub
release pages.

## Install

Use `rpm-get` to install `rpm-get`.

```bash
sudo rpm install curl
curl -sL https://raw.githubusercontent.com/wimpysworld/rpm-get/main/rpm-get | sudo -E bash -s install rpm-get
```

Alternatively, you can [download the `.rpm` of `rpm-get` from the releases page](https://github.com/wimpysworld/rpm-get/releases/latest)
and install it manually with

``` bash
sudo rpm-get install ./path/to/rpm-get_<version>.rpm
```

## Usage

Here's an example of how to install Google Chrome.

```bash
rpm-get install google-chrome-stable
```

You can see what applications are supported by using `rpm-get list` or you
can search the available applications with `rpm-get search <app>`.

You can upgrade packages installed using `rpm-get` by running
`rpm-get upgrade`.

### Why?

3rd party rpm repositories exist. They are not going away.

Some application vendors and projects express their support for Debian/Ubuntu
by publishing `.rpms` of their software as direct downloads or via
their own rpm repositories. `rpm-get` makes it easy to find, install and update
`.rpms` published in this way.

  - Perhaps you want to use software that is not (yet) officially packaged for Debian/Ubuntu.
  - Perhaps you want to use software that is fast moving and newer versions are available from the vendor/project.
  - Perhaps you want to use some non-free software that Debian/Ubuntu cannot distribute due to licensing restrictions.

`rpm-get` tries to remedy this by providing a curated index of software
available for Ubuntu that is published by the project or vendor. `rpm-get` is
inspired by Software Boutique, a graphical application I co-authored that does
much the same thing and has been included in [Ubuntu MATE](https://ubuntu-mate.org)
since 2015.

## Supported Software

The list of the software that can be installed, updated and removed using `rpm-get`'s main repository is available at [01-main/README](https://github.com/wimpysworld/rpm-get/blob/main/01-main/README.md).

### How do package updates work?

#### 3rd party rpm repositories and PPAs

If packages are available via a 3rd party `rpm` repository or a Launchpad PPA,
then those packages will be updated/upgraded when using `sudo rpm-get update`
and `sudo rpm-get upgrade`.

#### GitHub Releases and direct downloads

For `.rpm` packages that are only available via GitHub Releases or direct
download, then those packages can only be updated/upgraded by using
`rpm-get update` and `rpm-get upgrade`.

#### GitHub API Rate Limits

`rpm-get` uses the [GitHub REST API](https://docs.github.com/en/rest) for some functionality when applications are provided via  GitHub Releases
and for unauthenticated interactions this API is [rate-limited](https://docs.github.com/en/rest/overview/resources-in-the-rest-api#rate-limiting) to 60 calls per hour per source (IP Address). This is vital for keeping the API responsive and available to all users, but can be inconvenient if you have a lot of GitHub releases being handled by `rpm-get` (or need to update several times in a short period to test your [contribution](https://github.com/wimpysworld/rpm-get/blob/main/01-main/CONTRIBUTING.md)) and will result in, for example, temporary failures to be able to upgrade or install applications via  GitHub Releases.

If you have a GitHub account you can authenticate your GitHub API usage to increase your rate-limit to 5000 requests per hour per authenticated user.  To do this you will need to use a [Personal Access Token (PAT)](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token). Once you have created a token within GitHub (or identified an appropriate existing token) you should insert it into an environment variable (`RPMGET_TOKEN`) for `rpm-get` to pick up and use to authenticate to the GitHub API.

e.g.:

```
export RPMGET_TOKEN=<my-secret-token>
rpm-get update
rpm-get upgrade
```

## Adding Software

For information on what is acceptable as suggestion for new packages and instructions on how to open a PR to add a new package to the main repository, head to [01-main/CONTRIBUTING](https://github.com/wimpysworld/rpm-get/blob/main/01-main/CONTRIBUTING.md).

### Adding external repositories

It is possible to also add a `rpm-get`-compatible external repository, and supplement the list of supported packages, typically because you need to:

1. Add something which does not meet any of the general guidelines of the main repository; or
2. Change the definition of a package from the main repository.

For information on how to create and maintain a `rpm-get`-compatible external repository, head to [EXTREPO](https://github.com/wimpysworld/rpm-get/blob/main/EXTREPO.md).

How to use:

* Manually create the file `/etc/rpm-get/<priority>-<repo>.repo`, containing in its first line the base URL of the repository.
    * The `<priority>` value is a two-digit number between 00 and 99 that defines the order in which the repositories will be loaded (00 first, 99 last), so if any conflicting definitions are found, the one from the repository with the highest priority will be used (the builtin definitions from the `rpm-get` script itself have priority 00, the main repository has priority 01 and the custom user includes have priority 99).
    * The `<repo>` value can be anything, but it should preferably be unique and easy to remember.
* Run `rpm-get update`, so the manifest file and the package definition files are downloaded.

### Custom User Includes

As a more advanced feature, it is possible to also add your own local customizations or overrides, and supplement the list of packages supported by the main repository. This feature is especially useful so that your local copy of the main repository can remain unmodified and always be kept fully up to date by moving your customizations out in a seperate folder away from the main repository.

Typically because:

1. You are waiting on a pending request for a new software package, which has been submitted for consideration but has not been reviewed / accepted / merged / released yet; or
2. You need to add something which does not meet any of the general guidelines of the main repository, for whatever various reason(s).

How to use:

* Manually create the folder `/etc/rpm-get/99-local.d/` if it does not exist already. By default, `rpm-get` does not create this folder unless your specific distribution has packaged it that way.
* Any files directly within this folder will be bash sourced e.g. `. /etc/rpm-get/99-local.d/appname1`.
* The name of the added file **must** match **exactly** the name of the package being defined.
* For information on how to create a package definition file, head to [EXTREPO](https://github.com/wimpysworld/rpm-get/blob/main/EXTREPO.md#the-package-definition-files).
* Your user custom package definition files are then loaded after the package definitions from any added repository.
* A recommendation message is printed for any new user added definitions, with a URL link to open a request.
* Warning messages are then also printed for any conflicts detected for overriden definitions (of same name), which then take priority over existing ones.

For the last situation, this is most often meant as a helpful reminder to remove your custom definition once it has been successfully merged upstream into the main repository, so after the main repository updates itself you are properly notified. It also avoids keeping lots of duplicate definitions around.

We really hope that you will enjoy the convenience and flexibility of the user overrides feature, so please consider in return to open new issues and pull requests, for any new package definitions you create, so that we can share those back with the wider community. Many thanks for your consideration!

## Related projects

  * [App Outlet](https://app-outlet.github.io/): *A Universal linux app store*
  * [bin-get](https://github.com/OhMyMndy/bin-get): *Script to easily and safely fetch binaries from Github Releases/tags*
  * [makerpm](https://www.makerpm.org/): *A simplicity-focused packaging tool for Debian archives*
  * [MakeResolveDeb](https://www.danieltufvesson.com/makeresolverpm): *Install DaVinci Resolve or DaVinci Resolve Studio on Debian*
  * [pacstall](https://pacstall.dev/): *The AUR alternative for Ubuntu*
  * [Ubuntu Make](https://github.com/ubuntu/ubuntu-make): *Easy setup of common tools for developers on Ubuntu.*
  * [unsnap](https://github.com/popey/unsnap): *Quickly migrate from using snap packages to flatpaks*
## In the media

  * [The rpm-get tool helps Ubuntu (and derivative distro) fans grab extra apps](https://www.gamingonlinux.com/2022/05/the-rpm-get-tool-helps-ubuntu-and-derivative-distro-fans-grab-extra-apps/) - **GamingOnLinux**
  * [Deb-Get is ‘Apt-Get’ for 3rd-Party Ubuntu Software](https://www.omgubuntu.co.uk/2022/05/use-rpm-get-to-install-popular-linux-apps-on-ubuntu) - **omg! ubuntu!**
  * [How To Install Third-Party Linux Applications with Deb-Get](https://www.tomshardware.com/uk/how-to/install-linux-apps-with-rpm-get) - **toms HARDWARE**
  * [How to use rpm-get to install third-party software on Ubuntu](https://www.techrepublic.com/article/use-rpm-get-install-third-party-software-ubuntu/) - **TechRepublic**
  * [Ubuntu MATE’s Lead Creates a Nifty Tool to Help Install 3rd Party Deb Packages](https://news.itsfoss.com/rpm-get-ubuntu/) - **It's FOSS**
  * [Easily Get 3rd-Party Software on Ubuntu With Deb-Get](https://fosspost.org/rpm-get-install-3rd-party-software-ubuntu/) - **FOSS Post**
  * [rpm-get, an "rpm-get" to install third-party software on Ubuntu](https://ubunlog.com/en/rpm-get-un-rpm-get-para-instalar-software-de-terceros-en-ubuntu/) - **ubunlog**
  * [rpm-get, a utility for installing third-party software](https://www.linuxadictos.com/en/rpm-get-una-utilidad-para-instalar-software-de-terceros.html) - **Linux Addicted**

# BUGS

Submit bug reports online at:
<https://github.com/wimpysworld/rpm-get/issues>

# SEE ALSO

Full sources at: <https://github.com/wimpysworld/rpm-get>

rpmorah(1)
