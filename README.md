# Apple // emulator in go

An Apple //e emulator written in Go using [ebiten](https://github.com/hajimehoshi/ebiten).

## Features

* MOS 6502 CPU
* Keyboard
* 40 column text mode
* Low resolution monochrome and color graphics
* High resolution monochrome and color graphics
* Upper memory bank switching: $d000 page and ROM/RAM
* Main memory page1/page2 switching in text, lores and hires
* Disk image reading & writing
* Speaker audio

## Installation

The installation requires go modules go be installed

Build the executable

    go build

Download `apple2e.rom` from
[a2go.applearchives.com](http://a2go.applearchives.com/roms/) and put it in the root directory.

## Running it

    ./apple2-go
    ./apple2-go my_disk_image.dsk
    ./apple2-go -drive-head-click my_disk_image.dsk

## Keyboard shortcuts

* ctrl-alt-R reset
* ctrl-alt-M toggle monochrome/color display
* ctrl-alt-C caps lock
* ctrl-alt-F show FPS

## Running the tests
### Setup

The tests use DOS and Prodos disk images. Download them from

* dos33.dsk from [mirrors.apple2.org.za](https://mirrors.apple2.org.za/ftp.apple.asimov.net/images/masters/DOS33_blank_with_integer_basic.DSK)
* prodos19.dsk from [mirrors.apple2.org.za](https://mirrors.apple2.org.za/ftp.apple.asimov.net/images/masters/prodos/ProDOS_1_9.dsk)

### Running the tests

    go test -v

The CPU tests make use of [Klaus2m5's](https://github.com/Klaus2m5/6502_65C02_functional_tests)
 excellent 6502 functional tests.

### Creating the CPU test ROMs

The source files are `6502_functional_test.a65` and `6502_interrupt_test.a65`. They are assembled using `as65` into a binary file which contains a memory image of the test code. They are compressed into gzip files which are loaded into the apple memory by the unit tests.

Download [as65](http://www.kingswood-consulting.co.uk/assemblers/as65_142.zip) and unzip it to get the `as65` assembler binary.

Assemble the tests

    cd cpu
    as65 -l -m -w -h0 6502_functional_test.a65
    gzip 6502_functional_test.bin

    as65 -l -m -w -h0 6502_interrupt_test.a65
    gzip 6502_interrupt_test.bin


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

## Coding standards

Use `gofmt` to ensure standard go style consistency

     go fmt $(go list ./... | grep -v /vendor/)

Use `golint` to ensure Google's style consistency

    golint $(go list ./...)
