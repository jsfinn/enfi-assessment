# Enfi assessment - File Monitor

This project is a file monitor that watches a directory for changes and logs them to a file.

## Building and running

This project can be run without the need to explicitly build

`go run ./...`


## Overview of approach

This project contains an application to monitor a group of file ids and to "transfer" them if they are modified.  Per the instructions provided and ther communication with Scott via email, the following assumptions and descisions were made:
- In order to maintain the appropriate abstractions, the [monitor](monitor/monitor.go) is instantiated with implementations of the following abstractions:
  - [api](monitor/api.go) - An API as described in the dependencies section of the instructions
  - [cache](monitor/cache.go) - A cache to store the history of the files that have been processed 
-  *implementations of data input and data output are expected (ie: the watch list, the files themselves). Ultimately, we are most interested in the algo for managing very large watch lists that are both files and directory ids that obfuscate the files within*.  The files themselves are just represented by in memory data structures hidden in the [file_provider](mock/file_provider.go). The watchlist is automatically generated using that set of files. 
- History does not persist from run to run.  The cache is in memory and will be lost when the application is stopped.  This is a design decision to keep the application simple.  In a real world scenario, the cache would be either persisant or another service.  As a consequence, the application will always copy the all files on application startup, since it assumes they don't have a history.
- The application reads the initial file system structure, the watchlist, and the mutations from the file called "testdata.json".  A testdata generator is included in the mock package.  The filename is provided to the application via the configuration file.
- *"We would expect to be able to run this application locally and see output in real-time, such as files being processed, or watched."*  Any call to `Api.Copy` will be logged to the console.
- re: transferring a file: "Feel free to mock this step. A simple output that simulates that step is a-ok." Since "copy" and "transfer" are used interchangeably, the `Api.Copy` method is used to simulate the transfer of a file. 

## Monitoring Algorithm

The [monitor](monitor/monitor.go) itself is a service that accepts a list of file ids to monitor, an [Api](monitor/api.go) implementation, and a [Cache](monitor/cache.go) implementation.  All interactions with the file system are done through the Api.  

Periodically, the monitor's `EvaluateWatchlist` method is called.  This method will iterate over the watchlist to find the files that need to be evaluated.  These files are either:
- Files that have been explicity defined in the watchlist
- Files that reside in directories that have been explicity defined in the watchlist. * 

Any file that needs to be evaluated will be pushed onto the `evalutionChannel`, so that the evaluation can be done concurrently.  The evaluation is done by the `evaluateMetada` method.  

This method will:
- Retrieve the previous state of the file from the cache
- compare the previous state to the current state of the file
- If the file has been modified, it will be using the `Api.Copy` method.  This api call is not implemented to be asynchronous, but in a real world scenario, it would be.


*a note about directory recursion:
Files may be nested in subdirectories, therefore any directory that is on the watchlist needs to be fully evaluated. In the case that a subdirectory is already on a watchlist, it will not be evaluated again.  The api call to get a directories children returns the metadata for those children, so we don't need to call `api.getMetadata()` on those children.  

## Usage Examples

As configured, running the application will use the [testdata.json](testdata.json) for the initial state of the file system, the watchlist, and scheduled mutations.  The output looks like this:

```
$ go run ./...
2024/10/01 15:39:27 Copying file  file1  lastModified 1727811567346  version  1
2024/10/01 15:39:27 Copying file  file2  lastModified 1727811567346  version  1
2024/10/01 15:39:27 Copying file  file3  lastModified 1727811567346  version  1
2024/10/01 15:39:27 Copying file  file4  lastModified 1727811567346  version  1
2024/10/01 15:39:27 Copying file  file7  lastModified 1727811567346  version  1
2024/10/01 15:39:27 Copying file  file5  lastModified 1727811567346  version  1
2024/10/01 15:39:27 Copying file  file6  lastModified 1727811567346  version  1
2024/10/01 15:39:28 Copying file  file3  lastModified 1727811568347  version  2
2024/10/01 15:39:28 Copying file  file4  lastModified 1727811568347  version  2
2024/10/01 15:39:29 Copying file  file5  lastModified 1727811569348  version  2
2024/10/01 15:39:29 Copying file  file6  lastModified 1727811569348  version  2
2024/10/01 15:39:31 Copying file  file7  lastModified 1727811571349  version  2
2024/10/01 15:39:31 Copying file  file5  lastModified 1727811571349  version  3
```