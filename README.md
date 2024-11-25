# netangeslgo

This is partial implementation of netangels.ru dns provider's API. Any contribution is welcome in the 
form of PR's. Further documentation of the API can be found [here](https://api.netangels.ru/modules/gateway_api.api.dns.records/)).

## Usage 
Add this repository as go dependency.

``` go
import (
	"github.com/navigatore300/netangelsgo"
)
```
Create a new client with your API key.
``` go
client := CreateNetangelsClient("accountName", "apiKey")
```
Use the client to interact with the API.
``` go
// Get record for a domain
records, err := client.GetRecord("example.com")
```

## Implemented methods
- [x] GetRecord
- [x] AddRecord
- [x] RemoveRecord
- [x] UpdateRecord
- [x] UpdateDDNS

*Основа взята с https://github.com/RunnerM/simply-com-client*
