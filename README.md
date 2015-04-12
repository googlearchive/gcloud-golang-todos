## gcloud-golang-todos

> [TodoMVC](http://todomvc.com) backend using [gcloud-golang](//github.com/GoogleCloudPlatform/gcloud-golang).


### Prerequisites

1. Set up a [Go installation](https://golang.org/doc/install) and [workspace](https://golang.org/doc/code.html).
1. Install the [Cloud SDK](https://cloud.google.com/sdk/). If necessary, running the following will install
the Developer Preview commands and the App Engine SDK for Go.

    ```sh
    gcloud components update preview gae-go
    ```

1. Create a new cloud project on [console.developers.google.com](https://console.developers.google.com).
1. Export your project id:
    
    ```sh
    gcloud config set project <project id>
    ```

1. Clone the repository by running the following command:

    ```sh
    go get -u github.com/GoogleCloudPlatform/gcloud-golang-todos
    ```

1. Initialize the `todomvc` submodule. Since TodoMVC is linked within this repository as a git submodule, we need
to fetch its codebase separately:

    ```sh
    cd $GOPATH/src/github.com/GoogleCloudPlatform/gcloud-golang-todos
    git submodule update --init
    ```


### Running

#### [Locally via Managed VMs & Docker](https://developers.google.com/appengine/docs/managed-vms/)

```sh
# Check that Docker is running.
boot2docker up
$(boot2docker shellinit)

# Download the Docker runtime images for Managed VMs; make sure to select the Go runtime.
gcloud preview app setup-managed-vms

# Run the app.
gcloud preview app run main

# Open http://localhost:8080/examples/angularjs/index.html in the browser!
```

### Todo

* Determine a reasonable testing strategy. Either wait for an aetest port, or develop something more involved.

### Contributing changes

* See [CONTRIB.md](CONTRIB.md)


### Licensing

* See [LICENSE](LICENSE)
