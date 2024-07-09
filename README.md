# DEPRECATED

> **WARN:** Please note that this provider is now considered deprecated in favor of 
https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/redhat_openshift_cluster.


# Terraform Provider Azure Redhat Openshift

This is a provider to create [Azure Redhat Openshift](https://docs.microsoft.com/en-us/azure/openshift/)


## Prerequistes

### Authencation

* Subscripition id is required at Provider Level

    ```
    provider azureopenshift {
      subscription_id = "xxxx"
    }
    ```

    User can configure subscription id through environment variable

    ```
    ARM_SUBSCRIPTION_ID=xxxx
    ```

* Provider automatically honor Azure CLI Login credentials
* Provider Supports Service Principal Environment varialbles

    ```
    ARM_CLIENT_ID=xxxx
    ARM_CLIENT_SECRET=xxxx
    ARM_SUBSCRIPTION_ID=xxxx
    ARM_TENANT_ID=xxxx
    ```


### [Create Azure network with two empty subnets](https://docs.microsoft.com/en-us/azure/openshift/tutorial-create-cluster#create-a-virtual-network-containing-two-empty-subnets)
* Azure Resource Group
* Azure network
* Master subnet & Worker subnet

### [Create a service principal](https://docs.microsoft.com/en-us/azure/openshift/howto-create-service-principal?pivots=aro-azurecli)

### Give Red Hat Openshift Resource Provider network contributor role of Azure network

```bash
OPENSHIFT_RP_OBJECT_ID=$(az ad sp list --filter "displayname eq 'Azure Red Hat OpenShift RP'" --query "[?appDisplayName=='Azure Red Hat OpenShift RP'].objectId" --only-show-errors --output tsv)
az role assignment create --role "Network Contributor" --assignee-object-id ${OPENSHIFT_RP_OBJECT_ID} --scope [NETWORK_ID]
```

## Example Usage

```bash
terraform {
  required_providers {
    azureopenshift = {
      source  = "rh-mobb/azureopenshift"
      version = "~> 0.0.10"
    }
  }
}

provider azureopenshift {
}

resource "azureopenshift_redhatopenshift_cluster" "test" {
  name                = "tf-openshift"
  location            = var.location
  resource_group_name = var.resource_group_name

  master_profile {
    subnet_id = var.master_subnet_id
  }

  worker_profile {
    subnet_id = var.worker_subnet_id
  }

  service_principal {
    client_id     = var.client_id
    client_secret = var.client_secret
  }
}
```

