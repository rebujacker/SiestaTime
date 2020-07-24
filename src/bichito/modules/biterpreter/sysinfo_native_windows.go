// +build windows

package biterpreter

import (
  "bichito/modules/biterpreter/sysinfo_native_windows"
)


/*
Description: Sysinfo --> Windows. Retrieve Operating System key information from the foothold.
Flow:
A.Use C++ to retrieve this key information: sysinfo_native_windows.cpp
Note: These libraries need to be extracted from a darwin dev. kit, and are compiled with mingw32
*/
func Sysinfo() (bool,string){

  error,result := sysinfo_native_windows.SysinfoNativeWindows()
  return error,result
}