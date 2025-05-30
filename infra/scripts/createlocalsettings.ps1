$ErrorActionPreference = "Stop"

if (-not (Test-Path ".\app\local.settings.json")) {

    $output = azd env get-values

    # Parse the output to get the endpoint values
    foreach ($line in $output) {
        if ($line -match "AZURE_COSMOSDB_ENDPOINT"){
            $CosmosDBEndPoint = ($line -split "=")[1] -replace '"',''
        }
        if ($line -match "AZURE_OPENAI_ENDPOINT"){
            $OpenAIEndPoint = ($line -split "=")[1] -replace '"',''
        }
    }

    @{
        "IsEncrypted" = "false";
        "Values" = @{
            "AzureWebJobsStorage" = "UseDevelopmentStorage=true";
            "FUNCTIONS_WORKER_RUNTIME" = "custom";
            "AZURE_OPENAI_ENDPOINT" = "$OpenAIEndPoint";
            "CHAT_MODEL_DEPLOYMENT_NAME" = "chat";
            "cosmosDBNoSqlEndpoint__ENDPOINT" = "$CosmosDBEndPoint";
            "EMBEDDING_MODEL_DEPLOYMENT_NAME" = "embeddings";
            "SYSTEM_PROMPT" = "You must only use the provided documents to answer the question";
        }
    } | ConvertTo-Json | Out-File -FilePath ".\app\local.settings.json" -Encoding ascii
}