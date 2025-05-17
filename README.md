# LampGhost

English | [简体中文](./README.zh-CN.md)

## What is LampGhost?

![showcase](./doc/showcase.png)

`LampGhost` is an **offline** and **cross-platform** `beatoraja` save file viewer, favorite folder manager and provides lamp diff view . It supports:

- `I18n`:`English` and `zh-CN` are supported currently
- `Multiple User`: You can load your friends' save file and view their scores
- `Time Machine`: Compare lamp status between you and your friends, and could even "turn back" your friend's save file time(Suppose your friend is already ^^, using `LampGhost` allows you comparing the lamp status the moment that your friend achived ★5)
- `Favorite Folder`: Customize favorite folder and import it as a difficult table in `beatoraja`(tip: You can reload table by pressing F2 in game). More specifically, import `http://127.0.0.1:7391/table/lampghost.json` in `beatoraja` while keep `LampGhost` opening.

Because `LampGhost` is an **offline** tool, therefore:

- No need of an `ir.jar` file for connecting a centerialized server. Therefore won't be unusable due to server is down
- If you want to compare your friends' save file, you need to ask them and handle them manually, which could be a tedious task
- It's always here, when you want to see your progress, the only thing you need to do is pressing the sync button
- There might be a risk for some breaking changes that you have to remove `LampGhost`'s database file and have a clean start

> [!warning]
>
> Data is invaluable, backup your save file is highly recommanded. Although `LampGhost` would only **read** your save file but never **modify** it, but there still be some rare bug that ruin your save file

## How to download and use LampGhost?

Grab an executable at `Release` page, then double-click it. This should be working for `Windows` and `Linux`.

## Build from scratch

- Install [wails](https://github.com/wailsapp/wails)
- `wails build`

## Attentions

First of all, you need to import some difficult tables just like `beatoraja` to make the most features work.

Secondly,this project is still at early demo phase, you might encounter some mystery problem(Don't worry, won't ruin your save file most likely). You could delete the whole `LampGhost` data folder to have a clean restart:

For `Windows` user it's located at `%USERPROFILE%\.lampghost_wails`, for `Linux/OSX` user is `$HOME/.lampghost_wails`

### Unfortunate LR2 users

This project doesn't support load LR2 save file currently, although its precessor, the tui version of `LampGhost` does.

## Bugs reports and Advices

This project is made by my own(at least for now), so it's still very simple and biased. Any bug reports and advices are very welcomed!

## Thanks

- [@Yuntian](https://www.github.com/Yuntian52s) Providing ui design advices, without his selfless work `LampGhost` would still a be a very unfriendly tool
- [@Chuang](https://github.com/chuang1213), yzy, yf and other early version testers
- [Wails](https://github.com/wailsapp/wails): This project is based on `wails` framework, and benefits a lot from its amazing hot reload time and easy mind model
