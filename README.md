# TIM authentication package for GoLang

### Dependencies
github.com/nmcclain/ldap

### Installation
```
go get github.com/nmcclain/ldap
go get github.com/ausrasul/GO-TIM
```

Also you should have a config file (JSON formatted) with the following parameters:

	{
		"Ldap_server": "tds.domain.com",
		"Ldap_port": "389",
		"Base_dn": "ou=AppName,o=T2,DC=COM",
		"Ldap_user": "cn=APPUSER,ou=AppName,o=T2,DC=COM",
		"Ldap_pass": "AppPassword"
	}

### Usage:

You can use this package as stand-alone, see the example below.
It can be used also as a TIM driver for the User package (Look for GO-User package).

### Usage example:

```
package main

import (
   "github.com/ausrasul/GO-TIM"
)

func main(){
	user, err := tim.GetUser("username", "password")

	/* The result is:
		user map[string][]string
		example:
		user["mobile"][0] = "+46......."
		user["cn"][0] = "User Name"
		....
		user["Roles"] = ["role1", "role2", ...]
	*/
}

```
