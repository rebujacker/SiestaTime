#define _WIN32_DCOM

#include "windows_schtasks.h"
#include <windows.h>
#include <iostream>
#include <comdef.h>
#include <lmcons.h>
#include <taskschd.h>

#pragma comment(lib, "taskschd.lib")
#pragma comment(lib, "comsupp.lib")

using namespace std;

int SchtasksOnUserLogon(char* executablePath, char* taskName,char* error)
{

	//USES_CONVERSION;

	//  ------------------------------------------------------
	//  Initialize COM.
	HRESULT hr = CoInitializeEx(NULL, COINIT_MULTITHREADED);
	if (FAILED(hr))
	{
		sprintf(error, "\nCoInitializeEx failed: %x", hr);
		return 0;
	}

	//  Set general COM security levels.
	hr = CoInitializeSecurity(
		NULL,
		-1,
		NULL,
		NULL,
		RPC_C_AUTHN_LEVEL_PKT_PRIVACY,
		RPC_C_IMP_LEVEL_IMPERSONATE,
		NULL,
		0,
		NULL);

	if (FAILED(hr))
	{
		CoUninitialize();
		sprintf(error, "\nCoInitializeSecurity failed: %x", hr);
		return 0;
	}

	//  ------------------------------------------------------
	
	//Concert Strings
	const size_t cSizeTaskName = strlen(taskName) + 1;
	wchar_t* wcTaskName = new wchar_t[cSizeTaskName];
	mbstowcs(wcTaskName, taskName, cSizeTaskName);
	
	LPCWSTR wszTaskName = wcTaskName;

	//  ------------------------------------------------------
	//  Create an instance of the Task Service. 
	ITaskService* pService = NULL;
	hr = CoCreateInstance(CLSID_TaskScheduler,
		NULL,
		CLSCTX_INPROC_SERVER,
		IID_ITaskService,
		(void**)&pService);
	if (FAILED(hr))
	{
		CoUninitialize();
		sprintf(error, "Failed to create an instance of ITaskService: %x", hr);
		return 0;
	}

	//  Connect to the task service.
	hr = pService->Connect(_variant_t(), _variant_t(),
		_variant_t(), _variant_t());
	if (FAILED(hr))
	{
		pService->Release();
		CoUninitialize();
		sprintf(error, "ITaskService::Connect failed: %x", hr);
		return 0;
	}

	//  ------------------------------------------------------
	//  Get the pointer to the root task folder.  This folder will hold the
	//  new task that is registered.

	ITaskFolder* pRootFolder = NULL;
	hr = pService->GetFolder(_bstr_t(L"\\"), &pRootFolder);
	if (FAILED(hr))
	{
		pService->Release();
		CoUninitialize();
		sprintf(error, "Cannot get Root Folder pointer: %x", hr);
		return 0;
	}

	//  If the same task exists, remove it.

	pRootFolder->DeleteTask(_bstr_t(wszTaskName), 0);

	//  Create the task builder object to create the task.
	ITaskDefinition* pTask = NULL;
	hr = pService->NewTask(0, &pTask);

	pService->Release();  // COM clean up.  Pointer is no longer used.
	if (FAILED(hr))
	{

		pRootFolder->Release();
		CoUninitialize();
		sprintf(error, "Failed to create a task definition: %x", hr);
		return 0;
	}

	//  ------------------------------------------------------
	//  Get the registration info for setting the identification.
	IRegistrationInfo* pRegInfo = NULL;
	hr = pTask->get_RegistrationInfo(&pRegInfo);
	if (FAILED(hr))
	{

		pRootFolder->Release();
		pTask->Release();
		CoUninitialize();
		sprintf(error, "\nCannot get identification pointer: %x", hr);
		return 0;
	}

	hr = pRegInfo->put_Author(_bstr_t(L"Microsoft Corporation"));
	pRegInfo->Release();
	if (FAILED(hr))
	{

		pRootFolder->Release();
		pTask->Release();
		CoUninitialize();
		sprintf(error, "\nCannot put identification info: %x", hr);
		return 0;
	}

	//  ------------------------------------------------------
	//  Create the settings for the task
	ITaskSettings* pSettings = NULL;
	hr = pTask->get_Settings(&pSettings);
	if (FAILED(hr))
	{

		pRootFolder->Release();
		pTask->Release();
		CoUninitialize();
		sprintf(error, "\nCannot get settings pointer: %x", hr);
		return 0;
	}

	//  Set setting values for the task. 
	hr = pSettings->put_StartWhenAvailable(VARIANT_TRUE);
	pSettings->Release();
	if (FAILED(hr))
	{

		pRootFolder->Release();
		pTask->Release();
		CoUninitialize();
		sprintf(error, "\nCannot put setting info: %x", hr);
		return 0;
	}

	//  ------------------------------------------------------
	//  Get the trigger collection to insert the logon trigger.
	ITriggerCollection* pTriggerCollection = NULL;
	hr = pTask->get_Triggers(&pTriggerCollection);
	if (FAILED(hr))
	{

		pRootFolder->Release();
		pTask->Release();
		CoUninitialize();
		sprintf(error, "\nCannot get trigger collection: %x", hr);
		return 0;
	}

	//  Add the logon trigger to the task.
	ITrigger* pTrigger = NULL;
	hr = pTriggerCollection->Create(TASK_TRIGGER_LOGON, &pTrigger);
	pTriggerCollection->Release();
	if (FAILED(hr))
	{

		pRootFolder->Release();
		pTask->Release();
		CoUninitialize();
		sprintf(error, "\nCannot create the trigger: %x", hr);
		return 0;
	}

	ILogonTrigger* pLogonTrigger = NULL;
	hr = pTrigger->QueryInterface(
		IID_ILogonTrigger, (void**)&pLogonTrigger);
	pTrigger->Release();
	if (FAILED(hr))
	{

		pRootFolder->Release();
		pTask->Release();
		CoUninitialize();
		sprintf(error, "\nQueryInterface call failed for ILogonTrigger: %x", hr);
		return 0;
	}


	hr = pLogonTrigger->put_Id(_bstr_t(L"Trigger1"));
	if (FAILED(hr))
	{
		sprintf(error, "\nCannot put the trigger ID: %x", hr);
		return 0;
	}

	//  Define the user.  The task will execute when the user logs on.

	//Get Actual Process Username
	TCHAR  username[UNLEN + 1];
	DWORD username_len = UNLEN + 1;
	GetUserName(username, &username_len);
	//_tprintf(TEXT("\nComputer name:      %s"), username);

	//Convert Char* to BSTR 
	const size_t cSizeUsername = strlen(username) + 1;
	wchar_t* wcUsername = new wchar_t[cSizeUsername];
	//std::wstring wc(cSize, L'#');
	mbstowcs(wcUsername, username, cSizeUsername);

	hr = pLogonTrigger->put_UserId(_bstr_t(wcUsername));
	pLogonTrigger->Release();
	if (FAILED(hr))
	{

		pRootFolder->Release();
		pTask->Release();
		CoUninitialize();
		sprintf(error, "\nCannot add user ID to logon trigger: %x", hr);
		return 0;
	}


	//  ------------------------------------------------------
	//  Add an Action to the task. This task will execute notepad.exe.     
	IActionCollection* pActionCollection = NULL;

	//  Get the task action collection pointer.
	hr = pTask->get_Actions(&pActionCollection);
	if (FAILED(hr))
	{

		pRootFolder->Release();
		pTask->Release();
		CoUninitialize();
		sprintf(error, "\nCannot get Task collection pointer: %x", hr);
		return 0;
	}

	//  Create the action, specifying that it is an executable action.
	IAction* pAction = NULL;
	hr = pActionCollection->Create(TASK_ACTION_EXEC, &pAction);
	pActionCollection->Release();
	if (FAILED(hr))
	{
		pRootFolder->Release();
		pTask->Release();
		CoUninitialize();
		sprintf(error, "\nCannot create the action: %x", hr);
		return 0;
	}

	IExecAction* pExecAction = NULL;
	//  QI for the executable task pointer.
	hr = pAction->QueryInterface(
		IID_IExecAction, (void**)&pExecAction);
	pAction->Release();
	if (FAILED(hr))
	{
		pRootFolder->Release();
		pTask->Release();
		CoUninitialize();
		sprintf(error, "\nQueryInterface call failed for IExecAction: %x", hr);
		return 0;
	}


	//Convert Char* to BSTR 
	const size_t cSizeExePath = strlen(executablePath) + 1;
	wchar_t* wcExePath = new wchar_t[cSizeExePath];
	mbstowcs(wcExePath, executablePath, cSizeExePath);

	hr = pExecAction->put_Path(_bstr_t(wcExePath));
	pExecAction->Release();
	if (FAILED(hr))
	{
		pRootFolder->Release();
		pTask->Release();
		CoUninitialize();
		sprintf(error, "\nCannot set path of executable: %x", hr);
		return 0;
	}

	//  ------------------------------------------------------
	//  Save the task in the root folder.
	IRegisteredTask* pRegisteredTask = NULL;

	hr = pRootFolder->RegisterTaskDefinition(
		_bstr_t(wszTaskName),
		pTask,
		TASK_CREATE_OR_UPDATE,
		_variant_t(L""),
		_variant_t(),
		TASK_LOGON_INTERACTIVE_TOKEN,
		_variant_t(L""),
		&pRegisteredTask);
	if (FAILED(hr))
	{

		pRootFolder->Release();
		pTask->Release();
		CoUninitialize();
		sprintf(error, "\nError saving the Task : %x", hr);
		return 0;
	}

	// Clean up
	pRootFolder->Release();
	pTask->Release();
	pRegisteredTask->Release();
	CoUninitialize();

	sprintf(error, "\n Success! Task successfully registered.");
	return 1;
}

int SchtasksDelete(char* taskName, char* error)
{

	//  ------------------------------------------------------
	//  Initialize COM.
	HRESULT hr = CoInitializeEx(NULL, COINIT_MULTITHREADED);
	if (FAILED(hr))
	{
		sprintf(error, "\nCoInitializeEx failed: %x", hr);
		return 0;
	}

	//  Set general COM security levels.
	hr = CoInitializeSecurity(
		NULL,
		-1,
		NULL,
		NULL,
		RPC_C_AUTHN_LEVEL_PKT_PRIVACY,
		RPC_C_IMP_LEVEL_IMPERSONATE,
		NULL,
		0,
		NULL);

	if (FAILED(hr))
	{
		CoUninitialize();
		sprintf(error, "\nCoInitializeSecurity failed: %x", hr);
		return 0;
	}

	//  ------------------------------------------------------

	//  ------------------------------------------------------
	//  Create an instance of the Task Service. 
	ITaskService* pService = NULL;
	hr = CoCreateInstance(CLSID_TaskScheduler,
		NULL,
		CLSCTX_INPROC_SERVER,
		IID_ITaskService,
		(void**)&pService);
	if (FAILED(hr))
	{
		CoUninitialize();
		sprintf(error, "Failed to create an instance of ITaskService: %x", hr);
		return 0;
	}

	//  Connect to the task service.
	hr = pService->Connect(_variant_t(), _variant_t(),
		_variant_t(), _variant_t());
	if (FAILED(hr))
	{
		pService->Release();
		CoUninitialize();
		sprintf(error, "ITaskService::Connect failed: %x", hr);
		return 0;
	}

	//  ------------------------------------------------------
	//  Get the pointer to the root task folder.  This folder will hold the
	//  new task that is registered.

	ITaskFolder* pRootFolder = NULL;
	hr = pService->GetFolder(_bstr_t(L"\\"), &pRootFolder);
	if (FAILED(hr))
	{
		pService->Release();
		CoUninitialize();
		sprintf(error, "Cannot get Root Folder pointer: %x", hr);
		return 0;
	}


	//Convert Char* to BSTR 
	const size_t cSizeExePath = strlen(taskName) + 1;
	wchar_t* wcExePath = new wchar_t[cSizeExePath];
	mbstowcs(wcExePath, taskName, cSizeExePath);

	//  If the same task exists, remove it.
	pRootFolder->DeleteTask(_bstr_t(wcExePath), 0);

	sprintf(error, "\n Success! Task successfully Deleted.");
	return 1;
}