/*** Copyright (c) 2016, University of Florida Research Foundation, Inc. ***
 *** For more information please refer to the LICENSE.md file            ***/

package gorods

// #include "wrapper.h"
import "C"

import (
	"fmt"
	// "io/ioutil"
	// "path/filepath"
	// "strconv"
	// "strings"
	// "time"
	// "unsafe"
)

// Client structs are used to store connection options, and instatiate connections with those options
type Client struct {
	Options    *ConnectionOptions
	ConnectErr error
}

// OpenConnection will create a new connection using the previously configured iRODS client. It will execute the handler,
// and close *Collection and *Collection automatically when your handler finishes execution.
// Operations on a single connection are queued when shared between goroutines (iRODS C API
// doesn't support concurrent operations on a single connection), so be sure to open up new connections
// for long-running and concurrent operations to prevent blocking.
func (cli *Client) OpenConnection(opts CollectionOptions, handler func(*Collection, *Connection)) error {
	if cli.ConnectErr == nil {
		if con, err := NewConnection(cli.Options); err == nil {
			col, colEr := con.Collection(opts)

			if colEr != nil {
				return newError(Fatal, fmt.Sprintf("Can't open new connection: %v", colEr))
			}

			handler(col, con)

			if er := col.Close(); er != nil {
				return er
			}
			if er := con.Disconnect(); er != nil {
				return er
			}

			return nil
		} else {
			return newError(Fatal, fmt.Sprintf("Can't open new connection: %v", err))
		}
	}

	return newError(Fatal, fmt.Sprintf("Can't open new connection: %v", cli.ConnectErr))
}

// New creates a test connection to an iRods iCAT server, and returns a *Client struct if successful.
// EnvironmentDefined and UserDefined constants are used in ConnectionOptions{ Type: ... }).
// When EnvironmentDefined is specified, the options stored in ~/.irods/irods_environment.json will be used.
// When UserDefined is specified you must also pass Host, Port, Username, and Zone. Password
// should be set unless using an anonymous user account with tickets.
func New(opts ConnectionOptions) (*Client, error) {
	cli := new(Client)

	cli.Options = &opts

	if _, err := NewConnection(cli.Options); err != nil {
		cli.ConnectErr = err
		return nil, err
	}

	return cli, nil
}