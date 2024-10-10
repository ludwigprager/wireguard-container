# WireGuard Container

This repository provides a quickstart guide and configuration for setting up a WireGuard VPN server using Docker. WireGuard is a modern, high-performance VPN that is easy to configure and deploy.

## Features

- Easy setup with Docker
- Pre-configured WireGuard settings
- Quickstart scripts for starting/stopping the server
- Example configuration included

## Installation

To get started, you need to have Docker and Docker Compose installed on a linux machine.

1. Clone the repository:

    ```bash
    git clone https://github.com/ludwigprager/wireguard-container.git
    cd wireguard-container
    ```

2. create your config file  

    ```
    cp wg.yaml-example wg.yaml
    ```

3. Create Your Keys  
Replace every fields marked as 'XXXXXXXXXXXXXXXXX' in the example with keys your keys.  
Use the following command to print a key:

    ```
    wg genkey 
    ```

5. Create the config files  
The following command will create all required config file in a new folder `config-files':
    ```
    ./create-configs.sh 
    ```


4. Starting the Server  
To start the WireGuard server, run:

    ```bash
    ./start.sh
    ```

5. Add a mobile client
On an Android or Iphone install the wireguard mobile app and scan the QR code that was generated. You find it in the folder `config-files/` that was created by create-config.sh script.


