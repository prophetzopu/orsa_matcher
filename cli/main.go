package main

import (
    "fmt"
    "flag"
    "os"
    "math"
    "bufio"
    "strings"
    "container/heap"
    orsa "github.com/adamschaub/orsa_matcher"
    "github.com/gocarina/gocsv"
)

const StdGroupSize int = 3
const MaxGroupSize int = 4
const MinGroupSize int = 2

type GroupMap map[string]map[string]bool

func readGroups(fileName *string) GroupMap {
    groupFile, err := os.Open(*fileName)

    if err != nil {
        panic(err)
    }
    defer groupFile.Close()

    scanner := bufio.NewScanner(bufio.NewReader(groupFile))
    groups := make(GroupMap)
    for scanner.Scan() {
        teams := strings.Split(scanner.Text(), ",")
        for _, thisTeam := range teams {
            if thisTeam == "" {
                continue
            }

            existingTeam, ok := groups[thisTeam]
            if !ok {
                existingTeam = make(map[string]bool)
                groups[thisTeam] = existingTeam
            }

            for _, otherTeam := range teams {
                if otherTeam == "" || thisTeam == otherTeam {
                    continue
                }
                existingTeam[otherTeam] = true
            }
        }
    }
    return groups
}

var rankFileName = flag.String("ranks", "", "REQUIRED IF TEAM1 and TEAM2 are empty: CSV file with team ranks")
var groupsFileName = flag.String("groups", "", "REQUIRED IF TEAM1 and TEAM2 are empty: CSV file with previous team groups")
var platform = flag.String("platform", "pcps4", "Either pcps4 or xbox")
var format = flag.String("format", "3v3", "One of 3v3/2v2/1v1")
var week = flag.Int("week", 1, "The number of the week")
var team1 = flag.String("team1", "", "Team 1 name")
var team2 = flag.String("team2", "", "Team 2 name")

func checkFlags() bool {
    return (*team1 != "" && *team2 != "") || (*rankFileName != "" && *groupsFileName != "")
}

func PlaceExtraTeam(matches *[]orsa.MatchGroup, ranks *orsa.RankedList, groups GroupMap) {
    team := heap.Pop(ranks).(orsa.TeamRank)
    bestFit := 0
    bestFitConflicts := 0
    for i := len(*matches)-1; i >= 0; i-- {
        match := (*matches)[i]
        conflicts := 0
        for _, otherTeam := range match.Teams {
            if _, ok := groups[team.Name][otherTeam.Name]; ok {
                conflicts++
            }
        }
        if conflicts < match.N - 1 {
            if math.Abs((*matches)[bestFit].AvgPoints - float64(team.Points)) > math.Abs(match.AvgPoints - float64(team.Points)) {
                bestFit = i
                bestFitConflicts = conflicts
            }
        }
    }
    if bestFitConflicts == 0 {
        (*matches)[bestFit].Teams[(*matches)[bestFit].N] = team
    } else {
        insertInd := 0
        for i, currTeam := range (*matches)[bestFit].Teams {
            nextTeamInd := (i+1) % (*matches)[bestFit].N
            nextTeam := (*matches)[bestFit].Teams[nextTeamInd]
            _, okTeam1 := groups[team.Name][currTeam.Name]
            _, okTeam2 := groups[team.Name][nextTeam.Name]
            if !okTeam1 && !okTeam2 {
                insertInd = nextTeamInd
                break
            }
        }
        for i := len((*matches)[bestFit].Teams)-1; i >= 0; i-- {
            if i == insertInd {
                (*matches)[bestFit].Teams[i] = team
                break
            } else {
                (*matches)[bestFit].Teams[i] = (*matches)[bestFit].Teams[i-1]
            }
        }
    }
    (*matches)[bestFit].N++
}

func DisbandMatch (match orsa.MatchGroup, pool *orsa.RankedList) {
    for _, t := range match.Teams {
        if t.Name == "" { continue }
        heap.Push(pool, t)
    }
}

func MakeMatches(matches []orsa.MatchGroup, interleave []orsa.TeamRank, ranks *orsa.RankedList, groups GroupMap) []orsa.MatchGroup {
    deferredTeams := []orsa.TeamRank{}
    for ranks.Len() >= MinGroupSize && ranks.Len() > 1 {


        var firstTeam orsa.TeamRank
        if len(interleave) > 0 {
            firstTeam = interleave[0]
            interleave = interleave[1:len(interleave)]
        } else {
            firstTeam = heap.Pop(ranks).(orsa.TeamRank)
        }
        currentMatch := &orsa.MatchGroup{Teams: [4]orsa.TeamRank{firstTeam,}, N: 1, AvgPoints: float64(firstTeam.Points)}

        validTeams := true
        for ( ranks.Len() > 0 || len(interleave) > 0 ) && currentMatch.N < StdGroupSize && validTeams {
            var team orsa.TeamRank
            if ranks.Len() > 0 {
                team = heap.Pop(ranks).(orsa.TeamRank)
            } else {
                team = interleave[0]
                interleave = interleave[1:len(interleave)]
            }

            addTeam := true
            for _, otherTeam := range currentMatch.Teams {
                if otherTeam.Name == "" { continue }
                if _, ok := groups[team.Name][otherTeam.Name]; ok {
                    //These teams have already played. Defer the current team to a later match
                    addTeam = false
                    break
                }
            }
            if addTeam {
                currentMatch.Teams[currentMatch.N] = team
                currentMatch.AvgPoints += float64(team.Points)
                currentMatch.N++
            } else if ranks.Len() == 0 && len(interleave) == 0 {
                deferredTeams = append(deferredTeams, team)
                validTeams = false
            } else {
                deferredTeams = append(deferredTeams, team)
            }
        }

        //add deferred teams back into the pool
        for _,team := range deferredTeams {
            heap.Push(ranks, team)
        }
        deferredTeams = []orsa.TeamRank{}

        if validTeams {
            currentMatch.AvgPoints /= float64(currentMatch.N)
            matches = append(matches, *currentMatch)
        } else {
            for _, team := range currentMatch.Teams {
                if team.Name != "" { heap.Push(ranks, team) }
            }
            break
        }

        if ranks.Len() < MinGroupSize {
            //dump the interleave if rank pool gets too low, and try again
            for _,t := range interleave {
                heap.Push(ranks, t)
            }
            interleave = []orsa.TeamRank{}
        }
    }

    return matches
}

func PrintMatchesReddit (matches []orsa.MatchGroup) {
    const MatchHeader string = "**Groups for Week %d of the %s %s league are as follows:**\n\n"
    fmt.Printf(MatchHeader, *week, *format, strings.ToUpper(*platform))
    for group, match := range matches {
        fmt.Println(match.ToStringReddit(group+1, *week, *platform, *format))
    }
}

func PrintMatchesORSA (matches []orsa.MatchGroup) {
    for group, match := range matches {
        fmt.Println(match.ToStringORSA(group+1, *week, *platform, *format))
    }
}

func PrintMatchesFlat (matches []orsa.MatchGroup) {
    out := ""
    for _, match := range matches {
        for i, t := range match.Teams {
            out += t.Name
            if i+1 != match.N {
                out += ","
            }
        }
        out += "\n"
    }
    fmt.Print(out)
}

func main() {
    flag.Parse()
    if ok := checkFlags(); !ok {
        fmt.Println("One or more required flags were missing...")
        flag.PrintDefaults()
        return
    }

    if *team1 != "" && *team2 != "" {
        //Print markup for team names
        teamRank1 := orsa.TeamRank{Name: *team1}
        teamRank2 := orsa.TeamRank{Name: *team2}
        //match group with one empty match
        fakeMatch := []orsa.MatchGroup{{}}
        fakeMatch[0].Teams[0] = teamRank1
        fakeMatch[0].Teams[1] = teamRank2
        fakeMatch[0].N = 2

        fmt.Println("<----------REDDIT FORMAT---------->\n")
        PrintMatchesReddit(fakeMatch)
        fmt.Println("<----------ORSA FORMAT---------->\n")
        PrintMatchesORSA(fakeMatch)
        fmt.Println("\n\n<----------PLAIN GROUPS FORMAT---------->\n")
        PrintMatchesFlat(fakeMatch)
     
        return
    }

    rankFile, err := os.Open(*rankFileName)
    if err != nil {
        panic(err)
    }
    defer rankFile.Close()

    ranks := &orsa.RankedList{}
    heap.Init(ranks)

    unsortedRanks := []*orsa.TeamRank{}
    if err := gocsv.UnmarshalFile(rankFile, &unsortedRanks); err != nil {
        panic(err)
    }

    for _, r := range unsortedRanks {
        heap.Push(ranks, *r)
    }

    groups := readGroups(groupsFileName)
    matches := []orsa.MatchGroup{}

    maxAttempts := 11
    attemptNum := 1
    leftoverTeams := []orsa.TeamRank{}
    matches = MakeMatches(matches, []orsa.TeamRank{}, ranks, groups)
    for ranks.Len() > 1 && attemptNum < maxAttempts {
        //move leftover teams into an interleave group
        for ranks.Len() > 0 {
            leftoverTeams = append(leftoverTeams, heap.Pop(ranks).(orsa.TeamRank))
        }

        //disband a number of teams
        for i := 0; i < attemptNum; i++ {
            if len(matches) > 0 {
                DisbandMatch(matches[len(matches)-1], ranks)
                matches = matches[:len(matches)-1]
            }
        }
        fmt.Printf("%d leftover, %d pool\n", len(leftoverTeams), ranks.Len())
        matches = MakeMatches(matches, leftoverTeams, ranks, groups)

        leftoverTeams = []orsa.TeamRank{}
        attemptNum++
    }

    if ranks.Len() == 1 {
        //Add the odd man out to the first available match
        PlaceExtraTeam(&matches, ranks, groups)
    } else if ranks.Len() > 1 {
        fmt.Println("Attempts ran out, and we weren't able to group everyone :(")
    }

    fmt.Println("<----------REDDIT FORMAT---------->\n")
    PrintMatchesReddit(matches)
    fmt.Println("<----------ORSA FORMAT---------->\n")
    PrintMatchesORSA(matches)
    fmt.Println("\n\n<----------PLAIN GROUPS FORMAT---------->\n")
    PrintMatchesFlat(matches)
    fmt.Println("\n\n<----------RESULTS---------->\n")
    fmt.Printf("%d matches, %d teams left over\n\n", len(matches), ranks.Len())
}
