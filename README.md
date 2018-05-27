# Apple //e emulator in go

Ebiten based apple //e emulator written in go.

## Features

* MOS 6502 CPU
* keyboard
* 40 column text mode
* low resolution color graphics
* high resolution monochrome graphics
* upper memory bank switching: $d000 page and ROM/RAM switching
* main memory page1/page2 switching in text, lores and hires
* speaker audio

## Installation

Install prerequisites with glide

    glide up

Build the executable

    go build

Download `apple2e.rom` from
[a2go.applearchives.com](http://a2go.applearchives.com/roms/) and put it in the root directory.

## Running it

    ./apple2
    ./apple2 my_disk_image.dsk
    ./apple2 -drive-head-click my_disk_image.dsk

## Running the tests
### Setting up the tests

Some disk images are used for the tests

* dos33.dsk from e.g. [mirrors.apple2.org.za](https://mirrors.apple2.org.za/ftp.apple.asimov.net/images/masters/DOS33_blank_with_integer_basic.DSK)
* prodos19.dsk from e.g. [mirrors.apple2.org.za](https://mirrors.apple2.org.za/ftp.apple.asimov.net/images/masters/prodos/ProDOS_1_9.dsk)

### Running the tests

    go test -v

The CPU tests make use of [Klaus2m5's](https://github.com/Klaus2m5/6502_65C02_functional_tests)
 excellent 6502 functional tests.

## Known working images
* DOS 3.3
* Prodos 1.9
* Lemonade stand
* Montezuma's Revenge

## Keyboard shortcuts

* ctrl-alt-R reset
* ctrl-alt-M mute
* ctrl-alt-C capslock

## Remaining work

* 80 column card
* 48k aux memory
* double hires
* paddles

## Known issues

1. On MacOS, the initial beep is sometimes split into two little beeps. This appears to be an ebiten issue.
2. On MacOS, shutting down sometimes takes 30 seconds or so
