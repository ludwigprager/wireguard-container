# WireGuard Container

This repository provides a quickstart guide and configuration for setting up a WireGuard VPN server using Docker. WireGuard is a modern, high-performance VPN that is easy to configure and deploy.

## Features

- Easy setup with Docker
- Pre-configured WireGuard settings
- Quickstart scripts for starting/stopping the server
- Example configuration included

## Installation

To get started, you need to have Docker and Docker Compose installed on a linux machine. You also need the `wg` command. On mobile clients the wireguard app must be installed ([link to Android play store](https://play.google.com/store/apps/details?id=com.wireguard.android)). On Linux desktop machines the `wq-quick` command is required.

1. Clone the repository:

    ```bash
    git clone https://github.com/ludwigprager/wireguard-container.git
    cd wireguard-container
    ```

2. Create The Main Config File  

    ```
    cp wg.yaml-example wg.yaml
    ```

3. Create Your Keys  
Replace every fields marked as 'XXXXXXXXXXXXXXXXX' in the example with keys your keys.  
Use the following command to print a key:

    ```
    wg genkey 
    ```

5. Create the Client Config Files  
The following command creates the client config files and QR-codes in a folder `config-files':
    ```
    ./create-configs.sh 
    ```


4. Start the Server  
To start the WireGuard server, run:

    ```bash
    ./start.sh
    ```

5. Add a Mobile Client  
On an Android or Iphone install the wireguard mobile app and scan the QR code. These can be found it in the folder `config-files/` that was created in step 5. .


