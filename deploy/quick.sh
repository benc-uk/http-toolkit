#!/usr/bin/env bash
set -euo pipefail

REPO_ROOT=$(git rev-parse --show-toplevel)

# Quick deployment script for http-toolkit to Azure Container Apps
# This is not intended as a general purpose script, so modify as needed

which az > /dev/null || { echo -e "💥 Error! Azure CLI is not installed. https://aka.ms/azure-cli"; exit 1; }

SUB_NAME=$(az account show --query name -o tsv)
TENANT_ID=$(az account show --query tenantDisplayName -o tsv)

echo -e "\e[34m⛅ Azure details: \e[0m"
echo -e " 🔑 \e[34mSubscription: \e[33m$SUB_NAME\e[0m"
echo -e " 🌐 \e[34mTenant:       \e[33m$TENANT_ID\e[0m"

if [[ "${NOPROMPT-0}" != "1" ]]; then 
  read -r -p "🤔 Proceed with deployment? [Y/n] " response
  response=${response,,}    # tolower
  if [[ ! "$response" =~ ^(yes|y|"")$ ]]; then echo -e "\e[31m👋 Exiting...\e[0m"; exit 1; fi
fi

echo -e "\e[32m🔨 Deploying http-toolkit to an Azure Container App...\e[0m"

# We deploy from source as `az containerapp up` supports this, neat!
az containerapp up --name http-toolkit \
--source "$REPO_ROOT" \
--target-port 8000 \
--resource-group http-toolkit \
--environment http-toolkit \
--location uksouth