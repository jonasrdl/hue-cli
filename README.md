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