# Gocurrency

Gocurrency is a tool that analyses the uses of message-passing concurrency of Go projects and output its results in HTML and CSV. 
The tool has been developped for the paper "An empirical study of messaging passing concurrency in Go projects".

## Getting Started

To install and run Gocurrency, follow these instructions:

* Make sure you have [Golang](https://golang.org/doc/install) installed.
* Run ``` git clone https://github.com/nicolasdilley/gocurrency_tool ```
* Run ``` cd gocurrency_tool/analyser```
* Run ```go get ```
* Run ```go build```
* Put the projects you want to analyse in ```./projects.txt``` (comes prepopulated with 865 projects)
* Run ``` .analyser ./analyser/projects.txt```

Gocurrency will create and populate ``` ./analyser/results ``` with the HTML and CSV results.
Gocurrency will overwrite any existing results that were previously in that folder.

For the file ```analyser/projects.txt```: 

* Put each project name on a seperate line.
* The format of the name of the projects in  is name of the author/name of the project.

IE: The Go github project url is ```https://github.com/golang/go``` becomes ``` golang/go ```

## Results

Gocurrency will output two types of file : 

* HTML
* CSV 

The CSVs can be used to further analyse the results found by the tool. 
The HTML files illustrate in a more user-friendly way the findings of a particular project.

## Concurrency primitives

The tool analyses the concurrency primitives of go and certain patterns. 

Here is the list of those primitive and patterns as they appear in the resulting table: 

  * "Goroutine" -> The keyword "go" 
  * "Receive" -> Reception of a value on a channel "<-channel"
  * "Send" -> Send  on a channel "channel <- value"
  * "Synchronous chan" -> the declaration of a synchronous channel "make(chan type)"
  * "Go in for" -> A goroutine spawned in a for loop
  * "Range over chan" -> A range over a channel "for val := range channel {}"
  * "Goroutine in for with constant (constant)" -> A goroutine spawned in a for loop where the boudness is statically know "for i := 0; i < constant; i++ {go function()}". The constant in given in the "Number" field
  * "Known chan length (length)" -> The declaration of a channel where its size is statically known. The size found is given in the "Number" field.
  * "Unknown chan length" -> The declaration of a channel where its size is statically unknown. 
  * "Make chan in for" -> The declaration of a channel in a for loop
  * "Array of chans" -> The declaration of an array of type chan
  * "Constant array of chans (length)" -> The declaration of an array of type chan where the size is statically know. The size found is given in the "Number" field.
  * "Slice array of chans" -> The declaration of a slice of type chan
  * "Map of chans" -> The declaration of a map of type chan
  * "Close chan" -> A close statement on a channel "close(channel)"
  * "Select (number of branch)" -> A select statement. The number of cases are given by the field "Number"
  * "Select with default (number of branch)" -> A select statement. The number of cases are given by the field "Number" ("default" case included)
  * "Assign chan in for" -> The uses of a channel in a for loop. 
  * "Channel of channels" -> The declaration of a channel of type chan "var channel chan chan type"
  * "Receive only chan (<-chan)" -> A function where one of the parameter is a receive only channel.
  * "Send only chan (chan<-)" ->  A function where one of the parameter is a send only channel.
  * "chan used as a param" ->  A function where one of the parameter is channel without restrictions.
