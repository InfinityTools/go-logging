// +build windows

package logging
// Contains Windows-specific definitions.

import (
  "os"
)

// Stdnull redirects to the "Null" output file. Use this if you want to discard output of specific log levels.
var (
  Stdnull, _ = os.OpenFile("nul", os.O_WRONLY, 0200)
)
