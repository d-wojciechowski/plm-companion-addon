<!--- [![pipeline status](https://gitlab.com/wojciechowskid/wncplugin/badges/master/pipeline.svg)](https://gitlab.com/wojciechowskid/wncplugin-serveraddon/pipelines/latest) --->
<!--- [![report](https://goreportcard.com/badge/gitlab.com/wojciechowskid/wncplugin-serveraddon)](https://goreportcard.com/report/gitlab.com/wojciechowskid/wncplugin-serveraddon) --->

# Windchill Intellij Plugin

Server side addon for windchill to communicate with Intellij Plugin

## Getting Started

With cloned repository you can do following steps to interact with code:
````
    go build
````

to run it 
````
    go run wnc_plugin.go
````

### Prerequisites

To start you need:
* GOLang 1.13
* Protoc configured and installed [how to set up](https://developers.google.com/protocol-buffers/docs/gotutorial)

### Distribution

For windows machine there is already script [make.bat](make.bat), and for linux [make.sh](make.sh). Just run it, it will produce in distr folder executables for most common platforms.

## Authors

* **Dominik Wojciechowski** - [wojciechowskid](https://gitlab.com/wojciechowskid)
* **Michał Celniak** - [Michal1993r](https://gitlab.com/Michal1993r)

## License

This project is licensed under the GNU GENERAL PUBLIC LICENSE Version 3 - see the [LICENSE](LICENSE) file for details
