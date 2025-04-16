# RDPCloud
## RDP Control Panel

### 💫 Overview
RDPCloud is a simple and innovative RDP Control Panel. Designed to simplify the process of RDP Control for hosting providers and users
#### Videos:
- https://youtu.be/fGuu0ZIOxmw
- https://youtu.be/6SUfM0VSokI

### ⭐️ Features
* #### Portable
    RDPCloud is portable, works out of the box, and requires zero configuration
* #### Secure by Design
    RDPCloud is built with security in mind, all network requests are encrypted, and every core functionality is hardened and designed based on the principles of security by design
* #### Easy to Integrate
    RDPCloud comes with free plugins for easy integration with extensible software, such as WHMCS
* #### User Management
    RDPCloud lets you manage users (create, modify, delete, etc.) and their local groups (add user to local group, remove user from local group, etc.)
* #### Quota Management
    RDPCloud lets you manage disk quotas (enable disk quota, set default disk quota, disable disk quota, etc.) and disk quota entries (set disk quota entry, delete disk quota entry, etc.)
* #### WinAPI Support
    RDPCloud can be extended to support any Windows functionality exposed via WinAPI
* #### Client/Server Approach
    RDPCloud, unlike other solutions, adopts the client/server approach, giving you more flexibility and control
* #### Developer Friendly
    RDPCloud, being developer friendly, uses Protocol Buffers, making it easier to communicate with RDPCloud Server in any supported language
* #### Cross-Platform
    RDPCloud Client is cross-platform and runs on Linux, Windows, and macOS. RDPCloud Server runs on Windows only

### 💰 Price
This was $99 per server, but now it's complelty **free** yet I have left the license functionality. You can license your servers as much as you want.

### ⚒ Requirements
- Go
- npm
- protoc
- composer

### 🚀 Build / Install
`make SERVER_NAME=RDPCloud SERVER_LOCAL_IP=1.2.3.4 SERVER_PUBLIC_IP=1.2.3.4 IS_FREE_TRIAL=FALSE`
- Replace RDPCloud with your desired server name (business name)
- Replace 1.2.3.4 with yoour server IP

A bundle file will be created on the bundle folder. From here you can follow the linked videos above.
