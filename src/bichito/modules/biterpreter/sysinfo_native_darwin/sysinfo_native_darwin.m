//Commented to avoid cgo compiling errors
/*
#import <Foundation/Foundation.h>
#import <mach-o/arch.h>
#include <stdlib.h>

int arch(char * res){
    int n;
    NXArchInfo *info = NXGetLocalArchInfo();
    NSString *typeOfCpu = [NSString stringWithUTF8String:info->description];
    char *archch = strdup([typeOfCpu UTF8String]);
    n = sprintf(res,"%s",archch);
    return n;
}

int osv(char * res) {
    int n;
    NSProcessInfo *pInfo = [NSProcessInfo processInfo];
    NSString *version = [pInfo operatingSystemVersionString];
    char *versionch = strdup([version UTF8String]);
    n = sprintf(res,"%s",versionch);
    return n;
}

*/