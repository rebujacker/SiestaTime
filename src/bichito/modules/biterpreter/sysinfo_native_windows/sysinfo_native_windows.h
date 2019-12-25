#ifndef sysinfo_native_windows_H_
#define sysinfo_native_windows_H_

#ifdef __cplusplus
extern "C" {
#endif

	int ProcessIntegrity(char*, char*);
	int IsLocalAdmin(char*, char*);

#ifdef __cplusplus
}
#endif

#endif

