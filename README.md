# EnFi assessment - File Monitor

This project is a file monitor that watches a directory for changes and logs them to a file.

## Building and running

This project can be run without the need to explicitly build

```
go run ./...
```


## Overview of approach

This project contains an application to monitor a group of file ids and to "transfer" them if they are modified.  Per the instructions provided and the communication with Scott via email, the following assumptions and decisions were made:
- In order to maintain the appropriate abstractions, the [monitor](monitor/monitor.go) is instantiated with implementations of the following abstractions:
  - [api](monitor/api.go) - An API as described in the dependencies section of the instructions
  - [cache](monitor/cache.go) - A cache to store the history of the files that have been processed 
-  *implementations of data input and data output are expected (ie: the watch list, the files themselves). Ultimately, we are most interested in the algo for managing very large watch lists that are both files and directory ids that obfuscate the files within*.  The files themselves are represented by in-memory data structures hidden in the [file_provider](mock/file_provider.go).  
- History does not persist from run to run.  The cache is in memory and will be lost when the application is stopped.  This decision was made to keep the application simple.  In a real world scenario, the cache would be either persistent or another service.  As a consequence, the application will always copy the all files on application startup, since it assumes they don't have a history.
- The application reads the initial file system structure, the watchlist, and the mutations from the file called "testdatalarge.json".  A [testdata generator](mock/generate_testdata_test.go) is included in the mock package.  The filename is provided to the application via the configuration file. 
- *"We would expect to be able to run this application locally and see output in real-time, such as files being processed, or watched."*  Any call to `Api.Copy` will be logged to the console.
- re: transferring a file: *"Feel free to mock this step. A simple output that simulates that step is a-ok."* Since "copy" and "transfer" are used interchangeably, the `Api.Copy` method is used to simulate the transfer of a file. 

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
Files may be nested in subdirectories, therefore any directory that is on the watchlist needs to be fully evaluated. A subdirectory will only be evaluated once per each `EvaluateWatchlist` call, even if it's a subdirectory of another directory on the list.  The api call to get a directories children returns the metadata for those children, so we don't need to call `api.getMetadata()` on those children.  

### Separation of concerns

The algorithm at a high level, can be thought of as a series of steps:
```
 [ get file ids to evaluate ] -> [ perform evaluation ] -> [ copy files ]
```

Ideally, we would take a pipelined approach here so that each step is running in parallel.  In the implementation, the  first two steps are pipelined using channels, but the third step is synchronous.  In a real world scenario, the `Api.Copy` method would be asynchronous as well.  In addition, we would want to have multiple processors for each step in the pipeline, so that it's not completely blocked during long running operations.  This would be a good candidate for a worker pool pattern.

### Optimization Choices

There are three clear optimization choices that should be called out:
- Getting the metadata for the watched files.  There are actually two optimizations here - the first is only processing files that need to be processed.  In the case of directories, we only want to scan them once.  We do this by keeping a list of directories that have been evaluated for the current iteration.  The other would be to make the calls to `Api.GetMetadata` in parallel.  This would reduce the latency of the application, but would require a more complex implementation.

- Transferring the files that have been modified.  This should be done in parallel.  The current implementation is synchronous, but in a real world scenario, it would be asynchronous.

- Instead of polling the file system, we could use a file system watcher.  This would effectively eliminate the calls to `Api.GetMetadata`, although it's not clear that a file system watcher would be able to provide the metadata for the files that are in subdirectories of the watched directories.  This would require further investigation.

## Usage Examples

As configured, running the application will use the [testdata.json](testdata.json) for the initial state of the file system, the watchlist, and scheduled mutations.  The output looks like this:

```
$ go run ./...
2024/10/01 18:27:58 Copying file  file1  lastModified 1727821678953  version  1
2024/10/01 18:27:58 Copying file  file2  lastModified 1727821678953  version  1
2024/10/01 18:27:58 Copying file  file3  lastModified 1727821678953  version  1
2024/10/01 18:27:58 Copying file  file4  lastModified 1727821678953  version  1
2024/10/01 18:27:58 Copying file  file7  lastModified 1727821678953  version  1
2024/10/01 18:27:58 Copying file  file5  lastModified 1727821678953  version  1
2024/10/01 18:27:58 Copying file  file6  lastModified 1727821678953  version  1
2024/10/01 18:27:59 Copying file  file3  lastModified 1727821679954  version  2
2024/10/01 18:27:59 Copying file  file4  lastModified 1727821679954  version  2
2024/10/01 18:28:00 Copying file  file5  lastModified 1727821680956  version  2
2024/10/01 18:28:00 Copying file  file6  lastModified 1727821680956  version  2
2024/10/01 18:28:02 Copying file  file7  lastModified 1727821682958  version  2
2024/10/01 18:28:02 Copying file  file5  lastModified 1727821682958  version  3
2024/10/01 18:28:03 watch Log:
2024/10/01 18:28:03 File: file6   watchtype: implicit  version: 2   status: copied
2024/10/01 18:28:03 File: file1   watchtype: explicit  version: 1   status: copied
2024/10/01 18:28:03 File: file2   watchtype: explicit  version: 1   status: copied
2024/10/01 18:28:03 File: file3   watchtype: implicit  version: 2   status: copied
2024/10/01 18:28:03 File: file4   watchtype: implicit  version: 2   status: copied
2024/10/01 18:28:03 File: file7   watchtype: explicit  version: 2   status: copied
2024/10/01 18:28:03 File: file5   watchtype: implicit  version: 3   status: copied
2024/10/01 18:28:03 get_children_calls: 10
2024/10/01 18:28:03 evaluate_watchlist_calls: 5
2024/10/01 18:28:03 metadata_retrieved_calls: 25
2024/10/01 18:28:03 copy_file_calls: 13
```


## Evaluation

The current test runs with a simulated filesystem of 10000 files and 100 directoies, a watchlist of 500 files, with update to 5000 files happening for 10 iterations with a 1 second sleep per iteration. Running on my local system, it's no surprise that the there is very little latency since everything happens in memory (the run completes in 10 seconds).  

For a watch list that includes 4 directories and 496 files, the stats look like this:
```
evaluate_watchlist_calls: 10
metadata_retrieved_calls: 5040
get_children_calls: 80
copy_file_calls: 3126
````

#### Retrieve Metadata 
As we'd expect, the 5040 `RetrieveMetadata` calls scale linearly by the number of files in the watchlist (500 * 10 iterations).  While there is no latency in this particular test, we can expect that it will be present if these calls are made to a cloud service or remote file system.  One way to mitigate this would be to make these calls in parallel.  

#### Get Children
There are 4 directories in the watchlist.  The 80 calls to `GetChildren` is the result of the 10 iterations of the 4 directories and their subdirectories.  As mentioned before, this implementation does not need to evaluate each child since the metadata is returned by the call to `GetChildren`.  Also, each directory is only evaluated once.  This is a good example of how the algorithm is optimized to avoid redundant calls.  Similar to the retrieve metadata calls, these could be made in parallel to reduce latency.

#### Copy File
Of the 50k changes that are made to the files, 3126 are copied.  This makes sense since on average, we'd expect 1/2 the files to change, but we're only monitoring 10% of them.  500 * 10 * 0.5 = 2500.  The extra 626 could be the result of the directories that are being evaluated - there are "implicit" changes to 80 files.  As copy files are usually IO bound, these could be made in parallel to reduce latency.

### Latency Reduction

Above, we identified a way to reduce latency through parallelization. Another strategy for reducing latency is batch processing.

In the current implementation, the biggest source of latency remains the individual API calls made to handle file metadata and directory contents. If batch APIs were available — allowing the system to request metadata or directory contents for multiple files at once — this would significantly reduce latency. The overall time spent making API calls could be decreased by a factor corresponding to the batch size, since fewer individual calls would be needed.  By combining parallel processing and batch API calls, the system could operate much more efficiently.

## Cloud Implementation

Of course, to really handle this at scale we could leverage cloud technologies.  Services like AWS S3 already provide a lot of this functionality out of the box in terms of monitoring a filesystem and providing notifications on change via SQS, so we could look to see if we could leverage that.  DynamoDB could be used to store the cache, and Lambda could be used to process the changes.  This would allow us to scale the processing of the changes as needed.  We could also use SQS to queue the changes and process them in parallel.  This would be a good candidate for a serverless architecture.