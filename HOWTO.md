# GoRODS Usage / Examples

This document goes over a few common iRODS tasks that could be implemented with GoRODS.

### 1. How do I read data from a file stored in iRODS?

First we need to create our bare bones .go file, import the GoRODS package, and setup the client struct. Info on ConnectionOptions [can be found in the documentation](https://godoc.org/gopkg.in/jjacquay712/GoRODS.v0#New).

```go
package main

import (
	"fmt"
	"log"
	"github.com/jjacquay712/GoRODS"
)

func main() {
	
	client, conErr := gorods.New(gorods.ConnectionOptions{
		Type: gorods.UserDefined,

		Host: "localhost",
		Port: 1247,
		Zone: "tempZone",

		Username: "rods",
		Password: "password",
	})

	// Ensure the client initialized successfully and connected to the iCAT server
	if conErr != nil {
		log.Fatal(conErr)
	}

	// All example code in this document is written as if it was placed in this context / func scope, with the client already setup.

}

```

**Remember that the following examples leave out the import statement, main() function declaration, and client struct setup for the sake of brevity.**

This next example opens a connection to the iRODS server and fetches refrence to a data object, using "/tempZone/home/rods/hello.txt" as the data object's path. It then prints the contents. Read() returns a byte slice ([]byte) so it must be converted to a string.

**Example:**
```go

// Open a data object reference for /tempZone/home/rods/hello.txt
if openErr := client.OpenDataObject("/tempZone/home/rods/hello.txt", func(myFile *gorods.DataObj, con *gorods.Connection) {

	// read the contents
	if contents, readErr := myFile.Read(); readErr == nil {
		fmt.Printf("hello.txt file contents: '%v' \n", string(contents))
	} else {
		log.Fatal(readErr)
	}

}); openErr != nil {
	log.Fatal(openErr)
}

```

**Output:**

```
hello.txt file contents: 'Hello, World!' 
```

### 2. Can I selectively read sections of a file stored in iRODS (seek certain byte range)? If so, how?

This example is very similar to the one above, except it starts reading at an offset of 7 bytes, e.g. lseek(7) and reads the next 6 bytes. See [ReadChunk](https://godoc.org/gopkg.in/jjacquay712/GoRODS.v0#DataObj.ReadChunk) if you want to read the entire file in chucks.

**Example:**

```go

// Open a data object reference for /tempZone/home/rods/hello.txt
if openErr := client.OpenDataObject("/tempZone/home/rods/hello.txt", func(myFile *gorods.DataObj, con *gorods.Connection) {

	// Read 6 bytes starting with an offset of 7 bytes
	// myFile is closed automatically when this function goes out of scope
	if contents, readErr := myFile.ReadBytes(7, 6); readErr == nil {
		fmt.Printf("hello.txt file contents: '%v' \n", string(contents))
	} else {
		log.Fatal(readErr)
	}

}); openErr != nil {
	log.Fatal(openErr)
}

```

**Output:**

```
hello.txt file contents: 'World!' 
```

### 3. How do I write a file into iRODS?

There are a few ways to accomplish this, depending on whether the file (data object) already exists. This first example assumes you want to upload (iput) a new file into iRODS. To learn about the options available in DataObjOptions, [see the documentation](https://godoc.org/gopkg.in/jjacquay712/GoRODS.v0#DataObjOptions).

**Example:**

```go

// Open a collection reference for /tempZone/home/rods
if openErr := client.OpenCollection(gorods.CollectionOptions{
	Path: "/tempZone/home/rods",
}, func(col *gorods.Collection, con *gorods.Connection) {

	// Put local file hello.txt to the collection, using default options
	myFile, putErr := col.Put("hello.txt", gorods.DataObjOptions{})
	if putErr == nil {
		fmt.Printf("Successfully added %v to the collection\n", myFile.Name())
	} else {
		log.Fatal(putErr)
	}

}); openErr != nil {
	log.Fatal(openErr)
}

```

**Output:**
```
Successfully added hello.txt to the collection
```

You can also write to an existing data object in iRODS. Notice that Write() accepts a byte slice ([]byte) so you must convert strings prior to passing them.

**Example:**

```go

// Open a data object reference for /tempZone/home/rods/hello.txt
if openErr := client.OpenDataObject("/tempZone/home/rods/hello.txt", func(myFile *gorods.DataObj, con *gorods.Connection) {

	// Write sentence to data object
	if writeErr := myFile.Write([]byte("Has anyone really been far as decided to use even go want to do look more like?")); writeErr == nil {
		fmt.Printf("Successfully wrote to file\n")
	} else {
		log.Fatal(writeErr)
	}
}); openErr != nil {
	log.Fatal(openErr)
}

```

**Output:**

```
Successfully wrote to file
```

### 4. How can I get a list of files in a directory in iRODS?

GoRODS makes this very simple! The first example shows how to print the collection contents using it's String() interface, and the next example illustrates an iterator.

**Example:**

```go

// Open a collection reference for /tempZone/home/rods
if openErr := client.OpenCollection(gorods.CollectionOptions{
	Path: "/tempZone/home/rods",
}, func(col *gorods.Collection, con *gorods.Connection) {
	
	// Pass the collection struct to Printf
	fmt.Printf("%v \n", col)

}); openErr != nil {
	log.Fatal(openErr)
}

```

**Output:**

```
Collection: /tempZone/home/rods
	C: pemtest
	C: source-code
	C: test
	d: hello.txt
	d: mydir1.tar
```

Collections are denoted with "C:" and data objects with "d:". Here's another example using iterators:

**Example:**

```go

// Open a collection reference for /tempZone/home/rods
if openErr := client.OpenCollection(gorods.CollectionOptions{
	Path: "/tempZone/home/rods",
}, func(col *gorods.Collection, con *gorods.Connection) {

	// Loop over the data objects in the collection, print the file name
	col.EachDataObj(func(obj *gorods.DataObj) {
		fmt.Printf("%v \n", obj.Name())
	})

	// Loop over the subcollections in the collection, print the name
	col.EachCollection(func(subcol *gorods.Collection) {
		fmt.Printf("%v \n", subcol.Name())
	})

}); openErr != nil {
	log.Fatal(openErr)
}

```

**Output:**
```
hello.txt 
mydir1.tar 
pemtest 
source-code2 
test 
```

You can also access the slices directly, and write the loops yourself:

```go

// Open a collection reference for /tempZone/home/rods
if openErr := client.OpenCollection(gorods.CollectionOptions{
	Path: "/tempZone/home/rods",
}, func(col *gorods.Collection, con *gorods.Connection) {

	objs, _ := col.DataObjs()
	for _, obj := range objs {
		// Use obj here
	} 

	cols, _ := col.Collections()
	for _, col := range cols {
		// Use col here
	}  

	// All() returns a slice of both data objects and collections combined
	both, _ := col.All()
	for _, gObj := range both {
		// Use gObj (generic object) here
	}  

}); openErr != nil {
	log.Fatal(openErr)
}

```

### 5. How can I apply metadata to a file in iRODS?

Metadata can be associated with a data object in iRODS by calling the AddMeta function and passing a Meta struct. AddMeta works for collections also. You can add multiple AVUs to a data object or collection that share an attribute name, however the values must be unique.

**Example:**

```go

// Open a data object reference for /tempZone/home/rods/hello.txt
if openErr := client.OpenDataObject("/tempZone/home/rods/hello.txt", func(myFile *gorods.DataObj, con *gorods.Connection) {

	// Add meta AVU to data object
	if myAVU, metaErr := myFile.AddMeta(gorods.Meta{
		Attribute: "wordCount",
		Value:     "2",
		Units:     "int",
	}); metaErr == nil {
		fmt.Printf("Added meta AVU to data object: %v\n", myAVU)
	} else {
		log.Fatal(metaErr)
	}

}); openErr != nil {
	log.Fatal(openErr)
}

```

**Output:**

```
Added meta AVU to data object: wordCount: 2 (unit: int)
```

### 6. How can I retrieve metadata from a file in iRODS?

Because metadata AVUs can share attribute names, when fetching, Attribute() returns a slice of AVUs:

**Example:**

```go

// Open a data object reference for /tempZone/home/rods/hello.txt
if openErr := client.OpenDataObject("/tempZone/home/rods/hello.txt", func(myFile *gorods.DataObj, con *gorods.Connection) {

	// Fetch all AVUs where Attribute = "wordCount"
	if metaSlice, attrErr := myFile.Attribute("wordCount"); attrErr == nil {

		fmt.Printf("%v \n", metaSlice)

	} else {
		log.Fatal(attrErr)
	}

}); openErr != nil {
	log.Fatal(openErr)
}

```

**Output:**
```
[wordCount: 2 (unit: int)]
```

### 7. How can I search for a file by metadata and other attributes?

You can search for data objects and collections that have a particular AVU using QueryMeta(). The syntax for the query string is identical to what you'd use with [imeta qu](https://docs.irods.org/4.1.9/icommands/metadata/).

**Example:**

```go

// Open a connection to iCAT
if openErr := client.OpenConnection(func(con *gorods.Connection) {

	if result, queryErr := con.QueryMeta("wordCount = 2"); queryErr == nil {
		fmt.Printf("%v \n", result)
	} else {
		log.Fatal(queryErr)
	}

}); openErr != nil {
	log.Fatal(openErr)
}

```

**Output:**
```
[DataObject: /tempZone/home/rods/hello.txt]
```

### 8. How do I set access controls?

Access controls can be set on data objects and collections using a few different functions (Chmod, GrantAccess). Regardless of the function you choose, there are three things you must know: the user or group you are granting the access to, the access level (Null, Read, Write, or Own), and whether or not the operation is recursive. You must pass the recursive flag to chmod on data objects, but the value isn't used for anything.

GrantAccess accepts a *gorods.User or *gorods.Group instead of a string, but it is otherwise identical to Chmod. These user and group structs can be retrieved using [Connection.Groups()](https://godoc.org/gopkg.in/jjacquay712/GoRODS.v0#Connection.Groups) / [Connection.Users()](https://godoc.org/gopkg.in/jjacquay712/GoRODS.v0#Connection.Users) (returns slice of all groups/users in iCAT) or the data object or collection [Owner() property](https://godoc.org/gopkg.in/jjacquay712/GoRODS.v0#DataObj.Owner).

**Example:**

```go

// Open a data object reference for /tempZone/home/rods/hello.txt
if openErr := client.OpenDataObject("/tempZone/home/rods/hello.txt", func(myFile *gorods.DataObj, con *gorods.Connection) {

	if chmodErr := myFile.Chmod("developers", gorods.Write, false); chmodErr == nil {
		fmt.Printf("Chmod success!\n")
	} else {
		log.Fatal(chmodErr)
	}

}); openErr != nil {
	log.Fatal(openErr)
}

```

**Output:**
```
Chmod success!
```

### 8. How do I move / copy data objects and collections on the iRODS server?

The example below only illustrates move and copy operations on data objects, but you can use the same functions on collections too. The CopyTo and MoveTo functions accept both *Collection references and path relative strings. If the target collection does not exist when copying, it will be created recursively. This does not apply to move operations. Neither functions support using ".." to represent the parent directory, this feature might be implemented later.

**Example:**

```go

// Open a data object reference for /tempZone/home/rods/hello.txt
if openErr := client.OpenDataObject("/tempZone/home/rods/hello.txt", func(myFile *gorods.DataObj, con *gorods.Connection) {

	// Copy hello.txt to /tempZone/home/rods/newfolder
	if cpErr := myFile.CopyTo("newfolder"); cpErr == nil {
		fmt.Printf("Copy success!\n")
	} else {
		log.Fatal(cpErr)
	}

	// Move hello.txt to /tempZone/home/rods/newfolder2
	if mvErr := myFile.MoveTo("newfolder2"); mvErr == nil {
		fmt.Printf("Move success!\n")
	} else {
		log.Fatal(mvErr)
	}

}); openErr != nil {
	log.Fatal(openErr)
}

```

**Output:**
```
Copy success!
Move success!
```

# Advanced Topics

This section covers topics that are helpful to know when getting into the advanced usage of GoRODS.

### The different ways you can connect to iRODS using GoRODS, and open data objects / collections

Configuring and implementing the GoRODS Client struct is the easiest and most straight forward way to use GoRODS, as illustrated in the examples above. Here is the client initialization again, for reference:

```go

client, conErr := gorods.New(gorods.ConnectionOptions{
	Type: gorods.UserDefined,

	Host: "localhost",
	Port: 1247,
	Zone: "tempZone",

	Username: "rods",
	Password: "password",
})

```

These are the three functions available as a starting point, once you've initialized the client:

1. [client.OpenDataObject](https://godoc.org/gopkg.in/jjacquay712/GoRods.v0#Client.OpenDataObject)
2. [client.OpenCollection](https://godoc.org/gopkg.in/jjacquay712/GoRods.v0#Client.OpenCollection)
3. [client.OpenConnection](https://godoc.org/gopkg.in/jjacquay712/GoRods.v0#Client.OpenConnection)

These three functions create a new connection to iRODS each time they're called. This has importance for application concurrency, see [this section](#threading--goroutine-connection-concerns) for more details. If you'd rather manage the opening and closing of connections, collections, and data objects yourself, you can initialize the connection directly:

```go

con, err := gorods.NewConnection(&gorods.ConnectionOptions{
	Type: gorods.UserDefined,

	Host: "localhost",
	Port: 1247,
	Zone: "tempZone",

	Username: "rods",
	Password: "password",
})

// Ensure we connected to the iCAT server
if err != nil {
	log.Fatal(err)
}

obj, _ := con.DataObject("/tempZone/home/rods/hello.txt")

col, _ := con.Collection(gorods.CollectionOptions{Path: "/tempZone/home/rods",})

// Or

col := obj.Col()

```

From there you are free to use the connection [as described in the documentation](https://godoc.org/gopkg.in/jjacquay712/GoRods.v0#Connection), just remember to call Disconnect when you're finished.

### Data Object Replicas

By default, when you access a slice of data objects or use a collection iterator, you will only retrieve a single reference to a particular data object. Even if the data object is replicated to multiple resource servers. You can find out which resource the data object belongs to using the [DataObj.Resource()](https://godoc.org/gopkg.in/jjacquay712/GoRods.v0#DataObj.Resource) function.

To access a data object reference for every replica (resource it's stored in), simply set the GetRepls field to true in your CollectionOptions. Here's an example:

**Example:**

```go

// Open a connection to iCAT, get a specific collection
if openErr := client.OpenCollection(gorods.CollectionOptions{
	Path:     "/tempZone/home/rods",
	GetRepls: true,
}, func(col *gorods.Collection, con *gorods.Connection) {

	fmt.Printf("%v \n", col)

}); openErr != nil {
	log.Fatal(openErr)
}

```

**Output:**

```
Collection: /tempZone/home/rods
	C: pemtest
	C: source-code2
	C: test
	d: hello.txt
	d: hello.txt
	d: mydir1.tar

```

Notice the duplicate data object "hello.txt" appears, because it is replicated to multiple resource servers. You can find out which resource the data object belongs to using the [DataObj.Resource()](https://godoc.org/gopkg.in/jjacquay712/GoRods.v0#DataObj.Resource) function. In later versions of GoRODS, duplicates might be combined into a single data object reference.


### PAM Authentication

GoRODS currently supports standard iRODS password authentication as well as PAM. You must configure a few things server-side and setup SSL certs before you use PAM with GoRODS. [See the "PAM > Server Configuration" section in the iRODS documentation](https://docs.irods.org/4.1.8/manual/authentication/#pam). You can toggle between the two authentication mechanisms by setting the AuthType field in ConnectionOptions:

**Standard Password Auth Example:**

```go

client, conErr := gorods.New(gorods.ConnectionOptions{
	Type: gorods.UserDefined,

	Host: "localhost",
	Port: 1247,
	Zone: "tempZone",

	Username: "rods",
	Password: "password",
})

```

**PAM Auth Example:**

```go

client, conErr := gorods.New(gorods.ConnectionOptions{
	Type:          gorods.UserDefined,
	AuthType:      gorods.PAMAuth,
	PAMPassFile:   ".secure-pam-session",
    PAMPassExpire: 4, // hours

	Host: "localhost",
	Port: 1247,
	Zone: "tempZone",

	Username: "rods",
	Password: "pampassword",
})

```

The PAMPassFile and PAMPassExpire fields are not required when using PAM. By default, the session is set to expire in one hour and the PAM session password is stored exclusively in memory. If you choose to specify a PAMPassFile, the session password will be cached to the file system. This increases efficiency somewhat, by removing the need to re authenticate and fetch the session password for each new connection.

### Collection Lazy Loading vs Eager Loading

When accessing a collection using GoRODS, you will sometimes need to access a sub-collection and it's contents. You can choose to either recursively load all sub-collections in the tree (eager loading, the collection you're working with being the root node), or you can lazy load sub-collections. By default, collections are lazy loaded. Here's an example of eager loading using the Recursive field of CollectionOptions.

**Example:**

```go

// Open a collection reference for /tempZone/home/rods
if openErr := client.OpenCollection(gorods.CollectionOptions{
	Path:      "/tempZone/home/rods",
	Recursive: true, // eager load the collection and it's children recursively
}, func(col *gorods.Collection, con *gorods.Connection) {
	
	// Pass the collection struct to Printf
	fmt.Printf("%v \n", col)

}); openErr != nil {
	log.Fatal(openErr)
}

```

The benefit of eager loading is the quick traversal and and access of collection contents, after initialization. However, this requires more initialization time to complete. Lazy loading is quick to initialize, but is slower as you traverse through the sub-collections.

### iRODS Tickets

You can use tickets with GoRODS by specifying the ticket string in ConnectionOptions.

**Example:**

First you must create an anonymous user for ticket-based access:

```
$ iadmin mkuser anonymous#tempZone rodsuser
```

Then generate the ticket (creates ticket with write access to the test collection):

```
$ iticket create write test
```

Now you can use the generated ticket in GoRODS, no password required. One thing to consider is eager loading collections (using the Recursive flag) while authenticating with a ticket. This will fail, because tickets set on collections don't apply to their sub-collections ([this is a bug in iRODS](https://github.com/irods/irods/issues/3299)).

```go

client, conErr := gorods.New(gorods.ConnectionOptions{
	Type: gorods.UserDefined,

	Host: "localhost",
	Port: 1247,
	Zone: "tempZone",

	Username: "anonymous",

	Ticket: "J93XyzJ6UGIXrJN",
})

// Ensure the client initialized successfully and connected to the iCAT server
if conErr != nil {
	log.Fatal(conErr)
}

// Open a collection reference for /tempZone/home/rods/test
if openErr := client.OpenCollection(gorods.CollectionOptions{
	Path:      "/tempZone/home/rods/test",
	Recursive: false, // Lazy load sub-collections
}, func(col *gorods.Collection, con *gorods.Connection) {
	
	// Pass the collection struct to Printf
	fmt.Printf("%v \n", col)

}); openErr != nil {
	log.Fatal(openErr)
}

```


### Serving iRODS data objects (files) over HTTP


Here's an example using the built-in HTTP handler interface, this is the easiest way to implement GoRODS over HTTP.

**Example:**

```go
package main

import (
	"fmt"
	"github.com/jjacquay712/GoRODS"
	"log"
	"net/http"
)

func main() {

	client, conErr := gorods.New(gorods.ConnectionOptions{
		Type: gorods.UserDefined,

		Host: "localhost",
		Port: 1247,
		Zone: "tempZone",

		Username: "rods",
		Password: "password",
	})

	// Ensure the client initialized successfully and connected to the iCAT server
	if conErr != nil {
		log.Fatal(conErr)
	}

	mountPath := "/irods/"

	// Setup the GoRODS FileServer
	fs := gorods.FileServer(gorods.FSOptions{
		Path:   "/tempZone/home/rods",
		Client: client,
		Download: true,
		StripPrefix: mountPath,
	})

	// Create the URL router
	mux := http.NewServeMux()

	// Serve the iRODS collection at /irods/
	mux.Handle(mountPath, http.StripPrefix(mountPath, fs))

	// Start HTTP server on port 6060
	log.Fatal(http.ListenAndServe(":6060", mux))

}

```

This example uses gorilla/mux as the HTTP router and manually serves the data objects (more control).

**Example:**
```go

package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jjacquay712/GoRODS"
	"log"
	"net/http"
	"strconv"
)

func main() {

	client, conErr := gorods.New(gorods.ConnectionOptions{
		Type: gorods.UserDefined,

		Host: "localhost",
		Port: 1247,
		Zone: "tempZone",

		Username: "rods",
		Password: "password",
	})

	// Ensure the client initialized successfully and connected to the iCAT server
	if conErr != nil {
		log.Fatal(conErr)
	}

	r := mux.NewRouter()

	r.HandleFunc("/rods/hello.txt", func(w http.ResponseWriter, r *http.Request) {

		// Open /tempZone/home/rods/hello.txt for download
		if openErr := client.OpenDataObject("/tempZone/home/rods/hello.txt", func(obj *gorods.DataObj, con *gorods.Connection) {

			// Setup download headers
			w.Header().Set("Content-Disposition", "attachment; filename="+obj.Name())
			w.Header().Set("Content-type", "application/octet-stream")
			w.Header().Set("Content-Length", strconv.FormatInt(obj.Size(), 10))

			// Read file from iRODS, and write to response
			obj.ReadChunk(1024000, func(chunk []byte) {
				w.Write(chunk)
			})

		}); openErr != nil {
			log.Fatal(openErr)
		}

	})

	// Start HTTP server on port 6060
	log.Fatal(http.ListenAndServe(":6060", r))

}

```

#### Threading / goroutine Connection Concerns

In the example above, you'll notice that we call client.OpenDataObject within the route handler. This is important if you plan on serving many files concurrently. Every call to OpenDataObject, OpenCollection, or OpenCollection from the client struct will open up a new network connection to iRODS. Because these connections aren't shared between goroutines in the example (goroutines being spun up for every HTTP route handler), there's no operation blocking, enabling fast simultaneous downloads. You'll probably want to use this pattern in your application.

It is possible to share connections between goroutines, however the operations will be blocked if the connection is already in use in another goroutine. This limits concurrent download performance to the transfer speed of a single iRODS network connection. It also enables the use of long-running iRODS connections, which has been discouraged by the developers. It's recommended that you follow the "icommand pattern" where you connect to iRODS, get the data you need, and disconnect immediately afterwards.

Here's an example of "blocking" code, which you should avoid using:

```diff
- Don't use the code below, for reference only.
```

**Example:**

```go

package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jjacquay712/GoRODS"
	"log"
	"net/http"
	"strconv"
)

func main() {

	connection, conErr := gorods.NewConnection(&gorods.ConnectionOptions{
		Type: gorods.UserDefined,

		Host: "localhost",
		Port: 1247,
		Zone: "tempZone",

		Username: "rods",
		Password: "password",
	})

	// Ensure the client initialized successfully and connected to the iCAT server
	if conErr != nil {
		log.Fatal(conErr)
	}

	// Open /tempZone/home/rods/hello.txt for download
	myFileDownload, fileErr := connection.DataObject("/tempZone/home/rods/hello.txt")

	// Ensure we were able to get a data object reference
	if fileErr != nil {
		log.Fatal(fileErr)
	}

	r := mux.NewRouter()

	r.HandleFunc("/rods/hello.txt", func(w http.ResponseWriter, r *http.Request) {

		// Setup download headers
		w.Header().Set("Content-Disposition", "attachment; filename="+myFileDownload.Name())
		w.Header().Set("Content-type", "application/octet-stream")
		w.Header().Set("Content-Length", strconv.FormatInt(myFileDownload.Size(), 10))

		// Read file from iRODS, and write to response
		// This will block other goroutines (HTTP handlers) using the same connection until it's finished reading 1024000 bytes
		myFileDownload.ReadChunk(1024000, func(chunk []byte) {
			w.Write(chunk)
		})
	})

	// Start HTTP server on port 6060
	log.Fatal(http.ListenAndServe(":6060", r))

	// Disconnect after the HTTP server shuts down
	if disErr := connection.Disconnect(); disErr != nil {
		log.Fatal(disErr)
	}

}

```
