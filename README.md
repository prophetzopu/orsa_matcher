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
  -platform string
        Either pcps4 or xbox (default "pcps4")
  -team1 string
        Team 1 name
  -team2 string
        Team 2 name
  -groups string
        REQUIRED IF TEAM1 and TEAM2 are empty: CSV file with previous team groups
  -ranks string
        REQUIRED IF TEAM1 and TEAM2 are empty: CSV file with team ranks
  -week int
        The number of the week (default 1)
```

### Manual Matchmaking

Simply run the `cli`, specifying the team names, platform, format, and week:
```
./cli --week 3 --format 2v2 --platform pcps4 -team1 "Team Rocket" -team2 "No Boost No Problem"
```

This will output both `Team Rocket vs No Boost No Problems` and `No Boost No Problems vs Team Rocket` in Reddit and ORSA website markdown:

```
<----------REDDIT FORMAT---------->

**Groups for Week 3 of the 2v2 PCPS4 league are as follows:**

Group 1: Team Rocket, No Boost No Problem

[Team Rocket vs No Boost No Problem](https://muut.com/arlchampionships#!/week-3-2v2-pcps4:team-rocket-vs-no-boost-no)

[No Boost No Problem vs Team Rocket](https://muut.com/arlchampionships#!/week-3-2v2-pcps4:no-boost-no-problem-vs-team)

<----------ORSA FORMAT---------->

<button class=”orsa-accordion”>Team Rocket vs No Boost No Problem</button>
<div class=”orsa-muut-thread”>
<a class=”muut” href=”https://muut.com/arlchampionships#!/week-3-2v2-pcps4:team-rocket-vs-no-boost-no/comments”></a>
<script src=”//cdn.muut.com/1/moot.min.js”></script>
</div>
<button class=”orsa-accordion”>No Boost No Problem vs Team Rocket</button>
<div class=”orsa-muut-thread”>
<a class=”muut” href=”https://muut.com/arlchampionships#!/week-3-2v2-pcps4:no-boost-no-problem-vs-team/comments”></a>
<script src=”//cdn.muut.com/1/moot.min.js”></script>
```

### Automatic Matchmaking
Download your spreadsheets as individual CSV files (File -> Download As in google docs), and save them somewhere convenient.

Once you have a file for the previous groups and the current ranks, just run
```
C:\project>cli -week 3 -format 3v3 -platform pcps4 -ranks \path\to\ranks.csv -groups \path\to\groups.csv
```

Which will output a list of matches in Reddit markdown, as well as a simple CSV of all the matches created.
