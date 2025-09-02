:: Users logged into the slicer will need to specify their unique user folder id
:: replace default below with the id from:
:: 
:: OrcaSlicer
::     Mac: /Library/Application Support/OrcaSlicer/user/USERID#
::     Linux: /.config/OrcaSlicerOrcaSlicer/user/USERID#
::     Windows: /AppData/Roaming/OrcaSlicer/user/USERID#
:: 
:: CrealityPrint
::     Mac: /Library/Application Support/Creality/Creality Print/6.0/user/USERID#
::     Linux: /.config/Creality/Creality Print/6.0/user/USERID#
::     Windows: /AppData/Roaming/Creality/Creality Print/6.0/user/USERID#
:: 
:: Not logged in: 'default'

:: Set the Slicer UserID
set USER_ID=default

:: Override the ssh username & password                                                                                                                                                             
set SSH_USER=root                                                                                                                                                                                
set SSH_PASSWORD=creality_2024                                                                                                                                                                   
                                                                                                                                                                                                    
:: Set the Printer IP                                                                                                                                                                               
set PRINTER_IP=192.x.x.x                                                                                                                                                                            
                                                                                                                                                                                                    
:: Set the slicer you want to sync from 'orca' or 'creality'                                                                                                                                        
set SLICER=orca

"C:\Users\%userprofile%\Downloads\filament-sync-tool.exe" --printer-ip %PRINTER_IP% --user %SSH_USER% --password %SSH_PASSWORD% --slicer %SLICER% --userid %USER_ID%

:: Add a pause to check logs if anythings goes wrong in the execution, uncomment the line below.
:: pause
