---
services: service-bus
platforms: go
author: mcardosos
---

# Azure Service Bus Management Sample using Azure SDK for Go

This sample demonstrates how to manage Azure Service Bus using Go, and specifically how to:

- Create a namespace
- Create an authorization rule for the namespace
- List keys and connection strings for the namespace
- Create a queue
- Create a topic
- Create a subscription

If you don't have a Microsoft Azure subscription you can get a FREE trial account [here](https://azure.microsoft.com/pricing/free-trial).

**On this page**

- [Run this sample](#run)
- [More information](#info)

<a id="run"></a>

## Run this sample

1. If you don't already have it, [install Go](https://golang.org/dl/).

1. Clone the repository.

    ```
    git clone https://github.com:Azure-Samples/service-bus-go-management.git
    ```

    or simply `go get` it:

     ```
    go get github.com/Azure-Samples/service-bus-go-management
    ```

1. Install the dependencies using glide.

    ```
    cd service-bus-go-management
    glide install
    ```

1. Create an Azure service principal either through
    [Azure CLI](https://azure.microsoft.com/documentation/articles/resource-group-authenticate-service-principal-cli/),
    [PowerShell](https://azure.microsoft.com/documentation/articles/resource-group-authenticate-service-principal/)
    or [the portal](https://azure.microsoft.com/documentation/articles/resource-group-create-service-principal-portal/).

1. Set the following environment variables using the information from the service principle that you created.

    ```
    export AZURE_TENANT_ID={your tenant id}
    export AZURE_CLIENT_ID={your client id}
    export AZURE_CLIENT_SECRET={your client secret}
    export AZURE_SUBSCRIPTION_ID={your subscription id}
    ```

    > [AZURE.NOTE] On Windows, use `set` instead of `export`.

1. Run the sample.

    ```
    go run example.go
    ```

<a id="info"></a>

## More information

- [Service Bus documentation](https://docs.microsoft.com/en-us/azure/service-bus/)

***

This project has adopted the [Microsoft Open Source Code of Conduct](https://opensource.microsoft.com/codeofconduct/). For more information see the [Code of Conduct FAQ](https://opensource.microsoft.com/codeofconduct/faq/) or contact [opencode@microsoft.com](mailto:opencode@microsoft.com) with any additional questions or comments.