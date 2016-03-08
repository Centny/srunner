@echo off
set /p RUN_USR=Please enter user name:
set /p RUN_PWD=Please enter %RUN_USR% password:
reg add "HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows NT\CurrentVersion\Winlogon" /v AutoAdminLogon /t REG_SZ /d 1 /f
reg add "HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows NT\CurrentVersion\Winlogon" /v DefaultUserName /t REG_SZ /d %RUN_USR% /f
reg add "HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows NT\CurrentVersion\Winlogon" /v DefaultPassword /t REG_SZ /d %RUN_PWD% /f
set INSTALL="%SystemDrive%\Users\%RUN_USR%\AppData\Roaming\Microsoft\Windows\Start Menu\Programs\Startup\start_srun.bat"
echo @echo off >%INSTALL%
echo cd /d %~dp0 >>%INSTALL%
echo rundll32 user32.dll,LockWorkStation >>%INSTALL%
echo run_srun.bat >>%INSTALL%
echo pause >>%INSTALL%
pause
