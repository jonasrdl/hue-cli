# hue-cli

hue-cli is a command-line tool for controlling Philips Hue smart lighting systems. 
With hue-cli, you can discover your Hue Bridge, register your app with it, and perform various operations.

## Usage
### 1. Discover Philips Hue Bridge
Use the discover command to find your Philips Hue Bridge on the local network. This command will automatically save the Bridge's IP address in the configuration file for future use.   
```
hue-cli discover
```

### 2. Register with Philips Hue Bridge
Before you can control your lights, you need to register your app with the Hue Bridge. Use the register command.   

**Press the link button on your Hue Bridge before you try to register.**

```
hue-cli register
```

### 3. List Devices (Lights)
Once registered, you can list all devices connected to the Hue Bridge.   
```
hue-cli list
```

**More options such as controlling etc. will follow in the future!**

## Discovery Process

`hue-cli` uses mDNS (Multicast DNS) to discover Philips Hue Bridges on the local network. mDNS allows devices to announce and discover services within a local network without relying on a centralized DNS server.

When you run the `discover` command, `hue-cli` sends mDNS queries to locate Philips Hue Bridges available on your local network. This process is automatic and helps you find and configure your Hue Bridge seamlessly.

If a Hue Bridge is discovered, its IP address is automatically saved in the configuration file, eliminating the need for manual configuration. This simplifies the setup and makes it convenient to control your Hue lights through the CLI.