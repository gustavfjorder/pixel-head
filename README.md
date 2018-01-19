# Zombie Hunter 3000!
A realtime multiplayer zombie survival game created in Go.
Go kill zombies with 4 different weapons, knife, handgun, rifle and shotgun.
But be aware, there are four hardcore types of zombies, the slowpoke, the speedy gonzales, the exploder and the regular booring one.
Remember to place exploding barrels and pickup upgrades in the crates.

## Features
- Singleplayer
- Multiplayer on LAN
- Multiplayer on server (no server hosted at the current time)
- 4 different zombies
- 4 different weapons
- Exploding barrels
- Loot boxes for upgrades

## Controls
- `W` - Move up
- `A` - Move left
- `S` - Move down
- `D` - Move right
- `B` - Place an exploding barrel
- `Space` - Shoot/melee
- `1-4` - Change weapon

## Download
This game can be played on every OS.

- [Windows](https://drive.google.com/open?id=1VfA5MoDqfPXdWo6tLRg4FdFuc9yHy7eX)
- [Linux](https://drive.google.com/open?id=1PYUj8ldCjBPhFMDrC3-bLqq4Bm3ETX4r)
- [macOS](https://drive.google.com/open?id=1jj3dTT8CxyI7UlNsuCKlQp8-0_GAmLj5)

## Development
This project relies on [dep](https://github.com/golang/dep) for depedency management.
Please make sure it is installed before proceeding.

To start developing on this project, clone the repository and install the dependencies.

```bash
git clone https://github.com/gustavfjorder/pixel-head.git
dep ensure
```

To run the game engine, please see [pixel](https://github.com/faiface/pixel#requirements) for requirements.
If you are on Windows please see [this page](https://github.com/faiface/pixel/wiki/Building-Pixel-on-Windows) too.

## Known issues

### macOS is only single player
Probably due to issues with the tuple space `goSpace` macOS cannot play in a multiplayer game.
Other clients cannot connect to the macOS space either.

Giving follow error message:

```
ErrReceiveMessage: read tcp4 172.20.10.5:51575->172.20.10.14:31416: read: connection reset by peer
panic: 
  github.com/gustavfjorder/pixel-head/vendor/github.com/pspaces/gospace/space:
    Space(tcp://172.20.10.14:31416/lounge).Put("request", "b9gsu95nc8e4vq10chi0"): operation on this space failed.
```

### macOS livelocking
Playing the game with too many server requests will lock the tuple space for a some time.
It automatically resolves the deadlock and continues for 10-20 seconds.
