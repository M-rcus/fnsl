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
