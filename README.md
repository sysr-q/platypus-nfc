# PlatypusCON^H^H^HCAMP NFC workshop

If you're just interested in reading the data off of the passport tags, the keys are:

* `DEADBEEF1337` (key A, for all sectors, has read/write access)
* `FFEEDDCCBBAA` (key B, for all sectors except 1)
* `D3F7D3F7D3F7` (key B for sector 1)

The tag contains an NDEF format with a `text/plain` record containing the flag:
`$LJHS$1ca8c3ffd5d3617ce45a6152f77ea068$`.

A full dump of the master tag the rest were cloned from is checked into the repo
as `platypuscamp_master_nfc.txt`.

## `display/`

This runs as a CLI on the RasPI or whatever. Basically listens on `:1337` for
either `0` (clear) `1` (success/green) `2` (failure/red) `$` (dolla dolla bills).
Blanks the whole terminal that color. Very boring, but whipped together at the
last second since I realised the LEDs on the ACR122U don't work with libnfc
compiled on ARMvh7.

## `listener/`

Daemon that basically listens for libnfc events, grabs data off the card, then
depending on whether or not the authentication is successful (determined based
on current config) connects to the display service and flashes the screen.

## `printout/`

This is the things you saw printed at the workshop table. Includes the source
Markdown and the compiled PDF version.

## Notes, ideas, suggestions, etc.

The current bindings for libfreefare for Go don't work against libfreefare 0.4.0.
To remedy this, I figured out the solution and included a patch, so just apply
`go_freefare_tag.diff` over the top of your freefare binding folder in your GOPATH.
