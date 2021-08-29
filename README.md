# RSS notification
Informs you of unread updates to selected RSS feeds.

## Usage
1. Clone repository.
```shell
git clone https://github.com/qasterr/rss-notification
cd rss-notification
```
2. Create a `log.txt` file, leave it empty.
3. Create a `list.txt` file and write an url to an RSS feed on each line.
4. Install all dependencies.
```shell
go get -u .
```
5. Build the executable.
```shell
go build
```
6. OPTIONAL, add the executable to system startup (OS-dependant).