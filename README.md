**_Note:_** Because I'm running my own servers for serveral years, main development is done at at https://git.ypbind.de/cgit/redfish-tool/

----

# Preface
The [Redfish REST API](https://redfish.dmtf.org/) defines a de-facto standard for a defined access
to the management boards and information / operation on servers.

# Command line parameters
The general usage of `redfish-tool` is:

```
redfish-tool <global_options> <command> <command_options>
```

## Global options
Global options apply to all commands. Supported global options are:

| *Option* | *Description* | *Comment* |
|:---------|:--------------|:----------|
| `--ask` | Ask for password | Mutually exclusive with `--password` and `--password-file=` |
| `--debug` | Show debug information | :heavy_exclamation_mark: ***This will leak login credentials in the output*** :heavy_exclamation_mark: |
| `--format=<fmt>` | Output format | Valid values for `<fmt>` are: |
|                  |               |  `text` (*this is the default*) |
|                  |               |  `json` |
| `--help` | Shows the help text | |
| `--host=<host>[,<host>,...]` | Comma separated list of hosts/management boards to connect to | Assuming all listed hosts use the same user/password for authentication and the same setting for SSL verification |
| | | At least one host is mandatory |
| `--insecure` | Don't validate servers SSL certificate | |
| `--password=<pass>` | Authenticate with password `<pass>` | :heavy_exclamation_mark: *The password will show up in the process table and your shell history. Quotes and escapes may be needed depending on your shell* :heavy_exclamation_mark: |
| | | In a productive environment you should use the `--password-file` option instead |
| `--password-file=<file>` | Read password for authentication from `<file>` | Only the first line from `<file>` will be used as password |
| | | Use `-` as file name to read from standard input |
| `--port=<port>` | Connect to `<port>` | *Default:* 443 |
| | | **Note:** HTTPS will *always* be used because it is the mandatory protocol |
| `--user=<user>` | Authenticate as `<user>` | |
| `--timeout=<sec>` | HTTP connection timeout in seconds | *Default:* 60 |
| `--version` | Show version information | |

## Subcommands

### Account management
Account management is not supported by Lenovo (they provide no service endpoint) and Inspur (the provide an service endpoint but don't implement it).

**Note:** HP/HPE don't support roles but a map of predefined and immutable privileges (mostly for backward compatibility). HP/HPE privileges supported by this program are:

| *Name* | *Description* |
|:-------|:--------------|
| `iloconfig` | Allow configuration of the iLO management board |
| `login` | Allow iLO login (**required** for Redfish API access) |
| `remoteconsole` | Allow access to the remote console |
| `userconfig` | Allow creation/modification and deletion of user accounts on the iLO |
| `virtualmedia` | Allow to mount and unmount virtual media |
| `virtualpowerandreset` | Allow power management and reset of the server and iLO |

To simplify the use if HP/HPE privilege maps this programm recognizes the following aliases:

| *Name* | *HP/HPE privilege* |
|:-------|:-------------------|
| `none` | No privileges at all |
| `readonly` | `login` |
| `operator` | `login`, `remoteconsole`, `virtualmedia`, `virtualpowerandreset`|
| `administrator` | `login`, `remoteconsole`, `userconfig`, `virtualmedia`, `virtualpowerandreset`, `iloconfig` |

**Note:** DELL/EMC iDRAC uses a hardcoded, predefined number of account slots and the first slot is reserved and can't be used or modified.

#### List all users on the management board - `list-all-users`
To obtain a list of all users on the management board the `list-all-users` command can be used.
This command don't support any command specific options.

#### Get information about a particular user on the management board - `get-user`
The `get-user` command retrieves information about a particular user. The name or id of the user is mandatory.

| *Option* | *Description* | *Comment* |
|:---------|:--------------|:----------|
| `--name=<name>` | Name of the user account | `--name` and `--id` are mutually exclusive |
| `--id=<id>` | ID of the user account | `--name` and `--id` are mutually exclusive |

#### Get all roles from the service processor - `get-all-roles`
A list of all roles defined on the service processor can be obtained by the command `get-all-roles`
This command don't support any command specific options.

#### Get information about a specific role - `get-role`
Using the `get-role` command information about a particular role can be obtained.

| *Option* | *Description* | *Comment* |
|:---------|:--------------|:----------|
| `--id=<id>` | Get role information for role with ID `<id>` | |

***Note:*** *Because role names are not unique, roles can only be listed by ID instead of name.*

#### Add a new user on the service processor - `add-user`
To add a new user account on the service processor use the `add-user` command.

| *Option* | *Description* | *Comment* |
|:---------|:--------------|:----------|
| `--disabled` | New account will be created and disabled | |
| `--hpe-privileges=<privilege>[,<privilege>,...]` | Comma separated list of HP/HPE privileges | see note about HP/HPE privileges above |
| `--locked` | The new account will be created and locked | |
| `--name=<name>` | Name of the user account to create | **Mandatory** |
| `--password=<pass>` | Set password of the new account to `<pass>` | :heavy_exclamation_mark: *The password will show up in the process table and your shell history. Quotes and escapes may be needed depending on your shell* :heavy_exclamation_mark: |
| | | In a productive environment you should use the `--password-file` option instead |
| `--password-file=<file>` | Read password for the new account from `<file>` | Only the first line from `<file>` will be used as password |
| | | Use `-` as file name to read from standard input |
| `--role=<role>` | Assign new user to role `<role>` | on HP/HPE iLO use `--hpe-privileges` instead |

#### Delete a user on the management board - `del-user`
To remove an existing account on the service processor the `del-user` command can be used.

| *Option* | *Description* | *Comment* |
|:---------|:--------------|:----------|
| `--name=<name>` | Name of the user account to delete | **Mandatory** |

#### Modify an existing user on the management board - `modify-user`
A existing user account on the service processor can be modified by the `modify-user` command.

| *Option* | *Description* | *Comment* |
|:---------|:--------------|:----------|
| `--disabled` | Account will be disabled | |
| `--hpe-privileges=<privilege>[,<privilege>,...]` | Comma separated list of HP/HPE privileges | see note about HP/HPE privileges above |
| `--locked` | The account will be locked | |
| `--name=<name>` | Name of the user account to modify | **Mandatory** |
| `--password=<pass>` | Set password of the account to `<pass>` | :heavy_exclamation_mark: *The password will show up in the process table and your shell history. Quotes and escapes may be needed depending on your shell* :heavy_exclamation_mark: |
| | | In a productive environment you should use the `--password-file` option instead |
| | | If omitted the password will be asked and read from standard input |
| `--password-file=<file>` | Read new password for the account from `<file>` | Only the first line from `<file>` will be used as password |
| | | Use `-` as file name to read from standard input |
| `--rename=<newname>` | Rename user account to `<newname>` | |
| `--role=<role>` | Assign new role `<role>` to account | on HP/HPE iLO use `--hpe-privileges` instead |

#### Set a new password for an existing account - `passwd`
Using the `passwd` command the password of an existing account on the management processor can be changed.

| *Option* | *Description* | *Comment* |
|:---------|:--------------|:----------|
| `--name=<name>` | Name of the user account to change the password | **Mandatory** |
| `--password=<pass>` | Set password of the account to `<pass>` | :heavy_exclamation_mark: *The password will show up in the process table and your shell history. Quotes and escapes may be needed depending on your shell* :heavy_exclamation_mark: |
| | | In a productive environment you should use the `--password-file` option instead |
| | | If omitted the password will be asked and read from standard input |
| `--password-file=<file>` | Read new password for the account from `<file>` | Only the first line from `<file>` will be used as password |
| | | Use `-` as file name to read from standard input |

### Certificate management
Certificate management is not supported on DELL, Inspur, Lenovo and Supermicro because the don't provide the required endpoint (`/v1/redfish/SecurityService`).

**Note:** To generate and deploy certificates on DELL/EMC iDRAC use the `racadm` command

**Note:** For Supermicro upload the public and private SSL key using the web interface.

**Note:** For Lenovo upload the public and private SSL key using either the web interface or the `???` tool (*and hope the service processor
will start the webserver with the new SSL certificate or start the web server at all*)

**Note:** Inspur is a complete and utter failure because generation and uploading of a new certificate can't be done at all. (**Neither on the command line nor in the web interface.**)

#### Generate a certificate signing request - `gen-csr`
To start the generation of a new certificate signing request (and private SSL key) on the management board use the command `gen-csr`

| *Option* | *Description* | *Comment* |
|:---------|:--------------|:----------|
| `--country=<c>` | Set country name | |
| `--c=<c>` | Set country name | |
| `--state=<s>` | Set state name | |
| `--s=<s>` | Set state name | |
| `--locality=<l>` | Set locality/city name | |
| `--l=<l>` | Set locality/city name |
| `--organisation=<o>` | Set organisation name | |
| `--o=<o>` | Set organisation name |
| `--organisational-unit=<ou>` | Set organisational unit name | |
| `--ou=<ou>` | Set organisational unit name | |
| `--common-name=<cn>` | Set the common name | The hostname will be used if not set |
| `--cn=<cn>` | Set the common name | The hostname will be used if not set |

Depending on the service processor hardware the generation of the certificate signing request and the private SSL key can take a while (up to serveral minutes).

**Note:** If both the full option (e.g. `--organisational-unit`) and the abbreviation (e.g. `--ou`) are used the value of the abbreviation will be included in the certificate signing request.

**Note:** HP/HPE requires **all** of country (`C`), common name (`CN`), organisation name (`O`), organisational unit (`OU`), city or locality (`L`) and state (`S`) to be set in the certificate signing request.

**Note:** Huawei don't allow the forward slash (`/`) in **any** of the fields (even if escaped)

#### Fetch genearated certificate signing request - `fetch-csr`
Once the generation of the certificate signing request (and the private SSL key) has been finnished on the service processor it can be fetched by the `fetch-csr` command.

This command don't support any command specific options.
The content of the certificate signing request is printed to the standard output.

##### Import the public key of the SSL certificate - `import-cert`
After the certificate authority produced the public key of the SSL certificate by signing the certificate signing request the `import-cert` command can be used to upload the new public key to the service processor.

| *Option* | *Description* | *Comment* |
|:---------|:--------------|:----------|
| `--certificate=<file>` | Read public SSL key from `<file>` | Use `-` to read the data from standard input |
| | | **Mandatory** |

### Service processor operations
#### Get list of all management boards - `get-all-managers`
The Redfish standard allows for multiple management boards. To get the list of all management boards use the command `get-all-managers`

This command don't support any command specific options.

#### Get information about a specific management board - `get-manager`
The command `get-manager` lists information about a specific management board.

| *Option* | *Description* | *Comment* |
|:---------|:--------------|:----------|
| `--id=<id>` | Get information about management board with ID `<id>` | `--id` and `--uuid` are mutually exclusive |
| `--uuid=<uuid>` | Get information about management board with UUID `<uuid>` | `--id` and `--uuid` are mutually exclusive |

#### Reset service processor - `reset-sp`
To reset the service processor the command `reset-sp` can be used.

This command don't support any command specific options.

### Server/system operations
#### Get list of all systems - `get-all-systems`
The Redfish standard allows for multiple systems for each management board. To get the list of all systems the `get-all-systems`
is used.

This command don't support any command specific options.

#### Get information about a specific system - `get-system`
The `get-system` command retrieves information about a specific system.

| *Option* | *Description* | *Comment* |
|:---------|:--------------|:----------|
| `--id=<id>` | Get information about a system with ID `<id>` | `--id` and `--uuid` are mutually exclusive |
| `--uuid=<uuid>` | Get information about a system with UUID `<uuid>` | `--id` and `--uuid` are mutually exclusive |

#### Set power state of a systeme - `system-power`
The power state of a specific system can be set by using the `system-power` command.

| *Option* | *Description* | *Comment* |
|:---------|:--------------|:----------|
| `--id=<id>` | Set power state of the system with ID `<id>` | `--id` and `--uuid` are mutually exclusive |
| `--state=<state>` | Set the power state to `<state>` | for the name of the supported power states see the table below |
| `--uuid=<uuid>` | Set power state of the system with UUID `<uuid>` | `--id` and `--uuid` are mutually exclusive |

Names of the power states vary by vendor. Known states are:

| *Vendor* | *States* |
|:---------|:---------|
| DELL | `On`, `ForceOff`, `GracefulRestart`, `GracefulShutdown`, `PushPowerButton`, `Nmi` |
| HPE | `On`, `ForceOff`, `ForceRestart`, `Nmi`, `PushPowerButton` |
| Huawei | `On`, `ForceOff`, `GracefulShutdown`, `ForceRestart`, `Nmi`, `ForcePowerCycle` |
| Inspur | `On`, `ForceOff`, `GracefulShutdown`, `GracefulRestart`, `ForceRestart`, `Nmi`, `ForceOn`, `PushPowerButton` |
| Lenovo | `Nmi`, `ForceOff`, `ForceOn`, `GracefulShutdown`, `ForceRestart` |
| Supermicro | `On`, `ForceOff`, `GracefulShutdown`, `GracefulRestart`, `ForceRestart`, `Nmi`, `ForceOn` |

### License operations
**Note:** At the moment only HP/HPE is supported.

#### Get installed licenses - `get-license`
To get the list of installed licenses the `get-license` command can be used.

| *Option* | *Description* | *Comment* |
|:---------|:--------------|:----------|
| `--id=<id>` | Get license from management board identified by ID `<id>` | `--id` and `--uuid` are mutually exclusive |
| `--uuid=<uuid>` | Get license from management board identified by UUID `<uuid>` | `--id` and `--uuid` are mutually exclusive |

**Note:** Newer iLO versions will obfuscate the first bytes of the license key. This behavior can't be changed.

#### Add a license key - `add-license`
A new license can be added to the management board using the `add-license` command (e.g. on system board change).

| *Option* | *Description* | *Comment* |
|:---------|:--------------|:----------|
| `--id=<id>` | Add license to management board identified by ID `<id>` | `--id` and `--uuid` are mutually exclusive |
| `--license=<lic>` | License key to add to management board | This will expose the license key to the process table and the shell history |
| `--license-file=<file>` | Add the license key from a file | |
| `--uuid=<uuid>` | Add license to management board identified by UUID `<uuid>` | `--id` and `--uuid` are mutually exclusive |

# Vendor compatibility
Although the Redfish standard describes access to the API and the endpoints involved some vendors
omit mandatory endpoints (e.g. [Lenvo](https://www.lenovo.com/us/en/) for the `/v1/redfish/AccountService`), 
declare endpoints but don't implement the declared endpoints (e.g. [INSPUR](https://en.inspur.com/) for the `/v1/redfish/AccountService`) or just fail to understand the JSON simple format and return data structures incompatible with the Redfish specification (e.g. [INSPUR](https://en.inspur.com/) returns an array of strings instead of an array of `Member` objects on `/v1/redfish/Chassis`)

Additionally some vendors rely on vendor specific methods to implement a function, most notably the certificate handling for [Supermicro](https://www.supermicro.com/en/home/) and [DELL](https://www.dell.com/en-us).

Some vendors use their own vendor specific extensions, e.g. the use of a privilege map instead of roles for [HPE](https://www.hpe.com/) for backward compatibility.

The table below lists known limitations:

| *Command* | *DELL* | *HPE* | *Huawei* | *INSPUR* | *Lenovo* | *Supermicro* |
|:----------|:------:|:-----:|:--------:|:--------:|:--------:|:------------:|
| `get-all-users` | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark: | :no_entry: | :no_entry: | :heavy_check_mark: |
| `get-user` | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark: | :no_entry: | :no_entry: | :heavy_check_mark: |
| `get-all-roles` | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark: | :no_entry: | :no_entry: | :heavy_check_mark: |
| `get-role` | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark: | :no_entry: | :no_entry: | :heavy_check_mark: |
| `add-user` | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark: | :no_entry: | :no_entry: | :heavy_check_mark: |
| `del_user` | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark: | :no_entry: | :no_entry: | :heavy_check_mark: |
| `passwd` | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark: | :no_entry: | :no_entry: | :heavy_check_mark: |
| `modify-user` | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark: | :no_entry: | :no_entry: | :heavy_check_mark: |
| `gen-csr` | :no_entry: | :heavy_check_mark: | :heavy_check_mark: | :no_entry: | :no_entry: | :no_entry:  |
| `fetch-csr` | :no_entry: | :heavy_check_mark: | :heavy_check_mark: | :no_entry: | :no_entry: | :no_entry:  |
| `import-cert` | :no_entry: | :heavy_check_mark: | :heavy_check_mark: | :no_entry: | :no_entry: | :no_entry:  |
| `get-all-managers` | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark: |
| `get-manager` | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark: |
| `reset-sp` | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark: |
| `get-all-systems` | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark: |
| `get-system` | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark: |
| `system-power` | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark: |
| `get-license` | :no_entry: | :heavy_check_mark: | no additional licenses needed | no additional licenses needed | :no_entry: | no additional licenses needed |
| `add-license` | :no_entry: | :heavy_check_mark: | no additional licenses needed | no additional licenses needed | :no_entry: | no additional licenses needed |

----

# Licenses
## redfish-tool
This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.

## go-redfish (https://git.ypbind.de/cgit/go-redfish)
This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.

## logrus (https://github.com/sirupsen/logrus)
The MIT License (MIT)

Copyright (c) 2014 Simon Eskildsen

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.

## x/crypto (https://golang.org/x/crypto/)
Copyright (c) 2009 The Go Authors. All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are
met:

   * Redistributions of source code must retain the above copyright
notice, this list of conditions and the following disclaimer.
   * Redistributions in binary form must reproduce the above
copyright notice, this list of conditions and the following disclaimer
in the documentation and/or other materials provided with the
distribution.
   * Neither the name of Google Inc. nor the names of its
contributors may be used to endorse or promote products derived from
this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
"AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

