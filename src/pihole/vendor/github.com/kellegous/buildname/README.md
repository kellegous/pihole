# BuildName

Generate human names for your binaries.

For most programs, server and cli, that I create these days, I tend to use linker flags to assign the git sha to a variable. It's convenient to be able to ask a binary what source it was built from. But git shas and most versioning schemes are just so boring. So many years ago, I started generating names for my builds based on hashing the version string and selecting both an adjective and a name. I liked that the names were far more conversational than verison numbers:

> What build are you using?

> I'm using Quirky Badger

You can see this naming scheme in most of my projects. https://kellegous.com/version.

## installation

`go get -u github.com/kellegous/buildname`

## Quick Start

```go
import "github.com/kellegous/buildname"

...

name := buildname.FromVersion("e6a92ec2fe5fba022c31c32c97ea455cee4b2736")
fmt.Printf("name: %s\n", name) // Clever Turn
```