<p align="center">
  <img src="https://raw.githubusercontent.com/tiagorlampert/CHAOS/master/content/logo.png">
</p>

<h1 align="center">CHAOS Payload Generator</h1>
<p align="center">
  <a href="https://golang.org/">
    <img src="https://img.shields.io/badge/Golang-1.11-blue.svg">
  </a>
  <a href="https://github.com/tiagorlampert/CHAOS/blob/master/LICENSE">
    <img src="https://img.shields.io/badge/License-BSD%203-lightgrey.svg">
  </a>
  <a href="https://github.com/tiagorlampert/CHAOS/blob/master/main.go">
    <img src="https://img.shields.io/badge/Release-3.0-red.svg">
  </a>
    <a href="https://opensource.org">
    <img src="https://img.shields.io/badge/Open%20Source-%E2%9D%A4-brightgreen.svg">
  </a>
</p>

<p align="center">
  CHAOS is a PoC that allow generate payloads and control remote operating systems.
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
| `Keylogger`              |    X    |        |       |
| `Persistence`            |    X    |        |       |
| `Open URL`               |    X    |    X   |   X   |
| `Get OS Info`            |    X    |    X   |   X   |
| `Fork Bomb`              |    X    |    X   |   X   |
| `Run Hidden`             |    X    |        |       |

## Tested On
[![Kali)](https://www.google.com/s2/favicons?domain=https://www.kali.org/)](https://www.kali.org) **Kali Linux - ROLLING EDITION**

## How to Install
```bash
# Install dependencies
$ sudo apt install golang git go-dep -y

# Get this repository
$ go get github.com/tiagorlampert/CHAOS

# Go into the repository
$ cd ~/go/src/github.com/tiagorlampert/CHAOS

# Get project dependencies
$ dep ensure

# Run
$ go run main.go
```

## How to Use

Command     | On HOST does...
:-----      |:-----
`generate`  |Generate a payload (e.g. `generate lhost=192.168.0.100 lport=8080 fname=chaos --windows`)
`lhost=`    |Specify a ip for connection
`lport=`    |Specify a port for connection
`fname=`    |Specify a filename to output
`--windows` |Target Windows
`--macos`   |Target Mac OS
`--linux`   |Target Linux
`listen`    |Listen for a new connection (e.g. `listen lport=8080`)
`serve`     |Serve files
`exit`      |Quit this program

Command                 | On TARGET does...
:-----                  |:-----
`download`              |File Download
`upload`                |File Upload
`screenshot`            |Take a Screenshot
`keylogger_start`       |Start Keylogger session
`keylogger_show`        |Show Keylogger session logs
`persistence_enable`    |Install at Startup
`persistence_disable`   |Remove from Startup
`getos`                 |Get OS name
`lockscreen`            |Lock the OS screen
`openurl`               |Open the URL informed
`bomb`                  |Run Fork Bomb
`clear`                 |Clear the Screen
`back`                  |Close connection but keep running on target
`exit`                  |Close connection and exit on target

## Video
<p align="center">
<a href="https://www.youtube.com/watch?v=Fq_0yDPFjYE">
  <img src="https://img.youtube.com/vi/Fq_0yDPFjYE/maxresdefault.jpg" width="700"/>
</a></p>

## Gif
<p align="center">
<img src="https://github.com/tiagorlampert/CHAOS/blob/master/content/screenshot.gif">
</p>

## FAQ
> #### Why does Keylogger capture all uppercase letters?
> All the letters obtained using the keylogger are uppercase letters. It is a known issue, in case anyone knows how to fix the Keylogger function using golang, please contact me or open an issue.

> #### Why are necessary get and install external libraries?
> To implement the screenshot function i used a third-party library, you can check it in https://github.com/kbinani/screenshot and https://github.com/lxn/win. You must download and install it to generate the payload.

## Contact
:email: **tiagorlampert@gmail.com**

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
