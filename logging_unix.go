// +build !windows

package logging
// Contains Unix-specific definitions.

import (
  "os"
)

// Stdnull redirects to the "Null" output file. Use this if you want to discard output of specific log levels.
var (
  Stdnull, _ = os.OpenFile("/dev/null", os.O_WRONLY, 0200)
)
