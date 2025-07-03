Heavily inspired by https://github.com/HurricanePrint/Filament-Sync  

This one allow you to not have a service on the printer and work as a single binary.  
***No need to install dependencies on the printer.***

Run the binary as a post-processing script to sync the profiles to the printer.  


Please open an issue if you see anything wrong so we can improve the tool !  

## Download
Current build release is available here:  
https://github.com/zaggash/go-filament-sync/releases/latest  

Latest dev pre-release are available here:  
https://github.com/zaggash/go-filament-sync/releases/tag/dev-main


## Usage
```
Usage of ./filament-sync-tool :
  -password string
    	Password for SSH connection to printer (default "creality_2024")
  -printer-ip string
    	IP address of the Creality printer (required)
  -slicer string
    	Specify the slicer type: 'orca' or 'creality' (default "orca")
  -user string
    	Username for SSH connection to printer (default "root")
  -userid string
    	Specify the user ID for the slicer profile folder (default "default")
```


## Creating custom filament presets
If you want to get your presets ready just copy your settings into a new custom filament profile

Right side of the Filament section

Click "Set filaments to use"

Click "Custom Filament" at the top then click "Create New"

If you don't see this option on Creality Print, go into the options/preferences and set User Role to Professional

You will need to add this into the Notes section of the filament

{"id":"","vendor":"","type":"","name":""}
The "id" should be a unique 5 digit value that you will also match with your custom RFID tags if you are using them

An id is still required even if you are not using RFID tags as the tool searches by id when updating filament settings

Here is an example

{"id":"02345","vendor":"Elegoo","type":"PLA","name":"Fast PLA"}

## RFIS to CFS Android app
Use the advanced mode to set the id as "Material Code"

--------------------------

## How to build locally :
### Build the binaries
`docker build -f ./docker/Dockerfile.build -t filament-sync-tool-builder:latest .`

### Run temp container to extract binaries
`docker create --name temp_tool_container filament-sync-tool-builder:latest`
#### Copy Linux binary:
`docker cp temp_tool_container:/app/filament-sync-tool_linux_amd64 ./filament-sync-tool_linux_amd64`
#### Copy macOS binary (Intel/AMD64):
`docker cp temp_tool_container:/app/filament-sync-tool_macos_amd64 ./filament-sync-tool_macos_amd64`
#### Copy Windows binary:
`docker cp temp_tool_container:/app/filament-sync-tool_windows_amd64.exe ./filament-sync-tool_windows_amd64.exe`

### Remove temp container
`docker rm -f temp_tool_container`
