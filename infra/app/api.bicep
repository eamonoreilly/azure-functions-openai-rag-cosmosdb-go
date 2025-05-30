// Description: Bicep module to deploy an Azure Functions App with a Linux App Service Plan
param name string
param location string = resourceGroup().location
param tags object = {}
param applicationInsightsName string = ''
param appServicePlanId string
param appSettings object = {}
param runtimeName string 
param serviceName string = 'api'
param storageAccountName string
param virtualNetworkSubnetId string = ''
param identityId string = ''
param identityClientId string = ''
param aiServiceUrl string = ''
param cosmosDBEndpoint string = ''

var applicationInsightsIdentity = 'ClientId=${identityClientId};Authorization=AAD'

module api '../core/host/functions-windows-appservice.bicep' = {
  name: '${serviceName}-functions-module'
  params: {
    name: name
    location: location
    tags: union(tags, { 'azd-service-name': serviceName })
    identityType: 'UserAssigned'
    identityId: identityId
    appSettings: union(appSettings,
      {
        AzureWebJobsStorage__clientId : identityClientId
        APPLICATIONINSIGHTS_AUTHENTICATION_STRING: applicationInsightsIdentity
        AZURE_OPENAI_ENDPOINT: aiServiceUrl
        AZURE_CLIENT_ID: identityClientId
        cosmosDBNoSqlEndpoint__Endpoint: cosmosDBEndpoint
        cosmosDBNoSqlEndpoint__clientId: identityClientId
        cosmosDBNoSqlEndpoint__credential: 'managedidentity'
        SYSTEM_PROMPT: 'You must only use the provided documents to answer the question'
      })
    applicationInsightsName: applicationInsightsName
    appServicePlanId: appServicePlanId
    kind: 'functionapp'
    runtimeName: runtimeName
    storageAccountName: storageAccountName
    virtualNetworkSubnetId: virtualNetworkSubnetId
  }
}

output SERVICE_API_NAME string = api.outputs.name
output SERVICE_API_URI string = api.outputs.uri
output SERVICE_API_IDENTITY_PRINCIPAL_ID string = api.outputs.identityPrincipalId
