# Onedrive API

### Registering your app for Microsoft Graph

To connect with Microsoft Graph, you'll need a work/school account or a Microsoft account.

1. Go to the [Microsoft Application Registration Portal](https://aka.ms/appregistrations).
2. When prompted, sign in with your account credentials.
3. Find My applications and click Add an app.
4. Enter your app's name and click Create application.

### Token authentication flow

To start the sign-in process with the token flow, use a web browser or web-browser control to load a URL request.

GET https://login.microsoftonline.com/common/oauth2/v2.0/authorize?client_id={client_id}&scope={scope}&response_type=code&redirect_uri={redirect_uri}

Upon successful authentication and authorization of your application, the web browser is redirected to the redirect URL provided with additional parameters added to the URL.

https://myapp.com/auth-redirect#access_token={access_token}&authentication_token={authentication_token}&token_type={token_type}&expires_in={expires_in}&scope={scope}&user_id={user_id}

GET https://login.microsoftonline.com/common/oauth2/v2.0/logout?post_logout_redirect_uri={redirect-uri}

POST https://login.microsoftonline.com/common/oauth2/v2.0/token
Content-Type: application/x-www-form-urlencoded

client_id={client_id}&redirect_uri={redirect_uri}&client_secret={client_secret}
&code={code}&grant_type=authorization_code