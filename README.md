# LampGhost

English | [简体中文](./README.zh-CN.md)

## What is LampGhost?

`LampGhost` is an **offline** and **cross-platform** `beatoraja` save file viewer, or say "lamp viewer". It supports:

- View your save file, show your lamp status and course result
- Add your friends' save file, show your friends' lamp status
- Compare lamp status between you and your friends, and could even "turn back" your friend's save file time(Suppose your friend is already ^^, using `LampGhost` allows you comparing the lamp status the moment that your friend achived ★5)
- Customize favorite folder and import it as a difficult table in `beatoraja`(tip: You can reload table by pressing F2 in game). More specifically, import `http://127.0.0.1:7391/table/lampghost.json` in `beatoraja` while keep `LampGhost` opening.

Because `LampGhost` is an **offline** tool, therefore it:

- No need of an `ir.jar` file for connecting a centerialized server. Therefore won't be unusable due to server is down
- If you want to compare your friends' save file, you need to ask them and handle them manually
- There might be a risk for some breaking changes that you have to clear all `LampGhost` related contents

> [!warning]
>
> Data is invaluable, backup your save file is highly recommanded. Although `LampGhost` would only **read** your save file but never **modify** it, but there still be some rare bug that ruin your save file

## How to download and use LampGhost?

Grab an executable at `Release` page, then just open it.

## Attentions

This project is still at early demo phase, you might encounter some mystery problem(Don't worry, won't ruin your save file most likely). You could delete the whole `LampGhost` data folder to have a clean restart:

For `Windows` user it's located at `%USERPROFILE%\.lampghost_wails`, for `Linux/OSX` user is `$HOME/.lampghost_wails`

### Unfortunate LR2 users

This project might never implements the support of LR2 save file, although it's predecessor has. This is because LR2 won't record play time in its database while most features of `LampGhost` are based on it.

## Bugs reports and Advices

This project is made by my own(at least for now), so it's still very simple and biased. Any bug reports and advices are very welcomed!
