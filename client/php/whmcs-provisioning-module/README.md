## Requirements
- PHP 7.0 or higher
- grpc extension

### Webserver
This module is expected to work on Apache.
For use with other web servers such as NGINX you have to mimic the .htaccess file(s) behaviour for security reasons

## Install
1. Extract rdpcloud.zip archive into whmcs/modules/servers/
2. Create server in WHMCS (Make sure to check the SSL/TLS option)
3. Create product in WHMCS
    1. Set Local Groups (Default: Users and Remote Desktop Users)
    2. Set Disk Quotas. Format: VolumePath, Quota Threshold, Quota Limit (example: C:\\, 1.5GB, 2GB)
    3. Enable Single Sign-On and enter the RDP Control Panel URL (example: https://1.2.3.4:5028)

## Multiple Servers
To use the RDPCloud module with multiple servers you have to concatenate all the servers certificates into one single file. Example: `cat server-cert-1.pem server-cert-2.pem server-cert-3.pem > server-cert.pem`

The certificate file should be placed at `cert/server-cert.pem`

## Notes
- WHMCS Terminate function will delete the user and his profile directory. Only the profile directory will be deleted, if the user owns other files in other locations they will remain untouched and the terminate function may fail
- It's recommended to have only one volume in the system
- Make sure the quota is enabled (enforced)
- This module is expected to work on case insensitive file systems (That's normally the default behaviour in Windows)

## FREE INSTALLATION:
You can get a free installation for your first server via AnyDesk. Read /docs/SUPPORT.md file
