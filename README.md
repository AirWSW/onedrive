# OneDrive API client

OneDrive API client using Microsoft Graph API

### Feature

1. OneDrive API client wrapper as OneIndex, PyIndex etc.
2. Cache while first request income and auto refreshing.
3. In memory cache with permanent local cache file (Experiment).
4. Multi OneDrives within one endpoint, specify one drive by using query (Experiment).
5. Run without other requirements, just configure once and run it.

#### Danger testing

1. Volume mounts and mount type setting.

#### Known BUGs

### TODO list

1. File upload and ongoing download support.
2. Buildin web graph user interface support.
3. Administration controls support.

### Registering your app for Microsoft Graph

To connect with Microsoft Graph, you'll need a work/school account or a Microsoft account.

1. Go to the [Microsoft Application Registration Portal](https://aka.ms/appregistrations).
2. When prompted, sign in with your account credentials.
3. Find My applications and click Add an app.
4. Enter your app's name and click Create application.

### Setup your config file

**Minimize setting template**

```json
{
  "oneDrives": [
    {
      "microsoftEndPoints": {
        "azureAdPortalEndPointUrl": "https://portal.azure.com",
        "azureAdEndPointUrl": "https://login.microsoftonline.com",
        "microsoftgraphApiEndPointUrl": "https://graph.microsoft.com"
      },
      "azureAdAppRegistration": {
        "clientId": "Your Azure AD App Client ID",
        "redirectUris": ["Your Azure AD App Redirect URL"],
        "clientSecret": "Your Azure AD App Client Secret"
      },
      "azureAdAuthFlowContext": {
        "grantScope": "Files.ReadWrite User.Read offline_access"
      },
      "oneDriveDescription": {
        "rootPath": "root",
        "refreshInterval": 3600
      }
    }
  ]
}
```

### API endpoints of Microsoft

#### Azure AD portal endpoint

```
https://portal.azure.com
https://portal.azure.cn           (Azure AD China operated by 21Vianet)
https://portal.microsoftazure.de  (Azure AD Germany)
```

#### Azure AD endpoint

```
https://login.microsoftonline.com
https://login.chinacloudapi.cn     (Azure AD China operated by 21Vianet)
https://login.microsoftonline.de   (Azure AD Germany)
```

#### Microsoft Graph API

```
https://graph.microsoft.com
https://microsoftgraph.chinacloudapi.cn  (Microsoft Graph China operated by 21Vianet)
https://graph.microsoft.de               (Microsoft Graph Germany)
```
