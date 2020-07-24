// +build darwin

package biterpreter

import (
    "bichito/modules/biterpreter/sysinfo_native_darwin"
)

/*
Description: Sysinfo --> Darwin. Retrieve Operating System key information from the foothold.
Flow:
A.Call the Objective C wrapped source (sysinfo_native_darwin.go/sysinfo_native_darwin.m)
	A1. Call OSX native libraries (<Foundation/Foundation.h>,<mach-o/arch.h>) to extract information
Note: These libraries need to be extracted from a darwin dev. kit, and are compiled with OSXCROSS
*/
func Sysinfo() (bool,string){

  error,result := sysinfo_native_darwin.SysinfoNativeDarwin()
  return error,result
}