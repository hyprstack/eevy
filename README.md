# Description #
eevy is a micro weight go based message broker.

Sources listen on http ports, AWS SQS queues etc for incoming events.  Event listeners keep an eye out for events with specific names and invoke handlers when certain events are detected.

Handlers can invoke an AWS Lambda function, place the event data on an AWS SQS queue, call an http end point with OAuth2 authentication etc.
# Getting started #
## Installing ##
### From Source ###
Eevy is built using go so ensure that you have this properly installed first.
```
go get github.com/hevnly/eevy
go build github.com/hevnly/eevy
```
### Using binary ###
Download the binary file from the github release page(currently built using Ubuntu system)
## Configuration ##
A yaml file is used to configure the sources, handlers and listeners. It must be called "conf.yml" and sit beside the eevy binary file.  Below is an example configuration file;
```yaml
logs:
  event: stdout
  handler: stderr
  app: /var/log/eevy/app.log
sources:
  - type: sqs
    url: https://sqs.eu-west-1.amazonaws.com/999999999/lambda
    region: eu-west-1
  - type: http
    port: 8000
handlers:
  testSqs:
    type: sqs
    region: eu-west-1
    url: https://sqs.eu-west-1.amazonaws.com/999999999999/output
  testLambda:
    type: lambda
    region: eu-west-1
    function: test
  testSqs2:
    type: sqs
    region: eu-west-1
    url: https://sqs.eu-west-1.amazonaws.com/999999999999/output
listeners:
  "_":
    - testSqs
  application.*:
    - testLambda
  test.event2:
    - testSqs2
```
The above example creates two sources; an AWS SQS and starts listening on port 8000 for HTTP requests.
There are three handlers defined which can be used throughout eevy.
The "_" listener is a wild card and responds to every event received (this would make more sense if it was "\*" however yaml doesn't like special characters at start of a key).  It evokes the handlers that are defined it its array, in this example it places the events message received onto an AWS SQS queue via he "testSqs" handler. "application.*" invokes its handler when any event that begins with "application" is detected.  Finally "test.event2" will only be called when exactly that event is revived.  Without the trialling "*" it will not respond to any event that has a sub event name eg "test.event2.anything".
## Sources ##
### AWS SQS ###
#### Description ####
This listener pulls events from the specified AWS SQS queue.
#### Config ####
| Name | Type | Description                                |
| ---- | ---- | ------------------------------------------ |
|type  |string| This must be set to "sqs"                  |
|url   |string| The URL for the SQS                        |
|region|string| The aws region it belongs to eg "eu-west-1"|

### HTTP ###
#### Description ####
This source listens on a specified port for HTTP connections. A POST request should be sent to "http://hostname:<port>/event/".  The body of the post request is treated as the message.
#### Config ####
| Name | Type | Description                                |
| ---- | ---- | ------------------------------------------ |
|type  |string| This must be set to "http"                 |
|port  |int   | The port number to listen on               |
|bindIp|string| Ip address to bind to                      |

## Handlers ##
### AWS SQS ###
#### Description ####
When a relevant event is received place its message section onto the supplied AWS SQS.
#### Config ####
| Name   | Type | Description                                |
| ------ | ---- | ------------------------------------------ |
|type    |string| This must be set to "sqs"                  |
|url     |string| The URL for the SQS                        |
|region  |string| The aws region it belongs to eg "eu-west-1"|
|message |string| The body of the request sent to the listener |

### AWS Lambda ###
#### Description ####
This listener invokes an AWS Lambda function
#### Config ####
| Name   | Type | Description                                  |
| ------ | ---- | -------------------------------------------- |
|type    |string| This must be set to "lambda"                 |
|function|string| The name of the function to be invoked       |
|region  |string| The AWS region it belongs to eg "eu-west-1"  |
|message |string| The body of the request sent to the listener |

### cli ###
#### Description ####
This listener performs a local system call
#### Config ####
| Name   | Type | Description                                  |
| ------ | ---- | -------------------------------------------- |
|type    |string| This must be set to "lambda"                 |
|function|string| The name of the function to be invoked       |
|region  |string| The AWS region it belongs to eg "eu-west-1"  |
|message |string| The body of the request sent to the listener |

### OAuth2 ###
#### Description ####
Use this listener when sending events to an OAuth 2 end point such as the hevnly api.
#### Config ####
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

### Variables ###
Sometimes the end point of an OAuth call we want a handler to call is embedded in the event message.  To allow this we can use variables.  Variables are handler values that are surrounded by "${message.name}".  The following example shows how this should be used;
*Event*
```json
{
    "event":"api.testSqs",
    "message":{"verb":"post","body":"test","url":"http:\/\/hevnly.dev\/img\/r55435730c954c\/resized\"}
}
```
*Handler configuration*
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
These variables can be used in most handler configuration fields.
# Events #
An event is simply a json object that must have the following structure
```json
{
    "event": "event.name",
    "message": {"desc": "json object with the associated events data"} 
}
```