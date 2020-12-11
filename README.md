<p align="center">
  <img src="https://raw.githubusercontent.com/tiagorlampert/CHAOS/master/content/logo.png">
</p>

<p align="center">
  <a href="https://golang.org/">
    <img src="https://img.shields.io/badge/Golang-1.15-blue.svg">
  </a>
  <a href="https://github.com/tiagorlampert/CHAOS/blob/master/LICENSE">
    <img src="https://img.shields.io/badge/License-BSD%203-lightgrey.svg">
  </a>
  <a href="https://github.com/tiagorlampert/CHAOS/blob/master/main.go">
    <img src="https://img.shields.io/badge/Release-4.X-red.svg">
  </a>
    <a href="https://opensource.org">
    <img src="https://img.shields.io/badge/Open%20Source-%E2%9D%A4-brightgreen.svg">
  </a>
</p>

<p align="center">
  CHAOS is a Remote Administration Tool that allow generate binaries to control remote operating systems.
</p>

## Disclaimer
<p align="center">
  :books: This project was created only for learning purpose.
</p>

THIS SOFTWARE IS PROVIDED "AS IS" WITHOUT WARRANTY OF ANY KIND. YOU MAY USE THIS SOFTWARE AT YOUR OWN RISK. THE USE IS COMPLETE RESPONSIBILITY OF THE END-USER. THE DEVELOPERS ASSUME NO LIABILITY AND ARE NOT RESPONSIBLE FOR ANY MISUSE OR DAMAGE CAUSED BY THIS PROGRAM.

## Features

| Feature                  |  ![w]   |  ![m]  |  ![l] |
|:-------------------------|:-------:|:------:|:-----:|
| `Reverse Shell`          |    X    |    X   |   X   |
| `Download File`          |    X    |    X   |   X   |
| `Upload File`            |    X    |    X   |   X   |
| `Screenshot`             |    X    |    X   |   X   |
| `Persistence`            |    X    |        |       |
| `Open URL`               |    X    |    X   |   X   |
| `Get OS Info`            |    X    |    X   |   X   |
| `Run Hidden`             |    X    |        |       |

## How to Install
```bash
# Install dependencies
$ sudo apt install golang git -y

# Get this repository
$ git clone https://github.com/tiagorlampert/CHAOS

# Go into the repository
$ cd CHAOS/

# Run
$ go run cmd/chaos/main.go
```

## How to Use

Command     | On HOST does...
:-----      |:-----
`generate`  |Generate a binary (e.g. `generate address=192.168.0.100 port=8080 filename=chaos --windows --hidden`)
`address=`  |Specify a ip for connection
`port=`     |Specify a port for connection
`filename=` |Specify a filename to output binary
`--windows` |Target Windows
`--macos`   |Target Mac OS
`--linux`   |Target Linux
`--hidden`  |Run a hidden binary (only for Windows)
`listen`    |Listen for a new connection (e.g. `listen port=8080`)
`serve`     |Serve directory files
`exit`      |Quit this program

Command                 | On TARGET does...
:-----                  |:-----
`download {filePath}`   |File Download
`upload {filePath}`     |File Upload
`screenshot`            |Take a Screenshot
`persistence enable`    |Install at Startup
`persistence disable`   |Remove from Startup
`information`           |Get OS information
`lockscreen`            |Lock the OS screen
`open-url {url}`              |Open the URL informed
`exit`                  |Quit app

## Video (out of date)
<p align="center">
<a href="https://www.youtube.com/watch?v=Fq_0yDPFjYE">
  <img src="https://img.youtube.com/vi/Fq_0yDPFjYE/maxresdefault.jpg" width="700"/>
</a></p>

## Gif (out of date)
<p align="center">
<img src="https://github.com/tiagorlampert/CHAOS/blob/master/content/screenshot.gif">
</p>

## Donate
If you enjoyed this project, give me a cup of coffee. :)

[![Donate](https://www.paypalobjects.com/en_US/i/btn/btn_donateCC_LG.gif)](https://www.paypal.com/cgi-bin/webscr?cmd=_donations&business=SG83FSKPKCRJ6&currency_code=USD&source=url)


## License

>The [BSD 3-Clause License](https://opensource.org/licenses/BSD-3-Clause)
>
>Copyright (c) 2017, Tiago Rodrigo Lampert
>
>All rights reserved.
>
>Redistribution and use in source and binary forms, with or without modification, are permitted provided that the following conditions are met:
>
>* Redistributions of source code must retain the above copyright notice, this
  list of conditions and the following disclaimer.
>
>* Redistributions in binary form must reproduce the above copyright notice,
  this list of conditions and the following disclaimer in the documentation
and/or other materials provided with the distribution.
>
>* Neither the name of the copyright holder nor the names of its
  contributors may be used to endorse or promote products derived from
this software without specific prior written permission.
>
>THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

[w]:https://raw.githubusercontent.com/tiagorlampert/CHAOS/master/content/windows.png "Windows status"
[l]:https://raw.githubusercontent.com/tiagorlampert/CHAOS/master/content/linux.png "Linux status"
[m]:https://raw.githubusercontent.com/tiagorlampert/CHAOS/master/content/mac.png "Mac OS status"
