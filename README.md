# Bindman-DNS Webhook
[![Go Report Card](https://goreportcard.com/badge/github.com/labbsr0x/bindman-dns-webhook)](https://goreportcard.com/report/github.com/labbsr0x/bindman-dns-webhook)

This repository lays out the pieces of code of a Bindman-DNS webhook.

The libraries present here should be used in order to ease out integrations among clients and managers.

# Clients and Managers

The objective behind this project is to automate the management of DNS records of a number of different DNS Server providers.

- **Clients** are pieces of software that has the knowledge of services being added and removed from your cluster. Depending of the technology adopted, they may be able to identify which hostname was attached to that service. With that information in hands, the client can then delegate the modification of a DNS Record to a Bindman DNS manager.

- **Managers** are pieces of software that receives requests from clients to modify its DNS Server records.

# Samples
Two samples are provided in the `samples` folder.

- **client**: demonstrates how one can leverage the client library to communicate with the DNS manager webhook APIs. It is configured via the `BINDMAN_DNS_MANAGER_ADDRESS` and `BINDMAN_REVERSE_PROXY_ADDRESS` environment variables which, respectively, defines the address of the manager instance and the address of the reverse proxy that will handle requests to the Sandman services.

- **hook**: demonstrates how one can leverage the hook library to receive requests modifying the DNS records it manages. 

A postman collection is provided (`samples/bindman-dns-webhook-samples.postman_collection.json`) thats lays out the available apis and how to communicate with them.

To build and run the samples just type `docker-compose up` from the samples folder.
