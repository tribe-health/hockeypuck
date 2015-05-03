# Install on Ubuntu Server

# Prerequisites

## Ubuntu 14.04
The latest LTS is recommended. Hockeypuck 2.0 is currently packaged for trusty.

# Add the unstable Hockeypuck PPA

```
sudo apt-add-repository ppa:hockeypuck/unstable
sudo apt-get update
```

# Install the database of your choice
If you plan on connecting to a local database on the same server, install it now:

MongoDB:
`sudo apt-get install mongodb-server`

PostgreSQL 9.4:
You'll need to get it from the Postgres [Apt repository for LTS distributions](http://www.postgresql.org/download/linux/ubuntu/).

Otherwise, skip and install hockeypuck. You'll need to `configure </configuration.html>`_ your Hockeypuck instance to connect to your PostgreSQL database.

# Install Hockeypuck

`sudo apt-get install hockeypuck`

# Next steps

* [Configure](configuration.md) the Hockeypuck server.

