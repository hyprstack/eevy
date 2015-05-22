# README #
## Description ##
This code is the backbone for our event driven architecture.  It listens to sources which could be an AWS SQS, an HTTP port or anything else so long as we write the handler.  When an event is recieved by one of the sources, calls are made to registered listeners which could be a AWS Lambda function, HTTP OAuth2 call, placing event an an AWS SQS etc.
## Events ##
An event is simply a json object that must have the following structure
```json
{
    "event": "event.name",
    "message": "the data that belongs to this event"
}
```
## Config ##
A yaml file is used to configure the sources and listeners. Below is an example configuration file;
```yaml
sources:
  - type: sqs
    url: https://sqs.eu-west-1.amazonaws.com/999999999/lambda
    region: eu-west-1
  - type: http
    port: 8000
listeners:
  _:
    - type: sqs
      region: eu-west-1
      url: https://sqs.eu-west-1.amazonaws.com/999999999999/output
  api.*:
    - type: lambda
      region: eu-west-1
      function: test
  api.testSqs:
    - type: sqs
      region: eu-west-1
      url: https://sqs.eu-west-1.amazonaws.com/999999999999/output
```
The above example creates two sources; an AWS SQS and starts listening on port 8000 for HTTP requests.
There are three listeners. "_" listener is a wild card and responds to every event received (this would make more sense if it was "\*" however yaml doesn't like special characters at start of a key).  It places the message received onto an AWS SQS queue. "api.*" responds to any event that begins with "api".  When an event that begins with "api" is received the AWS Lambda function is executed.  Finally "api.testSqs" only fires when an event with exactly this name is revived.
### Sources ###
#### AWS SQS ####
##### Description #####
This listener pulls events from the specified AWS SQS queue.
##### Config #####
| Name | Type | Description                                |
| ---- | ---- | ------------------------------------------ |
|type  |string| This must be set to "sqs"                  |
|url   |string| The URL for the SQS                        |
|region|string| The aws region it belongs to eg "eu-west-1"|

#### HTTP ####
##### Description #####
This source listens on a specified port for HTTP connections. A POST request should be sent to "http://hostname:<port>/event/".  The body of the post request is treated as the message.
##### Config #####
| Name | Type | Description                                |
| ---- | ---- | ------------------------------------------ |
|type  |string| This must be set to "http"                 |
|port  |int   | The port number to listen on               |

### Listeners ###
#### AWS SQS ####
##### Description #####
When a relevant event is recieved place its message section onto the supplied AWS SQS.
##### Config #####
| Name   | Type | Description                                |
| ------ | ---- | ------------------------------------------ |
|type    |string| This must be set to "sqs"                  |
|url     |string| The URL for the SQS                        |
|region  |string| The aws region it belongs to eg "eu-west-1"|
|message |string| The body of the request sent to the listener |

#### AWS Lambda ####
##### Description #####
This listener invokes an AWS Lambda function
##### Config #####
| Name   | Type | Description                                  |
| ------ | ---- | -------------------------------------------- |
|type    |string| This must be set to "lambda"                 |
|function|string| The name of the function to be invoked       |
|region  |string| The AWS region it belongs to eg "eu-west-1"  |
|message |string| The body of the request sent to the listener |

#### OAuth2 ####
##### Description #####
Use this listener when sending events to an OAuth 2 end point such as the hevnly api.
##### Config #####
| Name        | Type   | Description                                           |
| ----------- | ------ | ----------------------------------------------------- |
|type         |string  | This must be set to "oauth2"                          |
|clientId     |string  | OAuth2 client id to use during authentication         |
|clientSecret |string  | The AWS region it belongs to eg "eu-west-1"           |
|scope        |[]string| Array of scopes to use during authentication          |
|tokenUrl     |string  | OAuth2 token URL                                      |
|endPoint     |string  | URL to call after OAuth2 authentication               |
|verb         |string  | ("get" or "post") HTTP verb used when calling endpoint|
|message      |string  | The body of the request sent to the listener          |

##### Variables #####
Some times the endpoint, verb or body can come from the event message.  If this is the case you can set these by using variables.
The example below will call the endpoint that is defined in the event message. "${message.<what ever>}" symbolises that this value is defined by a parameter in the event message.  This assumes that the message in the event is a json string.  If we had $(message.testUrl) it would assume that the json string in the event message will have an element "testUrl" and that its value should be used.
###### Event ######
```json
{
    "event":"api.testSqs",
    "message":"{\"verb\":\"post\",\"body\":\"test\",\"url\":\"http:\/\/hevnly.dev\/img\/r55435730c954c\/resized\"}"
}
```
###### Listener config ######
```yaml
- type: oauth2
  clientId: testClientID
  clientSecret: testClientSecret
  scope:
    - oauth2
  tokenUrl: http://test.dev/oauth/v2/token
  endPoint: ${message.url}
  body: ${message.body}
  verb: ${message.verb}
```