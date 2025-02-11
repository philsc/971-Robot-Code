package requests

import (
	"bytes"
	"io"
	"net/http"
	"reflect"
	"testing"

	"github.com/frc971/971-Robot-Code/scouting/db"
	"github.com/frc971/971-Robot-Code/scouting/scraping"
	"github.com/frc971/971-Robot-Code/scouting/webserver/requests/debug"
	"github.com/frc971/971-Robot-Code/scouting/webserver/requests/messages/error_response"
	"github.com/frc971/971-Robot-Code/scouting/webserver/requests/messages/refresh_match_list"
	"github.com/frc971/971-Robot-Code/scouting/webserver/requests/messages/refresh_match_list_response"
	"github.com/frc971/971-Robot-Code/scouting/webserver/requests/messages/request_all_driver_rankings"
	"github.com/frc971/971-Robot-Code/scouting/webserver/requests/messages/request_all_driver_rankings_response"
	"github.com/frc971/971-Robot-Code/scouting/webserver/requests/messages/request_all_matches"
	"github.com/frc971/971-Robot-Code/scouting/webserver/requests/messages/request_all_matches_response"
	"github.com/frc971/971-Robot-Code/scouting/webserver/requests/messages/request_all_notes"
	"github.com/frc971/971-Robot-Code/scouting/webserver/requests/messages/request_all_notes_response"
	"github.com/frc971/971-Robot-Code/scouting/webserver/requests/messages/request_data_scouting"
	"github.com/frc971/971-Robot-Code/scouting/webserver/requests/messages/request_data_scouting_response"
	"github.com/frc971/971-Robot-Code/scouting/webserver/requests/messages/request_notes_for_team"
	"github.com/frc971/971-Robot-Code/scouting/webserver/requests/messages/request_shift_schedule"
	"github.com/frc971/971-Robot-Code/scouting/webserver/requests/messages/request_shift_schedule_response"
	"github.com/frc971/971-Robot-Code/scouting/webserver/requests/messages/submit_data_scouting"
	"github.com/frc971/971-Robot-Code/scouting/webserver/requests/messages/submit_data_scouting_response"
	"github.com/frc971/971-Robot-Code/scouting/webserver/requests/messages/submit_driver_ranking"
	"github.com/frc971/971-Robot-Code/scouting/webserver/requests/messages/submit_notes"
	"github.com/frc971/971-Robot-Code/scouting/webserver/requests/messages/submit_shift_schedule"
	"github.com/frc971/971-Robot-Code/scouting/webserver/server"
	flatbuffers "github.com/google/flatbuffers/go"
)

// Validates that an unhandled address results in a 404.
func Test404(t *testing.T) {
	db := MockDatabase{}
	scoutingServer := server.NewScoutingServer()
	HandleRequests(&db, scrapeEmtpyMatchList, scoutingServer)
	scoutingServer.Start(8080)
	defer scoutingServer.Stop()

	resp, err := http.Get("http://localhost:8080/requests/foo")
	if err != nil {
		t.Fatalf("Failed to get data: %v", err)
	}
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("Expected error code 404, but got %d instead", resp.Status)
	}
}

// Validates that we can submit new data scouting data.
func TestSubmitDataScoutingError(t *testing.T) {
	db := MockDatabase{}
	scoutingServer := server.NewScoutingServer()
	HandleRequests(&db, scrapeEmtpyMatchList, scoutingServer)
	scoutingServer.Start(8080)
	defer scoutingServer.Stop()

	resp, err := http.Post("http://localhost:8080/requests/submit/data_scouting", "application/octet-stream", bytes.NewReader([]byte("")))
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatal("Unexpected status code. Got", resp.Status)
	}

	responseBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal("Failed to read response bytes:", err)
	}
	errorResponse := error_response.GetRootAsErrorResponse(responseBytes, 0)

	errorMessage := string(errorResponse.ErrorMessage())
	if errorMessage != "Failed to parse SubmitDataScouting: runtime error: index out of range [3] with length 0" {
		t.Fatal("Got mismatched error message:", errorMessage)
	}
}

// Validates that we can submit new data scouting data.
func TestSubmitDataScouting(t *testing.T) {
	db := MockDatabase{}
	scoutingServer := server.NewScoutingServer()
	HandleRequests(&db, scrapeEmtpyMatchList, scoutingServer)
	scoutingServer.Start(8080)
	defer scoutingServer.Stop()

	builder := flatbuffers.NewBuilder(1024)
	builder.Finish((&submit_data_scouting.SubmitDataScoutingT{
		Team:                  971,
		Match:                 1,
		SetNumber:             8,
		CompLevel:             "quals",
		StartingQuadrant:      2,
		AutoBall1:             true,
		AutoBall2:             false,
		AutoBall3:             false,
		AutoBall4:             false,
		AutoBall5:             false,
		MissedShotsAuto:       9971,
		UpperGoalAuto:         9971,
		LowerGoalAuto:         9971,
		MissedShotsTele:       9971,
		UpperGoalTele:         9971,
		LowerGoalTele:         9971,
		DefenseRating:         9971,
		DefenseReceivedRating: 4,
		ClimbLevel:            submit_data_scouting.ClimbLevelLow,
		Comment:               "this is a comment",
	}).Pack(builder))

	response, err := debug.SubmitDataScouting("http://localhost:8080", builder.FinishedBytes())
	if err != nil {
		t.Fatal("Failed to submit data scouting: ", err)
	}

	// We get an empty response back. Validate that.
	expected := submit_data_scouting_response.SubmitDataScoutingResponseT{}
	if !reflect.DeepEqual(expected, *response) {
		t.Fatal("Expected ", expected, ", but got:", *response)
	}
}

// Validates that we can request the full match list.
func TestRequestAllMatches(t *testing.T) {
	db := MockDatabase{
		matches: []db.TeamMatch{
			{
				MatchNumber: 1, SetNumber: 1, CompLevel: "qm",
				Alliance: "R", AlliancePosition: 1, TeamNumber: 5,
			},
			{
				MatchNumber: 1, SetNumber: 1, CompLevel: "qm",
				Alliance: "R", AlliancePosition: 2, TeamNumber: 42,
			},
			{
				MatchNumber: 1, SetNumber: 1, CompLevel: "qm",
				Alliance: "R", AlliancePosition: 3, TeamNumber: 600,
			},
			{
				MatchNumber: 1, SetNumber: 1, CompLevel: "qm",
				Alliance: "B", AlliancePosition: 1, TeamNumber: 971,
			},
			{
				MatchNumber: 1, SetNumber: 1, CompLevel: "qm",
				Alliance: "B", AlliancePosition: 2, TeamNumber: 400,
			},
			{
				MatchNumber: 1, SetNumber: 1, CompLevel: "qm",
				Alliance: "B", AlliancePosition: 3, TeamNumber: 200,
			},
			{
				MatchNumber: 2, SetNumber: 1, CompLevel: "qm",
				Alliance: "R", AlliancePosition: 1, TeamNumber: 6,
			},
			{
				MatchNumber: 2, SetNumber: 1, CompLevel: "qm",
				Alliance: "R", AlliancePosition: 2, TeamNumber: 43,
			},
			{
				MatchNumber: 2, SetNumber: 1, CompLevel: "qm",
				Alliance: "R", AlliancePosition: 3, TeamNumber: 601,
			},
			{
				MatchNumber: 2, SetNumber: 1, CompLevel: "qm",
				Alliance: "B", AlliancePosition: 1, TeamNumber: 972,
			},
			{
				MatchNumber: 2, SetNumber: 1, CompLevel: "qm",
				Alliance: "B", AlliancePosition: 2, TeamNumber: 401,
			},
			{
				MatchNumber: 2, SetNumber: 1, CompLevel: "qm",
				Alliance: "B", AlliancePosition: 3, TeamNumber: 201,
			},
			{
				MatchNumber: 3, SetNumber: 1, CompLevel: "qm",
				Alliance: "R", AlliancePosition: 1, TeamNumber: 7,
			},
			{
				MatchNumber: 3, SetNumber: 1, CompLevel: "qm",
				Alliance: "R", AlliancePosition: 2, TeamNumber: 44,
			},
			{
				MatchNumber: 3, SetNumber: 1, CompLevel: "qm",
				Alliance: "R", AlliancePosition: 3, TeamNumber: 602,
			},
			{
				MatchNumber: 3, SetNumber: 1, CompLevel: "qm",
				Alliance: "B", AlliancePosition: 1, TeamNumber: 973,
			},
			{
				MatchNumber: 3, SetNumber: 1, CompLevel: "qm",
				Alliance: "B", AlliancePosition: 2, TeamNumber: 402,
			},
			{
				MatchNumber: 3, SetNumber: 1, CompLevel: "qm",
				Alliance: "B", AlliancePosition: 3, TeamNumber: 202,
			},
		},
	}
	scoutingServer := server.NewScoutingServer()
	HandleRequests(&db, scrapeEmtpyMatchList, scoutingServer)
	scoutingServer.Start(8080)
	defer scoutingServer.Stop()

	builder := flatbuffers.NewBuilder(1024)
	builder.Finish((&request_all_matches.RequestAllMatchesT{}).Pack(builder))

	response, err := debug.RequestAllMatches("http://localhost:8080", builder.FinishedBytes())
	if err != nil {
		t.Fatal("Failed to request all matches: ", err)
	}

	expected := request_all_matches_response.RequestAllMatchesResponseT{
		MatchList: []*request_all_matches_response.MatchT{
			// MatchNumber, SetNumber, CompLevel
			// R1, R2, R3, B1, B2, B3
			{
				1, 1, "qm",
				5, 42, 600, 971, 400, 200,
			},
			{
				2, 1, "qm",
				6, 43, 601, 972, 401, 201,
			},
			{
				3, 1, "qm",
				7, 44, 602, 973, 402, 202,
			},
		},
	}
	if len(expected.MatchList) != len(response.MatchList) {
		t.Fatal("Expected ", expected, ", but got ", *response)
	}
	for i, match := range expected.MatchList {
		if !reflect.DeepEqual(*match, *response.MatchList[i]) {
			t.Fatal("Expected for match", i, ":", *match, ", but got:", *response.MatchList[i])
		}
	}

}

// Validates that we can request the stats.
func TestRequestDataScouting(t *testing.T) {
	db := MockDatabase{
		stats: []db.Stats{
			{
				TeamNumber: 971, MatchNumber: 1, SetNumber: 2, CompLevel: "quals",
				StartingQuadrant: 1,
				AutoBallPickedUp: [5]bool{true, false, false, false, true},
				ShotsMissed:      1, UpperGoalShots: 2, LowerGoalShots: 3,
				ShotsMissedAuto: 4, UpperGoalAuto: 5, LowerGoalAuto: 6,
				PlayedDefense: 7, DefenseReceivedScore: 3, Climbing: 2,
				Comment: "a lovely comment", CollectedBy: "john",
			},
			{
				TeamNumber: 972, MatchNumber: 1, SetNumber: 4, CompLevel: "extra",
				StartingQuadrant: 2,
				AutoBallPickedUp: [5]bool{false, false, true, false, false},
				ShotsMissed:      2, UpperGoalShots: 3, LowerGoalShots: 4,
				ShotsMissedAuto: 5, UpperGoalAuto: 6, LowerGoalAuto: 7,
				PlayedDefense: 8, DefenseReceivedScore: 1, Climbing: 4,
				Comment: "another lovely comment", CollectedBy: "andrea",
			},
		},
	}
	scoutingServer := server.NewScoutingServer()
	HandleRequests(&db, scrapeEmtpyMatchList, scoutingServer)
	scoutingServer.Start(8080)
	defer scoutingServer.Stop()

	builder := flatbuffers.NewBuilder(1024)
	builder.Finish((&request_data_scouting.RequestDataScoutingT{}).Pack(builder))

	response, err := debug.RequestDataScouting("http://localhost:8080", builder.FinishedBytes())
	if err != nil {
		t.Fatal("Failed to request all matches: ", err)
	}

	expected := request_data_scouting_response.RequestDataScoutingResponseT{
		StatsList: []*request_data_scouting_response.StatsT{
			{
				Team: 971, Match: 1, SetNumber: 2, CompLevel: "quals",
				MissedShotsAuto: 4, UpperGoalAuto: 5, LowerGoalAuto: 6,
				MissedShotsTele: 1, UpperGoalTele: 2, LowerGoalTele: 3,
				DefenseRating:         7,
				DefenseReceivedRating: 3,
				CollectedBy:           "john",
				AutoBall1:             true, AutoBall2: false, AutoBall3: false,
				AutoBall4: false, AutoBall5: true,
				StartingQuadrant: 1,
				ClimbLevel:       request_data_scouting_response.ClimbLevelFailedWithPlentyOfTime,
				Comment:          "a lovely comment",
			},
			{
				Team: 972, Match: 1, SetNumber: 4, CompLevel: "extra",
				MissedShotsAuto: 5, UpperGoalAuto: 6, LowerGoalAuto: 7,
				MissedShotsTele: 2, UpperGoalTele: 3, LowerGoalTele: 4,
				DefenseRating:         8,
				DefenseReceivedRating: 1,
				CollectedBy:           "andrea",
				AutoBall1:             false, AutoBall2: false, AutoBall3: true,
				AutoBall4: false, AutoBall5: false,
				StartingQuadrant: 2,
				ClimbLevel:       request_data_scouting_response.ClimbLevelMedium,
				Comment:          "another lovely comment",
			},
		},
	}
	if len(expected.StatsList) != len(response.StatsList) {
		t.Fatal("Expected ", expected, ", but got ", *response)
	}
	for i, match := range expected.StatsList {
		if !reflect.DeepEqual(*match, *response.StatsList[i]) {
			t.Fatal("Expected for stats", i, ":", *match, ", but got:", *response.StatsList[i])
		}
	}
}

func TestSubmitNotes(t *testing.T) {
	database := MockDatabase{}
	scoutingServer := server.NewScoutingServer()
	HandleRequests(&database, scrapeEmtpyMatchList, scoutingServer)
	scoutingServer.Start(8080)
	defer scoutingServer.Stop()

	builder := flatbuffers.NewBuilder(1024)
	builder.Finish((&submit_notes.SubmitNotesT{
		Team:         971,
		Notes:        "Notes",
		GoodDriving:  true,
		BadDriving:   false,
		SketchyClimb: true,
		SolidClimb:   false,
		GoodDefense:  true,
		BadDefense:   false,
	}).Pack(builder))

	_, err := debug.SubmitNotes("http://localhost:8080", builder.FinishedBytes())
	if err != nil {
		t.Fatal("Failed to submit notes: ", err)
	}

	expected := []db.NotesData{
		{
			TeamNumber:   971,
			Notes:        "Notes",
			GoodDriving:  true,
			BadDriving:   false,
			SketchyClimb: true,
			SolidClimb:   false,
			GoodDefense:  true,
			BadDefense:   false,
		},
	}

	if !reflect.DeepEqual(database.notes, expected) {
		t.Fatal("Submitted notes did not match", expected, database.notes)
	}
}

func TestRequestNotes(t *testing.T) {
	database := MockDatabase{
		notes: []db.NotesData{{
			TeamNumber:   971,
			Notes:        "Notes",
			GoodDriving:  true,
			BadDriving:   false,
			SketchyClimb: true,
			SolidClimb:   false,
			GoodDefense:  true,
			BadDefense:   false,
		}},
	}
	scoutingServer := server.NewScoutingServer()
	HandleRequests(&database, scrapeEmtpyMatchList, scoutingServer)
	scoutingServer.Start(8080)
	defer scoutingServer.Stop()

	builder := flatbuffers.NewBuilder(1024)
	builder.Finish((&request_notes_for_team.RequestNotesForTeamT{
		Team: 971,
	}).Pack(builder))
	response, err := debug.RequestNotes("http://localhost:8080", builder.FinishedBytes())
	if err != nil {
		t.Fatal("Failed to submit notes: ", err)
	}

	if response.Notes[0].Data != "Notes" {
		t.Fatal("requested notes did not match", response)
	}
}

func TestRequestShiftSchedule(t *testing.T) {
	db := MockDatabase{
		shiftSchedule: []db.Shift{
			{
				MatchNumber: 1,
				R1scouter:   "Bob",
				R2scouter:   "James",
				R3scouter:   "Robert",
				B1scouter:   "Alice",
				B2scouter:   "Mary",
				B3scouter:   "Patricia",
			},
			{
				MatchNumber: 2,
				R1scouter:   "Liam",
				R2scouter:   "Noah",
				R3scouter:   "Oliver",
				B1scouter:   "Emma",
				B2scouter:   "Charlotte",
				B3scouter:   "Amelia",
			},
		},
	}
	scoutingServer := server.NewScoutingServer()
	HandleRequests(&db, scrapeEmtpyMatchList, scoutingServer)
	scoutingServer.Start(8080)
	defer scoutingServer.Stop()

	builder := flatbuffers.NewBuilder(1024)
	builder.Finish((&request_shift_schedule.RequestShiftScheduleT{}).Pack(builder))

	response, err := debug.RequestShiftSchedule("http://localhost:8080", builder.FinishedBytes())
	if err != nil {
		t.Fatal("Failed to request shift schedule: ", err)
	}

	expected := request_shift_schedule_response.RequestShiftScheduleResponseT{
		ShiftSchedule: []*request_shift_schedule_response.MatchAssignmentT{
			{
				MatchNumber: 1,
				R1scouter:   "Bob",
				R2scouter:   "James",
				R3scouter:   "Robert",
				B1scouter:   "Alice",
				B2scouter:   "Mary",
				B3scouter:   "Patricia",
			},
			{
				MatchNumber: 2,
				R1scouter:   "Liam",
				R2scouter:   "Noah",
				R3scouter:   "Oliver",
				B1scouter:   "Emma",
				B2scouter:   "Charlotte",
				B3scouter:   "Amelia",
			},
		},
	}
	if len(expected.ShiftSchedule) != len(response.ShiftSchedule) {
		t.Fatal("Expected ", expected, ", but got ", *response)
	}
	for i, match := range expected.ShiftSchedule {
		if !reflect.DeepEqual(*match, *response.ShiftSchedule[i]) {
			t.Fatal("Expected for shift schedule", i, ":", *match, ", but got:", *response.ShiftSchedule[i])
		}
	}
}

func TestSubmitShiftSchedule(t *testing.T) {
	database := MockDatabase{}
	scoutingServer := server.NewScoutingServer()
	HandleRequests(&database, scrapeEmtpyMatchList, scoutingServer)
	scoutingServer.Start(8080)
	defer scoutingServer.Stop()

	builder := flatbuffers.NewBuilder(1024)
	builder.Finish((&submit_shift_schedule.SubmitShiftScheduleT{
		ShiftSchedule: []*submit_shift_schedule.MatchAssignmentT{
			{MatchNumber: 1,
				R1scouter: "Bob",
				R2scouter: "James",
				R3scouter: "Robert",
				B1scouter: "Alice",
				B2scouter: "Mary",
				B3scouter: "Patricia"},
		},
	}).Pack(builder))

	_, err := debug.SubmitShiftSchedule("http://localhost:8080", builder.FinishedBytes())
	if err != nil {
		t.Fatal("Failed to submit shift schedule: ", err)
	}

	expected := []db.Shift{
		{MatchNumber: 1,
			R1scouter: "Bob",
			R2scouter: "James",
			R3scouter: "Robert",
			B1scouter: "Alice",
			B2scouter: "Mary",
			B3scouter: "Patricia"},
	}
	if !reflect.DeepEqual(expected, database.shiftSchedule) {
		t.Fatal("Expected ", expected, ", but got:", database.shiftSchedule)
	}
}

// Validates that we can download the schedule from The Blue Alliance.
func TestRefreshMatchList(t *testing.T) {
	scrapeMockSchedule := func(int32, string) ([]scraping.Match, error) {
		return []scraping.Match{
			{
				CompLevel:   "qual",
				MatchNumber: 1,
				SetNumber:   2,
				Alliances: scraping.Alliances{
					Red: scraping.Alliance{
						TeamKeys: []string{
							"100",
							"200",
							"300",
						},
					},
					Blue: scraping.Alliance{
						TeamKeys: []string{
							"101",
							"201",
							"301",
						},
					},
				},
				WinningAlliance: "",
				EventKey:        "",
				Time:            0,
				PredictedTime:   0,
				ActualTime:      0,
				PostResultTime:  0,
				ScoreBreakdowns: scraping.ScoreBreakdowns{},
			},
		}, nil
	}

	database := MockDatabase{}
	scoutingServer := server.NewScoutingServer()
	HandleRequests(&database, scrapeMockSchedule, scoutingServer)
	scoutingServer.Start(8080)
	defer scoutingServer.Stop()

	builder := flatbuffers.NewBuilder(1024)
	builder.Finish((&refresh_match_list.RefreshMatchListT{}).Pack(builder))

	response, err := debug.RefreshMatchList("http://localhost:8080", builder.FinishedBytes())
	if err != nil {
		t.Fatal("Failed to request all matches: ", err)
	}

	// Validate the response.
	expected := refresh_match_list_response.RefreshMatchListResponseT{}
	if !reflect.DeepEqual(expected, *response) {
		t.Fatal("Expected ", expected, ", but got ", *response)
	}

	// Make sure that the data made it into the database.
	expectedMatches := []db.TeamMatch{
		{
			MatchNumber: 1, SetNumber: 2, CompLevel: "qual",
			Alliance: "R", AlliancePosition: 1, TeamNumber: 100,
		},
		{
			MatchNumber: 1, SetNumber: 2, CompLevel: "qual",
			Alliance: "R", AlliancePosition: 2, TeamNumber: 200,
		},
		{
			MatchNumber: 1, SetNumber: 2, CompLevel: "qual",
			Alliance: "R", AlliancePosition: 3, TeamNumber: 300,
		},
		{
			MatchNumber: 1, SetNumber: 2, CompLevel: "qual",
			Alliance: "B", AlliancePosition: 1, TeamNumber: 101,
		},
		{
			MatchNumber: 1, SetNumber: 2, CompLevel: "qual",
			Alliance: "B", AlliancePosition: 2, TeamNumber: 201,
		},
		{
			MatchNumber: 1, SetNumber: 2, CompLevel: "qual",
			Alliance: "B", AlliancePosition: 3, TeamNumber: 301,
		},
	}

	if !reflect.DeepEqual(expectedMatches, database.matches) {
		t.Fatal("Expected ", expectedMatches, ", but got ", database.matches)
	}
}

func TestSubmitDriverRanking(t *testing.T) {
	database := MockDatabase{}
	scoutingServer := server.NewScoutingServer()
	HandleRequests(&database, scrapeEmtpyMatchList, scoutingServer)
	scoutingServer.Start(8080)
	defer scoutingServer.Stop()

	builder := flatbuffers.NewBuilder(1024)
	builder.Finish((&submit_driver_ranking.SubmitDriverRankingT{
		MatchNumber: 36,
		Rank1:       1234,
		Rank2:       1235,
		Rank3:       1236,
	}).Pack(builder))

	_, err := debug.SubmitDriverRanking("http://localhost:8080", builder.FinishedBytes())
	if err != nil {
		t.Fatal("Failed to submit driver ranking: ", err)
	}

	expected := []db.DriverRankingData{
		{MatchNumber: 36, Rank1: 1234, Rank2: 1235, Rank3: 1236},
	}

	if !reflect.DeepEqual(database.driver_ranking, expected) {
		t.Fatal("Submitted notes did not match", expected, database.notes)
	}
}

// Validates that we can request the driver rankings.
func TestRequestDriverRankings(t *testing.T) {
	db := MockDatabase{
		driver_ranking: []db.DriverRankingData{
			{
				MatchNumber: 36,
				Rank1:       1234,
				Rank2:       1235,
				Rank3:       1236,
			},
			{
				MatchNumber: 36,
				Rank1:       101,
				Rank2:       202,
				Rank3:       303,
			},
		},
	}
	scoutingServer := server.NewScoutingServer()
	HandleRequests(&db, scrapeEmtpyMatchList, scoutingServer)
	scoutingServer.Start(8080)
	defer scoutingServer.Stop()

	builder := flatbuffers.NewBuilder(1024)
	builder.Finish((&request_all_driver_rankings.RequestAllDriverRankingsT{}).Pack(builder))

	response, err := debug.RequestAllDriverRankings("http://localhost:8080", builder.FinishedBytes())
	if err != nil {
		t.Fatal("Failed to request all driver rankings: ", err)
	}

	expected := request_all_driver_rankings_response.RequestAllDriverRankingsResponseT{
		DriverRankingList: []*request_all_driver_rankings_response.RankingT{
			{
				MatchNumber: 36,
				Rank1:       1234,
				Rank2:       1235,
				Rank3:       1236,
			},
			{
				MatchNumber: 36,
				Rank1:       101,
				Rank2:       202,
				Rank3:       303,
			},
		},
	}
	if len(expected.DriverRankingList) != len(response.DriverRankingList) {
		t.Fatal("Expected ", expected, ", but got ", *response)
	}
	for i, match := range expected.DriverRankingList {
		if !reflect.DeepEqual(*match, *response.DriverRankingList[i]) {
			t.Fatal("Expected for driver ranking", i, ":", *match, ", but got:", *response.DriverRankingList[i])
		}
	}
}

// Validates that we can request all notes.
func TestRequestAllNotes(t *testing.T) {
	db := MockDatabase{
		notes: []db.NotesData{
			{
				TeamNumber:   971,
				Notes:        "Notes",
				GoodDriving:  true,
				BadDriving:   false,
				SketchyClimb: true,
				SolidClimb:   false,
				GoodDefense:  true,
				BadDefense:   false,
			},
			{
				TeamNumber:   972,
				Notes:        "More Notes",
				GoodDriving:  false,
				BadDriving:   false,
				SketchyClimb: false,
				SolidClimb:   true,
				GoodDefense:  false,
				BadDefense:   true,
			},
		},
	}
	scoutingServer := server.NewScoutingServer()
	HandleRequests(&db, scrapeEmtpyMatchList, scoutingServer)
	scoutingServer.Start(8080)
	defer scoutingServer.Stop()

	builder := flatbuffers.NewBuilder(1024)
	builder.Finish((&request_all_notes.RequestAllNotesT{}).Pack(builder))

	response, err := debug.RequestAllNotes("http://localhost:8080", builder.FinishedBytes())
	if err != nil {
		t.Fatal("Failed to request all notes: ", err)
	}

	expected := request_all_notes_response.RequestAllNotesResponseT{
		NoteList: []*request_all_notes_response.NoteT{
			{
				Team:         971,
				Notes:        "Notes",
				GoodDriving:  true,
				BadDriving:   false,
				SketchyClimb: true,
				SolidClimb:   false,
				GoodDefense:  true,
				BadDefense:   false,
			},
			{
				Team:         972,
				Notes:        "More Notes",
				GoodDriving:  false,
				BadDriving:   false,
				SketchyClimb: false,
				SolidClimb:   true,
				GoodDefense:  false,
				BadDefense:   true,
			},
		},
	}
	if len(expected.NoteList) != len(response.NoteList) {
		t.Fatal("Expected ", expected, ", but got ", *response)
	}
	for i, note := range expected.NoteList {
		if !reflect.DeepEqual(*note, *response.NoteList[i]) {
			t.Fatal("Expected for note", i, ":", *note, ", but got:", *response.NoteList[i])
		}
	}
}

// A mocked database we can use for testing. Add functionality to this as
// needed for your tests.

type MockDatabase struct {
	matches        []db.TeamMatch
	stats          []db.Stats
	notes          []db.NotesData
	shiftSchedule  []db.Shift
	driver_ranking []db.DriverRankingData
}

func (database *MockDatabase) AddToMatch(match db.TeamMatch) error {
	database.matches = append(database.matches, match)
	return nil
}

func (database *MockDatabase) AddToStats(stats db.Stats) error {
	database.stats = append(database.stats, stats)
	return nil
}

func (database *MockDatabase) ReturnMatches() ([]db.TeamMatch, error) {
	return database.matches, nil
}

func (database *MockDatabase) ReturnStats() ([]db.Stats, error) {
	return database.stats, nil
}

func (database *MockDatabase) QueryStats(int) ([]db.Stats, error) {
	return []db.Stats{}, nil
}

func (database *MockDatabase) QueryNotes(requestedTeam int32) ([]string, error) {
	var results []string
	for _, data := range database.notes {
		if data.TeamNumber == requestedTeam {
			results = append(results, data.Notes)
		}
	}
	return results, nil
}

func (database *MockDatabase) AddNotes(data db.NotesData) error {
	database.notes = append(database.notes, data)
	return nil
}

func (database *MockDatabase) ReturnAllNotes() ([]db.NotesData, error) {
	return database.notes, nil
}

func (database *MockDatabase) AddToShift(data db.Shift) error {
	database.shiftSchedule = append(database.shiftSchedule, data)
	return nil
}

func (database *MockDatabase) ReturnAllShifts() ([]db.Shift, error) {
	return database.shiftSchedule, nil
}

func (database *MockDatabase) QueryAllShifts(int) ([]db.Shift, error) {
	return []db.Shift{}, nil
}

func (database *MockDatabase) AddDriverRanking(data db.DriverRankingData) error {
	database.driver_ranking = append(database.driver_ranking, data)
	return nil
}

func (database *MockDatabase) ReturnAllDriverRankings() ([]db.DriverRankingData, error) {
	return database.driver_ranking, nil
}

// Returns an empty match list from the fake The Blue Alliance scraping.
func scrapeEmtpyMatchList(int32, string) ([]scraping.Match, error) {
	return nil, nil
}
