# Home IO Packet monitorer

This script is used to monitor the packets from Home IO simulation and modified them.
The [Modbus interface for Home I/O simulation](https://github.com/Klagarge/Modbus2HomeIO) and his [Controller](https://github.com/Klagarge/ControllerHomeIo) have to be used.

This script have to be run on sudo mode on an MitM environment. (e.g. ARP poisoning with `ettercap`)

## Prerequisites
<p align="left">
<a href="https://www.kali.org/" target="_blank" rel="noreferrer"> <img src="https://upload.wikimedia.org/wikipedia/commons/thumb/2/2b/Kali-dragon-icon.svg/1200px-Kali-dragon-icon.svg.png" alt="kali linux" width="60" height="60"/> </a>
<a href="https://www.ettercap-project.org/" target="_blank" rel="noreferrer"><img src="https://www.kali.org/tools/ettercap/images/ettercap-logo.svg" alt="ettercap" width="60" height="60"/> </a>
<a href="https://linux.die.net/man/8/iptables" target="_blank" rel="noreferrer"><img src="https://projects.task.gda.pl/uploads/-/system/project/avatar/286/iptables-logo.png" alt="iptables" width="60" height="60"/> </a>
<a href="https://go.dev/" target="_blank" rel="noreferrer"> <img src="https://cdn.icon-icons.com/icons2/2107/PNG/512/file_type_go_gopher_icon_130571.png" alt="go" width="60" height="60"/> </a>

The following components must be installed:

- Ettercap
- iptables
- Go 1.22 or higher

Ettercap and iptables are installed by default on Kali Linux.
They can be installed on other Linux using apt:

```bash
sudo apt install ettercap-text-only iptables
```

## Usage
1. Start ARP poisoning with Ettercap (change the IP address):
   ```bash
   sudo ettercap -T -i eth0 -M arp /192.168.39.110// /192.168.37.163//
   ```
   exit with `q`
2. Redirect port for intercept packet
   ```bash
   sudo iptables -t nat -A PREROUTING -p tcp --dport 5802 -j REDIRECT --to-port 5803
   ```
3. Run the go scrypt on sudo mode
   ```bash
   sudo go run .
   ```

## How it works
Ettercap is used to perform ARP poisoning and redirect the traffic to the attacker.
ARP poising involves sending ARP messages to the targets machine, for fake their correspondent and make them send their packets to the attacker.

Iptables is used to put redirect packet to another port, so they can be modified by the Go script.

The Go script get connection from the controller and create another connection to the simulation.
This script used its own certificates, so it's work only if certificate aren't verified on controller and simulation parts.
On unencrypted packet, the script analyse and check if this is a request for the corresponding client id and adresse and put the transaction id on a list.
When this transaction id appear again, the script modify the packet to make it a response with the value the attacker want.

## Authors
- **Rémi Heredero** - _Initial work_ - [Klagarge](https://github.com/Klagarge)