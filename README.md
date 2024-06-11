# LampGhost: The ghost system for BMS lamp tracker

> [!CAUTION]
> Do not directly copy LampGhost to your oraja directory, this might ruin your own songdata.db and scorelog.db file. Copy them to other directory is recommanded.

> [!CAUTION]
> This project is still in heavy development. Any change to database related code would destroy your current database.

## Main Goals

- [x] Diff yours and other players' lamp status at some point in the past
- [ ] Support lr2 score import


LampGhost aims at diff with previous status, which means you can "revert" other players' play record after some point (e.g only using the record before the rival reaches ★4)

This project so called "lampghost" is because it's like the "ghost record" in car racing game. 

## Quick start

### Build this project

```bash
$ go build
```

### Initialize

Before you doing anything, you need to call init command to setup LampGhost's database.

```bash
$ ./lampghost init
```

> [!NOTE]
> LampGhost's data is stored in home directory by default, which is `$HOME/.lampghost` on linux/osx and `%USERPROFILE%\.lampghost\` on windows.  

> [!WARNING]
> Calling `./lampghost init` the second time won't delete any data file expect lampghost.db. In other words, `init` command doesn't really **init** everything, all data files expect lampghost.db still locate at your home directory. I haven't make a decision on `init` behaviour, so in the current situation, if you do believe there is a mistake, you can delete the `.lampghost` directory direclty at your home directory.

### Import difficulty table

```bash
$ ./lampghost table add --alias insane "http://zris.work/bmstable/insane/insane_header.json"
$ ./lampghost table add --alias insane2 "http://zris.work/bmstable/insane2/insane_header.json"
```

> [!NOTE]
> Just like setting up your oraja steps, you have to import "NEW GENERATION INSANE" table to get courses data. The insane table doesn't contain it!

### Copy database file

Take me(Catizard) and my rival(Red) as an example. Make a new directory and put scorelog.db and songdata.db file into it.  
After this step done, your directory may look like:
```
.
├── Catizard
│   ├── scorelog.db
│   └── songdata.db
├── Red
│   ├── red_scorelog.db
│   └── songdata.db
└── lampghost
```

> [!NOTE]
> This is not a strictly form you must to keep, you can put them all in one directory. It's only a recommand way as an example.

### Import rivals

Register Catizard and Red into LampGhost

```bash
$ ./lampghost rival add "Catizard" ./Catizard/scorelog.db ./Catizard/songdata.db
$ ./lampghost rival add "Red" ./Red/red_scorelog.db ./Red/songdata.db
```

### Build tags(Optional)

Tags are similar to milestone. This step is to generate some tags by lampghost (e.g first clear of ★5 course)

```bash
$ ./lampghost rival build tags --rival "Red"
```

### Start race

```bash
$ ./lampghost ghost Catizard Red --tags
```

> [!NOTE]
> If you haven't generate any tags, remove '--tag' flag.

## Further

- I'm new to go language, need some help for managing packages (especially service part)
- Better frustation (need advice & help)
- Add configuration file support
- Support "shrink" mode: only one songdata.db file is used across everywhere, only scorelog.db need be loaded. One songdata.db file would be around 60MB~150MB which is too large for collecting
- Form of LampGhost is not good: you need ask your friend to give you the scorelog.db and songdata.db file. Website is a better choice. like [cinnamon](http://cinnamon.link).
- Draw charts of your (or others) progression