# DFI

## DEAD
This project, after some advice, is no longer being actively worked on. The git history has been removed for a similar reason.

IRC: `#dfi` on Freenode

[Whitepaper](https://github.com/wjh/whitepaper)

A distributed file sharing and indexing network.

## An old UI prototype, built on Electron

![](http://i.giphy.com/26FmQ8QfY5zOglgGY.gif)


## What even is this?

- a queryable database of torrents and metadata
- a p2p network for the distribution of torrent info hashes and metadata

Peer discovery is done via an implementation of Kademlia, a DHT similar to the one Bittorrent uses. Peers are assigned addresses similar to Bitcoin addresses, except using Ed25519 and SHA3.

## What does it do?
DFI allows users to discover file indexes. Searches can then be performed on the remote peer for files that they are indexing, and file data is sent as a response. This includes a BitTorrent info-hash, which can then be used to download the file.

DFI also allows users to mirror the index of a peer. This massively enhances search speed, and allows anyone to take a complete backup of an index.

DFI can also be routed through any SOCKS proxy, and can create a Tor onion address automatically - this aids privacy and traverses the NAT, at the cost of performance.

## Sounds cool, when can I use it?

Now! It will likely have some bugs, but is mostly working.

You can download DFI [here](https://github.com/wjh/dfi/releases)

## Building

### MacOS, Linux, *nix
Make sure that you have `make` installed. You will also need `git`, and, of course, `go`.

```
git clone https://github.com/dfi/dfi
cd dfi
make
```

The resulting binary will be in `dfi/bin`.

### Windows
The Makefile doesn't seem to work so well with Windows, so you'll need to have Go properly setup and installed, this may require setting $GOPATH
```
go get github.com/dfi/dfi
```

The resulting "dfid" binary should be automatically installed into your $GOPATH.
## Usage

DFI presents a HTTP interface for usage, so you can interact with it using `curl`. There is a command program available [here](https://gitlab.com/PoroCYon/siv), which makes interaction easier. siv is also downloadable from the [release](https://github.com/dfi/dfi/releases) page.
### Connecting to the main network
Presently I am running several DFI nodes, you can bootstrap from this network as follows (assuming you have dfid running and listening on port 8080).

First you need Tor running: to do this, cd into the `tor` directory, and run `tor -f torrc`. You will also need to edit `dfid.toml` and change `tor.enabled` and `socks.enabled` to be true. The permissions on the DFI directory may well need to be `700` for Tor to be happy.

To get started, simply run dfid. The output will contain your DFI address, which will look something like this: `ZncGWimPZHWxjTMj51QNKg25PTCXphtLbh`

In order to connect to the rest of the network, you will need to bootstrap. This can either be done using the below API, or using [siv](https://gitlab.com/PoroCYon/siv).

```
curl localhost:8080/self/bootstrap/x4yknq5x7iijrmgy.onion/
curl localhost:8080/self/explore/
```

### API

By default, DFI listens on `localhost:8080`. This is configurable in `dfid.toml`. 

#### self

These routes affect the local peer, ie the client running on your machine. They're generally used to interact with your own database, or change settings, etc.

##### `/self/addpost/` POST
This is used to add a post to your database, a post is essentially a torrent infohash and some metadata. The POST body requires a parameter of `data` and `index`.

The former is JSON, and is specified as such:
```
InfoHash   string - the torrent infohash
Title      string - a name for the post
Size       int    - the size of the torrent in bytes 
FileCount  int    - number of files the torrent has
Seeders    int    - number of seeders the torrent has
Leechers   int    - number of leechers the torrent has
UploadDate int    - Unix timestamp in seconds
Tags       string - a comma-separated list of alphanumeric tags
Meta       string - a JSON-encoded object
```

The other parameter, `index`, should be either "true" or "false". This indicates whether or not DFI should add the post to the full text search index. If this is true, then the `Title` field will be indexed and the post will show up in search results.

##### `/self/index/{since}/` GET
This performs a full text search index on all posts that have an id greater than `{since}`.

##### `/self/resolve/{address}` GET
This resolves a DFI address into a JSON entry. Entries are specified as such:

```
address        Address 
name           string  
desc           string  
publicAddress  string  
publicKey      []byte  
postCount      int     
updated        uint64  
signature      []byte 
collectionHash []byte 
port           int   
seeds          [][]byte 
seeding        [][]byte 
seen           int      
```

##### `/self/bootstrap/{address}/` GET
Bootstraps the DFI node from the given address. This address must be a non-dfi address - for instance, a domain name, IP address, onion address, or anything else. Note that dfi can be configured to use a SOCKS proxy, see dfid.toml.

##### `/self/search/` POST
Perform a full text search on the local database.

This takes the parameters of `query` and `page`, where query is the search term and page is the page of results we want - this starts at 0.

##### `/self/recent/{page}/` GET
Gets the most recent posts. The page is given as the `{page}` parameter.

##### `/self/popular/{page}/` GET
Gets the most popular posts. The page is given as the `{page}` parameter.

##### `/self/peers/` GET
Returns a list of peers.

##### `/self/explore/` GET
Begin network exploration. This should happen automatically at start if you have peers in your routing table, otherwise it needs to be ran manually.

##### `/self/set/{name}/` POST
This is used to set various settings for the node. Here are possible values for `{name}`:
- name: This sets the name field of the entry and can be used to identify your node
- desc: A short description of your entry
- public: set the public address for your entry, this is what it's DFI address will resolve to


##### `/self/get/{name}/` GET
Gets values, much like set - supports the same values. Also supports:
- dfi: gets the DFI address
- postcount: the number of posts this node has
- entry: the full DFI DHT entry

#### peer
These routes allow you to query a remote peer. The `{address}` parameter refers to the encoded DFI address of the peer, like the example above.

##### `/peer/{address}/ping/`
Pings the peer.

##### `/peer/{address}/announce/`
Announces your own entry to the peer.

##### `/peer/{address}/rsearch/`
Performs a remote search on the peer.

##### `/peer/{address}/mirror/`
Download a local copy of the peer's post database, which can then be indexed and searched.

##### `/peer/{address}/search/`
Search the local copy of the peer's database, this only works after a successful `mirror`.

##### `/peer/{address}/recent/{page}/`
Get the `{page}` of most recent posts for the given peer.

##### `/peer/{address}/popular/{page}/`
Performs a remote search on the peer.

##### `/peer/{address}/index/{since}/`
Add all posts with an id larger than `{since}` to the FTS index.
