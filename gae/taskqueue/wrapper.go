package taskqueue

import "google.golang.org/appengine/taskqueue"

// named queue is not supported in aetest environment,
// So function can be configurable for testing purose.

// Add is a wapper for google.golang.org/appengine/taskqueue.Add
var Add = taskqueue.Add
