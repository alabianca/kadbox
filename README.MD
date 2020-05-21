# kadbox

#### Requirements
go version >= 1.14 

#### Install
`make` will install the cli under `bin/kadbox`. 
You can also install it into your `GOPATH` by navigating into cmd/kadbox and `go install`.

### Initializing
Before starting the server you need to initialize a directory as your `home` directory. Navigate to any directory and initialize it
with `kadbox init`. You can conveniently initialize our machines `$HOME` directory with
`kadobx init -g`. Initializing will create a public and private keypair and will assign you a unique identity.

 ### Starting the Server
 In the same directory you initialized OR any subdirectory you may run `kadbox server start`.
 The node attempts to connect to `bootstrap` peers and join the network.
 
 ### Bootstrap Peers
 What are bootstrap peers? These are pre configured nodes assisting new peers when they join the network. They provide the joining peer with the
 contact information of its closest neighbors and inform them about the joining node.
 Currently, a bootstrap peer also acts as a circuit relay in case two given peers can't connect directly to each other due to NATs and firewalls.
 
 The list of bootstrap peers can be found in `home/.kadbox/kadconfig`. You can actually run your own server as a bootstrap peer by running `kadbox server start -g`. The `-g` flag will configure the server as a bootstap peer.
 
 ### Sharing a file
 You can share a file in the network with `kadbox add <path to file>`. You must be in your initialized home directory or any sub directory of it.
 You know it worked if the command spits out an identifier like `kadbox://<hash_of_content>`. You can share this identifier with anyone that wants to download
 your content directly from your machine. 
 
 ### Downloading content
 Content can be downloaded using `kadbox get <identifier>`.
 Currently, it will write the content to your `pwd`.
 
 
 #### Todo
 - [ ] Re-advertise stored files when restarting a server
 - [ ] Finish Electron Desktop UI 
  