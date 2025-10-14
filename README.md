# ember
Game engine for building HTML5 games

## Installation instructions 

You'll need to install the Go programming language for this, either use a package manager to install or use the official website to download

### Clone this repo 
```bash 
git clone https://github.com/lordryns/ember.git
```

### Install dependencies 
```bash 
go mod tidy
```

This command should install everything you need to run this app (mostly the Fyne toolkit)

### Build from source

There are two ways to do this (maybe three if you just want to run once)

**Run the script**

```bash 
go run .
```

This will run the app without creating a binary 

**Build normally into a binary** 
```bash 
go build
```

This will give you a ./ember or ember.exe depending on your operating system 

**Turn into a Distributable binary**
You don't need this option if you're just trying it out

```bash 
go install fyne.io/fyne/v2/cmd/fyne@latest
```

```bash 
fyne install 
```

if this stage does not work, refer to the fyne documentation 
```
