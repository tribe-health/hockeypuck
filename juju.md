# Deploying Hockeypuck with Juju

## Prerequisites

Juju 1.23.2 or later is recommended to make full use of this charm's [juju
actions](https://jujucharms.com/docs/stable/actions). Install juju with:

`sudo apt-add-repository ppa:juju/stable`

`sudo apt-get update`

`sudo apt-get install juju-core`

Also `apt-get install juju-local` if you'd like to use the local provider.

Familiarity with Juju and a bootstrapped environment is assumed. Read the [Juju
Documentation](https://jujucharms.com/docs/) to get started.

## Deploying Hockeypuck

Deploy a Hockeypuck service:

`juju deploy cs:~hockeypuck/trusty/hockeypuck`

Deploy MongoDB and relate it:

`juju deploy mongodb`

`juju add-relation mongodb hockeypuck`

`juju expose hockeypuck`

## Accessing your new instance

`juju status hockeypuck` will show the public address of the Hockeypuck
workload. For example:

```
$ juju status hockeypuck
environment: azure
machines:
  "15":
    agent-state: started
    agent-version: 1.23.2
    dns-name: juju-azure-dev-y9157oo521.cloudapp.net
    instance-id: juju-azure-dev-y9157oo521-jujuw0fh43evjcmace0ol1gsg0kltv5dh9b7bs8chm9gjj4gmp
    instance-state: ReadyRole
    series: trusty
    hardware: arch=amd64 cpu-cores=1 mem=1792M root-disk=130048M
services:
  hockeypuck:
    charm: local:trusty/hockeypuck-9
    exposed: true
    relations:
      mongodb:
      - mongodb
    units:
      hockeypuck/0:
        agent-state: started
        agent-version: 1.23.2
        machine: "15"
        open-ports:
        - 11370/tcp
        - 11371/tcp
        public-address: juju-azure-dev-y9157oo521.cloudapp.net
```

You should be able to access the keyserver at the listed public address,
`juju-azure-dev-y9157oo521.cloudapp.net`

## HTTP reverse-proxy

Expose Hockeypuck on port 80 behind haproxy.

`juju deploy haproxy`

`juju add-relation hockeypuck:website haproxy:reverseproxy`

`juju expose haproxy`

Or behind squid for caching.

`juju deploy squid-reverseproxy`

`juju add-relation hockeypuck:website squid-reverseproxy`

`juju set squid-reverseproxy port=11371`

`juju expose squid`

## TODOs
TODO: using the actions
TODO: peering relations

