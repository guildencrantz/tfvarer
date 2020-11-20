# `tfvarser`

Ever want to convert an HCL tfvars file to JSON? Or a JSON tfvars file to HCL?
Or merge two tfvars files? No? Good for you!

If, however, you want (or just need) to do these `tfvarser` should help.

## Installation

### Go Binary

`go get github.com/guildencrantz/tfvarser`

### Docker Image

`docker pull guildencrantz/tfvarser:1`

## Usage

`tfvarser -hcl terraform.tfvars -json new_values.json`

This will read the `terraform.tfvars`, then add the contents of
`new_values.json` (overwriting any keys that are in `terraform.tfvars` with
values from `new_values.json`, or just adding new keys) and output, to `stdout`,
the merged config (in `hcl`).

You can specify `-hcl` and `-json` as many times as you want. All files are
processed in left-to-right flag order, regardless of format.

By default the output is `hcl`, however you can output `json` with `-output
json`.

Likewise by default values from latter configs will override earlier configs,
however you can suppress this behavior (so that once a key is set it's not
overwritten) by setting `-overwrite false`.

Refer to `tfvarser -help` for current flags.

## Why?

More than once now I've found I needed to merge configs. In this case after
importing existing infrastructure into terraform I need to build the config so
that it won't just destroy all that existing infrastructure. Well, I have most
of the needed config in JSON, and the few pieces that aren't in that JSON are in
an existing `tfvars` file. Since the tooling assumes environment specific
configurations will be in `tfvars` file it was easier to write this than to
update that to handle `tfvars.json` files.

## Caveats

In addition to the normal "heh, this was written very quickly after hours, you
get what you pay for and you're on your own" it should be noted that this tool,
though it generically reads HCL and JSON, and generically outputs HCL and JSON,
it _does_ expect that the inputs will be maps. If you pass in a file that's not
a map at the top level I _believe_ it'll silently skip that file? Not 100% sure
(note the current test coverage).
