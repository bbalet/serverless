# Serverless Function

Complete example of a serverless function using an external data source (MySQL). The function is triggered by an HTTP GET method and will return and increment a counter of visits.

## A note about dependancies

The MySQL driver handles connection pooling, so you don't have to worry about the extra time needed for establishing a DB connection between many requests if you function is always called after a warm startup.

logrus allows you to define fields in the log stream.

This project requires go 1.13+ so you don't need to vendor these deps into a vendor folder as we rely on modules.

## Prerequisites & Test

You need to access to a private or public MySQL instance with a basic schema (see *counter.sql*). If you deploy the code on Google Cloud, the MySQL server must be accessible by the underlying VM running your serverless function.

As a 12-factors application as a service, this sample code requires you to set the configuration of the MySQL database into environment variable. So before trying to test or benchmark the code, you must make a copy of either *setenv.ps1* or *setenv.sh* into, let's say, *mysetenv.ps1* or *mysetenv.sh* (and then launch one the files so as to populate your environnement variables). Then you can test your setup:

    go test

## Deploy

As you will notice with this example, you can use the same code base for OpenFaaS and Google Cloud Functions.

### OpenFaaS

Build the package using OpenFaaS template:

    faas-cli build --image counterexample --handler . --name GetCounterOpenFaas --lang go

Passing environment variables to the underlying VM:
https://docs.openfaas.com/reference/yaml/#function-secure-secrets

Of course, you'd need to rearrange your code so as to fit with the folder organization of a typical OpenFaas project.

### Google Cloud

Cloud Function doesn't require you to import any lib or implement a main function. However, you do need to implement a specific handler. The Handler for Google Cloud is into the file *google_cloud.go*.

Google Cloud supports various kind of triggers (Cloud Storage, Pub/Sub, etc.). For this example, we will use an HTTP trigger. We can pass the environment variables (in our case everything needed to connect to our MySQL instance: MYSQL_*) as a local YAML file or as key/value list of the argument `--set-env-vars`

    gcloud functions deploy GetCounter --runtime go113 --trigger-http --allow-unauthenticated --env-vars-file=mysetenv.yaml

Tip:the YAML file is expecting to contain only string values, so put integer literals in string like with my example with the MySQL port number.

Copy/paste REGION and PROJECT_ID from the terminal and Invoke the function:

    https://REGION-PROJECT_ID.cloudfunctions.net/GetCounter

For example: https://us-central1-awesomeproject.cloudfunctions.net/GetCounter

You can delete the function after your tests:

    gcloud functions delete HelloGet 

### What about AWS Lambda?

AWS requires additionnal steps such as implementing a main function waiting for a trigger (watchdog) and to include "github.com/aws/aws-lambda-go/lambda". For more information: https://docs.aws.amazon.com/lambda/latest/dg/golang-handler.html

### And Azure Functions?

**As of April 2020, the support of Go is still experimental**

Azure doesn't require you to include any lib. However, you need to implement a main function with a standard HTTP Serve Mux. See this example: https://github.com/Azure-Samples/functions-custom-handlers/blob/master/go/GoCustomHandlers.go

And you must provide a configuration file. Here an example for a Function returning a string: https://github.com/Azure-Samples/functions-custom-handlers/blob/master/go/HttpTriggerStringReturnValue/function.json

## About concurrent programing

In this example we have two examples of programing model:
 - Sequential programing where we read and update the database in a sequence.
 - Concurrent programing where we read and update the database more of less in parallel (remember that a serverless function is running on a VM with 1VCPU most the time, but while we are waiting for the result of the SQL query our code can do something else).

The example is provided with a benchmark that you could run on your computer:

    $ go test -bench=.
    BenchmarkSequential-4                153           7940957 ns/op
    BenchmarkConcurrent-4                231           5174905 ns/op

By default, the benchmark tool execute the test function in a loop for 1s. 
So we can interpret the result by saying that :
 - the sequential function is executed 153 times during a second
 - the concurrent function is executed 231 times during a second

The number for ns/op gives you the approximate time taken by one iteration.
So given the pricing for serverless functions, you should apply this approach.
