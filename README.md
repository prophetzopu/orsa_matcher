# orsa_matcher
Simple tool for matching teams in the ORSA competitive pool

## Installing

### Windows
[Download the go binaries for windows](https://golang.org/dl/) and run the installer. Default installation is at `C:\Go`.

Add `C:\Go\bin` to your system `%PATH`. Open the System control panel -> Advanced system settings -> Environment Variables.

From powershell or cmd prompt, create a directory from your project, and set your GOPATH (you’ll have to repeat this for each cmd prompt)
```
C:\Users\me>mkdir path\to\my\project
C:\Users\me>cd path\to\my\project
C:\Users\me\path\to\my\project>set GOPATH=C:\Users\me\path\to\my\project
```

Pull down this project
```
C:\project>go get github.com/adamschaub/orsa_matcher
```

This is a lazy project, so you have to install dependencies yourself :(
```
C:\project>go get -u github.com/gocarina/gocsv
```

Then build the project
```
C:\project>go build github.com/adamschaub/orsa_matcher
C:\project>go build github.com/adamschaub/orsa_matcher/cli
```

## Running

The CLI has a number of options
```
-format string
    One of 3v3/2v2/1v1 (default "3v3")
-groups string
    REQUIRED: CSV file with previous team groups
-platform string
    Either pcps4 or xbox (default "pcps4")
-ranks string
    REQUIRED: CSV file with team ranks
-week int
    The number of the week (default 1)
```

Download your spreadsheets as individual CSV files (File -> Download As in google docs), and save them somewhere convenient.

Once you have a file for the previous groups and the current ranks, just run
```
C:\project>cli -week 3 -format 3v3 -platform pcps4 -ranks \path\to\ranks.csv -groups \path\to\ranks.csv
```

Which will output a list of matches in Reddit markdown, as well as a simple CSV of all the matches created.
