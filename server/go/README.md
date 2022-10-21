## Requirements
- (none)

## Usage
1. Run the executable from CMD with the argument `-h`
    Example: `rdpcloud-server-windows.exe -h`

## Install As A Service (Recommended)
1. Download [NSSM - the Non-Sucking Service Manager](https://nssm.cc/) (2.24-101 or newer)
2. Run NSSM
3. Install RDPCloud Server as a service

## Notes
- The files should be placed in a protected location (where other users do not have access to)
- You may want to restrict access to Users (prevent users from writing directly to C:\\, etc.)
- You may want to delete/disable $RECYCLE.BIN folder(s) to avoid quota issues
- You may want to enable quota (enforced) and set default quotas (can be done using the RDPCloud Client)
- Ratelimit can be done using a reverse proxy (envoy, nginx, etc.)

## FREE INSTALLATION:
You can get a free installation for your first server via AnyDesk. Read /docs/SUPPORT.md file
