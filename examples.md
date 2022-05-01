# Local
Vaultctl support adding local directories and files to the vault. A `config` file will be generated in `~/.vault` to keep track of added directories and single files. To keep `~/.vault/config` as light as possible not all files in a directory are added. Rather in the directory itself there is another `.location` file which serves all images of this directory.

Example commands for the local subcommand:
```bash
vaultpaper local addDir path/to/folder
vaultpaper local deleteDir path/to/folder
vaultpaper local addFile path/to/picture.jpg
vaultpaper local deleteFile path/to/picture.jpg
vaultpaper local random
```

# Reddit
The `reddit` subcommand downloads pictures from Reddit and stores them inside their respective subreddit directory in the `~/.vault`, which contains a `.location` file for indexing.
```
~/.vault
  - config
  - itookapicture
    - .location
    - picture1.jpg
    - picture2.jpg
  - wallpaper
    - .location
    - picture1.jpg

```

Example commands for the `reddit` subcommand:
```bash
vaultpaper reddit download itookapicture 
vaultpaper reddit download itookapicture --count 5 --sort top --time today
vaultpaper reddit delete path/to/picture
vaultpaper reddit set itookapicture --sort top --time today
vaultpaper 
vaultpaper local random
```