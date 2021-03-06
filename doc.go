// Package harp is a go application deploy tool (or an easy way to start go daemon or run go programs on remote servers).
//
// Please consult up-to-date README from https://github.com/bom-d-van/harp. Docs here are made for good GoDoc Searches.
//
// What Harp Does
//
// Harp simply builds your application and upload it to your server. It brings you a complete solution for deploying common applications. It syncs, restarts, kills, and deploys your applications.
//
// The best way to learn what harp does and helps is to use it. (In test directory, there are docker files and harp configurations you can play with)
//
// Usage:
//
//     # Init harp.json
//     harp init
//
//     # Configure your servers and apps in harp.json. see section Configuration below.
//     harp -s dev deploy
//
//     # Or
//     harp -s prod deploy
//
//     # Restart server
//     harp -s prod restart
//
//     # Shut down server
//     harp -s prod kill
//
//     # Inspect server info
//     harp -s prod info
//
//     # Rollback release
//     harp -s prod rollback $version-tag
//
//     # Tail server logs
//     harp -s prod log
//
//     # Done. More flags and usages are in harp -v
//
// Configuration
//
// example:
//
//     {
//         "GOOS": "linux",   // for go build
//         "GOARCH": "amd64", // for go build
//         "App": {
//             "Name":       "app",
//             "ImportPath": "github.com/bom-d-van/harp/test",
//
//             // these are included in all file Excludeds
//             "DefaultExcludeds": [".git/", "tmp/", ".DS_Store", "node_modules/"],
//             "Files":      [
//                 // files here could be a string or an object
//                 "github.com/bom-d-van/harp/test/files",
//                 {
//                     "Path": "github.com/bom-d-van/harp/test/file",
//                     "Excludeds": ["builds"]
//                 }
//             ]
//         },
//         "Servers": {
//             "prod": [{
//                 "User": "app",
//                 "Host": "192.168.59.103",
//                 "Port": ":49155"
//             }, {
//                 "User": "app",
//                 "Host": "192.168.59.104",
//                 "Port": ":49156"
//             }],
//
//             "dev": [{
//                 "User": "app",
//                 "Host": "192.168.59.102",
//                 "Port": ":49155"
//             }]
//         }
//     }
//
// How to specify server or server sets:
//
// Using the configuration above as example, server set means the key in `Servers` object value, i.e. `prod`, `dev`.
// While server is elemnt in server set arrays, you can specify it by `{User}@{Host}{Port}`.
//
//     # deploy prod servers
//     harp -s prod deploy
//
//     # deploy dev servers
//     harp -s dev deploy
//
//     # deploy only one prod server:
//     harp -server app@192.168.59.102:49155 deploy
//
// Migration
//
// Or run a go package/file on remote server.
// You can specify server or server sets on which your migration need to be executed.
//
// Simple:
//
//     harp -server app@192.168.59.103:49153 run migration.go
//
// With env and arguments:
//
//     harp -server app@192.168.59.103:49153 run "AppEnv=prod migration2.go -arg1 val1 -arg2 val2"
//
// Multiple migrations:
//
//     harp -server app@192.168.59.103:49153 run migration.go "AppEnv=prod migration2.go -arg1 val1 -arg2 val2"
//
// __Note__: Harp saved the current migration files in `$HOME/harp/{{.App.Name}}/migrations.tar.gz`. You can uncompress it and execute the binary manually if you prefer or on special occasions.
//
// Rollback
//
// By default harp will save three most recent releases in `$HOME/harp/{{.App.Name}}/releases` directory. The current release is the newest release in the releases list.
//
//     # list all releases
//     harp -s prod rollback ls
//
//     # rollback
//     harp -s prod rollback 15-06-14-11:29:14
//
// And there is also a `rollback.sh` script in `$HOME/harp/{{.App.Name}}` that you can use to rollback release.
//
// You can change how many releases you want to keep by `RollbackCount` or disable rollback by `NoRollback` in configs.
//
//     {
//         "GOOS": "linux",   // for go build
//         "GOARCH": "amd64", // for go build
//
//         "NoRollback": true,
//         "RollbackCount": 10,
//
//         "App": {
//             ...
//         },
//         ...
//     }
//
// Build Override
//
// Add `BuildCmd` option in `App` as bellow:
//
//     "App": {
//         "Name":       "app",
//         "BuildCmd":   "docker run -t -v $GOPATH:/home/app golang  /bin/sh -c 'GOPATH=/home/app /usr/local/go/bin/go     build -o path/to/app/tmp/app project/import/path'"
//     }
//
// Build override is useful doing cross compilation for cgo-involved projects, e.g. using Mac OS X building Linux binaries by docker or any other tools etc.
//
// Note: harp expects build output appears in directory `tmp/{{app name}}` where you evoke harp command (i.e. pwd).
//
// Script Override
//
// harp supports you to override its default deploy script. Add configuration like bellow:
//
//     "App": {
//         "Name":         "app",
//         "DeployScript": "path-to-your-script-template"
//     },
//
// The script could be a `text/template.Template`, into which harp pass a data as bellow:
//
//     map[string]interface{}{
//         "App":           harp.App,
//         "Server":        harp.Server,
//         "SyncFiles":     syncFilesScript,
//         "RestartServer": restartScript,
//     }
//
//     type App struct {
//         Name       string
//         ImportPath string
//         Files      []string
//
//         Args []string
//         Envs map[string]string
//
//         BuildCmd string
//
//         KillSig string
//
//         // TODO: could override default deploy script for out-of-band deploy
//         DeployScript  string
//         RestartScript string
//     }
//
//     type Server struct {
//         Envs   map[string]string
//         GoPath string
//         LogDir string
//         PIDDir string
//
//         User string
//         Host string
//         Port string
//
//         Set string
//
//         client *ssh.Client
//     }
//
// A default deploy script is:
//
//     set -e
//     {{.SyncFiles}}
//     {{.RestartServer}}
//
// Similarly, restart script could be override too. And its default template is:
//
//     set -e
//     {{.RestartServer}}
//
// You can inspect your script by evoking command: `harp -s prod inspect deploy` or `harp -s prod inspect restart`.
//
// Cross Compilation
//
// If you need to initialize cross compilation environment, harp has a simple commend to help you:
//
//     harp xc
package main
