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

### Pagination

For some requests, you can use pagination to get the next/previous page of results.  
Check if the response has a `.Pagination` field,
and if so, you can use the `.NextPage()` and `.PreviousPage()` methods to get the next or previous page of results.

Be also sure to check if there are more pages available by using the `.HasNext()` and `.HasPrevious()` methods.  
If you don't, for `PreviousPage()` you might get an error,
and for `NextPage()` it'll try to query future pages that may,
or may not exist (at the time when the original query was made).

```go
// Get the first two payments
payment, err := cli.PaymentService.GetAllPayment(acc.ID, pagination.Count(2))
if err != nil { panic(err) }

// Do something with payment response //

// Check if there are older payments
if payment.Pagination.HasPrevious() {
    // By default, the count parameter is retained from the previous request,
    // but you can override it using model.Pagination.SetCount()
    payment, err = cli.PaymentService.GetAllPayment(acc.ID, payment.Pagination.SetCount(5).PreviousPage())
    if err != nil { panic(err) }

    // Do something with the previous five payments // 
} else {
    fmt.Println("No previous payments available")
}
```

If there is no pagination object yet, because it's the first request, you can use the functions in the `pagination` package.


```go
count := pagination.Count(5) // Limits the number of returned elements to 5
newerThan := pagination.NewerThan(4024672) // Return elements newer than the given ID
olderThan := pagination.OlderThan(6774768) // Return elements older than the given ID

payment, err := cli.PaymentService.GetAllPayment(acc.ID, count, olderThan)
if err != nil { panic(err) }

// Do something with the 5 payments that are older than 6774768 //
```

## Rate Limiting

There is a built-in functionality, which should prevent the rate limit from being exceeded.
The code is from before I forked the project, and it doesn't seem to properly function all the time.

I since added a sort-of backoff policy rate limiter, that is triggered when a 429 status code is returned.
With every failure, the wait time is increased, up to a maximum of 12 seconds.
Therefore, if a request fails due to rate limiting, it'll try again for up to 12 seconds before returning an error.
If you don't want this functionality, you can disable it by setting `cli.DisableBackoff` to `true`.
In that which case a potential rate-limiting error will be returned immediately.

I strongly advise against disabling this functionality, as in my experience it's favorable to potentially rather wait a few seconds for a request,
than get an error and do the re-try logic yourself.
