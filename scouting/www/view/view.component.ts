import {Component, OnInit} from '@angular/core';
import {
  Ranking,
  RequestAllDriverRankingsResponse,
} from '../../webserver/requests/messages/request_all_driver_rankings_response_generated';
import {
  Stats,
  RequestDataScoutingResponse,
} from '../../webserver/requests/messages/request_data_scouting_response_generated';
import {
  Note,
  RequestAllNotesResponse,
} from '../../webserver/requests/messages/request_all_notes_response_generated';

import {ViewDataRequestor} from '../rpc';

type Source = 'Notes' | 'Stats' | 'DriverRanking';

//TODO(Filip): Deduplicate
const COMP_LEVEL_LABELS = {
  qm: 'Qualifications',
  ef: 'Eighth Finals',
  qf: 'Quarter Finals',
  sf: 'Semi Finals',
  f: 'Finals',
};

@Component({
  selector: 'app-view',
  templateUrl: './view.ng.html',
  styleUrls: ['../app/common.css', './view.component.css'],
})
export class ViewComponent {
  constructor(private readonly viewDataRequestor: ViewDataRequestor) {}

  // Make COMP_LEVEL_LABELS available in view.ng.html.
  readonly COMP_LEVEL_LABELS = COMP_LEVEL_LABELS;

  // Progress and error messages to display to
  // the user when fetching data.
  progressMessage: string = '';
  errorMessage: string = '';

  // The current data source being displayed.
  currentSource: Source = 'Notes';

  // Current sort (ascending/descending match numbers).
  // noteList is sorted based on team number until match
  // number is added for note scouting.
  ascendingSort = true;

  // Stores the corresponding data.
  noteList: Note[] = [];
  driverRankingList: Ranking[] = [];
  statList: Stats[] = [];

  // Fetch notes on initialization.
  ngOnInit() {
    this.fetchCurrentSource();
  }

  // Called when a user changes the sort direction.
  // Changes the data order between ascending/descending.
  sortData() {
    this.ascendingSort = !this.ascendingSort;
    if (!this.ascendingSort) {
      this.driverRankingList.sort((a, b) => b.matchNumber() - a.matchNumber());
      this.noteList.sort((a, b) => b.team() - a.team());
      this.statList.sort((a, b) => b.match() - a.match());
    } else {
      this.driverRankingList.sort((a, b) => a.matchNumber() - b.matchNumber());
      this.noteList.sort((a, b) => a.team() - b.team());
      this.statList.sort((a, b) => a.match() - b.match());
    }
  }

  // Called when a user selects a new data source
  // from the dropdown.
  switchDataSource(target: Source) {
    this.currentSource = target;
    this.progressMessage = '';
    this.errorMessage = '';
    this.noteList = [];
    this.driverRankingList = [];
    this.statList = [];
    this.fetchCurrentSource();
  }

  // Call the method to fetch data for the current source.
  fetchCurrentSource() {
    switch (this.currentSource) {
      case 'Notes': {
        this.fetchNotes();
      }

      case 'Stats': {
        this.fetchStats();
      }

      case 'DriverRanking': {
        this.fetchDriverRanking();
      }
    }
  }

  // TODO(Filip): Add delete functionality.
  // Gets called when a user clicks the delete icon.
  async deleteData() {
    const block_alerts = document.getElementById(
      'block_alerts'
    ) as HTMLInputElement;
    if (!block_alerts.checked) {
      if (!window.confirm('Actually delete data?')) {
        this.errorMessage = 'Deleting data has not been implemented yet.';
        return;
      }
    }
  }

  // Fetch all driver ranking data and store in driverRankingList.
  async fetchDriverRanking() {
    this.progressMessage = 'Fetching driver ranking data. Please be patient.';
    this.errorMessage = '';

    try {
      this.driverRankingList =
        await this.viewDataRequestor.fetchDriverRankingList();
      this.progressMessage = 'Successfully fetched driver ranking data.';
    } catch (e) {
      this.errorMessage = e;
      this.progressMessage = '';
    }
  }

  // Fetch all data scouting (stats) data and store in statList.
  async fetchStats() {
    this.progressMessage = 'Fetching stats list. Please be patient.';
    this.errorMessage = '';

    try {
      this.statList = await this.viewDataRequestor.fetchStatsList();
      this.progressMessage = 'Successfully fetched stats list.';
    } catch (e) {
      this.errorMessage = e;
      this.progressMessage = '';
    }
  }

  // Fetch all notes data and store in noteList.
  async fetchNotes() {
    this.progressMessage = 'Fetching notes list. Please be patient.';
    this.errorMessage = '';

    try {
      this.noteList = await this.viewDataRequestor.fetchNoteList();
      this.progressMessage = 'Successfully fetched note list.';
    } catch (e) {
      this.errorMessage = e;
      this.progressMessage = '';
    }
  }

  // Parse all selected keywords for a note entry
  // into one string to be displayed in the table.
  parseKeywords(entry: Note) {
    let parsedKeywords = '';

    if (entry.goodDriving()) {
      parsedKeywords += 'Good Driving ';
    }
    if (entry.badDriving()) {
      parsedKeywords += 'Bad Driving ';
    }
    if (entry.solidClimb()) {
      parsedKeywords += 'Solid Climb ';
    }
    if (entry.sketchyClimb()) {
      parsedKeywords += 'Sketchy Climb ';
    }
    if (entry.goodDefense()) {
      parsedKeywords += 'Good Defense ';
    }
    if (entry.badDefense()) {
      parsedKeywords += 'Bad Defense ';
    }

    return parsedKeywords;
  }
}
