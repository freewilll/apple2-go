# Apple // emulator in go

An Apple //e emulator written in Go using [ebiten](https://github.com/hajimehoshi/ebiten).

## Features

* MOS 6502 CPU
* Keyboard
* 40 column text mode
* Low resolution color graphics
* High resolution monochrome graphics
* Upper memory bank switching: $d000 page and ROM/RAM
* Main memory page1/page2 switching in text, lores and hires
* Disk image reading & writing
* Speaker audio

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

## Keyboard shortcuts

* ctrl-alt-R reset
* ctrl-alt-M mute
* ctrl-alt-C caps lock

## Running the tests
### Setup

The tests use DOS and Prodos disk images. Download them from

* dos33.dsk from [mirrors.apple2.org.za](https://mirrors.apple2.org.za/ftp.apple.asimov.net/images/masters/DOS33_blank_with_integer_basic.DSK)
* prodos19.dsk from [mirrors.apple2.org.za](https://mirrors.apple2.org.za/ftp.apple.asimov.net/images/masters/prodos/ProDOS_1_9.dsk)

### Running the tests

    go test -v

The CPU tests make use of [Klaus2m5's](https://github.com/Klaus2m5/6502_65C02_functional_tests)
 excellent 6502 functional tests.

## Known working disk images
* DOS 3.3
* Prodos 1.9
* Lemonade stand
* Montezuma's Revenge

## Remaining work

* 80 column card
* 48k aux memory
* double hires
* joystick