package auth

import (
	"context"
	"errors"
	"log"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

const (
	credNameAzureCLI = "AroCLICredential"

	AzurePublicString       = "public"
	AzureUSGovernmentString = "usgovernment"
	AzureChinaString        = "china"
)

type Config struct {
	SubscriptionId string
	TenantId       string
	ClientId       string
	ClientSecret   string
	Environment    string
}

type DefaultAroCredential struct {
	chain *azidentity.ChainedTokenCredential
}

func NewDefaultAroCredential(config Config) (*DefaultAroCredential, error) {
	var creds []azcore.TokenCredential
	var errorMessages []string

	options := &azidentity.ClientSecretCredentialOptions{
		ClientOptions: azcore.ClientOptions{
			Cloud: getCloud(config),
		},
	}

	clientSecretCred, err := azidentity.NewClientSecretCredential(config.TenantId, config.ClientId, config.ClientSecret, options)
	if err == nil {
		creds = append(creds, clientSecretCred)
	} else {
		errorMessages = append(errorMessages, "AroClientSecretCredential: "+err.Error())
	}

	cliCred, err := azidentity.NewAzureCLICredential(nil)
	if err == nil {
		creds = append(creds, cliCred)
	} else {
		errorMessages = append(errorMessages, "AroCLICredential: "+err.Error())
	}

	err = defaultAroCredentialConstructorErrorHandler(len(creds), errorMessages)
	if err != nil {
		return nil, err
	}

	chain, err := azidentity.NewChainedTokenCredential(creds, nil)
	if err != nil {
		return nil, err
	}
	return &DefaultAroCredential{chain: chain}, nil
}

// GetToken requests an access token from Azure Active Directory. This method is called automatically by Azure SDK clients.
func (c *DefaultAroCredential) GetToken(ctx context.Context, opts policy.TokenRequestOptions) (azcore.AccessToken, error) {
	return c.chain.GetToken(ctx, opts)
}

func defaultAroCredentialConstructorErrorHandler(numberOfSuccessfulCredentials int, errorMessages []string) (err error) {
	errorMessage := strings.Join(errorMessages, "\n\t")

	if numberOfSuccessfulCredentials == 0 {
		return errors.New(errorMessage)
	}

	if len(errorMessages) != 0 {
		log.Printf("NewDefaultAroCredential failed to initialize some credentials:\n\t%s", errorMessage)
	}

	return nil
}

func getCloud(config Config) cloud.Configuration {
	switch config.Environment {
	case AzureChinaString:
		return cloud.AzureChina
	case AzureUSGovernmentString:
		return cloud.AzureGovernment
	default:
		return cloud.AzurePublic
	}
}