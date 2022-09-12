# Video store

An API for storing/retrieving videos from "any generic" video storage platform.

Although the only targeted platform is Youtube for the time being, the point of the generic approach
is to be able to change easily if needed be.

## Platforms

### Configuring Youtube

Google's Oauth is a bit of a pain to work with. 
To be able to upload videos, we have to include the scope *https://www.googleapis.com/auth/youtube.upload*.

However, this scope requires a three-legged oauth validation, and there is no real way to use a Service Account
for this specific API.

So, to make this work :
1. Go to the [Google API dashboard](https://console.cloud.google.com/apis/dashboard) and create a project
2. In the project dashboard, in **Credentials** create a new **Web application** client id/client secret pair
3. Next, go to [Google's oauth playground website](https://developers.google.com/oauthplayground/), click on the cogwheel icon and input your clientID/client secret
4. In scopes, select *https://www.googleapis.com/auth/youtube*, *https://www.googleapis.com/auth/youtube* and click on "Authorize APIs"
5. You'll get an authorization code, click on "Exchange for token" to get an access token and a refresh token. 

With this, we have all the values for the env variables needed for authenticating to Youtube :

- YT_CLIENT_ID
- YT_CLIENT_SECRET
- YT_REFRESH_TOKEN

The access token will be (re)generated from the refresh token automatically.