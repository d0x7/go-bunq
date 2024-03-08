# go-bunq

Go API client for [bunq](https://doc.bunq.com)

While this fork is still heavily based on the original [go-bunq](https://github.com/OGKevin/go-bunq)
and it's [fork](https://github.com/gjongenelen/go-bunq),
it's not a drop-in replacement as many structs have been moved to a different package
and some methods have different parameter signatures.

## Disclaimer

This was forked from a fork,
that was unmaintained for three years,
while the original repo is unmaintained for close to five years.
The code is very sparsely, if at all, documented.

I'd recommend checking out the original [bunq documentation](https://doc.bunq.com),
to find what you need and derive from there on.

## Installation

```bash
go get github.com/d0x7/go-bunq
```

## Usage

You first need to activate your API key and register your device, by creating an API Context.
You can easily create one by using `bunq.CreateContext`, which can later be re-used by using `bunq.LoadContext`.

When activating your key, you can either manually set a list of
permitted IP addresses and/or ranges, or either use `bunq.WildcardIP` to allow all IPs,
or `bunq.CurrentIP`  to allow only the current IP address.

```go
func main() {
  cli, err := bunq.CreateContext(
      context.Background(), // ctx: The context for the client.
      bunq.BaseURLSandbox,  // baseURL: Depending on the environment, should either be bunq.BaseURLProduction or bunq.BaseURLSandbox.
      "sandbox_5f6c7d8e9f0a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6q7r8s9t0u1v2w3x4y5z6a7b8c9d0e1f2", // apiKey: Your API key for accessing the bunq API.
      "go-bunq-client",       // deviceDescription: An device description that shows up in the app, to identify the API key
      bunq.WildcardIP,        // permittedIps: IP address or IP ranges allowed to use this API key. bunq.WildcardIP or bunq.CurrentIP may be used as well.
      "bunq_go_sandbox.json", // contextFile: Path to the file where the API Context will be stored.
  )
  if err != nil { panic(err) }
  // If this succeeds, the API Context has successfully been created and written to file.
  // You can now use the API client as usual, or exit.

  // Query all bank accounts
  resp, err := cli.AccountService.GetAllMonetaryAccountBank()
  if err != nil { panic(err) }

  // And print the response
  for _, r := range resp.Response {
    acc := r.MonetaryAccountBank
    fmt.Printf("Account %d has %s %s on %s\n", acc.ID, acc.Balance.Value, acc.Balance.Currency, acc.GetIBAN())
  }
  // -> Account 113131 has 19401.48 EUR on NL00BUNQ1234567890
}
```

After you created the API Context, you can now load it at any time using the `bunq.LoadContext` function.

```go
cli, err := bunq.LoadContext(context.Background(), "bunq_go_sandbox.json")
if err != nil { panic(err) }
// Again, if this succeeds, the client is now initialized and may be used as usual.
```
