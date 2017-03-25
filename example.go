// This package demonstrates how to manage Azure Service Bus using Go.
package main

import (
	"fmt"
	"os"

	"github.com/Azure/azure-sdk-for-go/arm/resources/resources"
	"github.com/Azure/azure-sdk-for-go/arm/servicebus"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/to"
)

// This example requires that the following environment vars are set:
//
// AZURE_TENANT_ID: contains your Azure Active Directory tenant ID or domain
// AZURE_CLIENT_ID: contains your Azure Active Directory Application Client ID
// AZURE_CLIENT_SECRET: contains your Azure Active Directory Application Secret
// AZURE_SUBSCRIPTION_ID: contains your Azure Subscription ID
//

var (
	groupsClient    resources.GroupsClient
	namespaceClient servicebus.NamespacesClient
	queuesClient    servicebus.QueuesClient
	topicsClient    servicebus.TopicsClient
	subClient       servicebus.SubscriptionsClient

	location          = "westus"
	resourceGroupName = "azure-sample"
	namespaceName     = "golangrocksonazure"
	authRuleName      = "authrule"
	queueName         = "queue1"
	topicName         = "topic1"
	subscriptionName  = "sub1"
)

func init() {
	subscriptionID := getEnvVarOrExit("AZURE_SUBSCRIPTION_ID")
	tenantID := getEnvVarOrExit("AZURE_TENANT_ID")

	oauthConfig, err := azure.PublicCloud.OAuthConfigForTenant(tenantID)
	onErrorFail(err, "OAuthConfigForTenant failed")

	clientID := getEnvVarOrExit("AZURE_CLIENT_ID")
	clientSecret := getEnvVarOrExit("AZURE_CLIENT_SECRET")
	spToken, err := azure.NewServicePrincipalToken(*oauthConfig, clientID, clientSecret, azure.PublicCloud.ResourceManagerEndpoint)
	onErrorFail(err, "NewServicePrincipalToken failed")

	createClients(subscriptionID, spToken)
}

func main() {
	createResourceGroup()

	createNamespace()
	createAuthRule()
	listKeys()
	createQueue()
	createTopic()
	createSubscription()

	fmt.Print("Press enter to delete all the resources created in this sample...")
	var input string
	fmt.Scanln(&input)
	deleteResourceGroup()
}

func createResourceGroup() {
	fmt.Printf("Creating resource group '%s'...\n", resourceGroupName)
	groupParams := resources.Group{
		Location: to.StringPtr(location),
	}
	_, err := groupsClient.CreateOrUpdate(resourceGroupName, groupParams)
	onErrorFail(err, "Group create failed")
}

func createNamespace() {
	fmt.Printf("Creating namespace '%s'...\n", namespaceName)
	namespaceParams := servicebus.NamespaceCreateOrUpdateParameters{
		Location: to.StringPtr(location),
		Sku: &servicebus.Sku{
			Tier: servicebus.SkuTierStandard,
			Name: servicebus.Standard,
		},
	}
	_, err := namespaceClient.CreateOrUpdate(resourceGroupName, namespaceName, namespaceParams, nil)
	onErrorFail(err, "Namespace create failed")
}

func createAuthRule() {
	fmt.Printf("Creating authorization rule '%s' for namespace '%s'...\n", authRuleName, namespaceName)
	ruleParams := servicebus.SharedAccessAuthorizationRuleCreateOrUpdateParameters{
		SharedAccessAuthorizationRuleProperties: &servicebus.SharedAccessAuthorizationRuleProperties{
			Rights: &[]servicebus.AccessRights{
				servicebus.Listen,
				servicebus.Manage,
				servicebus.Send,
			},
		},
	}
	_, err := namespaceClient.CreateOrUpdateAuthorizationRule(resourceGroupName, namespaceName, authRuleName, ruleParams)
	onErrorFail(err, "Authorization rule create failed")
}

func listKeys() {
	fmt.Printf("List keys for '%s' namespace...\n", namespaceName)
	keys, err := namespaceClient.ListKeys(resourceGroupName, namespaceName, authRuleName)
	onErrorFail(err, "List keys failed")
	fmt.Printf("\tKey name: %s\n", *keys.KeyName)
	fmt.Printf("\tPrimary key: %s\n", *keys.PrimaryKey)
	fmt.Printf("\tSecondary key: %s\n", *keys.SecondaryKey)
	fmt.Printf("\tPrimary connection string: %s\n", *keys.PrimaryConnectionString)
	fmt.Printf("\tSecondary connection string: %s\n", *keys.SecondaryConnectionString)
}

func createQueue() {
	fmt.Printf("Creating queue '%s'...\n", queueName)
	queueParams := servicebus.QueueCreateOrUpdateParameters{
		Location: to.StringPtr(location),
		QueueProperties: &servicebus.QueueProperties{
			EnablePartitioning: to.BoolPtr(true),
		},
	}
	_, err := queuesClient.CreateOrUpdate(resourceGroupName, namespaceName, queueName, queueParams)
	onErrorFail(err, "Queue create failed")
}

func createTopic() {
	fmt.Printf("Creating topic '%s'...\n", topicName)
	topicsParams := servicebus.TopicCreateOrUpdateParameters{
		Location: to.StringPtr(location),
		TopicProperties: &servicebus.TopicProperties{
			EnablePartitioning: to.BoolPtr(true),
		},
	}
	_, err := topicsClient.CreateOrUpdate(resourceGroupName, namespaceName, topicName, topicsParams)
	onErrorFail(err, "Topic create failed")
}

func createSubscription() {
	fmt.Printf("Creating subscription '%s'...\n", subscriptionName)
	subParams := servicebus.SubscriptionCreateOrUpdateParameters{
		Location: to.StringPtr(location),
	}
	_, err := subClient.CreateOrUpdate(resourceGroupName, namespaceName, topicName, subscriptionName, subParams)
	onErrorFail(err, "Subscription create failed")
}

func deleteResourceGroup() {
	fmt.Printf("Deleting resource group '%s'...\n", resourceGroupName)
	_, err := groupsClient.Delete(resourceGroupName, nil)
	onErrorFail(err, "Group delete failed")
}

func createClients(subID string, token *azure.ServicePrincipalToken) {
	groupsClient = resources.NewGroupsClient(subID)
	groupsClient.Authorizer = token

	namespaceClient = servicebus.NewNamespacesClient(subID)
	namespaceClient.Authorizer = token

	queuesClient = servicebus.NewQueuesClient(subID)
	queuesClient.Authorizer = token

	topicsClient = servicebus.NewTopicsClient(subID)
	topicsClient.Authorizer = token

	subClient = servicebus.NewSubscriptionsClient(subID)
	subClient.Authorizer = token
}

// getEnvVarOrExit returns the value of specified environment variable or terminates if it's not defined.
func getEnvVarOrExit(varName string) string {
	value := os.Getenv(varName)
	if value == "" {
		fmt.Printf("Missing environment variable %s\n", varName)
		os.Exit(1)
	}

	return value
}

// onErrorFail prints a failure message and exits the program if err is not nil.
func onErrorFail(err error, message string) {
	if err != nil {
		fmt.Printf("%s: %s\n", message, err)
		os.Exit(1)
	}
}
