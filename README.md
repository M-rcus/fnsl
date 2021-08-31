# Fansly downloader

> Version 1.0 - 2021-08-11

This downloads all photos and videos from Fansly users you're subscribed to. Note: you must have a subscription to a user (usually that means paying for it), this only saves the photos/videos you have access to.

You'll need Go 1.16 or higher installed.

## Configuration

Log in to Fansly with a web browser to get a session.

Rename `config.yaml.sample` to `config.yaml` then set:

- `authorization`: the value of the authorization token. To get it, use the browser's console (developer tools) and run this:  
  
  ```js
  JSON.parse(window.localStorage.getItem("session_active_session")).token
  ```

  Result will be similar to "Mjc5N[...]Y"
- `userAgent`: your browser's user agent. From the console:  
  
  ```js
  window.navigator.userAgent
  ```
  
  Result will be similar to "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.1.2 Safari/605.1.15"

Optionally set `outDir` to the directory where you want files to be saved (default is `out` in the current folder).

## Start scraper

In the folder where the source code is, open a terminal and run:

```sh
go run .
```

## Credit & license

All credit for version 1.0 goes to [reddit user /u/Winter-Elephant-2250](https://www.reddit.com/user/Winter-Elephant-2250/comments/p3j87m/released_fansly_scraper_app_open_source/) - Initial commit: dc35e55e21a77d423855917b062d975195577e2a

Based on the reddit post, version 1.0 is licensed under the "WTF Public License" (WTFPL)

```
      DO WHAT THE FUCK YOU WANT TO PUBLIC LICENSE 
                Version 2, December 2004

Everyone is permitted to copy and distribute verbatim or modified 
copies of this license document, and changing it is allowed as long 
as the name is changed. 

          DO WHAT THE FUCK YOU WANT TO PUBLIC LICENSE 
  TERMS AND CONDITIONS FOR COPYING, DISTRIBUTION AND MODIFICATION 

1. You just DO WHAT THE FUCK YOU WANT TO.
```

However, any changes afterwards will be licensed under "Unlicense": https://unlicense.org/  
Mainly as a "covering my ass" policy :)

```
This is free and unencumbered software released into the public domain.

Anyone is free to copy, modify, publish, use, compile, sell, or
distribute this software, either in source code form or as a compiled
binary, for any purpose, commercial or non-commercial, and by any
means.

In jurisdictions that recognize copyright laws, the author or authors
of this software dedicate any and all copyright interest in the
software to the public domain. We make this dedication for the benefit
of the public at large and to the detriment of our heirs and
successors. We intend this dedication to be an overt act of
relinquishment in perpetuity of all present and future rights to this
software under copyright law.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
IN NO EVENT SHALL THE AUTHORS BE LIABLE FOR ANY CLAIM, DAMAGES OR
OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE,
ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
OTHER DEALINGS IN THE SOFTWARE.

For more information, please refer to <http://unlicense.org/>
```