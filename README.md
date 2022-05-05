[`pneuma` is a breath of fresh air][penuma].

`pneuma` is an experiment in learning Go with a focus on networking.  It
is heavily inspired by projects like [`nornir`][nornir]
and [`gornir`][gornir], but with a heavier focus on screen-scraping.

`pneuma` leverages [`scrapligo`][scrapligo] for interacting with network
devices.  It provides two things that `scrapligo` doesn't contain out of
the box: an inventory feature and a series of standardized structs for
working with data present on a network device.

## Project Status

`pneuma` is an experiment.  It shouldn't even be considered Alpha.  It
lacks tests, has no commitments to stability or backwards-compatibility,
and should generally be considered unreliable.

`pneuma` is not used on any production network.  It is not developed or
tested against a production network.  It isn't even developed against a
model of a production network.

## Why?

I created `pneuma` to help me learn the basics of Go.  Its scope is
fairly small, but I'm also very familiar with the subject, which means
I'm not learning about things that I don't care about yet (like
databases, APIs, etc.).

Screen-scraping isn't sexy.  Everyone wants to use APIs to do
everything.  While it's a lofty goal, my personal experience has been
that these APIs are unstable, unreliable, and change in weird ways from
one version to the next.  A network device's CLI tends to be a little
more stable, changes less often, and isn't a weird Frankenstein mash-up
of OpenConfig + vendor-native models.  Some vendors do a much better
job of supporting OpenConfig.  Some do a better job at documenting their
models.  But for those vendors that do a bad job of this, I would much
rather spend my time translating my expertise in the CLI into a library
than I would trying to dig through opaque models to find something that
would hopefully be what I'm looking for.

## pneuma's Inventory

`pneuma`'s inventory is simple and influenced by the Kubernetes API
object structure.  There are devices and groups.  While the inventory is
designed in a way that _hopefully_ supports pluggable inventory sources,
that is currently untested.  So, we'll look at everything through the
lens of the default YAML-based source.

There are `groups` and `devices`.  `groups` are completely optional:
you do not need them at all.  They exist primarily to collect common
parts of `devices`.  A Device can be a member of one or more `groups`.
When loading a YAML inventory source, each Device is iterated upon.
When a Device lists Group memberships, those `groups` are looked up,
and their data is recursively merged into the Device's data.  This is
done in order, so if you have multiple `groups` with conflicting data,
the last item in the group takes precedence.  If a particular piece of
data is present in both the Device and the Group, the Device's data
takes precendence.

`devices` and `groups` each have the following structure:

```yaml
<kind>:
  id:
    name: str
    metadata:
      labels:
        k: str
      platform: str
      connection:
        hostname: str
        username: str
        password: str
        port: int
        timeout: str
      # if `<kind>` is `devices`, there's also a `groups[]`:
      groups: [str]
    spec: {}
```

Additionally, `devices` have a `.metadata.groups[]` that lists its group
membership.

The YAML below shows a basic example.  There are two groups: one to
contain settings specific to a datacenter (`slc-dc`), and one to contain
settings that apply to all edge routers.  We specify these as group
memberships under `.devices.edgeRouter1.metadata.groups[]`.

```yaml
---
devices:
  edgeRouter1:
    name: "edge1-slc1"
    metadata:
      labels:
        zone: 2
      platform: arista_eos
      connection:
        hostname: 172.20.20.201
        username: admin
        password: admin
      groups:
        - slc-dc
        - edge-routers
    spec:
      routerId: 192.168.254.10
groups:
  slc-dc:
    name: slc-dc
    metadata:
      labels:
        dc: slc1
      connection:
        timeout: 5s
    spec:
      ntpServers:
        - 10.10.10.5
        - 10.10.20.5
  edge-routers:
    name: edge-routers
    metadata:
      labels:
        role: edge
      connection:
        port: 22
    spec:
      bgp:
        globalPrefixLimit: 950000
```

After loading, you will effectively have the following `devices`:

```yaml
devices:
  edgeRouter1:
    metadata:
      labels:
        dc: slc1
        role: edge
        zone: "2"
      platform: arista_eos
      groups:
        - slc-dc
        - edge-routers
      connection:
        port: 22
        hostname: 172.20.20.201
        username: admin
        password: admin
        timeout: 5s
    groups:
      - slc-dc
      - edge-routers
    spec:
      bgp:
        globalPrefixLimit: 950000
      ntpServers:
        - 10.10.10.5
        - 10.10.20.5
      routerId: 192.168.254.10
```

Taking inspiration from [nornir][nornir] and [gornir][gornir], you can
also filter your inventory:

```go
package main

import (
	"fmt"
	"pneuma/pkg/device"
	"pneuma/pkg/inventory"

	"gopkg.in/yaml.v3"
)

func main() {
	inv := inventory.YamlInventory{Filename: "example.yml"}
	i, err := inv.NewInventory()
	newInv, _ := i.FilterMetadata(
		&device.Metadata{
			Labels: map[string]string{
				"dc": "slc1",
			},
			Platform: "arista_eos",
		})
	y, _ := yaml.Marshal(newInv.Devices["edgeRouter1"].Spec)
	fmt.Println(string(y))
}
```

Besides the built-in filters (`FilterMetadata`, `FilterNameRegex`), you
can also pass your own `CustomFilterFunc` to
`*inventory.Inventory.Filter`.

## Sending Commands

Each `Device` should have a `.metadata.platform` and a
`.metadata.connection` to define how to interact with a device.  This
can be provided per device or in a group.  The platform should be in the
same format as [`scrapligo` expects][scrapligo-format].

Once done, you can use the `Device.SendCommands` method.  For example:

```go
package main

import (
	"fmt"
	"pneuma/pkg/device"
	"pneuma/pkg/inventory"

	"gopkg.in/yaml.v3"
)

func main() {
	i := inventory.YamlInventory{Filename: "example.yml"}
	inv, err := i.NewInventory()

	// Send some commands
	inv.Devices["edgeRouter1"].Connection.Open()
	defer inv.Devices["edgeRouter1"].Connection.Close()
	r, _ = inv.Devices["edgeRouter1"].Connection.SendCommand("show version")
	fmt.Println(r)
}
```

## The Future

In the future, I want to bake in some common parsing, similar to what
something like [NAPALM][napalm], probably into an OpenConfig-esque data
structure.  I don't particularly love YANG or how (seemingly
unnecessarily) nested some of those models are, so there will be some
deviation.  Of course, that assumes I continue working on this.

[pneuma]: https://en.wikipedia.org/wiki/Pneuma
[nornir]: https://nornir.readthedocs.io/en/latest/
[gornir]: https://pkg.go.dev/github.com/nornir-automation/gornir
[scrapligo]: https://github.com/scrapli/scrapligo
[scrapligo-format]: https://github.com/scrapli/scrapligo/blob/main/driver/core/factory.go#L14-L25
[napalm]: https://napalm.readthedocs.io/en/latest/
