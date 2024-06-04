# LampGhost: The ghost system for BMS lamp tracker

> [!CAUTION]
> Do not directly copy LampGhost to your oraja directory, this might ruin your own songdata.db and scorelog.db file. Copy them to other directory is recommanded.

> [!CAUTION]
> This project is still in heavy development. Any change to database related code would destroy your current database.

## Main Goals

- Diff yours and other players' lamp status
- Diff yours and other players' previous lamp status
- Support lr2

LampGhost aims at diff with previous status. For example, you are a ★4 player while your friend is already ★8.If you ask for some lamp suggestions, you might not get too much useful infomation. LampGhost could do: only use the scorelog before your friend arrive ★5, construcing a previous status to diff.

This project so called "lampghost" is because it's like the "ghost record" in car racing game. 

## Further

- Use context package when command begins
- Better frustation (need advice & help)
- Add configuration file support
- Support "shrink" mode: only one songdata.db file is used across everywhere, only scorelog.db need be loaded. One songdata.db file would be around 60MB~150MB which is too large for collecting
- Form of LampGhost is not good: you need ask your friend to give you the scorelog.db and songdata.db file. Website is a better choice. like [cinnamon](http://cinnamon.link).
- Draw charts of your (or others) progression