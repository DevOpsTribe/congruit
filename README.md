## Congruit

![alt text](https://github.com/lucky-sideburn/congruit/blob/master/img/logo.png "Tux1")
![alt text](https://github.com/lucky-sideburn/congruit/blob/master/img/logo.png "Tux2")

### The configuration management tool that loves Bash
#### Simple, lightweight and fully customizable by you

### Table of Contents

**[Description](#description)**

**[Concepts](#concepts)**

**[Stockroom](#stockroom)**

**[Place](#place)**

**[Work](#work)**

**[Workplace](#workplace)**

**[Build your workplaces](#build-your-workplaces)**

**[Prerequisites](#prerequisites)**

**[Try Congruit With Vagrant](#try-congruit-with-vagrant)**

**[Usage](#usage)**

## Description
Congruit is a lightweight configuration management and automation tool. It is written in Go but works through Bash. It manages shell scripts you created to configure your Linux platforms.

## Concepts
The main concepts of Congruit are

* Stockroom repository
* Works
* Places
* Workplaces

## Stockroom
The Stockroom is the main repository that describes your platform. Congruit reads the stockroom and does things.

## Place
A place is a shell script that must return 0. You should be in a right place to do a work.
Example:

Is this Linux server a Centos 7?

```
[ ! -e /etc/redhat-release ] && exit 1
cat /etc/redhat-release | grep "Centos Linux release 7.*"
```

## Work
Work is a shell script that installs and configures programs or runs Docker containers like in the following example:

```
docker run --rm -p 8888:8080 tomcat:latest &> /dev/null &
```

## Workplace
Workplaces are the union between works and places and are JSON file.

Example:

```
[
  {
   "places": ["debian","screen_is_not_installed"],
   "works": ["screen_package_apt"]
  },
  {
   "places": ["centos7","screen_is_not_installed"],
   "works": ["screen_package_yum"]
  }
]
```
the workplace is able to decide which is the correct strategy to install software.
Congruit executes places and, if they return 0, it does works.

## Build you workplace
1. Create your stockroom. I would like create a public repository with common and useful workplaces. For now you can take a look at https://github.com/lucky-sideburn/congruit/tree/master/stockroom

2. You need to describe your places. Example:
  * is this server running a specific Linux distribution?
  * are there particular configuration files, installed software, environment variables that describe the role or the functionality of this server?
  * places are executed before works... You can copy files which contain environment variables that can be used by works.
  Put places in stockroom/places/ folder

3. Create works. Put the scripts in stockroom/works. Works install software, get configuration file from a repository, manage Docker containers ecc..

4. Create workplaces in stockroom/workplaces. Try to make them usable in more environments and follow this example:

Workplaces are array of hashes
```
[
  {
   "places": ["is_linux", "is_frontend","is_apache_not_istalled"],
   "works": ["install_apache"]
  },

  {
   "places": ["is_linux", "is_frontend", "has_additionl_vhost"],
   "works": ["additional_vhost","restart_apache"]
  },

  {
   "places": ["is_production"],
   "works": ["do_backup"]
  }

]
```

## Try Congruit With Vagrant

* List all virtual machine in current Vagrant project

```
eugenio@local:[~/WORK/GO/src/congruit]: vagrant status
Current machine states:

Centos7                   running (virtualbox)
```

* Provision and test your workplaces

```
export WORKPLACES_ENABLED=install_screen
vagrant provision Centos7
```

Remember to set WORKPLACES_ENABLED (example: WORKPLACES_ENABLED=do_this,do_this2,do_foobar) in order to execute workplaces

## Prerequisites
1. GO
2. a stockroom. You can take as example the stockroom present in this repo. Please, create symlink from stockroom/workplaces/foo to stockroom/workplaces_enabled/foo if you want apply the workplace "foo" during congruit execution


## Usage
1. `git clone https://github.com/lucky-sideburn/congruit.git`
2. `go build conguit.go`

3. Start Congruit

`./congruit  -stockroom-dir=./stockroom`

`./congruit  -stockroom-dir=./stockroom -debug`

