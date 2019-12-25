// +build windows

package biterpreter

import (
  "bichito/modules/biterpreter/sysinfo_native_windows"
)



func Sysinfo() (bool,string){

  error,result := sysinfo_native_windows.SysinfoNativeWindows()
  return error,result
}