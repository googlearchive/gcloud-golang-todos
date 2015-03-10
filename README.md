## gcloud-golang-todos

> [TodoMVC](http://todomvc.com) backend using [gcloud-golang](//github.com/GoogleCloudPlatform/gcloud-golang).


### Prerequisites

1. Create a new cloud project on [console.developers.google.com](https://console.developers.google.com)
1. Export your project id:
    
    ```sh
    gcloud config set project <project id>
    ```

1. go-get this code!

    ```sh
    go get -u github.com/GoogleCloudPlatform/gcloud-golang-todos
    ```

1. Initialize the todomvc subproject:

    ```sh
    cd $GOPATH/src/github.com/GoogleCloudPlatform/gcloud-golang-todos
    git submodule init
    ```


### Running

#### [Locally via Managed VMs & Docker](https://developers.google.com/appengine/docs/managed-vms/)

```sh
# Check that Docker is running
boot2docker up
$(boot2docker shellinit)

# Run the app
gcloud preview app run main

# Open http://localhost:8080/examples/angularjs/index.html in the browser!
```

### Contributing changes

* See [CONTRIB.md](CONTRIB.md)


### Licensing

* See [LICENSE](LICENSE)
