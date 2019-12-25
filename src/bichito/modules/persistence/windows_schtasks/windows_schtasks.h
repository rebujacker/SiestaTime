#ifndef windows_schtasks_H_
#define windows_schtasks_H_

#ifdef __cplusplus
extern "C" {
#endif

	int SchtasksOnUserLogon(char*,char*,char*);
	int SchtasksDelete(char*, char*);

#ifdef __cplusplus
}
#endif

#endif
