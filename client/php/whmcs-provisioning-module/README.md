This module is expected to work on Apache.
For use with other web servers such as NGINX you have to mimic the .htaccess files behaviour for security reasons

WHMCS Terminate function will delete the user and his profile directory
Only the profile directory will be deleted, if the user have other files in other folders they will remain untouched

Recommended to use only one volume
or at least to assign only one quota volume per user
Recommended to only set the quota on the volume where the user profiles exist

Make sure the quota is enabled (enforced)

This module is expected to work on case insensitive file systems (That's normally the default behaviour in Windows)
