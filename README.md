# IIS provider

The IIS provider enables [Terraform](https://terraform.io/)/[OpenTF](https://opentofu.org/) to manage IIS resources.

Still IIS reached his end-of-life (EOL), there are still many applications that run under IIS and managing manually or even via DevOps is really messy, organizing it with Terraform will make it cleaner!

## How does it work?

The provider relies on Powershell commands with [IIS.Administration](https://www.powershellgallery.com/packages/IISAdministration/) module executed through WinRM when managing remote servers, so be sure to have the ports 5985/5986 allowed at the remote server.

### Why powershell?

There is an available API called [IIS.Administration](https://github.com/microsoft/IIS.Administration) developed by Microsoft to enable managing IIS and relies o HTTP calls.
However, I didn't want to install it on every server of the company and some CI/CD scripts were already available deploying applications.

## Installing

> TBD

## How to use

> TBD

