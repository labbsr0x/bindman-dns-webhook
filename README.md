# Sandman DNS Webhook
This repository lays out the pieces of code of a sandman DNS manager webhook.

The libraries present here should be used in order to ease out integrations among listeners and managers.

# Listeners and Managers

The objective behind this project is to automate the management of DNS records of a number of DNS Server that are sitting above a Sandman cluster.

- **Listeners** are pieces of software that listens to events a specific Sandman cluster emits, filtering for events of new services being added. A Sandman cluster expects these events to be annotated with labels identifying which DNS hostname should be assigned to that a sandman cluster service. Once identified these labels, a request is sent to its registered DNS manager.

- **Managers** are pieces of software that receives requests from listeners to modify its DNS Server records. 

# Samples
Two samples are provided in the `samples` folder.

- **client**: demonstrates how one can leverage the client library to communicate with the DNS manager webhook APIs. It is configured via the `MANAGER_ADDRESS` and `REVERSE_PROXY_ADDRESS` environment variables, which respectively, defines the address of the manager instance and the address of the reverse proxy that will handle requests to the services the listener has identified as manageable by this DNS manager

- **hook**: demonstrates how one can leverage the hook library to receive requests modifying the DNS records it manages. It expects the `TAGS` environment variable to be defined.

A postman collection is provided (`samples/sandman-dns-webhook-samples.postman_collection.json`) thats lays out the available apis and how to communicate with them.

# Tags

When a Sandman cluster runs a new service, it is expected from this service to be annotaded with labels indicating what the hostname for that service should be.

More than that, it should also be labeled with information of its exposure. In other words, which DNS Server should handle DNS queries to that hostname the service is expected to be attached to.

Sandman accomplishes that by adopting a Tag system, where each service launched by a Sandman cluster is annotated with Tags and each DNS manager is also run with its own Tags.

**The intersection of these two Tags indicates to the sandman which DNS Managers will manage which service.**