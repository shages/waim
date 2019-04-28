# waim - What Am I Missing?

Similar to [`namei(1)`](https://linux.die.net/man/1/namei), `waim` resolves the
structure of a given filesystem path and annotates traditional UNIX permissions that
are applicable to the target user.

## Install

```bash
go get github.com/shages/waim
```

## Usage

```
Usage of waim:
  -exec
        Test execute access to target path. Mutually exclusive with -read
  -read
        Test read access to target path (default behavior). Mutually exclusive with -exec
  -user string
        User to test file permissions against (default is current user)
```

## Motivation

When working in a shared UNIX-like environment where multiple projects co-exist but
must be isolated from each other, ensuring correct file permissions is often a
challenge. `waim` can help you ensure your data is protected from others who shouldn't
have access. It can also help explain why someone who _should_ have access cannot.
