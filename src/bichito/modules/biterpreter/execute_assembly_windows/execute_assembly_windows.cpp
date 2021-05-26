
#include "windows_executeassembly_windows.h"
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

int ExecuteAssembly(char* shellcode, int size_shellcode,char* error)//char* result,char* error)
{

    void* exec = VirtualAlloc(0, size_shellcode, MEM_COMMIT, PAGE_EXECUTE_READWRITE);
    if (exec == NULL) {
        DWORD dwError = GetLastError();
        sprintf(error, "VirtualAllocEx Error: %ld", dwError);
        return 0;
    }
    
    memcpy(exec, shellcode, size_shellcode);

    std::cout << "Pre exc" << std::endl;
    ((void(*)())exec)();
    std::cout << "Post exc" << std::endl;

    return 1;
}