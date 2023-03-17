# UTMStackDatasourcesSalesForce

---
## Description
UTMStack Datasource for Salesforce is a tool developed in Go (using golang 1.20), to interact with Salesforce platform RESTFul API `v57.0`.
The main function is to extract logs from Salesforce and send them to UTMStack correlation endpoint.

## Contents
- [Configuration](#configuration)
- [Build with Docker](#build-with-docker)
- [Usage with docker for production](#usage-with-docker-for-production)
- [Salesforce documentation reference](#salesforce-documentation-reference)
    - [Instance URL](#instance-url)
    - [User information](#user-information)
    - [Customer key and secret](#customer-key-and-secret)


## Configuration
The first step to ensure communication between the tool and Salesforce platform is the
configuration of some environment variables listed as follows:

Variables marked as (`Required`) must be defined even if it has default value. (`Optional`) variables can be omitted in which
case will use the default value
- clientID: (`Required`) Represents the _customer key_ from the _connected app_ configuration. Default value `"not set"`.
- clientSecret: (`Required`) Represents the _customer secret_ from the _connected app_ configuration. Default value `"not set"`.
- username: (`Required`) Represents the username used to connect to the Salesforce platform. Default value `"not set"`.
- password: (`Required`) Represents the user's password used to connect to the Salesforce platform. Default value `"not set"`.
- securityToken: (`Required`) Represents the security token associated with the user above. Default value `"not set"`.
- instanceUrl: (`Required`) Represents the instance URL provided by Salesforce. Default value `"not set"`.
- OAuthService: (`Optional`) Represents the Salesforce base login URL. Default value `"https://login.salesforce.com"`.
- LoginEndpoint: (`Optional`) Represents the Salesforce endpoint to retrieve the session token. Default value `"/services/oauth2/token"`.
- EventsEndPoint: (`Optional`) Represents the Salesforce endpoint of the event log files, to get the data by Id. Default value `"/services/data/v57.0/sobjects/EventLogFile"`. 
- QueryEndPoint: (`Optional`) Represents the Salesforce endpoint to query and return all the event log files endpoints. Default value `"/services/data/v57.0/query"`
- siemURL: (`Optional`) Represents the log destination after extracted and transformed from Salesforce. Default value `"http://correlation:8080/v1/newlog"`

**Note:** For more information about how to get some variables values from Salesforce, check [here](#salesforce-documentation-reference) 

## Build with Docker

To build the application for your own use, you must build it based on the Dockerfile located at the root of the application files,
then run `docker build` command Ex: `docker build -t sforceds:latest -f Dockerfile .`, after that
you can run the image using `docker run` command and passing all the variables listed above with `-e varValue`.
<br>[Back to Contents](#contents)

## Usage with docker for production
To execute the tool with the stable release, you must create a docker container using this docker image `docker pull ghcr.io/atlasinsidecorp/sforceds:latest`,
with the environment variables configured as described before.
Also, to avoid processing the **same logs over and over again**, you must create a docker volume pointing to
`/local_storage` folder with read and write permission.

To do that, create the folder on the machine to map the volume: `mkdir -p /utmstack/sforceds/`
Then, add the volume to docker compose config file, as follows:
~~~
...
volumes:
      - /utmstack/sforceds/:/local_storage
...
~~~
Start the docker compose.
If you aren't using compose create the volume with `docker volume create` command and associate it to the container before executing.
<br>[Back to Contents](#contents)

## Salesforce documentation reference
If you don't know how to get the values of the environment variables related to Salesforce platform, 
here you can find some useful information about it.

### Instance URL
The instance url is provided by Salesforce when you create the account, but in case you've lost the email,
you can try one of:
- Click your profile icon at the right top of the screen, below the user information is the instance URL
- Type "My Domain" in the quick search at the top of the screen, then in the "My Domain Details" section
look for the value of "Current My Domain URL"

Example: `myproj-dev-ed.develop.my.salesforce.com`, ensure that you add the protocol `https` to the environment variable
value - > `https://myproj-dev-ed.develop.my.salesforce.com`
<br>[Back to Contents](#contents)

### User information
The user information and credentials are provided by Salesforce when you create the account, but in case you've lost the email,
 click on your profile icon at the right top of the screen, and go to settings, you will be redirected to you "Personal Information":

- Username: You should see your _username_ and more in the "Personal Information" option of the left menu.
- Password: You can change your password with the "Change My Password" option on the menu at the left, but be aware, when you change 
user information like password a new security token is generated
- Security Token: With the "Reset My Security Token" option on the menu at the left, you can get
  a new security token.

<br>[Back to Contents](#contents)

### Customer key and secret
Customer key and secret used for oauth authentication to the Salesforce API from external applications, require that you have a "Connected App" configured, 
to do that:

- Login to your instance
- Go to "Setup" -> "Setup for current app" option at the right top of the screen.
- Type "App Manager" in the quick search box and select the "App Manager" option.
- On the screen "Lightning Experience App Manager" click on "New Connected App" button
  and fill the fields
- In the "API (Enable OAuth Settings)" section, check "Enable OAuth Settings", you can use
  `https://login.salesforce.com/services/oauth2/success` as the "Callback URL" value.
- In the same section be sure that you select "Access Connect REST API resources (chatter_api)" 
  in "Selected OAuth Scopes".
- Save your changes, and select "continue" on the next screen.
- You will be redirected to a page with details of the Connected App you just created, 
  in the "API (Enable OAuth Settings)" section, click the button "Manage Consumer Details" next to
  "Consumer Key and Secret", you will be prompted to verify your identity by a code sent to your
  email. Once you verify your identity you must see your customer key and secret in the "Consumer Details"
  section.
- Save the customer key and secret in a safe place with the "copy" button.

If you already have a "Connected App" configured, to get the customer key and secret values
you have to:

- Login to your instance
- Go to "Setup" -> "Setup for current app" option at the right top of the screen.
- Type "App Manager" in the quick search box and select the "App Manager" option.
- Search your connected app by "App Name" column, then click the icon on the final column of the row
  and select "View" option.
- You will be redirected to a page with details of the Connected App,
  in the "API (Enable OAuth Settings)" section, click the button "Manage Consumer Details" next to
  "Consumer Key and Secret", you will be prompted to verify your identity by a code sent to your
  email. Once you verify your identity you must see your customer key and secret in the "Consumer Details"
  section.
- Save the customer key and secret in a safe place with the "copy" button.

<br>[Back to Contents](#contents)