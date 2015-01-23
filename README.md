## gcloud-golang-todos

> [TodoMVC](http://todomvc.com) backend using [gcloud-golang](//github.com/GoogleCloudPlatform/gcloud-golang).

[![Build Status](https://travis-ci.org/GoogleCloudPlatform/gcloud-golang-todos.svg?branch=master)](https://travis-ci.org/GoogleCloudPlatform/gcloud-golang-todos)


### Prerequisites

1. Create a new cloud project on [console.developers.google.com](http://console.developers.google.com)
2. [Enable](https://console.developers.google.com/flows/enableapi?apiid=datastore) the [Google Cloud Datastore API](https://developers.google.com/datastore)
3. Create a new service account and copy the JSON credentials to `key.json`
4. Export your project id:

    ```sh
    $ export PROJECT_ID=<project id>
    ```
5. Initialize the todomvc subproject:

	git submodule init


### Running


#### Locally

```sh
$ gcloud components update gae-go
$ go get google.golang.org/appengine
$ boot2docker up
$
```


#### [Docker](https://docker.com)


#### [Managed VMs](https://developers.google.com/appengine/docs/managed-vms/)


### Resources


### Contributing changes

* See [CONTRIB.md](CONTRIB.md)


### Licensing

* See [LICENSE](LICENSE)
