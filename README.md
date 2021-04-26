<!--- [![pipeline status](https://gitlab.com/wojciechowskid/wncplugin/badges/master/pipeline.svg)](https://gitlab.com/wojciechowskid/wncplugin-serveraddon/pipelines/latest) --->
<!--- [![report](https://goreportcard.com/badge/gitlab.com/wojciechowskid/wncplugin-serveraddon)](https://goreportcard.com/report/gitlab.com/wojciechowskid/wncplugin-serveraddon) --->

# PLM Companion Addon

Server side addon to work as a slave for Intellij IDEA Plugin (see here: [LINK](https://github.com/d-wojciechowski/plm-companion))

This program is designed to be a tunnel, and enable actions such as: disc discovery, action execution ,etc. 

### Prerequisites

To start you need:
* GOLang 1.16
* Protoc configured and installed [how to set up](https://developers.google.com/protocol-buffers/docs/gotutorial)

### Getting started

1. Clone repository
2. Execute init of submodule (about it you can learn [here](https://git-scm.com/book/en/v2/Git-Tools-Submodules))
````
git submodule update --init --recursive
````
3. Download protoc-gen-go.
````
go get -u github.com/golang/protobuf/protoc-gen-go
````
4. Execute make script, it will also generate go files based on proto.
````
make.ps1
````
### Distribution

For windows machine there is already script [make.ps1](make.ps1), and for linux [make.sh](make.sh). Just run it, it will produce in distr folder executables for most common platforms.

## Authors

* **Dominik Wojciechowski** - [d-wojciechowski](https://github.com/d-wojciechowski)
* **Michał Celniak** - [Michal1993r](https://gitlab.com/Michal1993r)

## License

This project is licensed under the GNU GENERAL PUBLIC LICENSE Version 3 - see the [LICENSE](LICENSE) file for details
