:: Set PROFILE_PATH to the full path of your slicer's filament profile directory.
::
:: OrcaSlicer (Windows):
::     %APPDATA%\OrcaSlicer\user\default\filament\base
::     (replace 'default' with your user ID if logged in)
::
:: Creality Print (Windows):
::     %APPDATA%\Creality\Creality Print\6.0\user\default\filament\base
::     (replace 6.0 with your installed version)
::     (replace 'default' with your user ID if logged in)

:: Set the full profile directory path
set PROFILE_PATH=%APPDATA%\OrcaSlicer\user\default\filament\base

:: Override the ssh username & password
set SSH_USER=root
set SSH_PASSWORD=creality_2024

:: Set the Printer IP
set PRINTER_IP=192.x.x.x

"%userprofile%\Downloads\filament-sync-tool.exe" --printer-ip %PRINTER_IP% --user %SSH_USER% --password %SSH_PASSWORD% --profile-path %PROFILE_PATH%

:: Add a pause to check logs if anythings goes wrong in the execution, uncomment the line below.
:: pause
