# areena-go

This is a tool that I made *very quickly* to download the 2019 speed skating world cup finals from Yle Areena, because [yle-dl](https://aajanki.github.io/yle-dl/) couldn't download them at the time. It's very likely that this tool won't work anymore, so it's is only here for reference. I recommend using `yle-dl` for regular Areena downloading.

## Usage

`areena-go 1-50083241`

`1-50083241` is the video ID that you can find in the URL.

Please note that this tool only downloads the raw `.ts` files, so you need to manually combine them into a single video file. `ffmpeg` is a good tool for this. I'll add step by step instructions the next time I need to update this tool.
