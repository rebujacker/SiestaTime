// +build darwin

package biterpreter

import (
    "bichito/modules/biterpreter/sysinfo_native_darwin"
)


func Sysinfo() (bool,string){

  error,result := sysinfo_native_darwin.SysinfoNativeDarwin()
  return error,result
}