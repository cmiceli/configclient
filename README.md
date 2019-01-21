# configclient
The client for cmiceli/configserver, supporting subscriptions to files.

The intention of this is to provide a simple client which with very little work can operate as a daemon on a host, subscribing to configs which it will refresh periodically.

This can then mean that applications and other programs can get their configs dynamically, or have them persisted to disk
