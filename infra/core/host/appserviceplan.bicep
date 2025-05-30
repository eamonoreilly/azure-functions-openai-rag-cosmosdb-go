param name string
param location string = resourceGroup().location
param tags object = {}
param kind string = ''
param sku object

resource appServicePlan 'Microsoft.Web/serverfarms@2024-04-01' = {
  name: name
  location: location
  tags: tags
  sku: sku
  kind: kind
}

output id string = appServicePlan.id
